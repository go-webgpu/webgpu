//go:build android

package wgpu

import (
	"unsafe"
)

// CreateSurfaceFromAndroidNativeWindow creates a surface from an
// ANativeWindow. The caller must keep its ANativeWindow reference alive until
// the returned Surface is released.
func (inst *Instance) CreateSurfaceFromAndroidNativeWindow(window uintptr) (*Surface, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if inst == nil || inst.handle == 0 {
		return nil, &WGPUError{Op: "CreateSurface", Message: "instance is nil or released"}
	}

	source, err := newSurfaceSourceAndroidNativeWindow(window)
	if err != nil {
		return nil, err
	}
	desc := surfaceDescriptor{
		nextInChain: uintptr(unsafe.Pointer(&source)),
		label:       EmptyStringView(),
	}

	handle, _, _ := procInstanceCreateSurface.Call(
		inst.handle,
		uintptr(unsafe.Pointer(&desc)),
	)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreateSurface", Message: "failed to create surface"}
	}

	trackResource(handle, "Surface")
	return &Surface{handle: handle}, nil
}
