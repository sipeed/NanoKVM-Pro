package network

import (
	"NanoKVM-Server/proto"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	StaticIPConfigFile = "/boot/eth.nodhcp"
)

func (s *Service) GetStaticIP(c *gin.Context) {
	var rsp proto.Response

	data := &proto.GetStaticIPRsp{}

	if isStaticIPEnabled() {
		content, err := os.ReadFile(StaticIPConfigFile)
		if err != nil {
			log.Errorf("failed to read static ip config: %v", err)
			rsp.ErrRsp(c, -1, "read static ip failed.")
			return
		}

		data.Enabled = true
		data.IP = string(content)
	}

	rsp.OkRspWithData(c, data)
	log.Debugf("get static ip: %v", data)
}

func (s *Service) SetStaticIP(c *gin.Context) {
	var req proto.SetStaticIPReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if (req.Enabled && req.IP == "") || (!req.Enabled && req.IP != "") {
		rsp.ErrRsp(c, -2, "invalid arguments")
		return
	}

	if req.Enabled {
		err := os.WriteFile(StaticIPConfigFile, []byte(req.IP), 0644)
		if err != nil {
			log.Errorf("failed to write %s: %v", StaticIPConfigFile, err)
			rsp.ErrRsp(c, -3, "enable static ip failed")
			return
		}
	} else {
		err := os.Remove(StaticIPConfigFile)
		if err != nil {
			log.Errorf("failed to remove %s: %v", StaticIPConfigFile, err)
			rsp.ErrRsp(c, -4, "disable static ip failed")
			return
		}
	}

	go restartNetwork()

	rsp.OkRsp(c)
	log.Debugf("set static ip: %v", req)
}

func restartNetwork() {
	for _, command := range []string{"down", "up"} {
		cmd := exec.Command("ip", "link", "set", "eth0", command)
		if err := cmd.Run(); err != nil {
			log.Errorf("failed to %s network: %v", command, err)
		}
	}
}

func isStaticIPEnabled() bool {
	_, err := os.Stat(StaticIPConfigFile)
	return err == nil
}
