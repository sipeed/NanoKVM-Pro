package router

import (
	"NanoKVM-Server/middleware"
	"NanoKVM-Server/service/ai"

	"github.com/gin-gonic/gin"
)

func aiRouter(r *gin.Engine) {
	service := ai.NewService()
	api := r.Group("/api/ai").Use(middleware.CheckToken())
	api.GET("/config", service.GetConfig)
	api.POST("/config", service.SetConfig)
	api.POST("/mcp/key", service.GetMCPAPIKey)
	api.POST("/mcp/key/copy", service.GetMCPAPIKey)
	api.POST("/mcp/key/regenerate", service.RegenerateMCPAPIKey)

	r.Any("/api/mcp", ai.MCPAPIKeyMiddleware(), gin.WrapH(ai.NewMCPHandler()))
}
