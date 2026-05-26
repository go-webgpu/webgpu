package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// TextureDescriptor describes a texture to create.
type TextureDescriptor struct {
	Label         string
	Usage         gputypes.TextureUsage
	Dimension     gputypes.TextureDimension
	Size          gputypes.Extent3D
	Format        gputypes.TextureFormat
	MipLevelCount uint32
	SampleCount   uint32
	ViewFormats   []gputypes.TextureFormat
}

// textureDescriptorWire is the FFI-compatible struct with wgpu-native enum values.
// CRITICAL: Usage is uint64 because wgpu-native defines WGPUTextureUsageFlags as uint64!
type textureDescriptorWire struct {
	NextInChain     uintptr
	Label           StringView
	Usage           uint64 // TextureUsage bitflags (uint64 in wgpu-native!)
	Dimension       uint32 // TextureDimension (needs +1 shift)
	Size            gputypes.Extent3D
	Format          uint32 // TextureFormat (converted via map)
	MipLevelCount   uint32
	SampleCount     uint32
	ViewFormatCount uintptr
	ViewFormats     uintptr
}

// TextureViewDescriptor describes a texture view to create.
type TextureViewDescriptor struct {
	Label           string
	Format          gputypes.TextureFormat
	Dimension       gputypes.TextureViewDimension
	BaseMipLevel    uint32
	MipLevelCount   uint32
	BaseArrayLayer  uint32
	ArrayLayerCount uint32
	Aspect          TextureAspect
	Usage           gputypes.TextureUsage
}

// textureViewDescriptorWire is the FFI-compatible struct with wgpu-native enum values.
// CRITICAL: Usage is uint64 because wgpu-native defines WGPUTextureUsageFlags as uint64!
type textureViewDescriptorWire struct {
	NextInChain     uintptr
	Label           StringView
	Format          uint32 // TextureFormat (converted)
	Dimension       uint32 // TextureViewDimension (needs +1 shift)
	BaseMipLevel    uint32
	MipLevelCount   uint32
	BaseArrayLayer  uint32
	ArrayLayerCount uint32
	Aspect          TextureAspect
	_pad            [4]byte
	Usage           uint64 // TextureUsage bitflags (uint64 in wgpu-native!)
}

// CreateView creates a view into this texture.
// Pass nil for default view parameters.
// Enum values are converted from gputypes to wgpu-native values before FFI call.
// Returns an error if the FFI call fails or the texture is nil.
func (t *Texture) CreateView(desc *TextureViewDescriptor) (*TextureView, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if t == nil || t.handle == 0 {
		return nil, &WGPUError{Op: "CreateView", Message: "texture is nil or released"}
	}

	var descPtr uintptr
	if desc != nil {
		// Convert Go-idiomatic descriptor to FFI wire format
		wireDesc := textureViewDescriptorWire{
			Label:           stringToStringView(desc.Label),
			Format:          uint32(desc.Format),
			Dimension:       uint32(desc.Dimension),
			BaseMipLevel:    desc.BaseMipLevel,
			MipLevelCount:   desc.MipLevelCount,
			BaseArrayLayer:  desc.BaseArrayLayer,
			ArrayLayerCount: desc.ArrayLayerCount,
			Aspect:          desc.Aspect,
			Usage:           uint64(desc.Usage), // bitflags, uint64 in wgpu-native
		}
		descPtr = uintptr(unsafe.Pointer(&wireDesc))
	}

	handle, _, _ := procTextureCreateView.Call(
		t.handle,
		descPtr,
	)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreateView", Message: "wgpu returned null handle"}
	}
	trackResource(handle, "TextureView")
	return &TextureView{handle: handle}, nil
}

// Destroy destroys the texture.
func (t *Texture) Destroy() {
	mustInit()
	if t.handle != 0 {
		procTextureDestroy.Call(t.handle) //nolint:errcheck
	}
}

