package mjpeg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func writeFrame(c *gin.Context, data []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = c.Request.Context().Err()
			if err == nil {
				err = fmt.Errorf("panic recovered in writeFrame: %v", r)
			}
		}
	}()

	if err = http.NewResponseController(c.Writer).SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		return err
	}
	header := "--frame\r\nContent-Type: image/jpeg\r\nContent-Length: " + strconv.Itoa(len(data)) + "\r\n\r\n"
	if _, err = c.Writer.WriteString(header); err != nil {
		return err
	}

	if _, err = c.Writer.Write(data); err != nil {
		return err
	}

	if _, err = c.Writer.Write([]byte("\r\n")); err != nil {
		return err
	}

	c.Writer.Flush()
	return nil
}
