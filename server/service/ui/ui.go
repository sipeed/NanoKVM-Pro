package ui

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type GetUIRsp struct {
	SleepS     int      `json:"sleep_s"`
	Sleep      bool     `json:"sleep"`
	Deprecated []string `json:"deprecated"`
	Direction  bool     `json:"dir"`
	Status     struct {
		Eth  bool `json:"eth"`
		Hdmi bool `json:"hdmi"`
		Hid  bool `json:"hid"`
		Wifi bool `json:"wifi"`
	} `json:"status"`
	UiType string `json:"ui_type"`
}

type SetUIReq struct {
	SleepS int `json:"sleep_s"`
}

type SetUIRsp struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}

func GetUI() (*GetUIRsp, error) {
	data, err := Get("/ui/get")
	if err != nil {
		return nil, err
	}

	var rsp GetUIRsp
	if err := json.Unmarshal(data, &rsp); err != nil {
		log.Errorf("failed to unmarshal response: %v", err)
		return nil, err
	}

	return &rsp, nil
}

func SetUI(sleep int) error {
	req := SetUIReq{
		SleepS: sleep,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		log.Errorf("failed to marshal request: %v", err)
		return err
	}

	data, err := Post("/ui/set", jsonReq)
	if err != nil {
		return err
	}

	var rsp SetUIRsp
	if err := json.Unmarshal(data, &rsp); err != nil {
		log.Errorf("failed to unmarshal response: %v", err)
		return err
	}

	if rsp.Status != "ok" {
		return fmt.Errorf("%s", rsp.Message)
	}

	return nil
}
