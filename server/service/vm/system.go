package vm

import (
	"NanoKVM-Server/proto"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Service) SwitchPiKVM(c *gin.Context) {
	var rsp proto.Response

	file := "/etc/kvm/server.txt"
	data := []byte("pikvm")

	if err := os.WriteFile(file, data, 0o644); err != nil {
		rsp.ErrRsp(c, -1, "operation failed")
		log.Errorf("failed to write %s: %s", file, err)
		return
	}

	if err := reboot(); err != nil {
		rsp.ErrRsp(c, -2, "reboot failed")
		return
	}

	rsp.OkRsp(c)
	log.Debug("switch system to PiKVM")
}

func (s *Service) Reboot(c *gin.Context) {
	var rsp proto.Response

	if err := reboot(); err != nil {
		rsp.ErrRsp(c, -1, "operation failed")
		return
	}

	rsp.OkRsp(c)
	log.Debug("system rebooted")
}

func reboot() error {
	log.Println("reboot system...")

	err := exec.Command("reboot").Run()
	if err != nil {
		log.Errorf("failed to reboot: %s", err)
		return err
	}

	return nil
}
