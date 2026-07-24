//go:build windows

package wgpu

import (
	"syscall"
	"testing"
)

func closeTimestampPeriodABILibrary(t *testing.T, library Library) {
	t.Helper()
	windowsLibrary, ok := library.(*windowsLibrary)
	if !ok {
		t.Fatalf("timestamp-period ABI library has type %T, want *windowsLibrary", library)
	}
	if err := syscall.FreeLibrary(syscall.Handle(windowsLibrary.dll.Handle())); err != nil {
		t.Fatalf("close timestamp-period ABI library: %v", err)
	}
}
