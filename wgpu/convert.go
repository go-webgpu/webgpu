// convert.go provides conversion functions between gputypes and wgpu-native v29 wire values.
//
// # Why conversions are needed
//
// gputypes follows the WebGPU JS specification numbering where enum values start at
// Undefined=0. wgpu-native v29 introduces BindingNotUsed=0 sentinels for binding-related
// enums, shifting all other values by +1. VertexFormat and VertexStepMode have additional
// structural differences (missing single-component variants, removed VertexBufferNotUsed).
//
// # Enums requiring explicit conversion (converters in this file)
//
//   - BufferBindingType: BindingNotUsed=0 in v29 shifts Undefined/Uniform/Storage/ReadOnlyStorage by +1.
//   - SamplerBindingType: same +1 shift due to BindingNotUsed=0.
//   - TextureSampleType: same +1 shift.
//   - StorageTextureAccess: same +1 shift.
//   - VertexFormat: gputypes omits single-component 8/16-bit variants (Uint8, Sint8,
//     Unorm8, Snorm8, Uint16, Sint16, Unorm16, Snorm16, Float16) added in v29,
//     causing a non-trivial numbering gap throughout the enum.
//   - VertexStepMode: gputypes has VertexBufferNotUsed=1 (removed in v29);
//     v29 maps Vertex=1, Instance=2 instead of gputypes Vertex=2, Instance=3.
//
// # Enums matching v29 exactly — use direct uint32 cast, no converter needed
//
//   - TextureFormat (gputypes v0.3.0 matches v29 exactly, including R16*/RG16*/RGBA16* Unorm/Snorm)
//   - TextureViewDimension, TextureDimension, TextureAspect
//   - LoadOp (Undefined=0, Load=1, Clear=2), StoreOp (Undefined=0, Store=1, Discard=2)
//   - BlendFactor values 0x00–0x0D match v29; gputypes lacks Src1* (0x0E–0x11) but those
//     are unused via gputypes API so no conversion is needed.
//   - BlendOperation, PrimitiveTopology, FrontFace, CullMode
//   - All bitflags: BufferUsage, TextureUsage, ShaderStage, ColorWriteMask, MapMode
//   - FilterMode, MipmapFilterMode, AddressMode, CompareFunction, StencilOperation
//   - IndexFormat, PresentMode, CompositeAlphaMode, PowerPreference

package wgpu

import "github.com/gogpu/gputypes"

// =============================================================================
// BufferBindingType conversion
// gputypes: Undefined=0, Uniform=1, Storage=2, ReadOnlyStorage=3
// wgpu-native v29: BindingNotUsed=0, Undefined=1, Uniform=2, Storage=3, ReadOnlyStorage=4
// Mapping: 0→0 (BindingNotUsed), others +1
// =============================================================================

func toWGPUBufferBindingType(t gputypes.BufferBindingType) uint32 {
	if t == 0 {
		return 0 // gputypes Undefined=0 maps to wgpu BindingNotUsed=0
	}
	return uint32(t) + 1
}

// =============================================================================
// SamplerBindingType conversion
// gputypes: Undefined=0, Filtering=1, NonFiltering=2, Comparison=3
// wgpu-native v29: BindingNotUsed=0, Undefined=1, Filtering=2, NonFiltering=3, Comparison=4
// Mapping: 0→0 (BindingNotUsed), others +1
// =============================================================================

func toWGPUSamplerBindingType(t gputypes.SamplerBindingType) uint32 {
	if t == 0 {
		return 0 // gputypes Undefined=0 maps to wgpu BindingNotUsed=0
	}
	return uint32(t) + 1
}

// =============================================================================
// TextureSampleType conversion
// gputypes: Undefined=0, Float=1, UnfilterableFloat=2, Depth=3, Sint=4, Uint=5
// wgpu-native v29: BindingNotUsed=0, Undefined=1, Float=2, UnfilterableFloat=3, Depth=4, Sint=5, Uint=6
// Mapping: 0→0 (BindingNotUsed), others +1
// =============================================================================

func toWGPUTextureSampleType(t gputypes.TextureSampleType) uint32 {
	if t == 0 {
		return 0 // gputypes Undefined=0 maps to wgpu BindingNotUsed=0
	}
	return uint32(t) + 1
}

// =============================================================================
// StorageTextureAccess conversion
// gputypes: Undefined=0, WriteOnly=1, ReadOnly=2, ReadWrite=3
// wgpu-native v29: BindingNotUsed=0, Undefined=1, WriteOnly=2, ReadOnly=3, ReadWrite=4
// Mapping: 0→0 (BindingNotUsed), others +1
// =============================================================================

func toWGPUStorageTextureAccess(t gputypes.StorageTextureAccess) uint32 {
	if t == 0 {
		return 0 // gputypes Undefined=0 maps to wgpu BindingNotUsed=0
	}
	return uint32(t) + 1
}

