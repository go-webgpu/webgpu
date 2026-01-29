package wgpu

import (
	"testing"

	"github.com/gogpu/gputypes"
)

func TestCreateRenderPipelineSimple(t *testing.T) {
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

	t.Log("Creating simple render pipeline...")
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

	if pipeline.Handle() == 0 {
		t.Fatal("RenderPipeline handle is zero")
	}

	t.Logf("RenderPipeline created: handle=%#x", pipeline.Handle())
}

func TestCreateRenderPipelineWithDescriptor(t *testing.T) {
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

	t.Log("Creating render pipeline with descriptor...")
	pipeline := device.CreateRenderPipeline(&RenderPipelineDescriptor{
		Vertex: VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
		},
		Fragment: &FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []ColorTargetState{{
				Format:    gputypes.TextureFormatBGRA8Unorm,
				WriteMask: gputypes.ColorWriteMaskAll,
			}},
		},
		Primitive: PrimitiveState{
			Topology:  gputypes.PrimitiveTopologyTriangleList,
			FrontFace: gputypes.FrontFaceCCW,
			CullMode:  gputypes.CullModeNone,
		},
		Multisample: MultisampleState{
			Count: 1,
			Mask:  0xFFFFFFFF,
		},
	})
	if pipeline == nil {
		t.Fatal("CreateRenderPipeline returned nil")
	}
	defer pipeline.Release()

	t.Logf("RenderPipeline with descriptor created: handle=%#x", pipeline.Handle())
}

func TestRenderPipelineGetBindGroupLayout(t *testing.T) {
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

	// Shader with uniform binding
	shaderCode := `
struct Uniforms {
    color: vec4<f32>,
};

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

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
    return uniforms.color;
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

	t.Log("Getting bind group layout from render pipeline...")
	layout := pipeline.GetBindGroupLayout(0)
	if layout == nil {
		t.Fatal("GetBindGroupLayout returned nil")
	}
	defer layout.Release()

	if layout.Handle() == 0 {
		t.Fatal("BindGroupLayout handle is zero")
	}

	t.Logf("BindGroupLayout from render pipeline: handle=%#x", layout.Handle())
}

func TestRenderPipelineWithDepth(t *testing.T) {
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

	shaderCode := `
@vertex
fn vs_main(@builtin(vertex_index) idx: u32) -> @builtin(position) vec4<f32> {
    var pos = array<vec2<f32>, 3>(
        vec2<f32>(0.0, 0.5),
        vec2<f32>(-0.5, -0.5),
        vec2<f32>(0.5, -0.5)
    );
    return vec4<f32>(pos[idx], 0.5, 1.0); // z = 0.5
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

	t.Log("Creating render pipeline with depth testing...")
	pipeline := device.CreateRenderPipeline(&RenderPipelineDescriptor{
		Vertex: VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
		},
		Fragment: &FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []ColorTargetState{{
				Format:    gputypes.TextureFormatBGRA8Unorm,
				WriteMask: gputypes.ColorWriteMaskAll,
			}},
		},
		DepthStencil: &DepthStencilState{
			Format:            gputypes.TextureFormatDepth24Plus,
			DepthWriteEnabled: true,
			DepthCompare:      gputypes.CompareFunctionLess,
		},
		Primitive: PrimitiveState{
			Topology: gputypes.PrimitiveTopologyTriangleList,
		},
		Multisample: MultisampleState{
			Count: 1,
			Mask:  0xFFFFFFFF,
		},
	})
	if pipeline == nil {
		t.Fatal("CreateRenderPipeline with depth returned nil")
	}
	defer pipeline.Release()

	t.Logf("RenderPipeline with depth: handle=%#x", pipeline.Handle())
}
