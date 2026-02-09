package wgpu

// Handle types — opaque wrappers around wgpu-native object pointers.
// Each type must be explicitly released via its Release method when no longer needed.

// Instance is the entry point to the WebGPU API.
// Create with [CreateInstance], release with [Instance.Release].
type Instance struct{ handle uintptr }

// Adapter represents a physical GPU and its capabilities.
// Obtained via [Instance.RequestAdapter], release with [Adapter.Release].
type Adapter struct{ handle uintptr }

// Device is the logical connection to a GPU, used to create all other resources.
// Obtained via [Adapter.RequestDevice], release with [Device.Release].
type Device struct{ handle uintptr }

// Queue is used to submit command buffers and write data to buffers/textures.
// Obtained via [Device.GetQueue], release with [Queue.Release].
type Queue struct{ handle uintptr }

// Buffer represents a block of GPU-accessible memory.
// Create with [Device.CreateBuffer], release with [Buffer.Release].
type Buffer struct{ handle uintptr }

// Texture represents a GPU texture resource (1D, 2D, or 3D).
// Create with [Device.CreateTexture], release with [Texture.Release].
type Texture struct{ handle uintptr }

// TextureView is a view into a subset of a [Texture], used in bind groups and render passes.
// Create with [Texture.CreateView], release with [TextureView.Release].
type TextureView struct{ handle uintptr }

// Sampler defines how a shader samples a [Texture].
// Create with [Device.CreateSampler], release with [Sampler.Release].
type Sampler struct{ handle uintptr }

// ShaderModule holds compiled shader code (WGSL or SPIR-V).
// Create with [Device.CreateShaderModuleWGSL], release with [ShaderModule.Release].
type ShaderModule struct{ handle uintptr }

// BindGroupLayout defines the layout of resource bindings for a shader stage.
// Create with [Device.CreateBindGroupLayout], release with [BindGroupLayout.Release].
type BindGroupLayout struct{ handle uintptr }

// BindGroup binds actual GPU resources (buffers, textures, samplers) to shader slots.
// Create with [Device.CreateBindGroup], release with [BindGroup.Release].
type BindGroup struct{ handle uintptr }

// PipelineLayout defines the bind group layouts used by a pipeline.
// Create with [Device.CreatePipelineLayout], release with [PipelineLayout.Release].
type PipelineLayout struct{ handle uintptr }

// RenderPipeline is a compiled render pipeline configuration (shaders, vertex layout, blend state).
// Create with [Device.CreateRenderPipeline], release with [RenderPipeline.Release].
type RenderPipeline struct{ handle uintptr }

// ComputePipeline is a compiled compute pipeline configuration.
// Create with [Device.CreateComputePipeline], release with [ComputePipeline.Release].
type ComputePipeline struct{ handle uintptr }

// CommandEncoder records GPU commands into a [CommandBuffer].
// Create with [Device.CreateCommandEncoder], finalize with [CommandEncoder.Finish].
type CommandEncoder struct{ handle uintptr }

// CommandBuffer holds encoded GPU commands ready for submission via [Queue.Submit].
// Obtained from [CommandEncoder.Finish], release with [CommandBuffer.Release].
type CommandBuffer struct{ handle uintptr }

// RenderPassEncoder records draw commands within a render pass.
// Begin with [CommandEncoder.BeginRenderPass], end with [RenderPassEncoder.End].
type RenderPassEncoder struct{ handle uintptr }

// ComputePassEncoder records dispatch commands within a compute pass.
// Begin with [CommandEncoder.BeginComputePass], end with [ComputePassEncoder.End].
type ComputePassEncoder struct{ handle uintptr }

// Surface represents a platform window surface for presenting rendered frames.
// Create with platform-specific CreateSurface, release with [Surface.Release].
type Surface struct{ handle uintptr }

// QuerySet holds a set of GPU queries (occlusion or timestamp).
// Create with [Device.CreateQuerySet], release with [QuerySet.Release].
type QuerySet struct{ handle uintptr }

// RenderBundle is a pre-recorded set of render commands for efficient replay.
// Obtained from [RenderBundleEncoder.Finish], release with [RenderBundle.Release].
type RenderBundle struct{ handle uintptr }

// RenderBundleEncoder records render commands into a [RenderBundle].
// Create with [Device.CreateRenderBundleEncoder], finalize with [RenderBundleEncoder.Finish].
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
