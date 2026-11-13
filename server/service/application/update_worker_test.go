package application

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func prepareUpdateWorkerTest(t *testing.T) {
	t.Helper()
	root, lock, inspector, starter := updateJobRoot, updateLockPath, updateWorkerInspector, updateWorkerStarter
	updateJobRoot = t.TempDir()
	updateLockPath = filepath.Join(t.TempDir(), "update.lock")
	updateWorkerInspector = func() (updateWorkerInfo, error) { return updateWorkerInfo{}, nil }
	t.Cleanup(func() {
		updateJobRoot, updateLockPath, updateWorkerInspector, updateWorkerStarter = root, lock, inspector, starter
	})
}

func testInstallingJob(t *testing.T) manualUpdateJob {
	t.Helper()
	job := manualUpdateJob{
		ID: "a52b258f-c832-4cf3-a761-409dfe22a710", Kind: artifactApplication,
		Version: "1.2.16", SHA256: "expected", State: "installing",
		CreatedAt: time.Now(), UpdatedAt: time.Now(), ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := saveManualJob(job); err != nil {
		t.Fatal(err)
	}
	return job
}

func TestRecoveryPreservesIndependentUpdateDuringServiceRestart(t *testing.T) {
	prepareUpdateWorkerTest(t)
	job := testInstallingJob(t)
	updateWorkerInspector = func() (updateWorkerInfo, error) {
		return updateWorkerInfo{Active: true, JobID: job.ID}, nil
	}
	if err := recoverInterruptedManualUpdates(); err != nil {
		t.Fatal(err)
	}
	got, err := loadManualJob(job.ID)
	if err != nil || got.State != "installing" {
		t.Fatalf("recovered live job: %#v, %v", got, err)
	}
	if release, err := acquireUpdate(); err == nil {
		release()
		t.Fatal("another update was admitted while the systemd worker was active")
	}
}

func TestRecoveryDoesNotRaceAdmissionOrWorkerFileLock(t *testing.T) {
	prepareUpdateWorkerTest(t)
	job := testInstallingJob(t)
	release, err := acquireUpdateFileLock()
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if err := recoverInterruptedManualUpdates(); err != nil {
		t.Fatal(err)
	}
	got, err := loadManualJob(job.ID)
	if err != nil || got.State != "installing" {
		t.Fatalf("recovered locked job: %#v, %v", got, err)
	}
}

func TestOnlineUpdateArtifactSurvivesRequestWorkspaceCleanup(t *testing.T) {
	prepareUpdateWorkerTest(t)
	workspace := t.TempDir()
	artifact := filepath.Join(workspace, "bundle.tar.gz")
	if err := os.WriteFile(artifact, []byte("validated archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	job, err := stageOnlineUpdate(artifact, artifactInspection{Kind: artifactApplication, Version: "1.2.16"})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(workspace); err != nil {
		t.Fatal(err)
	}
	checksum, err := sha256File(manualArtifactPath(job.ID))
	if err != nil || checksum != job.SHA256 {
		t.Fatalf("worker lost staged artifact: %q, %v", checksum, err)
	}
}

func TestConfirmWorkerStartFailureIsPersistedAndReleasesLock(t *testing.T) {
	prepareUpdateWorkerTest(t)
	job := testInstallingJob(t)
	job.State = "staged"
	if err := saveManualJob(job); err != nil {
		t.Fatal(err)
	}
	updateWorkerStarter = func(string) error { return errors.New("systemd-run unavailable") }
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)
	c.Request = httptest.NewRequest("POST", "/api/application/update/manual/confirm", strings.NewReader(`{"id":"`+job.ID+`"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	(&Service{}).ConfirmManualUpdate(c)
	got, err := loadManualJob(job.ID)
	if err != nil || got.State != "failed" || got.Error == "" {
		t.Fatalf("launch failure: %#v, %v", got, err)
	}
	release, err := acquireUpdate()
	if err != nil {
		t.Fatalf("failed launch retained update lock: %v", err)
	}
	release()
}

func TestWorkerRejectsChangedArtifactAndPersistsFailure(t *testing.T) {
	prepareUpdateWorkerTest(t)
	job := testInstallingJob(t)
	if err := os.WriteFile(manualArtifactPath(job.ID), []byte("changed archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := runManualUpdateJob(job.ID); err == nil {
		t.Fatal("changed archive was accepted")
	}
	got, err := loadManualJob(job.ID)
	if err != nil || got.State != "failed" {
		t.Fatalf("worker failure: %#v, %v", got, err)
	}
}

func TestWorkerInfoDistinguishesMissingUnitFromSystemdFailure(t *testing.T) {
	info, err := parseUpdateWorkerInfo("LoadState=not-found\nActiveState=inactive\n", errors.New("exit 1"))
	if err != nil || info.Active {
		t.Fatalf("missing unit: %#v, %v", info, err)
	}
	if _, err := parseUpdateWorkerInfo("", errors.New("cannot connect to systemd")); err == nil {
		t.Fatal("systemd failure must not mark a running job interrupted")
	}
}

func TestInstallationOrdersMainApplicationLastAndStopsAfterFailure(t *testing.T) {
	if _, err := exec.LookPath("dpkg-deb"); err != nil {
		t.Skip("dpkg-deb is required")
	}
	root, bin, logPath := t.TempDir(), t.TempDir(), filepath.Join(t.TempDir(), "dpkg.log")
	const version = "1.2.16"
	for _, name := range []string{"nanokvmpro", "pikvm", "kvmcomm"} {
		packageName := name
		if name == "nanokvmpro" {
			packageName = "nanokvm"
		}
		packageRoot := filepath.Join(t.TempDir(), "payload")
		if err := os.MkdirAll(filepath.Join(packageRoot, "DEBIAN"), 0o755); err != nil {
			t.Fatal(err)
		}
		control := fmt.Sprintf("Package: %s\nVersion: %s\nArchitecture: arm64\nMaintainer: Test <test@example.com>\nDescription: test\n", packageName, version)
		if err := os.WriteFile(filepath.Join(packageRoot, "DEBIAN", "control"), []byte(control), 0o644); err != nil {
			t.Fatal(err)
		}
		filename := name + "_" + version + "_arm64.deb"
		if output, err := exec.Command("dpkg-deb", "--build", packageRoot, filepath.Join(root, filename)).CombinedOutput(); err != nil {
			t.Fatalf("build test DEB: %s: %v", output, err)
		}
		data, err := os.ReadFile(filepath.Join(root, filename))
		if err != nil {
			t.Fatal(err)
		}
		hash := sha512.Sum512(data)
		metadata, err := json.Marshal(FileInfo{Name: filename, Version: version, SHA512: base64.StdEncoding.EncodeToString(hash[:]), Size: len(data)})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, name+"_"+version+".json"), metadata, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	script := "#!/bin/sh\nname=$(basename \"$2\")\nprintf '%s\\n' \"$name\" >>\"$TEST_DPKG_LOG\"\ncase \"$name\" in \"${TEST_DPKG_FAIL:-never}_\"*) exit 1;; esac\n"
	if err := os.WriteFile(filepath.Join(bin, "dpkg"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+":"+os.Getenv("PATH"))
	t.Setenv("TEST_DPKG_LOG", logPath)
	report := func(string, int) error { return nil }
	if err := install(root, version, report); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"pikvm_1.2.16_arm64.deb", "kvmcomm_1.2.16_arm64.deb", "nanokvmpro_1.2.16_arm64.deb"}
	if got := strings.Fields(string(data)); !reflect.DeepEqual(got, want) {
		t.Fatalf("install order: %v", got)
	}
	if err := os.WriteFile(logPath, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TEST_DPKG_FAIL", "kvmcomm")
	if err := install(root, version, report); err == nil {
		t.Fatal("failed package was ignored")
	}
	data, err = os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "nanokvmpro_") {
		t.Fatalf("main application restarted after a component failed: %s", data)
	}
}
