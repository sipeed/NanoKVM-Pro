package mjpeg

import (
	"NanoKVM-Server/common"
	"NanoKVM-Server/service/stream"
	"time"

	"github.com/gin-gonic/gin"
)

func Connect(c *gin.Context) {
	c.Header("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Pragma", "no-cache")
	c.Header("X-Server-Date", time.Now().Format(time.RFC1123))

	hub := stream.GetVideoHub()
	hub.Borrow(common.STREAM_TYPE_MJPEG)
	subscription := hub.Subscribe(common.STREAM_TYPE_MJPEG)
	defer subscription.Close()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case frame, ok := <-subscription.Frames:
			if !ok {
				return
			}
			if frame.Result < 0 || !hub.IsCurrent(frame) || !hub.IsMode(common.STREAM_TYPE_MJPEG) {
				continue
			}
			if err := writeFrame(c, frame.Data); err != nil {
				return
			}
		}
	}
}
