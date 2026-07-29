package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

const minutesPerDay = 24 * 60

type ScreenOff struct {
	Enabled     bool `json:"enabled"`
	StartMinute int  `json:"start_minute"`
	EndMinute   int  `json:"end_minute"`
}

func GetScreenOff() (*ScreenOff, error) {
	data, err := Get("/screenoff/get")
	if err != nil {
		return nil, err
	}

	var screenOff ScreenOff
	if err := decodeScreenOff(data, &screenOff); err != nil {
		return nil, err
	}
	if err := validateScreenOff(screenOff); err != nil {
		return nil, fmt.Errorf("invalid screen-off response: %w", err)
	}

	return &screenOff, nil
}

func SetScreenOff(screenOff ScreenOff) error {
	if err := validateScreenOff(screenOff); err != nil {
		return err
	}

	data, err := json.Marshal(screenOff)
	if err != nil {
		return err
	}

	rspData, err := Post("/screenoff/set", data)
	if err != nil {
		return err
	}
	return decodeSetScreenOffResponse(rspData)
}

func decodeSetScreenOffResponse(data []byte) error {
	var rsp struct {
		Status  string      `json:"status"`
		Message interface{} `json:"message"`
	}
	if err := json.Unmarshal(data, &rsp); err != nil {
		return fmt.Errorf("decode screen-off response: %w", err)
	}
	if rsp.Status != "ok" {
		return fmt.Errorf("set screen-off schedule failed: %v", rsp.Message)
	}
	return nil
}

func decodeScreenOff(data []byte, screenOff *ScreenOff) error {
	var wire struct {
		Enabled     *bool `json:"enabled"`
		StartMinute *int  `json:"start_minute"`
		EndMinute   *int  `json:"end_minute"`
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&wire); err != nil {
		return err
	}
	if wire.Enabled == nil || wire.StartMinute == nil || wire.EndMinute == nil {
		return fmt.Errorf("missing required screen-off fields")
	}
	var trailing interface{}
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON data")
		}
		return err
	}

	screenOff.Enabled = *wire.Enabled
	screenOff.StartMinute = *wire.StartMinute
	screenOff.EndMinute = *wire.EndMinute
	return nil
}

func validateScreenOff(screenOff ScreenOff) error {
	if screenOff.StartMinute < 0 || screenOff.StartMinute >= minutesPerDay ||
		screenOff.EndMinute < 0 || screenOff.EndMinute >= minutesPerDay {
		return fmt.Errorf("screen-off minutes must be between 0 and %d", minutesPerDay-1)
	}
	if screenOff.StartMinute == screenOff.EndMinute {
		return fmt.Errorf("screen-off start and end minutes must differ")
	}
	return nil
}
