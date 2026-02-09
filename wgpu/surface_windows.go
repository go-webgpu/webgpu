//go:build windows

package wgpu

import (
	"unsafe"
)

// surfaceSourceWindowsHWND is the native structure for Windows surface creation - 32 bytes.
type surfaceSourceWindowsHWND struct {
	chain     ChainedStruct // 16 bytes: next (8) + sType (4) + padding (4)
	hinstance uintptr       // 8 bytes - HINSTANCE
	hwnd      uintptr       // 8 bytes - HWND
}

// CreateSurfaceFromWindowsHWND creates a surface from a Windows HWND.
// hinstance should be the HINSTANCE of the application (can be 0).
// hwnd is the window handle to create the surface for.
func (inst *Instance) CreateSurfaceFromWindowsHWND(hinstance, hwnd uintptr) (*Surface, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	// Build WGPUSurfaceSourceWindowsHWND
	source := surfaceSourceWindowsHWND{
		chain: ChainedStruct{
			Next:  0,
			SType: uint32(STypeSurfaceSourceWindowsHWND),
		},
		hinstance: hinstance,
		hwnd:      hwnd,
	}

	// Build WGPUSurfaceDescriptor with source chained
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
