package tailscale

import (
	"NanoKVM-Server/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shirou/gopsutil/v4/host"
	log "github.com/sirupsen/logrus"
)

const (
	OriginalURL = "https://pkgs.tailscale.com/stable/"

	TempDir = "/root/.tailscale"

	TailscaleFile  = "/usr/bin/tailscale"
	TailscaledFile = "/usr/sbin/tailscaled"
	EnvFile        = "/etc/default/tailscaled"

	Tailscaled     = "tailscaled.service"
	TailscaledConf = "/etc/systemd/system/tailscaled.service"
)

func GetTailscaleName() string {
	ARCH := "arm64"
	VERSION := "latest"
	info, err := host.Info()
	if err != nil {
		return fmt.Sprintf("tailscale_%s_%s.tgz", ARCH, VERSION)
	}

	ARCH = info.KernelArch
	if ARCH == "aarch64" {
		ARCH = "arm64"
	}
	return fmt.Sprintf("tailscale_%s_%s.tgz", VERSION, ARCH)
}

func isInstalled() bool {
	for _, file := range []string{TailscaleFile, TailscaledFile} {
		if _, err := os.Stat(file); err != nil {
			return false
		}
	}

	return true
}

func install() error {
	_ = os.MkdirAll(TempDir, 0o755)
	defer func() {
		_ = os.RemoveAll(TempDir)
	}()

	// download
	tarFile := filepath.Join(TempDir, GetTailscaleName())
	if err := download(tarFile); err != nil {
		log.Errorf("failed to download tailscale: %s", err)
		return err
	}

	// decompress
	dir, err := utils.UnTarGz(tarFile, TempDir)
	if err != nil {
		log.Errorf("failed to decompress tailscale: %s", err)
		return err
	}

	// move
	if err := os.Rename(filepath.Join(dir, "tailscale"), TailscaleFile); err != nil {
		log.Errorf("failed to rename tailscale file: %s", err)
		return err
	}

	if err := os.Rename(filepath.Join(dir, "tailscaled"), TailscaledFile); err != nil {
		log.Errorf("failed to rename tailscaled file: %s", err)
		return err
	}

	// service config
	serviceConf := filepath.Join(dir, "systemd", Tailscaled)
	info, err := os.Stat(serviceConf)
	if err != nil {
		log.Warnf("failed to get file info for %s: %s", serviceConf, err)
		return err
	}

	if err := os.Chmod(serviceConf, info.Mode()|0111); err != nil {
		log.Printf("add service permission failed: %s", err)
		return err
	}

	if err := os.Rename(serviceConf, TailscaledConf); err != nil {
		log.Printf("failed to move %s to %s: %s", serviceConf, TailscaledConf, err)
		return err
	}

	// create required directories
	for _, directory := range []string{"/var/lib/tailscale", "/run/tailscale"} {
		if err := os.MkdirAll(directory, 0755); err != nil {
			log.Errorf("failed to create directory %s: %s", directory, err)
		}
		if err := os.Chmod(directory, 0755); err != nil {
			log.Errorf("failed to chmod directory %s: %s", directory, err)
		}
	}

	// create environment file
	envContent := `PORT=41641
FLAGS=""
`
	err = os.WriteFile(EnvFile, []byte(envContent), 0644)
	if err != nil {
		log.Errorf("failed to create tailscaled environment file: %s", err)
		return err
	}

	// start service
	if err := utils.DaemonReload(); err != nil {
		log.Errorf("failed to reload systemd: %s", err)
		return err
	}

	if err := utils.StartService(Tailscaled, true); err != nil {
		log.Errorf("failed to start kvmadmin: %s", err)
		return err
	}

	log.Debugf("install tailscale successfully")
	return nil
}

func uninstall() error {
	if !isInstalled() {
		log.Debugf("tailscale not install")
		return nil
	}

	// stop service
	if err := utils.StopService(Tailscaled, true); err != nil {
		log.Errorf("failed to stop kvmadmin: %s", err)
		return err
	}

	// delete files
	for _, file := range []string{TailscaleFile, TailscaledFile, EnvFile, TailscaledConf} {
		if err := os.Remove(file); err != nil {
			log.Errorf("failed to remove file %s: %s", file, err)
			return err
		}
	}

	// reload
	if err := utils.DaemonReload(); err != nil {
		log.Warnf("failed to reload systemd: %s", err)
	}

	log.Debugf("uninstall tailscale successfully")
	return nil
}

func download(target string) error {
	url, err := getDownloadURL()
	if err != nil {
		log.Errorf("failed to get Tailscale download url: %s", err)
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("failed to download Tailscale: %s", err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(target)
	if err != nil {
		log.Errorf("failed to create file: %s", err)
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Errorf("failed to copy response body to file: %s", err)
		return err
	}

	log.Debugf("download Tailscale successfully")
	return nil
}

func getDownloadURL() (string, error) {
	resp, err := (&http.Client{}).Get(OriginalURL + GetTailscaleName())
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Request.URL.String(), nil
}
