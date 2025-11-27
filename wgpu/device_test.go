package wgpu

import (
	"testing"
)

func TestRequestDevice(t *testing.T) {
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

	t.Log("Requesting device...")
	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	if device.Handle() == 0 {
		t.Fatal("Device handle is zero")
	}

	t.Logf("Device obtained: handle=%#x", device.Handle())
}

func TestDeviceGetQueue(t *testing.T) {
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

	t.Log("Getting queue...")
	queue := device.GetQueue()
	if queue == nil {
		t.Fatal("GetQueue returned nil")
	}
	defer queue.Release()

	if queue.Handle() == 0 {
		t.Fatal("Queue handle is zero")
	}

	t.Logf("Queue obtained: handle=%#x", queue.Handle())
}
