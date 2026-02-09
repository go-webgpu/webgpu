package wgpu

import (
	"testing"
)

// TestSurfaceGetCapabilities_NilSurface tests nil safety for surface.
func TestSurfaceGetCapabilities_NilSurface(t *testing.T) {
	instance, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("Failed to request adapter: %v", err)
	}
	defer adapter.Release()

	var surface *Surface
	_, err = surface.GetCapabilities(adapter)
	if err == nil {
		t.Error("Expected error for nil surface, got nil")
	}
}

// TestSurfaceGetCapabilities_NilAdapter tests nil safety for adapter.
func TestSurfaceGetCapabilities_NilAdapter(t *testing.T) {
	// We cannot create a real surface without a window, so we just test nil adapter
	surface := &Surface{handle: 1} // fake handle
	_, err := surface.GetCapabilities(nil)
	if err == nil {
		t.Error("Expected error for nil adapter, got nil")
	}
}

// Note: Full integration testing of GetCapabilities requires a real window surface,
// which is tested in the examples (e.g., examples/triangle).
