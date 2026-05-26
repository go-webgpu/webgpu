package wgpu

// RequestAdapterStatus is the status returned by RequestAdapter callback.
type RequestAdapterStatus uint32

const (
	// RequestAdapterStatusSuccess indicates the adapter was successfully obtained.
	RequestAdapterStatusSuccess RequestAdapterStatus = 0x00000001
	// RequestAdapterStatusCallbackCancelled indicates the operation was cancelled (e.g. instance dropped).
	// Renamed from InstanceDropped in v29.
	RequestAdapterStatusCallbackCancelled RequestAdapterStatus = 0x00000002
	// RequestAdapterStatusUnavailable indicates no suitable adapter is available.
	RequestAdapterStatusUnavailable RequestAdapterStatus = 0x00000003
	// RequestAdapterStatusError indicates an error occurred during adapter request.
	RequestAdapterStatusError RequestAdapterStatus = 0x00000004
	// RequestAdapterStatusInstanceDropped is a deprecated alias for CallbackCancelled.
	// Deprecated: Use RequestAdapterStatusCallbackCancelled.
	RequestAdapterStatusInstanceDropped = RequestAdapterStatusCallbackCancelled
)

// RequestDeviceStatus is the status returned by RequestDevice callback.
type RequestDeviceStatus uint32

const (
	// RequestDeviceStatusSuccess indicates the device was successfully obtained.
	RequestDeviceStatusSuccess RequestDeviceStatus = 0x00000001
	// RequestDeviceStatusCallbackCancelled indicates the operation was cancelled (e.g. instance dropped).
	// Renamed from InstanceDropped in v29.
	RequestDeviceStatusCallbackCancelled RequestDeviceStatus = 0x00000002
	// RequestDeviceStatusError indicates an error occurred during device request.
	RequestDeviceStatusError RequestDeviceStatus = 0x00000003
	// RequestDeviceStatusInstanceDropped is a deprecated alias for CallbackCancelled.
	// Deprecated: Use RequestDeviceStatusCallbackCancelled.
	RequestDeviceStatusInstanceDropped = RequestDeviceStatusCallbackCancelled
)

// FeatureLevel indicates the WebGPU feature level.
type FeatureLevel uint32

const (
	// FeatureLevelUndefined indicates no value is passed. Added in v29.
	FeatureLevelUndefined FeatureLevel = 0x00000000
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
	// STypeRenderPassMaxDrawCount identifies the max draw count chained struct.
	STypeRenderPassMaxDrawCount SType = 0x00000003

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
	// STypeSurfaceColorManagement identifies a surface color management chained struct. New in v29.
	STypeSurfaceColorManagement SType = 0x0000000A
	// STypeRequestAdapterWebXROptions identifies WebXR adapter options. New in v29.
	STypeRequestAdapterWebXROptions SType = 0x0000000B
	// STypeTextureComponentSwizzleDescriptor identifies a texture component swizzle descriptor. New in v29.
	STypeTextureComponentSwizzleDescriptor SType = 0x0000000C
	// STypeExternalTextureBindingLayout identifies an external texture binding layout. New in v29.
	STypeExternalTextureBindingLayout SType = 0x0000000D
	// STypeExternalTextureBindingEntry identifies an external texture binding entry. New in v29.
	STypeExternalTextureBindingEntry SType = 0x0000000E
	// STypeCompatibilityModeLimits identifies compat-mode limits. New in v29.
	STypeCompatibilityModeLimits SType = 0x0000000F
	// STypeTextureBindingViewDimension identifies a texture binding view dimension. New in v29.
	STypeTextureBindingViewDimension SType = 0x00000010

	// Native wgpu-native extension STypes (0x0003XXXX range)

	// STypeDeviceExtras identifies wgpu-native device extras chained struct.
	STypeDeviceExtras SType = 0x00030001
	// STypeNativeLimits identifies wgpu-native native limits chained struct.
	STypeNativeLimits SType = 0x00030002
	// STypePipelineLayoutExtras identifies wgpu-native pipeline layout extras chained struct.
	STypePipelineLayoutExtras SType = 0x00030003
	// STypeShaderSourceGLSL identifies a GLSL shader source chained struct (wgpu-native extension).
	STypeShaderSourceGLSL SType = 0x00030004
	// STypeInstanceExtras identifies wgpu-native instance extras chained struct.
	STypeInstanceExtras SType = 0x00030006
	// STypeBindGroupEntryExtras identifies wgpu-native bind group entry extras.
	STypeBindGroupEntryExtras SType = 0x00030007
	// STypeBindGroupLayoutEntryExtras identifies wgpu-native bind group layout entry extras.
	STypeBindGroupLayoutEntryExtras SType = 0x00030008
	// STypeQuerySetDescriptorExtras identifies wgpu-native query set descriptor extras.
	STypeQuerySetDescriptorExtras SType = 0x00030009
	// STypeSurfaceConfigurationExtras identifies wgpu-native surface configuration extras.
	STypeSurfaceConfigurationExtras SType = 0x0003000A
	// STypeSurfaceSourceSwapChainPanel identifies a WinUI SwapChainPanel surface source.
	STypeSurfaceSourceSwapChainPanel SType = 0x0003000B
	// STypePrimitiveStateExtras identifies wgpu-native primitive state extras.
	STypePrimitiveStateExtras SType = 0x0003000C
)

