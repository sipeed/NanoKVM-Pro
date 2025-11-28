package vm

import (
	"NanoKVM-Server/proto"
	"NanoKVM-Server/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	EdidDir            = "/kvmcomm/edid"
	LT6911Edid         = "/proc/lt6911_info/edid"
	LT6911EdidSnapshot = "/proc/lt6911_info/edid_snapshot"

	CustomEdidDir  = "/etc/kvm/edid"
	CustomEdidFlag = "/etc/kvm/edid/edid_flag"
)

var EDIDMap = map[byte]string{
	0x12: "E18-4K30FPS",
	0x30: "E48-4K39FPS",
	0x36: "E54-1080P60FPS",
	0x38: "E56-2K60FPS",
	0x3a: "E58-4K16-10",
	0x3f: "E63-Ultrawide",
}

func (s *Service) GetEdid(ctx *gin.Context) {
	var resp proto.Response
	content, err := os.ReadFile(LT6911EdidSnapshot)
	if err != nil {
		resp.ErrRsp(ctx, -1, "get edid failed")
		return
	}

	edid, ok := EDIDMap[content[12]]
	if !ok {
		// custom EDID
		if flag, err := os.ReadFile(CustomEdidFlag); err == nil {
			edid = strings.TrimSpace(string(flag))
		}
	}

	resp.OkRspWithData(ctx, gin.H{
		"edid": edid,
	})
}

func (s *Service) SwitchEdid(c *gin.Context) {
	var req proto.SwitchEdidReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	if req.Edid == "" {
		rsp.ErrRsp(c, -2, "invalid EDID")
		return
	}

	srcPath := filepath.Join(EdidDir, req.Edid+".bin")
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		// custom EDID
		srcPath = filepath.Join(CustomEdidDir, req.Edid)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			log.Debugf("unknown edid: %s", req.Edid)
			rsp.ErrRsp(c, -3, "invalid EDID")
			return
		}

		_ = os.WriteFile(CustomEdidFlag, []byte(req.Edid), 0644)
	}

	if err := copyFile(srcPath, LT6911Edid); err != nil {
		log.Errorf("failed to switch EDID %s: %s", req.Edid, err)
		rsp.ErrRsp(c, -4, "failed to switch EDID")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("switch edid %s", req.Edid)
}

func (s *Service) GetCustomEdidList(c *gin.Context) {
	var rsp proto.Response

	var files []string

	err := filepath.Walk(CustomEdidDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isBin(info.Name()) {
			files = append(files, info.Name())
		}

		return nil
	})
	if err != nil {
		rsp.ErrRsp(c, -1, "get EDID failed")
		return
	}

	rsp.OkRspWithData(c, &proto.GetCustomEdidListRsp{
		EdidList: files,
	})
}

func (s *Service) DeleteEdid(c *gin.Context) {
	var req proto.DeleteEdidReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	file := filepath.Join(CustomEdidDir, req.Edid)

	if err := os.Remove(file); err != nil {
		log.Errorf("delete edid %s failed: %s", req.Edid, err)
		rsp.ErrRsp(c, -2, "delete failed")
		return
	}

	rsp.OkRsp(c)
	log.Debugf("delete edid %s success", req.Edid)
}

func (s *Service) UploadEdid(c *gin.Context) {
	var rsp proto.Response

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		rsp.ErrRsp(c, -1, "bad request")
		return
	}
	defer file.Close()

	if _, err = os.Stat(CustomEdidDir); err != nil {
		_ = os.MkdirAll(CustomEdidDir, 0o755)
	}

	target := fmt.Sprintf("%s/%s", CustomEdidDir, header.Filename)
	dst, err := os.Create(target)
	if err != nil {
		rsp.ErrRsp(c, -2, "create file failed")
		return
	}
	defer dst.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(dst, file, buf)
	if err != nil {
		rsp.ErrRsp(c, -3, "save failed")
		return
	}

	_ = utils.EnsurePermission(target, 0o644)

	data := &proto.UploadEdidRsp{
		File: header.Filename,
	}
	rsp.OkRspWithData(c, data)
	log.Debugf("upload edid file: %s", header.Filename)
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

func isBin(name string) bool {
	nameLower := strings.ToLower(name)
	if strings.HasSuffix(nameLower, ".bin") {
		return true
	}

	return false
}
