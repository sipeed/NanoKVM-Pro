package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type ProgressWriter struct {
	Total      int64
	Downloaded int64
	Callback   func(int64, int64)
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Downloaded += int64(n)
	pw.Callback(pw.Total, pw.Downloaded)
	return n, nil
}

func Download(req *http.Request, target string) error {
	log.Debugf("downloading %s to %s", req.URL.String(), target)
	err := os.MkdirAll(filepath.Dir(target), 0o755)
	if err != nil {
		log.Errorf("create dir %s err: %s", filepath.Dir(target), err)
		return err
	}
	out, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		log.Errorf("cannot create file '%s', error: %s", target, err)
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Errorf("request error: %s", err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("request failed, status code: %d", resp.StatusCode)
		return errors.New("update website is inaccessible right now")
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Errorf("download file to %s err: %s", target, err)
		return err
	}

	return nil
}

func DownloadWithProgress(req *http.Request, target string, callback func(int64, int64)) error {
	log.Debugf("downloading %s to %s", req.URL.String(), target)

	err := os.MkdirAll(filepath.Dir(target), 0o755)
	if err != nil {
		log.Errorf("create dir %s err: %s", filepath.Dir(target), err)
		return err
	}

	out, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		log.Errorf("cannot create file '%s', error: %s", target, err)
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Errorf("request error: %s", err)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("request failed, status code: %d", resp.StatusCode)
		return errors.New("update website is inaccessible right now")
	}

	totalSize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		totalSize = -1
	}

	progressWriter := &ProgressWriter{
		Total:      totalSize,
		Downloaded: 0,
		Callback:   callback,
	}

	reader := io.TeeReader(resp.Body, progressWriter)

	_, err = io.Copy(out, reader)
	if err != nil {
		log.Errorf("download file to %s err: %s", target, err)
		return err
	}

	return nil
}
