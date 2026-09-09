package vm

import (
	"NanoKVM-Server/proto"
	"encoding/json"
	"testing"
)

func TestValidateDisplayPolicyRequest(t *testing.T) {
	mode := modeIdleClock
	if err := validateDisplayPolicyRequest(proto.SetLcdDisplayPolicyReq{Mode: &mode}); err != nil {
		t.Fatalf("valid mode rejected: %v", err)
	}

	badMode := "scheduledOff"
	if err := validateDisplayPolicyRequest(proto.SetLcdDisplayPolicyReq{Mode: &badMode}); err == nil {
		t.Fatal("unsupported mode accepted")
	}

	enabled := true
	start, end := 60, 60
	if err := validateDisplayPolicyRequest(proto.SetLcdDisplayPolicyReq{Schedule: &proto.SetLcdDisplaySchedule{
		Enabled:     &enabled,
		StartMinute: &start,
		EndMinute:   &end,
	}}); err == nil {
		t.Fatal("equal enabled schedule endpoints accepted")
	}
	validEnd := 1
	for name, schedule := range map[string]*proto.SetLcdDisplaySchedule{
		"all fields":  {},
		"enabled":     {StartMinute: &start, EndMinute: &validEnd},
		"startMinute": {Enabled: &enabled, EndMinute: &validEnd},
		"endMinute":   {Enabled: &enabled, StartMinute: &start},
	} {
		t.Run("missing "+name, func(t *testing.T) {
			if err := validateDisplayPolicyRequest(proto.SetLcdDisplayPolicyReq{Schedule: schedule}); err == nil {
				t.Fatal("incomplete schedule accepted")
			}
		})
	}

	disabled, zero := false, 0
	if err := validateDisplayPolicyRequest(proto.SetLcdDisplayPolicyReq{Schedule: &proto.SetLcdDisplaySchedule{
		Enabled: &disabled, StartMinute: &zero, EndMinute: &zero,
	}}); err != nil {
		t.Fatalf("explicit false and zero values rejected: %v", err)
	}

	if err := validateDisplayPolicyRequest(proto.SetLcdDisplayPolicyReq{}); err == nil {
		t.Fatal("empty update accepted")
	}
}

func TestValidateDisplayPolicyRequestRejectsEmptyJSONSchedule(t *testing.T) {
	var req proto.SetLcdDisplayPolicyReq
	if err := json.Unmarshal([]byte(`{"schedule":{}}`), &req); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	if err := validateDisplayPolicyRequest(req); err == nil {
		t.Fatal("empty JSON schedule accepted")
	}
}
