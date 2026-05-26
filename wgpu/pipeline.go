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

// PipelineLayoutDescriptor describes a pipeline layout to create.
// BindGroupLayouts is a slice of *BindGroupLayout; nil for auto layout.
type PipelineLayoutDescriptor struct {
	Label            string
	BindGroupLayouts []*BindGroupLayout
}

// pipelineLayoutDescriptorWire is the FFI-compatible C-layout struct for wgpu-native.
// v29: Added ImmediateSize field for immediate data (push constants replacement).
// CRITICAL: layout must match WGPUPipelineLayoutDescriptor exactly.
// nextInChain(8)+label(16)+bindGroupLayoutCount(8)+bindGroupLayouts(8)+immediateSize(4)+pad(4) = 48 bytes.
type pipelineLayoutDescriptorWire struct {
	NextInChain          uintptr // *ChainedStruct
	Label                StringView
	BindGroupLayoutCount uintptr // size_t
	BindGroupLayouts     uintptr // *WGPUBindGroupLayout
	ImmediateSize        uint32  // NEW in v29: bytes of immediate data allocated for shaders (requires NativeFeatureImmediates)
	_pad                 [4]byte //nolint:unused // padding for FFI alignment
}

// NativeLimits extends WGPULimits with wgpu-native specific limits.
// Chain via NextInChain in Limits with SType = STypeNativeLimits.
// v29 BREAKING: maxPushConstantSize renamed to maxImmediateSize; maxNonSamplerBindings and
// maxBindingArrayElementsPerShaderStage added.
type NativeLimits struct {
	Chain                                 ChainedStruct // chain.SType must be STypeNativeLimits
	MaxImmediateSize                      uint32        // was maxPushConstantSize in v27
	MaxNonSamplerBindings                 uint32        // max live non-sampler bindings (DX12 only)
	MaxBindingArrayElementsPerShaderStage uint32        // max resources in binding arrays per shader stage
}

// PipelineLayoutExtras provides wgpu-native specific pipeline layout extensions.
// Chain via NextInChain in PipelineLayoutDescriptor with SType = STypePipelineLayoutExtras.
// v29 BREAKING: pushConstantRangeCount/pushConstantRanges replaced by immediateDataSize.
type PipelineLayoutExtras struct {
	Chain             ChainedStruct // chain.SType must be STypePipelineLayoutExtras
	ImmediateDataSize uint32        // bytes of immediate data for shaders (requires NativeFeatureImmediates)
}

// ComputePipelineDescriptor describes a compute pipeline to create.
// Layout is nil for auto layout.
type ComputePipelineDescriptor struct {
	Label      string
	Layout     *PipelineLayout // nil for auto layout
	Module     *ShaderModule
	EntryPoint string
}

// computePipelineDescriptorWire is the FFI-compatible C-layout struct for wgpu-native.
// CRITICAL: layout must match WGPUComputePipelineDescriptor exactly.
// nextInChain(8)+label(16)+layout(8)+compute(48) = 80 bytes.
type computePipelineDescriptorWire struct {
	NextInChain uintptr // *ChainedStruct
	Label       StringView
	Layout      uintptr // WGPUPipelineLayout (nullable)
	Compute     ProgrammableStageDescriptor
}

// CreatePipelineLayout creates a pipeline layout.
// Returns an error if the FFI call fails or the device/descriptor is nil.
func (d *Device) CreatePipelineLayout(desc *PipelineLayoutDescriptor) (*PipelineLayout, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if d == nil || d.handle == 0 {
		return nil, &WGPUError{Op: "CreatePipelineLayout", Message: "device is nil or released"}
	}
	if desc == nil {
		return nil, &WGPUError{Op: "CreatePipelineLayout", Message: "descriptor is nil"}
	}

	// Convert []*BindGroupLayout → []uintptr handles
	var layoutsPtr uintptr
	var handles []uintptr
	if len(desc.BindGroupLayouts) > 0 {
		handles = make([]uintptr, len(desc.BindGroupLayouts))
		for i, l := range desc.BindGroupLayouts {
			if l != nil {
				handles[i] = l.handle
			}
		}
		layoutsPtr = uintptr(unsafe.Pointer(&handles[0]))
	}

	wire := pipelineLayoutDescriptorWire{
		Label:                stringToStringView(desc.Label),
		BindGroupLayoutCount: uintptr(len(desc.BindGroupLayouts)),
		BindGroupLayouts:     layoutsPtr,
	}

	handle, _, _ := procDeviceCreatePipelineLayout.Call(
		d.handle,
		uintptr(unsafe.Pointer(&wire)),
	)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreatePipelineLayout", Message: "wgpu returned null handle"}
	}
	trackResource(handle, "PipelineLayout")
	return &PipelineLayout{handle: handle}, nil
}

