// Package main demonstrates GPU timestamp queries for profiling.
// NOTE: Timestamp queries require the TIMESTAMP_QUERY feature which may not
// be enabled by default. This example shows how the API works.
package main

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
)

func main() {
	fmt.Println("=== Timestamp Query Example ===")
	fmt.Println()
	fmt.Println("GPU Timestamp Queries enable precise GPU profiling.")
	fmt.Println("They measure execution time in nanoseconds.")
	fmt.Println()
	fmt.Println("NOTE: This example requires TIMESTAMP_QUERY feature.")
	fmt.Println("      If not available, it demonstrates CPU timing instead.")
	fmt.Println()

	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	// Initialize WebGPU
	if err := wgpu.Init(); err != nil {
		return fmt.Errorf("init wgpu: %w", err)
	}

	inst, err := wgpu.CreateInstance(nil)
	if err != nil {
		return fmt.Errorf("create instance: %w", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		return fmt.Errorf("request adapter: %w", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		return fmt.Errorf("request device: %w", err)
	}
	defer device.Release()

	queue := device.GetQueue()
	defer queue.Release()

	// Try to create a timestamp query set
	// This will fail if TIMESTAMP_QUERY feature is not enabled
	querySet := device.CreateQuerySet(&wgpu.QuerySetDescriptor{
		Type:  wgpu.QueryTypeTimestamp,
		Count: 2,
	})

	if querySet != nil {
		defer querySet.Release()
		fmt.Println("TIMESTAMP_QUERY feature is available!")
		return runWithTimestamps(device, queue, querySet)
	}

	fmt.Println("TIMESTAMP_QUERY feature not available.")
	fmt.Println("Demonstrating with CPU timing instead.")
	fmt.Println()
	return runWithCPUTiming(device, queue)
}

// runWithTimestamps demonstrates actual GPU timestamp queries
func runWithTimestamps(device *wgpu.Device, queue *wgpu.Queue, querySet *wgpu.QuerySet) error {
	fmt.Println()
	fmt.Println("Using GPU timestamp queries for accurate profiling...")
	fmt.Println()

	// Create buffer to resolve query results (2 timestamps * 8 bytes each)
	queryResultSize := uint64(16)
	queryResultBuffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage: gputypes.BufferUsageQueryResolve | gputypes.BufferUsageCopySrc,
		Size:  queryResultSize,
	})
	if queryResultBuffer == nil {
		return fmt.Errorf("failed to create query result buffer")
	}
	defer queryResultBuffer.Release()

	// Create staging buffer for CPU read
	stagingBuffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage: gputypes.BufferUsageMapRead | gputypes.BufferUsageCopyDst,
		Size:  queryResultSize,
	})
	if stagingBuffer == nil {
		return fmt.Errorf("failed to create staging buffer")
	}
	defer stagingBuffer.Release()

	// Create compute pipeline for workload
	shaderCode := `
@group(0) @binding(0) var<storage, read_write> data: array<f32>;

@compute @workgroup_size(256)
fn main(@builtin(global_invocation_id) global_id: vec3<u32>) {
    let idx = global_id.x;
    if (idx < arrayLength(&data)) {
        var sum: f32 = data[idx];
        for (var i: u32 = 0u; i < 100u; i = i + 1u) {
            sum = sum * 1.01 + 0.01;
        }
        data[idx] = sum;
    }
}
`
	shader := device.CreateShaderModuleWGSL(shaderCode)
	if shader == nil {
		return fmt.Errorf("failed to create shader")
	}
	defer shader.Release()

	pipeline := device.CreateComputePipelineSimple(nil, shader, "main")
	if pipeline == nil {
		return fmt.Errorf("failed to create pipeline")
	}
	defer pipeline.Release()

	// Create data buffer
	const numElements = 1024 * 1024
	bufferSize := uint64(numElements * 4)

	dataBuffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage: gputypes.BufferUsageStorage,
		Size:  bufferSize,
	})
	if dataBuffer == nil {
		return fmt.Errorf("failed to create data buffer")
	}
	defer dataBuffer.Release()

	// Create bind group
	bindGroupLayout := pipeline.GetBindGroupLayout(0)
	defer bindGroupLayout.Release()

	bindGroup := device.CreateBindGroupSimple(bindGroupLayout, []wgpu.BindGroupEntry{
		wgpu.BufferBindingEntry(0, dataBuffer, 0, bufferSize),
	})
	if bindGroup == nil {
		return fmt.Errorf("failed to create bind group")
	}
	defer bindGroup.Release()

	// Record commands with timestamps
	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		return fmt.Errorf("failed to create command encoder")
	}

	// Write start timestamp
	encoder.WriteTimestamp(querySet, 0)

	// Execute compute pass
	pass := encoder.BeginComputePass(nil)
	pass.SetPipeline(pipeline)
	pass.SetBindGroup(0, bindGroup, nil)
	pass.DispatchWorkgroups(numElements/256, 1, 1)
	pass.End()
	pass.Release()

	// Write end timestamp
	encoder.WriteTimestamp(querySet, 1)

	// Resolve query results to buffer
	encoder.ResolveQuerySet(querySet, 0, 2, queryResultBuffer, 0)

	// Copy to staging buffer
	encoder.CopyBufferToBuffer(queryResultBuffer, 0, stagingBuffer, 0, queryResultSize)

	cmdBuffer := encoder.Finish(nil)
	encoder.Release()

	queue.Submit(cmdBuffer)
	cmdBuffer.Release()

	// Wait for GPU
	device.Poll(true)

	// Map staging buffer
	if err := stagingBuffer.MapAsync(device, wgpu.MapModeRead, 0, queryResultSize); err != nil {
		return fmt.Errorf("map staging buffer: %w", err)
	}

	// Read timestamp values
	ptr := stagingBuffer.GetMappedRange(0, queryResultSize)
	if ptr == nil {
		return fmt.Errorf("failed to get mapped range")
	}

	data := (*[16]byte)(ptr)
	startTimestamp := *(*uint64)(unsafe.Pointer(&data[0]))
	endTimestamp := *(*uint64)(unsafe.Pointer(&data[8]))

	stagingBuffer.Unmap()

	// Calculate elapsed ticks
	// Note: To convert to nanoseconds, you need the timestamp period
	// from the adapter (typically 1 ns/tick, but varies by GPU)
	elapsedTicks := endTimestamp - startTimestamp

	// Assume 1 ns/tick (common on most GPUs)
	// For accurate conversion, use adapter.GetTimestampPeriod() if available
	const assumedPeriodNs = 1.0
	elapsedNs := float64(elapsedTicks) * assumedPeriodNs

	fmt.Printf("Timestamp Query Results:\n")
	fmt.Printf("  Start timestamp:   %d ticks\n", startTimestamp)
	fmt.Printf("  End timestamp:     %d ticks\n", endTimestamp)
	fmt.Printf("  Elapsed ticks:     %d\n", elapsedTicks)
	fmt.Printf("  GPU execution time: ~%.3f ms (assuming 1 ns/tick)\n", elapsedNs/1_000_000)
	fmt.Println()
	fmt.Println("GPU timestamp queries provide accurate profiling!")

	return nil
}

