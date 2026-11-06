//go:build windows

package wgpu

import (
	"testing"
	"unsafe"
)

func TestSurfaceSourceSwapChainPanelLayout(t *testing.T) {
	var s surfaceSourceSwapChainPanel
	if got := unsafe.Sizeof(s); got != 24 {
		t.Errorf("sizeof(surfaceSourceSwapChainPanel) = %d, want 24", got)
	}
	if got := unsafe.Offsetof(s.panelNative); got != 16 {
		t.Errorf("offsetof(surfaceSourceSwapChainPanel.panelNative) = %d, want 16", got)
	}
}

func TestCreateSurfaceFromSwapChainPanelRejectsNil(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Skipf("CreateInstance failed: %v", err)
	}
	defer inst.Release()
	if _, err := inst.CreateSurfaceFromSwapChainPanel(0); err == nil {
		t.Fatal("expected an error for a nil swap chain panel")
	}
}
