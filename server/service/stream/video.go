package stream

import (
	"sync"
	"sync/atomic"
	"time"

	"NanoKVM-Server/common"
)

type videoSource interface {
	ReadMjpeg(uint16, uint16, uint16) ([]byte, int)
	ReadH264(uint16, uint16, uint16) ([]byte, int)
	ReadH265(uint16, uint16, uint16) ([]byte, int)
	CloseVideoStream()
	RequestIDR() int
	SetBitrate(uint16) int
	SetMjpegQuality(uint16) int
	SetRateControl(uint8) int
	SetFps(uint8) int
	SetGop(uint8) int
}

// VideoFrame owns its bytes. Timestamp is a monotonic, process-local
// microsecond clock. The legacy C API does not expose the hardware PTS.
type VideoFrame struct {
	Data       []byte
	Result     int
	KeyFrame   bool
	Duration   time.Duration
	Timestamp  int64
	Generation uint64
}

type VideoSubscription struct {
	Frames  <-chan VideoFrame
	Updates <-chan int
	Done    <-chan struct{}
	frames  chan VideoFrame
	updates chan int
	done    chan struct{}
	hub     *VideoHub
	mode    int
	waitKey bool // protected by hub.mutex
	status  int
}

func (s *VideoSubscription) Close() { s.hub.unsubscribe(s) }

// Release IPC promptly even if this peer's network writer is still blocked.
// The caller stops the watcher when its writer exits normally.
func (s *VideoSubscription) CloseOn(done <-chan struct{}) func() {
	stopped := make(chan struct{})
	go func() {
		select {
		case <-done:
			s.Close()
		case <-s.Done:
		case <-stopped:
		}
	}()
	return func() { close(stopped); s.Close() }
}

// VideoHub owns the only Server video reader. Unbuffered frame handoffs mirror
// NanoKVM-Go's reader/writer backpressure: a blocked network writer stops the
// next read, fills the bounded IPC socket, and stops kvm_vin's AX_VENC_GetStream.
// Axera's bounded input/output FIFOs own congestion. No Go frame queue or
// application-level frame dropping absorbs this pressure. All selected peers
// participate in the handoff; disconnecting one unblocks it immediately.
//
// Read/copy, settings and socket release are serialized. Network handoffs take
// place outside that mutex, so controls, mode changes and cancellation remain
// available even while a writer is blocked.
type VideoHub struct {
	mutex      sync.Mutex
	source     videoSource
	screen     *common.Screen
	mode       atomic.Int32
	generation atomic.Uint64
	idrPending atomic.Bool
	modeDone   chan struct{}
	subs       map[*VideoSubscription]struct{}
	running    bool
	started    time.Time
	lastFrame  time.Time
	lastIDR    time.Time
	onFrame    func()
}

var videoHubOnce sync.Once
var videoHub *VideoHub

func GetVideoHub() *VideoHub {
	videoHubOnce.Do(func() {
		videoHub = newVideoHub(common.GetKvmVision(), common.GetScreen())
		videoHub.onFrame = func() { GetFrameRateCounter().Update() }
	})
	return videoHub
}

func newVideoHub(source videoSource, screen *common.Screen) *VideoHub {
	h := &VideoHub{source: source, screen: screen, subs: make(map[*VideoSubscription]struct{}), started: time.Now(), modeDone: make(chan struct{})}
	h.mode.Store(int32(screen.Snapshot().StreamType))
	h.generation.Store(1)
	return h
}

func (h *VideoHub) IsMode(mode int) bool            { return int(h.mode.Load()) == mode }
func (h *VideoHub) IsCurrent(frame VideoFrame) bool { return h.generation.Load() == frame.Generation }

