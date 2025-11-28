package router

import (
	"NanoKVM-Server/middleware"
	"NanoKVM-Server/service/speedtest"

	"github.com/gin-gonic/gin"
)

func speedtestRouter(c *gin.Engine) {
	service := speedtest.NewServer()

	api := c.Group("/api/speed").Use(middleware.CheckToken())

	api.GET("ping", service.Ping)
	api.GET("download", service.DownloadSpeedTest)
	api.POST("upload", service.UploadSpeedTest)
}
