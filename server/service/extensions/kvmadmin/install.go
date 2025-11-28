package kvmadmin

import (
	"NanoKVM-Server/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	AppDir    = "/kvmadmin"
	TempDir   = "/root/.kvmadmin"
	ConfigDir = "/etc/kvmadmin"

	Kvmadmin     = "kvmadmin.service"
	KvmadminConf = "/etc/systemd/system/kvmadmin.service"

	DownloadURL = "https://cdn.sipeed.com/nanokvm/resources/kvmadmin.tar.gz"
)

func isInstalled() bool {
	_, err := os.Stat(filepath.Join(AppDir, "NanoKVM-Admin"))
	return err == nil
}

func install() error {
	if isInstalled() {
		log.Debugf("kvmadmin is already installed")
		return nil
	}

	_ = os.MkdirAll(TempDir, 0o755)
	defer func() {
		_ = os.RemoveAll(TempDir)
	}()

	// download
	tarFile := filepath.Join(TempDir, "kvmadmin.tar.gz")
	if err := download(tarFile); err != nil {
		log.Errorf("failed to download kvmadmin: %s", err)
		return err
	}

	// decompress
	dir, err := utils.UnTarGz(tarFile, TempDir)
	if err != nil {
		log.Errorf("failed to decompress kvmadmin: %s", err)
		return err
	}

	// create service config
	serviceConf := filepath.Join(dir, Kvmadmin)
	info, err := os.Stat(serviceConf)
	if err != nil {
		log.Errorf("failed to get file info for %s: %s", serviceConf, err)
		return err
	}

	if err := os.Chmod(serviceConf, info.Mode()|0111); err != nil {
		log.Printf("add service permission failed: %s", err)
		return err
	}

	if err := os.Rename(serviceConf, KvmadminConf); err != nil {
		log.Printf("failed to move %s to %s: %s", serviceConf, KvmadminConf, err)
		return err
	}

	// move
	if err := utils.MoveFilesRecursively(dir, AppDir); err != nil {
		log.Errorf("failed to move kvmadmin: %s", err)
		return err
	}

	// start service
	if err := utils.DaemonReload(); err != nil {
		log.Errorf("failed to reload systemd: %s", err)
		return err
	}

	if err := utils.StartService(Kvmadmin, true); err != nil {
		log.Errorf("failed to start kvmadmin: %s", err)
		return err
	}

	log.Debugf("kvmadmin is installed successfully")
	return nil
}

func uninstall() error {
	if !isInstalled() {
		log.Debugf("kvmadmin is not installed")
		return nil
	}

	// stop service
	if err := utils.StopService(Kvmadmin, true); err != nil {
		log.Warnf("failed to stop kvmadmin: %s", err)
	}

	// delete files
	for _, directory := range []string{AppDir, ConfigDir} {
		if err := os.RemoveAll(directory); err != nil {
			log.Errorf("failed to remove %s: %s", directory, err)
			return err
		}
	}

	if err := os.Remove(KvmadminConf); err != nil {
		log.Warnf("failed to remove %s: %s", KvmadminConf, err)
	}

	// reload
	if err := utils.DaemonReload(); err != nil {
		log.Warnf("failed to reload systemd: %s", err)
	}

	log.Debugf("uninstall kvmadmin successfully")
	return nil
}

func download(target string) error {
	resp, err := http.Get(DownloadURL)
	if err != nil {
		log.Errorf("failed to download kvmadmin: %s", err)
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

	log.Debugf("download kvmadmin successfully")
	return nil
}
