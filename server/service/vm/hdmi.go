package vm

import (
	"os"
	"strings"
	"time"

	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	LT6911Power        = "/proc/lt6911_info/power"
	LT6911HdmiPower    = "/proc/lt6911_info/hdmi_power"
	LT6911LoopoutPower = "/proc/lt6911_info/loopout_power"
)

func (s *Service) GetHdmiCapture(c *gin.Context) {
	var rsp proto.Response

	enabled, err := isHdmiEnabled(LT6911Power)
	if err != nil {
		rsp.ErrRsp(c, -1, "failed to get HDMI capture status")
		return
	}

	rsp.OkRspWithData(c, &proto.GetHdmiCaptureRsp{
		Enabled: enabled,
	})
	log.Debugf("get HDMI capture status: %t", enabled)
}

func (s *Service) SetHdmiCapture(c *gin.Context) {
	var req proto.SetHdmiCaptureReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	status := "off"
	if req.Enabled {
		status = "on"
	}

	if err := os.WriteFile(LT6911Power, []byte(status), 0644); err != nil {
		rsp.ErrRsp(c, -2, "failed to set HDMI capture status")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("set HDMI capture status: %s", status)
}

func (s *Service) GetHdmiPassthrough(c *gin.Context) {
	var rsp proto.Response

	enabled, err := isHdmiEnabled(LT6911LoopoutPower)
	if err != nil {
		rsp.ErrRsp(c, -1, "failed to get HDMI passthrough status")
		return
	}

	rsp.OkRspWithData(c, &proto.GetHdmiPassthroughRsp{
		Enabled: enabled,
	})
	log.Debugf("get HDMI passthrough status: %t", enabled)
}

func (s *Service) SetHdmiPassthrough(c *gin.Context) {
	var req proto.SetHdmiPassthroughReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	var err error
	if req.Enabled {
		err = enableHdmiPassthrough()
	} else {
		err = disableHdmiPassthrough()
	}

	if err != nil {
		rsp.ErrRsp(c, -2, "failed to set HDMI passthrough status")
		return
	}

	time.Sleep(10 * time.Millisecond)

	rsp.OkRsp(c)
	log.Debugf("set HDMI passthrough status: %t", req.Enabled)
}

func isHdmiEnabled(flag string) (bool, error) {
	content, err := os.ReadFile(flag)
	if err != nil {
		return false, err
	}

	enabled := strings.TrimSpace(string(content)) == "on"
	return enabled, nil
}

func enableHdmiPassthrough() error {
	if err := os.WriteFile(LT6911HdmiPower, []byte("0"), 0644); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(LT6911LoopoutPower, []byte("1"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(LT6911HdmiPower, []byte("1"), 0644); err != nil {
		return err
	}
	return nil
}

func disableHdmiPassthrough() error {
	if err := os.WriteFile(LT6911LoopoutPower, []byte("0"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(LT6911HdmiPower, []byte("0"), 0644); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(LT6911HdmiPower, []byte("1"), 0644); err != nil {
		return err
	}
	return nil
}
