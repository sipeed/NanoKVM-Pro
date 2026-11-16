package ai

import (
	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Service struct{}

type configResponse struct {
	MCPEnabled      bool   `json:"mcpEnabled"`
	MCPAPIKeyMasked string `json:"mcpApiKeyMasked,omitempty"`
}

func NewService() *Service { return &Service{} }

func configView(cfg Config) configResponse {
	view := configResponse{MCPEnabled: cfg.MCPEnabled}
	if cfg.MCPAPIKey != "" {
		if len(cfg.MCPAPIKey) > 18 {
			view.MCPAPIKeyMasked = cfg.MCPAPIKey[:12] + "..." + cfg.MCPAPIKey[len(cfg.MCPAPIKey)-6:]
		} else {
			view.MCPAPIKeyMasked = "********"
		}
	}
	return view
}

func (s *Service) GetConfig(c *gin.Context) {
	var rsp proto.Response
	cfg, err := loadConfig()
	if err != nil {
		log.Errorf("load MCP config: %v", err)
		rsp.ErrRsp(c, -1, "get AI config failed")
		return
	}
	rsp.OkRspWithData(c, configView(cfg))
}

func (s *Service) SetConfig(c *gin.Context) {
	var req struct {
		MCPEnabled *bool `json:"mcpEnabled" form:"mcpEnabled"`
	}
	var rsp proto.Response
	if err := c.ShouldBind(&req); err != nil || req.MCPEnabled == nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}
	cfg, err := updateConfig(func(cfg Config) (Config, error) {
		cfg.MCPEnabled = *req.MCPEnabled
		if cfg.MCPEnabled {
			return ensureAPIKey(cfg)
		}
		return cfg, nil
	})
	if err != nil {
		log.Errorf("save MCP config: %v", err)
		rsp.ErrRsp(c, -2, "operation failed")
		return
	}
	rsp.OkRspWithData(c, configView(cfg))
}

func (s *Service) GetMCPAPIKey(c *gin.Context) {
	var rsp proto.Response
	cfg, err := loadConfig()
	if err != nil {
		rsp.ErrRsp(c, -1, "get MCP API key failed")
		return
	}
	c.Header("Cache-Control", "no-store")
	rsp.OkRspWithData(c, gin.H{"mcpApiKey": cfg.MCPAPIKey})
}

func (s *Service) RegenerateMCPAPIKey(c *gin.Context) {
	var rsp proto.Response
	cfg, err := updateConfig(regenerateAPIKey)
	if err != nil {
		log.Errorf("regenerate MCP API key: %v", err)
		rsp.ErrRsp(c, -1, "operation failed")
		return
	}
	rsp.OkRspWithData(c, configView(cfg))
}
