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

	t.Log("Getting device limits...")
	deviceLimits, err := device.GetLimits()

	// NOTE: Currently wgpu-native v27 returns error from wgpuDeviceGetLimits
	// See issue #3 - waiting for wgpu-native webgpu-headers update
	if err != nil {
		t.Logf("Device GetLimits returned expected error (wgpu-native v27 limitation): %v", err)
		t.Skip("Skipping test - Device.GetLimits not yet supported in wgpu-native v27")
		return
	}

	// If we get here, wgpu-native was updated and the function works!
	// Verify limits are non-zero
	if deviceLimits.Limits.MaxTextureDimension2D == 0 {
		t.Error("MaxTextureDimension2D is zero")
	}

	t.Logf("Device limits: MaxTextureDimension2D=%d, MaxBindGroups=%d",
		deviceLimits.Limits.MaxTextureDimension2D,
		deviceLimits.Limits.MaxBindGroups)
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
	features := device.GetFeatures()

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
	_, err := d.GetLimits()
	if err == nil {
		t.Error("Expected error for nil device")
	}
	if err.Error() != "wgpu: Device.GetLimits: device is nil" {
		t.Errorf("Unexpected error message: %v", err)
	}
}
