//go:build !linux

package stream

import (
	log "github.com/sirupsen/logrus"
)

func SetCPUAffinity(cpuID int) error {
	log.Debugf("CPU affinity setting is not supported on this platform")
	return nil
}
