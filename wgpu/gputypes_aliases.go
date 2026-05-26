package wgpu

import "github.com/gogpu/gputypes"

// Type aliases from gputypes for single-import ergonomics.
// Importing "github.com/go-webgpu/webgpu/wgpu" is sufficient — no separate
// gputypes import required when using these aliases.

// Extent3D is a 3D extent (width/height/depth or array layers).
type Extent3D = gputypes.Extent3D

// Origin3D is a 3D origin (x/y/z or array layer offset).
type Origin3D = gputypes.Origin3D

// Color is an RGBA color with double precision.
// Note: wgpu package also defines Color struct for render pass clear values.
// This alias shadows it — use wgpu.Color directly for the render pass type.

// MapMode specifies buffer mapping mode.
// Note: MapMode is already defined as a native type in buffer.go (uint64).

// Texture types.
// TextureAspect is defined as a native enum in enums.go.
type TextureFormat = gputypes.TextureFormat
type TextureDimension = gputypes.TextureDimension
type TextureViewDimension = gputypes.TextureViewDimension

// Buffer types.
type BufferUsage = gputypes.BufferUsage

// Texture usage type.
type TextureUsage = gputypes.TextureUsage

// Shader stage type.
// ShaderStage is the uint32 bitflag for individual stages (vertex/fragment/compute).
// ShaderStages is an alias for the same type for use in pipeline descriptors.
type ShaderStage = gputypes.ShaderStage
type ShaderStages = gputypes.ShaderStages

// Primitive assembly types.
type PrimitiveTopology = gputypes.PrimitiveTopology
type FrontFace = gputypes.FrontFace
type CullMode = gputypes.CullMode
type IndexFormat = gputypes.IndexFormat

// Blend types.
type BlendFactor = gputypes.BlendFactor
type BlendOperation = gputypes.BlendOperation
type ColorWriteMask = gputypes.ColorWriteMask

// Depth/stencil types.
type CompareFunction = gputypes.CompareFunction
type StencilOperation = gputypes.StencilOperation

// Vertex types.
type VertexFormat = gputypes.VertexFormat
type VertexStepMode = gputypes.VertexStepMode

// Sampler types.
type FilterMode = gputypes.FilterMode
type MipmapFilterMode = gputypes.MipmapFilterMode
type AddressMode = gputypes.AddressMode

// Surface/presentation types.
type PresentMode = gputypes.PresentMode
type CompositeAlphaMode = gputypes.CompositeAlphaMode

// Adapter types.
type PowerPreference = gputypes.PowerPreference

// Render pass types.
type LoadOp = gputypes.LoadOp
type StoreOp = gputypes.StoreOp

// Features is a bitmask of enabled GPU features.
// Note: Limits is defined as a native FFI struct in adapter.go (matches wgpu-native ABI).
type Features = gputypes.Features

// --- BufferUsage constants ---

const (
	BufferUsageNone         = gputypes.BufferUsageNone
	BufferUsageMapRead      = gputypes.BufferUsageMapRead
	BufferUsageMapWrite     = gputypes.BufferUsageMapWrite
	BufferUsageCopySrc      = gputypes.BufferUsageCopySrc
	BufferUsageCopyDst      = gputypes.BufferUsageCopyDst
	BufferUsageIndex        = gputypes.BufferUsageIndex
	BufferUsageVertex       = gputypes.BufferUsageVertex
	BufferUsageUniform      = gputypes.BufferUsageUniform
	BufferUsageStorage      = gputypes.BufferUsageStorage
	BufferUsageIndirect     = gputypes.BufferUsageIndirect
	BufferUsageQueryResolve = gputypes.BufferUsageQueryResolve
)

// --- TextureUsage constants ---

const (
	TextureUsageNone             = gputypes.TextureUsageNone
	TextureUsageCopySrc          = gputypes.TextureUsageCopySrc
	TextureUsageCopyDst          = gputypes.TextureUsageCopyDst
	TextureUsageTextureBinding   = gputypes.TextureUsageTextureBinding
	TextureUsageStorageBinding   = gputypes.TextureUsageStorageBinding
	TextureUsageRenderAttachment = gputypes.TextureUsageRenderAttachment
)

// --- TextureFormat constants ---

