package wgpu

import (
	"unsafe"
)

// ShaderStage identifies which shader stages can access a binding.
type ShaderStage uint64

const (
	ShaderStageNone     ShaderStage = 0x0000000000000000
	ShaderStageVertex   ShaderStage = 0x0000000000000001
	ShaderStageFragment ShaderStage = 0x0000000000000002
	ShaderStageCompute  ShaderStage = 0x0000000000000004
)

// BufferBindingType describes how a buffer is bound.
type BufferBindingType uint32

const (
	BufferBindingTypeBindingNotUsed  BufferBindingType = 0x00000000
	BufferBindingTypeUndefined       BufferBindingType = 0x00000001
	BufferBindingTypeUniform         BufferBindingType = 0x00000002
	BufferBindingTypeStorage         BufferBindingType = 0x00000003
	BufferBindingTypeReadOnlyStorage BufferBindingType = 0x00000004
)

// SamplerBindingType describes how a sampler is bound.
type SamplerBindingType uint32

const (
	SamplerBindingTypeBindingNotUsed SamplerBindingType = 0x00000000
	SamplerBindingTypeUndefined      SamplerBindingType = 0x00000001
	SamplerBindingTypeFiltering      SamplerBindingType = 0x00000002
	SamplerBindingTypeNonFiltering   SamplerBindingType = 0x00000003
	SamplerBindingTypeComparison     SamplerBindingType = 0x00000004
)

// TextureSampleType describes how a texture is sampled.
type TextureSampleType uint32

const (
	TextureSampleTypeBindingNotUsed    TextureSampleType = 0x00000000
	TextureSampleTypeUndefined         TextureSampleType = 0x00000001
	TextureSampleTypeFloat             TextureSampleType = 0x00000002
	TextureSampleTypeUnfilterableFloat TextureSampleType = 0x00000003
	TextureSampleTypeDepth             TextureSampleType = 0x00000004
	TextureSampleTypeSint              TextureSampleType = 0x00000005
	TextureSampleTypeUint              TextureSampleType = 0x00000006
)

// TextureViewDimension describes the dimension of a texture view.
type TextureViewDimension uint32

const (
	TextureViewDimensionUndefined TextureViewDimension = 0x00000000
	TextureViewDimension1D        TextureViewDimension = 0x00000001
	TextureViewDimension2D        TextureViewDimension = 0x00000002
	TextureViewDimension2DArray   TextureViewDimension = 0x00000003
	TextureViewDimensionCube      TextureViewDimension = 0x00000004
	TextureViewDimensionCubeArray TextureViewDimension = 0x00000005
	TextureViewDimension3D        TextureViewDimension = 0x00000006
)

// BufferBindingLayout describes buffer binding properties.
type BufferBindingLayout struct {
	NextInChain      uintptr // *ChainedStruct
	Type             BufferBindingType
	HasDynamicOffset Bool
	MinBindingSize   uint64
}

// SamplerBindingLayout describes sampler binding properties.
type SamplerBindingLayout struct {
	NextInChain uintptr // *ChainedStruct
	Type        SamplerBindingType
}

// TextureBindingLayout describes texture binding properties.
type TextureBindingLayout struct {
	NextInChain   uintptr // *ChainedStruct
	SampleType    TextureSampleType
	ViewDimension TextureViewDimension
	Multisampled  Bool
}

// StorageTextureAccess describes storage texture access mode.
type StorageTextureAccess uint32

const (
	StorageTextureAccessBindingNotUsed StorageTextureAccess = 0x00000000
	StorageTextureAccessUndefined      StorageTextureAccess = 0x00000001
	StorageTextureAccessWriteOnly      StorageTextureAccess = 0x00000002
	StorageTextureAccessReadOnly       StorageTextureAccess = 0x00000003
	StorageTextureAccessReadWrite      StorageTextureAccess = 0x00000004
)

// TextureFormat describes the format of texture data.
type TextureFormat uint32

const (
	TextureFormatUndefined            TextureFormat = 0x00
	TextureFormatR8Unorm              TextureFormat = 0x01
	TextureFormatR8Snorm              TextureFormat = 0x02
	TextureFormatR8Uint               TextureFormat = 0x03
	TextureFormatR8Sint               TextureFormat = 0x04
	TextureFormatR16Uint              TextureFormat = 0x05
	TextureFormatR16Sint              TextureFormat = 0x06
	TextureFormatR16Float             TextureFormat = 0x07
	TextureFormatRG8Unorm             TextureFormat = 0x08
	TextureFormatRG8Snorm             TextureFormat = 0x09
	TextureFormatRG8Uint              TextureFormat = 0x0A
	TextureFormatRG8Sint              TextureFormat = 0x0B
	TextureFormatR32Float             TextureFormat = 0x0C
	TextureFormatR32Uint              TextureFormat = 0x0D
	TextureFormatR32Sint              TextureFormat = 0x0E
	TextureFormatRG16Uint             TextureFormat = 0x0F
	TextureFormatRG16Sint             TextureFormat = 0x10
	TextureFormatRG16Float            TextureFormat = 0x11
	TextureFormatRGBA8Unorm           TextureFormat = 0x12
	TextureFormatRGBA8UnormSrgb       TextureFormat = 0x13
	TextureFormatRGBA8Snorm           TextureFormat = 0x14
	TextureFormatRGBA8Uint            TextureFormat = 0x15
	TextureFormatRGBA8Sint            TextureFormat = 0x16
	TextureFormatBGRA8Unorm           TextureFormat = 0x17
	TextureFormatBGRA8UnormSrgb       TextureFormat = 0x18
	TextureFormatRGBA16Float          TextureFormat = 0x21
	TextureFormatRGBA32Float          TextureFormat = 0x23
	TextureFormatStencil8             TextureFormat = 0x26
	TextureFormatDepth16Unorm         TextureFormat = 0x27
	TextureFormatDepth24Plus          TextureFormat = 0x28
	TextureFormatDepth24PlusStencil8  TextureFormat = 0x29
	TextureFormatDepth32Float         TextureFormat = 0x2A
	TextureFormatDepth32FloatStencil8 TextureFormat = 0x2B
)

