package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type chunkWriter struct {
	buffer bytes.Buffer
	limit  int
}

func (w *chunkWriter) Write(payload []byte) (int, error) {
	if len(payload) > w.limit {
		payload = payload[:w.limit]
	}
	return w.buffer.Write(payload)
}

func successfulSnapshotHeader(size int) SnapshotHeader {
	return SnapshotHeader{
		Version:    snapshotProtocolVersion,
		OK:         true,
		RetCode:    0,
		Message:    "ok",
		SessionID:  "session",
		CaptureID:  "session:7",
		CapturedAt: "2026-07-14T11:00:00.123+08:00",
		Source:     "mcp",
		Width:      1920,
		Height:     1080,
		Size:       size,
	}
}

func encodeSnapshotResponse(t *testing.T, header SnapshotHeader, jpeg []byte) []byte {
	t.Helper()
	payload, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("marshal snapshot header: %v", err)
	}
	payload = append(payload, '\n')
	return append(payload, jpeg...)
}

func snapshotResponseDial(
	t *testing.T,
	header SnapshotHeader,
	jpeg []byte,
	requests chan<- []byte,
) snapshotDialContextFunc {
	t.Helper()
	response := encodeSnapshotResponse(t, header, jpeg)
	return func(context.Context, string, string) (net.Conn, error) {
		client, server := net.Pipe()
		go func() {
			defer server.Close()
			line, err := bufio.NewReader(server).ReadBytes('\n')
			if err != nil {
				return
			}
			if requests != nil {
				requests <- append([]byte(nil), line...)
			}
			_ = writeAll(server, response)
		}()
		return client, nil
	}
}

