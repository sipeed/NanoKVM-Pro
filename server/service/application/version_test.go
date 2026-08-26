package application

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"NanoKVM-Server/proto"
	"github.com/gin-gonic/gin"
)

type updateRoundTripFunc func(*http.Request) (*http.Response, error)

func (fn updateRoundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func TestGetVersionReportsCustomSourceFailure(t *testing.T) {
	previousPath := updateSourcePath
	previousClient := updateHTTPClient
	updateSourcePath = t.TempDir() + "/update-source.json"
	updateHTTPClient = &http.Client{Transport: updateRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       io.NopCloser(strings.NewReader("unavailable")),
			Header:     make(http.Header),
		}, nil
	})}
	t.Cleanup(func() {
		updateSourcePath = previousPath
		updateHTTPClient = previousClient
	})

	if err := saveUpdateSourceConfig(UpdateSourceConfig{Enabled: true, URL: "https://updates.example"}); err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	(&Service{}).GetVersion(context)

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"code":-1`) {
		t.Fatalf("unexpected response: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestGetVersionKeepsOfficialFallbackOnSourceFailure(t *testing.T) {
	previousPath := updateSourcePath
	previousClient := updateHTTPClient
	updateSourcePath = t.TempDir() + "/update-source.json"
	updateHTTPClient = &http.Client{Transport: updateRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       io.NopCloser(strings.NewReader("unavailable")),
			Header:     make(http.Header),
		}, nil
	})}
	t.Cleanup(func() {
		updateSourcePath = previousPath
		updateHTTPClient = previousClient
	})

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	(&Service{}).GetVersion(context)

	var response struct {
		Code int                 `json:"code"`
		Data proto.GetVersionRsp `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Code != 0 || response.Data.Current != response.Data.Latest {
		t.Fatalf("unexpected response: status=%d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestValidateLatestTreatsSizeAsAdvisory(t *testing.T) {
	latest := &Latest{
		Version: "1.2.15",
		Name:    "nanokvm_pro_1.2.15.tar.gz",
		Sha512:  base64.StdEncoding.EncodeToString(make([]byte, sha512.Size)),
		Size:    uint64(maxManualPackageBytes) + 1,
	}
	if err := validateLatest(latest); err != nil {
		t.Fatalf("validateLatest() rejected advisory size: %v", err)
	}
}
