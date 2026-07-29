package vm

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/service/ui"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const minutesPerDay = 24 * 60

func (s *Service) GetLcdScreenOff(c *gin.Context) {
	var rsp proto.Response

	screenType, err := getScreenType()
	if err != nil {
		rsp.ErrRsp(c, -1, "read screen type failed")
		return
	}
	if screenType == screenTypeATX {
		rsp.OkRspWithData(c, &proto.GetLcdScreenOffRsp{Supported: false})
		return
	}

	screenOff, err := ui.GetScreenOff()
	if err != nil {
		rsp.ErrRsp(c, -2, "get screen-off schedule failed")
		return
	}

	rsp.OkRspWithData(c, &proto.GetLcdScreenOffRsp{
		Supported:   true,
		Enabled:     screenOff.Enabled,
		StartMinute: screenOff.StartMinute,
		EndMinute:   screenOff.EndMinute,
	})
}

func (s *Service) SetLcdScreenOff(c *gin.Context) {
	var req proto.SetLcdScreenOffReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}
	if err := validateLcdScreenOff(req); err != nil {
		rsp.ErrRsp(c, -1, err.Error())
		return
	}

	screenType, err := getScreenType()
	if err != nil {
		rsp.ErrRsp(c, -2, "read screen type failed")
		return
	}
	if screenType != screenTypeDesk {
		rsp.ErrRsp(c, -3, "screen-off scheduling is only supported on Desk")
		return
	}

	if err := ui.SetScreenOff(ui.ScreenOff{
		Enabled:     *req.Enabled,
		StartMinute: *req.StartMinute,
		EndMinute:   *req.EndMinute,
	}); err != nil {
		rsp.ErrRsp(c, -4, "set screen-off schedule failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("set LCD screen-off schedule: enabled=%t start=%d end=%d", *req.Enabled, *req.StartMinute, *req.EndMinute)
}

func validateLcdScreenOff(req proto.SetLcdScreenOffReq) error {
	if req.Enabled == nil || req.StartMinute == nil || req.EndMinute == nil {
		return fmt.Errorf("enabled, startMinute and endMinute are required")
	}
	if *req.StartMinute < 0 || *req.StartMinute >= minutesPerDay ||
		*req.EndMinute < 0 || *req.EndMinute >= minutesPerDay {
		return fmt.Errorf("startMinute and endMinute must be between 0 and %d", minutesPerDay-1)
	}
	if *req.StartMinute == *req.EndMinute {
		return fmt.Errorf("startMinute and endMinute must differ")
	}
	return nil
}
