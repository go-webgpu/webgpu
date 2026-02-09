//go:build linux

package wgpu

import (
	"unsafe"
)

// surfaceSourceXlibWindow is the native structure for X11 via Xlib surface creation - 32 bytes.
type surfaceSourceXlibWindow struct {
	chain   ChainedStruct // 16 bytes: next (8) + sType (4) + padding (4)
	display uintptr       // 8 bytes - Display*
	window  uint64        // 8 bytes - Window (XID on 64-bit)
}

// surfaceSourceWaylandSurface is the native structure for Wayland surface creation - 32 bytes.
type surfaceSourceWaylandSurface struct {
	chain   ChainedStruct // 16 bytes: next (8) + sType (4) + padding (4)
	display uintptr       // 8 bytes - wl_display*
	surface uintptr       // 8 bytes - wl_surface*
}

// CreateSurfaceFromXlibWindow creates a surface from an X11 Xlib window.
// display is the X11 Display pointer.
// window is the X11 Window ID (XID).
func (inst *Instance) CreateSurfaceFromXlibWindow(display uintptr, window uint64) (*Surface, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	// Build WGPUSurfaceSourceXlibWindow
	source := surfaceSourceXlibWindow{
		chain: ChainedStruct{
			Next:  0,
			SType: uint32(STypeSurfaceSourceXlibWindow),
		},
		display: display,
		window:  window,
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

// CreateSurfaceFromWaylandSurface creates a surface from a Wayland surface.
// display is the wl_display pointer.
// surface is the wl_surface pointer.
func (inst *Instance) CreateSurfaceFromWaylandSurface(display, surface uintptr) (*Surface, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	// Build WGPUSurfaceSourceWaylandSurface
	source := surfaceSourceWaylandSurface{
		chain: ChainedStruct{
			Next:  0,
			SType: uint32(STypeSurfaceSourceWaylandSurface),
		},
		display: display,
		surface: surface,
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
