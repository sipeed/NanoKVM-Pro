package hid

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
)

// ReportWriter is valid only inside WithInputReports. Sharing the existing HID
// locks keeps browser reports and MCP shortcuts from interleaving mid-action.
type ReportWriter struct{ hid *Hid }

func (h *Hid) WithInputReports(action func(*ReportWriter) error) error {
	h.Lock()
	defer h.Unlock()
	if h.g0 == nil || h.g1 == nil || h.g2 == nil {
		h.OpenNoLock()
	}
	err := action(&ReportWriter{hid: h})
	if errors.Is(err, os.ErrClosed) || errors.Is(err, syscall.ENODEV) || errors.Is(err, syscall.ESHUTDOWN) {
		h.CloseNoLock()
		h.g0, h.g1, h.g2 = nil, nil, nil
	}
	return err
}

func writeReport(file *os.File, report []byte, size int) error {
	if len(report) != size {
		return fmt.Errorf("invalid HID report length: %d, want %d", len(report), size)
	}
	if file == nil {
		return fmt.Errorf("USB HID device unavailable")
	}
	_ = file.SetWriteDeadline(time.Now().Add(8 * time.Millisecond))
	n, err := file.Write(report)
	if err == nil && n != len(report) {
		err = io.ErrShortWrite
	}
	return err
}

func (w *ReportWriter) WriteKeyboardReport(report []byte) error {
	return writeReport(w.hid.g0, report, 8)
}

func (w *ReportWriter) WriteRelativeMouseReport(report []byte) error {
	return writeReport(w.hid.g1, report, 4)
}

func (w *ReportWriter) WriteAbsoluteMouseReport(report []byte) error {
	return writeReport(w.hid.g2, report, 6)
}