// SurfaceGetCurrentTextureStatus describes the result of GetCurrentTexture.
type SurfaceGetCurrentTextureStatus uint32

const (
	// SurfaceGetCurrentTextureStatusSuccessOptimal indicates the texture was obtained optimally.
	SurfaceGetCurrentTextureStatusSuccessOptimal SurfaceGetCurrentTextureStatus = 0x00000001
	// SurfaceGetCurrentTextureStatusSuccessSuboptimal indicates the texture was obtained but may not be optimal.
	SurfaceGetCurrentTextureStatusSuccessSuboptimal SurfaceGetCurrentTextureStatus = 0x00000002
	// SurfaceGetCurrentTextureStatusTimeout indicates the operation timed out.
	SurfaceGetCurrentTextureStatusTimeout SurfaceGetCurrentTextureStatus = 0x00000003
	// SurfaceGetCurrentTextureStatusOutdated indicates the surface needs reconfiguration.
	SurfaceGetCurrentTextureStatusOutdated SurfaceGetCurrentTextureStatus = 0x00000004
	// SurfaceGetCurrentTextureStatusLost indicates the surface was lost and must be recreated.
	SurfaceGetCurrentTextureStatusLost SurfaceGetCurrentTextureStatus = 0x00000005
	// SurfaceGetCurrentTextureStatusError indicates a deterministic error (e.g. surface not configured).
	// BREAKING: v27 had OutOfMemory=0x06, DeviceLost=0x07, Error=0x08; v29 collapsed to Error=0x06.
	SurfaceGetCurrentTextureStatusError SurfaceGetCurrentTextureStatus = 0x00000006

	// NativeSurfaceGetCurrentTextureStatusOccluded is a wgpu-native extension status.
	// Returned on macOS Metal when the window is occluded/minimized.
	NativeSurfaceGetCurrentTextureStatusOccluded SurfaceGetCurrentTextureStatus = 0x00030001
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
	// FeatureNameCoreFeaturesAndLimits indicates core features and limits support.
	FeatureNameCoreFeaturesAndLimits FeatureName = 0x00000001
	// FeatureNameDepthClipControl enables depth clip control.
	FeatureNameDepthClipControl FeatureName = 0x00000002
	// FeatureNameDepth32FloatStencil8 enables depth32float-stencil8 texture format.
	FeatureNameDepth32FloatStencil8 FeatureName = 0x00000003
	// FeatureNameTextureCompressionBC enables BC texture compression.
	FeatureNameTextureCompressionBC FeatureName = 0x00000004
	// FeatureNameTextureCompressionBCSliced3D enables sliced 3D BC compression.
	FeatureNameTextureCompressionBCSliced3D FeatureName = 0x00000005
	// FeatureNameTextureCompressionETC2 enables ETC2 texture compression.
	FeatureNameTextureCompressionETC2 FeatureName = 0x00000006
	// FeatureNameTextureCompressionASTC enables ASTC texture compression.
	FeatureNameTextureCompressionASTC FeatureName = 0x00000007
	// FeatureNameTextureCompressionASTCSliced3D enables sliced 3D ASTC compression.
	FeatureNameTextureCompressionASTCSliced3D FeatureName = 0x00000008
	// FeatureNameTimestampQuery enables timestamp query support.
	FeatureNameTimestampQuery FeatureName = 0x00000009
	// FeatureNameIndirectFirstInstance enables indirect first instance.
	FeatureNameIndirectFirstInstance FeatureName = 0x0000000A
	// FeatureNameShaderF16 enables f16 in shaders.
	FeatureNameShaderF16 FeatureName = 0x0000000B
	// FeatureNameRG11B10UfloatRenderable enables RG11B10Ufloat as render target.
	FeatureNameRG11B10UfloatRenderable FeatureName = 0x0000000C
	// FeatureNameBGRA8UnormStorage enables BGRA8Unorm storage textures.
	FeatureNameBGRA8UnormStorage FeatureName = 0x0000000D
	// FeatureNameFloat32Filterable enables filterable float32 textures.
	FeatureNameFloat32Filterable FeatureName = 0x0000000E
	// FeatureNameFloat32Blendable enables blendable float32 textures.
	FeatureNameFloat32Blendable FeatureName = 0x0000000F
	// FeatureNameClipDistances enables clip distances in shaders.
	FeatureNameClipDistances FeatureName = 0x00000010
	// FeatureNameDualSourceBlending enables dual source blending.
	FeatureNameDualSourceBlending FeatureName = 0x00000011
	// FeatureNameSubgroups enables subgroup operations.
	FeatureNameSubgroups FeatureName = 0x00000012
	// FeatureNameTextureFormatsTier1 enables tier 1 texture formats.
	FeatureNameTextureFormatsTier1 FeatureName = 0x00000013
	// FeatureNameTextureFormatsTier2 enables tier 2 texture formats.
	FeatureNameTextureFormatsTier2 FeatureName = 0x00000014
	// FeatureNamePrimitiveIndex enables primitive index in shaders.
	FeatureNamePrimitiveIndex FeatureName = 0x00000015
	// FeatureNameTextureComponentSwizzle enables texture component swizzle.
	FeatureNameTextureComponentSwizzle FeatureName = 0x00000016
)

