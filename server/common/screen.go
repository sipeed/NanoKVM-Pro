package common

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

type Screen struct {
	Width  uint16
	Height uint16

	StreamType      int
	RateControlMode uint8
	BitRate         uint16
	FPS             uint8
	GOP             uint8
	Quality         uint16

	RealFPS int
}

const (
	WidthPath  = "/proc/lt6911_info/width"
	HeightPath = "/proc/lt6911_info/height"
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
			Width:  readSize(WidthPath),
			Height: readSize(HeightPath),

			StreamType:      STREAM_TYPE_H264_WEBRTC,
			RateControlMode: RATE_CONTROL_VBR,
			BitRate:         8000,
			FPS:             0,
			GOP:             50,
			Quality:         80,

			RealFPS: 0,
		}
	})

	return screen
}

func (s *Screen) Check() {
	s.Width = readSize(WidthPath)
	s.Height = readSize(HeightPath)

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
