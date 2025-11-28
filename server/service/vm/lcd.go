package vm

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/service/ui"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Service) GetLcdTimeFormat(c *gin.Context) {
	var rsp proto.Response

	lvtime, err := ui.GetLvtime()
	if err != nil {
		rsp.ErrRsp(c, -1, "get time format failed")
		return
	}

	rsp.OkRspWithData(c, &proto.GetLcdTimeFormatRsp{
		Format: lvtime.Format,
	})
	log.Debugf("get LCD time format: %s", lvtime.Format)
}

func (s *Service) SetLcdTimeFormat(c *gin.Context) {
	var req proto.SetLcdTimeFormatReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid parameters")
		return
	}

	format := strings.ToLower(req.Format)
	if err := ui.SetLvtime(format); err != nil {
		rsp.ErrRsp(c, -2, "set time format failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("set LCD time format: %s", req.Format)
}
