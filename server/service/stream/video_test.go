package stream

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"NanoKVM-Server/common"
	"NanoKVM-Server/service/ai"

	"github.com/gin-gonic/gin"
	"github.com/pion/rtcp"
)

type fakeVideo struct {
	reads         atomic.Int32
	closes        atomic.Int32
	idrs          atomic.Int32
	controlResult atomic.Int32
	quality       atomic.Uint32
	bitrate       atomic.Uint32
	next          chan int
}

func (f *fakeVideo) read(result int) ([]byte, int) {
	f.reads.Add(1)
	select {
	case result = <-f.next:
	default:
	}
	if result < 0 {
		return nil, result
	}
	return []byte{byte(result)}, result
}
func (f *fakeVideo) ReadMjpeg(_, _, quality uint16) ([]byte, int) { return f.read(0) }
func (f *fakeVideo) ReadH264(_, _, quality uint16) ([]byte, int)  { return f.read(3) }
func (f *fakeVideo) ReadH265(_, _, quality uint16) ([]byte, int)  { return f.read(7) }
func (f *fakeVideo) CloseVideoStream()                            { f.closes.Add(1) }
func (f *fakeVideo) RequestIDR() int                              { f.idrs.Add(1); return 0 }
func (f *fakeVideo) SetBitrate(value uint16) int {
	f.bitrate.Store(uint32(value))
	return int(f.controlResult.Load())
}
func (f *fakeVideo) SetMjpegQuality(value uint16) int {
	f.quality.Store(uint32(value))
	return int(f.controlResult.Load())
}
func (f *fakeVideo) SetRateControl(uint8) int { return int(f.controlResult.Load()) }
func (f *fakeVideo) SetFps(uint8) int         { return int(f.controlResult.Load()) }
func (f *fakeVideo) SetGop(uint8) int         { return int(f.controlResult.Load()) }

func testHub() (*VideoHub, *fakeVideo) {
	f := &fakeVideo{next: make(chan int, 20)}
	screen := &common.Screen{ScreenState: common.ScreenState{StreamType: common.STREAM_TYPE_H264_WEBRTC, RateControlMode: 1, BitRate: 8000, Quality: 80, GOP: 50}}
	return newVideoHub(f, screen), f
}

func eventually(t *testing.T, predicate func() bool) {
	t.Helper()
	deadline := time.After(time.Second)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for !predicate() {
		select {
		case <-deadline:
			t.Fatal("condition timed out")
		case <-ticker.C:
		}
	}
}

func receive(t *testing.T, sub *VideoSubscription) VideoFrame {
	t.Helper()
	select {
	case frame := <-sub.Frames:
		return frame
	case <-time.After(time.Second):
		t.Fatal("no frame")
		return VideoFrame{}
	}
}

func TestVideoBackpressureStopsReadsAndControlsRemainAvailable(t *testing.T) {
	h, f := testHub()
	sub := h.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
	defer sub.Close()
	eventually(t, func() bool { return f.reads.Load() == 1 })
	time.Sleep(60 * time.Millisecond)
	if got := f.reads.Load(); got != 1 {
		t.Fatalf("blocked writer absorbed %d frames; expected one in-flight frame and no queue", got)
	}
	controlled := make(chan int, 1)
	go func() { controlled <- h.SetQuality(12000) }()
	select {
	case result := <-controlled:
		if result != 0 {
			t.Fatal(result)
		}
	case <-time.After(time.Second):
		t.Fatal("network handoff held the control mutex")
	}
	first := receive(t, sub)
	eventually(t, func() bool { return f.reads.Load() == 2 })
	time.Sleep(40 * time.Millisecond)
	if f.reads.Load() != 2 {
		t.Fatal("capture ran ahead of the writer")
	}
	second := receive(t, sub)
	if second.Timestamp <= first.Timestamp {
		t.Fatal("Direct timestamps did not advance")
	}
	sub.Close()
	eventually(t, func() bool { h.mutex.Lock(); defer h.mutex.Unlock(); return !h.running })
	if f.closes.Load() == 0 {
		t.Fatal("last client did not release IPC")
	}
}

