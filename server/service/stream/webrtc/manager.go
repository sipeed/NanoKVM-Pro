package webrtc

import (
	"runtime"
	"sync/atomic"
	"time"

	"NanoKVM-Server/common"
	"NanoKVM-Server/service/stream"
	"NanoKVM-Server/service/stream/opus"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/samplebuilder"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4/pkg/media"
	log "github.com/sirupsen/logrus"
)

func NewWebRTCManager() *WebRTCManager {
	m := &WebRTCManager{
		clients:      make(map[*websocket.Conn]*Client),
		audioSending: 0,
	}
	m.updateClientSnapshotLocked()

	return m
}

func (m *WebRTCManager) AddClient(ws *websocket.Conn, client *Client) {
	client.track.updateExtension()

	m.mutex.Lock()
	m.clients[ws] = client
	count := m.updateClientSnapshotLocked()
	m.mutex.Unlock()

	stream.GetVideoHub().SetMode(common.STREAM_TYPE_H264_WEBRTC)

	log.Debugf("added client %s, total clients: %d", ws.RemoteAddr(), count)
}

func (m *WebRTCManager) RemoveClient(ws *websocket.Conn) {
	m.mutex.Lock()
	if client, ok := m.clients[ws]; ok {
		client.stopVideoStream()
	}
	delete(m.clients, ws)
	count := m.updateClientSnapshotLocked()
	m.mutex.Unlock()

	log.Debugf("removed client %s, total clients: %d", ws.RemoteAddr(), count)
}

func (m *WebRTCManager) GetClientCount() int {
	return len(m.getClients())
}

func (m *WebRTCManager) updateClientSnapshotLocked() int {
	clients := make([]*Client, 0, len(m.clients))
	for _, client := range m.clients {
		clients = append(clients, client)
	}
	m.clientSnapshot.Store(&clients)

	return len(clients)
}

func (m *WebRTCManager) getClients() []*Client {
	clients := m.clientSnapshot.Load()
	if clients == nil {
		return nil
	}

	return *clients
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
		if m.GetClientCount() == 0 {
			return
		}

		if !stream.GetVideoHub().IsMode(common.STREAM_TYPE_H264_WEBRTC) {
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

		clients := m.getClients()
		for _, client := range clients {
			client.track.writeAudio(sample)
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

	log.Debugf("audio microphone track ended")
}
