package ai

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MCP always requires its own key, including when web authentication is disabled.
func MCPAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, err := loadConfig()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "MCP auth unavailable"})
			return
		}
		if !cfg.MCPEnabled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "MCP service is disabled"})
			return
		}
		header := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if !strings.HasPrefix(header, "Bearer ") || cfg.MCPAPIKey == "" ||
			subtle.ConstantTimeCompare([]byte(token), []byte(cfg.MCPAPIKey)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid MCP API key"})
			return
		}
		c.Next()
	}
}
