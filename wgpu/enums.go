package wgpu

// RequestAdapterStatus is the status returned by RequestAdapter callback.
type RequestAdapterStatus uint32

const (
	RequestAdapterStatusSuccess         RequestAdapterStatus = 0x00000001
	RequestAdapterStatusInstanceDropped RequestAdapterStatus = 0x00000002
	RequestAdapterStatusUnavailable     RequestAdapterStatus = 0x00000003
	RequestAdapterStatusError           RequestAdapterStatus = 0x00000004
)

// RequestDeviceStatus is the status returned by RequestDevice callback.
type RequestDeviceStatus uint32

const (
	RequestDeviceStatusSuccess         RequestDeviceStatus = 0x00000001
	RequestDeviceStatusInstanceDropped RequestDeviceStatus = 0x00000002
	RequestDeviceStatusError           RequestDeviceStatus = 0x00000003
	RequestDeviceStatusUnknown         RequestDeviceStatus = 0x00000004
)

// FeatureLevel indicates the WebGPU feature level.
type FeatureLevel uint32

const (
	FeatureLevelCompatibility FeatureLevel = 0x00000001
	FeatureLevelCore          FeatureLevel = 0x00000002
)

// CallbackMode controls how callbacks are fired.
type CallbackMode uint32

const (
	CallbackModeWaitAnyOnly        CallbackMode = 0x00000001
	CallbackModeAllowProcessEvents CallbackMode = 0x00000002
	CallbackModeAllowSpontaneous   CallbackMode = 0x00000003
)

// SType identifies chained struct types.
type SType uint32

const (
	// Standard WebGPU STypes
	STypeShaderSourceSPIRV SType = 0x00000001
	STypeShaderSourceWGSL  SType = 0x00000002

	// Surface source STypes
	STypeSurfaceSourceMetalLayer          SType = 0x00000004
	STypeSurfaceSourceWindowsHWND         SType = 0x00000005
	STypeSurfaceSourceXlibWindow          SType = 0x00000006
	STypeSurfaceSourceWaylandSurface      SType = 0x00000007
	STypeSurfaceSourceAndroidNativeWindow SType = 0x00000008
	STypeSurfaceSourceXCBWindow           SType = 0x00000009

	// Native extension STypes (from wgpu.h)
	STypeInstanceExtras SType = 0x00030006
)

// SurfaceGetCurrentTextureStatus describes the result of GetCurrentTexture.
type SurfaceGetCurrentTextureStatus uint32

const (
	SurfaceGetCurrentTextureStatusSuccessOptimal    SurfaceGetCurrentTextureStatus = 0x01
	SurfaceGetCurrentTextureStatusSuccessSuboptimal SurfaceGetCurrentTextureStatus = 0x02
	SurfaceGetCurrentTextureStatusTimeout           SurfaceGetCurrentTextureStatus = 0x03
	SurfaceGetCurrentTextureStatusOutdated          SurfaceGetCurrentTextureStatus = 0x04
	SurfaceGetCurrentTextureStatusLost              SurfaceGetCurrentTextureStatus = 0x05
	SurfaceGetCurrentTextureStatusOutOfMemory       SurfaceGetCurrentTextureStatus = 0x06
	SurfaceGetCurrentTextureStatusDeviceLost        SurfaceGetCurrentTextureStatus = 0x07
	SurfaceGetCurrentTextureStatusError             SurfaceGetCurrentTextureStatus = 0x08
)

// TextureAspect describes which aspect of a texture to access.
type TextureAspect uint32

const (
	TextureAspectUndefined   TextureAspect = 0x00
	TextureAspectAll         TextureAspect = 0x01
	TextureAspectStencilOnly TextureAspect = 0x02
	TextureAspectDepthOnly   TextureAspect = 0x03
)

// OptionalBool is a tri-state boolean for WebGPU.
type OptionalBool uint32

const (
	OptionalBoolFalse     OptionalBool = 0x00000000
	OptionalBoolTrue      OptionalBool = 0x00000001
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
	WGPUStatusSuccess WGPUStatus = 0x00000000
	WGPUStatusError   WGPUStatus = 0x00000001
)

// BufferMapState describes the mapping state of a buffer.
type BufferMapState uint32

const (
	BufferMapStateUnmapped BufferMapState = 0x00000001
	BufferMapStatePending  BufferMapState = 0x00000002
	BufferMapStateMapped   BufferMapState = 0x00000003
)

// BackendType describes the graphics backend being used.
type BackendType uint32

const (
	BackendTypeUndefined BackendType = 0x00000000
	BackendTypeNull      BackendType = 0x00000001
	BackendTypeWebGPU    BackendType = 0x00000002
	BackendTypeD3D11     BackendType = 0x00000003
	BackendTypeD3D12     BackendType = 0x00000004
	BackendTypeMetal     BackendType = 0x00000005
	BackendTypeVulkan    BackendType = 0x00000006
	BackendTypeOpenGL    BackendType = 0x00000007
	BackendTypeOpenGLES  BackendType = 0x00000008
)

// AdapterType describes the type of GPU adapter.
type AdapterType uint32

const (
	AdapterTypeDiscreteGPU   AdapterType = 0x00000001
	AdapterTypeIntegratedGPU AdapterType = 0x00000002
	AdapterTypeCPU           AdapterType = 0x00000003
	AdapterTypeUnknown       AdapterType = 0x00000004
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
