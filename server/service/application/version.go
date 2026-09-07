package application

import (
	"NanoKVM-Server/proto"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Latest is the bounded manifest data needed to locate and verify an update.
type Latest struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Sha512  string `json:"sha512"`
	Size    uint64 `json:"size"`
	Url     string `json:"url"`
}

const maxUpdateManifestBytes = 64 << 10

var latestApplicationNamePattern = regexp.MustCompile(`^nanokvm_pro_([A-Za-z0-9][A-Za-z0-9._+-]{0,127})\.tar\.gz$`)

// GetVersion reports the installed version and the latest version from the
// currently selected source.
func (s *Service) GetVersion(c *gin.Context) {
	var rsp proto.Response

	currentVersion := getCurrentVersion()

	latestVersion := currentVersion
	if latest, err := getLatest(); err == nil {
		latestVersion = latest.Version
	} else {
		cfg, cfgErr := loadUpdateSourceConfig()
		if cfgErr != nil || cfg.Enabled {
			log.Errorf("failed to query custom update source: %v", err)
			rsp.ErrRsp(c, -1, "failed to query update source")
			return
		}
	}

	rsp.OkRspWithData(c, &proto.GetVersionRsp{
		Current: currentVersion,
		Latest:  latestVersion,
	})
	log.Debugf("current version: %s, latest version: %s", currentVersion, latestVersion)
}

// getCurrentVersion reads the version exposed by the installed application.
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
	// The configured source is resolved for each check so resetting the source
	// takes effect without restarting the service.
	baseURL, err := resolveApplicationUpdateBaseURL()
	if err != nil {
		return nil, err
	}

	manifestURL := joinUpdateURL(baseURL, "nanokvm_pro_latest.json") + "?now=" + strconv.FormatInt(time.Now().Unix(), 10)
	body, err := readUpdateManifest(manifestURL, maxUpdateManifestBytes)
	if err != nil {
		log.Errorf("failed to read latest version manifest: %v", err)
		return nil, err
	}

	var latest Latest
	if err := json.Unmarshal(body, &latest); err != nil {
		log.Errorf("failed to unmarshal response: %s", err)
		return nil, err
	}
	if err := validateLatest(&latest); err != nil {
		return nil, err
	}

	latest.Url = joinUpdateURL(baseURL, latest.Name)

	log.Debugf("get application latest version: %s", latest.Version)
	return &latest, nil
}

func validateLatest(latest *Latest) error {
	// Validate names and hashes before constructing a download URL from a
	// manifest supplied by either the official or a custom source.
	if latest == nil || !validArtifactVersion(latest.Version) {
		return errors.New("invalid application update version")
	}
	matches := latestApplicationNamePattern.FindStringSubmatch(latest.Name)
	if len(matches) != 2 || matches[1] != latest.Version {
		return errors.New("invalid application update package name")
	}
	if err := validateSHA512String(latest.Sha512); err != nil {
		return errors.New("invalid application update checksum")
	}
	return nil
}
