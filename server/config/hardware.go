package config

import (
	"os"
	"regexp"
	"strings"
)

type HWVersion int

const (
	HWVersionUnknown = iota
	HWVersionDesk
	HWVersionATX

	HWVersionFile = "/proc/lt6911_info/version"
)

var HWPro = Hardware{
	Version:      HWVersionDesk,
	GPIOReset:    "/sys/class/gpio/gpio35/value",
	GPIOPower:    "/sys/class/gpio/gpio7/value",
	GPIOPowerLED: "/sys/class/gpio/gpio75/value",
	GPIOHDDLed:   "/sys/class/gpio/gpio74/value",
}

func (h HWVersion) String() string {
	switch h {
	case HWVersionDesk, HWVersionATX:
		return "Pro"
	default:
		return "Unknown"
	}
}

func DeviceNumber() string {
	content, err := os.ReadFile(HWVersionFile)
	if err != nil {
		return "unknown"
	}

	parts := strings.Split(string(content), " ")
	length := len(parts)

	if length > 0 {
		return parts[length-1]
	}

	return "unknown"
}

func GetHwVersion() HWVersion {
	content, err := os.ReadFile(HWVersionFile)
	if err != nil {
		return HWVersionUnknown
	}

	version := strings.ToLower(string(content))

	var atx = regexp.MustCompile(`(?i)atx`)
	var desk = regexp.MustCompile(`(?i)(desk)`)

	switch {
	case desk.MatchString(version):
		return HWVersionDesk
	case atx.MatchString(version):
		return HWVersionATX
	default:
		return HWVersionUnknown
	}
}

func getHardware() (h Hardware) {
	return HWPro
}
