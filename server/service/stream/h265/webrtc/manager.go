package webrtc

import (
	"NanoKVM-Server/common"
	"NanoKVM-Server/service/stream"
	"NanoKVM-Server/service/stream/opus"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/samplebuilder"

	"runtime"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4/pkg/media"
	log "github.com/sirupsen/logrus"
)

func NewWebRTCManager() *WebRTCManager {
	m := &WebRTCManager{
		tracks:       make(map[*websocket.Conn]*Track),
		audioSending: 0,
	}
	m.updateTrackSnapshotLocked()

	return m
}

func (m *WebRTCManager) AddTrack(ws *websocket.Conn, track *Track) {
	track.updateExtension()

	m.mutex.Lock()
	m.tracks[ws] = track
	count := m.updateTrackSnapshotLocked()
	m.mutex.Unlock()

	stream.GetVideoHub().SetMode(common.STREAM_TYPE_H265_WEBRTC)

	log.Debugf("added track %s, total tracks: %d", ws.RemoteAddr(), count)
}

func (m *WebRTCManager) RemoveTrack(ws *websocket.Conn) {
	m.mutex.Lock()
	if track, ok := m.tracks[ws]; ok {
		track.stopVideoStream()
	}
	delete(m.tracks, ws)
	count := m.updateTrackSnapshotLocked()
	m.mutex.Unlock()

	log.Debugf("removed track %s, total tracks: %d", ws.RemoteAddr(), count)
}

func (m *WebRTCManager) GetTrackCount() int {
	return len(m.getTracks())
}

func (m *WebRTCManager) updateTrackSnapshotLocked() int {
	tracks := make([]*Track, 0, len(m.tracks))
	for _, track := range m.tracks {
		tracks = append(tracks, track)
	}
	m.trackSnapshot.Store(&tracks)

	return len(tracks)
}

func (m *WebRTCManager) getTracks() []*Track {
	tracks := m.trackSnapshot.Load()
	if tracks == nil {
		return nil
	}

	return *tracks
}

func (m *WebRTCManager) StartAudioStream() {
	if atomic.CompareAndSwapInt32(&m.audioSending, 0, 1) {
		go m.sendAudioStream()
	}
}

func (m *WebRTCManager) StartMicStream(track *webrtc.TrackRemote) {
	go m.receiveAudio(track)
}

func (m *WebRTCManager) sendAudioStream() {
	defer atomic.StoreInt32(&m.audioSending, 0)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := stream.SetCPUAffinity(1); err != nil {
		return
	}

	vision := common.GetKvmVision()

	duration := time.Millisecond * 20
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		if m.GetTrackCount() == 0 {
			return
		}

		if !stream.GetVideoHub().IsMode(common.STREAM_TYPE_H265_WEBRTC) {
			continue
		}

		data, result := vision.ReadAudio()
		if result < 0 || len(data) == 0 {
			continue
		}

		sample := media.Sample{
			Data:     data,
			Duration: duration,
		}

		tracks := m.getTracks()
		for _, track := range tracks {
			track.writeAudio(sample)
		}
	}
}

func (m *WebRTCManager) receiveAudio(track *webrtc.TrackRemote) {
	player := opus.GetAudioInputPlayer()
	if err := player.Start(); err != nil {
		log.Errorf("failed to start audio player: %s", err)
		return
	}
	defer player.Stop()

	sampleBuilder := samplebuilder.New(20, &opus.Packet{}, 48000)

	buf := make([]byte, 1500)
	for {
		n, _, err := track.Read(buf)
		if err != nil {
			log.Debugf("audio track read error: %s", err)
			break
		}

		packet := &rtp.Packet{}
		if err := packet.Unmarshal(buf[:n]); err != nil {
			log.Debugf("failed to unmarshal RTP packet: %s", err)
			continue
		}

		sampleBuilder.Push(packet)

		for {
			sample := sampleBuilder.Pop()
			if sample == nil {
				break
			}

			if err := player.DecodeAndWrite(sample.Data); err != nil {
				log.Debugf("failed to decode and write: %s", err)
			}
		}
	}

	log.Debug("audio microphone track ended")
}
