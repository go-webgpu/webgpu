//go:build darwin

package wgpu

import (
	"unsafe"
)

// surfaceSourceMetalLayer is the native structure for macOS surface creation - 24 bytes.
type surfaceSourceMetalLayer struct {
	chain ChainedStruct // 16 bytes: next (8) + sType (4) + padding (4)
	layer uintptr       // 8 bytes - CAMetalLayer*
}

// CreateSurfaceFromMetalLayer creates a surface from a CAMetalLayer.
// layer is the CAMetalLayer pointer from the native macOS view.
func (inst *Instance) CreateSurfaceFromMetalLayer(layer uintptr) (*Surface, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	// Build WGPUSurfaceSourceMetalLayer
	source := surfaceSourceMetalLayer{
		chain: ChainedStruct{
			Next:  0,
			SType: uint32(STypeSurfaceSourceMetalLayer),
		},
		layer: layer,
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

	return &Surface{handle: handle}, nil
}
