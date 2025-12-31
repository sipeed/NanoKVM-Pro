package router

import (
	"NanoKVM-Server/middleware"
	"NanoKVM-Server/service/stream"
	"NanoKVM-Server/service/stream/direct"
	"NanoKVM-Server/service/stream/mjpeg"
	"NanoKVM-Server/service/stream/webrtc"

	"github.com/gin-gonic/gin"
)

func streamRouter(r *gin.Engine) {
	service := stream.NewService()

	api := r.Group("/api").Use(middleware.CheckToken())

	api.POST("/stream/rate-control", service.SetRateControlMode) // set rate control mode
	api.POST("/stream/mode", service.SetMode)                    // set stream mode
	api.POST("/stream/quality", service.SetQuality)              // set stream quality / bit-rate
	api.POST("/stream/gop", service.SetGop)                      // set stream gop
	api.POST("/stream/fps", service.SetFps)                      // set stream fps

	api.GET("/stream/mjpeg", mjpeg.Connect)        // mjpeg stream
	api.GET("/stream/h264/direct", direct.Connect) // h264 stream (direct)
	api.GET("/stream/h264/webrtc", webrtc.Connect) // h264 stream (webrtc)
}
