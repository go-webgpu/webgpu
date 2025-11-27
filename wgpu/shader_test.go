package wgpu

import (
	"testing"
)

const testComputeShader = `
@group(0) @binding(0) var<storage, read_write> data: array<f32>;

@compute @workgroup_size(64)
fn main(@builtin(global_invocation_id) global_id: vec3<u32>) {
    let idx = global_id.x;
    data[idx] = data[idx] * 2.0;
}
`

func TestCreateShaderModuleWGSL(t *testing.T) {
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

	t.Log("Creating shader module from WGSL...")
	shader := device.CreateShaderModuleWGSL(testComputeShader)
	if shader == nil {
		t.Fatal("CreateShaderModuleWGSL returned nil")
	}
	defer shader.Release()

	if shader.Handle() == 0 {
		t.Fatal("ShaderModule handle is zero")
	}

	t.Logf("ShaderModule created: handle=%#x", shader.Handle())
}

const simpleShader = `
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

func TestCreateShaderModuleVertex(t *testing.T) {
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

	t.Log("Creating vertex/fragment shader module...")
	shader := device.CreateShaderModuleWGSL(simpleShader)
	if shader == nil {
		t.Fatal("CreateShaderModuleWGSL returned nil")
	}
	defer shader.Release()

	t.Logf("Vertex/Fragment ShaderModule created: handle=%#x", shader.Handle())
}
