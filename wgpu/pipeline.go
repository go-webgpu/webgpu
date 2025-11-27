package wgpu

import (
	"unsafe"
)

// ProgrammableStageDescriptor describes a programmable shader stage.
type ProgrammableStageDescriptor struct {
	NextInChain   uintptr    // *ChainedStruct
	Module        uintptr    // WGPUShaderModule
	EntryPoint    StringView // Entry point function name
	ConstantCount uintptr    // size_t
	Constants     uintptr    // *ConstantEntry
}

// PipelineLayoutDescriptor describes a pipeline layout.
type PipelineLayoutDescriptor struct {
	NextInChain          uintptr // *ChainedStruct
	Label                StringView
	BindGroupLayoutCount uintptr // size_t
	BindGroupLayouts     uintptr // *WGPUBindGroupLayout
}

// ComputePipelineDescriptor describes a compute pipeline.
type ComputePipelineDescriptor struct {
	NextInChain uintptr // *ChainedStruct
	Label       StringView
	Layout      uintptr // WGPUPipelineLayout (nullable)
	Compute     ProgrammableStageDescriptor
}

// CreatePipelineLayout creates a pipeline layout.
func (d *Device) CreatePipelineLayout(desc *PipelineLayoutDescriptor) *PipelineLayout {
	mustInit()
	if desc == nil {
		return nil
	}
	handle, _, _ := procDeviceCreatePipelineLayout.Call(
		d.handle,
		uintptr(unsafe.Pointer(desc)),
	)
	if handle == 0 {
		return nil
	}
	return &PipelineLayout{handle: handle}
}

// CreatePipelineLayoutSimple creates a pipeline layout with the given bind group layouts.
func (d *Device) CreatePipelineLayoutSimple(layouts []*BindGroupLayout) *PipelineLayout {
	mustInit()
	if len(layouts) == 0 {
		// Create empty pipeline layout
		desc := PipelineLayoutDescriptor{
			Label:                EmptyStringView(),
			BindGroupLayoutCount: 0,
			BindGroupLayouts:     0,
		}
		return d.CreatePipelineLayout(&desc)
	}
	// Convert to handles
	handles := make([]uintptr, len(layouts))
	for i, l := range layouts {
		handles[i] = l.handle
	}
	desc := PipelineLayoutDescriptor{
		Label:                EmptyStringView(),
		BindGroupLayoutCount: uintptr(len(handles)),
		BindGroupLayouts:     uintptr(unsafe.Pointer(&handles[0])),
	}
	return d.CreatePipelineLayout(&desc)
}

// Release releases the pipeline layout.
func (pl *PipelineLayout) Release() {
	if pl.handle != 0 {
		procPipelineLayoutRelease.Call(pl.handle) //nolint:errcheck
		pl.handle = 0
	}
}

// Handle returns the underlying handle.
func (pl *PipelineLayout) Handle() uintptr { return pl.handle }

// CreateComputePipeline creates a compute pipeline.
func (d *Device) CreateComputePipeline(desc *ComputePipelineDescriptor) *ComputePipeline {
	mustInit()
	if desc == nil {
		return nil
	}
	handle, _, _ := procDeviceCreateComputePipeline.Call(
		d.handle,
		uintptr(unsafe.Pointer(desc)),
	)
	if handle == 0 {
		return nil
	}
	return &ComputePipeline{handle: handle}
}

// CreateComputePipelineSimple creates a compute pipeline with the given shader and entry point.
// If layout is nil, auto layout is used.
func (d *Device) CreateComputePipelineSimple(layout *PipelineLayout, shader *ShaderModule, entryPoint string) *ComputePipeline {
	mustInit()
	if shader == nil {
		return nil
	}
	entryBytes := []byte(entryPoint)
	desc := ComputePipelineDescriptor{
		Label: EmptyStringView(),
		Compute: ProgrammableStageDescriptor{
			Module: shader.handle,
			EntryPoint: StringView{
				Data:   uintptr(unsafe.Pointer(&entryBytes[0])),
				Length: uintptr(len(entryBytes)),
			},
		},
	}
	if layout != nil {
		desc.Layout = layout.handle
	}
	return d.CreateComputePipeline(&desc)
}

// GetBindGroupLayout returns the bind group layout at the given index.
// Useful for auto-layout pipelines.
func (cp *ComputePipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	mustInit()
	handle, _, _ := procComputePipelineGetBindGroupLayout.Call(
		cp.handle,
		uintptr(groupIndex),
	)
	if handle == 0 {
		return nil
	}
	return &BindGroupLayout{handle: handle}
}

// Release releases the compute pipeline.
func (cp *ComputePipeline) Release() {
	if cp.handle != 0 {
		procComputePipelineRelease.Call(cp.handle) //nolint:errcheck
		cp.handle = 0
	}
}

// Handle returns the underlying handle.
func (cp *ComputePipeline) Handle() uintptr { return cp.handle }
