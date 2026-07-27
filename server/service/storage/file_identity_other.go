//go:build !linux

package storage

import "os"

func fileIdentity(_ os.FileInfo) (uint64, uint64) {
	return 0, 0
}
