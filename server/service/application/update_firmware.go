package application

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var firmwareRebootScheduler = scheduleFirmwareReboot

// installFirmwarePackage intentionally delegates writing to the existing
// kvmcomm updater. The service has already validated the same staged archive;
// the script remains the sole component allowed to select boot media and write
// its partitions.
func installFirmwarePackage(packagePath string) error {
	info, err := os.Stat(packagePath)
	if err != nil {
		return fmt.Errorf("stat staged firmware package: %w", err)
	}
	if !info.Mode().IsRegular() {
		return errors.New("staged firmware package is not a regular file")
	}
	if _, err := os.Stat(FirmwareUpdateScript); err != nil {
		return fmt.Errorf("firmware updater is unavailable: %w", err)
	}

	output := &limitedCommandOutput{limit: 8 << 10}
	cmd := exec.Command(FirmwareUpdateScript, "update", packagePath)
	cmd.Stdout = output
	cmd.Stderr = output
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("firmware update failed: %s", output.String())
	}
	return nil
}

func scheduleFirmwareReboot(string) error {
	if err := exec.Command("systemctl", "reboot").Start(); err != nil {
		return fmt.Errorf("automatic reboot failed: %w", err)
	}
	return nil
}

// limitedCommandOutput bounds logs returned to the API while still consuming
// the child process output so the updater cannot block on a full pipe.
type limitedCommandOutput struct {
	buffer bytes.Buffer
	limit  int
}

func (w *limitedCommandOutput) Write(data []byte) (int, error) {
	if w.buffer.Len() < w.limit {
		remaining := w.limit - w.buffer.Len()
		if len(data) > remaining {
			_, _ = w.buffer.Write(data[:remaining])
		} else {
			_, _ = w.buffer.Write(data)
		}
	}
	return len(data), nil
}

func (w *limitedCommandOutput) String() string {
	if w.buffer.Len() == 0 {
		return "command returned no output"
	}
	return w.buffer.String()
}
