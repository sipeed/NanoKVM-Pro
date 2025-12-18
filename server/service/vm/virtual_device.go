package vm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"NanoKVM-Server/proto"
	"NanoKVM-Server/service/hid"
	"NanoKVM-Server/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	pathSdCard    = "/dev/mmcblk1"
	pathEmmcImage = "/exfat.img"

	pathVirtualDiskSD      = "/boot/usb.disk1.sd"
	pathVirtualDiskEmmc    = "/boot/usb.disk1.emmc"
	pathVirtualDiskMounted = "/sys/kernel/config/usb_gadget/g0/functions/mass_storage.disk1/lun.0/file"

	pathVirtualNetwork = "/boot/usb.ncm"
	pathNetworkFlag    = "/sys/kernel/config/usb_gadget/g0/configs/c.1/ncm.usb0"
)

const (
	scriptUsbDev    = "/kvmapp/scripts/usbdev.sh"
	scriptMountEmmc = "/kvmcomm/scripts/mount_emmc.py"
)

const (
	DiskTypeSDCard = "sdcard"
	DiskTypeEmmc   = "emmc"
)

type VirtualDevice interface {
	IsMounted() bool
	Mount() error
	Unmount() error
}

func (s *Service) GetVirtualDevice(c *gin.Context) {
	var rsp proto.Response

	rsp.OkRspWithData(c, &proto.GetVirtualDeviceRsp{
		IsNetworkEnabled: isVirtualNetworkMounted(),
		MountedDisk:      getMountedDiskType(),
		IsSdCardExist:    isSdCardPresent(),
		IsEmmcExist:      isEmmcImagePresent(),
	})

	log.Debug("get virtual device success")
}

func (s *Service) UpdateVirtualDevice(c *gin.Context) {
	var req proto.UpdateVirtualDeviceReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid argument")
		return
	}

	commands, err := resolveDeviceCommands(req.Device, req.Type)
	if err != nil {
		rsp.ErrRsp(c, -2, err.Error())
		return
	}

	if err := executeWithHidLock(commands); err != nil {
		log.Errorf("failed to execute virtual device commands: %s", err)
		rsp.ErrRsp(c, -3, "failed to mount device")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("update virtual device %s success", req.Device)
}

func (s *Service) RefreshVirtualDevice(c *gin.Context) {
	var req proto.RefreshVirtualDeviceReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil || req.Device != "emmc" {
		rsp.ErrRsp(c, -1, "invalid argument")
		return
	}

	err := exec.Command(scriptMountEmmc, "restart").Run()
	if err != nil {
		log.Errorf("failed to restart emmc image: %v", err)
		rsp.ErrRsp(c, -2, "failed to refresh")
		return
	}

	rsp.OkRsp(c)
	log.Debug("refresh virtual device success")
}

func resolveDeviceCommands(device, diskType string) ([]string, error) {
	switch device {
	case "network":
		return getNetworkCommands(), nil
	case "disk":
		return getDiskCommands(diskType)
	default:
		return nil, fmt.Errorf("invalid device: %s", device)
	}
}

func getNetworkCommands() []string {
	cmd := []string{
		scriptUsbDev + " stop",
	}

	if !isVirtualNetworkMounted() {
		cmd = append(cmd, "touch "+pathVirtualNetwork)
	} else {
		cmd = append(cmd, "rm -rf "+pathVirtualNetwork)
	}

	cmd = append(cmd, scriptUsbDev+" start")
	return cmd
}

func getDiskCommands(diskType string) ([]string, error) {
	// For emmc, ensure image exists before mounting
	if diskType == DiskTypeEmmc {
		if err := ensureEmmcImageExists(); err != nil {
			return nil, err
		}
	}

	cmd := []string{
		"rm -f " + pathVirtualDiskSD,
		"rm -f " + pathVirtualDiskEmmc,
	}

	// add mount command
	if getMountedDiskType() != diskType {
		if diskType == DiskTypeSDCard {
			cmd = append(cmd, "touch "+pathVirtualDiskSD)
		} else if diskType == DiskTypeEmmc {
			cmd = append(cmd, "touch "+pathVirtualDiskEmmc)
		}
	}

	cmd = append(cmd, scriptUsbDev+" restart")
	return cmd, nil
}

func ensureEmmcImageExists() error {
	if isFileExists(pathEmmcImage) {
		return nil
	}

	result, err := utils.RunShell(scriptMountEmmc + " start")
	if err != nil {
		return fmt.Errorf("failed to create emmc image: %v (exit code: %d, output: %s)",
			err, result.ExitCode, result.Stdout)
	}

	return nil
}

func executeWithHidLock(commands []string) error {
	h := hid.GetHid()
	h.Lock()
	h.CloseNoLock()
	defer func() {
		h.OpenNoLock()
		h.Unlock()
	}()

	result, err := utils.RunShell(commands...)
	if err != nil {
		return fmt.Errorf("command failed: %v (exit code: %d, stdout: %s, stderr: %s)",
			err, result.ExitCode, result.Stdout, result.Stderr)
	}

	return nil
}

func getMountedDiskType() string {
	if !isFileExists(pathVirtualDiskMounted) {
		return ""
	}

	content, err := os.ReadFile(pathVirtualDiskMounted)
	if err != nil {
		return ""
	}

	disk := strings.TrimSpace(string(content))

	switch disk {
	case "/dev/mmcblk1p1":
		return DiskTypeSDCard
	case pathEmmcImage:
		return DiskTypeEmmc
	default:
		return ""
	}
}

func isVirtualNetworkMounted() bool {
	return isFileExists(pathNetworkFlag)
}

func isSdCardPresent() bool {
	return isFileExists(pathSdCard)
}

func isEmmcImagePresent() bool {
	return isFileExists(pathEmmcImage)
}

func isFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
