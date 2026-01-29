package wgpu

import (
	"errors"
	"sync"
	"unsafe"

	"github.com/go-webgpu/goffi/ffi"
	"github.com/gogpu/gputypes"
)

// StringView represents a WebGPU string view (pointer + length).
type StringView struct {
	Data   uintptr // *char
	Length uintptr // size_t, use SIZE_MAX for null-terminated
}

// EmptyStringView returns an empty string view (null label).
func EmptyStringView() StringView {
	return StringView{Data: 0, Length: 0}
}

// Future represents an async operation handle.
type Future struct {
	ID uint64
}

// RequestAdapterOptions configures adapter selection.
type RequestAdapterOptions struct {
	NextInChain          uintptr // *ChainedStruct
	FeatureLevel         FeatureLevel
	PowerPreference      gputypes.PowerPreference
	ForceFallbackAdapter Bool
	CompatibilityMode    Bool
	CompatibleSurface    uintptr // WGPUSurface
}

// RequestAdapterCallbackInfo holds callback configuration.
type RequestAdapterCallbackInfo struct {
	NextInChain uintptr // *ChainedStruct
	Mode        CallbackMode
	Callback    uintptr // Function pointer
	Userdata1   uintptr
	Userdata2   uintptr
}

// adapterRequest holds state for an async adapter request.
type adapterRequest struct {
	done    chan struct{}
	adapter *Adapter
	status  RequestAdapterStatus
	message string
}

var (
	// Global registry for pending adapter requests
	adapterRequests   = make(map[uintptr]*adapterRequest)
	adapterRequestsMu sync.Mutex
	adapterRequestID  uintptr

	// Callback function pointer (created once)
	adapterCallbackPtr  uintptr
	adapterCallbackOnce sync.Once
)

// adapterCallbackHandler is the Go function called by C code via ffi.NewCallback.
// Windows x64 ABI: args in RCX, RDX, R8, R9, then stack.
// Signature: void(status uint32, adapter uintptr, message *StringView, userdata1 uintptr, userdata2 uintptr)
// Note: On Windows x64 ABI, structs > 8 bytes are passed by pointer.
// goffi v0.2.1+ requires all args to be uintptr and exactly one uintptr return.
func adapterCallbackHandler(status uintptr, adapter uintptr, message uintptr, userdata1, userdata2 uintptr) uintptr {
	// Extract message string (message is pointer to StringView on Windows)
	var msg string
	if message != 0 {
		// nolint:govet // message is uintptr from FFI callback - GC safe
		sv := (*StringView)(unsafe.Pointer(message))
		if sv.Data != 0 && sv.Length > 0 && sv.Length < 1<<20 {
			// nolint:govet // sv.Data is uintptr from C memory - GC safe
			msg = unsafe.String((*byte)(unsafe.Pointer(sv.Data)), int(sv.Length))
		}
	}

	// Find and complete the request
	adapterRequestsMu.Lock()
	req, ok := adapterRequests[userdata1]
	if ok {
		delete(adapterRequests, userdata1)
	}
	adapterRequestsMu.Unlock()

	if ok && req != nil {
		req.status = RequestAdapterStatus(status)
		if adapter != 0 {
			req.adapter = &Adapter{handle: adapter}
		}
		req.message = msg
		close(req.done)
	}
	return 0
}

// initAdapterCallback creates the C callback function pointer using goffi.
// goffi v0.2.1+ properly handles Windows x64 calling convention.
func initAdapterCallback() {
	adapterCallbackPtr = ffi.NewCallback(adapterCallbackHandler)
}

