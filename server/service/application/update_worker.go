package application

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	applicationUpdateUnit       = "nanokvm-self-update.service"
	applicationUpdateWorkerFlag = "--application-update-worker"
)

var (
	errUpdateInProgress   = errors.New("update already in progress")
	updateLockPath        = "/run/lock/nanokvm-application-update.lock"
	updateWorkerInspector = inspectUpdateWorker
	updateWorkerStarter   = startUpdateWorker
)

type updateWorkerInfo struct {
	Active bool
	JobID  string
}

// RunUpdateWorkerCommand is dispatched before hardware or HTTP initialization.
// systemd owns this process independently of nanokvm.service and its restarts.
func RunUpdateWorkerCommand(arguments []string) (bool, int) {
	if len(arguments) == 0 || arguments[0] != applicationUpdateWorkerFlag {
		return false, 0
	}
	if len(arguments) != 2 {
		return true, 2
	}
	if _, err := uuid.Parse(arguments[1]); err != nil {
		return true, 2
	}
	id := arguments[1]
	var release func()
	var err error
	// Type=exec returns once this image is running. The HTTP handler releases
	// its admission lock immediately after systemd-run returns.
	deadline := time.Now().Add(10 * time.Second)
	for {
		release, err = acquireUpdateFileLock()
		if err == nil || time.Now().After(deadline) {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		_ = failUpdateJob(id, err)
		log.Errorf("acquire independent update lock: %v", err)
		return true, 1
	}
	defer release()
	if err := runManualUpdateJob(id); err != nil {
		log.Errorf("independent update job %s: %v", id, err)
		return true, 1
	}
	return true, 0
}

func acquireUpdateFileLock() (func(), error) {
	file, err := os.OpenFile(updateLockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open update lock: %w", err)
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = file.Close()
		return nil, errUpdateInProgress
	}
	return func() { _ = file.Close() }, nil
}

func updateWorkerArguments(executable, id string) []string {
	return []string{
		"--unit=" + applicationUpdateUnit,
		"--collect",
		"--service-type=exec",
		"--property=RuntimeMaxSec=46min",
		"--setenv=NANOKVM_UPDATE_JOB_ID=" + id,
		"--setenv=DEBIAN_FRONTEND=noninteractive",
		executable, applicationUpdateWorkerFlag, id,
	}
}

func startUpdateWorker(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid update job")
	}
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	output, err := exec.Command("systemd-run", updateWorkerArguments(executable, id)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("start independent update: %w: %s", err, truncateUpdateOutput(string(output)))
	}
	return nil
}

func inspectUpdateWorker() (updateWorkerInfo, error) {
	if _, err := os.Stat("/run/systemd/system"); os.IsNotExist(err) {
		return updateWorkerInfo{}, nil
	}
	output, err := exec.Command("systemctl", "show", applicationUpdateUnit,
		"--property=LoadState", "--property=ActiveState", "--property=Environment").CombinedOutput()
	return parseUpdateWorkerInfo(string(output), err)
}

func parseUpdateWorkerInfo(output string, commandErr error) (updateWorkerInfo, error) {
	var info updateWorkerInfo
	for _, line := range strings.Split(output, "\n") {
		if line == "LoadState=not-found" {
			return info, nil
		}
		if strings.HasPrefix(line, "ActiveState=") {
			switch strings.TrimPrefix(line, "ActiveState=") {
			case "active", "activating", "reloading", "deactivating":
				info.Active = true
			}
		}
		if strings.HasPrefix(line, "Environment=") {
			for _, entry := range strings.Fields(strings.TrimPrefix(line, "Environment=")) {
				entry = strings.Trim(entry, "\"")
				if strings.HasPrefix(entry, "NANOKVM_UPDATE_JOB_ID=") {
					info.JobID = strings.TrimPrefix(entry, "NANOKVM_UPDATE_JOB_ID=")
				}
			}
		}
	}
	if commandErr != nil {
		return info, fmt.Errorf("inspect independent update: %w", commandErr)
	}
	return info, nil
}

func failUpdateJob(id string, cause error) error {
	job, err := loadManualJob(id)
	if err != nil {
		return err
	}
	job.State = "failed"
	job.Error = safeUpdateError(cause)
	job.UpdatedAt = time.Now().UTC()
	return saveManualJob(job)
}

func reportUpdateJob(id, step string, progress int) error {
	job, err := loadManualJob(id)
	if err != nil {
		return err
	}
	job.Step, job.Progress = step, progress
	job.UpdatedAt = time.Now().UTC()
	return saveManualJob(job)
}

func stageOnlineUpdate(path string, inspection artifactInspection) (manualUpdateJob, error) {
	if err := ensureUpdateJobRoot(); err != nil {
		return manualUpdateJob{}, err
	}
	id := uuid.NewString()
	if err := os.Mkdir(manualJobDirectory(id), 0o700); err != nil {
		return manualUpdateJob{}, err
	}
	success := false
	defer func() {
		if !success {
			_ = removeManualJob(id)
		}
	}()
	source, err := os.Open(path)
	if err != nil {
		return manualUpdateJob{}, err
	}
	defer source.Close()
	target, err := os.OpenFile(manualArtifactPath(id), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return manualUpdateJob{}, err
	}
	size, copyErr := io.Copy(target, source)
	if copyErr == nil {
		copyErr = target.Sync()
	}
	closeErr := target.Close()
	if err := errors.Join(copyErr, closeErr); err != nil {
		return manualUpdateJob{}, err
	}
	checksum, err := sha256File(manualArtifactPath(id))
	if err != nil {
		return manualUpdateJob{}, err
	}
	now := time.Now().UTC()
	job := manualUpdateJob{
		ID: id, Kind: artifactApplication, Version: inspection.Version,
		SHA256: checksum, Size: size, State: "installing", Step: "install",
		CreatedAt: now, UpdatedAt: now, ExpiresAt: now.Add(manualUpdateJobLifetime),
	}
	if err := saveManualJob(job); err != nil {
		return manualUpdateJob{}, err
	}
	success = true
	return job, nil
}

// GetUpdateStatus allows a reconnected browser to recover the current job.
func (s *Service) GetUpdateStatus(c *gin.Context) {
	var rsp proto.Response
	if err := recoverInterruptedManualUpdates(); err != nil {
		rsp.ErrRsp(c, -1, "failed to inspect update status")
		return
	}
	entries, err := os.ReadDir(updateJobRoot)
	if err != nil {
		rsp.ErrRsp(c, -1, "failed to read update status")
		return
	}
	var latest manualUpdateJob
	for _, entry := range entries {
		job, err := loadManualJob(entry.Name())
		if err != nil || job.State == "staged" {
			continue
		}
		if latest.ID == "" || job.UpdatedAt.After(latest.UpdatedAt) {
			latest = job
		}
	}
	if latest.ID == "" {
		rsp.OkRsp(c)
		return
	}
	rsp.OkRspWithData(c, manualJobResponse(latest))
}
