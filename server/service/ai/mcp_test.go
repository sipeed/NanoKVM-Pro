package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func isolatedConfig(t *testing.T) {
	t.Helper()
	old := configFilePath
	configFilePath = filepath.Join(t.TempDir(), "ai.json")
	t.Cleanup(func() { configFilePath = old })
}

func TestConfigurationLifecycleAndAuthentication(t *testing.T) {
	isolatedConfig(t)
	gin.SetMode(gin.TestMode)
	service := NewService()
	engine := gin.New()
	engine.GET("/config", service.GetConfig)
	engine.POST("/config", service.SetConfig)
	engine.POST("/key", service.GetMCPAPIKey)
	engine.POST("/regenerate", service.RegenerateMCPAPIKey)
	engine.Any("/mcp", MCPAPIKeyMiddleware(), func(c *gin.Context) { c.Status(http.StatusNoContent) })
	request := func(method, path, body, key string) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		if key != "" {
			r.Header.Set("Authorization", key)
		}
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, r)
		return w
	}
	if w := request("POST", "/mcp", "", ""); w.Code != 403 {
		t.Fatalf("default state: %d", w.Code)
	}
	if w := request("POST", "/config", `{}`, ""); !strings.Contains(w.Body.String(), "invalid arguments") {
		t.Fatal(w.Body.String())
	}
	if w := request("POST", "/config", `{"mcpEnabled":true}`, ""); !strings.Contains(w.Body.String(), `"code":0`) {
		t.Fatal(w.Body.String())
	}
	cfg, err := loadConfig()
	if err != nil || !cfg.MCPEnabled || !strings.HasPrefix(cfg.MCPAPIKey, "nap_mcp_") {
		t.Fatalf("config = %+v, %v", cfg, err)
	}
	oldKey := cfg.MCPAPIKey
	info, err := os.Stat(configFilePath)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("credential permissions: %v %v", info, err)
	}
	if w := request("GET", "/config", "", ""); strings.Contains(w.Body.String(), oldKey) {
		t.Fatal("ordinary config leaks full API key")
	}
	if w := request("POST", "/key", "", ""); !strings.Contains(w.Body.String(), oldKey) || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("protected key endpoint failed")
	}
	for _, header := range []string{"", oldKey, "Bearer wrong", "Basic " + oldKey, "Bearer "} {
		if w := request("POST", "/mcp", "", header); w.Code != 401 {
			t.Fatalf("invalid auth accepted: %q: %d", header, w.Code)
		}
	}
	if w := request("POST", "/mcp", "", "Bearer "+oldKey); w.Code != 204 {
		t.Fatalf("valid auth: %d", w.Code)
	}
	if w := request("POST", "/regenerate", "", ""); !strings.Contains(w.Body.String(), `"code":0`) {
		t.Fatal(w.Body.String())
	}
	cfg, err = loadConfig()
	if err != nil || cfg.MCPAPIKey == oldKey || !cfg.MCPEnabled {
		t.Fatal("key rotation failed")
	}
	if w := request("POST", "/mcp", "", "Bearer "+oldKey); w.Code != 401 {
		t.Fatal("old key remains valid")
	}
	if w := request("POST", "/mcp", "", "Bearer "+cfg.MCPAPIKey); w.Code != 204 {
		t.Fatal("new key rejected")
	}
	request("POST", "/config", `{"mcpEnabled":false}`, "")
	if w := request("POST", "/mcp", "", "Bearer "+cfg.MCPAPIKey); w.Code != 403 {
		t.Fatal("disable failed")
	}
	if err := os.WriteFile(configFilePath, []byte("{invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	if w := request("POST", "/mcp", "", "Bearer "+cfg.MCPAPIKey); w.Code != 500 {
		t.Fatal("corrupt config did not fail closed")
	}
}

type recordingHID struct {
	keyboard, relative, absolute [][]byte
	onKeyboard                   func()
}

func (w *recordingHID) WriteKeyboardReport(p []byte) error {
	w.keyboard = append(w.keyboard, bytes.Clone(p))
	if w.onKeyboard != nil {
		w.onKeyboard()
	}
	return nil
}
func (w *recordingHID) WriteRelativeMouseReport(p []byte) error {
	w.relative = append(w.relative, bytes.Clone(p))
	return nil
}
func (w *recordingHID) WriteAbsoluteMouseReport(p []byte) error {
	w.absolute = append(w.absolute, bytes.Clone(p))
	return nil
}
func testRemote(w HIDWriter) *Remote {
	return &Remote{gate: make(chan struct{}, 1), action: func(action func(HIDWriter) error) error { return action(w) }}
}

func TestProReportsAndCancellationRelease(t *testing.T) {
	writer := &recordingHID{}
	remote := testRemote(writer)
	if err := remote.MoveAbsolute(context.Background(), 0.5, 1); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(writer.absolute[0], []byte{0, 0, 0x40, 0xff, 0x7f, 0}) {
		t.Fatalf("Pro absolute report: %v", writer.absolute)
	}
	if err := remote.Click(context.Background(), "right", 1, true, 0.5, 1); err != nil {
		t.Fatal(err)
	}
	if writer.absolute[1][0] != 2 || writer.absolute[len(writer.absolute)-1][0] != 0 {
		t.Fatal("absolute click did not release")
	}
	if err := remote.MoveRelative(context.Background(), 300, -300); err != nil {
		t.Fatal(err)
	}
	var x, y int
	for _, report := range writer.relative {
		x += int(int8(report[1]))
		y += int(int8(report[2]))
	}
	if x != 300 || y != -300 {
		t.Fatalf("relative movement wrapped: %d,%d", x, y)
	}
	ctx, cancel := context.WithCancel(context.Background())
	writer.onKeyboard = cancel
	err := remote.PressKeys(ctx, []string{"ControlLeft", "KeyA"}, time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled action: %v", err)
	}
	if !bytes.Equal(writer.keyboard[len(writer.keyboard)-1], make([]byte, 8)) {
		t.Fatal("canceled shortcut left keys held")
	}
}

