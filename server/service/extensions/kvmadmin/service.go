package kvmadmin

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/utils"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Install(c *gin.Context) {
	var rsp proto.Response

	if err := install(); err != nil {
		rsp.ErrRsp(c, -1, "install failed")
		log.Errorf("install failed, err:%v", err)
		return
	}

	time.Sleep(5 * time.Second)

	rsp.OkRsp(c)
	log.Debugf("install success")
}

func (s *Service) Uninstall(c *gin.Context) {
	var rsp proto.Response

	if err := uninstall(); err != nil {
		rsp.ErrRsp(c, -1, "uninstall failed")
		log.Errorf("uninstall failed, err:%v", err)
		return
	}

	rsp.OkRsp(c)
	log.Debugf("uninstall success")
}

func (s *Service) Start(c *gin.Context) {
	var rsp proto.Response

	if err := utils.StartService(Kvmadmin, true); err != nil {
		rsp.ErrRsp(c, -1, "start failed")
		log.Errorf("failed to start kvmadmin service: %v", err)
		return
	}

	rsp.OkRsp(c)
	log.Debugf("kvmadmin service started")
}

func (s *Service) Stop(c *gin.Context) {
	var rsp proto.Response

	if err := utils.StopService(Kvmadmin, true); err != nil {
		rsp.ErrRsp(c, -1, "stop failed")
		log.Errorf("failed to stop kvmadmin service: %v", err)
		return
	}

	rsp.OkRsp(c)
	log.Debugf("kvmadmin service stopped")
}

func (s *Service) GetStatus(c *gin.Context) {
	var rsp proto.Response

	data := proto.GetKvmadminStatusRsp{}

	if !isInstalled() {
		data.State = "notInstall"
		rsp.OkRspWithData(c, &data)
		return
	}

	isRunning, _ := utils.IsServiceRunning(Kvmadmin)

	if isRunning {
		data.State = "running"
	} else {
		data.State = "notRunning"
	}

	rsp.OkRspWithData(c, &data)
}
