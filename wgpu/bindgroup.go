package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// BufferBindingLayout describes buffer binding properties.
type BufferBindingLayout struct {
	NextInChain      uintptr // *ChainedStruct
	Type             gputypes.BufferBindingType
	HasDynamicOffset Bool
	MinBindingSize   uint64
}

// SamplerBindingLayout describes sampler binding properties.
type SamplerBindingLayout struct {
	NextInChain uintptr // *ChainedStruct
	Type        gputypes.SamplerBindingType
}

// TextureBindingLayout describes texture binding properties.
type TextureBindingLayout struct {
	NextInChain   uintptr // *ChainedStruct
	SampleType    gputypes.TextureSampleType
	ViewDimension gputypes.TextureViewDimension
	Multisampled  Bool
}

// StorageTextureBindingLayout describes storage texture binding properties.
type StorageTextureBindingLayout struct {
	NextInChain   uintptr // *ChainedStruct
	Access        gputypes.StorageTextureAccess
	Format        gputypes.TextureFormat
	ViewDimension gputypes.TextureViewDimension
}

// BindGroupLayoutEntry describes a single binding in a bind group layout.
type BindGroupLayoutEntry struct {
	NextInChain    uintptr // *ChainedStruct
	Binding        uint32
	Visibility     gputypes.ShaderStage
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

// =============================================================================
// Wire structs for FFI (with converted enum values and uint64 ShaderStage)
// wgpu-native uses uint64 for WGPUShaderStageFlags (via WGPUFlags typedef)
// =============================================================================

// bufferBindingLayoutWire is the FFI-compatible struct with wgpu-native enum values.
type bufferBindingLayoutWire struct {
	NextInChain      uintptr
	Type             uint32 // wgpu-native value (converted from gputypes)
	HasDynamicOffset Bool
	MinBindingSize   uint64
}

// samplerBindingLayoutWire is the FFI-compatible struct with wgpu-native enum values.
// Size: 16 bytes (8 + 4 + 4 padding) - must match C struct padding
type samplerBindingLayoutWire struct {
	NextInChain uintptr
	Type        uint32  // wgpu-native value
	_pad        [4]byte // padding to 8-byte alignment (C struct padding)
}

// textureBindingLayoutWire is the FFI-compatible struct with wgpu-native enum values.
// Size: 24 bytes (8 + 4 + 4 + 4 + 4 padding) - must match C struct padding
type textureBindingLayoutWire struct {
	NextInChain   uintptr
	SampleType    uint32  // wgpu-native value
	ViewDimension uint32  // wgpu-native value
	Multisampled  Bool    // 4 bytes
	_pad          [4]byte // padding to 8-byte alignment
}

// storageTextureBindingLayoutWire is the FFI-compatible struct with wgpu-native enum values.
// Size: 24 bytes (8 + 4 + 4 + 4 + 4 padding) - must match C struct padding
type storageTextureBindingLayoutWire struct {
	NextInChain   uintptr
	Access        uint32  // wgpu-native value
	Format        uint32  // wgpu-native value
	ViewDimension uint32  // wgpu-native value
	_pad          [4]byte // padding to 8-byte alignment
}

// bindGroupLayoutEntryWire is the FFI-compatible struct with converted enums.
// CRITICAL: Visibility is uint64 because wgpu-native defines WGPUShaderStageFlags as uint64!
type bindGroupLayoutEntryWire struct {
	NextInChain    uintptr
	Binding        uint32
	_pad           [4]byte // padding to align Visibility to 8 bytes
	Visibility     uint64  // WGPUShaderStageFlags = uint64 in wgpu-native!
	Buffer         bufferBindingLayoutWire
	Sampler        samplerBindingLayoutWire
	Texture        textureBindingLayoutWire
	StorageTexture storageTextureBindingLayoutWire
}

// toWire converts a BindGroupLayoutEntry to its wire representation.
func (e *BindGroupLayoutEntry) toWire() bindGroupLayoutEntryWire {
	return bindGroupLayoutEntryWire{
		NextInChain: e.NextInChain,
		Binding:     e.Binding,
		Visibility:  uint64(e.Visibility), // widen uint32 to uint64
		Buffer: bufferBindingLayoutWire{
			NextInChain:      e.Buffer.NextInChain,
			Type:             toWGPUBufferBindingType(e.Buffer.Type),
			HasDynamicOffset: e.Buffer.HasDynamicOffset,
			MinBindingSize:   e.Buffer.MinBindingSize,
		},
		Sampler: samplerBindingLayoutWire{
			NextInChain: e.Sampler.NextInChain,
			Type:        toWGPUSamplerBindingType(e.Sampler.Type),
		},
		Texture: textureBindingLayoutWire{
			NextInChain:   e.Texture.NextInChain,
			SampleType:    toWGPUTextureSampleType(e.Texture.SampleType),
			ViewDimension: toWGPUTextureViewDimension(e.Texture.ViewDimension),
			Multisampled:  e.Texture.Multisampled,
		},
		StorageTexture: storageTextureBindingLayoutWire{
			NextInChain:   e.StorageTexture.NextInChain,
			Access:        toWGPUStorageTextureAccess(e.StorageTexture.Access),
			Format:        toWGPUTextureFormat(e.StorageTexture.Format),
			ViewDimension: toWGPUTextureViewDimension(e.StorageTexture.ViewDimension),
		},
	}
}

// bindGroupLayoutDescriptorWire is the FFI-compatible descriptor.
type bindGroupLayoutDescriptorWire struct {
	NextInChain uintptr
	Label       StringView
	EntryCount  uintptr
	Entries     uintptr // *bindGroupLayoutEntryWire
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
// Entries are converted from gputypes to wgpu-native enum values before FFI call.
func (d *Device) CreateBindGroupLayout(desc *BindGroupLayoutDescriptor) *BindGroupLayout {
	mustInit()
	if desc == nil {
		return nil
	}

	// If there are entries, we need to convert them to wire format
	var wireDesc bindGroupLayoutDescriptorWire
	wireDesc.NextInChain = desc.NextInChain
	wireDesc.Label = desc.Label
	wireDesc.EntryCount = desc.EntryCount

	if desc.EntryCount > 0 && desc.Entries != 0 {
		// Convert entries to wire format
		entries := unsafe.Slice((*BindGroupLayoutEntry)(unsafe.Pointer(desc.Entries)), desc.EntryCount)
		wireEntries := make([]bindGroupLayoutEntryWire, len(entries))
		for i := range entries {
			wireEntries[i] = entries[i].toWire()
		}
		wireDesc.Entries = uintptr(unsafe.Pointer(&wireEntries[0]))
	}

	handle, _, _ := procDeviceCreateBindGroupLayout.Call(
		d.handle,
		uintptr(unsafe.Pointer(&wireDesc)),
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

	// Convert entries to wire format
	wireEntries := make([]bindGroupLayoutEntryWire, len(entries))
	for i := range entries {
		wireEntries[i] = entries[i].toWire()
	}

	wireDesc := bindGroupLayoutDescriptorWire{
		Label:      EmptyStringView(),
		EntryCount: uintptr(len(entries)),
		Entries:    uintptr(unsafe.Pointer(&wireEntries[0])),
	}

	handle, _, _ := procDeviceCreateBindGroupLayout.Call(
		d.handle,
		uintptr(unsafe.Pointer(&wireDesc)),
	)
	if handle == 0 {
		return nil
	}
	return &BindGroupLayout{handle: handle}
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