func (h *VideoHub) SetMode(mode int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if !h.IsMode(mode) {
		// Complete the old read and release its socket before a new reader can
		// configure the codec. Reconnection replays libkvm's confirmed settings.
		h.source.CloseVideoStream()
		close(h.modeDone) // cancel an old, blocked network handoff
		h.modeDone = make(chan struct{})
		h.mode.Store(int32(mode))
		h.generation.Add(1)
		h.lastFrame = time.Time{}
		h.lastIDR = time.Time{}
		h.idrPending.Store(true)
		h.screen.Update(func(s *common.ScreenState) { s.StreamType = mode })
		for sub := range h.subs {
			sub.waitKey = true
			if sub.mode != mode {
				h.statusLocked(sub, -4)
			}
		}
	}
	h.startLocked()
}

// Subscribe does not change the mode: a delayed ICE callback from an old peer
// must not undo the user's newer mode selection.
func (h *VideoHub) Subscribe(mode int) *VideoSubscription {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	frames := make(chan VideoFrame)
	updates := make(chan int, 1) // coalesced status notifications, never frames
	done := make(chan struct{})
	sub := &VideoSubscription{Frames: frames, Updates: updates, Done: done, frames: frames, updates: updates, done: done, hub: h, mode: mode, waitKey: true}
	h.subs[sub] = struct{}{}
	if h.IsMode(mode) {
		h.idrPending.Store(true)
	} else {
		h.statusLocked(sub, -4)
	}
	h.startLocked()
	return sub
}

func (h *VideoHub) unsubscribe(sub *VideoSubscription) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if _, exists := h.subs[sub]; !exists {
		return
	}
	delete(h.subs, sub)
	close(sub.done) // Frames stays open: a handoff may be outside the mutex.
	if !h.hasActiveLocked() {
		h.source.CloseVideoStream()
		h.lastFrame = time.Time{}
	}
}

func (h *VideoHub) hasActiveLocked() bool {
	for sub := range h.subs {
		if h.IsMode(sub.mode) {
			return true
		}
	}
	return false
}

func (h *VideoHub) startLocked() {
	if h.running || !h.hasActiveLocked() {
		return
	}
	h.running = true
	go h.run()
}

func (h *VideoHub) run() {
	ticker := time.NewTicker(time.Second / 120)
	defer ticker.Stop()
	for range ticker.C {
		h.mutex.Lock()
		if !h.hasActiveLocked() {
			h.running = false // Subscribe holds this mutex: no lost restart.
			h.mutex.Unlock()
			return
		}
		frame, targets := h.readLocked()
		modeDone := h.modeDone
		h.mutex.Unlock()
		for _, sub := range targets {
			select {
			case <-sub.done:
			case <-modeDone:
			case sub.frames <- frame:
			}
		}
	}
}

func (h *VideoHub) statusLocked(sub *VideoSubscription, status int) {
	if sub.status == status {
		return
	}
	sub.status = status
	select {
	case <-sub.updates:
	default:
	}
	sub.updates <- status
}

