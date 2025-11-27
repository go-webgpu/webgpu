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

// PowerPreference indicates preference for GPU power usage.
type PowerPreference uint32

const (
	PowerPreferenceUndefined       PowerPreference = 0x00000000
	PowerPreferenceLowPower        PowerPreference = 0x00000001
	PowerPreferenceHighPerformance PowerPreference = 0x00000002
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

// TextureUsage describes how a texture will be used.
type TextureUsage uint64

const (
	TextureUsageNone             TextureUsage = 0x00
	TextureUsageCopySrc          TextureUsage = 0x01
	TextureUsageCopyDst          TextureUsage = 0x02
	TextureUsageTextureBinding   TextureUsage = 0x04
	TextureUsageStorageBinding   TextureUsage = 0x08
	TextureUsageRenderAttachment TextureUsage = 0x10
)

// PresentMode describes how frames are presented to the surface.
type PresentMode uint32

const (
	PresentModeUndefined   PresentMode = 0x00
	PresentModeFifo        PresentMode = 0x01 // VSync, always available
	PresentModeFifoRelaxed PresentMode = 0x02
	PresentModeImmediate   PresentMode = 0x03 // No VSync, may tear
	PresentModeMailbox     PresentMode = 0x04 // Triple buffering
)

// CompositeAlphaMode describes how alpha is handled for surface compositing.
type CompositeAlphaMode uint32

const (
	CompositeAlphaModeAuto            CompositeAlphaMode = 0x00
	CompositeAlphaModeOpaque          CompositeAlphaMode = 0x01
	CompositeAlphaModePremultiplied   CompositeAlphaMode = 0x02
	CompositeAlphaModeUnpremultiplied CompositeAlphaMode = 0x03
	CompositeAlphaModeInherit         CompositeAlphaMode = 0x04
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

// LoadOp describes what happens to render target at the beginning of a pass.
type LoadOp uint32

const (
	LoadOpUndefined LoadOp = 0x00
	LoadOpLoad      LoadOp = 0x01 // Keep existing content
	LoadOpClear     LoadOp = 0x02 // Clear to clear value
)

// StoreOp describes what happens to render target at the end of a pass.
type StoreOp uint32

const (
	StoreOpUndefined StoreOp = 0x00
	StoreOpStore     StoreOp = 0x01 // Write results to texture
	StoreOpDiscard   StoreOp = 0x02 // Discard (e.g., for depth after use)
)

// TextureAspect describes which aspect of a texture to access.
type TextureAspect uint32

const (
	TextureAspectUndefined   TextureAspect = 0x00
	TextureAspectAll         TextureAspect = 0x01
	TextureAspectStencilOnly TextureAspect = 0x02
	TextureAspectDepthOnly   TextureAspect = 0x03
)

// PrimitiveTopology describes how vertices form primitives.
type PrimitiveTopology uint32

const (
	PrimitiveTopologyUndefined     PrimitiveTopology = 0x00
	PrimitiveTopologyPointList     PrimitiveTopology = 0x01
	PrimitiveTopologyLineList      PrimitiveTopology = 0x02
	PrimitiveTopologyLineStrip     PrimitiveTopology = 0x03
	PrimitiveTopologyTriangleList  PrimitiveTopology = 0x04
	PrimitiveTopologyTriangleStrip PrimitiveTopology = 0x05
)

// IndexFormat describes the format of index buffer data.
type IndexFormat uint32

const (
	IndexFormatUndefined IndexFormat = 0x00
	IndexFormatUint16    IndexFormat = 0x01
	IndexFormatUint32    IndexFormat = 0x02
)

// FrontFace describes which winding order is considered front-facing.
type FrontFace uint32

const (
	FrontFaceUndefined FrontFace = 0x00
	FrontFaceCCW       FrontFace = 0x01 // Counter-clockwise
	FrontFaceCW        FrontFace = 0x02 // Clockwise
)

// CullMode describes which faces to cull.
type CullMode uint32

const (
	CullModeUndefined CullMode = 0x00
	CullModeNone      CullMode = 0x01
	CullModeFront     CullMode = 0x02
	CullModeBack      CullMode = 0x03
)

// VertexFormat describes the format of vertex attribute data.
type VertexFormat uint32

const (
	VertexFormatUint8     VertexFormat = 0x01
	VertexFormatUint8x2   VertexFormat = 0x02
	VertexFormatUint8x4   VertexFormat = 0x03
	VertexFormatSint8     VertexFormat = 0x04
	VertexFormatSint8x2   VertexFormat = 0x05
	VertexFormatSint8x4   VertexFormat = 0x06
	VertexFormatUnorm8    VertexFormat = 0x07
	VertexFormatUnorm8x2  VertexFormat = 0x08
	VertexFormatUnorm8x4  VertexFormat = 0x09
	VertexFormatSnorm8    VertexFormat = 0x0A
	VertexFormatSnorm8x2  VertexFormat = 0x0B
	VertexFormatSnorm8x4  VertexFormat = 0x0C
	VertexFormatUint16    VertexFormat = 0x0D
	VertexFormatUint16x2  VertexFormat = 0x0E
	VertexFormatUint16x4  VertexFormat = 0x0F
	VertexFormatSint16    VertexFormat = 0x10
	VertexFormatSint16x2  VertexFormat = 0x11
	VertexFormatSint16x4  VertexFormat = 0x12
	VertexFormatUnorm16   VertexFormat = 0x13
	VertexFormatUnorm16x2 VertexFormat = 0x14
	VertexFormatUnorm16x4 VertexFormat = 0x15
	VertexFormatSnorm16   VertexFormat = 0x16
	VertexFormatSnorm16x2 VertexFormat = 0x17
	VertexFormatSnorm16x4 VertexFormat = 0x18
	VertexFormatFloat16   VertexFormat = 0x19
	VertexFormatFloat16x2 VertexFormat = 0x1A
	VertexFormatFloat16x4 VertexFormat = 0x1B
	VertexFormatFloat32   VertexFormat = 0x1C
	VertexFormatFloat32x2 VertexFormat = 0x1D
	VertexFormatFloat32x3 VertexFormat = 0x1E
	VertexFormatFloat32x4 VertexFormat = 0x1F
	VertexFormatUint32    VertexFormat = 0x20
	VertexFormatUint32x2  VertexFormat = 0x21
	VertexFormatUint32x3  VertexFormat = 0x22
	VertexFormatUint32x4  VertexFormat = 0x23
	VertexFormatSint32    VertexFormat = 0x24
	VertexFormatSint32x2  VertexFormat = 0x25
	VertexFormatSint32x3  VertexFormat = 0x26
	VertexFormatSint32x4  VertexFormat = 0x27
)

// VertexStepMode describes how vertex buffer is advanced.
type VertexStepMode uint32

const (
	VertexStepModeUndefined           VertexStepMode = 0x00
	VertexStepModeVertexBufferNotUsed VertexStepMode = 0x01
	VertexStepModeVertex              VertexStepMode = 0x02
	VertexStepModeInstance            VertexStepMode = 0x03
)

// ColorWriteMask describes which color channels to write.
type ColorWriteMask uint32

const (
	ColorWriteMaskNone  ColorWriteMask = 0x00
	ColorWriteMaskRed   ColorWriteMask = 0x01
	ColorWriteMaskGreen ColorWriteMask = 0x02
	ColorWriteMaskBlue  ColorWriteMask = 0x04
	ColorWriteMaskAlpha ColorWriteMask = 0x08
	ColorWriteMaskAll   ColorWriteMask = 0x0F
)

// TextureDimension specifies the dimensionality of a texture.
type TextureDimension uint32

const (
	TextureDimensionUndefined TextureDimension = 0x00
	TextureDimension1D        TextureDimension = 0x01
	TextureDimension2D        TextureDimension = 0x02
	TextureDimension3D        TextureDimension = 0x03
)

// AddressMode determines how texture coordinates outside [0, 1] are handled.
type AddressMode uint32

const (
	AddressModeUndefined    AddressMode = 0x00
	AddressModeClampToEdge  AddressMode = 0x01
	AddressModeRepeat       AddressMode = 0x02
	AddressModeMirrorRepeat AddressMode = 0x03
)

// FilterMode determines how textures are sampled.
type FilterMode uint32

const (
	FilterModeUndefined FilterMode = 0x00
	FilterModeNearest   FilterMode = 0x01
	FilterModeLinear    FilterMode = 0x02
)

// MipmapFilterMode determines how mipmaps are sampled.
type MipmapFilterMode uint32

const (
	MipmapFilterModeUndefined MipmapFilterMode = 0x00
	MipmapFilterModeNearest   MipmapFilterMode = 0x01
	MipmapFilterModeLinear    MipmapFilterMode = 0x02
)

// CompareFunction for depth/stencil operations.
type CompareFunction uint32

const (
	CompareFunctionUndefined    CompareFunction = 0x00
	CompareFunctionNever        CompareFunction = 0x01
	CompareFunctionLess         CompareFunction = 0x02
	CompareFunctionEqual        CompareFunction = 0x03
	CompareFunctionLessEqual    CompareFunction = 0x04
	CompareFunctionGreater      CompareFunction = 0x05
	CompareFunctionNotEqual     CompareFunction = 0x06
	CompareFunctionGreaterEqual CompareFunction = 0x07
	CompareFunctionAlways       CompareFunction = 0x08
)

// OptionalBool is a tri-state boolean for WebGPU.
type OptionalBool uint32

const (
	OptionalBoolFalse     OptionalBool = 0x00000000
	OptionalBoolTrue      OptionalBool = 0x00000001
	OptionalBoolUndefined OptionalBool = 0x00000002
)

// StencilOperation describes stencil buffer operations.
type StencilOperation uint32

const (
	StencilOperationUndefined      StencilOperation = 0x00000000
	StencilOperationKeep           StencilOperation = 0x00000001
	StencilOperationZero           StencilOperation = 0x00000002
	StencilOperationReplace        StencilOperation = 0x00000003
	StencilOperationInvert         StencilOperation = 0x00000004
	StencilOperationIncrementClamp StencilOperation = 0x00000005
	StencilOperationDecrementClamp StencilOperation = 0x00000006
	StencilOperationIncrementWrap  StencilOperation = 0x00000007
	StencilOperationDecrementWrap  StencilOperation = 0x00000008
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
