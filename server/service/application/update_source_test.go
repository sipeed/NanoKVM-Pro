package application

import "testing"

func TestNormalizeUpdateSourceURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "directory", raw: " https://example.test/releases/ ", want: "https://example.test/releases"},
		{name: "credentials are not persisted", raw: "https://user:secret@example.test/releases", wantErr: true},
		{name: "empty resets", raw: "", want: ""},
		{name: "manifest is not root", raw: "https://example.test/nanokvm_pro_latest.json", wantErr: true},
		{name: "query is not allowed", raw: "https://example.test/releases?token=x", wantErr: true},
		{name: "file protocol is not allowed", raw: "file:///tmp/releases", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := normalizeUpdateSourceURL(test.raw)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeUpdateSourceURL() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("normalizeUpdateSourceURL() = %q, want %q", got, test.want)
			}
		})
	}
}
