package common

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

type Screen struct {
	Width   uint16
	Height  uint16
	FPS     int
	GOP     uint8
	Quality uint16
	BitRate uint16

	RealFPS    int
	StreamType int
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
			Width:      0,
			Height:     0,
			FPS:        30,
			GOP:        50,
			Quality:    80,
			BitRate:    5000,
			RealFPS:    0,
			StreamType: STREAM_TYPE_H264_WEBRTC,
		}
	})

	checkScreen()
	return screen
}

func checkScreen() {
	screen.Width = readSize(WidthPath)
	screen.Height = readSize(HeightPath)

	if screen.FPS <= 0 || screen.FPS > 200 {
		screen.FPS = 30
	}

	if screen.GOP <= 0 || screen.GOP > 200 {
		screen.GOP = 50
	}

	if screen.Quality <= 0 || screen.Quality > 100 {
		screen.Quality = 80
	}

	if screen.BitRate < 500 || screen.BitRate > 20000 {
		screen.BitRate = 5000
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
