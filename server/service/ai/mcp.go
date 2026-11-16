package ai

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TypeTextParams struct {
	Text    string `json:"text" jsonschema:"short ASCII text to type on the remote host"`
	DelayMS *int   `json:"delayMs,omitempty" jsonschema:"delay between HID reports in milliseconds, default 30, range 0..1000"`
}

type TypeTextResult struct {
	Typed        int    `json:"typed"`
	Skipped      string `json:"skipped,omitempty"`
	SkippedCount int    `json:"skippedCount"`
}

type PressKeysParams struct {
	Keys   []string `json:"keys" jsonschema:"KeyboardEvent.code key names, at most six non-modifier keys"`
	HoldMS *int     `json:"holdMs,omitempty" jsonschema:"hold time in milliseconds, default 50, range 0..1000"`
}

type MoveMouseParams struct {
	Mode            string   `json:"mode" jsonschema:"absolute or relative"`
	X               *float64 `json:"x,omitempty" jsonschema:"absolute normalized x, 0.0 to 1.0 inside the full source frame; required for absolute movement"`
	Y               *float64 `json:"y,omitempty" jsonschema:"absolute normalized y, 0.0 to 1.0 inside the full source frame; required for absolute movement"`
	DeltaX          int      `json:"deltaX,omitempty"`
	DeltaY          int      `json:"deltaY,omitempty"`
	CoordinateSpace string   `json:"coordinateSpace,omitempty" jsonschema:"content (default) or source; both refer to the full source frame on Pro"`
}

type ClickMouseParams struct {
	Button          string   `json:"button" jsonschema:"left, right, middle, back or forward"`
	Clicks          int      `json:"clicks,omitempty" jsonschema:"click count, default 1, maximum 5"`
	Mode            string   `json:"mode,omitempty" jsonschema:"absolute or relative, default relative"`
	X               *float64 `json:"x,omitempty" jsonschema:"absolute normalized x, 0.0 to 1.0 inside the full source frame"`
	Y               *float64 `json:"y,omitempty" jsonschema:"absolute normalized y, 0.0 to 1.0 inside the full source frame"`
	CoordinateSpace string   `json:"coordinateSpace,omitempty" jsonschema:"content (default) or source; a screenshot crop never changes this coordinate space"`
}

type ScrollMouseParams struct {
	DeltaY int `json:"deltaY"`
}
type ToolOKResult struct {
	OK bool `json:"ok"`
}

type ScreenshotParams struct {
	Quality   int  `json:"quality,omitempty" jsonschema:"output JPEG quality, default 85, range 1..100"`
	TimeoutMS *int `json:"timeoutMs,omitempty" jsonschema:"fresh-frame wait in milliseconds, default 1000, maximum 30000"`
	X         int  `json:"x,omitempty" jsonschema:"crop origin in full-frame pixels; requires w and h"`
	Y         int  `json:"y,omitempty" jsonschema:"crop origin in full-frame pixels; requires w and h"`
	W         int  `json:"w,omitempty" jsonschema:"crop width in pixels; zero with h zero returns the full frame"`
	H         int  `json:"h,omitempty" jsonschema:"crop height in pixels; cropping does not change absolute mouse coordinates"`
}

type ScreenshotContentRect struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ScreenshotResult struct {
	OK              bool                   `json:"ok"`
	RetCode         int                    `json:"retCode"`
	Message         string                 `json:"message"`
	SessionID       string                 `json:"sessionId,omitempty"`
	CaptureID       string                 `json:"captureId,omitempty"`
	CapturedAt      string                 `json:"capturedAt,omitempty"`
	Source          string                 `json:"source,omitempty"`
	Width           int                    `json:"width"`
	Height          int                    `json:"height"`
	Size            int                    `json:"size"`
	SourceWidth     int                    `json:"sourceWidth,omitempty"`
	SourceHeight    int                    `json:"sourceHeight,omitempty"`
	CoordinateSpace string                 `json:"coordinateSpace"`
	AreaSource      string                 `json:"areaSource"`
	Content         *ScreenshotContentRect `json:"content,omitempty"`
	Capture         *ScreenshotContentRect `json:"capture,omitempty"`
}

var captureSnapshot = CaptureSnapshotContext

func NewMCPHandler() http.Handler {
	return newMCPHandler(NewRemote(), newToolRateLimiter(20*time.Millisecond))
}

