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
