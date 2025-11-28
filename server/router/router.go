package router

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Init(r *gin.Engine) {
	web(r)
	server(r)
	log.Debugf("router init done")
}

func web(r *gin.Engine) {
	execPath, err := os.Executable()
	if err != nil {
		panic("invalid executable path")
	}

	execDir := filepath.Dir(execPath)
	webPath := fmt.Sprintf("%s/web", execDir)

	r.Use(static.Serve("/", static.LocalFile(webPath, true)))

	r.GET("/kvm", func(c *gin.Context) {
		c.Redirect(302, "/")
	})
}

func server(r *gin.Engine) {
	routers := []func(c *gin.Engine){
		authRouter,
		applicationRouter,
		vmRouter,
		streamRouter,
		storageRouter,
		networkRouter,
		hidRouter,
		wsRouter,
		localRouter,
		extensionsRouter,
		speedtestRouter,
	}

	for _, fn := range routers {
		fn(r)
	}
}
