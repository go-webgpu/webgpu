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
