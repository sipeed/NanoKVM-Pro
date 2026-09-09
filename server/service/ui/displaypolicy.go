package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// DisplayPolicy is the combined idle-display and scheduled screen-off state
// exposed by newer kvm_ui versions.
type DisplayPolicy struct {
	SupportedModes     []string        `json:"supported_modes"`
	Mode               string          `json:"mode"`
	ModeTimeoutSeconds map[string]int  `json:"mode_timeout_seconds,omitempty"`
	Schedule           DisplaySchedule `json:"schedule"`
}

type DisplaySchedule struct {
	Supported          bool `json:"supported"`
	Enabled            bool `json:"enabled"`
	StartMinute        int  `json:"start_minute"`
	EndMinute          int  `json:"end_minute"`
	WakeTimeoutSeconds int  `json:"wake_timeout_seconds,omitempty"`
}

type DisplayPolicyUpdate struct {
	Mode     *string          `json:"mode,omitempty"`
	Schedule *DisplaySchedule `json:"schedule,omitempty"`
}

func GetDisplayPolicy() (*DisplayPolicy, error) {
	data, err := Get("/display-policy/get")
	if err != nil {
		return nil, err
	}

	var policy DisplayPolicy
	if err := decodeDisplayPolicy(data, &policy); err != nil {
		return nil, fmt.Errorf("decode display policy: %w", err)
	}
	if err := validateDisplayPolicy(policy); err != nil {
		return nil, fmt.Errorf("invalid display policy response: %w", err)
	}
	return &policy, nil
}

func decodeDisplayPolicy(data []byte, policy *DisplayPolicy) error {
	var wire struct {
		Mode               *string        `json:"mode"`
		SupportedModes     *[]string      `json:"supported_modes"`
		ModeTimeoutSeconds map[string]int `json:"mode_timeout_seconds"`
		Schedule           *struct {
			Supported          *bool `json:"supported"`
			Enabled            *bool `json:"enabled"`
			StartMinute        *int  `json:"start_minute"`
			EndMinute          *int  `json:"end_minute"`
			WakeTimeoutSeconds *int  `json:"wake_timeout_seconds"`
		} `json:"schedule"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&wire); err != nil {
		return err
	}
	var trailing interface{}
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON data")
		}
		return err
	}
	if wire.Mode == nil || *wire.Mode == "" {
		return fmt.Errorf("missing required mode")
	}
	if wire.SupportedModes == nil || len(*wire.SupportedModes) == 0 {
		return fmt.Errorf("missing required supported_modes")
	}
	if wire.Schedule == nil || wire.Schedule.Supported == nil || wire.Schedule.Enabled == nil ||
		wire.Schedule.StartMinute == nil || wire.Schedule.EndMinute == nil ||
		wire.Schedule.WakeTimeoutSeconds == nil {
		return fmt.Errorf("missing required schedule fields")
	}
	seen := make(map[string]struct{}, len(*wire.SupportedModes))
	for _, mode := range *wire.SupportedModes {
		if mode == "" {
			return fmt.Errorf("supported_modes contains an empty mode")
		}
		if _, ok := seen[mode]; ok {
			return fmt.Errorf("supported_modes contains duplicate mode %q", mode)
		}
		seen[mode] = struct{}{}
	}
	if _, ok := seen[*wire.Mode]; !ok {
		return fmt.Errorf("mode %q is not supported", *wire.Mode)
	}
	policy.Mode = *wire.Mode
	policy.SupportedModes = append([]string(nil), (*wire.SupportedModes)...)
	policy.ModeTimeoutSeconds = wire.ModeTimeoutSeconds
	policy.Schedule = DisplaySchedule{
		Supported:          *wire.Schedule.Supported,
		Enabled:            *wire.Schedule.Enabled,
		StartMinute:        *wire.Schedule.StartMinute,
		EndMinute:          *wire.Schedule.EndMinute,
		WakeTimeoutSeconds: *wire.Schedule.WakeTimeoutSeconds,
	}
	return nil
}

func SetDisplayPolicy(update DisplayPolicyUpdate) error {
	if err := validateDisplayPolicyUpdate(update); err != nil {
		return err
	}
	data, err := json.Marshal(update)
	if err != nil {
		return err
	}
	rspData, err := Post("/display-policy/set", data)
	if err != nil {
		return err
	}
	var rsp struct {
		Status  string      `json:"status"`
		Message interface{} `json:"message"`
	}
	if err := json.Unmarshal(rspData, &rsp); err != nil {
		return fmt.Errorf("decode display policy response: %w", err)
	}
	if rsp.Status != "ok" {
		return fmt.Errorf("set display policy failed: %v", rsp.Message)
	}
	return nil
}

func validateDisplayPolicy(policy DisplayPolicy) error {
	if policy.Mode == "" || len(policy.SupportedModes) == 0 {
		return fmt.Errorf("mode and supported_modes are required")
	}
	found := false
	for _, mode := range policy.SupportedModes {
		if mode == policy.Mode {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("mode %q is not supported", policy.Mode)
	}
	if err := validateSchedule(policy.Schedule.Enabled, policy.Schedule.StartMinute, policy.Schedule.EndMinute); err != nil {
		return err
	}
	if policy.Schedule.WakeTimeoutSeconds < 0 ||
		(policy.Schedule.Supported && policy.Schedule.WakeTimeoutSeconds == 0) {
		return fmt.Errorf("invalid schedule wake timeout")
	}
	return nil
}

func validateDisplayPolicyUpdate(update DisplayPolicyUpdate) error {
	if update.Mode == nil && update.Schedule == nil {
		return fmt.Errorf("display policy update is empty")
	}
	if update.Mode != nil && *update.Mode == "" {
		return fmt.Errorf("mode is required")
	}
	if update.Schedule != nil {
		if err := validateSchedule(update.Schedule.Enabled, update.Schedule.StartMinute, update.Schedule.EndMinute); err != nil {
			return err
		}
	}
	return nil
}

func validateSchedule(enabled bool, startMinute int, endMinute int) error {
	if startMinute < 0 || startMinute >= minutesPerDay || endMinute < 0 || endMinute >= minutesPerDay {
		return fmt.Errorf("screen-off minutes must be between 0 and %d", minutesPerDay-1)
	}
	if enabled && startMinute == endMinute {
		return fmt.Errorf("screen-off start and end minutes must differ")
	}
	return nil
}
