package config

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const ConfigurationFile = "/etc/kvm/server.yaml"

var (
	instance Config
	once     sync.Once
)

func GetInstance() *Config {
	once.Do(initialize)

	return &instance
}

func initialize() {
	if err := readByFile(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			create()
		}

		if err = readByDefault(); err != nil {
			log.Fatalf("Failed to read default configuration!")
		}

		log.Println("using default configuration")
	}

	if err := validate(); err != nil {
		log.Fatalf("Failed to validate configuration!")
	}

	if err := viper.Unmarshal(&instance); err != nil {
		log.Fatalf("Failed to parse configuration: %s", err)
	}

	checkDefaultValue()

	if instance.Authentication == "disable" {
		log.Println("NOTICE: Authentication is disabled! Please ensure your service is secure!")
	}

	log.Println("config loaded successfully")
}

func readByFile() error {
	viper.AddConfigPath(filepath.Dir(ConfigurationFile))
	viper.SetConfigName(filepath.Base(ConfigurationFile))
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func readByDefault() error {
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		log.Printf("failed to marshal default config: %s", err)
		return err
	}

	return viper.ReadConfig(bytes.NewBuffer(data))
}

// Create configuration file.
func create() {
	var (
		file *os.File
		data []byte
		err  error
	)

	dir := filepath.Dir(ConfigurationFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0o644)
	}

	file, err = os.OpenFile(ConfigurationFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Printf("open config failed: %s", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	if data, err = yaml.Marshal(defaultConfig); err != nil {
		log.Printf("failed to marshal default config: %s", err)
		return
	}

	if _, err = file.Write(data); err != nil {
		log.Printf("failed to save config: %s", err)
		return
	}

	if err = file.Sync(); err != nil {
		log.Printf("failed to sync config: %s", err)
		return
	}

	log.Printf("create file %s with default configuration", ConfigurationFile)
}

func validate() error {
	if viper.GetInt("port.http") > 0 && viper.GetInt("port.https") > 0 {
		return nil
	}

	log.Println("configuration invalid, resetting to defaults")

	create()

	return readByDefault()
}
