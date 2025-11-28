package vm

import (
	"NanoKVM-Server/config"
	"fmt"
	"os"
	"runtime"
	"strings"

	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/host"
	log "github.com/sirupsen/logrus"
)

func GetSystemArch() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("failed to get host info: %v", err)
	}

	arch := info.KernelArch
	if arch == "" {
		arch = runtime.GOARCH
	}

	return arch, nil
}

func (s *Service) GetInfo(c *gin.Context) {
	var rsp proto.Response

	data := &proto.GetInfoRsp{
		IPs:          getIPs(),
		Mdns:         getMdns(),
		Image:        GetSystemVersion(),
		Application:  getApplicationVersion(),
		DeviceKey:    getDeviceKey(),
		DeviceNumber: config.DeviceNumber(),
	}

	arch, _ := GetSystemArch()
	data.Arch = arch

	rsp.OkRspWithData(c, data)
	log.Debug("get vm information success")
}

func GetSystemVersion() string {
	content, err := os.ReadFile("/boot/ver")
	if err != nil {
		return ""
	}

	imagesVer := strings.ReplaceAll(string(content), "\n", "")
	part := strings.Split(imagesVer, "-")
	ver := part[len(part)-1]
	if ver == "" {
		ver = "v0.0.0"
	}
	return ver
}

func getIPs() (ips []proto.IP) {
	interfaces, err := GetInterfaceInfos()
	if err != nil {
		return
	}

	for _, iface := range interfaces {
		if iface.IP.To4() != nil {
			ips = append(ips, proto.IP{
				Name:    iface.Name,
				Addr:    iface.IP.String(),
				Version: "IPv4",
				Type:    iface.Type,
			})
		}
	}

	return
}

func getMdns() string {
	if !isAvahiRunning() {
		return ""
	}

	content, err := os.ReadFile("/etc/hostname")
	if err != nil {
		return ""
	}

	mdns := strings.ReplaceAll(string(content), "\n", "")
	return fmt.Sprintf("%s.local", mdns)
}

func getApplicationVersion() string {
	content, err := os.ReadFile("/kvmapp/version")
	if err != nil {
		return "1.0.0"
	}

	return strings.ReplaceAll(string(content), "\n", "")
}

func getDeviceKey() string {
	content, err := os.ReadFile("/device_key")
	if err != nil {
		return ""
	}

	return strings.ReplaceAll(string(content), "\n", "")
}

func (s *Service) GetHardware(c *gin.Context) {
	var rsp proto.Response

	conf := config.GetInstance()
	version := conf.Hardware.Version.String()

	rsp.OkRspWithData(c, &proto.GetHardwareRsp{
		Version: version,
	})
}
