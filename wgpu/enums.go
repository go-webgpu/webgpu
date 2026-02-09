package wgpu

// RequestAdapterStatus is the status returned by RequestAdapter callback.
type RequestAdapterStatus uint32

const (
	// RequestAdapterStatusSuccess indicates the adapter was successfully obtained.
	RequestAdapterStatusSuccess RequestAdapterStatus = 0x00000001
	// RequestAdapterStatusInstanceDropped indicates the instance was dropped before completion.
	RequestAdapterStatusInstanceDropped RequestAdapterStatus = 0x00000002
	// RequestAdapterStatusUnavailable indicates no suitable adapter is available.
	RequestAdapterStatusUnavailable RequestAdapterStatus = 0x00000003
	// RequestAdapterStatusError indicates an error occurred during adapter request.
	RequestAdapterStatusError RequestAdapterStatus = 0x00000004
)

// RequestDeviceStatus is the status returned by RequestDevice callback.
type RequestDeviceStatus uint32

const (
	// RequestDeviceStatusSuccess indicates the device was successfully obtained.
	RequestDeviceStatusSuccess RequestDeviceStatus = 0x00000001
	// RequestDeviceStatusInstanceDropped indicates the instance was dropped before completion.
	RequestDeviceStatusInstanceDropped RequestDeviceStatus = 0x00000002
	// RequestDeviceStatusError indicates an error occurred during device request.
	RequestDeviceStatusError RequestDeviceStatus = 0x00000003
	// RequestDeviceStatusUnknown indicates an unknown error occurred.
	RequestDeviceStatusUnknown RequestDeviceStatus = 0x00000004
)

// FeatureLevel indicates the WebGPU feature level.
type FeatureLevel uint32

const (
	// FeatureLevelCompatibility indicates the compatibility feature level (WebGPU compat).
	FeatureLevelCompatibility FeatureLevel = 0x00000001
	// FeatureLevelCore indicates the core feature level (full WebGPU).
	FeatureLevelCore FeatureLevel = 0x00000002
)

// CallbackMode controls how callbacks are fired.
type CallbackMode uint32

const (
	// CallbackModeWaitAnyOnly fires callbacks only during WaitAny calls.
	CallbackModeWaitAnyOnly CallbackMode = 0x00000001
	// CallbackModeAllowProcessEvents fires callbacks during ProcessEvents calls.
	CallbackModeAllowProcessEvents CallbackMode = 0x00000002
	// CallbackModeAllowSpontaneous allows callbacks to fire at any time.
	CallbackModeAllowSpontaneous CallbackMode = 0x00000003
)

// SType identifies chained struct types.
type SType uint32

const (
	// STypeShaderSourceSPIRV identifies a SPIR-V shader source chained struct.
	STypeShaderSourceSPIRV SType = 0x00000001
	// STypeShaderSourceWGSL identifies a WGSL shader source chained struct.
	STypeShaderSourceWGSL SType = 0x00000002

	// STypeSurfaceSourceMetalLayer identifies a Metal layer surface source (macOS/iOS).
	STypeSurfaceSourceMetalLayer SType = 0x00000004
	// STypeSurfaceSourceWindowsHWND identifies a Windows HWND surface source.
	STypeSurfaceSourceWindowsHWND SType = 0x00000005
	// STypeSurfaceSourceXlibWindow identifies an Xlib window surface source (Linux).
	STypeSurfaceSourceXlibWindow SType = 0x00000006
	// STypeSurfaceSourceWaylandSurface identifies a Wayland surface source (Linux).
	STypeSurfaceSourceWaylandSurface SType = 0x00000007
	// STypeSurfaceSourceAndroidNativeWindow identifies an Android native window surface source.
	STypeSurfaceSourceAndroidNativeWindow SType = 0x00000008
	// STypeSurfaceSourceXCBWindow identifies an XCB window surface source (Linux).
	STypeSurfaceSourceXCBWindow SType = 0x00000009

	// STypeInstanceExtras identifies wgpu-native instance extras chained struct.
	STypeInstanceExtras SType = 0x00030006
)

// SurfaceGetCurrentTextureStatus describes the result of GetCurrentTexture.
type SurfaceGetCurrentTextureStatus uint32

