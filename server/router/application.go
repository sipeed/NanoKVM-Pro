package router

import (
	"NanoKVM-Server/middleware"
	"NanoKVM-Server/service/application"

	"github.com/gin-gonic/gin"
)

func applicationRouter(r *gin.Engine) {
	service := application.NewService()
	// Keep every application update endpoint behind the same token middleware.
	api := r.Group("/api").Use(middleware.CheckToken())

	api.GET("/application/version", service.GetVersion) // get application version
	api.POST("/application/update", service.Update)     // update application
	api.GET("/application/update-source", service.GetUpdateSource)        // get update source
	api.POST("/application/update-source", service.SetUpdateSource)      // set update source
	api.POST("/application/update-source/reset", service.ResetUpdateSource)  // restore official source

	api.POST("/application/update/manual/inspect", service.InspectManualUpdate) // inspect uploaded package
	api.POST("/application/update/manual/confirm", service.ConfirmManualUpdate) // confirm package install
	api.GET("/application/update/manual/:id", service.GetManualUpdate)           // get manual update state
	api.DELETE("/application/update/manual/:id", service.DeleteManualUpdate) // discard staged package

	api.GET("/application/preview", service.GetPreview)  // get preview updates state
	api.POST("/application/preview", service.SetPreview) // set preview updates state
}