// CreatePipelineLayoutSimple creates a pipeline layout with the given bind group layouts.
// Returns an error if the FFI call fails or the device is nil.
func (d *Device) CreatePipelineLayoutSimple(layouts []*BindGroupLayout) (*PipelineLayout, error) {
	return d.CreatePipelineLayout(&PipelineLayoutDescriptor{
		BindGroupLayouts: layouts,
	})
}

// Release releases the pipeline layout.
func (pl *PipelineLayout) Release() {
	if pl.handle != 0 {
		untrackResource(pl.handle)
		procPipelineLayoutRelease.Call(pl.handle) //nolint:errcheck
		pl.handle = 0
	}
}

// Handle returns the underlying handle.
func (pl *PipelineLayout) Handle() uintptr { return pl.handle }

// CreateComputePipeline creates a compute pipeline.
// Returns an error if the FFI call fails or the device/descriptor is nil.
func (d *Device) CreateComputePipeline(desc *ComputePipelineDescriptor) (*ComputePipeline, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if d == nil || d.handle == 0 {
		return nil, &WGPUError{Op: "CreateComputePipeline", Message: "device is nil or released"}
	}
	if desc == nil {
		return nil, &WGPUError{Op: "CreateComputePipeline", Message: "descriptor is nil"}
	}
	if desc.Module == nil {
		return nil, &WGPUError{Op: "CreateComputePipeline", Message: "shader module is nil"}
	}

	entryPointBytes := []byte(desc.EntryPoint)

	compute := ProgrammableStageDescriptor{
		Module: desc.Module.handle,
	}
	if len(entryPointBytes) > 0 {
		compute.EntryPoint = StringView{
			Data:   uintptr(unsafe.Pointer(&entryPointBytes[0])),
			Length: uintptr(len(entryPointBytes)),
		}
	} else {
		compute.EntryPoint = EmptyStringView()
	}

	var layoutHandle uintptr
	if desc.Layout != nil {
		layoutHandle = desc.Layout.handle
	}

	wire := computePipelineDescriptorWire{
		Label:   stringToStringView(desc.Label),
		Layout:  layoutHandle,
		Compute: compute,
	}

	handle, _, _ := procDeviceCreateComputePipeline.Call(
		d.handle,
		uintptr(unsafe.Pointer(&wire)),
	)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreateComputePipeline", Message: "wgpu returned null handle"}
	}
	trackResource(handle, "ComputePipeline")
	return &ComputePipeline{handle: handle}, nil
}

// CreateComputePipelineSimple creates a compute pipeline with the given shader and entry point.
// If layout is nil, auto layout is used.
// Returns an error if the FFI call fails or the device/shader is nil.
func (d *Device) CreateComputePipelineSimple(layout *PipelineLayout, shader *ShaderModule, entryPoint string) (*ComputePipeline, error) {
	return d.CreateComputePipeline(&ComputePipelineDescriptor{
		Layout:     layout,
		Module:     shader,
		EntryPoint: entryPoint,
	})
}

// GetBindGroupLayout returns the bind group layout at the given index.
// Useful for auto-layout pipelines.
func (cp *ComputePipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	mustInit()
	if cp == nil || cp.handle == 0 {
		return nil
	}
	handle, _, _ := procComputePipelineGetBindGroupLayout.Call(
		cp.handle,
		uintptr(groupIndex),
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "BindGroupLayout")
	return &BindGroupLayout{handle: handle}
}

// Release releases the compute pipeline.
func (cp *ComputePipeline) Release() {
	if cp.handle != 0 {
		untrackResource(cp.handle)
		procComputePipelineRelease.Call(cp.handle) //nolint:errcheck
		cp.handle = 0
	}
}

// Handle returns the underlying handle.
func (cp *ComputePipeline) Handle() uintptr { return cp.handle }