const (
	TextureFormatUndefined     = gputypes.TextureFormatUndefined
	TextureFormatR8Unorm       = gputypes.TextureFormatR8Unorm
	TextureFormatR8Snorm       = gputypes.TextureFormatR8Snorm
	TextureFormatR8Uint        = gputypes.TextureFormatR8Uint
	TextureFormatR8Sint        = gputypes.TextureFormatR8Sint
	TextureFormatR16Uint       = gputypes.TextureFormatR16Uint
	TextureFormatR16Sint       = gputypes.TextureFormatR16Sint
	TextureFormatR16Float      = gputypes.TextureFormatR16Float
	TextureFormatRG8Unorm      = gputypes.TextureFormatRG8Unorm
	TextureFormatRG8Snorm      = gputypes.TextureFormatRG8Snorm
	TextureFormatRG8Uint       = gputypes.TextureFormatRG8Uint
	TextureFormatRG8Sint       = gputypes.TextureFormatRG8Sint
	TextureFormatR32Float      = gputypes.TextureFormatR32Float
	TextureFormatR32Uint       = gputypes.TextureFormatR32Uint
	TextureFormatR32Sint       = gputypes.TextureFormatR32Sint
	TextureFormatRG16Uint      = gputypes.TextureFormatRG16Uint
	TextureFormatRG16Sint      = gputypes.TextureFormatRG16Sint
	TextureFormatRG16Float     = gputypes.TextureFormatRG16Float
	TextureFormatRGBA8Unorm    = gputypes.TextureFormatRGBA8Unorm
	TextureFormatRGBA8UnormSrgb = gputypes.TextureFormatRGBA8UnormSrgb
	TextureFormatRGBA8Snorm    = gputypes.TextureFormatRGBA8Snorm
	TextureFormatRGBA8Uint     = gputypes.TextureFormatRGBA8Uint
	TextureFormatRGBA8Sint     = gputypes.TextureFormatRGBA8Sint
	TextureFormatBGRA8Unorm    = gputypes.TextureFormatBGRA8Unorm
	TextureFormatBGRA8UnormSrgb = gputypes.TextureFormatBGRA8UnormSrgb
	TextureFormatRGB10A2Uint   = gputypes.TextureFormatRGB10A2Uint
	TextureFormatRGB10A2Unorm  = gputypes.TextureFormatRGB10A2Unorm
	TextureFormatRG11B10Ufloat = gputypes.TextureFormatRG11B10Ufloat
	TextureFormatRG32Float     = gputypes.TextureFormatRG32Float
	TextureFormatRG32Uint      = gputypes.TextureFormatRG32Uint
	TextureFormatRG32Sint      = gputypes.TextureFormatRG32Sint
	TextureFormatRGBA16Uint    = gputypes.TextureFormatRGBA16Uint
	TextureFormatRGBA16Sint    = gputypes.TextureFormatRGBA16Sint
	TextureFormatRGBA16Float   = gputypes.TextureFormatRGBA16Float
	TextureFormatRGBA32Float   = gputypes.TextureFormatRGBA32Float
	TextureFormatRGBA32Uint    = gputypes.TextureFormatRGBA32Uint
	TextureFormatRGBA32Sint    = gputypes.TextureFormatRGBA32Sint
	TextureFormatDepth32Float  = gputypes.TextureFormatDepth32Float
	TextureFormatDepth24Plus   = gputypes.TextureFormatDepth24Plus
	TextureFormatDepth24PlusStencil8 = gputypes.TextureFormatDepth24PlusStencil8
	TextureFormatDepth16Unorm  = gputypes.TextureFormatDepth16Unorm
)

// --- TextureDimension constants ---

const (
	TextureDimension1D = gputypes.TextureDimension1D
	TextureDimension2D = gputypes.TextureDimension2D
	TextureDimension3D = gputypes.TextureDimension3D
)

// --- ShaderStage constants ---

const (
	ShaderStageNone     = gputypes.ShaderStageNone
	ShaderStageVertex   = gputypes.ShaderStageVertex
	ShaderStageFragment = gputypes.ShaderStageFragment
	ShaderStageCompute  = gputypes.ShaderStageCompute
)

// --- PrimitiveTopology constants ---

