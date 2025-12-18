package router

import (
	"github.com/gin-gonic/gin"

	"NanoKVM-Server/middleware"
	"NanoKVM-Server/service/vm"
)

func vmRouter(r *gin.Engine) {
	service := vm.NewService()

	api := r.Group("/api").Use(middleware.CheckToken())

	api.GET("/vm/info", service.GetInfo)         // get device information
	api.GET("/vm/hardware", service.GetHardware) // get hardware version

	api.POST("/vm/gpio", service.SetGpio) // update gpio
	api.GET("/vm/gpio", service.GetGpio)  // get gpio
	api.GET("/vm/terminal/auth", service.TerminalAuth)
	api.GET("/vm/terminal", service.Terminal) // web terminal

	api.GET("/vm/script", service.GetScripts)           // get script
	api.POST("/vm/script/upload", service.UploadScript) // upload script
	api.POST("/vm/script/run", service.RunScript)       // run script
	api.DELETE("/vm/script", service.DeleteScript)      // delete script

	api.GET("/vm/device/virtual", service.GetVirtualDevice)              // get virtual device
	api.POST("/vm/device/virtual", service.UpdateVirtualDevice)          // update virtual device
	api.POST("/vm/device/virtual/refresh", service.RefreshVirtualDevice) // refresh virtual device

	api.GET("/vm/oled", service.GetOLED)                      // get OLED configuration
	api.POST("/vm/oled", service.SetOLED)                     // set OLED configuration
	api.GET("/vm/lcd/time/format", service.GetLcdTimeFormat)  // get LCD time format
	api.POST("/vm/lcd/time/format", service.SetLcdTimeFormat) // set LCD time format

	api.GET("/vm/hdmi/capture", service.GetHdmiCapture)          // get HDMI capture status
	api.POST("/vm/hdmi/capture", service.SetHdmiCapture)         // set HDMI capture status
	api.GET("/vm/hdmi/passthrough", service.GetHdmiPassthrough)  // get HDMI passthrough status
	api.POST("/vm/hdmi/passthrough", service.SetHdmiPassthrough) // set HDMI passthrough status

	api.GET("/vm/ssh", service.GetSSHState)         // get SSH state
	api.POST("/vm/ssh/enable", service.EnableSSH)   // enable SSH
	api.POST("/vm/ssh/disable", service.DisableSSH) // disable SSH

	api.GET("/vm/edid", service.GetEdid)                  // get current EDID
	api.POST("/vm/edid", service.SwitchEdid)              // switch EDID
	api.POST("/vm/edid/upload", service.UploadEdid)       // upload edid
	api.GET("/vm/edid/custom", service.GetCustomEdidList) // get custom EDID list
	api.POST("/vm/edid/delete", service.DeleteEdid)       // delete custom EDID

	api.GET("/vm/mouse-jiggler", service.GetMouseJiggler)   // get mouse jiggler
	api.POST("/vm/mouse-jiggler/", service.SetMouseJiggler) // set mouse jiggler

	api.GET("/vm/hostname", service.GetHostname)  // Get Hostname
	api.POST("/vm/hostname", service.SetHostname) // Set Hostname

	api.GET("/vm/web-title", service.GetWebTitle)  // Get web title
	api.POST("/vm/web-title", service.SetWebTitle) // Set web title

	api.GET("/vm/mdns", service.GetMdnsState)         // get mDNS state
	api.POST("/vm/mdns/enable", service.EnableMdns)   // enable mDNS
	api.POST("/vm/mdns/disable", service.DisableMdns) // disable mDNS

	api.POST("/vm/ledstrip/set", service.SetLedConfig)
	api.GET("/vm/ledstrip/get", service.GetLedConfig)

	api.POST("/vm/system/pikvm", service.SwitchPiKVM) // switch system to PiKVM
	api.POST("/vm/system/reboot", service.Reboot)     // reboot system

	api.POST("/vm/timezone", service.SetTimeZone)     // set time zone
	api.GET("/vm/timezone", service.GetTimeZone)      // get current time zone
	api.GET("/vm/time", service.GetTime)              // get time in real time
	api.GET("/vm/time/status", service.GetTimeStatus) // get time status
	api.POST("/vm/time/sync", service.SyncTime)       // synchronize time
}
