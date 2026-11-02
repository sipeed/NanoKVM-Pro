// Video transport probe: runs the real Server handlers on loopback only.
package main

import (
	"NanoKVM-Server/common"
	"NanoKVM-Server/service/stream"
	"NanoKVM-Server/service/stream/direct"
	h265direct "NanoKVM-Server/service/stream/h265/direct"
	h265webrtc "NanoKVM-Server/service/stream/h265/webrtc"
	"NanoKVM-Server/service/stream/mjpeg"
	"NanoKVM-Server/service/stream/webrtc"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	service := stream.NewService()
	router.POST("/api/stream/quality", service.SetQuality)
	router.POST("/api/stream/mode", service.SetMode)
	router.POST("/api/stream/fps", service.SetFps)
	router.POST("/api/stream/gop", service.SetGop)
	router.POST("/api/stream/rate-control", service.SetRateControlMode)
	router.GET("/api/stream/h264/direct", direct.Connect)
	router.GET("/api/stream/h265/direct", h265direct.Connect)
	router.GET("/api/stream/h264/webrtc", webrtc.Connect)
	router.GET("/api/stream/h265/webrtc", h265webrtc.Connect)
	router.GET("/api/stream/mjpeg", mjpeg.Connect)
	router.GET("/probe/state", func(c *gin.Context) { c.JSON(http.StatusOK, common.GetScreen().Snapshot()) })
	router.GET("/probe/hold", func(c *gin.Context) {
		milliseconds, _ := strconv.Atoi(c.DefaultQuery("ms", "200"))
		if milliseconds < 0 || milliseconds > 3000 {
			c.Status(http.StatusBadRequest)
			return
		}
		hub := stream.GetVideoHub()
		hub.SetMode(common.STREAM_TYPE_H264_DIRECT)
		sub := hub.Subscribe(common.STREAM_TYPE_H264_DIRECT)
		defer sub.Close()
		select {
		case <-c.Request.Context().Done():
			return
		case <-time.After(time.Duration(milliseconds) * time.Millisecond):
		}
		select {
		case <-c.Request.Context().Done():
			return
		case frame := <-sub.Frames:
			c.JSON(http.StatusOK, gin.H{"key": frame.KeyFrame, "size": len(frame.Data)})
		}
	})
	router.StaticFile("/probe.html", "probe.html")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	go func() { <-signals; common.GetKvmVision().Close(); os.Exit(0) }()
	if err := router.Run("127.0.0.1:18080"); err != nil {
		panic(err)
	}
}
