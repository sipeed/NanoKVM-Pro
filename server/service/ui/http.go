package ui

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	BaseURL = "http://127.0.0.1:65501/api"
)

var Online = new(atomic.Value)

func Get(url string) ([]byte, error) {
	resp, err := http.Get(BaseURL + url)
	if err != nil {
		log.Errorf("failed to request %s: %v", url, err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("wrong stauts code: %d", resp.StatusCode)
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response: %v", err)
		return nil, err
	}

	return body, nil
}

func Post(url string, data []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest("POST", BaseURL+url, bytes.NewBuffer(data))
	if err != nil {
		log.Errorf("failed to create request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("failed to send request: %v", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Errorf("wrong status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response: %v", err)
		return nil, err
	}

	return body, nil
}