func newMCPHandler(remote *Remote, limiter *toolRateLimiter) http.Handler {
	server := mcp.NewServer(&mcp.Implementation{Name: "nanokvm-remote-control", Version: "v1.0.0"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "type_text", Description: "Type short ASCII text into the remote host. Unsupported characters are reported as skipped."}, typeTextHandler(remote, limiter))
	mcp.AddTool(server, &mcp.Tool{Name: "press_keys", Description: "Press and release one keyboard shortcut on the remote host."}, pressKeysHandler(remote, limiter))
	mcp.AddTool(server, &mcp.Tool{Name: "move_mouse", Description: "Move the remote pointer. Absolute x/y are normalized 0..1 coordinates inside the full source frame, never pixels."}, moveMouseHandler(remote, limiter))
	mcp.AddTool(server, &mcp.Tool{Name: "click_mouse", Description: "Click a remote mouse button, optionally moving to an absolute position first. A cropped screenshot does not change the full-frame coordinate space."}, clickMouseHandler(remote, limiter))
	mcp.AddTool(server, &mcp.Tool{Name: "scroll_mouse", Description: "Scroll the remote mouse wheel vertically."}, scrollMouseHandler(remote, limiter))
	mcp.AddTool(server, &mcp.Tool{Name: "screenshot", Description: "Capture a fresh JPEG from the remote display without switching the video mode. Optional crop uses full-frame pixels; absolute mouse coordinates always refer to the full frame."}, screenshotHandler(limiter))
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, &mcp.StreamableHTTPOptions{SessionTimeout: 5 * time.Minute})
}

func inputDelay(value *int, fallback int) time.Duration {
	if value == nil {
		return time.Duration(fallback) * time.Millisecond
	}
	return time.Duration(clampInt(*value, 0, 1000)) * time.Millisecond
}

func typeTextHandler(remote *Remote, limiter *toolRateLimiter) mcp.ToolHandlerFor[TypeTextParams, TypeTextResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input TypeTextParams) (*mcp.CallToolResult, TypeTextResult, error) {
		if err := checkRate(limiter, "type_text"); err != nil {
			return nil, TypeTextResult{}, err
		}
		runes := []rune(input.Text)
		if len(runes) > 1024 {
			return nil, TypeTextResult{}, fmt.Errorf("text exceeds maximum length of 1024 runes")
		}
		skipped, err := remote.TypeText(ctx, input.Text, inputDelay(input.DelayMS, 30))
		if err != nil {
			return nil, TypeTextResult{}, err
		}
		return nil, TypeTextResult{Typed: len(runes) - len(skipped), Skipped: string(skipped), SkippedCount: len(skipped)}, nil
	}
}

func pressKeysHandler(remote *Remote, limiter *toolRateLimiter) mcp.ToolHandlerFor[PressKeysParams, ToolOKResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input PressKeysParams) (*mcp.CallToolResult, ToolOKResult, error) {
		if err := checkRate(limiter, "press_keys"); err != nil {
			return nil, ToolOKResult{}, err
		}
		if len(input.Keys) == 0 {
			return nil, ToolOKResult{}, fmt.Errorf("keys must not be empty")
		}
		if err := remote.PressKeys(ctx, input.Keys, inputDelay(input.HoldMS, 50)); err != nil {
			return nil, ToolOKResult{}, err
		}
		return nil, ToolOKResult{OK: true}, nil
	}
}

func absolutePoint(x, y *float64, space string) (float64, float64, error) {
	space = normalizedMode(space, "content")
	if space != "content" && space != "source" {
		return 0, 0, fmt.Errorf("unknown coordinate space: %s", space)
	}
	if x == nil || y == nil {
		return 0, 0, fmt.Errorf("absolute movement requires x and y")
	}
	if math.IsNaN(*x) || math.IsNaN(*y) || math.IsInf(*x, 0) || math.IsInf(*y, 0) || *x < 0 || *x > 1 || *y < 0 || *y > 1 {
		return 0, 0, fmt.Errorf("absolute x/y must be normalized coordinates in 0..1, never pixels")
	}
	return *x, *y, nil
}

func moveMouseHandler(remote *Remote, limiter *toolRateLimiter) mcp.ToolHandlerFor[MoveMouseParams, ToolOKResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input MoveMouseParams) (*mcp.CallToolResult, ToolOKResult, error) {
		if err := checkRate(limiter, "move_mouse"); err != nil {
			return nil, ToolOKResult{}, err
		}
		var err error
		switch normalizedMode(input.Mode, "absolute") {
		case "absolute":
			x, y, e := absolutePoint(input.X, input.Y, input.CoordinateSpace)
			if e != nil {
				return nil, ToolOKResult{}, e
			}
			err = remote.MoveAbsolute(ctx, x, y)
		case "relative":
			err = remote.MoveRelative(ctx, input.DeltaX, input.DeltaY)
		default:
			err = fmt.Errorf("unknown mouse mode: %s", input.Mode)
		}
		if err != nil {
			return nil, ToolOKResult{}, err
		}
		return nil, ToolOKResult{OK: true}, nil
	}
}

