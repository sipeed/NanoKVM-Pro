package ai

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	ConfigFile   = "/etc/kvm/ai.json"
	apiKeyPrefix = "nap_mcp_"
	apiKeyBytes  = 32
)

var (
	configMu       sync.Mutex
	configFilePath = ConfigFile
)

type Config struct {
	MCPEnabled bool   `json:"mcpEnabled"`
	MCPAPIKey  string `json:"mcpApiKey"`
}

func loadConfig() (Config, error) {
	configMu.Lock()
	defer configMu.Unlock()
	return loadConfigFromPath(configFilePath)
}

func saveConfig(cfg Config) error {
	configMu.Lock()
	defer configMu.Unlock()
	return saveConfigToPath(configFilePath, cfg)
}

func updateConfig(update func(Config) (Config, error)) (Config, error) {
	configMu.Lock()
	defer configMu.Unlock()

	cfg, err := loadConfigFromPath(configFilePath)
	if err != nil {
		return Config{}, err
	}

	updated, err := update(cfg)
	if err != nil {
		return Config{}, err
	}

	if err := saveConfigToPath(configFilePath, updated); err != nil {
		return Config{}, err
	}
	return updated, nil
}

func loadConfigFromPath(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func saveConfigToPath(path string, cfg Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(dir, ".ai.json.*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	return os.Chmod(path, 0o600)
}

func ensureAPIKey(cfg Config) (Config, error) {
	if cfg.MCPAPIKey != "" {
		return cfg, nil
	}

	key, err := generateAPIKey()
	if err != nil {
		return Config{}, err
	}
	cfg.MCPAPIKey = key
	return cfg, nil
}

func regenerateAPIKey(cfg Config) (Config, error) {
	key, err := generateAPIKey()
	if err != nil {
		return Config{}, err
	}
	cfg.MCPAPIKey = key
	return cfg, nil
}

func generateAPIKey() (string, error) {
	raw := make([]byte, apiKeyBytes)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate API key: %w", err)
	}

	return apiKeyPrefix + base64.RawURLEncoding.EncodeToString(raw), nil
}