const (
	// SurfaceGetCurrentTextureStatusSuccessOptimal indicates the texture was obtained optimally.
	SurfaceGetCurrentTextureStatusSuccessOptimal SurfaceGetCurrentTextureStatus = 0x01
	// SurfaceGetCurrentTextureStatusSuccessSuboptimal indicates the texture was obtained but may not be optimal.
	SurfaceGetCurrentTextureStatusSuccessSuboptimal SurfaceGetCurrentTextureStatus = 0x02
	// SurfaceGetCurrentTextureStatusTimeout indicates the operation timed out.
	SurfaceGetCurrentTextureStatusTimeout SurfaceGetCurrentTextureStatus = 0x03
	// SurfaceGetCurrentTextureStatusOutdated indicates the surface needs reconfiguration.
	SurfaceGetCurrentTextureStatusOutdated SurfaceGetCurrentTextureStatus = 0x04
	// SurfaceGetCurrentTextureStatusLost indicates the surface was lost and must be recreated.
	SurfaceGetCurrentTextureStatusLost SurfaceGetCurrentTextureStatus = 0x05
	// SurfaceGetCurrentTextureStatusOutOfMemory indicates GPU memory allocation failed.
	SurfaceGetCurrentTextureStatusOutOfMemory SurfaceGetCurrentTextureStatus = 0x06
	// SurfaceGetCurrentTextureStatusDeviceLost indicates the GPU device was lost.
	SurfaceGetCurrentTextureStatusDeviceLost SurfaceGetCurrentTextureStatus = 0x07
	// SurfaceGetCurrentTextureStatusError indicates an unspecified error occurred.
	SurfaceGetCurrentTextureStatusError SurfaceGetCurrentTextureStatus = 0x08
)

// TextureAspect describes which aspect of a texture to access.
type TextureAspect uint32

const (
	// TextureAspectUndefined leaves the texture aspect unspecified.
	TextureAspectUndefined TextureAspect = 0x00
	// TextureAspectAll accesses all aspects of the texture.
	TextureAspectAll TextureAspect = 0x01
	// TextureAspectStencilOnly accesses only the stencil aspect of a depth-stencil texture.
	TextureAspectStencilOnly TextureAspect = 0x02
	// TextureAspectDepthOnly accesses only the depth aspect of a depth-stencil texture.
	TextureAspectDepthOnly TextureAspect = 0x03
)

// OptionalBool is a tri-state boolean for WebGPU.
type OptionalBool uint32

const (
	// OptionalBoolFalse represents an explicit false value.
	OptionalBoolFalse OptionalBool = 0x00000000
	// OptionalBoolTrue represents an explicit true value.
	OptionalBoolTrue OptionalBool = 0x00000001
	// OptionalBoolUndefined represents an unset/default value.
	OptionalBoolUndefined OptionalBool = 0x00000002
)

// QueryType describes the type of queries in a QuerySet.
type QueryType uint32

const (
	// QueryTypeOcclusion specifies occlusion queries.
	QueryTypeOcclusion QueryType = 0x00000001
	// QueryTypeTimestamp specifies timestamp queries for GPU profiling.
	QueryTypeTimestamp QueryType = 0x00000002
)

// FeatureName describes a WebGPU feature that can be requested.
type FeatureName uint32

const (
	// FeatureNameTimestampQuery enables timestamp query support.
	FeatureNameTimestampQuery FeatureName = 0x00000003
	// Add more features as needed
)

// WGPUStatus describes the status returned from certain WebGPU operations.
type WGPUStatus uint32

const (
	// WGPUStatusSuccess indicates the operation completed successfully.
	WGPUStatusSuccess WGPUStatus = 0x00000000
	// WGPUStatusError indicates the operation failed.
	WGPUStatusError WGPUStatus = 0x00000001
)

// BufferMapState describes the mapping state of a buffer.
type BufferMapState uint32

const (
	// BufferMapStateUnmapped indicates the buffer is not mapped.
	BufferMapStateUnmapped BufferMapState = 0x00000001
	// BufferMapStatePending indicates a map operation is in progress.
	BufferMapStatePending BufferMapState = 0x00000002
	// BufferMapStateMapped indicates the buffer is currently mapped and accessible.
	BufferMapStateMapped BufferMapState = 0x00000003
)

// BackendType describes the graphics backend being used.
type BackendType uint32