func TestVideoMultiClientFanoutAndBlockedClientCancellation(t *testing.T) {
	h, f := testHub()
	a := h.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
	b := h.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
	defer a.Close()
	defer b.Close()
	var first, second VideoFrame
	select {
	case first = <-a.Frames:
		second = receive(t, b)
	case first = <-b.Frames:
		second = receive(t, a)
	case <-time.After(time.Second):
		t.Fatal("fanout timed out")
	}
	if first.Timestamp != second.Timestamp || first.Data[0] != second.Data[0] {
		t.Fatal("peers received different captures")
	}
	eventually(t, func() bool { return f.reads.Load() >= 2 })
	before := f.reads.Load()
	time.Sleep(40 * time.Millisecond)
	if f.reads.Load() != before {
		t.Fatal("fanout did not backpressure capture")
	}
	b.Close()
	if !receive(t, a).KeyFrame {
		t.Fatal("remaining peer lost its frame")
	}
	a.Close()
	eventually(t, func() bool { h.mutex.Lock(); defer h.mutex.Unlock(); return !h.running })
}

func TestVideoCancellationReleasesIPCWhileNetworkWriterIsBusy(t *testing.T) {
	h, f := testHub()
	sub := h.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
	done := make(chan struct{})
	stop := sub.CloseOn(done)
	defer stop()
	receive(t, sub) // The writer has taken a frame and now blocks on the network.
	eventually(t, func() bool { return f.reads.Load() == 2 })
	close(done)
	eventually(t, func() bool { return f.closes.Load() > 0 })
	eventually(t, func() bool { h.mutex.Lock(); defer h.mutex.Unlock(); return !h.running })
	if f.reads.Load() != 2 {
		t.Fatal("capture continued after the writer was cancelled")
	}
}

func TestVideoModeSwitchCancelsBlockedOldHandoff(t *testing.T) {
	h, f := testHub()
	old := h.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
	defer old.Close()
	eventually(t, func() bool { return f.reads.Load() == 1 })
	h.SetMode(common.STREAM_TYPE_H265_DIRECT)
	current := h.Subscribe(common.STREAM_TYPE_H265_DIRECT)
	defer current.Close()
	frame := receive(t, current)
	if frame.Result != 7 || !frame.KeyFrame || !h.IsCurrent(frame) {
		t.Fatalf("bad first HEVC frame: %+v", frame)
	}
	closed := f.closes.Load()
	old.Close()
	if f.closes.Load() != closed {
		t.Fatal("old peer closed the new reader")
	}
	// A late ICE callback may subscribe, but must not steal the selected mode.
	late := h.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
	defer late.Close()
	if !h.IsMode(common.STREAM_TYPE_H265_DIRECT) {
		t.Fatal("late subscriber changed mode")
	}
	h.SetMode(common.STREAM_TYPE_MJPEG)
	jpeg := h.Subscribe(common.STREAM_TYPE_MJPEG)
	defer jpeg.Close()
	if frame := receive(t, jpeg); frame.Result != 0 || !frame.KeyFrame {
		t.Fatal("bad JPEG recovery")
	}
}

func TestVideoErrorAndWrongCodecWaitForNewKeyframe(t *testing.T) {
	for _, script := range [][]int{{-3, 4, 4, 3}, {7, 4, 4, 3}, {-1, 4, 3}} {
		h, f := testHub()
		for _, result := range script {
			f.next <- result
		}
		sub := h.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
		frame := receive(t, sub)
		if frame.Result != 3 || !frame.KeyFrame || f.reads.Load() < int32(len(script)) {
			t.Fatal("broken prediction chain escaped to transport")
		}
		sub.Close()
	}
}

