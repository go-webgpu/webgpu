package wgpu

import (
	"testing"

	"github.com/gogpu/gputypes"
)

// TestCommandEncoderClearBuffer tests buffer clearing functionality.
func TestCommandEncoderClearBuffer(t *testing.T) {
	instance, err := CreateInstance(nil)
	if err != nil {
		t.Fatal("Failed to create instance:", err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		t.Fatal("Failed to request adapter:", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatal("Failed to request device:", err)
	}
	defer device.Release()

	// Create a buffer
	bufferDesc := BufferDescriptor{
		Label:            "",
		Usage:            gputypes.BufferUsageCopyDst | gputypes.BufferUsageCopySrc,
		Size:             256,
		MappedAtCreation: false,
	}
	buffer, err := device.CreateBuffer(&bufferDesc)
	if err != nil {
		t.Fatal("Failed to create buffer:", err)
	}
	defer buffer.Release()

	// Create command encoder
	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		t.Fatal("Failed to create command encoder:", err)
	}
	defer encoder.Release()

	// Test ClearBuffer
	encoder.ClearBuffer(buffer, 0, 256)

	// Finish command buffer
	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		t.Fatal("Failed to finish command encoder:", err)
	}
	defer cmdBuffer.Release()

	// Submit
	queue := device.Queue()
	if queue == nil {
		t.Fatal("Failed to get queue")
	}
	defer queue.Release()

	queue.Submit(cmdBuffer)
	device.Poll(true)
}

// TestCommandEncoderDebugMarkers tests debug marker functionality.
func TestCommandEncoderDebugMarkers(t *testing.T) {
	instance, err := CreateInstance(nil)
	if err != nil {
		t.Fatal("Failed to create instance:", err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		t.Fatal("Failed to request adapter:", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatal("Failed to request device:", err)
	}
	defer device.Release()

	// Create command encoder
	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		t.Fatal("Failed to create command encoder:", err)
	}
	defer encoder.Release()

	// Test debug markers
	encoder.PushDebugGroup("Test Group")
	encoder.InsertDebugMarker("Test Marker 1")
	encoder.InsertDebugMarker("Test Marker 2")
	encoder.PopDebugGroup()

	// Nested groups
	encoder.PushDebugGroup("Outer Group")
	encoder.PushDebugGroup("Inner Group")
	encoder.InsertDebugMarker("Nested Marker")
	encoder.PopDebugGroup()
	encoder.PopDebugGroup()

	// Finish command buffer
	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		t.Fatal("Failed to finish command encoder:", err)
	}
	defer cmdBuffer.Release()

	// Submit
	queue := device.Queue()
	if queue == nil {
		t.Fatal("Failed to get queue")
	}
	defer queue.Release()

	queue.Submit(cmdBuffer)
	device.Poll(true)
}

// TestTextureQueryAPIs tests texture query methods.
func TestTextureQueryAPIs(t *testing.T) {
	instance, err := CreateInstance(nil)
	if err != nil {
		t.Fatal("Failed to create instance:", err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		t.Fatal("Failed to request adapter:", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatal("Failed to request device:", err)
	}
	defer device.Release()

	// Create a texture
	textureDesc := TextureDescriptor{
		Label: "",
		Usage: gputypes.TextureUsageTextureBinding | gputypes.TextureUsageCopyDst,
		Size: gputypes.Extent3D{
			Width:              512,
			Height:             256,
			DepthOrArrayLayers: 4,
		},
		Format:        gputypes.TextureFormatRGBA8Unorm,
		MipLevelCount: 3,
		SampleCount:   1,
	}
	texture, err := device.CreateTexture(&textureDesc)
	if err != nil {
		t.Fatal("Failed to create texture:", err)
	}
	defer texture.Release()

	// Test query methods
	width := texture.Width()
	if width != 512 {
		t.Errorf("Expected width 512, got %d", width)
	}

	height := texture.Height()
	if height != 256 {
		t.Errorf("Expected height 256, got %d", height)
	}

	depthOrLayers := texture.DepthOrArrayLayers()
	if depthOrLayers != 4 {
		t.Errorf("Expected depth/layers 4, got %d", depthOrLayers)
	}

	mipLevels := texture.MipLevelCount()
	if mipLevels != 3 {
		t.Errorf("Expected mip levels 3, got %d", mipLevels)
	}

	format := texture.Format()
	if format != gputypes.TextureFormatRGBA8Unorm {
		t.Errorf("Expected format RGBA8Unorm, got %d", format)
	}
}

// TestTextureQueryAPIsNil tests texture query methods with nil texture.
func TestTextureQueryAPIsNil(t *testing.T) {
	var texture *Texture

	// All methods should return zero values and not panic
	if width := texture.Width(); width != 0 {
		t.Errorf("Expected width 0 for nil texture, got %d", width)
	}

	if height := texture.Height(); height != 0 {
		t.Errorf("Expected height 0 for nil texture, got %d", height)
	}

	if depth := texture.DepthOrArrayLayers(); depth != 0 {
		t.Errorf("Expected depth 0 for nil texture, got %d", depth)
	}

	if mips := texture.MipLevelCount(); mips != 0 {
		t.Errorf("Expected mip levels 0 for nil texture, got %d", mips)
	}

	if format := texture.Format(); format != gputypes.TextureFormatUndefined {
		t.Errorf("Expected format Undefined for nil texture, got %d", format)
	}
}

// TestClearBufferNil tests ClearBuffer with nil buffer.
func TestClearBufferNil(t *testing.T) {
	instance, err := CreateInstance(nil)
	if err != nil {
		t.Fatal("Failed to create instance:", err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		t.Fatal("Failed to request adapter:", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatal("Failed to request device:", err)
	}
	defer device.Release()

	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		t.Fatal("Failed to create command encoder:", err)
	}
	defer encoder.Release()

	// Should not panic with nil buffer
	encoder.ClearBuffer(nil, 0, 0)
}

// TestDebugMarkersEmptyStrings tests debug markers with empty strings.
func TestDebugMarkersEmptyStrings(t *testing.T) {
	instance, err := CreateInstance(nil)
	if err != nil {
		t.Fatal("Failed to create instance:", err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		t.Fatal("Failed to request adapter:", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatal("Failed to request device:", err)
	}
	defer device.Release()

	encoder, err := device.CreateCommandEncoder(nil)
	if err != nil {
		t.Fatal("Failed to create command encoder:", err)
	}
	defer encoder.Release()

	// Should not panic with empty strings
	encoder.InsertDebugMarker("")
	encoder.PushDebugGroup("")
	encoder.PopDebugGroup()
}
