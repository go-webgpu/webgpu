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
//
//	forceFallbackAdapter(4)+backendType(4)+compatibleSurface(8) = 32 bytes.
type requestAdapterOptionsWire struct {
	NextInChain          uintptr // *ChainedStruct
	FeatureLevel         FeatureLevel
	PowerPreference      gputypes.PowerPreference
	ForceFallbackAdapter Bool
	BackendType          BackendType // v29: select specific backend
	CompatibleSurface    uintptr     // WGPUSurface
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

// handleAdapterCallback completes a request after the platform callback entry
// normalizes the ABI-specific WGPUStringView representation.
// userdata2 is reserved by WebGPU and discarded by the platform entry.
func handleAdapterCallback(status uintptr, adapter uintptr, message StringView, userdata1 uintptr) uintptr {
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
		req.message = stringViewToString(message)
		close(req.done)
	}
	return 0
}

// initAdapterCallback creates the platform-correct C callback function pointer.
func initAdapterCallback() {
	adapterCallbackPtr = ffi.NewCallback(adapterCallbackEntry)
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
			// Cache limits at creation time so Limits() returns value without FFI.
			if req.adapter != nil {
				req.adapter.limits = fetchAdapterLimits(req.adapter.handle)
			}
			return req.adapter, nil
		default:
			// Process events to trigger callback
			i.ProcessEvents()
		}
	}
}

// fetchAdapterLimits calls wgpuAdapterGetLimits and converts the wire struct to public Limits.
// Returns zero-value Limits on failure (non-fatal: limits remain valid defaults).
func fetchAdapterLimits(handle uintptr) Limits {
	var wire limitsWire
	status, _, _ := procAdapterGetLimits.Call(
		handle,
		uintptr(unsafe.Pointer(&wire)),
	)
	if WGPUStatus(status) != WGPUStatusSuccess {
		return Limits{}
	}
	return limitsFromWire(&wire)
}

// Release releases the adapter resources.
func (a *Adapter) Release() {
	if a.handle != 0 {
		untrackResource(a.handle)
		procAdapterRelease.Call(a.handle) //nolint:errcheck
		a.handle = 0
	}
}

// Limits describes GPU resource limits for an adapter or device.
//
// This type matches the gogpu/wgpu API for cross-project compatibility.
// Limits are cached at creation time (RequestAdapter/RequestDevice) and
// returned by value — no FFI call is made on each access.
//
// Note: wgpu-native-specific fields (MaxImmediateSize, MaxNonSamplerBindings)
// are not exposed here. Use NativeLimits for native extensions.
type Limits struct {
	// MaxTextureDimension1D is the maximum 1D texture dimension.
	MaxTextureDimension1D uint32
	// MaxTextureDimension2D is the maximum 2D texture dimension.
	MaxTextureDimension2D uint32
	// MaxTextureDimension3D is the maximum 3D texture dimension.
	MaxTextureDimension3D uint32
	// MaxTextureArrayLayers is the maximum texture array layer count.
	MaxTextureArrayLayers uint32
	// MaxBindGroups is the maximum number of bind groups.
	MaxBindGroups uint32
	// MaxBindGroupsPlusVertexBuffers is the max bind groups + vertex buffers.
	MaxBindGroupsPlusVertexBuffers uint32
	// MaxBindingsPerBindGroup is the max bindings per bind group.
	MaxBindingsPerBindGroup uint32
	// MaxDynamicUniformBuffersPerPipelineLayout is the max dynamic uniform buffers.
	MaxDynamicUniformBuffersPerPipelineLayout uint32
	// MaxDynamicStorageBuffersPerPipelineLayout is the max dynamic storage buffers.
	MaxDynamicStorageBuffersPerPipelineLayout uint32
	// MaxSampledTexturesPerShaderStage is the max sampled textures per shader stage.
	MaxSampledTexturesPerShaderStage uint32
	// MaxSamplersPerShaderStage is the max samplers per shader stage.
	MaxSamplersPerShaderStage uint32
	// MaxStorageBuffersPerShaderStage is the max storage buffers per shader stage.
	MaxStorageBuffersPerShaderStage uint32
	// MaxStorageTexturesPerShaderStage is the max storage textures per shader stage.
	MaxStorageTexturesPerShaderStage uint32
	// MaxUniformBuffersPerShaderStage is the max uniform buffers per shader stage.
	MaxUniformBuffersPerShaderStage uint32
	// MaxUniformBufferBindingSize is the max uniform buffer binding size in bytes.
	MaxUniformBufferBindingSize uint64
	// MaxStorageBufferBindingSize is the max storage buffer binding size in bytes.
	MaxStorageBufferBindingSize uint64
	// MinUniformBufferOffsetAlignment is the minimum uniform buffer offset alignment.
	MinUniformBufferOffsetAlignment uint32
	// MinStorageBufferOffsetAlignment is the minimum storage buffer offset alignment.
	MinStorageBufferOffsetAlignment uint32
	// MaxVertexBuffers is the max vertex buffers in a pipeline.
	MaxVertexBuffers uint32
	// MaxBufferSize is the max buffer size in bytes.
	MaxBufferSize uint64
	// MaxVertexAttributes is the max vertex attributes in a pipeline.
	MaxVertexAttributes uint32
	// MaxVertexBufferArrayStride is the max vertex buffer array stride.
	MaxVertexBufferArrayStride uint32
	// MaxInterStageShaderVariables is the max inter-stage shader variables.
	MaxInterStageShaderVariables uint32
	// MaxColorAttachments is the max color attachments in a render pass.
	MaxColorAttachments uint32
	// MaxColorAttachmentBytesPerSample is the max bytes per sample for color attachments.
	MaxColorAttachmentBytesPerSample uint32
	// MaxComputeWorkgroupStorageSize is the max compute workgroup storage in bytes.
	MaxComputeWorkgroupStorageSize uint32
	// MaxComputeInvocationsPerWorkgroup is the max compute invocations per workgroup.
	MaxComputeInvocationsPerWorkgroup uint32
	// MaxComputeWorkgroupSizeX is the max compute workgroup size in X dimension.
	MaxComputeWorkgroupSizeX uint32
	// MaxComputeWorkgroupSizeY is the max compute workgroup size in Y dimension.
	MaxComputeWorkgroupSizeY uint32
	// MaxComputeWorkgroupSizeZ is the max compute workgroup size in Z dimension.
	MaxComputeWorkgroupSizeZ uint32
	// MaxComputeWorkgroupsPerDimension is the max compute workgroups per dimension.
	MaxComputeWorkgroupsPerDimension uint32
}

