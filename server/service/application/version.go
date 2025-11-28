package application

import (
	"NanoKVM-Server/proto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Latest struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Sha512  string `json:"sha512"`
	Size    uint   `json:"size"`
	Url     string `json:"url"`
}

func (s *Service) GetVersion(c *gin.Context) {
	var rsp proto.Response

	currentVersion := getCurrentVersion()

	latestVersion := currentVersion
	if latest, err := getLatest(); err == nil {
		latestVersion = latest.Version
	}

	rsp.OkRspWithData(c, &proto.GetVersionRsp{
		Current: currentVersion,
		Latest:  latestVersion,
	})
	log.Debugf("current version: %s, latest version: %s", currentVersion, latestVersion)
}

func getCurrentVersion() string {
	defaultVersion := "v1.0.0"

	versionFile := filepath.Join(AppDir, "version")
	content, err := os.ReadFile(versionFile)
	if err != nil {
		return defaultVersion
	}

	version := strings.ReplaceAll(string(content), "\n", "")
	if version == "" {
		return defaultVersion
	}

	return version
}

func getLatest() (*Latest, error) {
	baseURL := StableURL
	if isPreviewEnabled() {
		baseURL = PreviewURL
	}

	url := fmt.Sprintf("%s/nanokvm_pro_latest.json?now=%d", baseURL, time.Now().Unix())
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("failed to get latest version: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("server responded with status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var latest Latest
	if err := json.Unmarshal(body, &latest); err != nil {
		log.Errorf("failed to unmarshal response: %s", err)
		return nil, err
	}

	latest.Url = fmt.Sprintf("%s/%s", baseURL, latest.Name)

	log.Debugf("get application latest version: %s", latest.Version)
	return &latest, nil
}