// NativeFeature describes a wgpu-native extension feature.
type NativeFeature uint32

const (
	// NativeFeatureImmediates enables immediate data (push constants replacement).
	// Renamed from PushConstants in v29.
	NativeFeatureImmediates NativeFeature = 0x00030001
	// NativeFeaturePushConstants is a deprecated alias for Immediates.
	// Deprecated: Use NativeFeatureImmediates.
	NativeFeaturePushConstants = NativeFeatureImmediates
	// NativeFeatureTextureAdapterSpecificFormatFeatures enables device-specific texture format features.
	NativeFeatureTextureAdapterSpecificFormatFeatures NativeFeature = 0x00030002
	// NativeFeatureMultiDrawIndirectCount enables indirect draw count.
	NativeFeatureMultiDrawIndirectCount NativeFeature = 0x00030004
	// NativeFeatureVertexWritableStorage enables vertex shader writable storage.
	NativeFeatureVertexWritableStorage NativeFeature = 0x00030005
	// NativeFeatureTextureBindingArray enables texture binding arrays.
	NativeFeatureTextureBindingArray NativeFeature = 0x00030006
	// NativeFeatureSampledTextureAndStorageBufferArrayNonUniformIndexing enables non-uniform indexing.
	NativeFeatureSampledTextureAndStorageBufferArrayNonUniformIndexing NativeFeature = 0x00030007
	// NativeFeaturePipelineStatisticsQuery enables pipeline statistics queries.
	NativeFeaturePipelineStatisticsQuery NativeFeature = 0x00030008
	// NativeFeatureStorageResourceBindingArray enables storage resource binding arrays.
	NativeFeatureStorageResourceBindingArray NativeFeature = 0x00030009
	// NativeFeaturePartiallyBoundBindingArray enables partially bound binding arrays.
	NativeFeaturePartiallyBoundBindingArray NativeFeature = 0x0003000A
	// NativeFeatureTextureFormat16bitNorm enables normalized 16-bit texture formats.
	NativeFeatureTextureFormat16bitNorm NativeFeature = 0x0003000B
	// NativeFeatureTextureCompressionAstcHdr enables ASTC HDR compression.
	NativeFeatureTextureCompressionAstcHdr NativeFeature = 0x0003000C
	// NativeFeatureMappablePrimaryBuffers enables mappable primary buffers.
	NativeFeatureMappablePrimaryBuffers NativeFeature = 0x0003000E
	// NativeFeatureBufferBindingArray enables buffer binding arrays.
	NativeFeatureBufferBindingArray NativeFeature = 0x0003000F
	// NativeFeatureUniformBufferAndStorageTextureArrayNonUniformIndexing enables non-uniform indexing for these types.
	NativeFeatureUniformBufferAndStorageTextureArrayNonUniformIndexing NativeFeature = 0x00030010
	// NativeFeaturePolygonModeLine enables polygon line mode.
	NativeFeaturePolygonModeLine NativeFeature = 0x00030013
	// NativeFeaturePolygonModePoint enables polygon point mode.
	NativeFeaturePolygonModePoint NativeFeature = 0x00030014
	// NativeFeatureConservativeRasterization enables conservative rasterization.
	NativeFeatureConservativeRasterization NativeFeature = 0x00030015
	// NativeFeatureSpirvShaderPassthrough enables SPIR-V shader passthrough.
	NativeFeatureSpirvShaderPassthrough NativeFeature = 0x00030017
	// NativeFeatureVertexAttribute64bit enables 64-bit vertex attributes.
	NativeFeatureVertexAttribute64bit NativeFeature = 0x00030019
	// NativeFeatureTextureFormatNv12 enables NV12 texture format.
	NativeFeatureTextureFormatNv12 NativeFeature = 0x0003001A
	// NativeFeatureRayQuery enables ray query in shaders.
	NativeFeatureRayQuery NativeFeature = 0x0003001C
	// NativeFeatureShaderF64 enables f64 in shaders.
	NativeFeatureShaderF64 NativeFeature = 0x0003001D
	// NativeFeatureShaderI16 enables i16 in shaders.
	NativeFeatureShaderI16 NativeFeature = 0x0003001E
	// NativeFeatureShaderEarlyDepthTest enables early depth test attribute in shaders.
	NativeFeatureShaderEarlyDepthTest NativeFeature = 0x00030020
	// NativeFeatureSubgroup enables subgroup operations in compute/fragment shaders.
	NativeFeatureSubgroup NativeFeature = 0x00030021
	// NativeFeatureSubgroupVertex enables subgroup operations in vertex shaders.
	NativeFeatureSubgroupVertex NativeFeature = 0x00030022
	// NativeFeatureSubgroupBarrier enables subgroup barrier in compute shaders.
	NativeFeatureSubgroupBarrier NativeFeature = 0x00030023
	// NativeFeatureTimestampQueryInsideEncoders enables timestamp queries on command encoders.
	NativeFeatureTimestampQueryInsideEncoders NativeFeature = 0x00030024
	// NativeFeatureTimestampQueryInsidePasses enables timestamp queries inside render/compute passes.
	NativeFeatureTimestampQueryInsidePasses NativeFeature = 0x00030025
	// NativeFeatureShaderInt64 enables i64/u64 in shaders.
	NativeFeatureShaderInt64 NativeFeature = 0x00030026
)

