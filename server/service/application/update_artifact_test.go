package application

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/blake2b"
)

type testArchiveEntry struct {
	name     string
	data     []byte
	typeflag byte
}

func TestInspectApplicationArchive(t *testing.T) {
	version := "1.2.15"
	root := "nanokvm_pro_" + version
	entries := []testArchiveEntry{{name: root, typeflag: tar.TypeDir}}
	for _, app := range appNames {
		debName := app + "_" + version + "_arm64.deb"
		deb := []byte("not parsed during archive inspection: " + app)
		hash := sha512.Sum512(deb)
		metadata, err := json.Marshal(FileInfo{
			Version: version,
			Name:    debName,
			SHA512:  base64.StdEncoding.EncodeToString(hash[:]),
			Size:    len(deb),
		})
		if err != nil {
			t.Fatal(err)
		}
		entries = append(entries,
			testArchiveEntry{name: root + "/" + debName, data: deb, typeflag: tar.TypeReg},
			testArchiveEntry{name: root + "/" + app + "_" + version + ".json", data: metadata, typeflag: tar.TypeReg},
		)
	}
	path := writeGzipTar(t, entries)
	inspection, err := inspectUpdateArtifact(path, "bundle.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Kind != artifactApplication || inspection.Version != version {
		t.Fatalf("unexpected inspection: %#v", inspection)
	}
}

func TestInspectFirmwareArchiveChecksEveryFile(t *testing.T) {
	files := map[string][]byte{
		"firmware/u-boot_signed.bin":                                []byte("uboot"),
		"firmware/boot_signed.bin":                                  []byte("kernel"),
		"firmware/AX630C_emmc_arm64_k419_sipeed_nanokvm_signed.dtb": []byte("dtb"),
		"overlay/boot/ver":                                          []byte("custom build: v1.0.15\n"),
		"overlay/etc/example.conf":                                  []byte("enabled=true\n"),
	}
	entries := []testArchiveEntry{
		{name: "firmware", typeflag: tar.TypeDir},
		{name: "overlay", typeflag: tar.TypeDir},
		{name: "overlay/boot", typeflag: tar.TypeDir},
	}
	checksumLines := make([]byte, 0)
	for name, data := range files {
		hash := blake2b.Sum512(data)
		checksumLines = append(checksumLines, []byte(hex.EncodeToString(hash[:])+"  "+name+"\n")...)
		entries = append(entries, testArchiveEntry{name: name, data: data, typeflag: tar.TypeReg})
	}
	entries = append(entries, testArchiveEntry{name: "b2sum.txt", data: checksumLines, typeflag: tar.TypeReg})

	path := writeGzipTar(t, entries)
	inspection, err := inspectUpdateArtifact(path, "firmware.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Kind != artifactFirmware || inspection.Version != "v1.0.15" {
		t.Fatalf("unexpected inspection: %#v", inspection)
	}
}

func TestFirmwareRejectsUnexpectedFirmwareAndDuplicateChecksum(t *testing.T) {
	files := map[string][]byte{
		"firmware/u-boot_signed.bin":                                []byte("uboot"),
		"firmware/boot_signed.bin":                                  []byte("kernel"),
		"firmware/AX630C_emmc_arm64_k419_sipeed_nanokvm_signed.dtb": []byte("dtb"),
		"firmware/unexpected.bin":                                   []byte("unexpected"),
		"overlay/boot/ver":                                          []byte("v1.0.15\n"),
	}
	entries := firmwareTestEntries(files, false)
	path := writeGzipTar(t, entries)
	if inspection, err := inspectUpdateArtifact(path, "firmware.tar.gz"); err != nil || inspection.Kind != artifactFirmware {
		t.Fatalf("extra firmware member should remain compatible, got %#v, %v", inspection, err)
	}

	delete(files, "firmware/unexpected.bin")
	entries = firmwareTestEntries(files, true)
	path = writeGzipTar(t, entries)
	if _, err := inspectUpdateArtifact(path, "firmware.tar.gz"); err == nil {
		t.Fatal("duplicate b2sum entry must be rejected")
	}
}

