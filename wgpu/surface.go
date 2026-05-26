package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// surfaceDescriptor is the native structure for surface creation.
type surfaceDescriptor struct {
	nextInChain uintptr    // Pointer to platform-specific source
	label       StringView // 16 bytes
}

// surfaceConfigurationWire is the FFI-compatible structure for configuring a surface.
// Uses uint32 for format (converted from gputypes) and uint64 for usage.
type surfaceConfigurationWire struct {
	nextInChain     uintptr // 8 bytes
	device          uintptr // 8 bytes (WGPUDevice handle)
	format          uint32  // 4 bytes (converted from gputypes.TextureFormat)
	_pad1           [4]byte // 4 bytes padding
	usage           uint64  // 8 bytes (TextureUsage as uint64)
	width           uint32  // 4 bytes
	height          uint32  // 4 bytes
	viewFormatCount uintptr // 8 bytes (size_t)
	viewFormats     uintptr // 8 bytes (pointer)
	alphaMode       uint32  // 4 bytes (CompositeAlphaMode)
	presentMode     uint32  // 4 bytes (PresentMode)
}

// surfaceTexture is the native structure returned by GetCurrentTexture.
type surfaceTexture struct {
	nextInChain uintptr                        // 8 bytes
	texture     uintptr                        // 8 bytes (WGPUTexture)
	status      SurfaceGetCurrentTextureStatus // 4 bytes
	_pad        [4]byte                        // 4 bytes padding
}

// surfaceCapabilitiesWire is the FFI-compatible structure for WGPUSurfaceCapabilities.
// Matches C struct layout from wgpu-native v27.
type surfaceCapabilitiesWire struct {
	nextInChain      uintptr // 8 bytes (WGPUChainedStructOut*)
	usages           uint64  // 8 bytes (WGPUTextureUsage bitflags)
	formatCount      uintptr // 8 bytes (size_t)
	formats          uintptr // 8 bytes (WGPUTextureFormat* - pointer to array)
	presentModeCount uintptr // 8 bytes (size_t)
	presentModes     uintptr // 8 bytes (WGPUPresentMode* - pointer to array)
	alphaModeCount   uintptr // 8 bytes (size_t)
	alphaModes       uintptr // 8 bytes (WGPUCompositeAlphaMode* - pointer to array)
}

// SurfaceConfiguration describes how to configure a surface.
// Note: the Device field is deprecated — pass the device as a separate argument to Configure.
// It remains here for backward compatibility; if non-nil it takes precedence over the explicit arg.
type SurfaceConfiguration struct {
	// Device is deprecated: pass the device to Configure() directly instead.
	// Kept for backward compatibility. If non-nil, overrides the explicit device argument.
	Device      *Device
	Format      gputypes.TextureFormat
	Usage       gputypes.TextureUsage
	Width       uint32
	Height      uint32
	AlphaMode   gputypes.CompositeAlphaMode
	PresentMode gputypes.PresentMode
}

// SurfaceTexture holds the result of GetCurrentTexture.
type SurfaceTexture struct {
	Texture *Texture
	Status  SurfaceGetCurrentTextureStatus
}

// SurfaceCapabilities describes the capabilities of a surface for presentation.
// Returned by Surface.GetCapabilities() to query supported formats, present modes, etc.
type SurfaceCapabilities struct {
	Usages       gputypes.TextureUsage
	Formats      []gputypes.TextureFormat
	PresentModes []gputypes.PresentMode
	AlphaModes   []gputypes.CompositeAlphaMode
}

