package wgpu

import (
	"unsafe"
)

// CreateInstance creates a new WebGPU instance.
// Pass nil for default configuration.
func CreateInstance(desc *InstanceDescriptor) (*Instance, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	var descPtr uintptr
	if desc != nil {
		descPtr = uintptr(unsafe.Pointer(desc))
	}

	handle, _, _ := procCreateInstance.Call(descPtr)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreateInstance", Message: "failed to create instance"}
	}

	trackResource(handle, "Instance")
	return &Instance{handle: handle}, nil
}

// Release releases the instance resources.
func (i *Instance) Release() {
	if i.handle != 0 {
		untrackResource(i.handle)
		procInstanceRelease.Call(i.handle) //nolint:errcheck
		i.handle = 0
	}
}

// ProcessEvents processes pending async events.
func (i *Instance) ProcessEvents() {
	if i == nil || i.handle == 0 {
		return
	}
	procInstanceProcessEvents.Call(i.handle) //nolint:errcheck
}

// ChainedStruct is used for struct chaining (both input and output).
// In v29 ChainedStructOut was unified with ChainedStruct — use ChainedStruct everywhere.
type ChainedStruct struct {
	Next  uintptr // *ChainedStruct
	SType uint32
}

// ChainedStructOut is kept for backward compatibility.
// Deprecated: Use ChainedStruct. In v29 there is no separate ChainedStructOut in C header.
type ChainedStructOut = ChainedStruct

// InstanceDescriptor describes a WebGPU instance.
// v29 layout: nextInChain, requiredFeatureCount, requiredFeatures, requiredLimits.
// The v27 InstanceCapabilities/Features field is removed.
type InstanceDescriptor struct {
	NextInChain         uintptr // *ChainedStruct
	RequiredFeatureCount uintptr // size_t
	RequiredFeatures    uintptr // *InstanceFeatureName (const)
	RequiredLimits      uintptr // *InstanceLimits (const, nullable)
}

// InstanceLimits describes the limits required at instance creation.
// New in v29 — passed as RequiredLimits in InstanceDescriptor.
type InstanceLimits struct {
	NextInChain           uintptr // *ChainedStruct (nullable)
	TimedWaitAnyMaxCount  uint64
}

// Bool is a WebGPU boolean (uint32).
type Bool uint32

const (
	// False is the WebGPU boolean false value (0).
	False Bool = 0
	// True is the WebGPU boolean true value (1).
	True Bool = 1
)
