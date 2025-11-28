package ui

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type GetLvtimeRsp struct {
	Format string `json:"format"`
}

type SetLvtimeReq struct {
	Format string `json:"format"`
}

type SetLvtimeRsp struct {
	Status string `json:"status"`
}

const (
	LvtimeFormat12H = "12h"
	LvtimeFormat24H = "24h"
)

func GetLvtime() (*GetLvtimeRsp, error) {
	rsp, err := Get("/lvtime/get")
	if err != nil {
		log.Errorf("[UI]  get lvtime failed error=%s", err.Error())
		return nil, err
	}

	var result GetLvtimeRsp
	if err := json.Unmarshal(rsp, &result); err != nil {
		log.Errorf("[UI]  unmarshal result failed error=%s", err.Error())
		return nil, err
	}

	return &result, nil
}

func SetLvtime(format string) error {
	if format != LvtimeFormat12H && format != LvtimeFormat24H {
		return fmt.Errorf("invalid format %s", format)
	}

	req := &SetLvtimeReq{
		Format: format,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Errorf("[UI] marshal failed error=%s", err.Error())
		return err
	}

	rsp, err := Post("/lvtime/set", reqBytes)
	if err != nil {
		log.Errorf("[UI] set lvtime failed error=%s", err.Error())
		return err
	}

	var result SetLvtimeRsp
	err = json.Unmarshal(rsp, &result)
	if err != nil {
		log.Errorf("[UI] unmarshal rsp failed error=%s", err.Error())
		return nil
	}

	if result.Status != "ok" {
		log.Debugf("[UI]  response status not 'ok' current status is '%s'", result.Status)
		return fmt.Errorf("invalid status %s", result.Status)
	}

	return nil
}
