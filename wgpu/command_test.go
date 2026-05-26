package wgpu

import (
	"testing"
	"unsafe"

	"github.com/gogpu/gputypes"
)

const computeDoubleShader = `
@group(0) @binding(0) var<storage, read_write> data: array<f32>;

@compute @workgroup_size(64)
fn main(@builtin(global_invocation_id) global_id: vec3<u32>) {
    let idx = global_id.x;
    if (idx < arrayLength(&data)) {
        data[idx] = data[idx] * 2.0;
    }
}
`

func TestCreateCommandEncoder(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	t.Log("Creating command encoder...")
	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		t.Fatalf("CreateCommandEncoder failed: %v", err)
	}
	defer encoder.Release()

	if encoder.Handle() == 0 {
		t.Fatal("CommandEncoder handle is zero")
	}

	t.Logf("CommandEncoder created: handle=%#x", encoder.Handle())
}

func TestCommandEncoderFinish(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		t.Fatalf("CreateCommandEncoder failed: %v", err)
	}

	t.Log("Finishing command encoder...")
	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		t.Fatalf("Finish failed: %v", err)
	}
	defer cmdBuffer.Release()

	t.Logf("CommandBuffer created: handle=%#x", cmdBuffer.Handle())
}

func TestComputePassDispatch(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	queue := device.Queue()
	if queue == nil {
		t.Fatal("Queue returned nil")
	}
	defer queue.Release()

	// Create shader
	shader, err := device.CreateShaderModuleWGSL(computeDoubleShader)
	if err != nil {
		t.Fatalf("CreateShaderModuleWGSL failed: %v", err)
	}
	defer shader.Release()

	// Create compute pipeline with auto layout
	pipeline, err := device.CreateComputePipelineSimple(nil, shader, "main")
	if err != nil {
		t.Fatalf("CreateComputePipelineSimple failed: %v", err)
	}
	defer pipeline.Release()

	// Get bind group layout from pipeline
	bindGroupLayout := pipeline.GetBindGroupLayout(0)
	if bindGroupLayout == nil {
		t.Fatal("GetBindGroupLayout returned nil")
	}
	defer bindGroupLayout.Release()

	// Create buffer with initial data
	const numElements = 64
	bufferSize := uint64(numElements * 4) // 4 bytes per float32

	buffer, err := device.CreateBuffer(&BufferDescriptor{
		Label: "",
		Usage:            gputypes.BufferUsageStorage | gputypes.BufferUsageCopySrc | gputypes.BufferUsageCopyDst,
		Size:             bufferSize,
		MappedAtCreation: true,
	})
	if err != nil {
		t.Fatalf("CreateBuffer failed: %v", err)
	}
	defer buffer.Release()

	// Initialize buffer with test data (1.0, 2.0, 3.0, ...)
	ptr := buffer.GetMappedRange(0, bufferSize)
	if ptr == nil {
		t.Fatal("GetMappedRange returned nil")
	}
	data := (*[numElements]float32)(ptr)
	for i := range data {
		data[i] = float32(i + 1)
	}
	buffer.Unmap()
	t.Logf("Initialized buffer with %d elements", numElements)

	// Create bind group
	entries := []BindGroupEntry{
		BufferBindingEntry(0, buffer, 0, bufferSize),
	}
	bindGroup, err := device.CreateBindGroupSimple(bindGroupLayout, entries)
	if err != nil {
		t.Fatalf("CreateBindGroupSimple failed: %v", err)
	}
	defer bindGroup.Release()

	// Create and execute compute pass
	t.Log("Creating compute pass...")
	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		t.Fatalf("CreateCommandEncoder failed: %v", err)
	}

	pass, err := encoder.BeginComputePass(nil)
	if err != nil {
		t.Fatalf("BeginComputePass failed: %v", err)
	}

	pass.SetPipeline(pipeline)
	pass.SetBindGroup(0, bindGroup, nil)
	pass.DispatchWorkgroups(1, 1, 1) // 64 invocations (workgroup_size is 64)
	pass.End()
	pass.Release()

	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		t.Fatalf("Finish failed: %v", err)
	}

	t.Log("Submitting compute work...")
	queue.Submit(cmdBuffer)
	cmdBuffer.Release()

	t.Log("Compute pass dispatched successfully")
}

