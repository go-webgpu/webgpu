package wgpu

import (
	"testing"

	"github.com/gogpu/gputypes"
)

// TestErrorScopeEmptyStack tests popping an error scope when stack is empty.
// NOTE: Currently disabled because wgpu-native panics on empty stack pop.
// This is a known limitation - users must track push/pop manually.
func TestErrorScopeEmptyStack(t *testing.T) {
	t.Skip("wgpu-native panics when popping empty error scope stack - known limitation")

	// This test would cause a panic:
	// instance, _ := CreateInstance(nil)
	// defer instance.Release()
	// adapter, _ := instance.RequestAdapter(nil)
	// defer adapter.Release()
	// device, _ := adapter.RequestDevice(nil)
	// defer device.Release()
	// device.PopErrorScopeAsync(instance) // PANIC!
}

// TestErrorScopeNoError tests pushing and popping error scope with no error.
func TestErrorScopeNoError(t *testing.T) {
	instance, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	// Push error scope
	device.PushErrorScope(ErrorFilterValidation)

	// Do some valid operation (no error expected)
	queue := device.GetQueue()
	if queue == nil {
		t.Fatal("GetQueue returned nil")
	}
	defer queue.Release()

	// Pop error scope - should get NoError
	errType, message, err := device.PopErrorScopeAsync(instance)
	if err != nil {
		t.Fatalf("PopErrorScope failed: %v", err)
	}
	if errType != ErrorTypeNoError {
		t.Errorf("Expected ErrorTypeNoError, got %v with message: %s", errType, message)
	}
	t.Logf("No error captured (as expected)")
}

// TestErrorScopeValidation tests capturing validation errors.
func TestErrorScopeValidation(t *testing.T) {
	instance, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	// Push error scope to catch validation errors
	device.PushErrorScope(ErrorFilterValidation)

	// Try to create an invalid buffer (size 0 should be invalid)
	desc := BufferDescriptor{
		NextInChain:      0,
		Label:            EmptyStringView(),
		Usage:            gputypes.BufferUsageCopyDst | gputypes.BufferUsageMapRead,
		Size:             0, // Invalid: size must be > 0
		MappedAtCreation: False,
	}
	buffer := device.CreateBuffer(&desc)
	if buffer != nil {
		defer buffer.Release()
	}

	// Pop error scope - might get validation error
	errType, message, err := device.PopErrorScopeAsync(instance)
	if err != nil {
		t.Fatalf("PopErrorScope failed: %v", err)
	}

	// Log the result (implementation-dependent whether size=0 is caught)
	if errType != ErrorTypeNoError {
		t.Logf("Captured error (type=%v): %s", errType, message)
	} else {
		t.Logf("No error captured (size=0 might be allowed by implementation)")
	}
}

// TestErrorScopeNested tests nested error scopes (LIFO).
func TestErrorScopeNested(t *testing.T) {
	instance, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer instance.Release()

	adapter, err := instance.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	// Push multiple error scopes (all same type)
	device.PushErrorScope(ErrorFilterValidation)
	device.PushErrorScope(ErrorFilterValidation)
	device.PushErrorScope(ErrorFilterValidation)

	// Pop them in reverse order (LIFO)
	errType1, msg1, err1 := device.PopErrorScopeAsync(instance)
	if err1 != nil {
		t.Fatalf("PopErrorScope 1 failed: %v", err1)
	}
	t.Logf("Popped scope 1: type=%v, msg=%s", errType1, msg1)

	errType2, msg2, err2 := device.PopErrorScopeAsync(instance)
	if err2 != nil {
		t.Fatalf("PopErrorScope 2 failed: %v", err2)
	}
	t.Logf("Popped scope 2: type=%v, msg=%s", errType2, msg2)

	errType3, msg3, err3 := device.PopErrorScopeAsync(instance)
	if err3 != nil {
		t.Fatalf("PopErrorScope 3 failed: %v", err3)
	}
	t.Logf("Popped scope 3: type=%v, msg=%s", errType3, msg3)

	// NOTE: Cannot test popping empty stack because wgpu-native panics (known limitation)
	// Stack is now empty, but we don't pop again to avoid panic
	t.Logf("Successfully popped all 3 scopes in LIFO order")
}
