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
	// NativeBackends, when non-zero, is chained to wgpuCreateInstance as
	// WGPUInstanceExtras.backends and restricts the backends wgpu-native
	// enumerates (for example InstanceBackendDX12 to pick Direct3D 12 on
	// Windows, where the default adapter may otherwise be Vulkan).
	// Zero keeps the wgpu-native default (all backends).
	NativeBackends InstanceBackend
}

// nativeDisplayHandleWire matches WGPUNativeDisplayHandle in the wgpu-native
// v29 header: type(4)+pad(4)+union{xlib,xcb,wayland}(16) = 24 bytes.
type nativeDisplayHandleWire struct {
	Type uint32
	_    uint32
	Data [2]uintptr
}

// instanceExtrasWire matches WGPUInstanceExtras in the wgpu-native v29 header
// (112 bytes on 64-bit targets).
type instanceExtrasWire struct {
	Chain                   ChainedStruct           // 0
	Backends                uint64                  // 16: WGPUInstanceBackend
	Flags                   uint64                  // 24: WGPUInstanceFlag
	Dx12ShaderCompiler      uint32                  // 32: WGPUDx12Compiler
	Gles3MinorVersion       uint32                  // 36: WGPUGles3MinorVersion
	GLFenceBehaviour        uint32                  // 40: WGPUGLFenceBehaviour
	_                       uint32                  // 44: padding
	DxcPath                 StringView              // 48
	DxcMaxShaderModel       uint32                  // 64: WGPUDxcMaxShaderModel
	Dx12PresentationSystem  uint32                  // 68: WGPUDx12SwapchainKind
	BudgetForDeviceCreation uintptr                 // 72: const uint8_t* (nullable)
	BudgetForDeviceLoss     uintptr                 // 80: const uint8_t* (nullable)
	DisplayHandle           nativeDisplayHandleWire // 88
}

// instanceExtrasSink keeps the last chained WGPUInstanceExtras on the heap.
// The descriptor refers to it only through a uintptr, so a stack-allocated
// value could be moved by stack growth before the native call reads it.
var instanceExtrasSink *instanceExtrasWire

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
		wire := instanceDescriptorWire{} // zero = default
		if desc.NativeBackends != 0 {
			extras := &instanceExtrasWire{}
			instanceExtrasSink = extras
			extras.Chain.SType = uint32(STypeInstanceExtras)
			extras.Backends = uint64(desc.NativeBackends)
			wire.NextInChain = uintptr(unsafe.Pointer(extras))
		}
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
