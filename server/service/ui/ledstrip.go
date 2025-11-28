package ui

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type LedStrip struct {
	Brightness int  `json:"brightness"`
	Hor        int  `json:"hor"`
	On         bool `json:"on"`
	Ver        int  `json:"ver"`
}

type LedStripResponse struct {
	Status string `json:"status"`
}

func SetLedStrip(config map[string]any) *LedStripResponse {
	var rsp LedStripResponse

	data, err := json.Marshal(config)
	if err != nil {
		log.Errorf("[UI] config params marshal failed error=%s", err.Error())
		return nil
	}

	r, err := Post("/ledstrip/set", data)
	if err != nil {
		log.Errorf("[UI] calling set api failed, error=%s", err.Error())
	}

	err = json.Unmarshal(r, &rsp)
	if err != nil {
		log.Errorf("[UI]  %s", err.Error())
		return nil
	}

	if rsp.Status != "ok" {
		log.Debugf("[UI]  response status not 'ok' current status is '%s'", rsp.Status)
	}

	return &rsp
}

func GetLedStrip() (*LedStrip, error) {
	var result LedStrip

	r, err := Get("/ledstrip/get")
	if err != nil {
		log.Errorf("[UI]  get LED config failed error=%s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(r, &result)
	if err != nil {
		log.Errorf("[UI]  unmarshal result failed error=%s", err.Error())
		return nil, err
	}

	return &result, nil
}
