package ai

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const (
	snapshotSocketPath      = "/run/kvm/vin_snapshot.sock"
	snapshotProtocolVersion = 1
	snapshotHeaderMaxBytes  = 64 * 1024
	snapshotJPEGMaxBytes    = 32 * 1024 * 1024
	snapshotTimeoutMaxMS    = 30_000
	snapshotMaxConcurrent   = 8
	defaultSnapshotTimeout  = 1000
	snapshotDialTimeout     = 5 * time.Second
	snapshotDeadlineGrace   = time.Second
)

var errSnapshotClientBusy = errors.New("snapshot client concurrency limit reached")

var snapshotCaptureLimiter = newSnapshotClientLimiter(snapshotMaxConcurrent)

type SnapshotRequest struct {
	Quality   int  `json:"quality"`
	TimeoutMS *int `json:"timeoutMs,omitempty"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	W         int  `json:"w"`
	H         int  `json:"h"`
}

type socketSnapshotRequest struct {
	Version   int    `json:"version"`
	Cmd       string `json:"cmd"`
	Quality   int    `json:"quality"`
	TimeoutMS int    `json:"timeout_ms"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	W         int    `json:"w"`
	H         int    `json:"h"`
}

type SnapshotHeader struct {
	Version    int    `json:"version,omitempty"`
	OK         bool   `json:"ok"`
	RetCode    int    `json:"ret_code"`
	Message    string `json:"message"`
	SessionID  string `json:"session_id,omitempty"`
	Sequence   uint64 `json:"sequence,omitempty"`
	CaptureID  string `json:"capture_id,omitempty"`
	CapturedAt string `json:"captured_at,omitempty"`
	Source     string `json:"source,omitempty"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Size       int    `json:"size"`
}

type snapshotDialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)

type snapshotClientLimiter struct {
	slots chan struct{}
}

func newSnapshotClientLimiter(limit int) *snapshotClientLimiter {
	return &snapshotClientLimiter{slots: make(chan struct{}, limit)}
}

func (l *snapshotClientLimiter) acquire(ctx context.Context) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case l.slots <- struct{}{}:
		return func() { <-l.slots }, nil
	default:
		return nil, errSnapshotClientBusy
	}
}

func buildSocketSnapshotRequest(req SnapshotRequest) socketSnapshotRequest {
	quality := req.Quality
	if quality == 0 {
		quality = 85
	}
	quality = clampInt(quality, 1, 100)

	return socketSnapshotRequest{
		Version:   snapshotProtocolVersion,
		Cmd:       "capture_fresh",
		Quality:   quality,
		TimeoutMS: snapshotTimeoutMS(req),
		X:         req.X,
		Y:         req.Y,
		W:         req.W,
		H:         req.H,
	}
}

// CaptureSnapshotContext requests a new JPEG captured after the request reaches
// kvm_vin through the dedicated snapshot Socket.
func CaptureSnapshotContext(ctx context.Context, req SnapshotRequest) (SnapshotHeader, []byte, error) {
	dialer := &net.Dialer{Timeout: snapshotDialTimeout}
	return captureSnapshotContextWithDialer(
		ctx,
		req,
		snapshotCaptureLimiter,
		dialer.DialContext,
	)
}

func captureSnapshotContextWithDialer(
	ctx context.Context,
	req SnapshotRequest,
	limiter *snapshotClientLimiter,
	dialContext snapshotDialContextFunc,
) (SnapshotHeader, []byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	deadline := snapshotClientDeadline(req)
	operationCtx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	effectiveDeadline, _ := operationCtx.Deadline()

	release, err := limiter.acquire(operationCtx)
	if err != nil {
		if errors.Is(err, errSnapshotClientBusy) {
			header := SnapshotHeader{
				Version: snapshotProtocolVersion,
				OK:      false,
				RetCode: -16,
				Message: errSnapshotClientBusy.Error(),
			}
			return header, nil, errSnapshotClientBusy
		}
		return snapshotContextFailure(err)
	}
	defer release()

	header, jpeg, err := captureSnapshotAttempt(
		operationCtx,
		effectiveDeadline,
		snapshotSocketPath,
		req,
		dialContext,
	)
	if err != nil {
		if ctxErr := operationCtx.Err(); ctxErr != nil {
			return snapshotContextFailure(ctxErr)
		}
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return snapshotContextFailure(context.DeadlineExceeded)
		}
	}
	return header, jpeg, err
}

func captureSnapshotFromSocket(
	ctx context.Context,
	socketPath string,
	req SnapshotRequest,
	dialContext snapshotDialContextFunc,
) (SnapshotHeader, []byte, error) {
	deadline := snapshotClientDeadline(req)
	operationCtx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()
	effectiveDeadline, _ := operationCtx.Deadline()

	return captureSnapshotAttempt(
		operationCtx,
		effectiveDeadline,
		socketPath,
		req,
		dialContext,
	)
}

func captureSnapshotAttempt(
	ctx context.Context,
	deadline time.Time,
	socketPath string,
	req SnapshotRequest,
	dialContext snapshotDialContextFunc,
) (SnapshotHeader, []byte, error) {
	conn, err := dialContext(ctx, "unix", socketPath)
	if err != nil {
		return SnapshotHeader{}, nil, snapshotContextError(ctx, "dial snapshot socket", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(deadline); err != nil {
		return SnapshotHeader{}, nil, fmt.Errorf("set snapshot socket deadline: %w", err)
	}

	stopClose := context.AfterFunc(ctx, func() {
		_ = conn.Close()
	})
	defer stopClose()

	if err := writeSnapshotRequest(conn, buildSocketSnapshotRequest(req)); err != nil {
		return SnapshotHeader{}, nil, snapshotContextError(ctx, "write snapshot request", err)
	}

	header, jpeg, err := readSnapshotResponse(conn)
	if err != nil {
		return header, nil, snapshotContextError(ctx, "read snapshot response", err)
	}
	return header, jpeg, nil
}

func snapshotClientDeadline(req SnapshotRequest) time.Time {
	timeout := 5 * time.Second
	if req.TimeoutMS != nil {
		timeoutMS := snapshotTimeoutMS(req)
		if timeoutMS > 0 {
			timeout = time.Duration(timeoutMS)*time.Millisecond + snapshotDeadlineGrace
		}
	}
	return time.Now().Add(timeout)
}

func snapshotTimeoutMS(req SnapshotRequest) int {
	if req.TimeoutMS == nil {
		return defaultSnapshotTimeout
	}
	if *req.TimeoutMS <= 0 {
		return 0
	}
	return clampInt(*req.TimeoutMS, 1, snapshotTimeoutMaxMS)
}

func writeSnapshotRequest(writer io.Writer, request socketSnapshotRequest) error {
	return writeJSONLine(writer, request)
}

func writeJSONLine(writer io.Writer, request any) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal snapshot request: %w", err)
	}
	payload = append(payload, '\n')
	if err := writeAll(writer, payload); err != nil {
		return err
	}
	return nil
}

func writeAll(writer io.Writer, payload []byte) error {
	for len(payload) > 0 {
		n, err := writer.Write(payload)
		if n < 0 || n > len(payload) {
			return io.ErrShortWrite
		}
		payload = payload[n:]
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}

func readSnapshotResponse(reader io.Reader) (SnapshotHeader, []byte, error) {
	buffered := bufio.NewReaderSize(reader, snapshotHeaderMaxBytes+1)
	line, err := buffered.ReadSlice('\n')
	if errors.Is(err, bufio.ErrBufferFull) || len(line) > snapshotHeaderMaxBytes ||
		(errors.Is(err, io.EOF) && len(line) >= snapshotHeaderMaxBytes) {
		return SnapshotHeader{}, nil, fmt.Errorf("snapshot header too large: max %d bytes", snapshotHeaderMaxBytes)
	}
	if err != nil {
		return SnapshotHeader{}, nil, fmt.Errorf("read snapshot header: %w", err)
	}

	var header SnapshotHeader
	if err := json.Unmarshal(line, &header); err != nil {
		return SnapshotHeader{}, nil, fmt.Errorf("parse snapshot header: %w", err)
	}
	if header.Version != snapshotProtocolVersion {
		return header, nil, fmt.Errorf("unsupported snapshot protocol version: %d", header.Version)
	}
	if header.Size < 0 {
		return header, nil, fmt.Errorf("invalid snapshot JPEG size: %d", header.Size)
	}
	if !header.OK {
		return header, nil, nil
	}
	if header.RetCode != 0 {
		return header, nil, fmt.Errorf("successful snapshot response has ret_code %d", header.RetCode)
	}
	if strings.TrimSpace(header.SessionID) == "" {
		return header, nil, fmt.Errorf("snapshot response missing session_id")
	}
	if strings.TrimSpace(header.CaptureID) == "" {
		return header, nil, fmt.Errorf("snapshot response missing capture_id")
	}
	if strings.TrimSpace(header.CapturedAt) == "" {
		return header, nil, fmt.Errorf("snapshot response missing captured_at")
	}
	if strings.TrimSpace(header.Source) == "" {
		return header, nil, fmt.Errorf("snapshot response missing source")
	}
	if header.Width <= 0 || header.Height <= 0 {
		return header, nil, fmt.Errorf("invalid snapshot dimensions: %dx%d", header.Width, header.Height)
	}
	if header.Size == 0 {
		return header, nil, fmt.Errorf("invalid snapshot JPEG size: 0")
	}
	if header.Size > snapshotJPEGMaxBytes {
		return header, nil, fmt.Errorf("snapshot JPEG too large: %d bytes exceeds max %d", header.Size, snapshotJPEGMaxBytes)
	}

	jpeg := make([]byte, header.Size)
	if _, err := io.ReadFull(buffered, jpeg); err != nil {
		return header, nil, fmt.Errorf("read snapshot JPEG: %w", err)
	}

	return header, jpeg, nil
}

func snapshotContextFailure(err error) (SnapshotHeader, []byte, error) {
	retCode := -18
	message := "snapshot request cancelled"
	if errors.Is(err, context.DeadlineExceeded) {
		retCode = -13
		message = "snapshot timeout"
	}
	return SnapshotHeader{
		Version: snapshotProtocolVersion,
		OK:      false,
		RetCode: retCode,
		Message: message,
	}, nil, err
}

func snapshotContextError(ctx context.Context, operation string, err error) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return fmt.Errorf("%s: %w", operation, ctxErr)
	}
	return fmt.Errorf("%s: %w", operation, err)
}
