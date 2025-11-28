package vm

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *Service) SetTimeZone(c *gin.Context) {
	var req proto.SetTimeZoneReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if _, err := time.LoadLocation(req.Timezone); err != nil {
		rsp.ErrRsp(c, -2, "invalid timezone")
		return
	}

	cmd := exec.Command("timedatectl", "set-timezone", req.Timezone)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("set timezone failed, output: %s", string(output))
		rsp.ErrRsp(c, -2, "set timezone failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("set timezone %s success", req.Timezone)
}

func (s *Service) GetTimeZone(c *gin.Context) {
	var rsp proto.Response

	timezone, err := getIANATimezone()
	if err != nil {
		rsp.ErrRsp(c, -1, "failed to get timezone")
		return
	}

	rsp.OkRspWithData(c, &proto.GetTimeZoneRsp{
		Timezone: timezone,
	})
	log.Debugf("get timezone %s success", timezone)
}

func getIANATimezone() (string, error) {
	tzPath, err := os.Readlink("/etc/localtime")
	if err != nil {
		log.Errorf("Failed to read timezone from /etc/localtime: %s", err)
		return "", err
	}

	zoneinfo := "/usr/share/zoneinfo/"
	if i := strings.LastIndex(tzPath, zoneinfo); i != -1 {
		return tzPath[i+len(zoneinfo):], nil
	}

	if !filepath.IsAbs(tzPath) {
		tzPath = filepath.Join("/etc", tzPath)
	}

	relPath, err := filepath.Rel(zoneinfo, tzPath)
	if err != nil {
		log.Errorf("Failed to find timezone from %s: %s", tzPath, err)
		return "", err
	}

	return relPath, nil
}

func (s *Service) GetTime(c *gin.Context) {
	messageChan := make(chan int64)
	defer close(messageChan)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				messageChan <- t.UnixMilli()
			case <-c.Request.Context().Done():
				return
			}
		}
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-messageChan; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}

func (s *Service) GetTimeStatus(c *gin.Context) {
	var rsp proto.Response

	status, err := utils.GetChronySyncStatus()
	if err != nil {
		log.Errorf("get chrony sync status failed: %s", err)
		rsp.ErrRsp(c, -1, "get chrony sync status failed")
		return
	}

	rsp.OkRspWithData(c, &proto.GetTimeStatusRsp{
		IsSynchronized: status.IsSynchronized,
		LastSyncTime:   status.LastSyncTime,
	})
	log.Debugf("get chrony sync status success")
}

func (s *Service) SyncTime(c *gin.Context) {
	var rsp proto.Response

	err := utils.ChronycBurst()
	if err != nil {
		log.Errorf("chronyc makestep failed: %s", err)
		rsp.ErrRsp(c, -1, "sync time failed")
		return
	}

	time.Sleep(2 * time.Second)

	rsp.OkRsp(c)
	log.Debugf("sync time success")
}