const (
	PrimitiveTopologyPointList     = gputypes.PrimitiveTopologyPointList
	PrimitiveTopologyLineList      = gputypes.PrimitiveTopologyLineList
	PrimitiveTopologyLineStrip     = gputypes.PrimitiveTopologyLineStrip
	PrimitiveTopologyTriangleList  = gputypes.PrimitiveTopologyTriangleList
	PrimitiveTopologyTriangleStrip = gputypes.PrimitiveTopologyTriangleStrip
)

// --- FrontFace constants ---

const (
	FrontFaceCCW = gputypes.FrontFaceCCW
	FrontFaceCW  = gputypes.FrontFaceCW
)

// --- CullMode constants ---

const (
	CullModeNone  = gputypes.CullModeNone
	CullModeFront = gputypes.CullModeFront
	CullModeBack  = gputypes.CullModeBack
)

// --- IndexFormat constants ---

const (
	IndexFormatUint16 = gputypes.IndexFormatUint16
	IndexFormatUint32 = gputypes.IndexFormatUint32
)

// --- LoadOp constants ---

const (
	LoadOpLoad  = gputypes.LoadOpLoad
	LoadOpClear = gputypes.LoadOpClear
)

// --- StoreOp constants ---

const (
	StoreOpStore   = gputypes.StoreOpStore
	StoreOpDiscard = gputypes.StoreOpDiscard
)

// --- FilterMode constants ---

const (
	FilterModeNearest = gputypes.FilterModeNearest
	FilterModeLinear  = gputypes.FilterModeLinear
)

// --- AddressMode constants ---

const (
	AddressModeRepeat            = gputypes.AddressModeRepeat
	AddressModeMirrorRepeat      = gputypes.AddressModeMirrorRepeat
	AddressModeClampToEdge       = gputypes.AddressModeClampToEdge
)

// --- CompareFunction constants ---

const (
	CompareFunctionUndefined    = gputypes.CompareFunctionUndefined
	CompareFunctionNever        = gputypes.CompareFunctionNever
	CompareFunctionLess         = gputypes.CompareFunctionLess
	CompareFunctionEqual        = gputypes.CompareFunctionEqual
	CompareFunctionLessEqual    = gputypes.CompareFunctionLessEqual
	CompareFunctionGreater      = gputypes.CompareFunctionGreater
	CompareFunctionNotEqual     = gputypes.CompareFunctionNotEqual
	CompareFunctionGreaterEqual = gputypes.CompareFunctionGreaterEqual
	CompareFunctionAlways       = gputypes.CompareFunctionAlways
)

// --- PresentMode constants ---

const (
	PresentModeImmediate   = gputypes.PresentModeImmediate
	PresentModeMailbox     = gputypes.PresentModeMailbox
	PresentModeFifo        = gputypes.PresentModeFifo
	PresentModeFifoRelaxed = gputypes.PresentModeFifoRelaxed
)

// --- CompositeAlphaMode constants ---

const (
	CompositeAlphaModeAuto            = gputypes.CompositeAlphaModeAuto
	CompositeAlphaModeOpaque          = gputypes.CompositeAlphaModeOpaque
	CompositeAlphaModePremultiplied   = gputypes.CompositeAlphaModePremultiplied
	CompositeAlphaModeUnpremultiplied = gputypes.CompositeAlphaModeUnpremultiplied
	CompositeAlphaModeInherit         = gputypes.CompositeAlphaModeInherit
)

// --- PowerPreference constants ---

const (
	PowerPreferenceNone            = gputypes.PowerPreferenceNone
	PowerPreferenceLowPower        = gputypes.PowerPreferenceLowPower
	PowerPreferenceHighPerformance = gputypes.PowerPreferenceHighPerformance
)

// --- ColorWriteMask constants ---

const (
	ColorWriteMaskNone  = gputypes.ColorWriteMaskNone
	ColorWriteMaskRed   = gputypes.ColorWriteMaskRed
	ColorWriteMaskGreen = gputypes.ColorWriteMaskGreen
	ColorWriteMaskBlue  = gputypes.ColorWriteMaskBlue
	ColorWriteMaskAlpha = gputypes.ColorWriteMaskAlpha
	ColorWriteMaskAll   = gputypes.ColorWriteMaskAll
)

// --- VertexFormat constants ---

