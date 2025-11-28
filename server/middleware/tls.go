package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func Tls() gin.HandlerFunc {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect:          false,
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
	})

	secureFunc := func(c *gin.Context) {
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			c.Abort()
			return
		}

		// Check if redirect response was set
		if status := c.Writer.Status(); status >= 300 && status < 400 {
			c.Abort()
		}
	}

	return secureFunc
}