// Error values for surface operations.
// These are sentinel errors for programmatic error handling via errors.Is().
var (
	ErrSurfaceNeedsReconfigure = &WGPUError{Op: "Surface.GetCurrentTexture", Message: "surface needs reconfigure"}
	ErrSurfaceLost             = &WGPUError{Op: "Surface.GetCurrentTexture", Message: "surface lost"}
	ErrSurfaceTimeout          = &WGPUError{Op: "Surface.GetCurrentTexture", Message: "surface texture timeout"}
	// ErrSurfaceOccluded is returned on macOS Metal when the window is minimized or fully covered.
	// Applications should skip rendering for the current frame and try again when unoccluded.
	// New in wgpu-native v29.
	ErrSurfaceOccluded = &WGPUError{Op: "Surface.GetCurrentTexture", Message: "surface occluded (window minimized or covered)"}
	// ErrSurfaceOutOfMemory is kept for backward compatibility.
	// Deprecated: In v29 this is reported as generic Error status.
	ErrSurfaceOutOfMemory = &WGPUError{Op: "Surface.GetCurrentTexture", Message: "out of memory"}
	// ErrSurfaceDeviceLost is kept for backward compatibility.
	// Deprecated: In v29 this is reported as generic Error status.
	ErrSurfaceDeviceLost = &WGPUError{Op: "Surface.GetCurrentTexture", Message: "device lost"}
)

// Configure configures the surface for rendering.
// The device argument specifies which logical device to use for the surface.
// If config.Device is also set (deprecated usage), it takes precedence over the device arg.
// Returns nil on success. Errors are surfaced through the Device uncaptured-error callback
// in this FFI implementation; the error return matches the gogpu/wgpu API signature.
// This replaces the deprecated SwapChain API.
// Enum values are converted from gputypes to wgpu-native values before FFI call.
func (s *Surface) Configure(device *Device, config *SurfaceConfiguration) error {
	mustInit()
	if s == nil || s.handle == 0 || config == nil {
		return nil
	}

	// config.Device takes precedence (backward compat) over the device argument.
	dev := device
	if config.Device != nil {
		dev = config.Device
	}
	if dev == nil || dev.handle == 0 {
		return nil
	}

	nativeConfig := surfaceConfigurationWire{
		nextInChain:     0,
		device:          dev.handle,
		format:          uint32(config.Format),
		usage:           uint64(config.Usage),
		width:           config.Width,
		height:          config.Height,
		viewFormatCount: 0,
		viewFormats:     0,
		alphaMode:       uint32(config.AlphaMode),
		presentMode:     uint32(config.PresentMode),
	}

	procSurfaceConfigure.Call( //nolint:errcheck
		s.handle,
		uintptr(unsafe.Pointer(&nativeConfig)),
	)
	return nil
}

// ConfigureLegacy configures the surface using only the config struct (legacy API).
// Deprecated: use Configure(device, config) instead.
func (s *Surface) ConfigureLegacy(config *SurfaceConfiguration) {
	_ = s.Configure(nil, config)
}

// Unconfigure removes the surface configuration.
func (s *Surface) Unconfigure() {
	mustInit()
	if s == nil || s.handle == 0 {
		return
	}
	procSurfaceUnconfigure.Call(s.handle) //nolint:errcheck
}

// GetCurrentTexture gets the current texture to render to.
// Returns the texture, a suboptimal flag (true if the surface needs reconfiguration
// but is still usable this frame), and any error. This matches the gogpu/wgpu API.
func (s *Surface) GetCurrentTexture() (*SurfaceTexture, bool, error) {
	if err := checkInit(); err != nil {
		return nil, false, err
	}
	if s == nil || s.handle == 0 {
		return nil, false, &WGPUError{Op: "Surface.GetCurrentTexture", Message: "surface is nil or released"}
	}

	var surfTex surfaceTexture

	procSurfaceGetCurrentTexture.Call( //nolint:errcheck
		s.handle,
		uintptr(unsafe.Pointer(&surfTex)),
	)

	result := &SurfaceTexture{
		Texture: &Texture{handle: surfTex.texture},
		Status:  surfTex.status,
	}

	switch surfTex.status {
	case SurfaceGetCurrentTextureStatusSuccessOptimal:
		return result, false, nil
	case SurfaceGetCurrentTextureStatusSuccessSuboptimal:
		// Surface still usable but caller should reconfigure soon.
		return result, true, nil
	case SurfaceGetCurrentTextureStatusOutdated:
		return result, false, ErrSurfaceNeedsReconfigure
	case SurfaceGetCurrentTextureStatusLost:
		return nil, false, ErrSurfaceLost
	case SurfaceGetCurrentTextureStatusTimeout:
		return nil, false, ErrSurfaceTimeout
	case NativeSurfaceGetCurrentTextureStatusOccluded:
		// wgpu-native v29: window is occluded/minimized (Metal backend only).
		// No texture is returned; caller should skip this frame and try again.
		return nil, false, ErrSurfaceOccluded
	default:
		// v29: SurfaceGetCurrentTextureStatusError (0x06) covers all error cases
		// including former OutOfMemory (0x06) and DeviceLost (0x07).
		return nil, false, &WGPUError{Op: "Surface.GetCurrentTexture", Message: "failed to get surface texture"}
	}
}

