package network

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/utils"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const WiFiScript = "/kvmcomm/scripts/wifi.sh"

func (s *Service) GetWifi(c *gin.Context) {
	var rsp proto.Response
	var data proto.GetWifiRsp

	interfaces, err := findWifiInterfaces()
	if err != nil || interfaces == nil || len(interfaces) == 0 {
		data.Supported = false
		rsp.OkRspWithData(c, &data)
		return
	}

	data.Supported = true
	data.APMode = isAPMode()

	details, err := utils.GetWifiDetails()
	if err != nil || details == nil || details.WPAState != "COMPLETED" {
		data.Connected = false
		rsp.OkRspWithData(c, &data)
		return
	}

	data.Connected = true

	data.Wifi = &proto.WiFi{
		SSID:      details.SSID,
		BSSID:     details.BSSID,
		Signal:    0,
		Frequency: 0,
		Security:  details.KeyMgmt,
	}

	rsp.OkRspWithData(c, &data)
	log.Debugf("get wifi successfully")
}

func (s *Service) ConnectWifiNoAuth(c *gin.Context) {
	var req proto.ConnectWifiReq
	var rsp proto.Response

	// Check Wi-Fi configuration mode
	if _, err := os.Stat("/tmp/wifi_config"); err != nil {
		rsp.ErrRsp(c, -1, "invalid mode")
		return
	}

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -2, "invalid parameters")
		return
	}

	_ = exec.Command(WiFiScript, "connect_stop").Run()

	if err := connect(req.Ssid, req.Password); err != nil {
		rsp.ErrRsp(c, -3, "set wifi failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("connect wifi successfully")
}

func (s *Service) ConnectWifi(c *gin.Context) {
	var req proto.ConnectWifiReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid parameters")
		return
	}

	_ = exec.Command(WiFiScript, "connect_stop").Run()

	if err := connect(req.Ssid, req.Password); err != nil {
		rsp.ErrRsp(c, -2, "set wifi failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("connect wifi successfully")
}

func (s *Service) DisconnectWifi(c *gin.Context) {
	var rsp proto.Response

	_ = exec.Command(WiFiScript, "connect_stop").Run()

	for range 5 {
		time.Sleep(2 * time.Second)
		if !isWifiConnected() {
			rsp.OkRsp(c)
			return
		}
	}

	rsp.ErrRsp(c, -1, "wifi disconnect failed")
	return
}

func connect(ssid string, password string) error {
	if ssid == "" || len(password) > 128 {
		return errors.New("invalid ssid or password")
	}

	_ = exec.Command(WiFiScript, "ap_stop").Run()
	_ = exec.Command(WiFiScript, "connect_stop").Run()

	var cmd *exec.Cmd
	if password != "" {
		cmd = exec.Command(WiFiScript, "connect_start", ssid, password)
	} else {
		cmd = exec.Command(WiFiScript, "connect_start", ssid)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		log.Errorf("failed to start wifi connect command: %v", err)
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(30 * time.Second):
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			log.Errorf("failed to kill process group: %v", err)
		}
		_ = exec.Command(WiFiScript, "connect_stop").Run()
		return errors.New("connect wifi timeout")
	case err := <-done:
		if err != nil {
			log.Errorf("connect wifi err: %v, out: %s", err, out.String())
			return err
		}
	}

	for range 5 {
		time.Sleep(2 * time.Second)
		if isWifiConnected() {
			return nil
		}
	}

	return errors.New("connect wifi timeout")
}

func isWifiConnected() bool {
	details, err := utils.GetWifiDetails()
	if err != nil {
		return false
	}

	return details.WPAState == "COMPLETED" && details.IPAddress != ""
}

func findWifiInterfaces() ([]string, error) {
	const netDir = "/sys/class/net"
	var interfaces []string

	files, err := os.ReadDir(netDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", netDir, err)
	}

	for _, file := range files {
		if file.IsDir() || (file.Type()&os.ModeSymlink) != 0 {
			wirelessPath := filepath.Join(netDir, file.Name(), "wireless")
			if _, err := os.Stat(wirelessPath); err == nil {
				interfaces = append(interfaces, file.Name())
			}
		}
	}

	return interfaces, nil
}

func isAPMode() bool {
	cmd := exec.Command("pgrep", "-x", "hostapd")
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			log.Debugf("pgrep process not found with exit code: %v", exitErr)
			return false
		}

		return false
	}

	return true
}