func (h *VideoHub) readLocked() (VideoFrame, []*VideoSubscription) {
	mode := int(h.mode.Load())
	h.consumeIDRLocked(mode)
	var data []byte
	var result int
	// Explicit controls own quality/bitrate; legacy per-read arguments must not
	// undo HTTP controls or replay a stale Go variable after reconnection.
	switch mode {
	case common.STREAM_TYPE_MJPEG:
		data, result = h.source.ReadMjpeg(0, 0, 0)
	case common.STREAM_TYPE_H264_WEBRTC, common.STREAM_TYPE_H264_DIRECT:
		data, result = h.source.ReadH264(0, 0, 0)
	case common.STREAM_TYPE_H265_WEBRTC, common.STREAM_TYPE_H265_DIRECT:
		data, result = h.source.ReadH265(0, 0, 0)
	}
	validType := mode == common.STREAM_TYPE_MJPEG && result == int(common.IMG_MJPEG_TYPE) ||
		(mode == common.STREAM_TYPE_H264_WEBRTC || mode == common.STREAM_TYPE_H264_DIRECT) && (result == int(common.IMG_H264_TYPE_IF) || result == int(common.IMG_H264_TYPE_PF)) ||
		(mode == common.STREAM_TYPE_H265_WEBRTC || mode == common.STREAM_TYPE_H265_DIRECT) && (result == int(common.IMG_H265_TYPE_IF) || result == int(common.IMG_H265_TYPE_PF))
	if result < 0 || len(data) == 0 || !validType {
		if result >= 0 {
			result = -1
		}
		h.lastFrame = time.Time{}
		h.idrPending.Store(true)
		for sub := range h.subs {
			if sub.mode != mode {
				continue
			}
			sub.waitKey = true
			h.statusLocked(sub, result)
		}
		return VideoFrame{}, nil
	}
	key := result == int(common.IMG_MJPEG_TYPE) || result == int(common.IMG_H264_TYPE_IF) || result == int(common.IMG_H265_TYPE_IF)
	now := time.Now()
	duration := now.Sub(h.lastFrame)
	if h.lastFrame.IsZero() {
		fps := h.screen.Snapshot().FPS
		if fps == 0 {
			fps = 60
		}
		duration = time.Second / time.Duration(fps)
	}
	h.lastFrame = now
	frame := VideoFrame{Data: data, Result: result, KeyFrame: key, Duration: duration, Timestamp: now.Sub(h.started).Microseconds(), Generation: h.generation.Load()}
	targets := make([]*VideoSubscription, 0, len(h.subs))
	for sub := range h.subs {
		if sub.mode != mode {
			continue
		}
		if sub.waitKey && !key {
			h.idrPending.Store(true)
			continue
		}
		sub.waitKey = false
		h.statusLocked(sub, 1)
		targets = append(targets, sub)
	}
	if h.onFrame != nil {
		h.onFrame()
	}
	return frame, targets
}

// RTCP readers only set a flag. Socket reads/control deadlines must never block
// the feedback path. The capture worker consumes and rate-limits requests.
func (h *VideoHub) RequestKeyFrame(mode int) {
	if h.IsMode(mode) && mode != common.STREAM_TYPE_MJPEG {
		h.idrPending.Store(true)
	}
}

func (h *VideoHub) consumeIDRLocked(mode int) {
	if mode == common.STREAM_TYPE_MJPEG || !h.idrPending.Load() || time.Since(h.lastIDR) < 100*time.Millisecond {
		return
	}
	if h.idrPending.Swap(false) && h.source.RequestIDR() == 0 {
		h.lastIDR = time.Now()
	} else {
		h.idrPending.Store(true)
	}
}

// Commit UI state only after libkvm confirms a setting. Before first connection
// libkvm caches it, then restores the cache on mode changes/producer restarts.
func (h *VideoHub) SetQuality(value uint16) int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	var result int
	if value <= 100 {
		result = h.source.SetMjpegQuality(value)
		if result == 0 {
			h.screen.Update(func(s *common.ScreenState) { s.Quality = value })
		}
	} else {
		result = h.source.SetBitrate(value)
		if result == 0 {
			h.screen.Update(func(s *common.ScreenState) { s.BitRate = value })
		}
	}
	return result
}

func (h *VideoHub) SetRateControl(mode uint8) int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	result := h.source.SetRateControl(mode)
	if result == 0 {
		h.screen.Update(func(s *common.ScreenState) { s.RateControlMode = mode })
	}
	return result
}

func (h *VideoHub) SetFps(fps uint8) int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	result := h.source.SetFps(fps)
	if result == 0 {
		h.screen.Update(func(s *common.ScreenState) { s.FPS = fps })
	}
	return result
}

func (h *VideoHub) SetGop(gop uint8) int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	result := h.source.SetGop(gop)
	if result == 0 {
		h.screen.Update(func(s *common.ScreenState) { s.GOP = gop })
	}
	return result
}