func TestRTCPRequestsAreNonblockingAndCoalesced(t *testing.T) {
	h, f := testHub()
	sub := h.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
	defer sub.Close()
	eventually(t, func() bool { return f.reads.Load() == 1 })
	finished := make(chan struct{})
	go func() {
		for range 10000 {
			h.RequestKeyFrame(common.STREAM_TYPE_H264_WEBRTC)
		}
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("RTCP blocked behind capture/writer")
	}
	if f.idrs.Load() != 1 {
		t.Fatal("RTCP called controls directly")
	}
	time.Sleep(110 * time.Millisecond)
	receive(t, sub)
	eventually(t, func() bool { return f.idrs.Load() == 2 })
	if !needsKeyFrame([]rtcp.Packet{&rtcp.ReceiverReport{}, &rtcp.PictureLossIndication{}, &rtcp.FullIntraRequest{}}) {
		t.Fatal("PLI/FIR ignored")
	}
	if !needsKeyFrame([]rtcp.Packet{&rtcp.FullIntraRequest{}}) {
		t.Fatal("FIR ignored")
	}
	if needsKeyFrame([]rtcp.Packet{&rtcp.ReceiverReport{}, &rtcp.TransportLayerNack{}}) {
		t.Fatal("ordinary RTCP requested IDR")
	}
}

