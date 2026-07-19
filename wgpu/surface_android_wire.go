package wgpu

// surfaceSourceAndroidNativeWindow matches
// WGPUSurfaceSourceAndroidNativeWindow in the WebGPU native v29 header.
// It is host-buildable so ordinary CI can verify the Android ABI layout.
type surfaceSourceAndroidNativeWindow struct {
	chain  ChainedStruct
	window uintptr // ANativeWindow*
}

func newSurfaceSourceAndroidNativeWindow(window uintptr) (surfaceSourceAndroidNativeWindow, error) {
	if window == 0 {
		return surfaceSourceAndroidNativeWindow{}, &WGPUError{
			Op:      "CreateSurface",
			Message: "Android native window is nil",
		}
	}

	return surfaceSourceAndroidNativeWindow{
		chain: ChainedStruct{
			SType: uint32(STypeSurfaceSourceAndroidNativeWindow),
		},
		window: window,
	}, nil
}
