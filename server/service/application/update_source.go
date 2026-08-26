package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const maxUpdateSourceURLLength = 2048

var (
	updateSourceMu   sync.Mutex
	updateSourcePath = UpdateSourceConfigFile // replaced by tests only
)

// UpdateSourceConfig is deliberately shared with kvmcomm. URL is the common
// root, not an application or firmware sub-directory. The firmware updater
// derives URL/pro itself, while this service derives URL/nanokvm_pro_latest.json.
type UpdateSourceConfig struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
}

func defaultUpdateSourceConfig() UpdateSourceConfig {
	return UpdateSourceConfig{Enabled: false, URL: StableURL}
}

// GetUpdateSource returns the persisted source without probing the network.
func (s *Service) GetUpdateSource(c *gin.Context) {
	var rsp proto.Response
	cfg, err := loadUpdateSourceConfig()
	if err != nil {
		log.Errorf("load update source configuration: %v", err)
		rsp.ErrRsp(c, -1, "failed to load update source")
		return
	}

	rsp.OkRspWithData(c, updateSourceResponse(cfg))
}

// SetUpdateSource validates and atomically persists a custom source.
func (s *Service) SetUpdateSource(c *gin.Context) {
	var req proto.SetUpdateSourceReq
	var rsp proto.Response

	// Do not use proto.ParseFormRequest here: it logs the request in debug mode.
	// Update-source URLs are persisted configuration and must never contain
	// credentials; authenticated HTTP remains defensively handled by the client.
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	baseURL, err := normalizeUpdateSourceURL(req.URL)
	if err != nil {
		rsp.ErrRsp(c, -2, err.Error())
		return
	}

	cfg := defaultUpdateSourceConfig()
	if baseURL != "" && baseURL != StableURL {
		cfg.Enabled = true
		cfg.URL = baseURL
	}
	if err := saveUpdateSourceConfig(cfg); err != nil {
		log.Errorf("save update source configuration: %v", err)
		rsp.ErrRsp(c, -3, "failed to save update source")
		return
	}

	rsp.OkRspWithData(c, updateSourceResponse(cfg))
}

// ResetUpdateSource disables the custom source and restores the official root.
func (s *Service) ResetUpdateSource(c *gin.Context) {
	var rsp proto.Response
	cfg := defaultUpdateSourceConfig()
	if err := saveUpdateSourceConfig(cfg); err != nil {
		log.Errorf("reset update source configuration: %v", err)
		rsp.ErrRsp(c, -1, "failed to reset update source")
		return
	}
	rsp.OkRspWithData(c, updateSourceResponse(cfg))
}

func updateSourceResponse(cfg UpdateSourceConfig) *proto.GetUpdateSourceRsp {
	return &proto.GetUpdateSourceRsp{
		URL:        cfg.URL,
		IsOfficial: !cfg.Enabled,
		Enabled:    cfg.Enabled,
	}
}

func resolveApplicationUpdateBaseURL() (string, error) {
	cfg, err := loadUpdateSourceConfig()
	if err != nil {
		return "", err
	}
	if cfg.Enabled {
		return cfg.URL, nil
	}
	if isPreviewEnabled() {
		return PreviewURL, nil
	}
	return StableURL, nil
}

// normalizeUpdateSourceURL accepts only a directory-like HTTP(S) URL so all
// derived manifest and package paths stay under the selected update root.
func normalizeUpdateSourceURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}
	if len(raw) > maxUpdateSourceURLLength {
		return "", errors.New("update server URL is too long")
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", errors.New("invalid update server URL")
	}
	if parsed.User != nil {
		return "", errors.New("update server URL must not contain credentials")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", errors.New("update server URL must not contain a query or fragment")
	}
	if strings.HasSuffix(strings.ToLower(strings.TrimRight(parsed.Path, "/")), ".json") {
		return "", errors.New("enter the update server directory, not a manifest URL")
	}

	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return parsed.String(), nil
}

func joinUpdateURL(base, suffix string) string {
	parsed, err := url.Parse(base)
	if err != nil {
		return ""
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/" + strings.TrimLeft(suffix, "/")
	return parsed.String()
}

func loadUpdateSourceConfig() (UpdateSourceConfig, error) {
	updateSourceMu.Lock()
	defer updateSourceMu.Unlock()
	return loadUpdateSourceConfigFromPath(updateSourcePath)
}

func loadUpdateSourceConfigFromPath(path string) (UpdateSourceConfig, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return defaultUpdateSourceConfig(), nil
	}
	if err != nil {
		return UpdateSourceConfig{}, err
	}

	var cfg UpdateSourceConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return UpdateSourceConfig{}, fmt.Errorf("decode update source: %w", err)
	}
	if !cfg.Enabled {
		return defaultUpdateSourceConfig(), nil
	}
	normalized, err := normalizeUpdateSourceURL(cfg.URL)
	if err != nil || normalized == "" {
		if err == nil {
			err = errors.New("custom update source URL is required")
		}
		return UpdateSourceConfig{}, err
	}
	cfg.URL = normalized
	return cfg, nil
}

func saveUpdateSourceConfig(cfg UpdateSourceConfig) error {
	updateSourceMu.Lock()
	defer updateSourceMu.Unlock()
	return saveUpdateSourceConfigToPath(updateSourcePath, cfg)
}

// saveUpdateSourceConfigToPath writes via a private temporary file and rename,
// preventing readers from observing a partially written configuration.
func saveUpdateSourceConfigToPath(path string, cfg UpdateSourceConfig) error {
	if !cfg.Enabled {
		cfg = defaultUpdateSourceConfig()
	}
	if cfg.Enabled {
		normalized, err := normalizeUpdateSourceURL(cfg.URL)
		if err != nil || normalized == "" {
			if err == nil {
				err = errors.New("custom update source URL is required")
			}
			return err
		}
		cfg.URL = normalized
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create update source directory: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode update source: %w", err)
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(dir, ".update-source.*")
	if err != nil {
		return fmt.Errorf("create temporary update source: %w", err)
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("set update source permissions: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write update source: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("sync update source: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close update source: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replace update source: %w", err)
	}
	if directory, err := os.Open(dir); err == nil {
		_ = directory.Sync()
		_ = directory.Close()
	}
	return nil
}
