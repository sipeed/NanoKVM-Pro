package router

import (
	"NanoKVM-Server/middleware"
	"NanoKVM-Server/service/extensions/assistant"
	"NanoKVM-Server/service/extensions/kvmadmin"
	"NanoKVM-Server/service/extensions/tailscale"

	"github.com/gin-gonic/gin"
)

func extensionsRouter(r *gin.Engine) {
	api := r.Group("/api/extensions").Use(middleware.CheckToken())

	ts := tailscale.NewService()
	assist := assistant.NewService()
	admin := kvmadmin.NewService()

	api.POST("/tailscale/install", ts.Install)     // install tailscale
	api.POST("/tailscale/uninstall", ts.Uninstall) // uninstall tailscale
	api.GET("/tailscale/status", ts.GetStatus)     // get tailscale status
	api.POST("/tailscale/up", ts.Up)               // run tailscale up
	api.POST("/tailscale/down", ts.Down)           // run tailscale down
	api.POST("/tailscale/login", ts.Login)         // tailscale login
	api.POST("/tailscale/logout", ts.Logout)       // tailscale logout
	api.POST("/tailscale/start", ts.Start)         // tailscale start
	api.POST("/tailscale/stop", ts.Stop)           // tailscale stop
	api.POST("/tailscale/restart", ts.Restart)     // tailscale restart

	api.POST("/assistant/install", assist.Install) // install assistant dependencies
	api.POST("/assistant/start", assist.Start)     // start assistant dependencies

	api.POST("/kvmadmin/install", admin.Install)     // install kvmadmin
	api.POST("/kvmadmin/uninstall", admin.Uninstall) // uninstall kvmadmin
	api.POST("/kvmadmin/start", admin.Start)         // run kvmadmin
	api.POST("/kvmadmin/stop", admin.Stop)           // stop kvmadmin
	api.GET("/kvmadmin/status", admin.GetStatus)     // get kvmadmin status
}
