package hid

import (
	"NanoKVM-Server/proto"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	ModeNormal  = "normal"
	ModeHidOnly = "hid-only"
	HidOnlyFlag = "/dev/shm/tmp/hid_only"
	USBScript   = "/kvmapp/scripts/usbdev.sh"
)

var hidModeCmdMap = map[string]string{
	ModeNormal:  "restart",
	ModeHidOnly: "hid-only",
}

func (s *Service) GetHidMode(c *gin.Context) {
	var rsp proto.Response

	mode := getHidMode()

	rsp.OkRspWithData(c, &proto.GetHidModeRsp{
		Mode: mode,
	})
	log.Debugf("get hid mode: %s", mode)
}

func (s *Service) SetHidMode(c *gin.Context) {
	var req proto.SetHidModeReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	arg, ok := hidModeCmdMap[req.Mode]
	if !ok {
		rsp.ErrRsp(c, -2, "invalid arguments")
		return
	}

	if mode := getHidMode(); req.Mode == mode {
		rsp.OkRsp(c)
		return
	}

	h := GetHid()
	h.Lock()
	h.CloseNoLock()
	defer func() {
		h.OpenNoLock()
		h.Unlock()
	}()

	if err := exec.Command("bash", USBScript, arg).Run(); err != nil {
		rsp.ErrRsp(c, -3, "failed to set hid mode")
		log.Errorf("Failed to execute script: %v: %s", USBScript, err)
		return
	}

	time.Sleep(3 * time.Second)

	rsp.OkRsp(c)
	log.Debugf("set hid mode: %s", req.Mode)
}

func (s *Service) Reset(c *gin.Context) {
	var rsp proto.Response

	mode := getHidMode()
	arg, ok := hidModeCmdMap[mode]
	if !ok {
		rsp.ErrRsp(c, -1, "invalid hid mode")
		log.Errorf("invalid hid mode: %s", mode)
		return
	}

	h := GetHid()
	h.Lock()
	h.CloseNoLock()
	defer func() {
		h.OpenNoLock()
		h.Unlock()
	}()

	if err := exec.Command("bash", USBScript, arg).Run(); err != nil {
		rsp.ErrRsp(c, -2, "failed to reset")
		log.Errorf("Failed to execute script: %v: %s", USBScript, err)
		return
	}

	time.Sleep(3 * time.Second)

	rsp.OkRsp(c)
	log.Debugf("reset hid success")
}

func getHidMode() string {
	_, err := os.Stat(HidOnlyFlag)
	if err != nil {
		return ModeNormal
	}
	return ModeHidOnly
}