const (
	VertexFormatUint8x2  = gputypes.VertexFormatUint8x2
	VertexFormatUint8x4  = gputypes.VertexFormatUint8x4
	VertexFormatSint8x2  = gputypes.VertexFormatSint8x2
	VertexFormatSint8x4  = gputypes.VertexFormatSint8x4
	VertexFormatFloat32  = gputypes.VertexFormatFloat32
	VertexFormatFloat32x2 = gputypes.VertexFormatFloat32x2
	VertexFormatFloat32x3 = gputypes.VertexFormatFloat32x3
	VertexFormatFloat32x4 = gputypes.VertexFormatFloat32x4
	VertexFormatUint32   = gputypes.VertexFormatUint32
	VertexFormatUint32x2 = gputypes.VertexFormatUint32x2
	VertexFormatUint32x3 = gputypes.VertexFormatUint32x3
	VertexFormatUint32x4 = gputypes.VertexFormatUint32x4
	VertexFormatSint32   = gputypes.VertexFormatSint32
	VertexFormatSint32x2 = gputypes.VertexFormatSint32x2
	VertexFormatSint32x3 = gputypes.VertexFormatSint32x3
	VertexFormatSint32x4 = gputypes.VertexFormatSint32x4
)

// --- VertexStepMode constants ---

const (
	VertexStepModeVertex   = gputypes.VertexStepModeVertex
	VertexStepModeInstance = gputypes.VertexStepModeInstance
)

// Binding layout types.
type BufferBindingType = gputypes.BufferBindingType
type SamplerBindingType = gputypes.SamplerBindingType
type TextureSampleType = gputypes.TextureSampleType

// --- BufferBindingType constants ---

const (
	BufferBindingTypeUndefined       = gputypes.BufferBindingTypeUndefined
	BufferBindingTypeUniform         = gputypes.BufferBindingTypeUniform
	BufferBindingTypeStorage         = gputypes.BufferBindingTypeStorage
	BufferBindingTypeReadOnlyStorage = gputypes.BufferBindingTypeReadOnlyStorage
)

// --- SamplerBindingType constants ---

const (
	SamplerBindingTypeUndefined   = gputypes.SamplerBindingTypeUndefined
	SamplerBindingTypeFiltering   = gputypes.SamplerBindingTypeFiltering
	SamplerBindingTypeNonFiltering = gputypes.SamplerBindingTypeNonFiltering
	SamplerBindingTypeComparison  = gputypes.SamplerBindingTypeComparison
)

// --- TextureSampleType constants ---

const (
	TextureSampleTypeUndefined         = gputypes.TextureSampleTypeUndefined
	TextureSampleTypeFloat             = gputypes.TextureSampleTypeFloat
	TextureSampleTypeUnfilterableFloat = gputypes.TextureSampleTypeUnfilterableFloat
	TextureSampleTypeDepth             = gputypes.TextureSampleTypeDepth
	TextureSampleTypeSint              = gputypes.TextureSampleTypeSint
	TextureSampleTypeUint              = gputypes.TextureSampleTypeUint
)

// --- TextureViewDimension constants ---

const (
	TextureViewDimensionUndefined = gputypes.TextureViewDimensionUndefined
	TextureViewDimension1D        = gputypes.TextureViewDimension1D
	TextureViewDimension2D        = gputypes.TextureViewDimension2D
	TextureViewDimension2DArray   = gputypes.TextureViewDimension2DArray
	TextureViewDimensionCube      = gputypes.TextureViewDimensionCube
	TextureViewDimensionCubeArray = gputypes.TextureViewDimensionCubeArray
	TextureViewDimension3D        = gputypes.TextureViewDimension3D
)

// --- StencilOperation constants ---

const (
	StencilOperationUndefined      = gputypes.StencilOperationUndefined
	StencilOperationKeep           = gputypes.StencilOperationKeep
	StencilOperationZero           = gputypes.StencilOperationZero
	StencilOperationReplace        = gputypes.StencilOperationReplace
	StencilOperationInvert         = gputypes.StencilOperationInvert
	StencilOperationIncrementClamp = gputypes.StencilOperationIncrementClamp
	StencilOperationDecrementClamp = gputypes.StencilOperationDecrementClamp
	StencilOperationIncrementWrap  = gputypes.StencilOperationIncrementWrap
	StencilOperationDecrementWrap  = gputypes.StencilOperationDecrementWrap
)
