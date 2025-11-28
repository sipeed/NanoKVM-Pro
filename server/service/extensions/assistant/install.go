package assistant

import (
	"NanoKVM-Server/proto"
	"bufio"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const Requirements = "/tmp/requirements.txt"

func (s *Service) Install(c *gin.Context) {
	var rsp proto.Response

	if err := createRequirements(); err != nil {
		rsp.ErrRsp(c, -1, "failed to create requirements")
		return
	}

	rsp.OkRsp(c)
}

func createRequirements() error {
	requirements := []string{
		"Flask>=3.1.1",
		"flask-cors>=6.0.1",
		"Flask-SocketIO==5.5.1",
		"pillow>=11.3.0",
		"openai>=1.107.3",
		"requests>=2.32.3",
		"dashscope>=1.24.5",
	}

	file, err := os.OpenFile(Requirements, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Errorf("Failed to open requirements.txt: %v", err)
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	writer := bufio.NewWriter(file)
	defer func() {
		_ = writer.Flush()
	}()

	for _, requirement := range requirements {
		_, err := writer.WriteString(requirement + "\n")
		if err != nil {
			log.Errorf("Failed to write to requirements.txt: %v", err)
			return err
		}
	}

	log.Debugf("Created requirements.txt")
	return nil
}
