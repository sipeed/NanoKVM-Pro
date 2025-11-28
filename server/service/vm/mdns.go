package vm

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const AvahiService = "avahi-daemon.service"

func (s *Service) GetMdnsState(c *gin.Context) {
	var rsp proto.Response

	rsp.OkRspWithData(c, &proto.GetMdnsStateRsp{
		Enabled: isAvahiRunning(),
	})
}

func (s *Service) EnableMdns(c *gin.Context) {
	var rsp proto.Response

	err := utils.StartService(AvahiService, true)
	if err != nil {
		log.Errorf("failed to start avahi: %s", err)
		rsp.ErrRsp(c, -1, "failed to enable mdns")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("%s started", AvahiService)
}

func (s *Service) DisableMdns(c *gin.Context) {
	var rsp proto.Response

	err := utils.StopService(AvahiService, true)
	if err != nil {
		log.Errorf("failed to stop avahi: %s", err)
		rsp.ErrRsp(c, -1, "failed to disable mdns")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("%s stopped", AvahiService)
}

func isAvahiRunning() bool {
	isRunning, err := utils.IsServiceRunning(AvahiService)
	if err != nil {
		log.Errorf("failed to get avahi status: %s", err)
		return false
	}
	return isRunning
}
