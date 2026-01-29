// Example: Instanced Rendering
// Demonstrates drawing many objects efficiently using GPU instancing.
// Instead of issuing separate draw calls, we draw all instances in one call.
package main

import (
	"fmt"
	"log"
	"math"
	"unsafe"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
)

// Vertex data: position (x, y) + color (r, g, b)
// Size: 20 bytes per vertex (2*4 + 3*4)
type Vertex struct {
	Position [2]float32
	Color    [3]float32
}

// Instance data: offset (x, y) + scale + padding
// Size: 16 bytes per instance (must be aligned for GPU)
type InstanceData struct {
	Offset [2]float32
	Scale  float32
	_pad   float32 // Padding to align to 16 bytes
}

// Triangle vertices (centered at origin)
var triangleVertices = []Vertex{
	{Position: [2]float32{0.0, 0.1}, Color: [3]float32{1.0, 0.0, 0.0}},   // Top - red
	{Position: [2]float32{-0.1, -0.1}, Color: [3]float32{0.0, 1.0, 0.0}}, // Bottom-left - green
	{Position: [2]float32{0.1, -0.1}, Color: [3]float32{0.0, 0.0, 1.0}},  // Bottom-right - blue
}

// Generate instance data for a grid of triangles
func generateInstanceData(gridSize int) []InstanceData {
	instances := make([]InstanceData, gridSize*gridSize)
	spacing := 2.0 / float32(gridSize+1)

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			idx := y*gridSize + x
			// Position instances in a grid from -1 to 1
			instances[idx] = InstanceData{
				Offset: [2]float32{
					-1.0 + spacing*float32(x+1),
					-1.0 + spacing*float32(y+1),
				},
				Scale: 0.5 + 0.5*float32(math.Sin(float64(idx)*0.5)), // Varying scale
			}
		}
	}
	return instances
}

const shaderCode = `
struct VertexInput {
    @location(0) position: vec2<f32>,
    @location(1) color: vec3<f32>,
    // Per-instance attributes
    @location(2) offset: vec2<f32>,
    @location(3) scale: f32,
}

struct VertexOutput {
    @builtin(position) position: vec4<f32>,
    @location(0) color: vec3<f32>,
}

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    // Apply scale and offset to position
    let scaled_pos = in.position * in.scale;
    let final_pos = scaled_pos + in.offset;
    out.position = vec4<f32>(final_pos, 0.0, 1.0);
    out.color = in.color;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    return vec4<f32>(in.color, 1.0);
}
`

