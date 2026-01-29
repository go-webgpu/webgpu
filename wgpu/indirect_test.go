package wgpu

import (
	"testing"
	"unsafe"

	"github.com/gogpu/gputypes"
)

func TestDispatchWorkgroupsIndirect(t *testing.T) {
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

	queue := device.GetQueue()
	defer queue.Release()

	// Create compute shader
	shaderCode := `
@group(0) @binding(0) var<storage, read_write> data: array<f32>;

@compute @workgroup_size(64)
fn main(@builtin(global_invocation_id) global_id: vec3<u32>) {
    let idx = global_id.x;
    if (idx < arrayLength(&data)) {
        data[idx] = data[idx] * 2.0;
    }
}
`
	shader := device.CreateShaderModuleWGSL(shaderCode)
	if shader == nil {
		t.Fatal("CreateShaderModuleWGSL returned nil")
	}
	defer shader.Release()

	pipeline := device.CreateComputePipelineSimple(nil, shader, "main")
	if pipeline == nil {
		t.Fatal("CreateComputePipelineSimple returned nil")
	}
	defer pipeline.Release()

	// Create storage buffer with test data
	const numElements = 64
	inputData := make([]float32, numElements)
	for i := range inputData {
		inputData[i] = float32(i + 1)
	}
	bufferSize := uint64(numElements * 4)

	storageBuffer := device.CreateBuffer(&BufferDescriptor{
		Usage:            gputypes.BufferUsageStorage | gputypes.BufferUsageCopySrc | gputypes.BufferUsageCopyDst,
		Size:             bufferSize,
		MappedAtCreation: True,
	})
	if storageBuffer == nil {
		t.Fatal("CreateBuffer for storage returned nil")
	}
	defer storageBuffer.Release()

	ptr := storageBuffer.GetMappedRange(0, bufferSize)
	if ptr != nil {
		mapped := unsafe.Slice((*float32)(ptr), numElements)
		copy(mapped, inputData)
	}
	storageBuffer.Unmap()

	// Create indirect buffer with dispatch args
	indirectArgs := DispatchIndirectArgs{
		WorkgroupCountX: 1, // 1 workgroup of 64 threads = 64 elements
		WorkgroupCountY: 1,
		WorkgroupCountZ: 1,
	}
	indirectSize := uint64(unsafe.Sizeof(indirectArgs))

	indirectBuffer := device.CreateBuffer(&BufferDescriptor{
		Usage:            gputypes.BufferUsageIndirect | gputypes.BufferUsageCopyDst,
		Size:             indirectSize,
		MappedAtCreation: True,
	})
	if indirectBuffer == nil {
		t.Fatal("CreateBuffer for indirect returned nil")
	}
	defer indirectBuffer.Release()

	indirectPtr := indirectBuffer.GetMappedRange(0, indirectSize)
	if indirectPtr != nil {
		*(*DispatchIndirectArgs)(indirectPtr) = indirectArgs
	}
	indirectBuffer.Unmap()

	// Create bind group
	bindGroupLayout := pipeline.GetBindGroupLayout(0)
	if bindGroupLayout == nil {
		t.Fatal("GetBindGroupLayout returned nil")
	}
	defer bindGroupLayout.Release()

	bindGroup := device.CreateBindGroupSimple(bindGroupLayout, []BindGroupEntry{
		BufferBindingEntry(0, storageBuffer, 0, bufferSize),
	})
	if bindGroup == nil {
		t.Fatal("CreateBindGroupSimple returned nil")
	}
	defer bindGroup.Release()

	// Create and submit command buffer
	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		t.Fatal("CreateCommandEncoder returned nil")
	}

	computePass := encoder.BeginComputePass(nil)
	computePass.SetPipeline(pipeline)
	computePass.SetBindGroup(0, bindGroup, nil)

	t.Log("Dispatching workgroups indirectly...")
	computePass.DispatchWorkgroupsIndirect(indirectBuffer, 0)

	computePass.End()
	computePass.Release()

	cmdBuffer := encoder.Finish(nil)
	encoder.Release()
	queue.Submit(cmdBuffer)
	cmdBuffer.Release()

	t.Log("DispatchWorkgroupsIndirect executed successfully")
}

func TestDrawIndirectArgs(t *testing.T) {
	// Test that DrawIndirectArgs has correct size and layout
	args := DrawIndirectArgs{
		VertexCount:   100,
		InstanceCount: 10,
		FirstVertex:   0,
		FirstInstance: 0,
	}

	expectedSize := uintptr(16) // 4 * uint32
	actualSize := unsafe.Sizeof(args)

	if actualSize != expectedSize {
		t.Errorf("DrawIndirectArgs size: expected %d, got %d", expectedSize, actualSize)
	}

	t.Logf("DrawIndirectArgs size: %d bytes", actualSize)
}

