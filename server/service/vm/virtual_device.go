package vm

import (
	"NanoKVM-Server/service/hid"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"NanoKVM-Server/proto"
)

const (
	sdCard          = "/dev/mmcblk1"
	virtualDisk     = "/boot/usb.disk1"
	virtualDiskFlag = "/sys/kernel/config/usb_gadget/g0/functions/mass_storage.disk1"

	virtualNetwork     = "/boot/usb.ncm"
	virtualNetworkFlag = "/sys/kernel/config/usb_gadget/g0/configs/c.1/ncm.usb0"
)

var (
	mountNetworkCommands = []string{
		"/kvmapp/scripts/usbdev.sh stop",
		"touch /boot/usb.ncm",
		"/kvmapp/scripts/usbdev.sh start",
	}

	unmountNetworkCommands = []string{
		"/kvmapp/scripts/usbdev.sh stop",
		"rm -rf /boot/usb.ncm",
		"/kvmapp/scripts/usbdev.sh start",
	}

	mountDiskCommands = []string{
		"touch /boot/usb.disk1",
		"/kvmapp/scripts/usbdev.sh stop",
		"/kvmapp/scripts/usbdev.sh start",
	}

	unmountDiskCommands = []string{
		"/kvmapp/scripts/usbdev.sh stop",
		"rm /boot/usb.disk1",
		"/kvmapp/scripts/usbdev.sh start",
	}
)

func executeShell(command ...string) (stdout string, stderr string, exitCode int, err error) {
	if len(command) == 0 {
		return "", "no command provided", -1, fmt.Errorf("no command provided")
	}
	defaultShell := "/bin/sh"
	var cmd *exec.Cmd
	if len(command) == 1 {
		cmd = exec.Command(defaultShell, "-c", command[0])
	} else {
		fullCmd := strings.Join(command, " && ")
		cmd = exec.Command(defaultShell, "-c", fullCmd)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	dir, err := os.Getwd()
	if err != nil {
		dir = "/"
	}
	cmd.Dir = dir
	cmd.Env = os.Environ()

	err = cmd.Run()
	stdout = strings.TrimSpace(stdoutBuf.String())
	stderr = strings.TrimSpace(stderrBuf.String())

	exitCode = 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return stdout, stderr, exitCode, err
}

func (s *Service) GetVirtualDevice(c *gin.Context) {
	var rsp proto.Response

	rsp.OkRspWithData(c, &proto.GetVirtualDeviceRsp{
		Disk:    isVirtualDiskExist(),
		Network: isVirtualNetworkExist(),
		SdCard:  isSdCardExist(),
	})
	log.Debugf("get virtual device success")
}

func (s *Service) UpdateVirtualDevice(c *gin.Context) {
	var req proto.UpdateVirtualDeviceReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid argument")
		return
	}

	if req.Device == "disk" {
		if !isVirtualDiskExist() && !isSdCardExist() {
			rsp.ErrRsp(c, -2, "SD card not exist")
			return
		}
	}

	var device string
	var commands []string

	switch req.Device {
	case "disk":
		device = virtualDisk

		if !isVirtualDiskExist() {
			commands = mountDiskCommands
		} else {
			commands = unmountDiskCommands
		}
	case "network":
		device = virtualNetwork

		if !isVirtualNetworkExist() {
			commands = mountNetworkCommands
		} else {
			commands = unmountNetworkCommands
		}
	default:
		rsp.ErrRsp(c, -3, "invalid arguments")
		return
	}

	h := hid.GetHid()
	h.Lock()
	h.CloseNoLock()
	defer func() {
		h.OpenNoLock()
		h.Unlock()
	}()

	stdout, _, code, err := executeShell(commands...)
	if code != 0 || err != nil {
		log.Errorf("failed to execute virutal device commands: %s, exit code: %d,the out put %s", err, code, stdout)
		rsp.ErrRsp(c, -4, "failed to mount device")
		return
	}

	on, _ := isDeviceExist(device)
	rsp.OkRspWithData(c, &proto.UpdateVirtualDeviceRsp{
		On: on,
	})

	log.Debugf("update virtual device %s success", req.Device)
}

func isDeviceExist(device string) (bool, error) {
	_, err := os.Stat(device)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	log.Errorf("check file %s err: %s", device, err)
	return false, err
}

func isVirtualDiskExist() bool {
	exist, _ := isDeviceExist(virtualDiskFlag)
	return exist
}

func isVirtualNetworkExist() bool {
	exist, _ := isDeviceExist(virtualNetworkFlag)
	return exist
}

func isSdCardExist() bool {
	exist, _ := isDeviceExist(sdCard)
	return exist
}
