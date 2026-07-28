package storage

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"hash"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"NanoKVM-Server/proto"
)

const (
	checksumCacheVersion = 1
	checksumBufferSize   = 1024 * 1024
	checksumCacheFile    = "/etc/kvm/cache/image-checksums.json"
)

var checksumMutex sync.Mutex

type checksumCache struct {
	Version int                            `json:"version"`
	Entries map[string]*checksumCacheEntry `json:"entries"`
}

type checksumCacheEntry struct {
	Size       int64             `json:"size"`
	ModifiedAt int64             `json:"modifiedAt"`
	Device     uint64            `json:"device,omitempty"`
	Inode      uint64            `json:"inode,omitempty"`
	Checksums  map[string]string `json:"checksums"`
}

func (s *Service) ChecksumImage(c *gin.Context) {
	var req proto.ChecksumImageReq
	var rsp proto.Response

	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}

	path, err := resolveImagePath(req.File, []string{imageDirectory, sdCardDirectory})
	if err != nil {
		log.Errorf("invalid image path %s: %v", req.File, err)
		rsp.ErrRsp(c, -2, "invalid image")
		return
	}

	if !checksumMutex.TryLock() {
		rsp.ErrRsp(c, -4, "checksum in progress")
		return
	}
	defer checksumMutex.Unlock()

	checksum, err := checksumImage(path, req.Algorithm, checksumCacheFile)
	if err != nil {
		log.Errorf("failed to checksum image %s: %v", path, err)
		rsp.ErrRsp(c, -3, "checksum image failed")
		return
	}

	rsp.OkRspWithData(c, &proto.ChecksumImageRsp{
		File:      req.File,
		Algorithm: req.Algorithm,
		Checksum:  checksum,
	})
}

func resolveImagePath(path string, roots []string) (string, error) {
	if !filepath.IsAbs(path) {
		return "", errors.New("image path must be absolute")
	}

	resolved, err := filepath.EvalSymlinks(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	resolved, err = filepath.Abs(resolved)
	if err != nil {
		return "", err
	}

	insideRoot := false
	for _, root := range roots {
		absoluteRoot, err := filepath.Abs(root)
		if err != nil {
			continue
		}

		relative, err := filepath.Rel(absoluteRoot, resolved)
		if err == nil && relative != "." && relative != ".." &&
			!strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			insideRoot = true
			break
		}
	}
	if !insideRoot {
		return "", errors.New("image path is outside an allowed directory")
	}

	extension := strings.ToLower(filepath.Ext(resolved))
	if extension != ".iso" && extension != ".img" {
		return "", errors.New("unsupported image extension")
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", errors.New("image is not a regular file")
	}

	return resolved, nil
}

func checksumImage(path string, algorithm string, cachePath string) (string, error) {
	hasher, err := newImageHasher(algorithm)
	if err != nil {
		return "", err
	}

	before, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	cache, err := loadChecksumCache(cachePath)
	if err != nil {
		log.Warnf("failed to load image checksum cache: %v", err)
		cache = newChecksumCache()
	}

	if entry, ok := cache.Entries[path]; ok && cacheEntryMatches(entry, before) {
		if checksum := entry.Checksums[algorithm]; checksum != "" {
			return checksum, nil
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	opened, err := file.Stat()
	if err != nil {
		return "", err
	}
	if !sameImageFile(before, opened) {
		return "", errors.New("image changed before checksum started")
	}

	buffer := make([]byte, checksumBufferSize)
	if _, err := io.CopyBuffer(hasher, file, buffer); err != nil {
		return "", err
	}

	afterPath, err := os.Stat(path)
	if err != nil {
		return "", errors.New("image changed during checksum")
	}
	afterFile, err := file.Stat()
	if err != nil {
		return "", err
	}
	if !sameImageFile(before, afterPath) || !sameImageFile(before, afterFile) {
		return "", errors.New("image changed during checksum")
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))
	entry, ok := cache.Entries[path]
	if !ok || !cacheEntryMatches(entry, afterPath) {
		device, inode := fileIdentity(afterPath)
		entry = &checksumCacheEntry{
			Size:       afterPath.Size(),
			ModifiedAt: afterPath.ModTime().UnixNano(),
			Device:     device,
			Inode:      inode,
			Checksums:  make(map[string]string),
		}
		cache.Entries[path] = entry
	}
	if entry.Checksums == nil {
		entry.Checksums = make(map[string]string)
	}
	entry.Checksums[algorithm] = checksum

	if err := saveChecksumCache(cachePath, cache); err != nil {
		log.Warnf("failed to save image checksum cache: %v", err)
	}

	return checksum, nil
}

func newImageHasher(algorithm string) (hash.Hash, error) {
	switch algorithm {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	default:
		return nil, errors.New("unsupported checksum algorithm")
	}
}

func sameImageFile(left os.FileInfo, right os.FileInfo) bool {
	return os.SameFile(left, right) && left.Size() == right.Size() &&
		left.ModTime().UnixNano() == right.ModTime().UnixNano()
}

func cacheEntryMatches(entry *checksumCacheEntry, info os.FileInfo) bool {
	device, inode := fileIdentity(info)
	return entry != nil && entry.Size == info.Size() &&
		entry.ModifiedAt == info.ModTime().UnixNano() &&
		entry.Device == device && entry.Inode == inode
}

func newChecksumCache() *checksumCache {
	return &checksumCache{
		Version: checksumCacheVersion,
		Entries: make(map[string]*checksumCacheEntry),
	}
}

func loadChecksumCache(path string) (*checksumCache, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return newChecksumCache(), nil
	}
	if err != nil {
		return nil, err
	}

	var cache checksumCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	if cache.Version != checksumCacheVersion {
		return newChecksumCache(), nil
	}
	if cache.Entries == nil {
		cache.Entries = make(map[string]*checksumCacheEntry)
	}

	return &cache, nil
}

func saveChecksumCache(path string, cache *checksumCache) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	temp, err := os.CreateTemp(filepath.Dir(path), ".image-checksums-*.tmp")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if err := temp.Chmod(0o644); err != nil {
		_ = temp.Close()
		return err
	}
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	if err := os.Rename(tempPath, path); err == nil {
		return nil
	} else if runtime.GOOS != "windows" {
		return err
	}

	// Windows does not replace an existing destination with os.Rename. The
	// device runs Linux, where the rename above is atomic; this fallback keeps
	// local development and tests functional.
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.Rename(tempPath, path)
}
