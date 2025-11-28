package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ChronyStatusInfo struct {
	IsSynchronized bool
	LastSyncTime   int64
	LeapStatus     string
	Stratum        string
	ReferenceID    string
}

func ChronycBurst() error {
	cmd := exec.Command("chronyc", "burst", "4/4")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		return fmt.Errorf("[chronyc burst] Execution failed: %w, output: %s", err, outputStr)
	}

	if !strings.Contains(outputStr, "200 OK") {
		return fmt.Errorf("[chronyc burst] Not the expected result: %s", outputStr)
	}

	return nil
}

func GetChronySyncStatus() (*ChronyStatusInfo, error) {
	cmd := exec.Command("chronyc", "tracking")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("[chronyc tracking] Execution failed: %w, output: %s", err, stderr.String())
	}

	return parseChronyTrackingOutput(out.String())
}

func parseChronyTrackingOutput(output string) (*ChronyStatusInfo, error) {
	info := &ChronyStatusInfo{}
	lines := strings.Split(output, "\n")
	parsedFields := 0

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Reference ID":
			info.ReferenceID = value
			parsedFields++
		case "Stratum":
			info.Stratum = value
			parsedFields++
		case "Ref time (UTC)":
			parsedTime, err := time.Parse(time.ANSIC, value)
			if err != nil {
				return nil, fmt.Errorf("[chronyc tracking] Parse Ref time (UTC) failed: %w", err)
			}
			info.LastSyncTime = parsedTime.UnixMilli()
			parsedFields++
		case "Leap status":
			info.LeapStatus = value
			parsedFields++
		}
	}

	if parsedFields < 4 {
		return nil, errors.New("[chronyc tracking] Parse output failed")
	}

	if info.LeapStatus == "Normal" {
		info.IsSynchronized = true
	}

	return info, nil
}
