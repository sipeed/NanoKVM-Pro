package stream

import (
	"net/http"
	"strconv"
	"strings"

	"NanoKVM-Server/common"
	"NanoKVM-Server/proto"
	"NanoKVM-Server/service/ai"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var captureJPEG = ai.CaptureSnapshotContext

type Service struct {
	video func() *VideoHub
}

func NewService() *Service {
	return &Service{video: GetVideoHub}
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

	if err := s.video().SetRateControl(mode); err != 0 {
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

	s.video().SetMode(streamType)

	rsp.OkRsp(c)
	log.Debugf("update video mode to %s", req.Mode)
}

func (s *Service) GetMode(c *gin.Context) {
	var rsp proto.Response
	hub := s.video()
	rsp.OkRspWithData(c, &proto.GetModeRsp{
		Mode:   common.StreamTypeName(hub.Selected()),
		Active: common.StreamTypeName(hub.Mode()),
	})
}

func (s *Service) Snapshot(c *gin.Context) {
	var rsp proto.Response

	quality := 85
	if raw := c.Query("quality"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 100 {
			rsp.ErrRsp(c, -1, "invalid arguments")
			return
		}
		quality = parsed
	}

	req := ai.SnapshotRequest{Quality: quality}
	if raw := c.Query("timeoutMs"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 0 {
			rsp.ErrRsp(c, -1, "invalid arguments")
			return
		}
		req.TimeoutMS = &parsed
	}

	header, jpeg, err := captureJPEG(c.Request.Context(), req)
	if err != nil || !header.OK || len(jpeg) == 0 {
		msg := header.Message
		if msg == "" && err != nil {
			msg = err.Error()
		}
		if msg == "" {
			msg = "snapshot failed"
		}
		code := header.RetCode
		if code == 0 {
			code = -1
		}
		log.Debugf("stream snapshot failed: %s", msg)
		rsp.Err(code, msg)
		c.JSON(http.StatusServiceUnavailable, rsp)
		return
	}

	c.Header("X-Width", strconv.Itoa(header.Width))
	c.Header("X-Height", strconv.Itoa(header.Height))
	c.Data(http.StatusOK, "image/jpeg", jpeg)
}

func (s *Service) SetQuality(c *gin.Context) {
	var req proto.SetQualityReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if req.Quality < 1 || req.Quality > 20000 {
		rsp.ErrRsp(c, -2, "invalid arguments")
		return
	}

	if s.video().SetQuality(uint16(req.Quality)) != 0 {
		rsp.ErrRsp(c, -3, "set quality failed")
		return
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

	if req.Gop < 1 || req.Gop > 200 {
		rsp.ErrRsp(c, -2, "invalid arguments")
		return
	}

	gop := uint8(req.Gop)
	result := s.video().SetGop(gop)
	if result < 0 {
		log.Errorf("failed to set GOP: %d", result)
		rsp.ErrRsp(c, -3, "failed to set GOP")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("update GOP to %d", req.Gop)
}

func (s *Service) SetFps(c *gin.Context) {
	var req proto.SetFpsReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if req.Fps < 0 || req.Fps > 120 {
		rsp.ErrRsp(c, -2, "invalid arguments")
		return
	}

	fps := uint8(req.Fps)
	result := s.video().SetFps(fps)
	if result < 0 {
		log.Errorf("failed to set FPS: %d", result)
		rsp.ErrRsp(c, -3, "set FPS failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("update FPS to %d", req.Fps)
}