// Release releases the texture reference.
func (t *Texture) Release() {
	if t.handle != 0 {
		untrackResource(t.handle)
		procTextureRelease.Call(t.handle) //nolint:errcheck
		t.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (t *Texture) Handle() uintptr { return t.handle }

// Width returns the width of the texture in texels.
func (t *Texture) Width() uint32 {
	mustInit()
	if t == nil || t.handle == 0 {
		return 0
	}
	result, _, _ := procTextureGetWidth.Call(t.handle)
	return uint32(result)
}

// Height returns the height of the texture in texels.
func (t *Texture) Height() uint32 {
	mustInit()
	if t == nil || t.handle == 0 {
		return 0
	}
	result, _, _ := procTextureGetHeight.Call(t.handle)
	return uint32(result)
}

// DepthOrArrayLayers returns the depth (for 3D textures) or array layer count.
func (t *Texture) DepthOrArrayLayers() uint32 {
	mustInit()
	if t == nil || t.handle == 0 {
		return 0
	}
	result, _, _ := procTextureGetDepthOrArrayLayers.Call(t.handle)
	return uint32(result)
}

// MipLevelCount returns the number of mip levels.
func (t *Texture) MipLevelCount() uint32 {
	mustInit()
	if t == nil || t.handle == 0 {
		return 0
	}
	result, _, _ := procTextureGetMipLevelCount.Call(t.handle)
	return uint32(result)
}

// Format returns the texture format.
// TextureFormat values match between gputypes v0.3.0 and wgpu-native v29 exactly.
func (t *Texture) Format() gputypes.TextureFormat {
	mustInit()
	if t == nil || t.handle == 0 {
		return gputypes.TextureFormatUndefined
	}
	result, _, _ := procTextureGetFormat.Call(t.handle)
	return gputypes.TextureFormat(result)
}

// Release releases the texture view reference.
func (tv *TextureView) Release() {
	if tv.handle != 0 {
		untrackResource(tv.handle)
		procTextureViewRelease.Call(tv.handle) //nolint:errcheck
		tv.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (tv *TextureView) Handle() uintptr { return tv.handle }

// CreateTexture creates a texture with the specified descriptor.
// Enum values are converted from gputypes to wgpu-native values before FFI call.
// Returns an error if the FFI call fails or the device/descriptor is nil.
func (d *Device) CreateTexture(desc *TextureDescriptor) (*Texture, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if d == nil || d.handle == 0 {
		return nil, &WGPUError{Op: "CreateTexture", Message: "device is nil or released"}
	}
	if desc == nil {
		return nil, &WGPUError{Op: "CreateTexture", Message: "descriptor is nil"}
	}

	// wgpu-native requires MipLevelCount >= 1 and SampleCount >= 1
	mipLevelCount := desc.MipLevelCount
	if mipLevelCount == 0 {
		mipLevelCount = 1
	}
	sampleCount := desc.SampleCount
	if sampleCount == 0 {
		sampleCount = 1
	}

	// Convert []TextureFormat → []uint32 for FFI (values match, but wire struct needs uint32 pointer)
	var viewFormatCount uintptr
	var viewFormatsPtr uintptr
	if len(desc.ViewFormats) > 0 {
		// Convert to uint32 slice (gputypes values equal wgpu-native values)
		wireFormats := make([]uint32, len(desc.ViewFormats))
		for i, f := range desc.ViewFormats {
			wireFormats[i] = uint32(f)
		}
		viewFormatCount = uintptr(len(wireFormats))
		viewFormatsPtr = uintptr(unsafe.Pointer(&wireFormats[0]))
	}

	// Convert to wire format with wgpu-native enum values
	wireDesc := textureDescriptorWire{
		Label:           stringToStringView(desc.Label),
		Usage:           uint64(desc.Usage), // bitflags, uint64 in wgpu-native
		Dimension:       uint32(desc.Dimension),
		Size:            desc.Size,
		Format:          uint32(desc.Format),
		MipLevelCount:   mipLevelCount,
		SampleCount:     sampleCount,
		ViewFormatCount: viewFormatCount,
		ViewFormats:     viewFormatsPtr,
	}

	handle, _, _ := procDeviceCreateTexture.Call(
		d.handle,
		uintptr(unsafe.Pointer(&wireDesc)),
	)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreateTexture", Message: "wgpu returned null handle"}
	}
	trackResource(handle, "Texture")
	return &Texture{handle: handle}, nil
}

// TexelCopyTextureInfo describes a texture for WriteTexture.
type TexelCopyTextureInfo struct {
	Texture  uintptr
	MipLevel uint32
	Origin   gputypes.Origin3D
	Aspect   TextureAspect
}

// TexelCopyBufferLayout describes buffer layout for WriteTexture.
type TexelCopyBufferLayout struct {
	Offset       uint64
	BytesPerRow  uint32
	RowsPerImage uint32
}

// TexelCopyBufferInfo describes a buffer source/destination for copy operations.
type TexelCopyBufferInfo struct {
	Layout TexelCopyBufferLayout
	Buffer uintptr // Buffer handle
}

// WriteTexture writes data to a texture.
func (q *Queue) WriteTexture(dest *TexelCopyTextureInfo, data []byte, layout *TexelCopyBufferLayout, size *gputypes.Extent3D) {
	mustInit()
	if q == nil || q.handle == 0 || dest == nil || layout == nil || size == nil || len(data) == 0 {
		return
	}
	procQueueWriteTexture.Call( //nolint:errcheck
		q.handle,
		uintptr(unsafe.Pointer(dest)),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(layout)),
		uintptr(unsafe.Pointer(size)),
	)
}