func main() {
	// Initialize
	if err := wgpu.Init(); err != nil {
		log.Fatal(err)
	}

	instance, err := wgpu.CreateInstance(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer device.Release()

	queue := device.GetQueue()
	defer queue.Release()

	// Create shader module
	shader := device.CreateShaderModuleWGSL(shaderCode)
	if shader == nil {
		log.Fatal("failed to create shader module")
	}
	defer shader.Release()

	// Generate instance data (5x5 grid = 25 instances)
	const gridSize = 5
	instanceData := generateInstanceData(gridSize)
	instanceCount := uint32(len(instanceData))

	fmt.Printf("Drawing %d triangles using instancing\n", instanceCount)

	// Create vertex buffer (per-vertex data)
	vertexBufferSize := uint64(len(triangleVertices)) * uint64(unsafe.Sizeof(triangleVertices[0]))
	vertexBuffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage:            gputypes.BufferUsageVertex | gputypes.BufferUsageCopyDst,
		Size:             vertexBufferSize,
		MappedAtCreation: wgpu.True,
	})
	if vertexBuffer == nil {
		log.Fatal("failed to create vertex buffer")
	}
	defer vertexBuffer.Release()

	// Copy vertex data
	ptr := vertexBuffer.GetMappedRange(0, vertexBufferSize)
	if ptr != nil {
		mappedSlice := unsafe.Slice((*Vertex)(ptr), len(triangleVertices))
		copy(mappedSlice, triangleVertices)
	}
	vertexBuffer.Unmap()

	// Create instance buffer (per-instance data)
	instanceBufferSize := uint64(len(instanceData)) * uint64(unsafe.Sizeof(instanceData[0]))
	instanceBuffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage:            gputypes.BufferUsageVertex | gputypes.BufferUsageCopyDst,
		Size:             instanceBufferSize,
		MappedAtCreation: wgpu.True,
	})
	if instanceBuffer == nil {
		log.Fatal("failed to create instance buffer")
	}
	defer instanceBuffer.Release()

	// Copy instance data
	ptr = instanceBuffer.GetMappedRange(0, instanceBufferSize)
	if ptr != nil {
		mappedSlice := unsafe.Slice((*InstanceData)(ptr), len(instanceData))
		copy(mappedSlice, instanceData)
	}
	instanceBuffer.Unmap()

	// Define vertex attributes for per-vertex buffer
	vertexAttributes := []wgpu.VertexAttribute{
		{Format: gputypes.VertexFormatFloat32x2, Offset: 0, ShaderLocation: 0}, // position
		{Format: gputypes.VertexFormatFloat32x3, Offset: 8, ShaderLocation: 1}, // color
	}

	// Define vertex attributes for per-instance buffer
	instanceAttributes := []wgpu.VertexAttribute{
		{Format: gputypes.VertexFormatFloat32x2, Offset: 0, ShaderLocation: 2}, // offset
		{Format: gputypes.VertexFormatFloat32, Offset: 8, ShaderLocation: 3},   // scale
	}

	// Create render pipeline with two vertex buffer layouts:
	// - Slot 0: Per-vertex data (position, color) with VertexStepModeVertex
	// - Slot 1: Per-instance data (offset, scale) with VertexStepModeInstance
	pipeline := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers: []wgpu.VertexBufferLayout{
				// Per-vertex buffer (slot 0)
				{
					ArrayStride:    uint64(unsafe.Sizeof(Vertex{})),
					StepMode:       gputypes.VertexStepModeVertex,
					AttributeCount: uintptr(len(vertexAttributes)),
					Attributes:     &vertexAttributes[0],
				},
				// Per-instance buffer (slot 1)
				{
					ArrayStride:    uint64(unsafe.Sizeof(InstanceData{})),
					StepMode:       gputypes.VertexStepModeInstance, // Key: advances per instance, not per vertex
					AttributeCount: uintptr(len(instanceAttributes)),
					Attributes:     &instanceAttributes[0],
				},
			},
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{
				{Format: gputypes.TextureFormatRGBA8Unorm, WriteMask: gputypes.ColorWriteMaskAll},
			},
		},
		Primitive: wgpu.PrimitiveState{
			Topology:  gputypes.PrimitiveTopologyTriangleList,
			FrontFace: gputypes.FrontFaceCCW,
			CullMode:  gputypes.CullModeNone,
		},
		Multisample: wgpu.MultisampleState{
			Count: 1,
			Mask:  0xFFFFFFFF,
		},
	})
	if pipeline == nil {
		log.Fatal("failed to create render pipeline")
	}
	defer pipeline.Release()

	// Create output texture
	outputTexture := device.CreateTexture(&wgpu.TextureDescriptor{
		Size:   gputypes.Extent3D{Width: 256, Height: 256, DepthOrArrayLayers: 1},
		Format: gputypes.TextureFormatRGBA8Unorm,
		Usage:  gputypes.TextureUsageRenderAttachment | gputypes.TextureUsageCopySrc,
	})
	if outputTexture == nil {
		log.Fatal("failed to create output texture")
	}
	defer outputTexture.Release()

	outputView := outputTexture.CreateView(nil)
	defer outputView.Release()

	// Create command encoder
	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		log.Fatal("failed to create command encoder")
	}

	// Begin render pass
	renderPass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{
			{
				View:       outputView,
				LoadOp:     gputypes.LoadOpClear,
				StoreOp:    gputypes.StoreOpStore,
				ClearValue: wgpu.Color{R: 0.1, G: 0.1, B: 0.1, A: 1.0},
			},
		},
	})

	renderPass.SetPipeline(pipeline)
	renderPass.SetVertexBuffer(0, vertexBuffer, 0, vertexBufferSize)     // Per-vertex data
	renderPass.SetVertexBuffer(1, instanceBuffer, 0, instanceBufferSize) // Per-instance data

	// Draw all instances in ONE call!
	// vertexCount=3 (triangle), instanceCount=25 (grid)
	renderPass.Draw(3, instanceCount, 0, 0)

	renderPass.End()
	renderPass.Release()

	// Submit
	cmdBuffer := encoder.Finish(nil)
	encoder.Release()

	queue.Submit(cmdBuffer)
	cmdBuffer.Release()

	fmt.Println("Instanced rendering complete!")
	fmt.Println()
	fmt.Println("Key concepts demonstrated:")
	fmt.Println("  - VertexStepModeVertex: data advances per vertex")
	fmt.Println("  - VertexStepModeInstance: data advances per instance")
	fmt.Printf("  - Draw(vertexCount=3, instanceCount=%d): draws %d triangles in one call\n", instanceCount, instanceCount)
	fmt.Println()
	fmt.Printf("Without instancing: %d draw calls\n", instanceCount)
	fmt.Println("With instancing: 1 draw call")
}
