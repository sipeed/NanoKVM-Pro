package webrtc

import (
	"NanoKVM-Server/common"
	"NanoKVM-Server/service/stream"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4/pkg/media"
	log "github.com/sirupsen/logrus"
)

func NewWebRTCManager() *WebRTCManager {
	return &WebRTCManager{
		clients:      make(map[*websocket.Conn]*Client),
		videoSending: 0,
		audioSending: 0,
		videoStatus:  0,
		mutex:        sync.RWMutex{},
	}
}

func (m *WebRTCManager) AddClient(ws *websocket.Conn, client *Client) {
	client.track.updateExtension()

	m.mutex.Lock()
	m.clients[ws] = client
	m.mutex.Unlock()

	common.GetKvmVision().SetStreamType(common.STREAM_TYPE_H264_WEBRTC)

	log.Debugf("added client %s, total clients: %d", ws.RemoteAddr(), len(m.clients))
}

func (m *WebRTCManager) RemoveClient(ws *websocket.Conn) {
	m.mutex.Lock()
	delete(m.clients, ws)
	m.mutex.Unlock()

	log.Debugf("removed client %s, total clients: %d", ws.RemoteAddr(), len(m.clients))
}

func (m *WebRTCManager) GetClientCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.clients)
}

func (m *WebRTCManager) StartVideoStream() {
	if atomic.CompareAndSwapInt32(&m.videoSending, 0, 1) {
		go m.sendVideoStream()
		log.Debugf("start sending h264 stream")
	}
}

func (m *WebRTCManager) StartAudioStream() {
	if atomic.CompareAndSwapInt32(&m.audioSending, 0, 1) {
		go m.sendAudioStream()
	}
}

func (m *WebRTCManager) sendVideoStream() {
	defer atomic.StoreInt32(&m.videoSending, 0)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := stream.SetCPUAffinity(0); err != nil {
		return
	}

	screen := common.GetScreen()
	vision := common.GetKvmVision()
	m.updateStatus(0)

	duration := time.Second / time.Duration(120)
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	startTime := time.Now()

	for range ticker.C {
		if m.GetClientCount() == 0 {
			log.Debugf("stop sending h264 stream")
			return
		}

		if vision.StreamType != common.STREAM_TYPE_H264_WEBRTC {
			m.updateStatus(-4)
			continue
		}

		data, result := vision.ReadH264(screen.Width, screen.Height, screen.BitRate)
		m.updateStatus(result)

		if result < 0 || len(data) == 0 {
			continue
		}

		timestamp := time.Since(startTime)
		startTime = time.Now()

		sample := media.Sample{
			Data:     data,
			Duration: timestamp,
		}

		for _, client := range m.clients {
			client.track.writeVideo(sample)
		}

		stream.GetFrameRateCounter().Update()
	}
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

		if vision.StreamType != common.STREAM_TYPE_H264_WEBRTC {
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

		for _, client := range m.clients {
			client.track.writeAudio(sample)
		}
	}
}

func (m *WebRTCManager) updateStatus(videoStatus int) {
	var newStatus int32
	if videoStatus >= 0 {
		newStatus = 1
	} else {
		newStatus = int32(videoStatus)
	}

	if atomic.LoadInt32(&m.videoStatus) == newStatus {
		return
	}

	data := strconv.Itoa(int(newStatus))

	m.mutex.RLock()
	clients := make([]*Client, 0, len(m.clients))
	for _, c := range m.clients {
		clients = append(clients, c)
	}
	m.mutex.RUnlock()

	for _, client := range clients {
		_ = client.WriteMessage("video-status", data)
	}

	atomic.StoreInt32(&m.videoStatus, newStatus)
}
