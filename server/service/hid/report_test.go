package hid

import (
	"bytes"
	"os"
	"testing"
)

func TestInputReportsUseProFormatsAndReturnErrors(t *testing.T) {
	h := &Hid{}
	files := []**os.File{&h.g0, &h.g1, &h.g2}
	for _, dest := range files {
		file, err := os.CreateTemp(t.TempDir(), "hid")
		if err != nil {
			t.Fatal(err)
		}
		*dest = file
		t.Cleanup(func() { _ = file.Close() })
	}
	err := h.WithInputReports(func(writer *ReportWriter) error {
		if err := writer.WriteKeyboardReport(make([]byte, 7)); err == nil {
			t.Fatal("short keyboard report accepted")
		}
		if err := writer.WriteAbsoluteMouseReport(make([]byte, 7)); err == nil {
			t.Fatal("Go mouse format accepted by Pro")
		}
		if err := writer.WriteKeyboardReport([]byte{1, 0, 4, 0, 0, 0, 0, 0}); err != nil {
			return err
		}
		return writer.WriteAbsoluteMouseReport([]byte{0, 0, 0x40, 0xff, 0x7f, 0})
	})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(h.g2.Name())
	if err != nil || !bytes.Equal(data, []byte{0, 0, 0x40, 0xff, 0x7f, 0}) {
		t.Fatalf("absolute output: %v %v", data, err)
	}
	writer := &ReportWriter{hid: &Hid{}}
	if err := writer.WriteKeyboardReport(make([]byte, 8)); err == nil {
		t.Fatal("unavailable device reported success")
	}
	_ = h.g0.Close()
	if err := (&ReportWriter{hid: h}).WriteKeyboardReport(make([]byte, 8)); err == nil {
		t.Fatal("closed device reported success")
	}
}

func TestInputReportsInvalidateDisconnectedDevice(t *testing.T) {
	h := &Hid{}
	for _, dest := range []**os.File{&h.g0, &h.g1, &h.g2} {
		file, err := os.CreateTemp(t.TempDir(), "hid")
		if err != nil {
			t.Fatal(err)
		}
		*dest = file
		t.Cleanup(func() { _ = file.Close() })
	}
	_ = h.g0.Close()
	err := h.WithInputReports(func(writer *ReportWriter) error { return writer.WriteKeyboardReport(make([]byte, 8)) })
	if err == nil {
		t.Fatal("disconnected device reported success")
	}
	if h.g0 != nil || h.g1 != nil || h.g2 != nil {
		t.Fatal("disconnected descriptors were kept for the next action")
	}
}
