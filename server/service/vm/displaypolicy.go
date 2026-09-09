package vm

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/service/ui"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	modeAlwaysOn  = "alwaysOn"
	modeIdleClock = "idleClock"
	modeIdleOff   = "idleOff"
)

func (s *Service) GetLcdDisplayPolicy(c *gin.Context) {
	var rsp proto.Response
	screenType, err := getScreenType()
	if err != nil {
		rsp.ErrRsp(c, -1, "read screen type failed")
		return
	}
	if screenType == screenTypeATX {
		rsp.OkRspWithData(c, &proto.GetLcdDisplayPolicyRsp{ScreenType: string(screenType)})
		return
	}

	policy, err := ui.GetDisplayPolicy()
	if err == nil {
		rsp.OkRspWithData(c, displayPolicyResponse(screenType, policy))
		return
	}
	if !ui.IsNotFound(err) {
		rsp.ErrRsp(c, -2, "get LCD display policy failed")
		return
	}

	// Older kvm_ui exposes only the scheduled screen-off endpoint.
	screenOff, legacyErr := ui.GetScreenOff()
	if legacyErr != nil {
		rsp.ErrRsp(c, -2, "get LCD display policy failed")
		return
	}
	rsp.OkRspWithData(c, &proto.GetLcdDisplayPolicyRsp{
		ScreenType: string(screenType),
		Schedule: proto.LcdDisplaySchedule{
			Supported:   true,
			Enabled:     screenOff.Enabled,
			StartMinute: screenOff.StartMinute,
			EndMinute:   screenOff.EndMinute,
		},
	})
}

func (s *Service) SetLcdDisplayPolicy(c *gin.Context) {
	var req proto.SetLcdDisplayPolicyReq
	var rsp proto.Response
	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}
	if err := validateDisplayPolicyRequest(req); err != nil {
		rsp.ErrRsp(c, -1, err.Error())
		return
	}
	screenType, err := getScreenType()
	if err != nil {
		rsp.ErrRsp(c, -2, "read screen type failed")
		return
	}
	if screenType != screenTypeDesk {
		rsp.ErrRsp(c, -3, "LCD display policy is only supported on Desk")
		return
	}

	update := ui.DisplayPolicyUpdate{Mode: req.Mode}
	if req.Schedule != nil {
		update.Schedule = &ui.DisplaySchedule{
			Enabled:     *req.Schedule.Enabled,
			StartMinute: *req.Schedule.StartMinute,
			EndMinute:   *req.Schedule.EndMinute,
		}
	}
	if err := ui.SetDisplayPolicy(update); err == nil {
		rsp.OkRsp(c)
		return
	} else if !ui.IsNotFound(err) {
		rsp.ErrRsp(c, -4, "set LCD display policy failed")
		return
	}

	if req.Mode != nil {
		rsp.ErrRsp(c, -4, "LCD display modes are not supported by this device")
		return
	}
	if req.Schedule == nil {
		rsp.ErrRsp(c, -1, "display policy update is empty")
		return
	}
	if err := ui.SetScreenOff(ui.ScreenOff{
		Enabled:     *req.Schedule.Enabled,
		StartMinute: *req.Schedule.StartMinute,
		EndMinute:   *req.Schedule.EndMinute,
	}); err != nil {
		rsp.ErrRsp(c, -4, "set screen-off schedule failed")
		return
	}
	rsp.OkRsp(c)
	log.Debugf("set LCD display policy using legacy screen-off endpoint: enabled=%t start=%d end=%d", *req.Schedule.Enabled, *req.Schedule.StartMinute, *req.Schedule.EndMinute)
}

func displayPolicyResponse(screenType screenType, policy *ui.DisplayPolicy) *proto.GetLcdDisplayPolicyRsp {
	return &proto.GetLcdDisplayPolicyRsp{
		ScreenType:         string(screenType),
		SupportedModes:     policy.SupportedModes,
		Mode:               policy.Mode,
		ModeTimeoutSeconds: policy.ModeTimeoutSeconds,
		Schedule: proto.LcdDisplaySchedule{
			Supported:          policy.Schedule.Supported,
			Enabled:            policy.Schedule.Enabled,
			StartMinute:        policy.Schedule.StartMinute,
			EndMinute:          policy.Schedule.EndMinute,
			WakeTimeoutSeconds: policy.Schedule.WakeTimeoutSeconds,
		},
	}
}

func validateDisplayPolicyRequest(req proto.SetLcdDisplayPolicyReq) error {
	if req.Mode == nil && req.Schedule == nil {
		return fmt.Errorf("display policy update is empty")
	}
	if req.Mode != nil {
		switch *req.Mode {
		case modeAlwaysOn, modeIdleClock, modeIdleOff:
		default:
			return fmt.Errorf("unsupported LCD display mode %q", *req.Mode)
		}
	}
	if req.Schedule != nil {
		if req.Schedule.Enabled == nil || req.Schedule.StartMinute == nil || req.Schedule.EndMinute == nil {
			return fmt.Errorf("schedule requires enabled, startMinute, and endMinute")
		}
		if *req.Schedule.StartMinute < 0 || *req.Schedule.StartMinute >= minutesPerDay || *req.Schedule.EndMinute < 0 || *req.Schedule.EndMinute >= minutesPerDay {
			return fmt.Errorf("startMinute and endMinute must be between 0 and %d", minutesPerDay-1)
		}
		if *req.Schedule.Enabled && *req.Schedule.StartMinute == *req.Schedule.EndMinute {
			return fmt.Errorf("startMinute and endMinute must differ")
		}
	}
	return nil
}
