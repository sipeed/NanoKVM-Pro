package common

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Screen struct {
	StreamType      int
	RateControlMode uint8
	BitRate         uint16
	FPS             uint8
	GOP             uint8
	Quality         uint16

	RealFPS int

	// width and height track the live HDMI input resolution. They are read from
	// the streaming loops and written by RefreshSize, so they are kept atomic.
	width  atomic.Uint32
	height atomic.Uint32

	sizeMu        sync.Mutex
	lastSizeCheck time.Time
}

const (
	WidthPath  = "/proc/lt6911_info/width"
	HeightPath = "/proc/lt6911_info/height"

	sizeCheckInterval = 5 * time.Second
)

var (
	screen     *Screen
	screenOnce sync.Once
)

var StreamTypeMap = map[string]int{
	"mjpeg":       STREAM_TYPE_MJPEG,
	"h264-webrtc": STREAM_TYPE_H264_WEBRTC,
	"h264-direct": STREAM_TYPE_H264_DIRECT,
	"h265-webrtc": STREAM_TYPE_H265_WEBRTC,
	"h265-direct": STREAM_TYPE_H265_DIRECT,
}

func GetScreen() *Screen {
	screenOnce.Do(func() {
		screen = &Screen{
			StreamType:      STREAM_TYPE_H264_WEBRTC,
			RateControlMode: RATE_CONTROL_VBR,
			BitRate:         8000,
			FPS:             0,
			GOP:             50,
			Quality:         80,

			RealFPS: 0,
		}

		screen.width.Store(uint32(readSize(WidthPath)))
		screen.height.Store(uint32(readSize(HeightPath)))
		screen.lastSizeCheck = time.Now()
	})

	return screen
}

func (s *Screen) Width() uint16 {
	return uint16(s.width.Load())
}

func (s *Screen) Height() uint16 {
	return uint16(s.height.Load())
}

// RefreshSize re-reads the HDMI input resolution, at most once every
// sizeCheckInterval. The encoder is configured from Width/Height, and the
// capture path links VIN straight to VENC without a scaler, so a stale value
// makes VENC reject every incoming frame and the stream never produces an
// image. Every streaming loop must call this to follow resolution changes.
func (s *Screen) RefreshSize() {
	s.sizeMu.Lock()
	defer s.sizeMu.Unlock()

	if time.Since(s.lastSizeCheck) < sizeCheckInterval {
		return
	}
	s.lastSizeCheck = time.Now()

	// Keep the last known resolution when the driver reports nothing, so a
	// transient read failure does not reconfigure the encoder to 0x0.
	if width := readSize(WidthPath); width > 0 {
		s.width.Store(uint32(width))
	}
	if height := readSize(HeightPath); height > 0 {
		s.height.Store(uint32(height))
	}
}

func (s *Screen) Check() {
	s.RefreshSize()

	if s.FPS < 0 || s.FPS > 120 {
		s.FPS = 0
	}

	if s.BitRate < 1000 || s.BitRate > 20000 {
		if s.RateControlMode == RATE_CONTROL_CBR {
			s.BitRate = 5000
		} else {
			s.BitRate = 8000
		}
	}

	if s.GOP < 1 || s.GOP > 200 {
		s.GOP = 50
	}

	if s.Quality < 1 || s.Quality > 100 {
		s.Quality = 80
	}
}

func readSize(filePath string) uint16 {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0
	}

	width, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 16)
	if err != nil {
		return 0
	}

	return uint16(width)
}
