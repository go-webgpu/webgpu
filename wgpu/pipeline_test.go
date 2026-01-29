package wgpu

import (
	"testing"

	"github.com/gogpu/gputypes"
)

const computeShaderDouble = `
@group(0) @binding(0) var<storage, read_write> data: array<f32>;

@compute @workgroup_size(64)
fn main(@builtin(global_invocation_id) global_id: vec3<u32>) {
    let idx = global_id.x;
    data[idx] = data[idx] * 2.0;
}
`

func TestCreatePipelineLayout(t *testing.T) {
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

	// Create bind group layout
	layoutEntries := []BindGroupLayoutEntry{
		{
			Binding:    0,
			Visibility: gputypes.ShaderStageCompute,
			Buffer: BufferBindingLayout{
				Type:           gputypes.BufferBindingTypeStorage,
				MinBindingSize: 0,
			},
		},
	}
	bindGroupLayout := device.CreateBindGroupLayoutSimple(layoutEntries)
	if bindGroupLayout == nil {
		t.Fatal("CreateBindGroupLayoutSimple returned nil")
	}
	defer bindGroupLayout.Release()

	// Create pipeline layout
	t.Log("Creating pipeline layout...")
	pipelineLayout := device.CreatePipelineLayoutSimple([]*BindGroupLayout{bindGroupLayout})
	if pipelineLayout == nil {
		t.Fatal("CreatePipelineLayoutSimple returned nil")
	}
	defer pipelineLayout.Release()

	if pipelineLayout.Handle() == 0 {
		t.Fatal("PipelineLayout handle is zero")
	}

	t.Logf("PipelineLayout created: handle=%#x", pipelineLayout.Handle())
}

func TestCreateComputePipeline(t *testing.T) {
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

	// Create shader
	shader := device.CreateShaderModuleWGSL(computeShaderDouble)
	if shader == nil {
		t.Fatal("CreateShaderModuleWGSL returned nil")
	}
	defer shader.Release()

	// Create pipeline with auto layout
	t.Log("Creating compute pipeline with auto layout...")
	pipeline := device.CreateComputePipelineSimple(nil, shader, "main")
	if pipeline == nil {
		t.Fatal("CreateComputePipelineSimple returned nil")
	}
	defer pipeline.Release()

	if pipeline.Handle() == 0 {
		t.Fatal("ComputePipeline handle is zero")
	}

	t.Logf("ComputePipeline created: handle=%#x", pipeline.Handle())
}

func TestComputePipelineGetBindGroupLayout(t *testing.T) {
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

	// Create shader
	shader := device.CreateShaderModuleWGSL(computeShaderDouble)
	if shader == nil {
		t.Fatal("CreateShaderModuleWGSL returned nil")
	}
	defer shader.Release()

	// Create pipeline with auto layout
	pipeline := device.CreateComputePipelineSimple(nil, shader, "main")
	if pipeline == nil {
		t.Fatal("CreateComputePipelineSimple returned nil")
	}
	defer pipeline.Release()

	// Get bind group layout
	t.Log("Getting bind group layout from pipeline...")
	layout := pipeline.GetBindGroupLayout(0)
	if layout == nil {
		t.Fatal("GetBindGroupLayout returned nil")
	}
	defer layout.Release()

	t.Logf("Auto BindGroupLayout: handle=%#x", layout.Handle())
}

func TestCreateComputePipelineWithExplicitLayout(t *testing.T) {
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

	// Create shader
	shader := device.CreateShaderModuleWGSL(computeShaderDouble)
	if shader == nil {
		t.Fatal("CreateShaderModuleWGSL returned nil")
	}
	defer shader.Release()

	// Create bind group layout
	layoutEntries := []BindGroupLayoutEntry{
		{
			Binding:    0,
			Visibility: gputypes.ShaderStageCompute,
			Buffer: BufferBindingLayout{
				Type:           gputypes.BufferBindingTypeStorage,
				MinBindingSize: 0,
			},
		},
	}
	bindGroupLayout := device.CreateBindGroupLayoutSimple(layoutEntries)
	if bindGroupLayout == nil {
		t.Fatal("CreateBindGroupLayoutSimple returned nil")
	}
	defer bindGroupLayout.Release()

	// Create pipeline layout
	pipelineLayout := device.CreatePipelineLayoutSimple([]*BindGroupLayout{bindGroupLayout})
	if pipelineLayout == nil {
		t.Fatal("CreatePipelineLayoutSimple returned nil")
	}
	defer pipelineLayout.Release()

	// Create pipeline with explicit layout
	t.Log("Creating compute pipeline with explicit layout...")
	pipeline := device.CreateComputePipelineSimple(pipelineLayout, shader, "main")
	if pipeline == nil {
		t.Fatal("CreateComputePipelineSimple returned nil")
	}
	defer pipeline.Release()

	t.Logf("ComputePipeline with explicit layout: handle=%#x", pipeline.Handle())
}
