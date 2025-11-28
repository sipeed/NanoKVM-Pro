package router

import (
	"github.com/gin-gonic/gin"

	"NanoKVM-Server/middleware"
	"NanoKVM-Server/service/storage"
)

func storageRouter(r *gin.Engine) {
	service := storage.NewService()
	api := r.Group("/api").Use(middleware.CheckToken())

	api.GET("/storage/image", service.GetImages)               // get image list
	api.GET("/storage/image/mounted", service.GetMountedImage) // get mounted image
	api.POST("/storage/image/upload", service.UploadImage)     // upload image
	api.POST("/storage/image/mount", service.MountImage)       // mount image
	api.POST("/storage/image/delete", service.DeleteImage)     // delete image

	api.POST("/storage/download/image", service.DownloadImage)       // download image
	api.GET("/storage/download/image/status", service.StatusImage)   // download image
	api.GET("/storage/download/image/enabled", service.ImageEnabled) // download image
}
