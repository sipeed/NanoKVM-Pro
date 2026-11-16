package ai

import (
	"fmt"
	"math"
)

const maxRelativeMouseReports = 1024

func absoluteCoordinate(value float64) uint16 {
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return uint16(math.Round(0x7fff * value))
}

func buildAbsolutePointerReport(x uint16, y uint16, touching bool, secondary bool) []byte {
	buttons := byte(0)
	if touching {
		buttons |= 1 << 0
	}
	if secondary {
		buttons |= 1 << 1
	}

	return []byte{
		buttons,
		byte(x & 0xff),
		byte((x >> 8) & 0xff),
		byte(y & 0xff),
		byte((y >> 8) & 0xff),
		0, // wheel
	}
}

func buildRelativeMouseReport(deltaX int, deltaY int, buttons byte, wheel int) []byte {
	return []byte{
		buttons,
		byte(int8(clampInt(deltaX, -127, 127))),
		byte(int8(clampInt(deltaY, -127, 127))),
		byte(int8(clampInt(wheel, -127, 127))),
	}
}

func buildRelativeMoveReports(deltaX int, deltaY int) ([][]byte, error) {
	if relativeReportCount(deltaX, deltaY) > maxRelativeMouseReports {
		return nil, fmt.Errorf("relative movement exceeds %d reports", maxRelativeMouseReports)
	}

	reports := make([][]byte, 0)
	for deltaX != 0 || deltaY != 0 {
		x := clampInt(deltaX, -127, 127)
		y := clampInt(deltaY, -127, 127)
		reports = append(reports, buildRelativeMouseReport(x, y, 0, 0))
		deltaX -= x
		deltaY -= y
	}
	return reports, nil
}

func mouseButtonBit(button string) (byte, error) {
	switch button {
	case "left":
		return 1 << 0, nil
	case "right":
		return 1 << 1, nil
	case "middle":
		return 1 << 2, nil
	case "back":
		return 1 << 3, nil
	case "forward":
		return 1 << 4, nil
	default:
		return 0, fmt.Errorf("unknown mouse button: %s", button)
	}
}

func clampInt(value int, min int, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func relativeReportCount(deltaX int, deltaY int) int {
	xReports := reportsForDelta(deltaX)
	yReports := reportsForDelta(deltaY)
	if xReports > yReports {
		return xReports
	}
	return yReports
}

func reportsForDelta(delta int) int {
	value := int64(delta)
	if value == 0 {
		return 0
	}
	if value < 0 {
		return int((-(value + 1))/127 + 1)
	}
	return int((value-1)/127 + 1)
}
