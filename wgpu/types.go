package wgpu

// Handle types - opaque pointers to WebGPU objects

// Instance represents a WebGPU instance.
type Instance struct{ handle uintptr }

// Adapter represents a WebGPU adapter (GPU).
type Adapter struct{ handle uintptr }

// Device represents a WebGPU device.
type Device struct{ handle uintptr }

// Queue represents a WebGPU command queue.
type Queue struct{ handle uintptr }

// Buffer represents a WebGPU buffer.
type Buffer struct{ handle uintptr }

// Texture represents a WebGPU texture.
type Texture struct{ handle uintptr }

// TextureView represents a view into a WebGPU texture.
type TextureView struct{ handle uintptr }

// Sampler represents a WebGPU sampler.
type Sampler struct{ handle uintptr }

// ShaderModule represents a compiled shader.
type ShaderModule struct{ handle uintptr }

// BindGroupLayout represents a bind group layout.
type BindGroupLayout struct{ handle uintptr }

// BindGroup represents a bind group.
type BindGroup struct{ handle uintptr }

// PipelineLayout represents a pipeline layout.
type PipelineLayout struct{ handle uintptr }

// RenderPipeline represents a render pipeline.
type RenderPipeline struct{ handle uintptr }

// ComputePipeline represents a compute pipeline.
type ComputePipeline struct{ handle uintptr }

// CommandEncoder represents a command encoder.
type CommandEncoder struct{ handle uintptr }

// CommandBuffer represents an encoded command buffer.
type CommandBuffer struct{ handle uintptr }

// RenderPassEncoder represents a render pass encoder.
type RenderPassEncoder struct{ handle uintptr }

// ComputePassEncoder represents a compute pass encoder.
type ComputePassEncoder struct{ handle uintptr }

// Surface represents a surface for presenting rendered images.
type Surface struct{ handle uintptr }

// QuerySet represents a query set.
type QuerySet struct{ handle uintptr }

// RenderBundle represents a pre-recorded bundle of render commands.
type RenderBundle struct{ handle uintptr }

// RenderBundleEncoder represents a render bundle encoder.
type RenderBundleEncoder struct{ handle uintptr }

// DrawIndirectArgs contains arguments for indirect (GPU-driven) draw calls.
// This struct must be written to a Buffer for use with DrawIndirect.
// Size: 16 bytes, must be aligned to 4 bytes.
type DrawIndirectArgs struct {
	VertexCount   uint32 // Number of vertices to draw
	InstanceCount uint32 // Number of instances to draw
	FirstVertex   uint32 // First vertex index
	FirstInstance uint32 // First instance index
}

// DrawIndexedIndirectArgs contains arguments for indirect indexed draw calls.
// This struct must be written to a Buffer for use with DrawIndexedIndirect.
// Size: 20 bytes, must be aligned to 4 bytes.
type DrawIndexedIndirectArgs struct {
	IndexCount    uint32 // Number of indices to draw
	InstanceCount uint32 // Number of instances to draw
	FirstIndex    uint32 // First index in the index buffer
	BaseVertex    int32  // Value added to vertex index before indexing into vertex buffer
	FirstInstance uint32 // First instance index
}

// DispatchIndirectArgs contains arguments for indirect compute dispatch.
// This struct must be written to a Buffer for use with DispatchWorkgroupsIndirect.
// Size: 12 bytes, must be aligned to 4 bytes.
type DispatchIndirectArgs struct {
	WorkgroupCountX uint32 // Number of workgroups in X dimension
	WorkgroupCountY uint32 // Number of workgroups in Y dimension
	WorkgroupCountZ uint32 // Number of workgroups in Z dimension
}

// Handle returns the underlying handle. For advanced use only.
func (i *Instance) Handle() uintptr { return i.handle }

// Handle returns the underlying handle. For advanced use only.
func (a *Adapter) Handle() uintptr { return a.handle }

// Handle returns the underlying handle. For advanced use only.
func (d *Device) Handle() uintptr { return d.handle }

// Handle returns the underlying handle. For advanced use only.
func (q *Queue) Handle() uintptr { return q.handle }
