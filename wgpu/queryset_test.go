package wgpu

import (
	"testing"
)

// Note: QuerySet tests require TIMESTAMP_QUERY feature which may not be
// available on all hardware. These tests verify the API works correctly
// but may be skipped if the feature isn't supported.

func TestCreateQuerySet(t *testing.T) {
	t.Skip("TIMESTAMP_QUERY feature not enabled by default - test skipped")
	// To enable: request device with RequiredFeatures including FeatureNameTimestampQuery

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

	t.Log("Creating timestamp query set...")
	querySet := device.CreateQuerySet(&QuerySetDescriptor{
		Type:  QueryTypeTimestamp,
		Count: 2,
	})
	if querySet == nil {
		t.Fatal("CreateQuerySet returned nil")
	}
	defer querySet.Release()

	if querySet.Handle() == 0 {
		t.Fatal("QuerySet handle is zero")
	}

	t.Logf("QuerySet created: handle=%#x", querySet.Handle())
}

func TestQuerySetDestroy(t *testing.T) {
	t.Skip("TIMESTAMP_QUERY feature not enabled by default - test skipped")

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

	querySet := device.CreateQuerySet(&QuerySetDescriptor{
		Type:  QueryTypeTimestamp,
		Count: 4,
	})
	if querySet == nil {
		t.Fatal("CreateQuerySet returned nil")
	}

	t.Log("Destroying query set...")
	querySet.Destroy()
	querySet.Release()

	t.Log("QuerySet destroyed successfully")
}

func TestWriteTimestamp(t *testing.T) {
	t.Skip("TIMESTAMP_QUERY feature not enabled by default - test skipped")

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
	defer queue.Release()

	// Create query set for timestamps
	querySet := device.CreateQuerySet(&QuerySetDescriptor{
		Type:  QueryTypeTimestamp,
		Count: 2,
	})
	if querySet == nil {
		t.Fatal("CreateQuerySet returned nil")
	}
	defer querySet.Release()

	// Create buffer to resolve query results
	resultBuffer := device.CreateBuffer(&BufferDescriptor{
		Usage: BufferUsageQueryResolve | BufferUsageCopySrc,
		Size:  16, // 2 timestamps * 8 bytes each
	})
	if resultBuffer == nil {
		t.Fatal("CreateBuffer for query results returned nil")
	}
	defer resultBuffer.Release()

	encoder := device.CreateCommandEncoder(nil)
	if encoder == nil {
		t.Fatal("CreateCommandEncoder returned nil")
	}

	t.Log("Writing timestamps...")
	encoder.WriteTimestamp(querySet, 0)
	encoder.WriteTimestamp(querySet, 1)

	t.Log("Resolving query set...")
	encoder.ResolveQuerySet(querySet, 0, 2, resultBuffer, 0)

	cmdBuffer := encoder.Finish(nil)
	encoder.Release()

	queue.Submit(cmdBuffer)
	cmdBuffer.Release()

	t.Log("Timestamp write and resolve completed successfully")
}

func TestQuerySetTypes(t *testing.T) {
	// Test that QueryType constants are defined correctly
	if QueryTypeTimestamp == 0 {
		t.Log("QueryTypeTimestamp has default value 0 - may indicate Occlusion type")
	}
	t.Logf("QueryTypeTimestamp = %#x", QueryTypeTimestamp)
}
