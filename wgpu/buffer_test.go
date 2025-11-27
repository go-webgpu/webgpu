package wgpu

import (
	"testing"
	"unsafe"
)

func TestCreateBuffer(t *testing.T) {
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

	t.Log("Creating buffer...")
	desc := &BufferDescriptor{
		Label:            EmptyStringView(),
		Usage:            BufferUsageCopyDst | BufferUsageMapRead,
		Size:             256,
		MappedAtCreation: False,
	}
	buffer := device.CreateBuffer(desc)
	if buffer == nil {
		t.Fatal("CreateBuffer returned nil")
	}
	defer buffer.Release()

	if buffer.Handle() == 0 {
		t.Fatal("Buffer handle is zero")
	}

	t.Logf("Buffer created: handle=%#x", buffer.Handle())

	// Test GetSize
	size := buffer.GetSize()
	if size != 256 {
		t.Errorf("Buffer size = %d, want 256", size)
	}
	t.Logf("Buffer size: %d bytes", size)
}

func TestBufferMappedAtCreation(t *testing.T) {
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

	t.Log("Creating buffer with MappedAtCreation...")
	desc := &BufferDescriptor{
		Label:            EmptyStringView(),
		Usage:            BufferUsageCopySrc,
		Size:             64,
		MappedAtCreation: True,
	}
	buffer := device.CreateBuffer(desc)
	if buffer == nil {
		t.Fatal("CreateBuffer returned nil")
	}
	defer buffer.Release()

	// Get mapped range
	ptr := buffer.GetMappedRange(0, 64)
	if ptr == nil {
		t.Fatal("GetMappedRange returned nil")
	}
	t.Logf("Mapped range: %p", ptr)

	// Write some data
	data := (*[64]byte)(ptr)
	for i := range data {
		data[i] = byte(i)
	}
	t.Log("Wrote test data to buffer")

	// Unmap
	buffer.Unmap()
	t.Log("Buffer unmapped")
}

func TestQueueWriteBuffer(t *testing.T) {
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
	if queue == nil {
		t.Fatal("GetQueue returned nil")
	}
	defer queue.Release()

	t.Log("Creating buffer for WriteBuffer test...")
	desc := &BufferDescriptor{
		Label:            EmptyStringView(),
		Usage:            BufferUsageCopyDst,
		Size:             128,
		MappedAtCreation: False,
	}
	buffer := device.CreateBuffer(desc)
	if buffer == nil {
		t.Fatal("CreateBuffer returned nil")
	}
	defer buffer.Release()

	// Write data to buffer
	testData := make([]byte, 64)
	for i := range testData {
		testData[i] = byte(i * 2)
	}

	t.Log("Writing data to buffer via queue...")
	queue.WriteBuffer(buffer, 0, testData)
	t.Log("WriteBuffer completed")
}

func TestQueueWriteBufferRaw(t *testing.T) {
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
	if queue == nil {
		t.Fatal("GetQueue returned nil")
	}
	defer queue.Release()

	t.Log("Creating buffer for WriteBufferRaw test...")
	desc := &BufferDescriptor{
		Label:            EmptyStringView(),
		Usage:            BufferUsageCopyDst,
		Size:             128,
		MappedAtCreation: False,
	}
	buffer := device.CreateBuffer(desc)
	if buffer == nil {
		t.Fatal("CreateBuffer returned nil")
	}
	defer buffer.Release()

	// Test with typed data (float32 array)
	floatData := []float32{1.0, 2.0, 3.0, 4.0}
	t.Log("Writing float32 data to buffer via queue...")
	queue.WriteBufferRaw(buffer, 0, unsafe.Pointer(&floatData[0]), uint64(len(floatData)*4))
	t.Log("WriteBufferRaw completed")
}