func TestDrawIndexedIndirectArgs(t *testing.T) {
	// Test that DrawIndexedIndirectArgs has correct size and layout
	args := DrawIndexedIndirectArgs{
		IndexCount:    100,
		InstanceCount: 10,
		FirstIndex:    0,
		BaseVertex:    0,
		FirstInstance: 0,
	}

	expectedSize := uintptr(20) // 5 * uint32/int32
	actualSize := unsafe.Sizeof(args)

	if actualSize != expectedSize {
		t.Errorf("DrawIndexedIndirectArgs size: expected %d, got %d", expectedSize, actualSize)
	}

	t.Logf("DrawIndexedIndirectArgs size: %d bytes", actualSize)
}

func TestDispatchIndirectArgs(t *testing.T) {
	// Test that DispatchIndirectArgs has correct size and layout
	args := DispatchIndirectArgs{
		WorkgroupCountX: 16,
		WorkgroupCountY: 16,
		WorkgroupCountZ: 1,
	}

	expectedSize := uintptr(12) // 3 * uint32
	actualSize := unsafe.Sizeof(args)

	if actualSize != expectedSize {
		t.Errorf("DispatchIndirectArgs size: expected %d, got %d", expectedSize, actualSize)
	}

	t.Logf("DispatchIndirectArgs size: %d bytes", actualSize)
}

func TestRenderBundleDrawIndirect(t *testing.T) {
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

	// Create indirect buffer
	indirectArgs := DrawIndirectArgs{
		VertexCount:   3,
		InstanceCount: 1,
		FirstVertex:   0,
		FirstInstance: 0,
	}
	indirectSize := uint64(unsafe.Sizeof(indirectArgs))

	indirectBuffer := device.CreateBuffer(&BufferDescriptor{
		Usage:            gputypes.BufferUsageIndirect | gputypes.BufferUsageCopyDst,
		Size:             indirectSize,
		MappedAtCreation: True,
	})
	if indirectBuffer == nil {
		t.Fatal("CreateBuffer for indirect returned nil")
	}
	defer indirectBuffer.Release()

	indirectPtr := indirectBuffer.GetMappedRange(0, indirectSize)
	if indirectPtr != nil {
		*(*DrawIndirectArgs)(indirectPtr) = indirectArgs
	}
	indirectBuffer.Unmap()

	// Create shader and pipeline
	shaderCode := `
@vertex
fn vs_main(@builtin(vertex_index) idx: u32) -> @builtin(position) vec4<f32> {
    var pos = array<vec2<f32>, 3>(
        vec2<f32>(0.0, 0.5),
        vec2<f32>(-0.5, -0.5),
        vec2<f32>(0.5, -0.5)
    );
    return vec4<f32>(pos[idx], 0.0, 1.0);
}

@fragment
fn fs_main() -> @location(0) vec4<f32> {
    return vec4<f32>(1.0, 0.0, 0.0, 1.0);
}
`
	shader := device.CreateShaderModuleWGSL(shaderCode)
	if shader == nil {
		t.Fatal("CreateShaderModuleWGSL returned nil")
	}
	defer shader.Release()

	pipeline := device.CreateRenderPipelineSimple(
		nil,
		shader, "vs_main",
		shader, "fs_main",
		gputypes.TextureFormatBGRA8Unorm,
	)
	if pipeline == nil {
		t.Fatal("CreateRenderPipelineSimple returned nil")
	}
	defer pipeline.Release()

	// Create render bundle with indirect draw
	colorFormats := []gputypes.TextureFormat{gputypes.TextureFormatBGRA8Unorm}
	bundleEncoder := device.CreateRenderBundleEncoderSimple(colorFormats, gputypes.TextureFormatUndefined, 1)
	if bundleEncoder == nil {
		t.Fatal("CreateRenderBundleEncoderSimple returned nil")
	}

	t.Log("Recording DrawIndirect to render bundle...")
	bundleEncoder.SetPipeline(pipeline)
	bundleEncoder.DrawIndirect(indirectBuffer, 0)

	bundle := bundleEncoder.Finish(nil)
	if bundle == nil {
		t.Fatal("Finish returned nil")
	}
	defer bundle.Release()

	t.Logf("RenderBundle with DrawIndirect created: handle=%#x", bundle.Handle())
}