func testJPEG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 64, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 3), G: uint8(y * 5), B: 42, A: 255})
		}
	}
	var out bytes.Buffer
	if err := jpeg.Encode(&out, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

type bearerTransport struct{ key string }

func (transport bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	cloned.Header.Set("Authorization", "Bearer "+transport.key)
	return http.DefaultTransport.RoundTrip(cloned)
}

func TestMCPStreamableHTTPToolsAndJPEG(t *testing.T) {
	isolatedConfig(t)
	if err := saveConfig(Config{MCPEnabled: true, MCPAPIKey: "test-credential"}); err != nil {
		t.Fatal(err)
	}
	writer := &recordingHID{}
	previous := captureSnapshot
	t.Cleanup(func() { captureSnapshot = previous })
	data := testJPEG(t)
	captureSnapshot = func(ctx context.Context, req SnapshotRequest) (SnapshotHeader, []byte, error) {
		if req.X != 0 || req.Y != 0 || req.W != 0 || req.H != 0 {
			t.Error("software crop changed capture request")
		}
		return SnapshotHeader{OK: true, Version: 1, SessionID: "producer-1", CaptureID: "producer-1-3", CapturedAt: "2026-11-13T00:00:00Z", Source: "pro-venc-mjpeg", Width: 64, Height: 32, Size: len(data)}, data, nil
	}
	engine := gin.New()
	engine.Any("/api/mcp", MCPAPIKeyMiddleware(), gin.WrapH(newMCPHandler(testRemote(writer), nil)))
	server := httptest.NewServer(engine)
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := mcp.NewClient(&mcp.Implementation{Name: "pro-port-test", Version: "1"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: server.URL + "/api/mcp", HTTPClient: &http.Client{Transport: bearerTransport{key: "test-credential"}}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	list, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, tool := range list.Tools {
		names = append(names, tool.Name)
	}
	sort.Strings(names)
	want := []string{"click_mouse", "move_mouse", "press_keys", "screenshot", "scroll_mouse", "type_text"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("tools = %v", names)
	}
	call := func(name string, args map[string]any) *mcp.CallToolResult {
		t.Helper()
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
		if err != nil || result.IsError {
			t.Fatalf("%s: %+v %v", name, result, err)
		}
		return result
	}
	call("type_text", map[string]any{"text": "aA!中", "delayMs": 0})
	call("press_keys", map[string]any{"keys": []string{"ControlLeft", "AltLeft", "Delete"}, "holdMs": 0})
	call("move_mouse", map[string]any{"mode": "absolute", "x": 0.5, "y": 1})
	call("click_mouse", map[string]any{"button": "left", "mode": "absolute", "x": 0.5, "y": 1})
	call("scroll_mouse", map[string]any{"deltaY": 300})
	result := call("screenshot", map[string]any{"quality": 75, "x": 8, "y": 4, "w": 16, "h": 8})
	var output ScreenshotResult
	encoded, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(encoded, &output); err != nil {
		t.Fatal(err)
	}
	if !output.OK || output.SourceWidth != 64 || output.SourceHeight != 32 || output.Width != 16 || output.Height != 8 || output.Capture.Left != 8 || output.Capture.Top != 4 {
		t.Fatalf("crop metadata: %+v", output)
	}
	var found bool
	for _, content := range result.Content {
		if img, ok := content.(*mcp.ImageContent); ok {
			decoded, err := jpeg.DecodeConfig(bytes.NewReader(img.Data))
			if err != nil || decoded.Width != 16 || decoded.Height != 8 || img.MIMEType != "image/jpeg" {
				t.Fatalf("MCP image: %+v %v", decoded, err)
			}
			found = true
		}
	}
	if !found {
		t.Fatal("screenshot did not return an MCP image")
	}
	for _, args := range []map[string]any{{"mode": "absolute", "x": 1920, "y": 1080}, {"mode": "absolute", "x": 0.5}, {"mode": "relative", "deltaX": 1000000000}} {
		before := len(writer.absolute) + len(writer.relative)
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "move_mouse", Arguments: args})
		if err != nil {
			t.Fatal(err)
		}
		if !result.IsError {
			t.Fatalf("invalid mouse input accepted: %v", args)
		}
		if before != len(writer.absolute)+len(writer.relative) {
			t.Fatal("invalid input wrote to USB")
		}
	}
	result, err = session.CallTool(ctx, &mcp.CallToolParams{Name: "screenshot", Arguments: map[string]any{"x": 63, "y": 0, "w": 2, "h": 1}})
	if err != nil || !result.IsError {
		t.Fatal("out-of-bounds crop accepted")
	}
}

func TestCanceledInputDoesNotWaitForBusyMCPAction(t *testing.T) {
	writer := &recordingHID{}
	remote := testRemote(writer)
	remote.gate <- struct{}{}
	defer func() { <-remote.gate }()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan error, 1)
	go func() { done <- remote.PressKeys(ctx, []string{"KeyA"}, 0) }()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("queued cancellation: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("canceled input remained queued behind another action")
	}
	if len(writer.keyboard) != 0 {
		t.Fatal("canceled queued input reached USB")
	}
}
