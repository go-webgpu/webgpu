package wgpu

import (
	"errors"
	"unsafe"
)

// CreateInstance creates a new WebGPU instance.
// Pass nil for default configuration.
func CreateInstance(desc *InstanceDescriptor) (*Instance, error) {
	mustInit()

	var descPtr uintptr
	if desc != nil {
		descPtr = uintptr(unsafe.Pointer(desc))
	}

	handle, _, _ := procCreateInstance.Call(descPtr)
	if handle == 0 {
		return nil, errors.New("wgpu: failed to create instance")
	}

	return &Instance{handle: handle}, nil
}

// Release releases the instance resources.
func (i *Instance) Release() {
	if i.handle != 0 {
		procInstanceRelease.Call(i.handle) //nolint:errcheck
		i.handle = 0
	}
}

// ProcessEvents processes pending async events.
func (i *Instance) ProcessEvents() {
	procInstanceProcessEvents.Call(i.handle) //nolint:errcheck
}

// ChainedStruct is used for struct chaining (input).
type ChainedStruct struct {
	Next  uintptr // *ChainedStruct
	SType uint32
}

// ChainedStructOut is used for struct chaining (output).
type ChainedStructOut struct {
	Next  uintptr // *ChainedStructOut
	SType uint32
}

// InstanceCapabilities describes instance capabilities.
// Note: This struct has specific padding requirements to match C layout.
type InstanceCapabilities struct {
	NextInChain          uintptr // *ChainedStructOut
	TimedWaitAnyEnable   Bool
	_pad                 uint32 // padding to align TimedWaitAnyMaxCount
	TimedWaitAnyMaxCount uint64
}

// InstanceDescriptor describes an Instance.
type InstanceDescriptor struct {
	NextInChain uintptr // *ChainedStruct
	Features    InstanceCapabilities
}

// Bool is a WebGPU boolean (uint32).
type Bool uint32

const (
	False Bool = 0
	True  Bool = 1
)
