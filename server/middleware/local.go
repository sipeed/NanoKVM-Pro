package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func LocalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ClientIP() == "127.0.0.1" || c.ClientIP() == "::1" || strings.HasPrefix(c.ClientIP(), "127.") {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, "unauthorized")
		c.Abort()
	}
}
