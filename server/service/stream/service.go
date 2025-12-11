package stream

import (
	"strings"

	"NanoKVM-Server/common"
	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SetRateControlMode(c *gin.Context) {
	var req proto.SetRateControlModeReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	var mode uint8

	switch strings.ToLower(req.Mode) {
	case "cbr":
		mode = common.RATE_CONTROL_CBR
	case "vbr":
		mode = common.RATE_CONTROL_VBR
	default:
		rsp.ErrRsp(c, -2, "invalid arguments")
		return
	}

	if err := common.GetKvmVision().SetRateControl(mode); err != 0 {
		rsp.ErrRsp(c, -3, "set rate control failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("update rate control mode to %s", req.Mode)
}

func (s *Service) SetMode(c *gin.Context) {
	var req proto.SetModeReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	streamType, ok := common.StreamTypeMap[strings.ToLower(req.Mode)]
	if !ok {
		rsp.ErrRsp(c, -2, "invalid arguments")
		return
	}

	common.GetScreen().StreamType = streamType
	common.GetKvmVision().SetStreamType(uint8(streamType))

	rsp.OkRsp(c)
	log.Debugf("update video mode to %s", req.Mode)
}

func (s *Service) SetQuality(c *gin.Context) {
	var req proto.SetQualityReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	screen := common.GetScreen()
	value := uint16(req.Quality)

	if value < 1 || value > 20000 {
		rsp.ErrRsp(c, -2, "invalid arguments")
		return
	}

	if req.Quality <= 100 {
		screen.Quality = value
	} else {
		screen.BitRate = value
	}

	rsp.OkRsp(c)
	log.Debugf("update quality to %d", req.Quality)
}

func (s *Service) SetGop(c *gin.Context) {
	var req proto.SetGopReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	gop := uint8(req.Gop)
	if gop < 1 || gop > 200 {
		rsp.ErrRsp(c, -2, "invalid arguments")
		return
	}

	common.GetScreen().GOP = gop
	common.GetKvmVision().SetGop(gop)

	rsp.OkRsp(c)
	log.Debugf("update GOP to %d", req.Gop)
}