func TestStreamHTTPControlsCommitOnlySuccessfulSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, f := testHub()
	service := &Service{video: func() *VideoHub { return h }}
	router := gin.New()
	router.POST("/quality", service.SetQuality)
	router.POST("/fps", service.SetFps)
	router.POST("/gop", service.SetGop)
	router.POST("/rc", service.SetRateControlMode)
	post := func(path, form string) int {
		t.Helper()
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		router.ServeHTTP(recorder, request)
		var response struct {
			Code int `json:"code"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		return response.Code
	}
	for _, value := range []string{"1", "100", "101", "20000"} {
		if post("/quality", "quality="+value) != 0 {
			t.Fatalf("legacy quality range rejected: %s", value)
		}
	}
	if f.quality.Load() != 100 || f.bitrate.Load() != 20000 {
		t.Fatal("quality/bitrate mapped to wrong controls")
	}
	if post("/fps", "fps=120") != 0 || post("/gop", "gop=200") != 0 || post("/rc", "mode=cbr") != 0 {
		t.Fatal("valid control failed")
	}
	before := h.screen.Snapshot()
	f.controlResult.Store(-1)
	for _, pair := range [][2]string{{"/quality", "quality=70"}, {"/quality", "quality=12000"}, {"/fps", "fps=30"}, {"/gop", "gop=50"}, {"/rc", "mode=vbr"}} {
		if post(pair[0], pair[1]) == 0 {
			t.Fatal("control failure was hidden")
		}
	}
	if after := h.screen.Snapshot(); before != after {
		t.Fatal("failed controls changed confirmed UI state")
	}
	for _, pair := range [][2]string{{"/quality", "quality=20001"}, {"/quality", "quality=-1"}, {"/fps", "fps=121"}, {"/gop", "gop=201"}, {"/rc", "mode=bad"}} {
		if post(pair[0], pair[1]) == 0 {
			t.Fatal("invalid argument accepted")
		}
	}
}

func TestStreamSnapshotReturnsJPEGWithoutChangingMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := testHub()
	previous := captureJPEG
	t.Cleanup(func() { captureJPEG = previous })
	captureJPEG = func(_ context.Context, req ai.SnapshotRequest) (ai.SnapshotHeader, []byte, error) {
		if req.Quality != 70 {
			t.Fatalf("quality = %d, want 70", req.Quality)
		}
		return ai.SnapshotHeader{OK: true, Width: 1920, Height: 1080}, []byte{0xff, 0xd8, 0xff}, nil
	}

	service := &Service{video: func() *VideoHub { return h }}
	router := gin.New()
	router.GET("/snapshot", service.Snapshot)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/snapshot?quality=70", nil))
	if recorder.Code != http.StatusOK {
		t.Fatal(recorder.Body.String())
	}
	if recorder.Header().Get("Content-Type") != "image/jpeg" {
		t.Fatal(recorder.Header().Get("Content-Type"))
	}
	if recorder.Header().Get("X-Width") != "1920" || recorder.Header().Get("X-Height") != "1080" {
		t.Fatal(recorder.Header())
	}
	if !bytes.Equal(recorder.Body.Bytes(), []byte{0xff, 0xd8, 0xff}) {
		t.Fatal("missing JPEG body")
	}
	if !h.IsMode(common.STREAM_TYPE_H264_WEBRTC) {
		t.Fatal("snapshot changed the live stream mode")
	}
}

func TestStreamSnapshotRejectsInvalidQualityAndCaptureErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := testHub()
	previous := captureJPEG
	t.Cleanup(func() { captureJPEG = previous })
	captureJPEG = func(context.Context, ai.SnapshotRequest) (ai.SnapshotHeader, []byte, error) {
		return ai.SnapshotHeader{OK: false, RetCode: -16, Message: "busy"}, nil, errors.New("busy")
	}

	service := &Service{video: func() *VideoHub { return h }}
	router := gin.New()
	router.GET("/snapshot", service.Snapshot)

	bad := httptest.NewRecorder()
	router.ServeHTTP(bad, httptest.NewRequest(http.MethodGet, "/snapshot?quality=0", nil))
	if bad.Code != http.StatusOK {
		t.Fatal(bad.Code)
	}
	var invalid struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(bad.Body.Bytes(), &invalid); err != nil || invalid.Code != -1 {
		t.Fatal(bad.Body.String())
	}

	failed := httptest.NewRecorder()
	router.ServeHTTP(failed, httptest.NewRequest(http.MethodGet, "/snapshot", nil))
	if failed.Code != http.StatusServiceUnavailable {
		t.Fatal(failed.Code)
	}
	var unavailable struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(failed.Body.Bytes(), &unavailable); err != nil || unavailable.Code != -16 || unavailable.Msg != "busy" {
		t.Fatal(failed.Body.String())
	}
}

func TestVideoMjpegBorrowRestoresSelectedMode(t *testing.T) {
	h, f := testHub()
	live := h.Subscribe(common.STREAM_TYPE_H264_WEBRTC)
	defer live.Close()
	eventually(t, func() bool { return f.reads.Load() == 1 })
	h.Borrow(common.STREAM_TYPE_MJPEG)
	if h.Selected() != common.STREAM_TYPE_H264_WEBRTC {
		t.Fatal("borrow changed selected mode")
	}
	if h.screen.Snapshot().StreamType != common.STREAM_TYPE_H264_WEBRTC {
		t.Fatal("borrow persisted MJPEG in screen state")
	}
	mjpeg := h.Subscribe(common.STREAM_TYPE_MJPEG)
	if frame := receive(t, mjpeg); frame.Result != 0 || !frame.KeyFrame {
		t.Fatal("bad JPEG frame")
	}
	select {
	case status := <-live.Updates:
		if status != -4 {
			t.Fatalf("status %d, want -4", status)
		}
	case <-time.After(time.Second):
		t.Fatal("H.264 viewer was not paused")
	}
	mjpeg.Close()
	eventually(t, func() bool { return h.IsMode(common.STREAM_TYPE_H264_WEBRTC) })
	if receive(t, live).Result != 3 {
		t.Fatal("H.264 did not resume")
	}
}

func TestVideoExplicitMjpegModeDoesNotRestore(t *testing.T) {
	h, _ := testHub()
	h.SetMode(common.STREAM_TYPE_MJPEG)
	sub := h.Subscribe(common.STREAM_TYPE_MJPEG)
	receive(t, sub)
	sub.Close()
	eventually(t, func() bool { h.mutex.Lock(); defer h.mutex.Unlock(); return !h.running })
	if !h.IsMode(common.STREAM_TYPE_MJPEG) || h.Selected() != common.STREAM_TYPE_MJPEG {
		t.Fatal("explicit MJPEG selection was rolled back")
	}
}

func TestStreamGetModeReportsSelectedAndActive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := testHub()
	service := &Service{video: func() *VideoHub { return h }}
	router := gin.New()
	router.GET("/mode", service.GetMode)
	h.Borrow(common.STREAM_TYPE_MJPEG)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/mode", nil))
	var response struct {
		Code int `json:"code"`
		Data struct {
			Mode   string `json:"mode"`
			Active string `json:"active"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Code != 0 || response.Data.Mode != "h264-webrtc" || response.Data.Active != "mjpeg" {
		t.Fatal(recorder.Body.String())
	}
}