// RequestAdapter requests a GPU adapter from the instance.
// This is a synchronous wrapper that blocks until the adapter is available.
func (i *Instance) RequestAdapter(options *RequestAdapterOptions) (*Adapter, error) {
	mustInit()

	// Initialize callback once
	adapterCallbackOnce.Do(initAdapterCallback)

	// Create request state
	req := &adapterRequest{
		done: make(chan struct{}),
	}

	// Register request
	adapterRequestsMu.Lock()
	adapterRequestID++
	reqID := adapterRequestID
	adapterRequests[reqID] = req
	adapterRequestsMu.Unlock()

	// Prepare options
	var optionsPtr uintptr
	if options != nil {
		optionsPtr = uintptr(unsafe.Pointer(options))
	}

	// Prepare callback info
	callbackInfo := RequestAdapterCallbackInfo{
		NextInChain: 0,
		Mode:        CallbackModeAllowProcessEvents,
		Callback:    adapterCallbackPtr,
		Userdata1:   reqID,
		Userdata2:   0,
	}

	// Call wgpuInstanceRequestAdapter
	// Returns WGPUFuture (uint64) but we use callback mode
	procInstanceRequestAdapter.Call( //nolint:errcheck
		i.handle,
		optionsPtr,
		uintptr(unsafe.Pointer(&callbackInfo)),
	)

	// Process events until callback fires
	for {
		select {
		case <-req.done:
			// Callback completed
			if req.status != RequestAdapterStatusSuccess {
				msg := req.message
				if msg == "" {
					msg = "adapter request failed"
				}
				return nil, errors.New("wgpu: " + msg)
			}
			return req.adapter, nil
		default:
			// Process events to trigger callback
			i.ProcessEvents()
		}
	}
}

// Release releases the adapter resources.
func (a *Adapter) Release() {
	if a.handle != 0 {
		procAdapterRelease.Call(a.handle) //nolint:errcheck
		a.handle = 0
	}
}

// Limits describes resource limits for an adapter or device.
// This contains the most commonly used limits. WebGPU spec defines ~50 limits total.
type Limits struct {
	MaxTextureDimension1D                     uint32
	MaxTextureDimension2D                     uint32
	MaxTextureDimension3D                     uint32
	MaxTextureArrayLayers                     uint32
	MaxBindGroups                             uint32
	MaxBindGroupsPlusVertexBuffers            uint32
	MaxBindingsPerBindGroup                   uint32
	MaxDynamicUniformBuffersPerPipelineLayout uint32
	MaxDynamicStorageBuffersPerPipelineLayout uint32
	MaxSampledTexturesPerShaderStage          uint32
	MaxSamplersPerShaderStage                 uint32
	MaxStorageBuffersPerShaderStage           uint32
	MaxStorageTexturesPerShaderStage          uint32
	MaxUniformBuffersPerShaderStage           uint32
	MaxUniformBufferBindingSize               uint64
	MaxStorageBufferBindingSize               uint64
	MaxVertexBuffers                          uint32
	MaxBufferSize                             uint64
	MaxVertexAttributes                       uint32
	MaxVertexBufferArrayStride                uint32
	MaxInterStageShaderVariables              uint32
	MaxColorAttachments                       uint32
	MaxColorAttachmentBytesPerSample          uint32
	MaxComputeWorkgroupStorageSize            uint32
	MaxComputeInvocationsPerWorkgroup         uint32
	MaxComputeWorkgroupSizeX                  uint32
	MaxComputeWorkgroupSizeY                  uint32
	MaxComputeWorkgroupSizeZ                  uint32
	MaxComputeWorkgroupsPerDimension          uint32
	MinUniformBufferOffsetAlignment           uint32
	MinStorageBufferOffsetAlignment           uint32
}

// SupportedLimits contains adapter limits.
type SupportedLimits struct {
	NextInChain uintptr // *ChainedStructOut
	Limits      Limits
}

// AdapterInfo contains information about the adapter.
type AdapterInfo struct {
	NextInChain  uintptr // *ChainedStructOut
	Vendor       StringView
	Architecture StringView
	Device       StringView
	Description  StringView
	BackendType  BackendType
	AdapterType  AdapterType
	VendorID     uint32
	DeviceID     uint32
}

// AdapterInfoGo is the Go-friendly version of AdapterInfo with actual strings.
type AdapterInfoGo struct {
	Vendor       string
	Architecture string
	Device       string
	Description  string
	BackendType  BackendType
	AdapterType  AdapterType
	VendorID     uint32
	DeviceID     uint32
}