func TestFullComputeExample(t *testing.T) {
	// Full end-to-end compute example with result verification
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	queue := device.Queue()
	if queue == nil {
		t.Fatal("Queue returned nil")
	}
	defer queue.Release()

	// Create shader
	shader, err := device.CreateShaderModuleWGSL(computeDoubleShader)
	if err != nil {
		t.Fatalf("CreateShaderModuleWGSL failed: %v", err)
	}
	defer shader.Release()

	// Create compute pipeline
	pipeline, err := device.CreateComputePipelineSimple(nil, shader, "main")
	if err != nil {
		t.Fatalf("CreateComputePipelineSimple failed: %v", err)
	}
	defer pipeline.Release()

	bindGroupLayout := pipeline.GetBindGroupLayout(0)
	if bindGroupLayout == nil {
		t.Fatal("GetBindGroupLayout returned nil")
	}
	defer bindGroupLayout.Release()

	// Create storage buffer
	const numElements = 64
	bufferSize := uint64(numElements * 4)

	storageBuffer, err := device.CreateBuffer(&BufferDescriptor{
		Label: "",
		Usage:            gputypes.BufferUsageStorage | gputypes.BufferUsageCopySrc | gputypes.BufferUsageCopyDst,
		Size:             bufferSize,
		MappedAtCreation: true,
	})
	if err != nil {
		t.Fatalf("CreateBuffer (storage) failed: %v", err)
	}
	defer storageBuffer.Release()

	// Initialize with test data
	ptr := storageBuffer.GetMappedRange(0, bufferSize)
	inputData := (*[numElements]float32)(ptr)
	for i := range inputData {
		inputData[i] = float32(i + 1)
	}
	storageBuffer.Unmap()

	// Create readback buffer
	readbackBuffer, err := device.CreateBuffer(&BufferDescriptor{
		Label: "",
		Usage:            gputypes.BufferUsageMapRead | gputypes.BufferUsageCopyDst,
		Size:             bufferSize,
		MappedAtCreation: false,
	})
	if err != nil {
		t.Fatalf("CreateBuffer (readback) failed: %v", err)
	}
	defer readbackBuffer.Release()

	// Create bind group
	entries := []BindGroupEntry{
		BufferBindingEntry(0, storageBuffer, 0, bufferSize),
	}
	bindGroup, err := device.CreateBindGroupSimple(bindGroupLayout, entries)
	if err != nil {
		t.Fatalf("CreateBindGroupSimple failed: %v", err)
	}
	defer bindGroup.Release()

	// Run compute
	t.Log("Running compute shader...")
	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		t.Fatalf("CreateCommandEncoder failed: %v", err)
	}
	pass, err := encoder.BeginComputePass(nil)
	if err != nil {
		t.Fatalf("BeginComputePass failed: %v", err)
	}
	pass.SetPipeline(pipeline)
	pass.SetBindGroup(0, bindGroup, nil)
	pass.DispatchWorkgroups(1, 1, 1)
	pass.End()
	pass.Release()

	// Copy result to readback buffer
	encoder.CopyBufferToBuffer(storageBuffer, 0, readbackBuffer, 0, bufferSize)

	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		t.Fatalf("Finish failed: %v", err)
	}
	queue.Submit(cmdBuffer)
	cmdBuffer.Release()

	t.Log("Compute shader executed, result copied to readback buffer")

	// Map the readback buffer and verify results
	t.Log("Mapping readback buffer...")
	mapPending, err := readbackBuffer.MapAsync(MapModeRead, 0, bufferSize)
	if err != nil {
		t.Fatalf("MapAsync failed: %v", err)
	}
	// Drive polling until resolved (requires GPU — test is GPU-gated).
	for {
		if ready, werr := mapPending.Status(); ready {
			if werr != nil {
				t.Fatalf("MapAsync resolved with error: %v", werr)
			}
			break
		}
		device.Poll(false)
	}
	mapPending.Release()

	resultPtr := readbackBuffer.GetMappedRange(0, bufferSize)
	if resultPtr == nil {
		t.Fatal("GetMappedRange returned nil for readback buffer")
	}

	results := (*[numElements]float32)(resultPtr)

	// Verify: each value should be doubled (1*2=2, 2*2=4, ..., 64*2=128)
	t.Log("Verifying compute results...")
	for i := 0; i < numElements; i++ {
		expected := float32((i + 1) * 2)
		if results[i] != expected {
			t.Errorf("results[%d] = %f, want %f", i, results[i], expected)
		}
	}

	readbackBuffer.Unmap()
	t.Log("Full compute example with verification completed successfully!")
}

