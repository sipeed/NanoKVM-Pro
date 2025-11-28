package router

import (
	"NanoKVM-Server/middleware"
	"NanoKVM-Server/service/ui"

	"github.com/gin-gonic/gin"
)

func localRouter(r *gin.Engine) {
	api := r.Group("/api").Use(middleware.LocalAuth())

	api.GET("/streamer/local", ui.GetStreamer)
}