func clickMouseHandler(remote *Remote, limiter *toolRateLimiter) mcp.ToolHandlerFor[ClickMouseParams, ToolOKResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input ClickMouseParams) (*mcp.CallToolResult, ToolOKResult, error) {
		if err := checkRate(limiter, "click_mouse"); err != nil {
			return nil, ToolOKResult{}, err
		}
		mode := normalizedMode(input.Mode, "relative")
		if mode != "absolute" && mode != "relative" {
			return nil, ToolOKResult{}, fmt.Errorf("unknown mouse mode: %s", input.Mode)
		}
		var x, y float64
		if mode == "absolute" {
			var err error
			x, y, err = absolutePoint(input.X, input.Y, input.CoordinateSpace)
			if err != nil {
				return nil, ToolOKResult{}, err
			}
		}
		if err := remote.Click(ctx, normalizedMode(input.Button, "left"), input.Clicks, mode == "absolute", x, y); err != nil {
			return nil, ToolOKResult{}, err
		}
		return nil, ToolOKResult{OK: true}, nil
	}
}

func scrollMouseHandler(remote *Remote, limiter *toolRateLimiter) mcp.ToolHandlerFor[ScrollMouseParams, ToolOKResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input ScrollMouseParams) (*mcp.CallToolResult, ToolOKResult, error) {
		if err := checkRate(limiter, "scroll_mouse"); err != nil {
			return nil, ToolOKResult{}, err
		}
		if err := remote.Scroll(ctx, input.DeltaY); err != nil {
			return nil, ToolOKResult{}, err
		}
		return nil, ToolOKResult{OK: true}, nil
	}
}

func screenshotHandler(limiter *toolRateLimiter) mcp.ToolHandlerFor[ScreenshotParams, ScreenshotResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input ScreenshotParams) (*mcp.CallToolResult, ScreenshotResult, error) {
		if err := checkRate(limiter, "screenshot"); err != nil {
			return nil, ScreenshotResult{}, err
		}
		if input.X < 0 || input.Y < 0 || input.W < 0 || input.H < 0 || (input.W == 0) != (input.H == 0) || (input.W == 0 && (input.X != 0 || input.Y != 0)) {
			return nil, ScreenshotResult{}, fmt.Errorf("invalid screenshot crop: provide non-negative x/y and positive w/h, or omit all four")
		}
		header, data, err := captureSnapshot(ctx, SnapshotRequest{TimeoutMS: input.TimeoutMS})
		output := ScreenshotResult{OK: header.OK, RetCode: header.RetCode, Message: header.Message, SessionID: header.SessionID, CaptureID: header.CaptureID, CapturedAt: header.CapturedAt, Source: header.Source, CoordinateSpace: "content", AreaSource: "none"}
		if err != nil {
			output.OK = false
			output.Message = err.Error()
			return toolError(err.Error()), output, nil
		}
		if !header.OK {
			return toolError(header.Message), output, nil
		}
		encoded, rect, err := prepareScreenshot(header, data, input)
		if err != nil {
			output.OK = false
			output.RetCode = -14
			output.Message = err.Error()
			return toolError(err.Error()), output, nil
		}
		output.Width, output.Height, output.Size = rect.Dx(), rect.Dy(), len(encoded)
		output.SourceWidth, output.SourceHeight = header.Width, header.Height
		output.Content = &ScreenshotContentRect{Width: header.Width, Height: header.Height}
		output.Capture = &ScreenshotContentRect{Left: rect.Min.X, Top: rect.Min.Y, Width: rect.Dx(), Height: rect.Dy()}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.ImageContent{Data: encoded, MIMEType: "image/jpeg"}}}, output, nil
	}
}

func prepareScreenshot(header SnapshotHeader, data []byte, input ScreenshotParams) ([]byte, image.Rectangle, error) {
	config, err := jpeg.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, image.Rectangle{}, fmt.Errorf("invalid snapshot JPEG: %w", err)
	}
	if config.Width != header.Width || config.Height != header.Height || config.Width > 3840 || config.Height > 2400 {
		return nil, image.Rectangle{}, fmt.Errorf("snapshot dimensions disagree with JPEG or exceed Pro limits")
	}
	rect := image.Rect(0, 0, config.Width, config.Height)
	if input.W > 0 {
		if input.X >= config.Width || input.Y >= config.Height || input.W > config.Width-input.X || input.H > config.Height-input.Y {
			return nil, image.Rectangle{}, fmt.Errorf("screenshot crop exceeds source dimensions %dx%d", config.Width, config.Height)
		}
		rect = image.Rect(input.X, input.Y, input.X+input.W, input.Y+input.H)
	}
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, rect, err
	}
	img = img.(interface {
		SubImage(image.Rectangle) image.Image
	}).SubImage(rect)
	quality := input.Quality
	if quality == 0 {
		quality = 85
	}
	var out bytes.Buffer
	if err := jpeg.Encode(&out, img, &jpeg.Options{Quality: clampInt(quality, 1, 100)}); err != nil {
		return nil, rect, err
	}
	return out.Bytes(), rect, nil
}

func checkRate(limiter *toolRateLimiter, tool string) error {
	if limiter != nil && !limiter.Allow(tool) {
		return fmt.Errorf("rate limit exceeded for %s", tool)
	}
	return nil
}

func toolError(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: message}}}
}

func normalizedMode(mode, fallback string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		return fallback
	}
	return mode
}
