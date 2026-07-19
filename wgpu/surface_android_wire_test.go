package wgpu

import (
	"testing"
	"unsafe"
)

func TestABISurfaceSourceAndroidNativeWindow(t *testing.T) {
	source, err := newSurfaceSourceAndroidNativeWindow(0x1234)
	if err != nil {
		t.Fatalf("newSurfaceSourceAndroidNativeWindow: %v", err)
	}

	if got := unsafe.Sizeof(source); got != 24 {
		t.Fatalf("sizeof(surfaceSourceAndroidNativeWindow) = %d, want 24", got)
	}
	if got := unsafe.Offsetof(source.window); got != 16 {
		t.Fatalf("offsetof(window) = %d, want 16", got)
	}
	if source.chain.Next != 0 {
		t.Fatalf("chain.Next = %#x, want 0", source.chain.Next)
	}
	if source.chain.SType != uint32(STypeSurfaceSourceAndroidNativeWindow) {
		t.Fatalf("chain.SType = %#x, want %#x", source.chain.SType, uint32(STypeSurfaceSourceAndroidNativeWindow))
	}
	if source.window != 0x1234 {
		t.Fatalf("window = %#x, want 0x1234", source.window)
	}
}

func TestSurfaceSourceAndroidNativeWindowRejectsZero(t *testing.T) {
	_, err := newSurfaceSourceAndroidNativeWindow(0)
	if err == nil {
		t.Fatal("newSurfaceSourceAndroidNativeWindow(0) succeeded")
	}
}