// runWithCPUTiming demonstrates timing with CPU (fallback)
func runWithCPUTiming(device *wgpu.Device, queue *wgpu.Queue) error {
	fmt.Println("Demonstrating compute workload with CPU timing...")
	fmt.Println()

	// Create compute pipeline
	shaderCode := `
@group(0) @binding(0) var<storage, read_write> data: array<f32>;

@compute @workgroup_size(256)
fn main(@builtin(global_invocation_id) global_id: vec3<u32>) {
    let idx = global_id.x;
    if (idx < arrayLength(&data)) {
        var sum: f32 = data[idx];
        for (var i: u32 = 0u; i < 100u; i = i + 1u) {
            sum = sum * 1.01 + 0.01;
        }
        data[idx] = sum;
    }
}
`
	shader := device.CreateShaderModuleWGSL(shaderCode)
	if shader == nil {
		return fmt.Errorf("failed to create shader")
	}
	defer shader.Release()

	pipeline := device.CreateComputePipelineSimple(nil, shader, "main")
	if pipeline == nil {
		return fmt.Errorf("failed to create pipeline")
	}
	defer pipeline.Release()

	// Create buffer
	const numElements = 1024 * 1024
	bufferSize := uint64(numElements * 4)

	buffer := device.CreateBuffer(&wgpu.BufferDescriptor{
		Usage: gputypes.BufferUsageStorage,
		Size:  bufferSize,
	})
	if buffer == nil {
		return fmt.Errorf("failed to create buffer")
	}
	defer buffer.Release()

	// Create bind group
	bindGroupLayout := pipeline.GetBindGroupLayout(0)
	defer bindGroupLayout.Release()

	bindGroup := device.CreateBindGroupSimple(bindGroupLayout, []wgpu.BindGroupEntry{
		wgpu.BufferBindingEntry(0, buffer, 0, bufferSize),
	})
	if bindGroup == nil {
		return fmt.Errorf("failed to create bind group")
	}
	defer bindGroup.Release()

	// Time the GPU work with CPU timer
	start := time.Now()

	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		return fmt.Errorf("failed to create command encoder")
	}

	pass := encoder.BeginComputePass(nil)
	pass.SetPipeline(pipeline)
	pass.SetBindGroup(0, bindGroup, nil)
	pass.DispatchWorkgroups(numElements/256, 1, 1)
	pass.End()
	pass.Release()

	cmdBuffer := encoder.Finish(nil)
	encoder.Release()

	queue.Submit(cmdBuffer)
	cmdBuffer.Release()

	// Wait for GPU to complete
	device.Poll(true)

	elapsed := time.Since(start)

	fmt.Printf("Compute Workload Results (CPU timing):\n")
	fmt.Printf("  Elements processed: %d\n", numElements)
	fmt.Printf("  CPU-measured time:  %v\n", elapsed)
	fmt.Println()
	fmt.Println("NOTE: CPU timing includes submission overhead.")
	fmt.Println("      GPU timestamp queries provide more accurate GPU-only timing.")
	fmt.Println()
	fmt.Println("To enable GPU timestamps, request TIMESTAMP_QUERY feature")
	fmt.Println("when creating the device.")

	return nil
}