const (
	// BackendTypeUndefined indicates no backend is specified.
	BackendTypeUndefined BackendType = 0x00000000
	// BackendTypeNull indicates a null/mock backend (for testing).
	BackendTypeNull BackendType = 0x00000001
	// BackendTypeWebGPU indicates the native WebGPU backend.
	BackendTypeWebGPU BackendType = 0x00000002
	// BackendTypeD3D11 indicates the Direct3D 11 backend (Windows).
	BackendTypeD3D11 BackendType = 0x00000003
	// BackendTypeD3D12 indicates the Direct3D 12 backend (Windows).
	BackendTypeD3D12 BackendType = 0x00000004
	// BackendTypeMetal indicates the Metal backend (macOS/iOS).
	BackendTypeMetal BackendType = 0x00000005
	// BackendTypeVulkan indicates the Vulkan backend (cross-platform).
	BackendTypeVulkan BackendType = 0x00000006
	// BackendTypeOpenGL indicates the OpenGL backend.
	BackendTypeOpenGL BackendType = 0x00000007
	// BackendTypeOpenGLES indicates the OpenGL ES backend (mobile/embedded).
	BackendTypeOpenGLES BackendType = 0x00000008
)

// AdapterType describes the type of GPU adapter.
type AdapterType uint32

const (
	// AdapterTypeDiscreteGPU indicates a dedicated/discrete GPU (best performance).
	AdapterTypeDiscreteGPU AdapterType = 0x00000001
	// AdapterTypeIntegratedGPU indicates an integrated GPU (shared with CPU).
	AdapterTypeIntegratedGPU AdapterType = 0x00000002
	// AdapterTypeCPU indicates a software/CPU-based renderer.
	AdapterTypeCPU AdapterType = 0x00000003
	// AdapterTypeUnknown indicates the adapter type could not be determined.
	AdapterTypeUnknown AdapterType = 0x00000004
)

// ErrorFilter filters error types in error scopes.
type ErrorFilter uint32

const (
	// ErrorFilterValidation catches validation errors.
	ErrorFilterValidation ErrorFilter = 0x00000001
	// ErrorFilterOutOfMemory catches out-of-memory errors.
	ErrorFilterOutOfMemory ErrorFilter = 0x00000002
	// ErrorFilterInternal catches internal errors.
	ErrorFilterInternal ErrorFilter = 0x00000003
)

// ErrorType describes the type of error that occurred.
type ErrorType uint32

const (
	// ErrorTypeNoError indicates no error occurred.
	ErrorTypeNoError ErrorType = 0x00000001
	// ErrorTypeValidation indicates a validation error.
	ErrorTypeValidation ErrorType = 0x00000002
	// ErrorTypeOutOfMemory indicates an out-of-memory error.
	ErrorTypeOutOfMemory ErrorType = 0x00000003
	// ErrorTypeInternal indicates an internal error.
	ErrorTypeInternal ErrorType = 0x00000004
	// ErrorTypeUnknown indicates an unknown error.
	ErrorTypeUnknown ErrorType = 0x00000005
)

// PopErrorScopeStatus describes the result of PopErrorScope operation.
type PopErrorScopeStatus uint32

const (
	// PopErrorScopeStatusSuccess indicates the error scope was successfully popped.
	PopErrorScopeStatusSuccess PopErrorScopeStatus = 0x00000001
	// PopErrorScopeStatusInstanceDropped indicates the instance was dropped.
	PopErrorScopeStatusInstanceDropped PopErrorScopeStatus = 0x00000002
	// PopErrorScopeStatusEmptyStack indicates the error scope stack was empty.
	PopErrorScopeStatusEmptyStack PopErrorScopeStatus = 0x00000003
)

// DeviceLostReason describes why a device was lost.
type DeviceLostReason uint32

const (
	// DeviceLostReasonUnknown indicates the device was lost for an unknown reason.
	DeviceLostReasonUnknown DeviceLostReason = 0x00000001
	// DeviceLostReasonDestroyed indicates the device was explicitly destroyed.
	DeviceLostReasonDestroyed DeviceLostReason = 0x00000002
	// DeviceLostReasonInstanceDropped indicates the instance was dropped.
	DeviceLostReasonInstanceDropped DeviceLostReason = 0x00000003
	// DeviceLostReasonFailedCreation indicates device creation failed.
	DeviceLostReasonFailedCreation DeviceLostReason = 0x00000004
)
