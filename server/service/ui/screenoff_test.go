package ui

import (
	"testing"
)

func TestDecodeScreenOff(t *testing.T) {
	var screenOff ScreenOff
	if err := decodeScreenOff([]byte(`{"enabled":true,"start_minute":1020,"end_minute":540}`), &screenOff); err != nil {
		t.Fatalf("decode screen-off: %v", err)
	}
	if !screenOff.Enabled || screenOff.StartMinute != 1020 || screenOff.EndMinute != 540 {
		t.Fatalf("unexpected screen-off value: %#v", screenOff)
	}
}

func TestDecodeScreenOffRejectsInvalidPayloads(t *testing.T) {
	tests := []string{
		`{"enabled":true,"start_minute":0}`,
		`{"enabled":true,"start_minute":0,"end_minute":1,"extra":true}`,
		`{"enabled":true,"start_minute":0,"end_minute":1}{}`,
	}

	for _, input := range tests {
		var screenOff ScreenOff
		if err := decodeScreenOff([]byte(input), &screenOff); err == nil {
			t.Errorf("decodeScreenOff(%s) unexpectedly succeeded", input)
		}
	}
}

func TestValidateScreenOff(t *testing.T) {
	tests := []struct {
		name    string
		screen  ScreenOff
		wantErr bool
	}{
		{name: "valid crossing midnight", screen: ScreenOff{StartMinute: 1020, EndMinute: 540}},
		{name: "valid day", screen: ScreenOff{StartMinute: 540, EndMinute: 1020}},
		{name: "same minute", screen: ScreenOff{StartMinute: 540, EndMinute: 540}, wantErr: true},
		{name: "negative start", screen: ScreenOff{StartMinute: -1, EndMinute: 1}, wantErr: true},
		{name: "end outside day", screen: ScreenOff{StartMinute: 1, EndMinute: 1440}, wantErr: true},
	}

	for _, test := range tests {
		if err := validateScreenOff(test.screen); (err != nil) != test.wantErr {
			t.Errorf("%s: validateScreenOff(%#v) error = %v, wantErr %t", test.name, test.screen, err, test.wantErr)
		}
	}
}

func TestDecodeSetScreenOffResponse(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{name: "success", data: `{"status":"ok"}`},
		{name: "internal error", data: `{"status":"err","message":"write failed"}`, wantErr: true},
		{name: "missing status", data: `{}`, wantErr: true},
		{name: "invalid JSON", data: `{`, wantErr: true},
	}

	for _, test := range tests {
		if err := decodeSetScreenOffResponse([]byte(test.data)); (err != nil) != test.wantErr {
			t.Errorf("%s: decodeSetScreenOffResponse() error = %v, wantErr %t", test.name, err, test.wantErr)
		}
	}
}