// InstanceFeatureName describes features that can be required at instance creation.
// New in v29.
type InstanceFeatureName uint32

const (
	// InstanceFeatureNameTimedWaitAny enables wgpuInstanceWaitAny with timeoutNS > 0.
	InstanceFeatureNameTimedWaitAny InstanceFeatureName = 0x00000001
	// InstanceFeatureNameShaderSourceSPIRV enables SPIR-V shader sources.
	InstanceFeatureNameShaderSourceSPIRV InstanceFeatureName = 0x00000002
	// InstanceFeatureNameMultipleDevicesPerAdapter allows multiple devices per adapter.
	InstanceFeatureNameMultipleDevicesPerAdapter InstanceFeatureName = 0x00000003
)

// ComponentSwizzle describes texture component swizzle mapping.
// New in v29.
type ComponentSwizzle uint32

const (
	// ComponentSwizzleUndefined indicates no value is passed.
	ComponentSwizzleUndefined ComponentSwizzle = 0x00000000
	// ComponentSwizzleZero forces the component value to 0.
	ComponentSwizzleZero ComponentSwizzle = 0x00000001
	// ComponentSwizzleOne forces the component value to 1.
	ComponentSwizzleOne ComponentSwizzle = 0x00000002
	// ComponentSwizzleR takes from the red channel.
	ComponentSwizzleR ComponentSwizzle = 0x00000003
	// ComponentSwizzleG takes from the green channel.
	ComponentSwizzleG ComponentSwizzle = 0x00000004
	// ComponentSwizzleB takes from the blue channel.
	ComponentSwizzleB ComponentSwizzle = 0x00000005
	// ComponentSwizzleA takes from the alpha channel.
	ComponentSwizzleA ComponentSwizzle = 0x00000006
)

