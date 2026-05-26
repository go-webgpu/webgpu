package wgpu

import (
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
// Matches the gogpu/wgpu API for cross-project compatibility.
type RequestAdapterOptions struct {
	// PowerPreference indicates power consumption preference.
	PowerPreference gputypes.PowerPreference
	// ForceFallbackAdapter forces the use of a software adapter.
	ForceFallbackAdapter bool
	// CompatibleSurface, if non-nil, restricts adapter selection to those
	// compatible with rendering to the given surface.
	CompatibleSurface *Surface
}

// requestAdapterOptionsWire is the FFI-compatible C-layout struct for wgpuInstanceRequestAdapter.
// v29 layout: nextInChain(8)+featureLevel(4)+powerPreference(4)+
//   forceFallbackAdapter(4)+backendType(4)+compatibleSurface(8) = 32 bytes.
type requestAdapterOptionsWire struct {
	NextInChain          uintptr // *ChainedStruct
	FeatureLevel         FeatureLevel
	PowerPreference      gputypes.PowerPreference
	ForceFallbackAdapter Bool
	BackendType          BackendType // v29: select specific backend
	CompatibleSurface    uintptr    // WGPUSurface
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
	// adapterRequests is the global registry for pending adapter requests.
	// Protected by adapterRequestsMu for concurrent access.
	adapterRequests   = make(map[uintptr]*adapterRequest)
	adapterRequestsMu sync.Mutex
	adapterRequestID  uintptr

	// adapterCallbackPtr is the callback function pointer (created once).
	// Protected by adapterCallbackOnce for concurrent initialization.
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
		sv := (*StringView)(ptrFromUintptr(message))
		if sv.Data != 0 && sv.Length > 0 && sv.Length < 1<<20 {
			msg = unsafe.String((*byte)(ptrFromUintptr(sv.Data)), int(sv.Length))
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
			trackResource(adapter, "Adapter")
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
	if err := checkInit(); err != nil {
		return nil, err
	}
	if i == nil || i.handle == 0 {
		return nil, &WGPUError{Op: "RequestAdapter", Message: "instance is nil or released"}
	}

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

	// Convert Go-idiomatic options to wire format.
	var optionsPtr uintptr
	if options != nil {
		var surfaceHandle uintptr
		if options.CompatibleSurface != nil {
			surfaceHandle = options.CompatibleSurface.handle
		}
		wire := requestAdapterOptionsWire{
			FeatureLevel:         FeatureLevelCore,
			PowerPreference:      options.PowerPreference,
			ForceFallbackAdapter: boolToWGPU(options.ForceFallbackAdapter),
			CompatibleSurface:    surfaceHandle,
		}
		optionsPtr = uintptr(unsafe.Pointer(&wire))
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
				return nil, &WGPUError{Op: "RequestAdapter", Message: msg}
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
		untrackResource(a.handle)
		procAdapterRelease.Call(a.handle) //nolint:errcheck
		a.handle = 0
	}
}

// Limits describes resource limits for an adapter or device.
// This corresponds to WGPULimits in webgpu.h v29.
// IMPORTANT: Field order must exactly match C struct layout for correct ABI.
// v29 BREAKING changes vs v27:
//   - NextInChain added as FIRST field (was absent in v27 Limits)
//   - MinUniformBufferOffsetAlignment and MinStorageBufferOffsetAlignment moved
//     from end to after MaxStorageBufferBindingSize (before MaxVertexBuffers)
//   - MaxImmediateSize added as LAST field (new in v29)
type Limits struct {
	NextInChain                               uintptr // *ChainedStruct — NEW in v29
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
	MinUniformBufferOffsetAlignment           uint32 // MOVED: now after MaxStorageBufferBindingSize
	MinStorageBufferOffsetAlignment           uint32 // MOVED: now after MinUniformBufferOffsetAlignment
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
	MaxImmediateSize                          uint32 // NEW in v29 (push constants replacement)
}

// SupportedFeatures contains features supported by adapter or device.
// This is the wire format for wgpuAdapterGetFeatures/wgpuDeviceGetFeatures (v29 single-call API).
// Call SupportedFeaturesFreeMembers after use to release C-allocated memory.
type SupportedFeatures struct {
	FeatureCount uintptr // size_t
	Features     uintptr // *FeatureName (C-allocated, must free with SupportedFeaturesFreeMembers)
}

// AdapterInfo contains information about the adapter.
// v29: NextInChain type changed from *ChainedStructOut to *ChainedStruct.
// v29: SubgroupMinSize and SubgroupMaxSize fields added.
type AdapterInfo struct {
	NextInChain   uintptr // *ChainedStruct (was *ChainedStructOut in v27)
	Vendor        StringView
	Architecture  StringView
	Device        StringView
	Description   StringView
	BackendType   BackendType
	AdapterType   AdapterType
	VendorID      uint32
	DeviceID      uint32
	SubgroupMinSize uint32 // NEW in v29
	SubgroupMaxSize uint32 // NEW in v29
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

// Limits retrieves the limits of this adapter.
// v29: WGPULimits now has nextInChain as first field; pass *Limits directly (no SupportedLimits wrapper).
// Returns nil if the adapter is nil or if the operation fails.
func (a *Adapter) Limits() (*Limits, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if a == nil || a.handle == 0 {
		return nil, &WGPUError{Op: "Adapter.Limits", Message: "adapter is nil"}
	}

	limits := &Limits{}
	status, _, _ := procAdapterGetLimits.Call(
		a.handle,
		uintptr(unsafe.Pointer(limits)),
	)

	if WGPUStatus(status) != WGPUStatusSuccess {
		return nil, &WGPUError{Op: "Adapter.Limits", Message: "operation failed"}
	}

	return limits, nil
}

// Features retrieves all features supported by this adapter.
// v29: Uses single-call wgpuAdapterGetFeatures with SupportedFeatures struct (replaces two-call EnumerateFeatures).
// The returned slice is copied from C memory; the underlying C allocation is freed automatically.
func (a *Adapter) Features() []FeatureName {
	mustInit()
	if a == nil || a.handle == 0 {
		return nil
	}

	// Single call: wgpu fills WGPUSupportedFeatures with C-allocated array
	var sf SupportedFeatures
	procAdapterGetFeatures.Call( //nolint:errcheck
		a.handle,
		uintptr(unsafe.Pointer(&sf)),
	)

	if sf.FeatureCount == 0 || sf.Features == 0 {
		return nil
	}

	// Copy features from C memory to Go slice
	count := int(sf.FeatureCount)
	features := make([]FeatureName, count)
	for i := range features {
		// Each FeatureName is uint32 (4 bytes)
		ptr := (*FeatureName)(ptrFromUintptr(sf.Features + uintptr(i)*4))
		features[i] = *ptr
	}

	// Free C-allocated memory
	procSupportedFeaturesFreeMembers.Call(uintptr(unsafe.Pointer(&sf))) //nolint:errcheck

	return features
}

// EnumerateFeatures is a deprecated alias for Features.
// Deprecated: Use Features instead. This method was renamed in wgpu-native v29.
func (a *Adapter) EnumerateFeatures() []FeatureName {
	return a.Features()
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

// Info retrieves information about this adapter.
// The returned AdapterInfoGo contains Go strings copied from C memory.
// Returns nil if the adapter is nil or if the operation fails.
func (a *Adapter) Info() (*AdapterInfoGo, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if a == nil || a.handle == 0 {
		return nil, &WGPUError{Op: "Adapter.Info", Message: "adapter is nil"}
	}

	// Get native adapter info
	var nativeInfo AdapterInfo
	status, _, _ := procAdapterGetInfo.Call(
		a.handle,
		uintptr(unsafe.Pointer(&nativeInfo)),
	)

	if WGPUStatus(status) != WGPUStatusSuccess {
		return nil, &WGPUError{Op: "Adapter.Info", Message: "operation failed"}
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
	return unsafe.String((*byte)(ptrFromUintptr(sv.Data)), int(sv.Length))
}
