//go:build linux || darwin

package wgpu

import (
	"testing"

	"github.com/go-webgpu/goffi/ffi"
)

func closeTimestampPeriodABILibrary(t *testing.T, library Library) {
	t.Helper()
	unixLibrary, ok := library.(*unixLibrary)
	if !ok {
		t.Fatalf("timestamp-period ABI library has type %T, want *unixLibrary", library)
	}
	if err := ffi.FreeLibrary(unixLibrary.handle); err != nil {
		t.Fatalf("close timestamp-period ABI library: %v", err)
	}
}