// PredefinedColorSpace describes a color space for surface color management.
// New in v29.
type PredefinedColorSpace uint32

const (
	// PredefinedColorSpaceSRGB is the standard sRGB color space.
	PredefinedColorSpaceSRGB PredefinedColorSpace = 0x00000001
	// PredefinedColorSpaceDisplayP3 is the Display P3 wide-gamut color space.
	PredefinedColorSpaceDisplayP3 PredefinedColorSpace = 0x00000002
)

// ToneMappingMode describes tone mapping for HDR surfaces.
// New in v29.
type ToneMappingMode uint32

const (
	// ToneMappingModeStandard is standard tone mapping (sRGB).
	ToneMappingModeStandard ToneMappingMode = 0x00000001
	// ToneMappingModeExtended is extended (HDR) tone mapping.
	ToneMappingModeExtended ToneMappingMode = 0x00000002
)

// WGPUStatus describes the status returned from certain WebGPU operations.
// Note: v29 changed Success from 0x00 to 0x01 and Error from 0x01 to 0x02.
type WGPUStatus uint32

const (
	// WGPUStatusSuccess indicates the operation completed successfully.
	WGPUStatusSuccess WGPUStatus = 0x00000001
	// WGPUStatusError indicates the operation failed.
	WGPUStatusError WGPUStatus = 0x00000002
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
	// PopErrorScopeStatusCallbackCancelled indicates the operation was cancelled (e.g. instance dropped).
	// Renamed from InstanceDropped in v29.
	PopErrorScopeStatusCallbackCancelled PopErrorScopeStatus = 0x00000002
	// PopErrorScopeStatusError indicates the error scope stack was empty or another error occurred.
	// Renamed from EmptyStack in v29.
	PopErrorScopeStatusError PopErrorScopeStatus = 0x00000003
	// PopErrorScopeStatusInstanceDropped is a deprecated alias for CallbackCancelled.
	// Deprecated: Use PopErrorScopeStatusCallbackCancelled.
	PopErrorScopeStatusInstanceDropped = PopErrorScopeStatusCallbackCancelled
	// PopErrorScopeStatusEmptyStack is a deprecated alias for Error.
	// Deprecated: Use PopErrorScopeStatusError.
	PopErrorScopeStatusEmptyStack = PopErrorScopeStatusError
)

// DeviceLostReason describes why a device was lost.
type DeviceLostReason uint32

const (
	// DeviceLostReasonUnknown indicates the device was lost for an unknown reason.
	DeviceLostReasonUnknown DeviceLostReason = 0x00000001
	// DeviceLostReasonDestroyed indicates the device was explicitly destroyed.
	DeviceLostReasonDestroyed DeviceLostReason = 0x00000002
	// DeviceLostReasonCallbackCancelled indicates the operation was cancelled (e.g. instance dropped).
	// Renamed from InstanceDropped in v29.
	DeviceLostReasonCallbackCancelled DeviceLostReason = 0x00000003
	// DeviceLostReasonFailedCreation indicates device creation failed.
	DeviceLostReasonFailedCreation DeviceLostReason = 0x00000004
	// DeviceLostReasonInstanceDropped is a deprecated alias for CallbackCancelled.
	// Deprecated: Use DeviceLostReasonCallbackCancelled.
	DeviceLostReasonInstanceDropped = DeviceLostReasonCallbackCancelled
)

// InstanceBackend is a bitflag selecting which graphics backends to enable.
// Used in InstanceExtras.Backends.
type InstanceBackend uint64

