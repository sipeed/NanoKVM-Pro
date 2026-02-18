package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"NanoKVM-Server/common"
	"NanoKVM-Server/config"
	"NanoKVM-Server/logger"
	"NanoKVM-Server/middleware"
	"NanoKVM-Server/router"
	"NanoKVM-Server/service/vm/jiggler"
	"NanoKVM-Server/service/vm"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	initialize()
	defer dispose()

	run()
}

func initialize() {
	logger.Init()

	// init screen parameters
	_ = common.GetScreen()

	// restore HDMI capture/passthrough settings
	vm.RestoreHdmiState()

	// run mouse jiggler
	jiggler.GetJiggler().Run()

	// waiting for exit signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sigChan
		log.Printf("\nReceived signal: %v\n", sig)

		dispose()
		os.Exit(0)
	}()
}

func run() {
	conf := config.GetInstance()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())

	if conf.Authentication == "disable" {
		r.Use(cors.AllowAll())
	}

	if conf.Proto != "http" {
		r.Use(middleware.Tls())
	}

	router.Init(r)
	httpAddr := fmt.Sprintf(":%d", conf.Port.Http)

	if conf.Proto == "http" {
		log.Printf("Starting HTTP server on %s\n", httpAddr)
		if err := r.Run(httpAddr); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	} else {
		httpsAddr := fmt.Sprintf(":%d", conf.Port.Https)
		log.Printf("Starting HTTPS server on %s, %s\n", httpAddr, httpsAddr)

		go runRedirect(httpAddr, httpsAddr)

		if err := r.RunTLS(httpsAddr, conf.Cert.Crt, conf.Cert.Key); err != nil {
			log.Fatalf("HTTPS server failed: %v", err)
		}
	}
}

func runRedirect(httpPort string, httpsPort string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		host, _, _ := net.SplitHostPort(req.Host)
		if host == "" {
			host = req.Host
		}

		targetURL := "https://" + host
		if httpsPort != ":443" {
			targetURL += httpsPort
		}
		targetURL += req.URL.String()

		http.Redirect(w, req, targetURL, http.StatusTemporaryRedirect)
	})

	if err := http.ListenAndServe(httpPort, handler); err != nil {
		log.Fatalf("HTTP redirect server failed: %v", err)
	}
}

func dispose() {
	common.GetKvmVision().Close()
}