// limitsWire is the FFI-compatible C-layout struct for wgpu-native v29 WGPULimits.
// IMPORTANT: Field order must exactly match C struct layout for correct ABI.
// v29 BREAKING changes vs v27:
//   - NextInChain added as FIRST field (was absent in v27 Limits)
//   - MinUniformBufferOffsetAlignment and MinStorageBufferOffsetAlignment moved
//     from end to after MaxStorageBufferBindingSize (before MaxVertexBuffers)
//   - MaxImmediateSize added as LAST field (new in v29)
type limitsWire struct {
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

// limitsFromWire converts a limitsWire (FFI struct with NextInChain) to public Limits.
func limitsFromWire(w *limitsWire) Limits {
	return Limits{
		MaxTextureDimension1D:                     w.MaxTextureDimension1D,
		MaxTextureDimension2D:                     w.MaxTextureDimension2D,
		MaxTextureDimension3D:                     w.MaxTextureDimension3D,
		MaxTextureArrayLayers:                     w.MaxTextureArrayLayers,
		MaxBindGroups:                             w.MaxBindGroups,
		MaxBindGroupsPlusVertexBuffers:            w.MaxBindGroupsPlusVertexBuffers,
		MaxBindingsPerBindGroup:                   w.MaxBindingsPerBindGroup,
		MaxDynamicUniformBuffersPerPipelineLayout: w.MaxDynamicUniformBuffersPerPipelineLayout,
		MaxDynamicStorageBuffersPerPipelineLayout: w.MaxDynamicStorageBuffersPerPipelineLayout,
		MaxSampledTexturesPerShaderStage:          w.MaxSampledTexturesPerShaderStage,
		MaxSamplersPerShaderStage:                 w.MaxSamplersPerShaderStage,
		MaxStorageBuffersPerShaderStage:           w.MaxStorageBuffersPerShaderStage,
		MaxStorageTexturesPerShaderStage:          w.MaxStorageTexturesPerShaderStage,
		MaxUniformBuffersPerShaderStage:           w.MaxUniformBuffersPerShaderStage,
		MaxUniformBufferBindingSize:               w.MaxUniformBufferBindingSize,
		MaxStorageBufferBindingSize:               w.MaxStorageBufferBindingSize,
		MinUniformBufferOffsetAlignment:           w.MinUniformBufferOffsetAlignment,
		MinStorageBufferOffsetAlignment:           w.MinStorageBufferOffsetAlignment,
		MaxVertexBuffers:                          w.MaxVertexBuffers,
		MaxBufferSize:                             w.MaxBufferSize,
		MaxVertexAttributes:                       w.MaxVertexAttributes,
		MaxVertexBufferArrayStride:                w.MaxVertexBufferArrayStride,
		MaxInterStageShaderVariables:              w.MaxInterStageShaderVariables,
		MaxColorAttachments:                       w.MaxColorAttachments,
		MaxColorAttachmentBytesPerSample:          w.MaxColorAttachmentBytesPerSample,
		MaxComputeWorkgroupStorageSize:            w.MaxComputeWorkgroupStorageSize,
		MaxComputeInvocationsPerWorkgroup:         w.MaxComputeInvocationsPerWorkgroup,
		MaxComputeWorkgroupSizeX:                  w.MaxComputeWorkgroupSizeX,
		MaxComputeWorkgroupSizeY:                  w.MaxComputeWorkgroupSizeY,
		MaxComputeWorkgroupSizeZ:                  w.MaxComputeWorkgroupSizeZ,
		MaxComputeWorkgroupsPerDimension:          w.MaxComputeWorkgroupsPerDimension,
	}
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
	NextInChain     uintptr // *ChainedStruct (was *ChainedStructOut in v27)
	Vendor          StringView
	Architecture    StringView
	Device          StringView
	Description     StringView
	BackendType     BackendType
	AdapterType     AdapterType
	VendorID        uint32
	DeviceID        uint32
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

// Limits returns the resource limits of this adapter.
//
// Limits are cached at adapter creation time and returned by value.
// No FFI call is made. Returns zero-value Limits if the adapter is nil.
// This matches the gogpu/wgpu API signature for cross-project compatibility.
func (a *Adapter) Limits() Limits {
	if a == nil || a.handle == 0 {
		return Limits{}
	}
	return a.limits
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
