package wgpu

import (
	"testing"
)

func TestRequestAdapter(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	t.Log("Requesting adapter...")
	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	if adapter.Handle() == 0 {
		t.Fatal("Adapter handle is zero")
	}

	t.Logf("Adapter obtained: handle=%#x", adapter.Handle())
}

func TestRequestAdapterWithOptions(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	options := &RequestAdapterOptions{
		FeatureLevel:    FeatureLevelCore,
		PowerPreference: PowerPreferenceHighPerformance,
	}

	t.Log("Requesting high-performance adapter...")
	adapter, err := inst.RequestAdapter(options)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	t.Logf("Adapter obtained: handle=%#x", adapter.Handle())
}