func TestSnapshotClientLimiterRejectsNinthCaptureWithoutBlocking(t *testing.T) {
	limiter := newSnapshotClientLimiter(snapshotMaxConcurrent)
	releases := make([]func(), 0, snapshotMaxConcurrent)
	for i := 0; i < snapshotMaxConcurrent; i++ {
		release, err := limiter.acquire(context.Background())
		if err != nil {
			t.Fatalf("acquire slot %d: %v", i, err)
		}
		releases = append(releases, release)
	}
	defer func() {
		for _, release := range releases {
			release()
		}
	}()

	type acquireResult struct {
		release func()
		err     error
	}
	resultCh := make(chan acquireResult, 1)
	go func() {
		release, err := limiter.acquire(context.Background())
		resultCh <- acquireResult{release: release, err: err}
	}()

	select {
	case result := <-resultCh:
		if result.release != nil {
			result.release()
		}
		if !errors.Is(result.err, errSnapshotClientBusy) {
			t.Fatalf("ninth acquire error = %v, want client busy", result.err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("ninth acquire blocked instead of failing immediately")
	}
}

func TestSnapshotClientLimiterPrefersCanceledContextToBusy(t *testing.T) {
	for _, full := range []bool{false, true} {
		t.Run(fmt.Sprintf("full=%t", full), func(t *testing.T) {
			limiter := newSnapshotClientLimiter(1)
			if full {
				release, err := limiter.acquire(context.Background())
				if err != nil {
					t.Fatalf("fill limiter: %v", err)
				}
				defer release()
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			release, err := limiter.acquire(ctx)
			if release != nil {
				release()
			}
			if !errors.Is(err, context.Canceled) || errors.Is(err, errSnapshotClientBusy) {
				t.Fatalf("acquire error = %v, want context canceled", err)
			}
		})
	}
}

func TestCaptureSnapshotContextReturnsBusyBeforeDial(t *testing.T) {
	previousLimiter := snapshotCaptureLimiter
	limiter := newSnapshotClientLimiter(1)
	snapshotCaptureLimiter = limiter
	defer func() { snapshotCaptureLimiter = previousLimiter }()

	release, err := limiter.acquire(context.Background())
	if err != nil {
		t.Fatalf("fill global limiter: %v", err)
	}
	defer release()

	started := time.Now()
	header, jpeg, err := CaptureSnapshotContext(context.Background(), SnapshotRequest{})
	if !errors.Is(err, errSnapshotClientBusy) {
		t.Fatalf("capture error = %v, want client busy", err)
	}
	if header.OK || header.RetCode != -16 || jpeg != nil {
		t.Fatalf("unexpected busy response: header=%#v jpeg=%v", header, jpeg)
	}
	if time.Since(started) > 100*time.Millisecond {
		t.Fatal("busy capture attempted socket I/O instead of returning immediately")
	}
}

func TestCaptureSnapshotContextPrefersCanceledContextToBusy(t *testing.T) {
	for _, full := range []bool{false, true} {
		t.Run(fmt.Sprintf("full=%t", full), func(t *testing.T) {
			previousLimiter := snapshotCaptureLimiter
			limiter := newSnapshotClientLimiter(1)
			snapshotCaptureLimiter = limiter
			defer func() { snapshotCaptureLimiter = previousLimiter }()

			if full {
				release, err := limiter.acquire(context.Background())
				if err != nil {
					t.Fatalf("fill global limiter: %v", err)
				}
				defer release()
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			header, jpeg, err := CaptureSnapshotContext(ctx, SnapshotRequest{})
			if !errors.Is(err, context.Canceled) || errors.Is(err, errSnapshotClientBusy) {
				t.Fatalf("capture error = %v, want context canceled", err)
			}
			if header.OK || header.RetCode != -18 || jpeg != nil {
				t.Fatalf("unexpected canceled response: header=%#v jpeg=%v", header, jpeg)
			}
		})
	}
}

func TestCaptureSnapshotContextNormalizesCancellationAfterAcquire(t *testing.T) {
	client, server := net.Pipe()
	defer server.Close()
	requestRead := make(chan struct{})
	go func() {
		defer server.Close()
		_, _ = bufio.NewReader(server).ReadBytes('\n')
		close(requestRead)
		_, _ = io.Copy(io.Discard, server)
	}()
	dial := func(context.Context, string, string) (net.Conn, error) {
		return client, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan struct {
		header SnapshotHeader
		err    error
	}, 1)
	go func() {
		header, _, err := captureSnapshotContextWithDialer(
			ctx, SnapshotRequest{}, newSnapshotClientLimiter(1), dial)
		resultCh <- struct {
			header SnapshotHeader
			err    error
		}{header: header, err: err}
	}()

	select {
	case <-requestRead:
	case <-time.After(time.Second):
		t.Fatal("snapshot request was not written")
	}
	cancel()

	select {
	case result := <-resultCh:
		if !errors.Is(result.err, context.Canceled) {
			t.Fatalf("capture error = %v, want context canceled", result.err)
		}
		if result.header.OK || result.header.RetCode != -18 {
			t.Fatalf("unexpected canceled header: %#v", result.header)
		}
	case <-time.After(time.Second):
		t.Fatal("context cancellation did not interrupt capture")
	}
}

func TestCaptureSnapshotContextNormalizesDeadlineAfterAcquire(t *testing.T) {
	client, server := net.Pipe()
	defer server.Close()
	go func() {
		defer server.Close()
		_, _ = bufio.NewReader(server).ReadBytes('\n')
		_, _ = io.Copy(io.Discard, server)
	}()
	dial := func(context.Context, string, string) (net.Conn, error) {
		return client, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	header, jpeg, err := captureSnapshotContextWithDialer(
		ctx, SnapshotRequest{}, newSnapshotClientLimiter(1), dial)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("capture error = %v, want deadline exceeded", err)
	}
	if header.OK || header.RetCode != -13 || jpeg != nil {
		t.Fatalf("unexpected timeout response: header=%#v jpeg=%v", header, jpeg)
	}
}

func TestCaptureSnapshotContextUsesDedicatedSocketOnly(t *testing.T) {
	const dialFailure = "dedicated snapshot socket unavailable"
	var addresses []string
	dial := func(_ context.Context, network, address string) (net.Conn, error) {
		if network != "unix" {
			t.Fatalf("network = %q, want unix", network)
		}
		addresses = append(addresses, address)
		return nil, errors.New(dialFailure)
	}

	_, _, err := captureSnapshotContextWithDialer(
		context.Background(), SnapshotRequest{}, newSnapshotClientLimiter(1), dial)
	if err == nil || !strings.Contains(err.Error(), dialFailure) {
		t.Fatalf("capture error = %v, want dedicated socket failure", err)
	}
	if len(addresses) != 1 || addresses[0] != snapshotSocketPath {
		t.Fatalf("dial addresses = %v, want only %q", addresses, snapshotSocketPath)
	}
}

func TestBuildSocketSnapshotRequestUsesCaptureFreshProtocol(t *testing.T) {
	timeout := snapshotTimeoutMaxMS + 1000
	request := buildSocketSnapshotRequest(SnapshotRequest{
		Quality:   101,
		TimeoutMS: &timeout,
		X:         12,
		Y:         34,
		W:         640,
		H:         480,
	})

	if request.Version != snapshotProtocolVersion || request.Cmd != "capture_fresh" {
		t.Fatalf("unexpected protocol request: %#v", request)
	}
	if request.Quality != 100 || request.TimeoutMS != snapshotTimeoutMaxMS {
		t.Fatalf("unexpected normalized capture values: %#v", request)
	}
	if request.X != 12 || request.Y != 34 || request.W != 640 || request.H != 480 {
		t.Fatalf("ROI changed while building request: %#v", request)
	}
}

func TestBuildSocketSnapshotRequestPreservesDefaults(t *testing.T) {
	request := buildSocketSnapshotRequest(SnapshotRequest{})
	if request.Quality != 85 || request.TimeoutMS != defaultSnapshotTimeout {
		t.Fatalf("unexpected defaults: %#v", request)
	}
	if request.X != 0 || request.Y != 0 || request.W != 0 || request.H != 0 {
		t.Fatalf("unexpected default ROI: %#v", request)
	}
}

func TestWriteSnapshotRequestHandlesShortWritesAndAddsNewline(t *testing.T) {
	writer := &chunkWriter{limit: 3}
	request := buildSocketSnapshotRequest(SnapshotRequest{Quality: 76, X: 1, Y: 2, W: 3, H: 4})
	if err := writeSnapshotRequest(writer, request); err != nil {
		t.Fatalf("write request: %v", err)
	}

	payload := writer.buffer.Bytes()
	if len(payload) == 0 || payload[len(payload)-1] != '\n' {
		t.Fatalf("request is not newline framed: %q", payload)
	}

	var decoded socketSnapshotRequest
	if err := json.Unmarshal(bytes.TrimSuffix(payload, []byte{'\n'}), &decoded); err != nil {
		t.Fatalf("decode request: %v", err)
	}
	if decoded != request {
		t.Fatalf("request changed across short writes: got %#v want %#v", decoded, request)
	}
}

func TestReadSnapshotResponseAcceptsMetadataAndIgnoresUnknownFields(t *testing.T) {
	jpeg := []byte{0xff, 0xd8, 0x01, 0x02, 0xff, 0xd9}
	headerJSON := fmt.Sprintf(
		`{"version":1,"ok":true,"ret_code":0,"message":"ok","session_id":"session","capture_id":"session:9","captured_at":"2026-07-14T11:00:00+08:00","source":"mcp","width":1280,"height":720,"size":%d,"future_field":{"enabled":true}}`+"\n",
		len(jpeg),
	)

	header, received, err := readSnapshotResponse(io.MultiReader(strings.NewReader(headerJSON), bytes.NewReader(jpeg)))
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if header.CaptureID != "session:9" || header.Source != "mcp" || header.Width != 1280 || header.Height != 720 {
		t.Fatalf("unexpected header: %#v", header)
	}
	if !bytes.Equal(received, jpeg) {
		t.Fatalf("unexpected JPEG: %x", received)
	}
}

func TestReadSnapshotResponseRejectsInvalidSuccessHeaders(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*SnapshotHeader)
		want   string
	}{
		{
			name: "version",
			mutate: func(header *SnapshotHeader) {
				header.Version = 2
			},
			want: "unsupported snapshot protocol version",
		},
		{
			name: "negative size",
			mutate: func(header *SnapshotHeader) {
				header.Size = -1
			},
			want: "invalid snapshot JPEG size",
		},
		{
			name: "nonzero success ret code",
			mutate: func(header *SnapshotHeader) {
				header.RetCode = -14
			},
			want: "has ret_code",
		},
		{
			name: "missing session id",
			mutate: func(header *SnapshotHeader) {
				header.SessionID = ""
			},
			want: "missing session_id",
		},
		{
			name: "missing capture id",
			mutate: func(header *SnapshotHeader) {
				header.CaptureID = ""
			},
			want: "missing capture_id",
		},
		{
			name: "missing captured at",
			mutate: func(header *SnapshotHeader) {
				header.CapturedAt = ""
			},
			want: "missing captured_at",
		},
		{
			name: "missing source",
			mutate: func(header *SnapshotHeader) {
				header.Source = ""
			},
			want: "missing source",
		},
		{
			name: "invalid dimensions",
			mutate: func(header *SnapshotHeader) {
				header.Width = 0
			},
			want: "invalid snapshot dimensions",
		},
		{
			name: "zero size",
			mutate: func(header *SnapshotHeader) {
				header.Size = 0
			},
			want: "invalid snapshot JPEG size",
		},
		{
			name: "oversized",
			mutate: func(header *SnapshotHeader) {
				header.Size = snapshotJPEGMaxBytes + 1
			},
			want: "snapshot JPEG too large",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			header := successfulSnapshotHeader(1)
			test.mutate(&header)
			_, _, err := readSnapshotResponse(bytes.NewReader(encodeSnapshotResponse(t, header, nil)))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("got error %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestReadSnapshotResponseRejectsTruncatedJPEG(t *testing.T) {
	header := successfulSnapshotHeader(10)
	_, _, err := readSnapshotResponse(bytes.NewReader(encodeSnapshotResponse(t, header, []byte{1, 2, 3})))
	if err == nil || !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("got error %v, want unexpected EOF", err)
	}
}

func TestReadSnapshotResponseRejectsOversizedHeader(t *testing.T) {
	payload := bytes.Repeat([]byte{'x'}, snapshotHeaderMaxBytes+1)
	_, _, err := readSnapshotResponse(bytes.NewReader(payload))
	if err == nil || !strings.Contains(err.Error(), "snapshot header too large") {
		t.Fatalf("got error %v, want oversized header error", err)
	}
}

func TestReadSnapshotResponseDoesNotReadFailurePayload(t *testing.T) {
	header := SnapshotHeader{
		Version: snapshotProtocolVersion,
		OK:      false,
		RetCode: -16,
		Message: "snapshot busy",
		Size:    123,
	}
	headerPayload := encodeSnapshotResponse(t, header, nil)
	reader := &headerOnlyReader{header: headerPayload}

	got, jpeg, err := readSnapshotResponse(reader)
	if err != nil {
		t.Fatalf("read failure response: %v", err)
	}
	if got.OK || jpeg != nil {
		t.Fatalf("unexpected failure response: header=%#v jpeg=%v", got, jpeg)
	}
	if reader.reads != 1 {
		t.Fatalf("failure response performed %d reads, want exactly the header read", reader.reads)
	}
}

type headerOnlyReader struct {
	header []byte
	reads  int
}

func (r *headerOnlyReader) Read(payload []byte) (int, error) {
	r.reads++
	if r.reads > 1 {
		return 0, errors.New("unexpected payload read")
	}
	return copy(payload, r.header), nil
}

func TestCaptureSnapshotFromSocketUsesCaptureFresh(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "snapshot.sock")
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	requestCh := make(chan socketSnapshotRequest, 1)
	serverErr := make(chan error, 1)
	jpeg := []byte{0xff, 0xd8, 0x03, 0x04, 0xff, 0xd9}
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			serverErr <- err
			return
		}
		defer conn.Close()

		line, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			serverErr <- err
			return
		}
		var request socketSnapshotRequest
		if err := json.Unmarshal(line, &request); err != nil {
			serverErr <- err
			return
		}
		requestCh <- request

		header := successfulSnapshotHeader(len(jpeg))
		if err := writeAll(conn, encodeSnapshotResponse(t, header, jpeg)); err != nil {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	timeout := 1500
	dialer := &net.Dialer{Timeout: time.Second}
	header, received, err := captureSnapshotFromSocket(
		context.Background(),
		socketPath,
		SnapshotRequest{Quality: 77, TimeoutMS: &timeout, X: 10, Y: 20, W: 300, H: 200},
		dialer.DialContext,
	)
	if err != nil {
		t.Fatalf("capture: %v", err)
	}
	if !header.OK || !bytes.Equal(received, jpeg) {
		t.Fatalf("unexpected response: header=%#v jpeg=%x", header, received)
	}

	request := <-requestCh
	if request.Version != snapshotProtocolVersion || request.Cmd != "capture_fresh" ||
		request.Quality != 77 || request.TimeoutMS != timeout ||
		request.X != 10 || request.Y != 20 || request.W != 300 || request.H != 200 {
		t.Fatalf("unexpected wire request: %#v", request)
	}
	if err := <-serverErr; err != nil {
		t.Fatalf("snapshot server: %v", err)
	}
}

func TestCaptureSnapshotContextCancellationInterruptsIO(t *testing.T) {
	client, server := net.Pipe()
	defer server.Close()

	requestRead := make(chan struct{})
	go func() {
		defer server.Close()
		_, _ = bufio.NewReader(server).ReadBytes('\n')
		close(requestRead)
		_, _ = io.Copy(io.Discard, server)
	}()

	dial := func(context.Context, string, string) (net.Conn, error) {
		return client, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan error, 1)
	go func() {
		_, _, err := captureSnapshotFromSocket(ctx, "ignored", SnapshotRequest{}, dial)
		resultCh <- err
	}()

	select {
	case <-requestRead:
	case <-time.After(time.Second):
		t.Fatal("snapshot request was not written")
	}
	cancel()

	select {
	case err := <-resultCh:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("got error %v, want context canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("context cancellation did not interrupt snapshot IO")
	}
}

func TestCaptureSnapshotFromSocketAllowsConcurrentConnections(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "snapshot.sock")
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	const clients = 2
	accepted := make(chan struct{}, clients)
	release := make(chan struct{})
	serverErr := make(chan error, clients)
	for range clients {
		go func() {
			conn, err := listener.Accept()
			if err != nil {
				serverErr <- err
				return
			}
			defer conn.Close()
			if _, err := bufio.NewReader(conn).ReadBytes('\n'); err != nil {
				serverErr <- err
				return
			}
			accepted <- struct{}{}
			<-release
			jpeg := []byte{0xff, 0xd8, 0xff, 0xd9}
			serverErr <- writeAll(conn, encodeSnapshotResponse(t, successfulSnapshotHeader(len(jpeg)), jpeg))
		}()
	}

	dialer := &net.Dialer{Timeout: time.Second}
	var wg sync.WaitGroup
	results := make(chan error, clients)
	for range clients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := captureSnapshotFromSocket(
				context.Background(), socketPath, SnapshotRequest{}, dialer.DialContext)
			results <- err
		}()
	}

	for range clients {
		select {
		case <-accepted:
		case <-time.After(time.Second):
			close(release)
			t.Fatal("snapshot clients were serialized")
		}
	}
	close(release)
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatalf("concurrent capture: %v", err)
		}
	}
	for range clients {
		if err := <-serverErr; err != nil {
			t.Fatalf("snapshot server: %v", err)
		}
	}
}
