package vm

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	lowPowerService = "kvm-sleep.service"
)

func (s *Service) GetLowPower(c *gin.Context) {
	var rsp proto.Response

	enabled, err := utils.IsServiceRunning(lowPowerService)
	if err != nil {
		rsp.ErrRsp(c, -1, "failed to get status")
		return
	}

	data := &proto.GetLowPowerRsp{
		Enabled: enabled,
	}

	rsp.OkRspWithData(c, data)
	log.Debugf("get low power mode status: %v", enabled)
}

func (s *Service) SetLowPower(c *gin.Context) {
	var req proto.SetLowPowerReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if req.Enable {
		if err := utils.StartService(lowPowerService, true); err != nil {
			log.Errorf("failed to start service: %v", err)
			rsp.ErrRsp(c, -2, "failed to start service")
			return
		}
	} else {
		if err := utils.StopService(lowPowerService, true); err != nil {
			log.Errorf("failed to stop service: %v", err)
			rsp.ErrRsp(c, -3, "failed to stop service")
			return
		}
	}

	rsp.OkRsp(c)
	log.Debugf("set low power mode: %v", req.Enable)
}
