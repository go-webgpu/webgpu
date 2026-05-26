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
	queue := device.Queue()
	if queue == nil {
		t.Fatal("Queue returned nil")
	}
	defer queue.Release()

	if queue.Handle() == 0 {
		t.Fatal("Queue handle is zero")
	}

	t.Logf("Queue obtained: handle=%#x", queue.Handle())
}

func TestDeviceGetLimits(t *testing.T) {
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

	t.Log("Getting device limits (cached, no FFI call)...")
	deviceLimits := device.Limits()

	// Limits are cached at RequestDevice time. On wgpu-native v29+ they should be non-zero.
	// On older wgpu-native (v27) GetLimits was not fully implemented; limits may be zero —
	// that is not a test failure, just log the result.
	t.Logf("Device limits: MaxTextureDimension2D=%d, MaxBindGroups=%d",
		deviceLimits.MaxTextureDimension2D,
		deviceLimits.MaxBindGroups)
}

func TestDeviceGetFeatures(t *testing.T) {
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

	t.Log("Getting device features...")
	features := device.Features()

	// Result can be nil (no optional features) - this is ok
	if features == nil {
		t.Log("Device has no optional features enabled")
	} else {
		t.Logf("Device has %d features enabled", len(features))
		for i, f := range features {
			t.Logf("  Feature %d: %d", i, f)
		}
	}
}

func TestDeviceHasFeature(t *testing.T) {
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

	t.Log("Checking for TimestampQuery feature...")
	// Test with FeatureNameTimestampQuery (may or may not be supported)
	// Just verify no crash
	hasTimestamp := device.HasFeature(FeatureNameTimestampQuery)
	t.Logf("Device has TimestampQuery: %v", hasTimestamp)
}

func TestDeviceGetLimits_Nil(t *testing.T) {
	var d *Device
	// Nil device returns zero-value Limits (no error, no panic).
	limits := d.Limits()
	if limits.MaxTextureDimension2D != 0 {
		t.Errorf("Expected zero-value Limits for nil device, got MaxTextureDimension2D=%d", limits.MaxTextureDimension2D)
	}
}