func TestCopyBufferToBuffer(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	queue := device.Queue()
	if queue == nil {
		t.Fatal("Queue returned nil")
	}
	defer queue.Release()

	// Create source buffer with data
	srcBuffer, err := device.CreateBuffer(&BufferDescriptor{
		Label: "",
		Usage:            gputypes.BufferUsageCopySrc | gputypes.BufferUsageCopyDst,
		Size:             256,
		MappedAtCreation: true,
	})
	if err != nil {
		t.Fatalf("CreateBuffer (src) failed: %v", err)
	}
	defer srcBuffer.Release()

	ptr := srcBuffer.GetMappedRange(0, 256)
	srcData := (*[256]byte)(ptr)
	for i := range srcData {
		srcData[i] = byte(i)
	}
	srcBuffer.Unmap()

	// Create destination buffer
	dstBuffer, err := device.CreateBuffer(&BufferDescriptor{
		Label: "",
		Usage:            gputypes.BufferUsageCopyDst | gputypes.BufferUsageMapRead,
		Size:             256,
		MappedAtCreation: false,
	})
	if err != nil {
		t.Fatalf("CreateBuffer (dst) failed: %v", err)
	}
	defer dstBuffer.Release()

	// Copy buffer
	t.Log("Copying buffer...")
	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		t.Fatalf("CreateCommandEncoder failed: %v", err)
	}
	encoder.CopyBufferToBuffer(srcBuffer, 0, dstBuffer, 0, 256)
	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		t.Fatalf("Finish failed: %v", err)
	}
	queue.Submit(cmdBuffer)
	cmdBuffer.Release()

	t.Log("Buffer copied successfully")
}

func TestQueueSubmitMultiple(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	queue := device.Queue()
	if queue == nil {
		t.Fatal("Queue returned nil")
	}
	defer queue.Release()

	// Create multiple command buffers
	t.Log("Creating multiple command buffers...")
	var cmdBuffers []*CommandBuffer
	for i := 0; i < 3; i++ {
		encoder, err := device.CreateCommandEncoder(nil)
		if err != nil {
			t.Fatalf("CreateCommandEncoder failed for buffer %d: %v", i, err)
		}
		cmdBuffer, err := encoder.Finish(nil)
		if err != nil {
			t.Fatalf("Finish failed for buffer %d: %v", i, err)
		}
		cmdBuffers = append(cmdBuffers, cmdBuffer)
	}

	// Submit all at once
	t.Log("Submitting multiple command buffers...")
	queue.Submit(cmdBuffers...)

	// Release
	for _, cb := range cmdBuffers {
		cb.Release()
	}

	t.Log("Multiple command buffers submitted successfully")
}

// Helper to check if two float slices are approximately equal
func floatsEqual(a, b []float32, tolerance float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		diff := a[i] - b[i]
		if diff < 0 {
			diff = -diff
		}
		if diff > tolerance {
			return false
		}
	}
	return true
}

// Unused but kept for future MapAsync implementation
var _ = unsafe.Pointer(nil)
