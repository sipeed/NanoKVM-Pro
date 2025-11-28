package application

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type FileInfo struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	SHA512  string `json:"sha512"`
	Size    int    `json:"size"`
}

func getFileInfo(filePath string) (*FileInfo, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Errorf("failed to read file %s: %v", filePath, err)
		return nil, err
	}

	var info FileInfo
	err = json.Unmarshal(jsonData, &info)
	if err != nil {
		log.Errorf("failed to unmarshal file %s: %v", filePath, err)
		return nil, err
	}

	return &info, nil
}

func checksum(filePath string, expectedHash string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("failed to open file %s: %v", filePath, err)
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	hasher := sha512.New()

	_, err = io.Copy(hasher, file)
	if err != nil {
		log.Errorf("failed to copy file contents to hasher: %v", err)
		return err
	}

	hash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	if hash != expectedHash {
		log.Errorf("invalid sha512 %s", hash)
		return fmt.Errorf("invalid sha512 %s", hash)
	}

	log.Debugf("%s checksum success", filePath)
	return nil
}
