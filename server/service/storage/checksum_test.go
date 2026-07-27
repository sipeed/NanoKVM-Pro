package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestResolveImagePath(t *testing.T) {
	root := t.TempDir()
	image := filepath.Join(root, "test.iso")
	if err := os.WriteFile(image, []byte("image"), 0o644); err != nil {
		t.Fatal(err)
	}

	resolved, err := resolveImagePath(image, []string{root})
	if err != nil {
		t.Fatalf("resolve image path: %v", err)
	}
	if resolved != image {
		t.Fatalf("resolved path = %q, want %q", resolved, image)
	}

	outside := filepath.Join(t.TempDir(), "outside.iso")
	if err := os.WriteFile(outside, []byte("image"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := resolveImagePath(outside, []string{root}); err == nil {
		t.Fatal("expected path outside root to be rejected")
	}

	unsupported := filepath.Join(root, "test.bin")
	if err := os.WriteFile(unsupported, []byte("image"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := resolveImagePath(unsupported, []string{root}); err == nil {
		t.Fatal("expected unsupported extension to be rejected")
	}
}

func TestChecksumImageAlgorithmsAndCache(t *testing.T) {
	root := t.TempDir()
	image := filepath.Join(root, "test.iso")
	cachePath := filepath.Join(root, "cache", "checksums.json")
	if err := os.WriteFile(image, []byte("NanoKVM-Pro"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := map[string]string{
		"md5":    "ecaeb431c31559e86fa341f8ba8268f2",
		"sha1":   "62fc27f1df288415e89bf72799d99730291d64b1",
		"sha256": "a6f6fbc3166a855a97abc73e3fab78ef6b2653e5832b5d9cc4bb72b98caa8c67",
	}

	for algorithm, expected := range tests {
		actual, err := checksumImage(image, algorithm, cachePath)
		if err != nil {
			t.Fatalf("checksum %s: %v", algorithm, err)
		}
		if actual != expected {
			t.Errorf("%s checksum = %q, want %q", algorithm, actual, expected)
		}
	}

	cache, err := loadChecksumCache(cachePath)
	if err != nil {
		t.Fatalf("load checksum cache: %v", err)
	}
	entry := cache.Entries[image]
	if entry == nil || len(entry.Checksums) != len(tests) {
		t.Fatalf("cache entry = %#v, want %d checksums", entry, len(tests))
	}
}

func TestChecksumImageInvalidatesChangedFile(t *testing.T) {
	root := t.TempDir()
	image := filepath.Join(root, "test.iso")
	cachePath := filepath.Join(root, "checksums.json")
	if err := os.WriteFile(image, []byte("first"), 0o644); err != nil {
		t.Fatal(err)
	}

	first, err := checksumImage(image, "sha256", cachePath)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(image, []byte("second"), 0o644); err != nil {
		t.Fatal(err)
	}
	modified := time.Now().Add(time.Second)
	if err := os.Chtimes(image, modified, modified); err != nil {
		t.Fatal(err)
	}

	second, err := checksumImage(image, "sha256", cachePath)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("checksum cache was not invalidated after the image changed")
	}
}

func TestNewImageHasherRejectsUnsupportedAlgorithm(t *testing.T) {
	if _, err := newImageHasher("sha512"); err == nil {
		t.Fatal("expected unsupported algorithm to be rejected")
	}
}
