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
		Label:            EmptyStringView(),
		Usage:            gputypes.BufferUsageCopyDst | gputypes.BufferUsageCopySrc,
		Size:             256,
		MappedAtCreation: False,
	}
	buffer := device.CreateBuffer(&bufferDesc)
	if buffer == nil {
		t.Fatal("Failed to create buffer")
	}
	defer buffer.Release()

	// Create command encoder
	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		t.Fatal("Failed to create command encoder")
	}
	defer encoder.Release()

	// Test ClearBuffer
	encoder.ClearBuffer(buffer, 0, 256)

	// Finish command buffer
	cmdBuffer := encoder.Finish(nil)
	if cmdBuffer == nil {
		t.Fatal("Failed to finish command encoder")
	}
	defer cmdBuffer.Release()

	// Submit
	queue := device.GetQueue()
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
	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		t.Fatal("Failed to create command encoder")
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
	cmdBuffer := encoder.Finish(nil)
	if cmdBuffer == nil {
		t.Fatal("Failed to finish command encoder")
	}
	defer cmdBuffer.Release()

	// Submit
	queue := device.GetQueue()
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
		Label: EmptyStringView(),
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
	texture := device.CreateTexture(&textureDesc)
	if texture == nil {
		t.Fatal("Failed to create texture")
	}
	defer texture.Release()

	// Test query methods
	width := texture.GetWidth()
	if width != 512 {
		t.Errorf("Expected width 512, got %d", width)
	}

	height := texture.GetHeight()
	if height != 256 {
		t.Errorf("Expected height 256, got %d", height)
	}

	depthOrLayers := texture.GetDepthOrArrayLayers()
	if depthOrLayers != 4 {
		t.Errorf("Expected depth/layers 4, got %d", depthOrLayers)
	}

	mipLevels := texture.GetMipLevelCount()
	if mipLevels != 3 {
		t.Errorf("Expected mip levels 3, got %d", mipLevels)
	}

	format := texture.GetFormat()
	if format != gputypes.TextureFormatRGBA8Unorm {
		t.Errorf("Expected format RGBA8Unorm, got %d", format)
	}
}

// TestTextureQueryAPIsNil tests texture query methods with nil texture.
func TestTextureQueryAPIsNil(t *testing.T) {
	var texture *Texture

	// All methods should return zero values and not panic
	if width := texture.GetWidth(); width != 0 {
		t.Errorf("Expected width 0 for nil texture, got %d", width)
	}

	if height := texture.GetHeight(); height != 0 {
		t.Errorf("Expected height 0 for nil texture, got %d", height)
	}

	if depth := texture.GetDepthOrArrayLayers(); depth != 0 {
		t.Errorf("Expected depth 0 for nil texture, got %d", depth)
	}

	if mips := texture.GetMipLevelCount(); mips != 0 {
		t.Errorf("Expected mip levels 0 for nil texture, got %d", mips)
	}

	if format := texture.GetFormat(); format != gputypes.TextureFormatUndefined {
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

	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		t.Fatal("Failed to create command encoder")
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

	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		t.Fatal("Failed to create command encoder")
	}
	defer encoder.Release()

	// Should not panic with empty strings
	encoder.InsertDebugMarker("")
	encoder.PushDebugGroup("")
	encoder.PopDebugGroup()
}