// =============================================================================
// VertexStepMode conversion
// gputypes v0.3.0: Undefined=0, VertexBufferNotUsed=1, Vertex=2, Instance=3
// wgpu-native v29: Undefined=0, Vertex=1, Instance=2
//   (VertexBufferNotUsed was removed in v29; Undefined is the sentinel for "not used")
//
// Mapping:
//   gputypes Undefined(0)           → v29 Undefined(0)
//   gputypes VertexBufferNotUsed(1) → v29 Undefined(0)  [removed, treat as not used]
//   gputypes Vertex(2)              → v29 Vertex(1)
//   gputypes Instance(3)            → v29 Instance(2)
// =============================================================================

func toWGPUVertexStepMode(m gputypes.VertexStepMode) uint32 {
	switch m {
	case gputypes.VertexStepModeVertex:
		return 1 // v29 Vertex
	case gputypes.VertexStepModeInstance:
		return 2 // v29 Instance
	default:
		// VertexStepModeUndefined(0) and VertexStepModeVertexBufferNotUsed(1)
		// both map to v29 Undefined(0) — buffer slot not used
		return 0
	}
}

// =============================================================================
// VertexFormat conversion
//
// gputypes v0.3.0 omits single-component 8-bit and 16-bit variants.
// v29 adds: Uint8(1), Sint8(4), Unorm8(7), Snorm8(10), Uint16(13), Sint16(16),
//            Unorm16(19), Snorm16(22), Float16(25), Unorm8x4BGRA(41).
//
// gputypes → v29 mapping (explicit table):
//
//	gputypes  v29   Format
//	     0     0    Undefined
//	     1     2    Uint8x2
//	     2     3    Uint8x4
//	     3     5    Sint8x2
//	     4     6    Sint8x4
//	     5     8    Unorm8x2
//	     6     9    Unorm8x4
//	     7    11    Snorm8x2
//	     8    12    Snorm8x4
//	     9    14    Uint16x2
//	    10    15    Uint16x4
//	    11    17    Sint16x2
//	    12    18    Sint16x4
//	    13    20    Unorm16x2
//	    14    21    Unorm16x4
//	    15    23    Snorm16x2
//	    16    24    Snorm16x4
//	    17    26    Float16x2
//	    18    27    Float16x4
//	    19    28    Float32
//	    20    29    Float32x2
//	    21    30    Float32x3
//	    22    31    Float32x4
//	    23    32    Uint32
//	    24    33    Uint32x2
//	    25    34    Uint32x3
//	    26    35    Uint32x4
//	    27    36    Sint32
//	    28    37    Sint32x2
//	    29    38    Sint32x3
//	    30    39    Sint32x4
//	    31    40    Unorm10_10_10_2  (gputypes: Unorm1010102)
// =============================================================================

func toWGPUVertexFormat(f gputypes.VertexFormat) uint32 {
	switch f {
	case gputypes.VertexFormatUndefined:
		return 0

	// 8-bit packed formats (gputypes lacks single-component variants)
	case gputypes.VertexFormatUint8x2:
		return 2
	case gputypes.VertexFormatUint8x4:
		return 3
	case gputypes.VertexFormatSint8x2:
		return 5
	case gputypes.VertexFormatSint8x4:
		return 6
	case gputypes.VertexFormatUnorm8x2:
		return 8
	case gputypes.VertexFormatUnorm8x4:
		return 9
	case gputypes.VertexFormatSnorm8x2:
		return 11
	case gputypes.VertexFormatSnorm8x4:
		return 12

	// 16-bit packed formats (gputypes lacks single-component variants)
	case gputypes.VertexFormatUint16x2:
		return 14
	case gputypes.VertexFormatUint16x4:
		return 15
	case gputypes.VertexFormatSint16x2:
		return 17
	case gputypes.VertexFormatSint16x4:
		return 18
	case gputypes.VertexFormatUnorm16x2:
		return 20
	case gputypes.VertexFormatUnorm16x4:
		return 21
	case gputypes.VertexFormatSnorm16x2:
		return 23
	case gputypes.VertexFormatSnorm16x4:
		return 24

	// 16-bit float packed formats (gputypes lacks single-component Float16)
	case gputypes.VertexFormatFloat16x2:
		return 26
	case gputypes.VertexFormatFloat16x4:
		return 27

	// 32-bit float formats
	case gputypes.VertexFormatFloat32:
		return 28
	case gputypes.VertexFormatFloat32x2:
		return 29
	case gputypes.VertexFormatFloat32x3:
		return 30
	case gputypes.VertexFormatFloat32x4:
		return 31

	// 32-bit unsigned integer formats
	case gputypes.VertexFormatUint32:
		return 32
	case gputypes.VertexFormatUint32x2:
		return 33
	case gputypes.VertexFormatUint32x3:
		return 34
	case gputypes.VertexFormatUint32x4:
		return 35

	// 32-bit signed integer formats
	case gputypes.VertexFormatSint32:
		return 36
	case gputypes.VertexFormatSint32x2:
		return 37
	case gputypes.VertexFormatSint32x3:
		return 38
	case gputypes.VertexFormatSint32x4:
		return 39

	// Packed normalized format
	case gputypes.VertexFormatUnorm1010102:
		return 40 // v29: Unorm10_10_10_2

	default:
		// Unknown gputypes format value — pass through as-is.
		// This handles future gputypes additions gracefully.
		return uint32(f)
	}
}
