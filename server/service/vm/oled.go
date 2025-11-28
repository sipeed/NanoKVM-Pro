package vm

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/service/ui"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Service) SetOLED(c *gin.Context) {
	var req proto.SetOledReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if err := ui.SetUI(req.Sleep); err != nil {
		rsp.ErrRsp(c, -2, "operation failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("set OLED sleep: %d", req.Sleep)
}

func (s *Service) GetOLED(c *gin.Context) {
	var rsp proto.Response

	data := &proto.GetOLEDRsp{
		Exist: true,
	}

	content, err := os.ReadFile("/proc/lt6911_info/version")
	if err != nil {
		rsp.ErrRsp(c, -1, "read file failed")
		return
	}

	fileContent := string(content)
	if strings.Contains(fileContent, "ATX") {
		data.Type = "ATX"
	} else if strings.Contains(fileContent, "Desk") {
		data.Type = "DESK"
	} else {
		rsp.ErrRsp(c, -2, "invalid file content")
		return
	}

	uiArgs, err := ui.GetUI()
	if err != nil {
		rsp.OkRspWithData(c, data)
		return
	}

	data.Sleep = uiArgs.SleepS
	rsp.OkRspWithData(c, data)

	log.Debugf("get OLED config successful, sleep %d", uiArgs.SleepS)
}