const (
	// InstanceBackendAll enables all available backends (default when zero).
	InstanceBackendAll InstanceBackend = 0x00000000
	// InstanceBackendVulkan enables the Vulkan backend.
	InstanceBackendVulkan InstanceBackend = 1 << 0
	// InstanceBackendGL enables the OpenGL/OpenGL ES backend.
	InstanceBackendGL InstanceBackend = 1 << 1
	// InstanceBackendMetal enables the Metal backend (macOS/iOS).
	InstanceBackendMetal InstanceBackend = 1 << 2
	// InstanceBackendDX12 enables the Direct3D 12 backend (Windows).
	InstanceBackendDX12 InstanceBackend = 1 << 3
	// InstanceBackendBrowserWebGPU enables the browser WebGPU backend (WASM).
	InstanceBackendBrowserWebGPU InstanceBackend = 1 << 5
	// InstanceBackendPrimary enables primary tier backends: Vulkan, Metal, DX12, BrowserWebGPU.
	InstanceBackendPrimary InstanceBackend = (1 << 0) | (1 << 2) | (1 << 3) | (1 << 5)
	// InstanceBackendSecondary enables secondary tier backends: GL only.
	// BREAKING: v27 had Secondary = GL | DX11; DX11 was removed in v29.
	InstanceBackendSecondary InstanceBackend = (1 << 1)
)

// InstanceFlag is a bitflag controlling instance debugging and validation behavior.
// Used in InstanceExtras.Flags.
type InstanceFlag uint64

const (
	// InstanceFlagEmpty has no flags set. Zero-initialization default.
	// BREAKING: v27 had Default=0; v29 renamed to Empty=0, Default moved to 1<<24.
	InstanceFlagEmpty InstanceFlag = 0x00000000
	// InstanceFlagDebug generates debug information in shaders and objects.
	InstanceFlagDebug InstanceFlag = 1 << 0
	// InstanceFlagValidation enables validation in the backend API.
	InstanceFlagValidation InstanceFlag = 1 << 1
	// InstanceFlagDiscardHalLabels suppresses label passing to backend HAL.
	InstanceFlagDiscardHalLabels InstanceFlag = 1 << 2
	// InstanceFlagAllowUnderlyingNoncompliantAdapter exposes adapters on non-compliant drivers.
	InstanceFlagAllowUnderlyingNoncompliantAdapter InstanceFlag = 1 << 3
	// InstanceFlagGPUBasedValidation enables GPU-based validation (implies Validation).
	InstanceFlagGPUBasedValidation InstanceFlag = 1 << 4
	// InstanceFlagValidationIndirectCall validates indirect buffer content before indirect draws.
	InstanceFlagValidationIndirectCall InstanceFlag = 1 << 5
	// InstanceFlagAutomaticTimestampNormalization normalizes timestamps to nanoseconds.
	InstanceFlagAutomaticTimestampNormalization InstanceFlag = 1 << 6
	// InstanceFlagDefault uses the default flags for the current build configuration.
	// BREAKING: v27 had Default=0; v29 Default=1<<24 (debug builds enable Debug+Validation).
	InstanceFlagDefault InstanceFlag = 1 << 24
	// InstanceFlagDebugging enables Debug and Validation.
	InstanceFlagDebugging InstanceFlag = 1 << 25
	// InstanceFlagAdvancedDebugging enables Debug, Validation, and GPUBasedValidation.
	InstanceFlagAdvancedDebugging InstanceFlag = 1 << 26
	// InstanceFlagWithEnv reads flag overrides from environment variables.
	InstanceFlagWithEnv InstanceFlag = 1 << 27
)

// NativeDisplayHandleType identifies the platform display connection type.
// Used in NativeDisplayHandle.Type. New in v29.
type NativeDisplayHandleType uint32

const (
	// NativeDisplayHandleTypeNone indicates no display handle. Default (zero-init).
	NativeDisplayHandleTypeNone NativeDisplayHandleType = 0x00000000
	// NativeDisplayHandleTypeXlib is an X11 display via Xlib.
	NativeDisplayHandleTypeXlib NativeDisplayHandleType = 0x00000001
	// NativeDisplayHandleTypeXcb is an X11 display via XCB.
	NativeDisplayHandleTypeXcb NativeDisplayHandleType = 0x00000002
	// NativeDisplayHandleTypeWayland is a Wayland display connection.
	NativeDisplayHandleTypeWayland NativeDisplayHandleType = 0x00000003
)
