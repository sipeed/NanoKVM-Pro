package utils

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	log "github.com/sirupsen/logrus"
)

var (
	systemctlClient     *dbus.Conn
	systemctlClientOnce sync.Once
)

func init() {
	var err error

	systemctlClientOnce.Do(func() {
		ctx := context.Background()
		systemctlClient, err = dbus.NewSystemConnectionContext(ctx)
		if err != nil {
			log.Errorf("connect systemctl failed error=%s", err)
			systemctlClient = nil
		}
	})
}

func IsServiceRunning(servicename string) (bool, error) {
	// Check if the service exists
	if systemctlClient == nil {
		return false, fmt.Errorf("failed to connect to systemd bus")
	}

	properties, err := systemctlClient.GetUnitPropertiesContext(context.Background(), servicename)
	if err != nil {
		return false, err
	}

	loadState, exists := properties["LoadState"]
	if !exists {
		return false, fmt.Errorf("LoadState property not found")
	}

	if loadStateStr, ok := loadState.(string); ok && loadStateStr == "not-found" {
		return false, fmt.Errorf("service not found")
	}

	activeState, exists := properties["ActiveState"]
	if !exists {
		return false, fmt.Errorf("ActiveState property not found")
	}

	activeStateStr, ok := activeState.(string)
	if !ok {
		return false, fmt.Errorf("ActiveState is not a string")
	}

	subState, exists := properties["SubState"]
	if !exists {
		return false, fmt.Errorf("SubState property not found")
	}

	subStateStr, ok := subState.(string)
	if !ok {
		return false, fmt.Errorf("SubState is not a string")
	}

	subStateActive := subStateStr == "listening" || subStateStr == "running"
	// is running
	return activeStateStr == "active" && subStateActive, nil
}

func DaemonReload() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		log.Errorf("Failed to connect systemd D-Bus: %v", err)
		return err
	}
	defer conn.Close()

	if err := conn.ReloadContext(ctx); err != nil {
		log.Errorf("Failed to execute daemon-reload: %v", err)
		return err
	}

	log.Debugf("Systemd daemon-reload completed successfully.")
	return nil
}

func StartService(name string, enable bool) error {
	if enable {
		_, _, err := systemctlClient.EnableUnitFilesContext(context.Background(), []string{name}, false, true)
		if err != nil {
			return fmt.Errorf("failed to enable service: %v", err)
		}

		if err := RestartService(name); err != nil {
			log.Debugf("restart service failed %v", err)
			return err
		}

		return nil
	}

	ch := make(chan string)

	if _, err := systemctlClient.StartUnitContext(context.Background(), name, "replace", ch); err != nil {
		return fmt.Errorf("failed to start service: %v", err)
	}

	select {
	case result := <-ch:
		if result != "done" {
			return fmt.Errorf("service start failed: %s", result)
		}
	case <-time.After(6 * time.Second):
		return fmt.Errorf("service start timed out")
	}

	if err := RestartService(name); err != nil {
		log.Debugf("restart service failed %v", err)
		return err
	}

	return nil
}

func RestartService(serviceName string) error {
	ch := make(chan string)

	if _, err := systemctlClient.RestartUnitContext(context.Background(), serviceName, "replace", ch); err != nil {
		return fmt.Errorf("failed to restart service: %v", err)
	}

	select {
	case result := <-ch:
		if result != "done" {
			return fmt.Errorf("service restart failed: %s", result)
		}
	case <-time.After(6 * time.Second):
		return fmt.Errorf("service restart timed out")
	}

	return nil
}

func StopService(name string, disable bool) error {
	if _, err := execute("systemctl stop " + name); err != nil {
		return err
	}

	if disable {
		if _, err := execute("systemctl disable " + name); err != nil {
			return err
		}
	}

	return nil
}

func execute(command string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("failed to execute %s: %s", command, err)
		return output, err
	}

	return output, err
}
