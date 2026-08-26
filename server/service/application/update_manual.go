package application

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"NanoKVM-Server/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v4/disk"
	log "github.com/sirupsen/logrus"
)

const (
	manualUpdateJobLifetime = 24 * time.Hour
	manualUploadOverhead    = int64(1 << 20)
	minimumFreeUpdateSpace  = uint64(16 << 20)
)

var (
	updateJobRoot = UpdateJobDir // replaced by tests only
	manualJobMu   sync.Mutex
)

type manualUpdateJob struct {
	ID        string       `json:"id"`
	Kind      artifactKind `json:"kind"`
	Version   string       `json:"version"`
	SHA256    string       `json:"sha256"`
	Size      int64        `json:"size"`
	State     string       `json:"state"`
	Error     string       `json:"error,omitempty"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	ExpiresAt time.Time    `json:"expiresAt"`
}

// InspectManualUpdate stages one upload and returns its validated job details.
func (s *Service) InspectManualUpdate(c *gin.Context) {
	var rsp proto.Response
	manualJobMu.Lock()
	defer manualJobMu.Unlock()
	if err := cleanupExpiredManualJobs(); err != nil {
		log.Warnf("clean expired update jobs: %v", err)
	}
	if err := removeStagedManualJobs(); err != nil {
		log.Errorf("clean staged update jobs: %v", err)
		rsp.ErrRsp(c, -1, "failed to prepare update staging")
		return
	}

	job, err := stageManualUpdate(c)
	if err != nil {
		log.Errorf("inspect uploaded update package: %v", err)
		rsp.ErrRsp(c, -1, err.Error())
		return
	}
	rsp.OkRspWithData(c, manualJobResponse(job))
}

// ConfirmManualUpdate changes a staged job to installing before starting its
// worker, so only an explicitly reviewed package can reach the installer.
func (s *Service) ConfirmManualUpdate(c *gin.Context) {
	var req proto.ConfirmManualUpdateReq
	var rsp proto.Response
	if err := proto.ParseFormRequest(c, &req); err != nil {
		rsp.ErrRsp(c, -1, "invalid arguments")
		return
	}
	if _, err := uuid.Parse(req.ID); err != nil {
		rsp.ErrRsp(c, -1, "invalid update job")
		return
	}

	manualJobMu.Lock()
	defer manualJobMu.Unlock()
	job, err := loadManualJob(req.ID)
	if err != nil {
		rsp.ErrRsp(c, -2, "update job not found")
		return
	}
	if time.Now().After(job.ExpiresAt) {
		_ = removeManualJob(job.ID)
		rsp.ErrRsp(c, -3, "update job expired; upload the package again")
		return
	}
	if job.State != "staged" {
		rsp.ErrRsp(c, -4, "update job is not ready to install")
		return
	}

	releaseLock, err := acquireUpdate()
	if err != nil {
		rsp.ErrRsp(c, -5, err.Error())
		return
	}
	job.State = "installing"
	job.Error = ""
	job.UpdatedAt = time.Now().UTC()
	if err := saveManualJob(job); err != nil {
		releaseLock()
		log.Errorf("persist installing update job: %v", err)
		rsp.ErrRsp(c, -6, "failed to start update")
		return
	}

	go func() {
		defer releaseLock()
		runManualUpdateJob(job.ID)
	}()
	rsp.OkRspWithData(c, manualJobResponse(job))
}

// DeleteManualUpdate discards an inspected package that has not started
// installation. It is intentionally unavailable once confirm has changed the
// job state, so a client cannot remove the bytes an installer is using.
func (s *Service) DeleteManualUpdate(c *gin.Context) {
	var rsp proto.Response
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		rsp.ErrRsp(c, -1, "invalid update job")
		return
	}

	manualJobMu.Lock()
	defer manualJobMu.Unlock()
	job, err := loadManualJob(id)
	if err != nil {
		rsp.ErrRsp(c, -2, "update job not found")
		return
	}
	if job.State != "staged" {
		rsp.ErrRsp(c, -3, "only an unconfirmed update job can be removed")
		return
	}
	if err := removeManualJob(id); err != nil {
		log.Errorf("remove staged update job %s: %v", id, err)
		rsp.ErrRsp(c, -4, "failed to remove update job")
		return
	}
	rsp.OkRsp(c)
}

// GetManualUpdate returns the persisted state of a confirmed manual job.
func (s *Service) GetManualUpdate(c *gin.Context) {
	var rsp proto.Response
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		rsp.ErrRsp(c, -1, "invalid update job")
		return
	}
	job, err := loadManualJob(id)
	if err != nil {
		rsp.ErrRsp(c, -2, "update job not found")
		return
	}
	rsp.OkRspWithData(c, manualJobResponse(job))
}

func stageManualUpdate(c *gin.Context) (manualUpdateJob, error) {
	// Staging performs all structural, checksum, and (for applications) Debian
	// metadata checks before exposing a job that can be confirmed.
	if err := ensureUpdateJobRoot(); err != nil {
		return manualUpdateJob{}, err
	}
	maxRequestBytes := maxManualPackageBytes + manualUploadOverhead
	if c.Request.ContentLength > maxRequestBytes {
		return manualUpdateJob{}, fmt.Errorf("uploaded request exceeds %d bytes", maxRequestBytes)
	}
	if c.Request.ContentLength > 0 {
		if err := ensureUpdateFreeSpace(uint64(c.Request.ContentLength)); err != nil {
			return manualUpdateJob{}, err
		}
	} else if err := ensureUpdateFreeSpace(minimumFreeUpdateSpace); err != nil {
		return manualUpdateJob{}, err
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBytes)

	reader, err := c.Request.MultipartReader()
	if err != nil {
		return manualUpdateJob{}, fmt.Errorf("invalid multipart data: %w", err)
	}
	id := uuid.NewString()
	jobDir := manualJobDirectory(id)
	if err := os.Mkdir(jobDir, 0o700); err != nil {
		return manualUpdateJob{}, fmt.Errorf("create update staging directory: %w", err)
	}
	success := false
	defer func() {
		if !success {
			_ = os.RemoveAll(jobDir)
		}
	}()

	originalName, artifactPath, checksum, size, err := saveManualUpload(reader, jobDir)
	if err != nil {
		return manualUpdateJob{}, err
	}
	if size <= 0 {
		return manualUpdateJob{}, errors.New("uploaded update package is empty")
	}

	inspection, err := inspectUpdateArtifact(artifactPath, originalName)
	if err != nil {
		return manualUpdateJob{}, err
	}
	if err := ensureUpdateFreeSpace(inspection.ExpandedBytes + minimumFreeUpdateSpace); err != nil {
		return manualUpdateJob{}, err
	}
	if inspection.Kind == artifactApplication {
		extractedDir, err := extractApplicationArchive(artifactPath, filepath.Join(jobDir, "extracted"), *inspection.Application)
		if err != nil {
			return manualUpdateJob{}, err
		}
		if err := validateApplicationDebMetadata(extractedDir, inspection.Version); err != nil {
			return manualUpdateJob{}, err
		}
	}

	now := time.Now().UTC()
	job := manualUpdateJob{
		ID:        id,
		Kind:      inspection.Kind,
		Version:   inspection.Version,
		SHA256:    checksum,
		Size:      size,
		State:     "staged",
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(manualUpdateJobLifetime),
	}
	if err := saveManualJob(job); err != nil {
		return manualUpdateJob{}, err
	}
	success = true
	return job, nil
}

func saveManualUpload(reader *multipart.Reader, jobDir string) (string, string, string, int64, error) {
	// Accept exactly one multipart field named file and retain its original name
	// only for format diagnostics; installation uses the staged path.
	var originalName, artifactPath, checksum string
	var size int64
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", "", "", 0, fmt.Errorf("read multipart data: %w", err)
		}
		if part.FormName() != "file" {
			continue
		}
		if artifactPath != "" {
			return "", "", "", 0, errors.New("multiple update files uploaded")
		}
		originalName = part.FileName()
		if err := validateUploadName(originalName); err != nil {
			return "", "", "", 0, err
		}
		artifactPath = filepath.Join(jobDir, "artifact")
		checksum, size, err = writeStagedUpload(part, artifactPath)
		if err != nil {
			return "", "", "", 0, err
		}
	}
	if artifactPath == "" {
		return "", "", "", 0, errors.New("no update file uploaded")
	}
	return originalName, artifactPath, checksum, size, nil
}

func validateUploadName(name string) error {
	if name == "" || len(name) > 255 || filepath.Base(name) != name || strings.ContainsAny(name, "/\\\x00") {
		return errors.New("invalid update file name")
	}
	return nil
}

func writeStagedUpload(part *multipart.Part, target string) (string, int64, error) {
	// Hash while writing and cap the stream so a missing Content-Length cannot
	// bypass the upload limit.
	output, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", 0, fmt.Errorf("create update staging file: %w", err)
	}
	success := false
	defer func() {
		_ = output.Close()
		if !success {
			_ = os.Remove(target)
		}
	}()

	hash := sha256.New()
	written, err := io.Copy(io.MultiWriter(output, hash), io.LimitReader(part, maxManualPackageBytes+1))
	if err != nil {
		return "", 0, fmt.Errorf("write update staging file: %w", err)
	}
	if written > maxManualPackageBytes {
		return "", 0, fmt.Errorf("uploaded package exceeds %d bytes", maxManualPackageBytes)
	}
	if err := output.Sync(); err != nil {
		return "", 0, fmt.Errorf("sync update staging file: %w", err)
	}
	if err := output.Close(); err != nil {
		return "", 0, fmt.Errorf("close update staging file: %w", err)
	}
	success = true
	return hex.EncodeToString(hash.Sum(nil)), written, nil
}

func runManualUpdateJob(id string) {
	// The job is persisted before this goroutine starts; every terminal result is
	// written back so polling clients can distinguish failure from interruption.
	job, err := loadManualJob(id)
	if err != nil {
		log.Errorf("load update job %s: %v", id, err)
		return
	}
	if err = installManualJob(job); err != nil {
		job.State = "failed"
		job.Error = safeUpdateError(err)
		log.Errorf("manual update job %s failed: %v", id, err)
		job.UpdatedAt = time.Now().UTC()
		if saveErr := saveManualJob(job); saveErr != nil {
			log.Errorf("persist update job %s result: %v", id, saveErr)
		}
		return
	}

	if job.Kind == artifactFirmware {
		if err := completeFirmwareRebootScheduling(&job); err != nil {
			log.Errorf("manual firmware update job %s: %v", id, err)
			return
		}
		log.Infof("manual firmware update job %s reboot scheduled", id)
		return
	}

	job.State = "succeeded"
	job.Error = ""
	job.UpdatedAt = time.Now().UTC()
	if saveErr := saveManualJob(job); saveErr != nil {
		log.Errorf("persist update job %s result: %v", id, saveErr)
		return
	}
	log.Infof("manual update job %s succeeded", id)
}

// Persist the terminal state before asking systemd to reboot. The reboot may
// stop this process immediately, so persisting afterward can leave the UI
// polling a state that can never be completed.
func completeFirmwareRebootScheduling(job *manualUpdateJob) error {
	job.State = "reboot_scheduled"
	job.Error = ""
	job.UpdatedAt = time.Now().UTC()
	if err := saveManualJob(*job); err != nil {
		return fmt.Errorf("persist reboot state: %w", err)
	}

	if err := firmwareRebootScheduler(job.ID); err != nil {
		job.State = "failed"
		job.Error = "firmware installed; automatic reboot scheduling failed; reboot manually"
		job.UpdatedAt = time.Now().UTC()
		if saveErr := saveManualJob(*job); saveErr != nil {
			return fmt.Errorf("schedule reboot: %w; persist recovery state: %v", err, saveErr)
		}
		return fmt.Errorf("schedule reboot: %w", err)
	}
	return nil
}

func installManualJob(job manualUpdateJob) error {
	// Re-hash and re-inspect the staged bytes to detect replacement after upload.
	artifactPath := manualArtifactPath(job.ID)
	checksum, err := sha256File(artifactPath)
	if err != nil {
		return err
	}
	if checksum != job.SHA256 {
		return errors.New("staged update package changed after inspection")
	}
	inspection, err := inspectUpdateArtifact(artifactPath, "")
	if err != nil {
		return err
	}
	if inspection.Kind != job.Kind || inspection.Version != job.Version {
		return errors.New("staged update package changed after inspection")
	}
	if err := ensureUpdateFreeSpace(inspection.ExpandedBytes + minimumFreeUpdateSpace); err != nil {
		return err
	}

	switch inspection.Kind {
	case artifactApplication:
		extractDir := filepath.Join(manualJobDirectory(job.ID), "install")
		if err := os.RemoveAll(extractDir); err != nil {
			return fmt.Errorf("clear application install directory: %w", err)
		}
		root, err := extractApplicationArchive(artifactPath, extractDir, *inspection.Application)
		if err != nil {
			return err
		}
		_ = sendMessage("install", 0)
		if err := install(root, inspection.Version); err != nil {
			return err
		}
		_ = sendMessage("restart", 0)
		return nil
	case artifactFirmware:
		_ = sendMessage("install", 0)
		return installFirmwarePackage(artifactPath)
	default:
		return errors.New("unsupported update package")
	}
}

func manualJobResponse(job manualUpdateJob) *proto.ManualUpdateRsp {
	return &proto.ManualUpdateRsp{
		ID:      job.ID,
		Version: job.Version,
		Size:    job.Size,
		State:   job.State,
		Error:   job.Error,
	}
}

func ensureUpdateJobRoot() error {
	if err := os.MkdirAll(updateJobRoot, 0o700); err != nil {
		return fmt.Errorf("create update job directory: %w", err)
	}
	return nil
}

func ensureUpdateFreeSpace(required uint64) error {
	if err := ensureUpdateJobRoot(); err != nil {
		return err
	}
	return ensureFreeSpaceAt(updateJobRoot, required)
}

func ensureFreeSpaceAt(directory string, required uint64) error {
	usage, err := disk.Usage(directory)
	if err != nil {
		return fmt.Errorf("check update free space: %w", err)
	}
	if usage.Free < required {
		return fmt.Errorf("insufficient free space for update: need %d bytes", required)
	}
	return nil
}

func manualJobDirectory(id string) string {
	return filepath.Join(updateJobRoot, id)
}

func manualArtifactPath(id string) string {
	return filepath.Join(manualJobDirectory(id), "artifact")
}

func manualJobPath(id string) string {
	return filepath.Join(manualJobDirectory(id), "job.json")
}

func saveManualJob(job manualUpdateJob) error {
	if _, err := uuid.Parse(job.ID); err != nil {
		return errors.New("invalid update job")
	}
	dir := manualJobDirectory(job.ID)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create update job directory: %w", err)
	}
	data, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		return fmt.Errorf("encode update job: %w", err)
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(dir, ".job.*")
	if err != nil {
		return fmt.Errorf("create temporary update job: %w", err)
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, manualJobPath(job.ID)); err != nil {
		return fmt.Errorf("replace update job: %w", err)
	}
	return nil
}

func loadManualJob(id string) (manualUpdateJob, error) {
	if _, err := uuid.Parse(id); err != nil {
		return manualUpdateJob{}, errors.New("invalid update job")
	}
	data, err := os.ReadFile(manualJobPath(id))
	if err != nil {
		return manualUpdateJob{}, err
	}
	var job manualUpdateJob
	if err := json.Unmarshal(data, &job); err != nil {
		return manualUpdateJob{}, err
	}
	if job.ID != id || job.State == "" || job.SHA256 == "" {
		return manualUpdateJob{}, errors.New("invalid update job")
	}
	return job, nil
}

func removeManualJob(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid update job")
	}
	return os.RemoveAll(manualJobDirectory(id))
}

func cleanupExpiredManualJobs() error {
	// Expired staged/terminal jobs are disposable; installing jobs stay until a
	// restart recovery pass can report their ambiguous installation state.
	if err := ensureUpdateJobRoot(); err != nil {
		return err
	}
	entries, err := os.ReadDir(updateJobRoot)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := uuid.Parse(entry.Name()); err != nil {
			continue
		}
		job, err := loadManualJob(entry.Name())
		if err != nil || (job.State != "installing" && time.Now().After(job.ExpiresAt)) {
			if err := os.RemoveAll(filepath.Join(updateJobRoot, entry.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

func recoverInterruptedManualUpdates() error {
	// Only "installing" is ambiguous after restart. A reboot_scheduled firmware
	// job already records its terminal pre-reboot state and must be preserved.
	if err := ensureUpdateJobRoot(); err != nil {
		return err
	}
	entries, err := os.ReadDir(updateJobRoot)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		job, err := loadManualJob(entry.Name())
		if err != nil || job.State != "installing" {
			continue
		}
		job.State = "failed"
		job.Error = "update process was interrupted; installation state is unknown; check the current version before retrying"
		job.UpdatedAt = time.Now().UTC()
		if err := saveManualJob(job); err != nil {
			return fmt.Errorf("recover update job %s: %w", job.ID, err)
		}
	}
	return nil
}

// removeStagedManualJobs keeps at most one unconfirmed upload. It runs while
// InspectManualUpdate holds manualJobMu, and ConfirmManualUpdate takes the
// same mutex before turning a staged job into an installing job.
func removeStagedManualJobs() error {
	if err := ensureUpdateJobRoot(); err != nil {
		return err
	}
	entries, err := os.ReadDir(updateJobRoot)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := uuid.Parse(entry.Name()); err != nil {
			continue
		}
		job, err := loadManualJob(entry.Name())
		if err != nil {
			return fmt.Errorf("read staged update job %s: %w", entry.Name(), err)
		}
		if job.State == "staged" {
			if err := removeManualJob(job.ID); err != nil {
				return fmt.Errorf("remove staged update job %s: %w", job.ID, err)
			}
		}
	}
	return nil
}

func sha256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func safeUpdateError(err error) string {
	if err == nil {
		return ""
	}
	message := strings.TrimSpace(err.Error())
	if len(message) > 512 {
		return message[:512]
	}
	return message
}
