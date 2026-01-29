// Package wgpu provides Zero-CGO WebGPU bindings for Go.
//
// This package wraps wgpu-native (Rust WebGPU implementation) using pure Go FFI
// via syscall on Windows and dlopen on Unix. No CGO is required.
//
// # Quick Start
//
//	// Initialize the library
//	if err := wgpu.Init(); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Create instance
//	instance, err := wgpu.CreateInstance(nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer instance.Release()
//
//	// Request adapter (GPU)
//	adapter, err := instance.RequestAdapter(nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer adapter.Release()
//
//	// Request device
//	device, err := adapter.RequestDevice(nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer device.Release()
//
// # Core Objects
//
// The WebGPU API is structured around several core object types:
//
//   - [Instance]: Entry point to the WebGPU API
//   - [Adapter]: Represents a physical GPU
//   - [Device]: Logical device for creating resources
//   - [Queue]: Command submission queue
//   - [Buffer]: GPU memory buffer
//   - [Texture]: GPU texture resource
//   - [Sampler]: Texture sampling configuration
//   - [ShaderModule]: Compiled shader code (WGSL)
//   - [BindGroup]: Resource bindings for shaders
//   - [RenderPipeline]: Configuration for rendering
//   - [ComputePipeline]: Configuration for compute
//   - [CommandEncoder]: Records GPU commands
//   - [RenderPassEncoder]: Records render commands
//   - [ComputePassEncoder]: Records compute commands
//
// # Resource Management
//
// All WebGPU objects must be released when no longer needed:
//
//	buffer := device.CreateBuffer(&wgpu.BufferDescriptor{...})
//	defer buffer.Release()
//
// # Render Pipeline
//
// A typical render pipeline setup:
//
//	// Create shader module
//	shader := device.CreateShaderModuleWGSL(shaderCode)
//	defer shader.Release()
//
//	// Create render pipeline
//	pipeline := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
//	    Vertex: wgpu.VertexState{
//	        Module:     vsModule,
//	        EntryPoint: "main",
//	        Buffers:    []wgpu.VertexBufferLayout{vertexBufferLayout},
//	    },
//	    Fragment: &wgpu.FragmentState{
//	        Module:     fsModule,
//	        EntryPoint: "main",
//	        Targets:    []wgpu.ColorTargetState{{Format: format, WriteMask: gputypes.ColorWriteMaskAll}},
//	    },
//	    // ... other configuration
//	})
//	defer pipeline.Release()
//
// # Compute Pipeline
//
// For GPU compute operations:
//
//	pipeline := device.CreateComputePipelineSimple(nil, shader, "main")
//	defer pipeline.Release()
//
//	// Dispatch compute work
//	computePass := encoder.BeginComputePass(nil)
//	computePass.SetPipeline(pipeline)
//	computePass.SetBindGroup(0, bindGroup, nil)
//	computePass.DispatchWorkgroups(workgroupCount, 1, 1)
//	computePass.End()
//
// # Indirect Drawing
//
// GPU-driven rendering using indirect buffers:
//
//	// Create indirect buffer with draw args
//	args := wgpu.DrawIndirectArgs{
//	    VertexCount:   3,
//	    InstanceCount: 100,
//	}
//	// Write to buffer...
//
//	// Draw using GPU-specified parameters
//	renderPass.DrawIndirect(indirectBuffer, 0)
//
// # RenderBundle
//
// Pre-record render commands for efficient replay:
//
//	bundleEncoder := device.CreateRenderBundleEncoderSimple(
//	    []gputypes.TextureFormat{surfaceFormat},
//	    gputypes.TextureFormatUndefined,
//	    1,
//	)
//	bundleEncoder.SetPipeline(pipeline)
//	bundleEncoder.Draw(3, 1, 0, 0)
//	bundle := bundleEncoder.Finish(nil)
//	defer bundle.Release()
//
//	// Later, in a render pass:
//	renderPass.ExecuteBundles([]*wgpu.RenderBundle{bundle})
//
// # Platform Support
//
// Supported platforms:
//   - Windows (x64) - uses syscall.LazyDLL
//   - Linux (x64, arm64) - uses goffi/dlopen
//   - macOS (x64, arm64) - uses goffi/dlopen
//
// # Dependencies
//
// This package requires wgpu-native library:
//   - Windows: wgpu_native.dll
//   - Linux: libwgpu_native.so
//   - macOS: libwgpu_native.dylib
//
// Download from: https://github.com/gfx-rs/wgpu-native/releases
package wgpu
