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
	if inst == nil || inst.handle == 0 {
		return nil, &WGPUError{Op: "CreateSurface", Message: "instance is nil or released"}
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

// surfaceSourceSwapChainPanel matches WGPUSurfaceSourceSwapChainPanel in the
// wgpu-native v29 header - 24 bytes.
type surfaceSourceSwapChainPanel struct {
	chain       ChainedStruct // 16 bytes
	panelNative uintptr       // 8 bytes - ISwapChainPanelNative*
}

// swapChainPanelSourceSink keeps the last chained source on the heap: the
// descriptor refers to it only through a uintptr.
var swapChainPanelSourceSink *surfaceSourceSwapChainPanel

// CreateSurfaceFromSwapChainPanel creates a surface from an
// ISwapChainPanelNative pointer (WGPUSurfaceSourceSwapChainPanel, DX12 only).
//
// wgpu-native creates the swap chain with CreateSwapChainForComposition and
// hands it to ISwapChainPanelNative::SetSwapChain. Besides a WinUI
// SwapChainPanel, the pointer can be any object implementing that interface,
// which lets the caller place the swap chain in its own DirectComposition
// tree (IDCompositionVisual::SetContent), e.g. to keep other DComp content
// such as a WebView2 composition controller visible above the swap chain in
// fullscreen, where an HWND swap chain would be promoted to independent flip.
// panelNative must stay valid while the surface lives: wgpu-native keeps a
// reference to it.
func (inst *Instance) CreateSurfaceFromSwapChainPanel(panelNative uintptr) (*Surface, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if inst == nil || inst.handle == 0 {
		return nil, &WGPUError{Op: "CreateSurface", Message: "instance is nil or released"}
	}
	if panelNative == 0 {
		return nil, &WGPUError{Op: "CreateSurface", Message: "swap chain panel is nil"}
	}

	source := &surfaceSourceSwapChainPanel{
		chain: ChainedStruct{
			Next:  0,
			SType: uint32(STypeSurfaceSourceSwapChainPanel),
		},
		panelNative: panelNative,
	}
	swapChainPanelSourceSink = source

	desc := surfaceDescriptor{
		nextInChain: uintptr(unsafe.Pointer(source)),
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
