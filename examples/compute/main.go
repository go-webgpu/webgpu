// Example: Compute Shader
// Demonstrates GPU parallel processing using compute shaders.
// This example doubles all values in an array using the GPU.
package main

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/go-webgpu/webgpu/wgpu"
)

// Compute shader that doubles each element in the array
const computeShader = `
@group(0) @binding(0) var<storage, read_write> data: array<f32>;

@compute @workgroup_size(64)
fn main(@builtin(global_invocation_id) global_id: vec3<u32>) {
    let index = global_id.x;
    if (index < arrayLength(&data)) {
        data[index] = data[index] * 2.0;
    }
}
`

func main() {
	// Initialize WebGPU
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
	shader := device.CreateShaderModuleWGSL(computeShader)
	if shader == nil {
		log.Fatal("failed to create compute shader")
	}
	defer shader.Release()

	// Create compute pipeline with auto layout
	pipeline := device.CreateComputePipelineSimple(nil, shader, "main")
	if pipeline == nil {
		log.Fatal("failed to create compute pipeline")
	}
	defer pipeline.Release()

	// Prepare input data
	const numElements = 256
	inputData := make([]float32, numElements)
	for i := range inputData {
		inputData[i] = float32(i + 1) // 1, 2, 3, ..., 256
	}

	fmt.Println("=== Compute Shader Example ===")
	fmt.Printf("Processing %d elements on GPU\n", numElements)
	fmt.Printf("Input (first 10): %v\n", inputData[:10])

	// Create storage buffer
	bufferSize := uint64(numElements * 4) // 4 bytes per float32
	storageBuffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage:            wgpu.BufferUsageStorage | wgpu.BufferUsageCopySrc | wgpu.BufferUsageCopyDst,
		Size:             bufferSize,
		MappedAtCreation: wgpu.True,
	})
	if storageBuffer == nil {
		log.Fatal("failed to create storage buffer")
	}
	defer storageBuffer.Release()

	// Copy input data to buffer
	ptr := storageBuffer.GetMappedRange(0, bufferSize)
	if ptr != nil {
		mappedSlice := unsafe.Slice((*float32)(ptr), numElements)
		copy(mappedSlice, inputData)
	}
	storageBuffer.Unmap()

	// Create readback buffer for results
	readbackBuffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage:            wgpu.BufferUsageMapRead | wgpu.BufferUsageCopyDst,
		Size:             bufferSize,
		MappedAtCreation: wgpu.False,
	})
	if readbackBuffer == nil {
		log.Fatal("failed to create readback buffer")
	}
	defer readbackBuffer.Release()

	// Get bind group layout from pipeline
	bindGroupLayout := pipeline.GetBindGroupLayout(0)
	if bindGroupLayout == nil {
		log.Fatal("failed to get bind group layout")
	}
	defer bindGroupLayout.Release()

	// Create bind group
	bindGroup := device.CreateBindGroupSimple(bindGroupLayout, []wgpu.BindGroupEntry{
		wgpu.BufferBindingEntry(0, storageBuffer, 0, bufferSize),
	})
	if bindGroup == nil {
		log.Fatal("failed to create bind group")
	}
	defer bindGroup.Release()

	// Create command encoder
	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		log.Fatal("failed to create command encoder")
	}

	// Begin compute pass
	computePass := encoder.BeginComputePass(nil)
	computePass.SetPipeline(pipeline)
	computePass.SetBindGroup(0, bindGroup, nil)

	// Dispatch workgroups
	// With workgroup_size(64), we need ceil(256/64) = 4 workgroups
	workgroupCount := uint32((numElements + 63) / 64)
	computePass.DispatchWorkgroups(workgroupCount, 1, 1)
	computePass.End()
	computePass.Release()

	// Copy results to readback buffer
	encoder.CopyBufferToBuffer(storageBuffer, 0, readbackBuffer, 0, bufferSize)

	// Submit commands
	cmdBuffer := encoder.Finish(nil)
	encoder.Release()
	queue.Submit(cmdBuffer)
	cmdBuffer.Release()

	// Map readback buffer and read results
	err = readbackBuffer.MapAsync(device, wgpu.MapModeRead, 0, bufferSize)
	if err != nil {
		log.Fatalf("MapAsync failed: %v", err)
	}

	resultPtr := readbackBuffer.GetMappedRange(0, bufferSize)
	if resultPtr != nil {
		results := unsafe.Slice((*float32)(resultPtr), numElements)
		fmt.Printf("Output (first 10): %v\n", results[:10])

		// Verify results
		correct := true
		for i := 0; i < numElements; i++ {
			expected := float32((i + 1) * 2)
			if results[i] != expected {
				fmt.Printf("Mismatch at %d: expected %f, got %f\n", i, expected, results[i])
				correct = false
				break
			}
		}
		if correct {
			fmt.Println("All results correct!")
		}
	}
	readbackBuffer.Unmap()

	fmt.Println()
	fmt.Println("Key concepts demonstrated:")
	fmt.Println("  - Storage buffer with read_write access")
	fmt.Println("  - Compute shader with @workgroup_size(64)")
	fmt.Println("  - DispatchWorkgroups for parallel execution")
	fmt.Println("  - Buffer mapping for CPU/GPU data transfer")
}
