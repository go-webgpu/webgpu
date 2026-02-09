package wgpu

import (
	"testing"
	"unsafe"
)

func TestInit(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	t.Log("Library initialized successfully")
}

func TestStructSizes(t *testing.T) {
	t.Logf("sizeof(ChainedStruct) = %d (expected 16)", unsafe.Sizeof(ChainedStruct{}))
	t.Logf("sizeof(InstanceCapabilities) = %d (expected 24)", unsafe.Sizeof(InstanceCapabilities{}))
	t.Logf("sizeof(InstanceDescriptor) = %d (expected 32)", unsafe.Sizeof(InstanceDescriptor{}))
}

func TestCreateInstanceWithNil(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance(nil) failed: %v", err)
	}
	defer inst.Release()

	if inst.Handle() == 0 {
		t.Fatal("Instance handle is zero")
	}

	t.Logf("Instance created successfully: handle=%#x", inst.Handle())
}

func TestCreateInstanceWithDescriptor(t *testing.T) {
	desc := &InstanceDescriptor{
		NextInChain: 0,
		Features: InstanceCapabilities{
			NextInChain:          0,
			TimedWaitAnyEnable:   False,
			TimedWaitAnyMaxCount: 0,
		},
	}

	inst, err := CreateInstance(desc)
	if err != nil {
		t.Fatalf("CreateInstance(desc) failed: %v", err)
	}
	defer inst.Release()

	if inst.Handle() == 0 {
		t.Fatal("Instance handle is zero")
	}

	t.Logf("Instance created successfully: handle=%#x", inst.Handle())
}

func TestInstanceRelease(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}

	handle := inst.Handle()
	t.Logf("Before release: handle=%#x", handle)

	inst.Release()

	if inst.Handle() != 0 {
		t.Fatal("Handle should be zero after release")
	}
	t.Log("Instance released successfully")
}

func TestCheckInitAfterLoad(t *testing.T) {
	// After library is loaded (which happens in TestInit), checkInit should return nil
	err := checkInit()
	if err != nil {
		t.Fatalf("checkInit() failed after successful library load: %v", err)
	}
	t.Log("checkInit() passed after library initialization")
}

func TestCreateInstanceReturnsErrLibraryNotLoaded(t *testing.T) {
	// This test documents the expected error type, but cannot test the actual
	// uninitialized state due to sync.Once - the library is already loaded.
	// The error value is defined and will be returned if Init() fails.
	if ErrLibraryNotLoaded == nil {
		t.Fatal("ErrLibraryNotLoaded should be defined")
	}
	t.Logf("ErrLibraryNotLoaded is defined: %v", ErrLibraryNotLoaded)
}