// GetLimits retrieves the limits of this adapter.
// Returns nil if the adapter is nil or if the operation fails.
func (a *Adapter) GetLimits() (*SupportedLimits, error) {
	mustInit()
	if a == nil || a.handle == 0 {
		return nil, errors.New("wgpu: adapter is nil")
	}

	limits := &SupportedLimits{}
	status, _, _ := procAdapterGetLimits.Call(
		a.handle,
		uintptr(unsafe.Pointer(limits)),
	)

	if WGPUStatus(status) != WGPUStatusSuccess {
		return nil, errors.New("wgpu: failed to get adapter limits")
	}

	return limits, nil
}

// EnumerateFeatures retrieves all features supported by this adapter.
// Returns a slice of FeatureName values.
func (a *Adapter) EnumerateFeatures() []FeatureName {
	mustInit()
	if a == nil || a.handle == 0 {
		return nil
	}

	// First call: get count
	count, _, _ := procAdapterEnumerateFeatures.Call(
		a.handle,
		0, // null pointer to get count
	)

	if count == 0 {
		return nil
	}

	// Second call: get features
	features := make([]FeatureName, count)
	procAdapterEnumerateFeatures.Call( //nolint:errcheck
		a.handle,
		uintptr(unsafe.Pointer(&features[0])),
	)

	return features
}

// HasFeature checks if the adapter supports a specific feature.
func (a *Adapter) HasFeature(feature FeatureName) bool {
	mustInit()
	if a == nil || a.handle == 0 {
		return false
	}

	result, _, _ := procAdapterHasFeature.Call(
		a.handle,
		uintptr(feature),
	)

	return Bool(result) == True
}

// GetInfo retrieves information about this adapter.
// The returned AdapterInfoGo contains Go strings copied from C memory.
// Returns nil if the adapter is nil or if the operation fails.
func (a *Adapter) GetInfo() (*AdapterInfoGo, error) {
	mustInit()
	if a == nil || a.handle == 0 {
		return nil, errors.New("wgpu: adapter is nil")
	}

	// Get native adapter info
	var nativeInfo AdapterInfo
	status, _, _ := procAdapterGetInfo.Call(
		a.handle,
		uintptr(unsafe.Pointer(&nativeInfo)),
	)

	if WGPUStatus(status) != WGPUStatusSuccess {
		return nil, errors.New("wgpu: failed to get adapter info")
	}

	// Convert StringViews to Go strings
	info := &AdapterInfoGo{
		BackendType: nativeInfo.BackendType,
		AdapterType: nativeInfo.AdapterType,
		VendorID:    nativeInfo.VendorID,
		DeviceID:    nativeInfo.DeviceID,
	}

	// Copy strings from C memory to Go memory
	if nativeInfo.Vendor.Data != 0 && nativeInfo.Vendor.Length > 0 {
		info.Vendor = stringViewToString(nativeInfo.Vendor)
	}
	if nativeInfo.Architecture.Data != 0 && nativeInfo.Architecture.Length > 0 {
		info.Architecture = stringViewToString(nativeInfo.Architecture)
	}
	if nativeInfo.Device.Data != 0 && nativeInfo.Device.Length > 0 {
		info.Device = stringViewToString(nativeInfo.Device)
	}
	if nativeInfo.Description.Data != 0 && nativeInfo.Description.Length > 0 {
		info.Description = stringViewToString(nativeInfo.Description)
	}

	// Free C memory allocated by wgpu-native
	procAdapterInfoFreeMembers.Call(uintptr(unsafe.Pointer(&nativeInfo))) //nolint:errcheck

	return info, nil
}

// stringViewToString converts a StringView to a Go string.
// This copies the data from C memory to Go memory.
func stringViewToString(sv StringView) string {
	if sv.Data == 0 || sv.Length == 0 {
		return ""
	}
	// Sanity check: limit string length to prevent issues
	if sv.Length > 1<<20 { // 1MB max
		return ""
	}
	// nolint:govet // sv.Data is uintptr from C memory - safe to convert
	return unsafe.String((*byte)(unsafe.Pointer(sv.Data)), int(sv.Length))
}
