package application

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var updateHTTPClient = &http.Client{Timeout: 15 * time.Minute}

func readUpdateManifest(rawURL string, maxBytes int64) ([]byte, error) {
	// Limit both Content-Length and streamed bytes because mirrors may omit the
	// header or provide an intentionally misleading value.
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := updateHTTPClient.Do(req)
	if err != nil {
		return nil, errors.New("update server is inaccessible")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update server returned status %d", resp.StatusCode)
	}
	if resp.ContentLength > maxBytes {
		return nil, fmt.Errorf("update manifest exceeds %d bytes", maxBytes)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > maxBytes {
		return nil, fmt.Errorf("update manifest exceeds %d bytes", maxBytes)
	}
	return body, nil
}

func downloadUpdatePackage(rawURL, target string, maxBytes int64, beforeWrite func(int64) error) error {
	// Write with exclusive creation and remove partial output on every failure.
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	resp, err := updateHTTPClient.Do(req)
	if err != nil {
		return errors.New("update server is inaccessible")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update server returned status %d", resp.StatusCode)
	}
	if resp.ContentLength > maxBytes {
		return fmt.Errorf("update package exceeds %d bytes", maxBytes)
	}
	if resp.ContentLength > 0 && beforeWrite != nil {
		if err := beforeWrite(resp.ContentLength); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return err
	}
	output, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	success := false
	defer func() {
		_ = output.Close()
		if !success {
			_ = os.Remove(target)
		}
	}()

	progress := &updateProgressWriter{total: resp.ContentLength}
	written, err := io.Copy(io.MultiWriter(output, progress), io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return err
	}
	if written > maxBytes {
		return fmt.Errorf("update package exceeds %d bytes", maxBytes)
	}
	if err := output.Sync(); err != nil {
		return err
	}
	if err := output.Close(); err != nil {
		return err
	}
	success = true
	return nil
}

type updateProgressWriter struct {
	total   int64
	written int64
}

// Write reports download progress without changing the bytes written to disk.
func (w *updateProgressWriter) Write(data []byte) (int, error) {
	w.written += int64(len(data))
	if w.total > 0 {
		progress := int(w.written * 100 / w.total)
		if progress > 100 {
			progress = 100
		}
		_ = sendMessage("download", progress)
	}
	return len(data), nil
}
