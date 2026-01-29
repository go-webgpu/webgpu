package wgpu

import (
	"testing"

	"github.com/gogpu/gputypes"
)

func TestCreateBindGroupLayout(t *testing.T) {
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

	t.Log("Creating bind group layout...")
	entries := []BindGroupLayoutEntry{
		{
			Binding:    0,
			Visibility: gputypes.ShaderStageCompute,
			Buffer: BufferBindingLayout{
				Type:           gputypes.BufferBindingTypeStorage,
				MinBindingSize: 0,
			},
		},
	}
	layout := device.CreateBindGroupLayoutSimple(entries)
	if layout == nil {
		t.Fatal("CreateBindGroupLayoutSimple returned nil")
	}
	defer layout.Release()

	if layout.Handle() == 0 {
		t.Fatal("BindGroupLayout handle is zero")
	}

	t.Logf("BindGroupLayout created: handle=%#x", layout.Handle())
}

func TestCreateBindGroup(t *testing.T) {
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

	// Create buffer
	bufferDesc := &BufferDescriptor{
		Label:            EmptyStringView(),
		Usage:            gputypes.BufferUsageStorage | gputypes.BufferUsageCopyDst,
		Size:             256,
		MappedAtCreation: False,
	}
	buffer := device.CreateBuffer(bufferDesc)
	if buffer == nil {
		t.Fatal("CreateBuffer returned nil")
	}
	defer buffer.Release()

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
	layout := device.CreateBindGroupLayoutSimple(layoutEntries)
	if layout == nil {
		t.Fatal("CreateBindGroupLayoutSimple returned nil")
	}
	defer layout.Release()

	// Create bind group
	t.Log("Creating bind group...")
	entries := []BindGroupEntry{
		BufferBindingEntry(0, buffer, 0, 256),
	}
	bindGroup := device.CreateBindGroupSimple(layout, entries)
	if bindGroup == nil {
		t.Fatal("CreateBindGroupSimple returned nil")
	}
	defer bindGroup.Release()

	if bindGroup.Handle() == 0 {
		t.Fatal("BindGroup handle is zero")
	}

	t.Logf("BindGroup created: handle=%#x", bindGroup.Handle())
}

func TestBindGroupWithMultipleBindings(t *testing.T) {
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

	// Create buffers
	inputBuffer := device.CreateBuffer(&BufferDescriptor{
		Label:            EmptyStringView(),
		Usage:            gputypes.BufferUsageStorage | gputypes.BufferUsageCopyDst,
		Size:             256,
		MappedAtCreation: False,
	})
	if inputBuffer == nil {
		t.Fatal("CreateBuffer (input) returned nil")
	}
	defer inputBuffer.Release()

	outputBuffer := device.CreateBuffer(&BufferDescriptor{
		Label:            EmptyStringView(),
		Usage:            gputypes.BufferUsageStorage | gputypes.BufferUsageCopySrc,
		Size:             256,
		MappedAtCreation: False,
	})
	if outputBuffer == nil {
		t.Fatal("CreateBuffer (output) returned nil")
	}
	defer outputBuffer.Release()

	// Create layout with two bindings
	layoutEntries := []BindGroupLayoutEntry{
		{
			Binding:    0,
			Visibility: gputypes.ShaderStageCompute,
			Buffer: BufferBindingLayout{
				Type:           gputypes.BufferBindingTypeReadOnlyStorage,
				MinBindingSize: 0,
			},
		},
		{
			Binding:    1,
			Visibility: gputypes.ShaderStageCompute,
			Buffer: BufferBindingLayout{
				Type:           gputypes.BufferBindingTypeStorage,
				MinBindingSize: 0,
			},
		},
	}
	layout := device.CreateBindGroupLayoutSimple(layoutEntries)
	if layout == nil {
		t.Fatal("CreateBindGroupLayoutSimple returned nil")
	}
	defer layout.Release()

	// Create bind group
	t.Log("Creating bind group with multiple bindings...")
	entries := []BindGroupEntry{
		BufferBindingEntry(0, inputBuffer, 0, 256),
		BufferBindingEntry(1, outputBuffer, 0, 256),
	}
	bindGroup := device.CreateBindGroupSimple(layout, entries)
	if bindGroup == nil {
		t.Fatal("CreateBindGroupSimple returned nil")
	}
	defer bindGroup.Release()

	t.Logf("BindGroup with %d bindings created: handle=%#x", len(entries), bindGroup.Handle())
}
