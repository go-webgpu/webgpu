package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// InstanceDescriptor configures instance creation.
// Matches the gogpu/wgpu API for cross-project compatibility.
//
// Pass nil to CreateInstance for default configuration (all primary backends enabled).
type InstanceDescriptor struct {
	// Backends selects which GPU backends to enable.
	// Use gputypes.BackendsPrimary (default) or specific backends.
	Backends gputypes.Backends
	// Flags controls instance features like debug layers and validation.
	// Use gputypes.InstanceFlagsDebug to enable GPU debug layer.
	Flags gputypes.InstanceFlags
}

// instanceDescriptorWire is the FFI-compatible C-layout struct for wgpuCreateInstance.
// v29 layout: nextInChain(8)+requiredFeatureCount(8)+requiredFeatures(8)+requiredLimits(8) = 32 bytes.
// The v27 InstanceCapabilities/Features field is removed in v29.
type instanceDescriptorWire struct {
	NextInChain          uintptr // *ChainedStruct
	RequiredFeatureCount uintptr // size_t
	RequiredFeatures     uintptr // *InstanceFeatureName (const)
	RequiredLimits       uintptr // *InstanceLimits (const, nullable)
}

// InstanceLimits describes the limits required at instance creation.
// New in v29 — passed as RequiredLimits in instanceDescriptorWire.
type InstanceLimits struct {
	NextInChain          uintptr // *ChainedStruct (nullable)
	TimedWaitAnyMaxCount uint64
}

// Bool is a WebGPU boolean (uint32).
type Bool uint32

const (
	// False is the WebGPU boolean false value (0).
	False Bool = 0
	// True is the WebGPU boolean true value (1).
	True Bool = 1
)

// ChainedStruct is used for struct chaining (both input and output).
// In v29 ChainedStructOut was unified with ChainedStruct — use ChainedStruct everywhere.
type ChainedStruct struct {
	Next  uintptr // *ChainedStruct
	SType uint32
}

// ChainedStructOut is kept for backward compatibility.
// Deprecated: Use ChainedStruct. In v29 there is no separate ChainedStructOut in C header.
type ChainedStructOut = ChainedStruct

// CreateInstance creates a new WebGPU instance.
// Pass nil for default configuration (all primary backends enabled).
func CreateInstance(desc *InstanceDescriptor) (*Instance, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}

	// Convert Go-idiomatic descriptor to wire format.
	// When desc is nil, pass null to wgpu-native for default behavior.
	var wirePtr uintptr
	if desc != nil {
		wire := instanceDescriptorWire{} // zero = default, backends/flags handled by wgpu-native extensions
		wirePtr = uintptr(unsafe.Pointer(&wire))
	}

	handle, _, _ := procCreateInstance.Call(wirePtr)
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
