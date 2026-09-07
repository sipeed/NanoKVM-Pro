package application

import (
	"errors"
	"testing"
)

func TestCompleteFirmwareRebootSchedulingPersistsScheduledState(t *testing.T) {
	previousRoot := updateJobRoot
	previousScheduler := firmwareRebootScheduler
	updateJobRoot = t.TempDir()
	defer func() {
		updateJobRoot = previousRoot
		firmwareRebootScheduler = previousScheduler
	}()

	job := manualUpdateJob{
		ID:     "03f2760e-0d79-497f-b9a1-13b64a8c4a06",
		Kind:   artifactFirmware,
		SHA256: "test-checksum",
		State:  "installing",
	}
	if err := saveManualJob(job); err != nil {
		t.Fatal(err)
	}
	called := false
	firmwareRebootScheduler = func(id string) error {
		called = id == job.ID
		return nil
	}

	if err := completeFirmwareRebootScheduling(&job); err != nil {
		t.Fatalf("completeFirmwareRebootScheduling() error = %v", err)
	}
	if !called {
		t.Fatal("firmware reboot scheduler was not called")
	}
	persisted, err := loadManualJob(job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if persisted.State != "reboot_scheduled" {
		t.Fatalf("state = %q, want reboot_scheduled", persisted.State)
	}
	if persisted.Error != "" {
		t.Fatalf("error = %q, want empty", persisted.Error)
	}
}

func TestCompleteFirmwareRebootSchedulingPersistsManualRecoveryState(t *testing.T) {
	previousRoot := updateJobRoot
	previousScheduler := firmwareRebootScheduler
	updateJobRoot = t.TempDir()
	defer func() {
		updateJobRoot = previousRoot
		firmwareRebootScheduler = previousScheduler
	}()

	job := manualUpdateJob{
		ID:     "bc5a7ae3-0e0b-4f23-8e51-17609cf0f5dd",
		Kind:   artifactFirmware,
		SHA256: "test-checksum",
		State:  "installing",
	}
	if err := saveManualJob(job); err != nil {
		t.Fatal(err)
	}
	firmwareRebootScheduler = func(string) error {
		return errors.New("reboot unavailable")
	}

	if err := completeFirmwareRebootScheduling(&job); err == nil {
		t.Fatal("completeFirmwareRebootScheduling() unexpectedly succeeded")
	}
	persisted, err := loadManualJob(job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if persisted.State != "failed" {
		t.Fatalf("state = %q, want failed", persisted.State)
	}
	const want = "firmware installed; automatic reboot scheduling failed; reboot manually"
	if persisted.Error != want {
		t.Fatalf("error = %q, want %q", persisted.Error, want)
	}
}

func TestRecoverInterruptedManualUpdate(t *testing.T) {
	previousRoot := updateJobRoot
	updateJobRoot = t.TempDir()
	defer func() { updateJobRoot = previousRoot }()

	job := manualUpdateJob{
		ID:      "d636979a-427e-49b7-8f7f-d880a902de91",
		Kind:    artifactApplication,
		Version: "1.2.15",
		SHA256:  "test-checksum",
		State:   "installing",
	}
	if err := saveManualJob(job); err != nil {
		t.Fatal(err)
	}
	if service := NewService(); service == nil {
		t.Fatal("NewService() returned nil")
	}

	recovered, err := loadManualJob(job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if recovered.State != "failed" || recovered.Error == "" {
		t.Fatalf("unexpected recovered job: %#v", recovered)
	}
}
