package wgpu

// surfaceSourceAndroidNativeWindow matches
// WGPUSurfaceSourceAndroidNativeWindow in the WebGPU native v29 header.
// It is host-buildable so ordinary CI can verify the Android ABI layout.
type surfaceSourceAndroidNativeWindow struct {
	chain  ChainedStruct // 16 bytes: next (8) + sType (4) + padding (4)
	window uintptr       // 8 bytes - ANativeWindow*
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
			Next:  0,
			SType: uint32(STypeSurfaceSourceAndroidNativeWindow),
		},
		window: window,
	}, nil
}
