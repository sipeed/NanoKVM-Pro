package assistant

import (
	"NanoKVM-Server/proto"
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Service) Start(ctx *gin.Context) {
	var rsp proto.Response

	output, _ := exec.Command("sh", "-c", "ps -ef | grep cua").CombinedOutput()
	if bytes.Contains(output, []byte("cua_webapp.py")) {
		rsp.OkRspWithData(ctx, gin.H{
			"status": "running",
			"code":   0,
		})
		log.Debugf("%s already running", "cua")
		return
	}

	cmd := exec.Command("sh", "-c", "python /kvmapp/cua/cua_webapp.py --auth 1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		rsp.ErrRsp(ctx, -1, "Failed to start python server")
		log.Errorf("Failed to start python server: %v", err)
		return
	}
	defer func() {
		_ = stdout.Close()
	}()

	go func() {
		if err := cmd.Run(); err != nil {
			log.Errorf("failed to start cua server: %v", err)
		}
	}()

	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			rsp.ErrRsp(ctx, -2, "Failed to start python server")
			log.Errorf("Failed to start python server: %v", err)
			return
		}

		if strings.Contains(line, "Authentication") {
			rsp.OkRspWithData(ctx, gin.H{
				"status": "starting",
				"code":   0,
			})
			return
		}
	}
}
