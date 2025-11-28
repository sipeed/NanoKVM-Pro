package vm

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/service/ui"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Service) SetLedConfig(c *gin.Context) {
	var req proto.SetLedStripReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid argument")
		return
	}

	args := map[string]any{
		"on":         req.On,
		"hor":        req.Hor,
		"ver":        req.Ver,
		"brightness": req.Brightness,
	}

	status := ui.SetLedStrip(args)
	if status == nil || status.Status != "ok" {
		rsp.ErrRsp(c, -2, "failed to set led strip")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("set led strip success: %v", req)
}

func (s *Service) GetLedConfig(c *gin.Context) {
	var rsp proto.Response

	res, err := ui.GetLedStrip()
	if err != nil {
		rsp.ErrRsp(c, -1, err.Error())
		return
	}

	rsp.OkRspWithData(c, res)
}
