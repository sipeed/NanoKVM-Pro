package vm

import (
	"fmt"
	"os"
	"strings"
)

const screenVersionPath = "/proc/lt6911_info/version"

type screenType string

const (
	screenTypeATX  screenType = "ATX"
	screenTypeDesk screenType = "DESK"
)

var readScreenVersion = os.ReadFile

func getScreenType() (screenType, error) {
	content, err := readScreenVersion(screenVersionPath)
	if err != nil {
		return "", err
	}
	return screenTypeFromVersion(content)
}

func screenTypeFromVersion(content []byte) (screenType, error) {
	version := string(content)
	switch {
	case strings.Contains(version, "ATX"):
		return screenTypeATX, nil
	case strings.Contains(version, "Desk"):
		return screenTypeDesk, nil
	default:
		return "", fmt.Errorf("unknown screen type")
	}
}
