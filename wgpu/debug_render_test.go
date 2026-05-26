package wgpu

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/gogpu/gputypes"
)

func TestDebugColorTargetState(t *testing.T) {
	// Create a colorTargetStateWire with known values
	target := colorTargetStateWire{
		nextInChain: 0,
		format:      uint32(gputypes.TextureFormatBGRA8Unorm), // Should be 27 (0x1B)
		writeMask:   uint64(gputypes.ColorWriteMaskAll),       // Should be 15 (0xF)
	}

	t.Logf("colorTargetStateWire size: %d", unsafe.Sizeof(target))
	t.Logf("format value: %d (0x%X)", target.format, target.format)
	t.Logf("writeMask value: %d (0x%X)", target.writeMask, target.writeMask)

	// Check field offsets
	t.Logf("nextInChain offset: %d", unsafe.Offsetof(target.nextInChain))
	t.Logf("format offset: %d", unsafe.Offsetof(target.format))
	t.Logf("blend offset: %d", unsafe.Offsetof(target.blend))
	t.Logf("writeMask offset: %d", unsafe.Offsetof(target.writeMask))

	// Dump raw bytes
	ptr := unsafe.Pointer(&target)
	bytes := unsafe.Slice((*byte)(ptr), unsafe.Sizeof(target))
	t.Logf("Raw bytes: %v", bytes)

	// Verify the format at expected position
	formatPtr := (*uint32)(unsafe.Pointer(uintptr(ptr) + 8))
	t.Logf("Value at offset 8 (format): %d (0x%X)", *formatPtr, *formatPtr)

	if *formatPtr != 27 {
		t.Errorf("Format at offset 8 should be 27 but is %d", *formatPtr)
	}
}

func TestDebugRenderPipelineBytes(t *testing.T) {
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

	// Check gputypes values
	// gputypes v0.3.0 and wgpu-native v29 have identical TextureFormat values.
	// BGRA8Unorm = 0x1B = 27 in both gputypes and v29 webgpu.h.
	t.Logf("gputypes.TextureFormatBGRA8Unorm = %d (0x%X)", gputypes.TextureFormatBGRA8Unorm, gputypes.TextureFormatBGRA8Unorm)
	t.Logf("gputypes.TextureFormatRG11B10Ufloat = %d (0x%X)", gputypes.TextureFormatRG11B10Ufloat, gputypes.TextureFormatRG11B10Ufloat)

	// Verify the direct cast: gputypes v0.3.0 BGRA8Unorm = 0x1B = 27, v29 BGRA8Unorm = 0x1B = 27 (match)
	converted := uint32(gputypes.TextureFormatBGRA8Unorm)
	t.Logf("uint32(BGRA8Unorm) = %d (0x%X)", converted, converted)

	if converted != 27 {
		t.Errorf("uint32(BGRA8Unorm) should return 27 (v29 wgpu-native value) but returned %d", converted)
	}

	// Manually create the structs to see what's happening
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
	shader, err := device.CreateShaderModuleWGSL(shaderCode)
	if err != nil {
		t.Fatalf("CreateShaderModuleWGSL failed: %v", err)
	}
	defer shader.Release()

	t.Logf("Shader module handle: 0x%X", shader.Handle())

	// Create the color target manually
	nativeTarget := colorTargetStateWire{
		nextInChain: 0,
		format:      27, // Hardcoded BGRA8Unorm
		writeMask:   15, // Hardcoded All
	}

	// Verify the bytes
	targetBytes := (*[32]byte)(unsafe.Pointer(&nativeTarget))
	t.Logf("nativeTarget bytes: %v", targetBytes[:])
	t.Logf("  Bytes 8-11 (format): %v", targetBytes[8:12])

	// Check what value is at offset 8
	formatVal := *(*uint32)(unsafe.Pointer(&targetBytes[8]))
	t.Logf("  Format at offset 8: %d (0x%X)", formatVal, formatVal)

	// Print struct sizes for comparison
	t.Logf("\nStruct sizes:")
	t.Logf("  colorTargetStateWire: %d", unsafe.Sizeof(colorTargetStateWire{}))
	t.Logf("  fragmentState: %d", unsafe.Sizeof(fragmentState{}))
	t.Logf("  vertexState: %d", unsafe.Sizeof(vertexState{}))
	t.Logf("  primitiveState: %d", unsafe.Sizeof(primitiveState{}))
	t.Logf("  multisampleState: %d", unsafe.Sizeof(multisampleState{}))
	t.Logf("  renderPipelineDescriptor: %d", unsafe.Sizeof(renderPipelineDescriptor{}))

	fmt.Println("Debug test complete - not calling CreateRenderPipeline to avoid crash")
}
