package ui

import (
	"fmt"
	"net/http"

	"NanoKVM-Server/common"

	"github.com/gin-gonic/gin"
)

func GetStreamer(c *gin.Context) {
	screen := common.GetScreen()

	resolution := Resolution{
		Width:  int(screen.Width),
		Height: int(screen.Height),
	}
	clients := 0
	if screen.StreamType == 0 {
		clients = 1
	}

	rsp := &Response{
		Ok: true,
		Result: Result{
			Delay: 130,
			Params: Params{
				Resolution: fmt.Sprintf("%dx%d", screen.Width, screen.Height),
			},
			Streamer: Streamer{
				H264: H264{
					Bitrate:     int(screen.BitRate),
					Fps:         screen.FPS,
					Gop:         int(screen.GOP),
					RealBitrate: int(screen.BitRate),
				},
				Source: Source{
					CapturedFps: screen.RealFPS,
					Resolution:  resolution,
					Online:      false,
				},
				Stream: Stream{
					Clients:   clients,
					QueuedFps: screen.RealFPS,
				},
			},
		},
	}

	c.JSON(http.StatusOK, rsp)
}
