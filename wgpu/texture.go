package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// TextureDescriptor describes a texture to create.
type TextureDescriptor struct {
	NextInChain     uintptr
	Label           StringView
	Usage           gputypes.TextureUsage
	Dimension       gputypes.TextureDimension
	Size            gputypes.Extent3D
	Format          gputypes.TextureFormat
	MipLevelCount   uint32
	SampleCount     uint32
	ViewFormatCount uintptr
	ViewFormats     uintptr
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
	NextInChain     uintptr
	Label           StringView
	Format          gputypes.TextureFormat
	Dimension       gputypes.TextureViewDimension
	BaseMipLevel    uint32
	MipLevelCount   uint32
	BaseArrayLayer  uint32
	ArrayLayerCount uint32
	Aspect          TextureAspect
	_pad            [4]byte //nolint:unused // padding for FFI alignment
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
func (t *Texture) CreateView(desc *TextureViewDescriptor) *TextureView {
	mustInit()

	var descPtr uintptr
	if desc != nil {
		// Convert to wire format with wgpu-native enum values
		wireDesc := textureViewDescriptorWire{
			NextInChain:     desc.NextInChain,
			Label:           desc.Label,
			Format:          toWGPUTextureFormat(desc.Format),
			Dimension:       toWGPUTextureViewDimension(desc.Dimension),
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
		return nil
	}
	return &TextureView{handle: handle}
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
		procTextureRelease.Call(t.handle) //nolint:errcheck
		t.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (t *Texture) Handle() uintptr { return t.handle }

// GetWidth returns the width of the texture in texels.
func (t *Texture) GetWidth() uint32 {
	mustInit()
	if t == nil || t.handle == 0 {
		return 0
	}
	result, _, _ := procTextureGetWidth.Call(t.handle)
	return uint32(result)
}

// GetHeight returns the height of the texture in texels.
func (t *Texture) GetHeight() uint32 {
	mustInit()
	if t == nil || t.handle == 0 {
		return 0
	}
	result, _, _ := procTextureGetHeight.Call(t.handle)
	return uint32(result)
}

// GetDepthOrArrayLayers returns the depth (for 3D textures) or array layer count.
func (t *Texture) GetDepthOrArrayLayers() uint32 {
	mustInit()
	if t == nil || t.handle == 0 {
		return 0
	}
	result, _, _ := procTextureGetDepthOrArrayLayers.Call(t.handle)
	return uint32(result)
}

// GetMipLevelCount returns the number of mip levels.
func (t *Texture) GetMipLevelCount() uint32 {
	mustInit()
	if t == nil || t.handle == 0 {
		return 0
	}
	result, _, _ := procTextureGetMipLevelCount.Call(t.handle)
	return uint32(result)
}

// GetFormat returns the texture format.
// The format is converted from wgpu-native enum to gputypes enum.
func (t *Texture) GetFormat() gputypes.TextureFormat {
	mustInit()
	if t == nil || t.handle == 0 {
		return gputypes.TextureFormatUndefined
	}
	result, _, _ := procTextureGetFormat.Call(t.handle)
	// Convert from wgpu-native enum to gputypes
	return fromWGPUTextureFormat(uint32(result))
}

// Release releases the texture view reference.
func (tv *TextureView) Release() {
	if tv.handle != 0 {
		procTextureViewRelease.Call(tv.handle) //nolint:errcheck
		tv.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (tv *TextureView) Handle() uintptr { return tv.handle }

// CreateTexture creates a texture with the specified descriptor.
// Enum values are converted from gputypes to wgpu-native values before FFI call.
func (d *Device) CreateTexture(desc *TextureDescriptor) *Texture {
	mustInit()
	if desc == nil {
		return nil
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

	// Convert to wire format with wgpu-native enum values
	wireDesc := textureDescriptorWire{
		NextInChain:     desc.NextInChain,
		Label:           desc.Label,
		Usage:           uint64(desc.Usage), // bitflags, uint64 in wgpu-native
		Dimension:       toWGPUTextureDimension(desc.Dimension),
		Size:            desc.Size,
		Format:          toWGPUTextureFormat(desc.Format),
		MipLevelCount:   mipLevelCount,
		SampleCount:     sampleCount,
		ViewFormatCount: desc.ViewFormatCount,
		ViewFormats:     desc.ViewFormats,
	}

	handle, _, _ := procDeviceCreateTexture.Call(
		d.handle,
		uintptr(unsafe.Pointer(&wireDesc)),
	)
	if handle == 0 {
		return nil
	}
	return &Texture{handle: handle}
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
	if dest == nil || layout == nil || size == nil || len(data) == 0 {
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
