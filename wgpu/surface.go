package wgpu

import (
	"errors"
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
	nextInChain     uintptr // 8 bytes (WGPUChainedStructOut*)
	usages          uint64  // 8 bytes (WGPUTextureUsage bitflags)
	formatCount     uintptr // 8 bytes (size_t)
	formats         uintptr // 8 bytes (WGPUTextureFormat* - pointer to array)
	presentModeCount uintptr // 8 bytes (size_t)
	presentModes    uintptr // 8 bytes (WGPUPresentMode* - pointer to array)
	alphaModeCount  uintptr // 8 bytes (size_t)
	alphaModes      uintptr // 8 bytes (WGPUCompositeAlphaMode* - pointer to array)
}

// SurfaceConfiguration describes how to configure a surface.
type SurfaceConfiguration struct {
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
var (
	ErrSurfaceNeedsReconfigure = errors.New("wgpu: surface needs reconfigure")
	ErrSurfaceLost             = errors.New("wgpu: surface lost")
	ErrSurfaceTimeout          = errors.New("wgpu: surface texture timeout")
	ErrSurfaceOutOfMemory      = errors.New("wgpu: out of memory")
	ErrSurfaceDeviceLost       = errors.New("wgpu: device lost")
)

// Configure configures the surface for rendering.
// This replaces the deprecated SwapChain API.
// Enum values are converted from gputypes to wgpu-native values before FFI call.
func (s *Surface) Configure(config *SurfaceConfiguration) {
	mustInit()

	nativeConfig := surfaceConfigurationWire{
		nextInChain:     0,
		device:          config.Device.handle,
		format:          toWGPUTextureFormat(config.Format),
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
}

// Unconfigure removes the surface configuration.
func (s *Surface) Unconfigure() {
	mustInit()
	procSurfaceUnconfigure.Call(s.handle) //nolint:errcheck
}

// GetCurrentTexture gets the current texture to render to.
// Returns the texture and its status. Check status before using the texture.
func (s *Surface) GetCurrentTexture() (*SurfaceTexture, error) {
	mustInit()

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
	case SurfaceGetCurrentTextureStatusSuccessOptimal,
		SurfaceGetCurrentTextureStatusSuccessSuboptimal:
		return result, nil
	case SurfaceGetCurrentTextureStatusOutdated:
		return result, ErrSurfaceNeedsReconfigure
	case SurfaceGetCurrentTextureStatusLost:
		return nil, ErrSurfaceLost
	case SurfaceGetCurrentTextureStatusTimeout:
		return nil, ErrSurfaceTimeout
	case SurfaceGetCurrentTextureStatusOutOfMemory:
		return nil, ErrSurfaceOutOfMemory
	case SurfaceGetCurrentTextureStatusDeviceLost:
		return nil, ErrSurfaceDeviceLost
	default:
		return nil, errors.New("wgpu: failed to get surface texture")
	}
}

// Present presents the current frame to the surface.
func (s *Surface) Present() {
	mustInit()
	procSurfacePresent.Call(s.handle) //nolint:errcheck
}

// Release releases the surface.
func (s *Surface) Release() {
	if s.handle != 0 {
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
	mustInit()

	if s == nil || s.handle == 0 {
		return nil, errors.New("wgpu: surface is nil")
	}
	if adapter == nil || adapter.handle == 0 {
		return nil, errors.New("wgpu: adapter is nil")
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
		rawFormats := unsafe.Slice((*uint32)(unsafe.Pointer(wire.formats)), wire.formatCount)
		caps.Formats = make([]gputypes.TextureFormat, len(rawFormats))
		for i, f := range rawFormats {
			caps.Formats[i] = fromWGPUTextureFormat(f)
		}
	}

	// Convert present modes array
	if wire.presentModeCount > 0 && wire.presentModes != 0 {
		rawPresentModes := unsafe.Slice((*uint32)(unsafe.Pointer(wire.presentModes)), wire.presentModeCount)
		caps.PresentModes = make([]gputypes.PresentMode, len(rawPresentModes))
		for i, pm := range rawPresentModes {
			caps.PresentModes[i] = gputypes.PresentMode(pm)
		}
	}

	// Convert alpha modes array
	if wire.alphaModeCount > 0 && wire.alphaModes != 0 {
		rawAlphaModes := unsafe.Slice((*uint32)(unsafe.Pointer(wire.alphaModes)), wire.alphaModeCount)
		caps.AlphaModes = make([]gputypes.CompositeAlphaMode, len(rawAlphaModes))
		for i, am := range rawAlphaModes {
			caps.AlphaModes[i] = gputypes.CompositeAlphaMode(am)
		}
	}

	// Free C memory allocated by wgpu-native
	procSurfaceCapabilitiesFreeMembers.Call(uintptr(unsafe.Pointer(&wire))) //nolint:errcheck

	return caps, nil
}
