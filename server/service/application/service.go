package application

import log "github.com/sirupsen/logrus"

const (
	StableURL  = "https://cdn.sipeed.com/nanokvm"
	PreviewURL = "https://cdn.sipeed.com/nanokvm/preview"

	AppDir = "/kvmapp"

	// TempDir is kept for compatibility with the existing application updater.
	// New updates always use a private per-operation workspace below it rather
	// than removing this directory as a whole.
	TempDir = "/root/.kvmcache"

	UpdateSourceConfigFile = "/etc/kvm/update-source.json"
	UpdateJobDir           = "/var/lib/kvm/update-jobs"
	FirmwareUpdateScript   = "/kvmcomm/scripts/firmware_update.sh"
)

var appNames = []string{"nanokvmpro", "pikvm", "kvmcomm"}

type Service struct{}

func NewService() *Service {
	// Recover jobs left in "installing" when a service restart interrupted the
	// worker goroutine; this keeps the persisted status actionable for the UI.
	if err := recoverInterruptedManualUpdates(); err != nil {
		log.Warnf("recover interrupted manual updates: %v", err)
	}
	return &Service{}
}