// StorageTextureBindingLayout describes storage texture binding properties.
type StorageTextureBindingLayout struct {
	NextInChain   uintptr // *ChainedStruct
	Access        StorageTextureAccess
	Format        TextureFormat
	ViewDimension TextureViewDimension
}

// BindGroupLayoutEntry describes a single binding in a bind group layout.
type BindGroupLayoutEntry struct {
	NextInChain    uintptr // *ChainedStruct
	Binding        uint32
	Visibility     ShaderStage
	Buffer         BufferBindingLayout
	Sampler        SamplerBindingLayout
	Texture        TextureBindingLayout
	StorageTexture StorageTextureBindingLayout
}

// BindGroupLayoutDescriptor describes a bind group layout.
type BindGroupLayoutDescriptor struct {
	NextInChain uintptr // *ChainedStruct
	Label       StringView
	EntryCount  uintptr // size_t
	Entries     uintptr // *BindGroupLayoutEntry
}

// BindGroupEntry describes a single binding in a bind group.
type BindGroupEntry struct {
	NextInChain uintptr // *ChainedStruct
	Binding     uint32
	Buffer      uintptr // WGPUBuffer (nullable)
	Offset      uint64
	Size        uint64
	Sampler     uintptr // WGPUSampler (nullable)
	TextureView uintptr // WGPUTextureView (nullable)
}

// BindGroupDescriptor describes a bind group.
type BindGroupDescriptor struct {
	NextInChain uintptr // *ChainedStruct
	Label       StringView
	Layout      uintptr // WGPUBindGroupLayout
	EntryCount  uintptr // size_t
	Entries     uintptr // *BindGroupEntry
}

// CreateBindGroupLayout creates a bind group layout.
func (d *Device) CreateBindGroupLayout(desc *BindGroupLayoutDescriptor) *BindGroupLayout {
	mustInit()
	if desc == nil {
		return nil
	}
	handle, _, _ := procDeviceCreateBindGroupLayout.Call(
		d.handle,
		uintptr(unsafe.Pointer(desc)),
	)
	if handle == 0 {
		return nil
	}
	return &BindGroupLayout{handle: handle}
}

// CreateBindGroupLayoutSimple creates a bind group layout with the given entries.
func (d *Device) CreateBindGroupLayoutSimple(entries []BindGroupLayoutEntry) *BindGroupLayout {
	mustInit()
	if len(entries) == 0 {
		return nil
	}
	desc := BindGroupLayoutDescriptor{
		Label:      EmptyStringView(),
		EntryCount: uintptr(len(entries)),
		Entries:    uintptr(unsafe.Pointer(&entries[0])),
	}
	return d.CreateBindGroupLayout(&desc)
}

// Release releases the bind group layout.
func (bgl *BindGroupLayout) Release() {
	if bgl.handle != 0 {
		procBindGroupLayoutRelease.Call(bgl.handle) //nolint:errcheck
		bgl.handle = 0
	}
}

// Handle returns the underlying handle.
func (bgl *BindGroupLayout) Handle() uintptr { return bgl.handle }

// CreateBindGroup creates a bind group.
func (d *Device) CreateBindGroup(desc *BindGroupDescriptor) *BindGroup {
	mustInit()
	if desc == nil {
		return nil
	}
	handle, _, _ := procDeviceCreateBindGroup.Call(
		d.handle,
		uintptr(unsafe.Pointer(desc)),
	)
	if handle == 0 {
		return nil
	}
	return &BindGroup{handle: handle}
}

// CreateBindGroupSimple creates a bind group with buffer entries.
func (d *Device) CreateBindGroupSimple(layout *BindGroupLayout, entries []BindGroupEntry) *BindGroup {
	mustInit()
	if layout == nil || len(entries) == 0 {
		return nil
	}
	desc := BindGroupDescriptor{
		Label:      EmptyStringView(),
		Layout:     layout.handle,
		EntryCount: uintptr(len(entries)),
		Entries:    uintptr(unsafe.Pointer(&entries[0])),
	}
	return d.CreateBindGroup(&desc)
}

// Release releases the bind group.
func (bg *BindGroup) Release() {
	if bg.handle != 0 {
		procBindGroupRelease.Call(bg.handle) //nolint:errcheck
		bg.handle = 0
	}
}

// Handle returns the underlying handle.
func (bg *BindGroup) Handle() uintptr { return bg.handle }

// BufferBindingEntry creates a BindGroupEntry for a buffer.
func BufferBindingEntry(binding uint32, buffer *Buffer, offset, size uint64) BindGroupEntry {
	return BindGroupEntry{
		Binding: binding,
		Buffer:  buffer.handle,
		Offset:  offset,
		Size:    size,
	}
}

// TextureBindingEntry creates a BindGroupEntry for a texture view.
func TextureBindingEntry(binding uint32, textureView *TextureView) BindGroupEntry {
	return BindGroupEntry{
		Binding:     binding,
		TextureView: textureView.handle,
	}
}

// SamplerBindingEntry creates a BindGroupEntry for a sampler.
func SamplerBindingEntry(binding uint32, sampler *Sampler) BindGroupEntry {
	return BindGroupEntry{
		Binding: binding,
		Sampler: sampler.handle,
	}
}
