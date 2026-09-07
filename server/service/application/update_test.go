package application

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadUsesIndependentSafetyLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("package"))
	}))
	defer server.Close()

	target := filepath.Join(t.TempDir(), "package.tar.gz")
	if err := download(server.URL, target); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "package" {
		t.Fatalf("downloaded content = %q", content)
	}
}
