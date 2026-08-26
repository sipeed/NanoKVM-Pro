package application

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ulikunitz/xz"
	"golang.org/x/crypto/blake2b"
)

const (
	maxManualPackageBytes   int64 = 512 << 20
	maxArchiveEntries             = 1024
	maxExpandedArchiveBytes       = uint64(1 << 30)
	maxArchiveFileBytes           = uint64(512 << 20)
	maxArchivePathLength          = 512
	maxMetadataBytes              = 1 << 20
)

type artifactKind string

const (
	artifactApplication artifactKind = "application"
	artifactFirmware    artifactKind = "firmware"
)

var (
	artifactVersionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+-]{0,127}$`)
	appPackagePattern      = regexp.MustCompile(`^nanokvm_pro_([A-Za-z0-9][A-Za-z0-9._+-]{0,127})$`)
	firmwareVersionPattern = regexp.MustCompile(`v?[0-9]+\.[0-9]+\.[0-9]+`)
)

type artifactInspection struct {
	Kind          artifactKind
	Version       string
	ExpandedBytes uint64
	Application   *applicationArchive
}

type archiveEntry struct {
	Name     string
	Typeflag byte
	Size     int64
	SHA512   string
	BLAKE2b  string
	Content  []byte
}

type archiveLayout struct {
	Entries       map[string]archiveEntry
	ExpandedBytes uint64
}

type applicationArchive struct {
	Root    string
	Version string
	Layout  archiveLayout
}

// A single archive inspection selects the application or firmware path; raw
// images are rejected by the same generic unsupported-package error family.
func inspectUpdateArtifact(artifactPath, originalName string) (artifactInspection, error) {
	lowerName := strings.ToLower(strings.TrimSpace(originalName))
	if strings.HasSuffix(lowerName, ".axp") {
		return artifactInspection{}, errors.New("AXP images are not supported by in-device update")
	}
	if strings.HasSuffix(lowerName, ".img") {
		return artifactInspection{}, errors.New("raw IMG images are not supported by in-device update")
	}
	layout, err := inspectTarArchive(artifactPath)
	if err != nil {
		return artifactInspection{}, err
	}
	if app, err := inspectApplicationArchive(layout); err == nil {
		return artifactInspection{
			Kind:          artifactApplication,
			Version:       app.Version,
			ExpandedBytes: layout.ExpandedBytes,
			Application:   &app,
		}, nil
	}
	if version, err := inspectFirmwareArchive(layout); err == nil {
		return artifactInspection{
			Kind:          artifactFirmware,
			Version:       version,
			ExpandedBytes: layout.ExpandedBytes,
		}, nil
	}

	return artifactInspection{}, errors.New("unsupported update package layout")
}

func inspectTarArchive(archivePath string) (archiveLayout, error) {
	reader, closeReader, err := openUpdateArchive(archivePath)
	if err != nil {
		return archiveLayout{}, err
	}
	defer closeReader()

	return readArchiveLayout(reader)
}

// openUpdateArchive selects gzip, xz, or plain tar without treating ZIP/raw
// image bytes as an installable update archive.
func openUpdateArchive(archivePath string) (io.Reader, func(), error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, nil, err
	}
	closeFile := func() { _ = file.Close() }

	var magic [6]byte
	read, err := io.ReadFull(file, magic[:])
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		closeFile()
		return nil, nil, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		closeFile()
		return nil, nil, err
	}
	if read >= 4 && bytes.Equal(magic[:4], []byte{'P', 'K', 3, 4}) {
		closeFile()
		return nil, nil, errors.New("AXP and ZIP images are not supported by in-device update")
	}
	if read >= 2 && magic[0] == 0x1f && magic[1] == 0x8b {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			closeFile()
			return nil, nil, fmt.Errorf("open gzip update archive: %w", err)
		}
		return gzipReader, func() {
			_ = gzipReader.Close()
			closeFile()
		}, nil
	}
	if read == len(magic) && bytes.Equal(magic[:], []byte{0xfd, '7', 'z', 'X', 'Z', 0x00}) {
		xzReader, err := xz.NewReader(file)
		if err != nil {
			closeFile()
			return nil, nil, fmt.Errorf("open xz update archive: %w", err)
		}
		return xzReader, closeFile, nil
	}

	// Plain tar is accepted for development packages, but raw .img files do not
	// have a valid tar layout and are rejected below as unsupported packages.
	return file, closeFile, nil
}

func readArchiveLayout(reader io.Reader) (archiveLayout, error) {
	// Hash entries while reading and retain only bounded metadata needed for
	// later format-specific validation.
	layout := archiveLayout{Entries: make(map[string]archiveEntry)}
	capturedApplicationMetadata := make(map[string]struct{}, len(appNames))
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return archiveLayout{}, fmt.Errorf("read update archive: %w", err)
		}
		name, err := validateArchivePath(header.Name)
		if err != nil {
			return archiveLayout{}, err
		}
		if _, exists := layout.Entries[name]; exists {
			return archiveLayout{}, fmt.Errorf("duplicate archive entry %q", name)
		}
		if len(layout.Entries) >= maxArchiveEntries {
			return archiveLayout{}, fmt.Errorf("update archive has more than %d entries", maxArchiveEntries)
		}
		if header.Typeflag != tar.TypeDir && !isRegularArchiveType(header.Typeflag) {
			return archiveLayout{}, fmt.Errorf("unsupported archive entry type for %q", name)
		}
		if header.Typeflag == tar.TypeDir && header.Size != 0 {
			return archiveLayout{}, fmt.Errorf("invalid archive directory %q", name)
		}
		entry := archiveEntry{Name: name, Typeflag: header.Typeflag, Size: header.Size}
		if isRegularArchiveType(header.Typeflag) {
			if header.Size < 0 {
				return archiveLayout{}, errors.New("negative archive file size")
			}
			size := uint64(header.Size)
			if size > maxArchiveFileBytes {
				return archiveLayout{}, fmt.Errorf("update archive member exceeds %d bytes", maxArchiveFileBytes)
			}
			if size > maxExpandedArchiveBytes-layout.ExpandedBytes {
				return archiveLayout{}, fmt.Errorf("update archive expanded size exceeds %d bytes", maxExpandedArchiveBytes)
			}
			layout.ExpandedBytes += size
			captureContent, appName := shouldCaptureArchiveMetadata(name, header.Size)
			if appName != "" {
				if _, alreadyCaptured := capturedApplicationMetadata[appName]; alreadyCaptured {
					captureContent = false
				} else {
					capturedApplicationMetadata[appName] = struct{}{}
				}
			}
			entry.SHA512, entry.BLAKE2b, entry.Content, err = digestArchiveEntry(tarReader, header.Size, name, captureContent)
			if err != nil {
				return archiveLayout{}, err
			}
		}
		layout.Entries[name] = entry
	}
	return layout, nil
}

// shouldCaptureArchiveMetadata permits only the metadata needed by the two
// supported formats. The application name check means no more than one JSON
// document for each of the three application components can be retained.
func shouldCaptureArchiveMetadata(name string, size int64) (bool, string) {
	if size > maxMetadataBytes {
		return false, ""
	}
	if name == "b2sum.txt" || name == "overlay/boot/ver" {
		return true, ""
	}
	parts := strings.Split(name, "/")
	if len(parts) != 2 {
		return false, ""
	}
	matches := appPackagePattern.FindStringSubmatch(parts[0])
	if len(matches) != 2 {
		return false, ""
	}
	for _, appName := range appNames {
		if parts[1] == appName+"_"+matches[1]+".json" {
			return true, appName
		}
	}
	return false, ""
}

func digestArchiveEntry(reader io.Reader, size int64, name string, captureContent bool) (string, string, []byte, error) {
	sha := sha512.New()
	blake, err := blake2b.New512(nil)
	if err != nil {
		return "", "", nil, err
	}
	var content bytes.Buffer
	var output io.Writer = io.MultiWriter(sha, blake)
	if captureContent {
		output = io.MultiWriter(sha, blake, &content)
	}
	written, err := io.CopyN(output, reader, size)
	if err != nil {
		return "", "", nil, fmt.Errorf("read archive entry %q: %w", name, err)
	}
	if written != size {
		return "", "", nil, fmt.Errorf("truncated archive entry %q", name)
	}
	return base64.StdEncoding.EncodeToString(sha.Sum(nil)), hex.EncodeToString(blake.Sum(nil)), content.Bytes(), nil
}

func validateArchivePath(raw string) (string, error) {
	// Reject absolute, traversal, separator, and control-character paths before
	// any archive member can be written to the filesystem.
	if raw == "." || raw == "./" {
		return ".", nil
	}
	if raw == "" || len(raw) > maxArchivePathLength || path.IsAbs(raw) || strings.Contains(raw, "\\") {
		return "", fmt.Errorf("invalid archive path %q", raw)
	}
	for _, character := range raw {
		if character < 32 || character == 127 {
			return "", fmt.Errorf("invalid archive path %q", raw)
		}
	}
	name := raw
	for strings.HasPrefix(name, "./") {
		name = strings.TrimPrefix(name, "./")
	}
	name = strings.TrimRight(name, "/")
	if name == "" {
		return ".", nil
	}
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return "", fmt.Errorf("invalid archive path %q", raw)
		}
	}
	if path.IsAbs(name) {
		return "", fmt.Errorf("invalid archive path %q", raw)
	}
	return name, nil
}

func isRegularArchiveType(typeFlag byte) bool {
	return typeFlag == tar.TypeReg || typeFlag == tar.TypeRegA
}

func inspectApplicationArchive(layout archiveLayout) (applicationArchive, error) {
	// Application bundles must contain exactly the expected DEB/JSON pairs under
	// one versioned top-level directory.
	var root, version string
	for _, entry := range layout.Entries {
		if !isRegularArchiveType(entry.Typeflag) {
			continue
		}
		parts := strings.Split(entry.Name, "/")
		if len(parts) != 2 || parts[0] == "" {
			return applicationArchive{}, errors.New("invalid application archive layout")
		}
		if root == "" {
			root = parts[0]
			matches := appPackagePattern.FindStringSubmatch(root)
			if len(matches) != 2 || !validArtifactVersion(matches[1]) {
				return applicationArchive{}, errors.New("invalid application package root")
			}
			version = matches[1]
		}
		if parts[0] != root {
			return applicationArchive{}, errors.New("application archive must have one top-level directory")
		}
	}
	if root == "" {
		return applicationArchive{}, errors.New("empty application archive")
	}
	if rootEntry, ok := layout.Entries[root]; !ok || rootEntry.Typeflag != tar.TypeDir {
		return applicationArchive{}, errors.New("application archive is missing its top-level directory")
	}

	expected := make(map[string]struct{}, len(appNames)*2)
	for _, appName := range appNames {
		expected[appName+"_"+version+"_arm64.deb"] = struct{}{}
		expected[appName+"_"+version+".json"] = struct{}{}
	}
	for name, entry := range layout.Entries {
		if !isRegularArchiveType(entry.Typeflag) {
			if name != root {
				return applicationArchive{}, fmt.Errorf("unexpected application archive directory %q", name)
			}
			continue
		}
		parts := strings.Split(name, "/")
		if len(parts) != 2 || parts[0] != root {
			return applicationArchive{}, errors.New("invalid application archive layout")
		}
		if _, ok := expected[parts[1]]; !ok {
			return applicationArchive{}, fmt.Errorf("unexpected application archive file %q", parts[1])
		}
	}
	for fileName := range expected {
		entry, ok := layout.Entries[root+"/"+fileName]
		if !ok || !isRegularArchiveType(entry.Typeflag) || entry.Size == 0 {
			return applicationArchive{}, fmt.Errorf("application archive missing %q", fileName)
		}
	}

	if err := validateApplicationManifestChecksums(layout, root, version); err != nil {
		return applicationArchive{}, err
	}
	return applicationArchive{Root: root, Version: version, Layout: layout}, nil
}

func validateApplicationManifestChecksums(layout archiveLayout, root, version string) error {
	for _, appName := range appNames {
		jsonEntry := layout.Entries[root+"/"+appName+"_"+version+".json"]
		if len(jsonEntry.Content) == 0 {
			return fmt.Errorf("application metadata %s is too large or missing", appName)
		}
		var info FileInfo
		if err := json.Unmarshal(jsonEntry.Content, &info); err != nil {
			return fmt.Errorf("decode application metadata %s: %w", appName, err)
		}
		debName := appName + "_" + version + "_arm64.deb"
		debEntry := layout.Entries[root+"/"+debName]
		if info.Version != version {
			return fmt.Errorf("application metadata version mismatch for %s", appName)
		}
		if info.Name != "" && info.Name != debName {
			return fmt.Errorf("application metadata name mismatch for %s", appName)
		}
		if info.Size > 0 && int64(info.Size) != debEntry.Size {
			return fmt.Errorf("application metadata size mismatch for %s", appName)
		}
		if err := validateSHA512String(info.SHA512); err != nil {
			return fmt.Errorf("invalid application metadata checksum for %s", appName)
		}
		if info.SHA512 != debEntry.SHA512 {
			return fmt.Errorf("application checksum mismatch for %s", appName)
		}
	}
	return nil
}

func validateSHA512String(value string) error {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil || len(decoded) != sha512.Size {
		return errors.New("invalid sha512")
	}
	return nil
}

func inspectFirmwareArchive(layout archiveLayout) (string, error) {
	// Firmware archives are restricted to the files understood by kvmcomm and a
	// complete b2sum manifest covering every regular member.
	required := []string{
		"firmware/u-boot_signed.bin",
		"firmware/boot_signed.bin",
		"firmware/AX630C_emmc_arm64_k419_sipeed_nanokvm_signed.dtb",
		"overlay/boot/ver",
		"b2sum.txt",
	}
	for _, name := range required {
		entry, ok := layout.Entries[name]
		if !ok || !isRegularArchiveType(entry.Typeflag) || entry.Size == 0 {
			return "", fmt.Errorf("firmware archive missing %q", name)
		}
	}
	for name, entry := range layout.Entries {
		if err := validateFirmwareArchiveMember(name, isRegularArchiveType(entry.Typeflag)); err != nil {
			return "", err
		}
	}
	if err := validateFirmwareChecksums(layout); err != nil {
		return "", err
	}
	version := firmwareVersionPattern.FindString(string(layout.Entries["overlay/boot/ver"].Content))
	if version == "" {
		return "", errors.New("invalid firmware version")
	}
	return version, nil
}

func validateFirmwareArchiveMember(name string, isFile bool) error {
	if name == "." {
		if isFile {
			return errors.New("firmware archive root must be a directory")
		}
		return nil
	}
	if name == "b2sum.txt" {
		if !isFile {
			return errors.New("firmware checksum manifest must be a regular file")
		}
		return nil
	}
	for _, required := range []string{
		"firmware/u-boot_signed.bin",
		"firmware/boot_signed.bin",
		"firmware/AX630C_emmc_arm64_k419_sipeed_nanokvm_signed.dtb",
		"overlay/boot/ver",
	} {
		if name == required {
			if !isFile {
				return fmt.Errorf("required firmware archive member is not a regular file: %q", name)
			}
			return nil
		}
	}
	if strings.HasPrefix(name, "overlay/") {
		return nil
	}
	if strings.HasPrefix(name, "firmware/") {
		return nil
	}
	if !isFile && (name == "firmware" || name == "overlay") {
		return nil
	}
	return fmt.Errorf("unexpected firmware archive member %q", name)
}

func validateFirmwareChecksums(layout archiveLayout) error {
	sums := layout.Entries["b2sum.txt"].Content
	if len(sums) == 0 || len(sums) > maxMetadataBytes {
		return errors.New("invalid firmware checksum manifest")
	}
	covered := make(map[string]struct{})
	for _, line := range strings.Split(strings.TrimSpace(string(sums)), "\n") {
		line = strings.TrimSuffix(line, "\r")
		if len(line) < 131 || line[128] != ' ' {
			return errors.New("invalid firmware checksum manifest")
		}
		expected, err := hex.DecodeString(line[:128])
		if err != nil || len(expected) != blake2b.Size {
			return errors.New("invalid firmware checksum manifest")
		}
		if line[129] != ' ' && line[129] != '*' {
			return errors.New("invalid firmware checksum manifest")
		}
		name := line[130:]
		name, err = validateArchivePath(name)
		if err != nil || name == "b2sum.txt" {
			return errors.New("invalid firmware checksum file name")
		}
		entry, ok := layout.Entries[name]
		if !ok || !isRegularArchiveType(entry.Typeflag) || entry.BLAKE2b != strings.ToLower(line[:128]) {
			return fmt.Errorf("firmware checksum mismatch for %q", name)
		}
		if _, duplicate := covered[name]; duplicate {
			return fmt.Errorf("duplicate firmware checksum entry for %q", name)
		}
		covered[name] = struct{}{}
	}
	for name, entry := range layout.Entries {
		if isRegularArchiveType(entry.Typeflag) && name != "b2sum.txt" {
			if _, ok := covered[name]; !ok {
				return fmt.Errorf("firmware checksum manifest does not cover %q", name)
			}
		}
	}
	return nil
}

func validArtifactVersion(version string) bool {
	return artifactVersionPattern.MatchString(version)
}

func extractApplicationArchive(archivePath, destination string, app applicationArchive) (string, error) {
	// Re-read and compare each member against the inspection snapshot to detect
	// replacement between upload inspection and installation.
	if err := os.MkdirAll(destination, 0o700); err != nil {
		return "", fmt.Errorf("create update extraction directory: %w", err)
	}
	reader, closeReader, err := openUpdateArchive(archivePath)
	if err != nil {
		return "", err
	}
	defer closeReader()

	tarReader := tar.NewReader(reader)
	seen := make(map[string]struct{})
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("read application archive: %w", err)
		}
		name, err := validateArchivePath(header.Name)
		if err != nil {
			return "", err
		}
		wanted, ok := app.Layout.Entries[name]
		if !ok || wanted.Typeflag != header.Typeflag || wanted.Size != header.Size {
			return "", errors.New("application archive changed after inspection")
		}
		if _, duplicate := seen[name]; duplicate {
			return "", fmt.Errorf("duplicate archive entry %q", name)
		}
		seen[name] = struct{}{}
		if header.Typeflag != tar.TypeDir && !isRegularArchiveType(header.Typeflag) {
			return "", fmt.Errorf("unsupported archive entry type for %q", name)
		}
		outputPath := filepath.Join(destination, filepath.FromSlash(name))
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(outputPath, 0o700); err != nil {
				return "", fmt.Errorf("create application archive directory: %w", err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(outputPath), 0o700); err != nil {
			return "", fmt.Errorf("create application archive parent: %w", err)
		}
		output, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			return "", fmt.Errorf("create application archive file: %w", err)
		}
		written, copyErr := io.CopyN(output, tarReader, header.Size)
		closeErr := output.Close()
		if copyErr != nil || written != header.Size {
			return "", fmt.Errorf("write application archive file %q: %w", name, copyErr)
		}
		if closeErr != nil {
			return "", fmt.Errorf("close application archive file: %w", closeErr)
		}
	}
	if len(seen) != len(app.Layout.Entries) {
		return "", errors.New("application archive changed after inspection")
	}
	return filepath.Join(destination, app.Root), nil
}

func validateApplicationDebMetadata(root, version string) error {
	// Read one dpkg-deb field per command: the multi-field form prefixes values
	// with labels, which would make a naïve Fields-based comparison incorrect.
	for _, appName := range appNames {
		debPath := filepath.Join(root, appName+"_"+version+"_arm64.deb")
		fields := make([]string, 0, 3)
		for _, field := range []string{"Package", "Version", "Architecture"} {
			output, err := exec.Command("dpkg-deb", "-f", debPath, field).Output()
			if err != nil {
				return fmt.Errorf("inspect %s package metadata: %w", appName, err)
			}
			fields = append(fields, strings.TrimSpace(string(output)))
		}
		expectedPackage := appName
		if appName == "nanokvmpro" {
			expectedPackage = "nanokvm"
		}
		if len(fields) != 3 || fields[0] != expectedPackage || fields[1] != version || fields[2] != "arm64" {
			return fmt.Errorf("invalid Debian package metadata for %s", appName)
		}
	}
	return nil
}
