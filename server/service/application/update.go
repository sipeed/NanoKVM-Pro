package application

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const maxTries = 3

var (
	updateMu         sync.Mutex
	updateInProgress bool
)

// acquireUpdate serializes online and manual installation paths and makes the
// release operation idempotent for error and goroutine cleanup paths.
func acquireUpdate() (func(), error) {
	updateMu.Lock()
	defer updateMu.Unlock()
	if updateInProgress {
		return nil, errors.New("update already in progress")
	}
	updateInProgress = true
	var once sync.Once
	return func() {
		once.Do(func() {
			updateMu.Lock()
			updateInProgress = false
			updateMu.Unlock()
		})
	}, nil
}

// Update runs the online application update under the shared update lock.
func (s *Service) Update(c *gin.Context) {
	var rsp proto.Response
	releaseLock, err := acquireUpdate()
	if err != nil {
		rsp.ErrRsp(c, -1, err.Error())
		return
	}
	defer releaseLock()

	if err := update(); err != nil {
		log.Errorf("failed to update application: %v", err)
		rsp.ErrRsp(c, -2, "failed to update service")
		return
	}

	rsp.OkRsp(c)
	log.Debug("application update succeeded")
}

// update is the online application path. It deliberately uses the same safe
// package inspection and Debian installer as a manually uploaded package.
func update() error {
	// Query the manifest before creating a workspace so a failed source check
	// does not leave a cache directory behind.
	latest, err := getLatest()
	if err != nil {
		return err
	}
	workspace, err := newOnlineUpdateWorkspace()
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(workspace) }()

	// Download into a private workspace; the manifest size is advisory, while
	// downloadUpdatePackage enforces the fixed local byte limit.
	_ = sendMessage("download", 0)
	packagePath := filepath.Join(workspace, "application.tar.gz")
	if err := download(latest.Url, packagePath); err != nil {
		return err
	}
	if err := checksum(packagePath, latest.Sha512); err != nil {
		return fmt.Errorf("check application package checksum: %w", err)
	}

	// Inspect the archive before extracting it, then require its version to
	// match the manifest selected above.
	inspection, err := inspectUpdateArtifact(packagePath, latest.Name)
	if err != nil {
		return err
	}
	if inspection.Kind != artifactApplication || inspection.Version != latest.Version {
		return errors.New("application package does not match its manifest")
	}
	if err := ensureFreeSpaceAt(workspace, inspection.ExpandedBytes+minimumFreeUpdateSpace); err != nil {
		return err
	}
	// Extract only the validated application members into the same workspace.
	_ = sendMessage("install", 0)
	root, err := extractApplicationArchive(packagePath, filepath.Join(workspace, "extracted"), *inspection.Application)
	if err != nil {
		return err
	}
	if err := install(root, inspection.Version); err != nil {
		return err
	}
	_ = sendMessage("restart", 0)
	return nil
}

func newOnlineUpdateWorkspace() (string, error) {
	if err := os.MkdirAll(TempDir, 0o700); err != nil {
		return "", fmt.Errorf("create update cache directory: %w", err)
	}
	workspace, err := os.MkdirTemp(TempDir, "application-update-*")
	if err != nil {
		return "", fmt.Errorf("create update workspace: %w", err)
	}
	if err := os.Chmod(workspace, 0o700); err != nil {
		_ = os.RemoveAll(workspace)
		return "", fmt.Errorf("set update workspace permissions: %w", err)
	}
	return workspace, nil
}

func download(rawURL, target string) error {
	for attempt := 0; attempt < maxTries; attempt++ {
		if attempt > 0 {
			log.Debugf("retry application download %d/%d", attempt+1, maxTries)
			time.Sleep(time.Duration(attempt) * time.Second)
		}
		if err := downloadUpdatePackage(rawURL, target, maxManualPackageBytes, func(length int64) error {
			return ensureFreeSpaceAt(filepath.Dir(target), uint64(length)+minimumFreeUpdateSpace)
		}); err != nil {
			log.Errorf("download application package failed: %v", err)
			_ = os.Remove(target)
			continue
		}
		info, err := os.Stat(target)
		if err != nil {
			return err
		}
		if info.Size() == 0 {
			_ = os.Remove(target)
			return errors.New("application package is empty")
		}
		return nil
	}
	return errors.New("downloading application package failed")
}

func install(dir string, version string) error {
	// Validate every package and checksum before invoking dpkg for any package.
	if err := validateApplicationInstallFiles(dir, version); err != nil {
		return err
	}

	for index, appName := range appNames {
		_ = sendMessage("install", index*30+10)
		filePath := filepath.Join(dir, fmt.Sprintf("%s_%s_arm64.deb", appName, version))
		installed := false
		var packageInstallErr error
		for attempt := 0; attempt < maxTries; attempt++ {
			if attempt > 0 {
				log.Debugf("retry installing %s %d/%d", appName, attempt+1, maxTries)
				time.Sleep(time.Duration(attempt) * time.Second)
			}
			cmd := exec.Command("dpkg", "-i", filePath)
			cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
			outputBuffer := &limitedCommandOutput{limit: 8 << 10}
			cmd.Stdout = outputBuffer
			cmd.Stderr = outputBuffer
			err := cmd.Run()
			if err != nil {
				packageInstallErr = fmt.Errorf("install %s: %w: %s", appName, err, outputBuffer.String())
				log.Errorf("%v", packageInstallErr)
				continue
			}
			installed = true
			log.Infof("installed %s", appName)
			break
		}
		if !installed {
			// Do not continue with later packages after a failure. This still cannot
			// make dpkg transactional, but it prevents intentionally advancing the
			// rest of a partially applied application bundle.
			if packageInstallErr == nil {
				packageInstallErr = fmt.Errorf("install %s failed", appName)
			}
			return packageInstallErr
		}
	}
	_ = sendMessage("install", 100)
	return nil
}

func validateApplicationInstallFiles(dir string, version string) error {
	for _, appName := range appNames {
		jsonPath := filepath.Join(dir, fmt.Sprintf("%s_%s.json", appName, version))
		debPath := filepath.Join(dir, fmt.Sprintf("%s_%s_arm64.deb", appName, version))
		fileInfo, err := getFileInfo(jsonPath)
		if err != nil {
			return fmt.Errorf("read %s metadata: %w", appName, err)
		}
		if fileInfo.Version != version || (fileInfo.Name != "" && fileInfo.Name != filepath.Base(debPath)) {
			return fmt.Errorf("invalid %s metadata", appName)
		}
		if err := validateSHA512String(fileInfo.SHA512); err != nil {
			return fmt.Errorf("invalid %s checksum", appName)
		}
		if err := checksum(debPath, fileInfo.SHA512); err != nil {
			return fmt.Errorf("check %s checksum: %w", appName, err)
		}
	}
	// Debian metadata is checked last so filenames and manifest hashes are
	// validated before asking dpkg-deb to inspect package control fields.
	return validateApplicationDebMetadata(dir, version)
}

// truncateUpdateOutput bounds installer diagnostics stored in a job response.
func truncateUpdateOutput(output string) string {
	const maxOutput = 4096
	if len(output) > maxOutput {
		return output[:maxOutput]
	}
	return output
}