func TestValidateArchivePathMatchesFirmwareHelper(t *testing.T) {
	valid := map[string]string{
		"./firmware/u-boot_signed.bin": "firmware/u-boot_signed.bin",
		"overlay/boot/":                "overlay/boot",
		"././overlay/boot/ver":         "overlay/boot/ver",
		"./":                           ".",
	}
	for raw, want := range valid {
		got, err := validateArchivePath(raw)
		if err != nil || got != want {
			t.Fatalf("validateArchivePath(%q) = %q, %v; want %q", raw, got, err, want)
		}
	}
	for _, raw := range []string{"firmware//boot.bin", "./../firmware/boot.bin", `firmware\\boot.bin`, "overlay/\x01bad"} {
		if _, err := validateArchivePath(raw); err == nil {
			t.Fatalf("validateArchivePath(%q) unexpectedly succeeded", raw)
		}
	}
}

func TestArchiveOnlyCachesRequiredMetadata(t *testing.T) {
	entries := []testArchiveEntry{
		{name: "nanokvm_pro_1.2.15", typeflag: tar.TypeDir},
		{name: "nanokvm_pro_1.2.15/nanokvmpro_1.2.15.json", data: []byte("{}"), typeflag: tar.TypeReg},
		{name: "nanokvm_pro_1.2.15/unrelated.json", data: bytes.Repeat([]byte("x"), 1024), typeflag: tar.TypeReg},
	}
	path := writeGzipTar(t, entries)
	layout, err := inspectTarArchive(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(layout.Entries["nanokvm_pro_1.2.15/nanokvmpro_1.2.15.json"].Content) == 0 {
		t.Fatal("required application metadata was not retained")
	}
	if len(layout.Entries["nanokvm_pro_1.2.15/unrelated.json"].Content) != 0 {
		t.Fatal("unrelated metadata must not be retained in memory")
	}
}

func TestInspectArchiveRejectsTraversalAndUnsupportedImage(t *testing.T) {
	path := writeGzipTar(t, []testArchiveEntry{{name: "../outside", data: []byte("bad"), typeflag: tar.TypeReg}})
	if _, err := inspectUpdateArtifact(path, "bundle.tar.gz"); err == nil {
		t.Fatal("expected traversal archive to be rejected")
	}
	if _, err := inspectUpdateArtifact(path, "system.img"); err == nil {
		t.Fatal("expected IMG to be rejected")
	}
}

func writeGzipTar(t *testing.T, entries []testArchiveEntry) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "update.tar.gz")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)
	for _, entry := range entries {
		typeflag := entry.typeflag
		if typeflag == 0 {
			typeflag = tar.TypeReg
		}
		header := &tar.Header{Name: entry.name, Typeflag: typeflag, Mode: 0o600}
		if typeflag != tar.TypeDir {
			header.Size = int64(len(entry.data))
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatal(err)
		}
		if len(entry.data) > 0 {
			if _, err := tarWriter.Write(entry.data); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

func firmwareTestEntries(files map[string][]byte, duplicateFirstChecksum bool) []testArchiveEntry {
	entries := []testArchiveEntry{
		{name: "firmware", typeflag: tar.TypeDir},
		{name: "overlay", typeflag: tar.TypeDir},
		{name: "overlay/boot", typeflag: tar.TypeDir},
	}
	checksumLines := make([]byte, 0)
	for name, data := range files {
		hash := blake2b.Sum512(data)
		line := hex.EncodeToString(hash[:]) + "  " + name + "\n"
		checksumLines = append(checksumLines, []byte(line)...)
		if duplicateFirstChecksum {
			checksumLines = append(checksumLines, []byte(line)...)
			duplicateFirstChecksum = false
		}
		entries = append(entries, testArchiveEntry{name: name, data: data, typeflag: tar.TypeReg})
	}
	entries = append(entries, testArchiveEntry{name: "b2sum.txt", data: checksumLines, typeflag: tar.TypeReg})
	return entries
}

func TestValidateFirmwareChecksumRejectsMalformedLine(t *testing.T) {
	layout := archiveLayout{Entries: map[string]archiveEntry{
		"b2sum.txt": {Name: "b2sum.txt", Typeflag: tar.TypeReg, Content: bytes.Repeat([]byte("0"), 128)},
	}}
	if err := validateFirmwareChecksums(layout); err == nil {
		t.Fatal("expected malformed manifest rejection")
	}
}
