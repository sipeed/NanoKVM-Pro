package vm

import (
	"encoding/json"
	"os"
	"sync"

	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	MenuBarConfigFile = "/etc/kvm/menubar"
)

var menuBarMutex sync.RWMutex

func (s *Service) GetMenuBarConfig(c *gin.Context) {
	var rsp proto.Response

	config, err := readMenuBarConfig()
	if err != nil {
		rsp.ErrRsp(c, -1, "get menu bar config failed")
		return
	}

	rsp.OkRspWithData(c, &proto.GetMenuBarConfigRsp{
		DisabledItems: config.DisabledItems,
	})

	log.Debugf("get menu bar config successful")
}

func (s *Service) SetMenuBarConfig(c *gin.Context) {
	var req proto.SetMenuBarConfigReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	config, _ := readMenuBarConfig()
	if config == nil {
		config = &proto.MenuBarConfig{}
	}

	config.DisabledItems = req.DisabledItems

	if err := writeMenuBarConfig(config); err != nil {
		rsp.ErrRsp(c, -2, "write config failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("set menu bar config: %+v", req)
}

func readMenuBarConfig() (*proto.MenuBarConfig, error) {
	menuBarMutex.RLock()
	defer menuBarMutex.RUnlock()

	data, err := os.ReadFile(MenuBarConfigFile)
	if err != nil {
		log.Errorf("failed to read %s: %v", MenuBarConfigFile, err)
		return nil, err
	}

	var config proto.MenuBarConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Errorf("failed to unmarshal menu bar: %v", err)
		return nil, err
	}

	return &config, nil
}

func writeMenuBarConfig(config *proto.MenuBarConfig) error {
	menuBarMutex.Lock()
	defer menuBarMutex.Unlock()

	data, err := json.Marshal(config)
	if err != nil {
		log.Errorf("failed to marshal menu bar: %v", err)
		return err
	}

	if err := os.WriteFile(MenuBarConfigFile, data, 0o644); err != nil {
		log.Errorf("failed to write %s: %v", MenuBarConfigFile, err)
		return err
	}

	log.Debugf("write to %s success", MenuBarConfigFile)
	return nil
}
