package application

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"NanoKVM-Server/proto"
	"NanoKVM-Server/utils"
)

const maxTries = 3

var (
	updateMutex sync.Mutex
	isUpdating  bool
)

func (s *Service) Update(c *gin.Context) {
	var rsp proto.Response

	updateMutex.Lock()
	if isUpdating {
		updateMutex.Unlock()
		rsp.ErrRsp(c, -1, "update already in progress")
		return
	}
	isUpdating = true
	updateMutex.Unlock()

	defer func() {
		updateMutex.Lock()
		isUpdating = false
		updateMutex.Unlock()
	}()

	if err := update(); err != nil {
		rsp.ErrRsp(c, -2, "failed to update service")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("update service success")
}

func update() error {
	_ = os.RemoveAll(TempDir)
	if err := os.MkdirAll(TempDir, 0o755); err != nil {
		log.Errorf("failed to create dir %s: %v", TempDir, err)
		return err
	}

	latest, err := getLatest()
	if err != nil {
		log.Errorf("failed to get latest service: %v", err)
		return err
	}

	// download
	_ = sendMessage("download", 0)
	tarFile := filepath.Join(TempDir, latest.Name)

	if err := download(latest.Url, tarFile); err != nil {
		log.Errorf("download app failed: %s", err)
		return err
	}

	// check tar sha512
	_ = sendMessage("install", 0)
	if err := checksum(tarFile, latest.Sha512); err != nil {
		log.Errorf("check sha512 failed: %s", err)
		return err
	}

	// decompress
	dir, err := utils.UnTarGz(tarFile, TempDir)
	if err != nil {
		log.Errorf("failed to decompress application: %s", err)
		return err
	}

	// install
	err = install(dir, latest.Version)
	if err != nil {
		log.Errorf("failed to install application: %s", err)
		return err
	}

	_ = sendMessage("restart", 0)

	log.Debugf("update service success")
	return nil
}

func download(url string, target string) (err error) {
	for i := range maxTries {
		if i > 0 {
			log.Debugf("attempt to download application #%d/%d", i+1, maxTries)
			time.Sleep(time.Second * 3)
		}

		var req *http.Request
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			log.Errorf("new request err: %s", err)
			continue
		}

		err = utils.DownloadWithProgress(req, target, progressHandler)
		if err != nil {
			log.Errorf("downloading latest application failed: %s", err)
			continue
		}

		log.Debugf("installation package downloaded to %s", target)
		return
	}

	log.Errorf("downloading latest application failed")
	return
}

func install(dir string, version string) error {
	// check sha512
	for _, appName := range appNames {
		jsonPath := filepath.Join(dir, fmt.Sprintf("%s_%s.json", appName, version))
		debPath := filepath.Join(dir, fmt.Sprintf("%s_%s_arm64.deb", appName, version))

		fileInfo, err := getFileInfo(jsonPath)
		if err != nil {
			log.Errorf("failed to get %s info: %s", appName, err)
			return err
		}

		if err := checksum(debPath, fileInfo.SHA512); err != nil {
			log.Errorf("failed to checksum %s: %s", appName, err)
			return err
		}
	}

	var failedApps []string

	// install
	for i, appName := range appNames {
		_ = sendMessage("install", i*30+10)

		installSuccess := false
		for try := range maxTries {
			if try > 0 {
				log.Debugf("try installing %s again: %d/%d", appName, try, maxTries)
				time.Sleep(time.Duration(try) * time.Second)
			}

			fileName := fmt.Sprintf("%s_%s_arm64.deb", appName, version)
			filePath := filepath.Join(dir, fileName)

			cmd := fmt.Sprintf("DEBIAN_FRONTEND=noninteractive dpkg -i %s", filePath)
			output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
			if err != nil {
				log.Errorf("failed to install %s: %s", filePath, string(output))
				continue
			}

			installSuccess = true
			log.Infof("install %s success", appName)
			break
		}

		if !installSuccess {
			failedApps = append(failedApps, appName)
		}
	}

	if len(failedApps) > 0 {
		log.Errorf("failed to install applications: %v", failedApps)
	}

	_ = sendMessage("install", 100)
	return nil
}
