package wgpu

import (
	"testing"
	"unsafe"

	"github.com/gogpu/gputypes"
)

func TestCreateRenderBundleEncoder(t *testing.T) {
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

	t.Log("Creating render bundle encoder...")
	colorFormats := []gputypes.TextureFormat{gputypes.TextureFormatBGRA8Unorm}
	encoder := device.CreateRenderBundleEncoderSimple(colorFormats, gputypes.TextureFormatUndefined, 1)
	if encoder == nil {
		t.Fatal("CreateRenderBundleEncoderSimple returned nil")
	}
	defer encoder.Release()

	if encoder.Handle() == 0 {
		t.Fatal("RenderBundleEncoder handle is zero")
	}

	t.Logf("RenderBundleEncoder created: handle=%#x", encoder.Handle())
}

func TestRenderBundleEncoderFinish(t *testing.T) {
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

	colorFormats := []gputypes.TextureFormat{gputypes.TextureFormatBGRA8Unorm}
	encoder := device.CreateRenderBundleEncoderSimple(colorFormats, gputypes.TextureFormatUndefined, 1)
	if encoder == nil {
		t.Fatal("CreateRenderBundleEncoderSimple returned nil")
	}

	t.Log("Finishing render bundle encoder (empty bundle)...")
	bundle := encoder.Finish(nil)
	if bundle == nil {
		t.Fatal("Finish returned nil")
	}
	defer bundle.Release()

	if bundle.Handle() == 0 {
		t.Fatal("RenderBundle handle is zero")
	}

	t.Logf("RenderBundle created: handle=%#x", bundle.Handle())
}

func TestRenderBundleWithPipeline(t *testing.T) {
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

	// Create pipeline
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

	// Create render bundle encoder
	colorFormats := []gputypes.TextureFormat{gputypes.TextureFormatBGRA8Unorm}
	bundleEncoder := device.CreateRenderBundleEncoderSimple(colorFormats, gputypes.TextureFormatUndefined, 1)
	if bundleEncoder == nil {
		t.Fatal("CreateRenderBundleEncoderSimple returned nil")
	}

	t.Log("Recording render commands to bundle...")
	bundleEncoder.SetPipeline(pipeline)
	bundleEncoder.Draw(3, 1, 0, 0)

	bundle := bundleEncoder.Finish(nil)
	if bundle == nil {
		t.Fatal("Finish returned nil")
	}
	defer bundle.Release()

	t.Logf("RenderBundle with draw commands created: handle=%#x", bundle.Handle())
}

func TestRenderBundleWithVertexBuffer(t *testing.T) {
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

	// Create vertex buffer
	vertices := []float32{
		0.0, 0.5,
		-0.5, -0.5,
		0.5, -0.5,
	}
	bufferSize := uint64(len(vertices) * 4)
	vertexBuffer := device.CreateBuffer(&BufferDescriptor{
		Usage:            gputypes.BufferUsageVertex | gputypes.BufferUsageCopyDst,
		Size:             bufferSize,
		MappedAtCreation: True,
	})
	if vertexBuffer == nil {
		t.Fatal("CreateBuffer returned nil")
	}
	defer vertexBuffer.Release()

	ptr := vertexBuffer.GetMappedRange(0, bufferSize)
	if ptr != nil {
		mapped := unsafe.Slice((*float32)(ptr), len(vertices))
		copy(mapped, vertices)
	}
	vertexBuffer.Unmap()

	// Use simple pipeline without vertex input (vertex_index)
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

	// Create render bundle with vertex buffer
	colorFormats := []gputypes.TextureFormat{gputypes.TextureFormatBGRA8Unorm}
	bundleEncoder := device.CreateRenderBundleEncoderSimple(colorFormats, gputypes.TextureFormatUndefined, 1)
	if bundleEncoder == nil {
		t.Fatal("CreateRenderBundleEncoderSimple returned nil")
	}

	t.Log("Recording render commands with vertex buffer to bundle...")
	bundleEncoder.SetPipeline(pipeline)
	// Set vertex buffer even though shader doesn't use it - testing the API
	bundleEncoder.SetVertexBuffer(0, vertexBuffer, 0, bufferSize)
	bundleEncoder.Draw(3, 1, 0, 0)

	bundle := bundleEncoder.Finish(nil)
	if bundle == nil {
		t.Fatal("Finish returned nil")
	}
	defer bundle.Release()

	t.Logf("RenderBundle with vertex buffer created: handle=%#x", bundle.Handle())
}
