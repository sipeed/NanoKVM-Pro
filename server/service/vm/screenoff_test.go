package vm

import (
	"NanoKVM-Server/proto"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestScreenTypeFromVersion(t *testing.T) {
	tests := []struct {
		version string
		want    screenType
		wantErr bool
	}{
		{version: "NanoKVM-ATX", want: screenTypeATX},
		{version: "NanoKVM-Desk", want: screenTypeDesk},
		{version: "NanoKVM-Unknown", wantErr: true},
	}

	for _, test := range tests {
		got, err := screenTypeFromVersion([]byte(test.version))
		if (err != nil) != test.wantErr || got != test.want {
			t.Errorf("screenTypeFromVersion(%q) = %q, %v; want %q, error=%t", test.version, got, err, test.want, test.wantErr)
		}
	}
}

func TestGetScreenTypeUsesReader(t *testing.T) {
	original := readScreenVersion
	t.Cleanup(func() { readScreenVersion = original })

	readScreenVersion = func(path string) ([]byte, error) {
		if path != screenVersionPath {
			t.Fatalf("path = %q, want %q", path, screenVersionPath)
		}
		return []byte("NanoKVM-Desk"), nil
	}

	got, err := getScreenType()
	if err != nil || got != screenTypeDesk {
		t.Fatalf("getScreenType() = %q, %v", got, err)
	}

	readScreenVersion = func(string) ([]byte, error) { return nil, errors.New("read failed") }
	if _, err := getScreenType(); err == nil {
		t.Fatal("expected read error")
	}
}

func TestValidateLcdScreenOff(t *testing.T) {
	boolPtr := func(value bool) *bool { return &value }
	intPtr := func(value int) *int { return &value }
	tests := []struct {
		req     proto.SetLcdScreenOffReq
		wantErr bool
	}{
		{req: proto.SetLcdScreenOffReq{Enabled: boolPtr(false), StartMinute: intPtr(0), EndMinute: intPtr(1)}},
		{req: proto.SetLcdScreenOffReq{Enabled: boolPtr(true), StartMinute: intPtr(1438), EndMinute: intPtr(1439)}},
		{req: proto.SetLcdScreenOffReq{Enabled: boolPtr(true), StartMinute: intPtr(1), EndMinute: intPtr(1)}, wantErr: true},
		{req: proto.SetLcdScreenOffReq{Enabled: boolPtr(true), StartMinute: intPtr(-1), EndMinute: intPtr(1)}, wantErr: true},
		{req: proto.SetLcdScreenOffReq{Enabled: boolPtr(true), StartMinute: intPtr(1), EndMinute: intPtr(1440)}, wantErr: true},
		{req: proto.SetLcdScreenOffReq{Enabled: boolPtr(true), EndMinute: intPtr(1)}, wantErr: true},
	}

	for _, test := range tests {
		if err := validateLcdScreenOff(test.req); (err != nil) != test.wantErr {
			t.Errorf("validateLcdScreenOff(%#v) error = %v, wantErr %t", test.req, err, test.wantErr)
		}
	}
}

func TestSetLcdScreenOffRejectsATX(t *testing.T) {
	gin.SetMode(gin.TestMode)
	original := readScreenVersion
	t.Cleanup(func() { readScreenVersion = original })
	readScreenVersion = func(string) ([]byte, error) { return []byte("NanoKVM-ATX"), nil }

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/vm/lcd/screen-off", strings.NewReader(`{"enabled":true,"startMinute":1020,"endMinute":540}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	NewService().SetLcdScreenOff(ctx)

	var response proto.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if recorder.Code != http.StatusOK || response.Code != -3 || response.Msg != "screen-off scheduling is only supported on Desk" {
		t.Fatalf("unexpected response: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestGetLcdScreenOffReportsATXAsUnsupported(t *testing.T) {
	gin.SetMode(gin.TestMode)
	original := readScreenVersion
	t.Cleanup(func() { readScreenVersion = original })
	readScreenVersion = func(string) ([]byte, error) { return []byte("NanoKVM-ATX"), nil }

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/vm/lcd/screen-off", nil)

	NewService().GetLcdScreenOff(ctx)

	var response struct {
		Code int `json:"code"`
		Data struct {
			Supported bool `json:"supported"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if recorder.Code != http.StatusOK || response.Code != 0 || response.Data.Supported {
		t.Fatalf("unexpected response: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}
