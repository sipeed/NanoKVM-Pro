package ui

import "testing"

func TestDecodeDisplayPolicyRejectsIncompleteResponses(t *testing.T) {
	tests := []string{
		`{}`,
		`{"mode":"idleClock","supported_modes":[],"schedule":{"supported":true,"enabled":false,"start_minute":0,"end_minute":0,"wake_timeout_seconds":60}}`,
		`{"mode":"idleClock","supported_modes":["alwaysOn"],"schedule":{"supported":true,"enabled":false,"start_minute":0,"end_minute":0,"wake_timeout_seconds":60}}`,
		`{"mode":"idleClock","supported_modes":["idleClock"],"schedule":{"supported":true,"enabled":false,"start_minute":0,"end_minute":0}}`,
		`{"mode":"idleClock","supported_modes":["idleClock"],"schedule":{"supported":true,"enabled":true,"start_minute":10,"end_minute":10,"wake_timeout_seconds":60}}`,
		`{"mode":"idleClock","supported_modes":["idleClock"],"schedule":{"supported":true,"enabled":false,"start_minute":0,"end_minute":0,"wake_timeout_seconds":60}}{}`,
	}
	for _, input := range tests {
		var policy DisplayPolicy
		err := decodeDisplayPolicy([]byte(input), &policy)
		if err == nil {
			err = validateDisplayPolicy(policy)
		}
		if err == nil {
			t.Errorf("display policy response %s unexpectedly passed validation", input)
		}
	}
}

func TestDecodeDisplayPolicyAcceptsCompleteResponse(t *testing.T) {
	var policy DisplayPolicy
	input := `{"mode":"idleClock","supported_modes":["alwaysOn","idleClock","idleOff"],"mode_timeout_seconds":{"idleClock":20,"idleOff":60},"schedule":{"supported":true,"enabled":true,"start_minute":1020,"end_minute":540,"wake_timeout_seconds":60}}`
	if err := decodeDisplayPolicy([]byte(input), &policy); err != nil {
		t.Fatalf("decodeDisplayPolicy() error = %v", err)
	}
	if err := validateDisplayPolicy(policy); err != nil {
		t.Fatalf("validateDisplayPolicy() error = %v", err)
	}
}