// Present presents the current frame to the surface.
// The texture argument is accepted for API compatibility with gogpu/wgpu but
// is unused in the FFI implementation (wgpuSurfacePresent takes no texture arg).
// Returns nil on success.
func (s *Surface) Present(texture ...*SurfaceTexture) error {
	mustInit()
	if s == nil || s.handle == 0 {
		return nil
	}
	procSurfacePresent.Call(s.handle) //nolint:errcheck
	return nil
}

// Release releases the surface.
func (s *Surface) Release() {
	if s.handle != 0 {
		untrackResource(s.handle)
		procSurfaceRelease.Call(s.handle) //nolint:errcheck
		s.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (s *Surface) Handle() uintptr { return s.handle }

// GetCapabilities queries the surface capabilities for the given adapter.
// This determines which texture formats, present modes, and alpha modes are supported.
// The caller must provide a valid adapter that will be used with this surface.
func (s *Surface) GetCapabilities(adapter *Adapter) (*SurfaceCapabilities, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	if s == nil || s.handle == 0 {
		return nil, &WGPUError{Op: "Surface.GetCapabilities", Message: "surface is nil"}
	}
	if adapter == nil || adapter.handle == 0 {
		return nil, &WGPUError{Op: "Surface.GetCapabilities", Message: "adapter is nil"}
	}

	// Call wgpuSurfaceGetCapabilities
	var wire surfaceCapabilitiesWire
	procSurfaceGetCapabilities.Call( //nolint:errcheck
		s.handle,
		adapter.handle,
		uintptr(unsafe.Pointer(&wire)),
	)

	// Convert wire struct to Go struct
	caps := &SurfaceCapabilities{
		Usages: gputypes.TextureUsage(wire.usages),
	}

	// Convert formats array
	if wire.formatCount > 0 && wire.formats != 0 {
		rawFormats := unsafe.Slice((*uint32)(ptrFromUintptr(wire.formats)), wire.formatCount)
		caps.Formats = make([]gputypes.TextureFormat, len(rawFormats))
		for i, f := range rawFormats {
			caps.Formats[i] = gputypes.TextureFormat(f)
		}
	}

	// Convert present modes array
	if wire.presentModeCount > 0 && wire.presentModes != 0 {
		rawPresentModes := unsafe.Slice((*uint32)(ptrFromUintptr(wire.presentModes)), wire.presentModeCount)
		caps.PresentModes = make([]gputypes.PresentMode, len(rawPresentModes))
		for i, pm := range rawPresentModes {
			caps.PresentModes[i] = gputypes.PresentMode(pm)
		}
	}

	// Convert alpha modes array
	if wire.alphaModeCount > 0 && wire.alphaModes != 0 {
		rawAlphaModes := unsafe.Slice((*uint32)(ptrFromUintptr(wire.alphaModes)), wire.alphaModeCount)
		caps.AlphaModes = make([]gputypes.CompositeAlphaMode, len(rawAlphaModes))
		for i, am := range rawAlphaModes {
			caps.AlphaModes[i] = gputypes.CompositeAlphaMode(am)
		}
	}

	// Free C memory allocated by wgpu-native
	procSurfaceCapabilitiesFreeMembers.Call(uintptr(unsafe.Pointer(&wire))) //nolint:errcheck

	return caps, nil
}
