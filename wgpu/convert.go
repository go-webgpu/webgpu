// convert.go provides conversion functions between gputypes (webgpu.h spec)
// and wgpu-native internal values.
//
// Background: gputypes follows an older webgpu.h schema where enums start at 0.
// wgpu-native v24+ uses a newer schema with BindingNotUsed=0, shifting other values by +1.
// TextureFormat also differs due to removal of R16Unorm/R16Snorm from the spec.

package wgpu

import "github.com/gogpu/gputypes"

// =============================================================================
// BufferBindingType conversion
// gputypes: Undefined=0, Uniform=1, Storage=2, ReadOnlyStorage=3
// wgpu-native: BindingNotUsed=0, Undefined=1, Uniform=2, Storage=3, ReadOnlyStorage=4
// =============================================================================

func toWGPUBufferBindingType(t gputypes.BufferBindingType) uint32 {
	// gputypes: Undefined=0, Uniform=1, Storage=2, ReadOnlyStorage=3
	// wgpu-native: BindingNotUsed=0, Undefined=1, Uniform=2, Storage=3, ReadOnlyStorage=4
	// Keep 0 as 0 (BindingNotUsed), shift others by +1
	if t == 0 {
		return 0 // BindingNotUsed
	}
	return uint32(t) + 1
}

// =============================================================================
// SamplerBindingType conversion
// gputypes: Undefined=0, Filtering=1, NonFiltering=2, Comparison=3
// wgpu-native: BindingNotUsed=0, Undefined=1, Filtering=2, NonFiltering=3, Comparison=4
// =============================================================================

func toWGPUSamplerBindingType(t gputypes.SamplerBindingType) uint32 {
	// gputypes: Undefined=0, Filtering=1, NonFiltering=2, Comparison=3
	// wgpu-native: BindingNotUsed=0, Undefined=1, Filtering=2, NonFiltering=3, Comparison=4
	// Keep 0 as 0 (BindingNotUsed), shift others by +1
	if t == 0 {
		return 0 // BindingNotUsed
	}
	return uint32(t) + 1
}

// =============================================================================
// TextureSampleType conversion
// gputypes: Undefined=0, Float=1, UnfilterableFloat=2, Depth=3, Sint=4, Uint=5
// wgpu-native: BindingNotUsed=0, Undefined=1, Float=2, UnfilterableFloat=3, Depth=4, Sint=5, Uint=6
// =============================================================================

func toWGPUTextureSampleType(t gputypes.TextureSampleType) uint32 {
	// gputypes: Undefined=0, Float=1, UnfilterableFloat=2, Depth=3, Sint=4, Uint=5
	// wgpu-native: BindingNotUsed=0, Undefined=1, Float=2, UnfilterableFloat=3, Depth=4, Sint=5, Uint=6
	// Keep 0 as 0 (BindingNotUsed), shift others by +1
	if t == 0 {
		return 0 // BindingNotUsed
	}
	return uint32(t) + 1
}

// =============================================================================
// TextureViewDimension conversion
// gputypes: Undefined=0, 1D=1, 2D=2, 2DArray=3, Cube=4, CubeArray=5, 3D=6
// wgpu-native v27 (bac5208): Undefined=0, 1D=1, 2D=2, 2DArray=3, Cube=4, CubeArray=5, 3D=6
// Values match! No conversion needed.
// =============================================================================

func toWGPUTextureViewDimension(t gputypes.TextureViewDimension) uint32 {
	// Values match between gputypes and wgpu-native v27 - no conversion needed
	return uint32(t)
}

// =============================================================================
// StorageTextureAccess conversion
// gputypes: Undefined=0, WriteOnly=1, ReadOnly=2, ReadWrite=3
// wgpu-native: BindingNotUsed=0, Undefined=1, WriteOnly=2, ReadOnly=3, ReadWrite=4
// =============================================================================

func toWGPUStorageTextureAccess(t gputypes.StorageTextureAccess) uint32 {
	// gputypes: Undefined=0, WriteOnly=1, ReadOnly=2, ReadWrite=3
	// wgpu-native: BindingNotUsed=0, Undefined=1, WriteOnly=2, ReadOnly=3, ReadWrite=4
	// Keep 0 as 0 (BindingNotUsed), shift others by +1
	if t == 0 {
		return 0 // BindingNotUsed
	}
	return uint32(t) + 1
}

// =============================================================================
// TextureFormat conversion
// gputypes follows older webgpu.h with R16Unorm/R16Snorm/RG16Unorm/RG16Snorm (4 formats).
// wgpu-native v24+ uses newer spec where these were moved to extensions.
// Result: formats after R8Sint are shifted by -2, and after RG8Sint by another -2.
// =============================================================================

func toWGPUTextureFormat(f gputypes.TextureFormat) uint32 {
	// gputypes values (old spec with R16/RG16 Unorm/Snorm in core):
	//   0=Undefined, 1-4=R8*, 5-6=R16Unorm/Snorm, 7-9=R16Uint/Sint/Float,
	//   10-13=RG8*, 14-16=R32*, 17-18=RG16Unorm/Snorm, 19-21=RG16Uint/Sint/Float,
	//   22-26=RGBA8*, 27-28=BGRA8*, 29=RGB10A2Uint, 30=RGB10A2Unorm,
	//   31=RG11B10Ufloat, 32=RGB9E5Ufloat, 33-35=RG32*, 36-37=RGBA16Unorm/Snorm,
	//   38-40=RGBA16Uint/Sint/Float, 41-43=RGBA32*, 44=Stencil8, 45=Depth16Unorm,
	//   46=Depth24Plus, 47=Depth24PlusStencil8, 48=Depth32Float, 49=Depth32FloatStencil8
	//
	// webgpu-headers values (new spec - R16/RG16 Unorm/Snorm moved to extensions):
	//   0=Undefined, 1-4=R8*, 5-7=R16Uint/Sint/Float (NO R16Unorm/Snorm!),
	//   8-11=RG8*, 12-14=R32*, 15-17=RG16Uint/Sint/Float (NO RG16Unorm/Snorm!),
	//   18-22=RGBA8*, 23-24=BGRA8*, 25=RGB10A2Uint, 26=RGB10A2Unorm,
	//   27=RG11B10Ufloat, 28=RGB9E5Ufloat, 29-31=RG32*, 32-34=RGBA16Uint/Sint/Float,
	//   (NO RGBA16Unorm/Snorm!), 35-37=RGBA32*, 38=Stencil8, 39=Depth16Unorm,
	//   40=Depth24Plus, 41=Depth24PlusStencil8, 42=Depth32Float, 43=Depth32FloatStencil8

	// Use a lookup table for common formats
	switch f {
	case gputypes.TextureFormatUndefined:
		return 0

	// 8-bit R formats (1-4 → 1-4, same)
	case gputypes.TextureFormatR8Unorm:
		return 1
	case gputypes.TextureFormatR8Snorm:
		return 2
	case gputypes.TextureFormatR8Uint:
		return 3
	case gputypes.TextureFormatR8Sint:
		return 4

	// R16 formats: gputypes has Unorm/Snorm at 5-6, wgpu-native doesn't
	// R16Uint/Sint/Float: gputypes 7-9 → wgpu-native 5-7
	case gputypes.TextureFormatR16Uint:
		return 5
	case gputypes.TextureFormatR16Sint:
		return 6
	case gputypes.TextureFormatR16Float:
		return 7

	// RG8 formats: gputypes 10-13 → wgpu-native 8-11
	case gputypes.TextureFormatRG8Unorm:
		return 8
	case gputypes.TextureFormatRG8Snorm:
		return 9
	case gputypes.TextureFormatRG8Uint:
		return 10
	case gputypes.TextureFormatRG8Sint:
		return 11

	// R32 formats: gputypes 14-16 → wgpu-native 12-14
	case gputypes.TextureFormatR32Float:
		return 12
	case gputypes.TextureFormatR32Uint:
		return 13
	case gputypes.TextureFormatR32Sint:
		return 14

	// RG16 formats: gputypes has Unorm/Snorm at 17-18, wgpu-native doesn't
	// RG16Uint/Sint/Float: gputypes 19-21 → wgpu-native 15-17
	case gputypes.TextureFormatRG16Uint:
		return 15
	case gputypes.TextureFormatRG16Sint:
		return 16
	case gputypes.TextureFormatRG16Float:
		return 17

	// RGBA8 formats: gputypes 22-26 → wgpu-native 18-22
	case gputypes.TextureFormatRGBA8Unorm:
		return 18
	case gputypes.TextureFormatRGBA8UnormSrgb:
		return 19
	case gputypes.TextureFormatRGBA8Snorm:
		return 20
	case gputypes.TextureFormatRGBA8Uint:
		return 21
	case gputypes.TextureFormatRGBA8Sint:
		return 22

	// BGRA8 formats: gputypes 27-28 → wgpu-native 23-24
	case gputypes.TextureFormatBGRA8Unorm:
		return 23
	case gputypes.TextureFormatBGRA8UnormSrgb:
		return 24

	// Packed formats: gputypes 29-32 → wgpu-native 25-28
	case gputypes.TextureFormatRGB10A2Uint:
		return 25
	case gputypes.TextureFormatRGB10A2Unorm:
		return 26
	case gputypes.TextureFormatRG11B10Ufloat:
		return 27
	case gputypes.TextureFormatRGB9E5Ufloat:
		return 28

	// RG32 formats: gputypes 33-35 → wgpu-native 29-31
	case gputypes.TextureFormatRG32Float:
		return 29
	case gputypes.TextureFormatRG32Uint:
		return 30
	case gputypes.TextureFormatRG32Sint:
		return 31

	// RGBA16 formats: gputypes has Unorm/Snorm at 36-37, wgpu-native doesn't
	// RGBA16Uint/Sint/Float: gputypes 38-40 → wgpu-native 32-34
	case gputypes.TextureFormatRGBA16Uint:
		return 32
	case gputypes.TextureFormatRGBA16Sint:
		return 33
	case gputypes.TextureFormatRGBA16Float:
		return 34

	// RGBA32 formats: gputypes 41-43 → wgpu-native 35-37
	case gputypes.TextureFormatRGBA32Float:
		return 35
	case gputypes.TextureFormatRGBA32Uint:
		return 36
	case gputypes.TextureFormatRGBA32Sint:
		return 37

	// Depth/Stencil formats: gputypes 44-49 → wgpu-native 38-43
	case gputypes.TextureFormatStencil8:
		return 38
	case gputypes.TextureFormatDepth16Unorm:
		return 39
	case gputypes.TextureFormatDepth24Plus:
		return 40
	case gputypes.TextureFormatDepth24PlusStencil8:
		return 41
	case gputypes.TextureFormatDepth32Float:
		return 42
	case gputypes.TextureFormatDepth32FloatStencil8:
		return 43

	default:
		// For compressed formats and others, try simple mapping
		// Compressed formats start at higher values and may need individual mapping
		return uint32(f)
	}
}

// fromWGPUTextureFormat converts a wgpu-native TextureFormat value to gputypes.
// This is the reverse of toWGPUTextureFormat.
func fromWGPUTextureFormat(f uint32) gputypes.TextureFormat {
	switch f {
	case 0:
		return gputypes.TextureFormatUndefined

	// 8-bit R formats (1-4 → 1-4, same)
	case 1:
		return gputypes.TextureFormatR8Unorm
	case 2:
		return gputypes.TextureFormatR8Snorm
	case 3:
		return gputypes.TextureFormatR8Uint
	case 4:
		return gputypes.TextureFormatR8Sint

	// R16 formats: wgpu-native 5-7 → gputypes 7-9
	case 5:
		return gputypes.TextureFormatR16Uint
	case 6:
		return gputypes.TextureFormatR16Sint
	case 7:
		return gputypes.TextureFormatR16Float

	// RG8 formats: wgpu-native 8-11 → gputypes 10-13
	case 8:
		return gputypes.TextureFormatRG8Unorm
	case 9:
		return gputypes.TextureFormatRG8Snorm
	case 10:
		return gputypes.TextureFormatRG8Uint
	case 11:
		return gputypes.TextureFormatRG8Sint

	// R32 formats: wgpu-native 12-14 → gputypes 14-16
	case 12:
		return gputypes.TextureFormatR32Float
	case 13:
		return gputypes.TextureFormatR32Uint
	case 14:
		return gputypes.TextureFormatR32Sint

	// RG16 formats: wgpu-native 15-17 → gputypes 19-21
	case 15:
		return gputypes.TextureFormatRG16Uint
	case 16:
		return gputypes.TextureFormatRG16Sint
	case 17:
		return gputypes.TextureFormatRG16Float

	// RGBA8 formats: wgpu-native 18-22 → gputypes 22-26
	case 18:
		return gputypes.TextureFormatRGBA8Unorm
	case 19:
		return gputypes.TextureFormatRGBA8UnormSrgb
	case 20:
		return gputypes.TextureFormatRGBA8Snorm
	case 21:
		return gputypes.TextureFormatRGBA8Uint
	case 22:
		return gputypes.TextureFormatRGBA8Sint

	// BGRA8 formats: wgpu-native 23-24 → gputypes 27-28
	case 23:
		return gputypes.TextureFormatBGRA8Unorm
	case 24:
		return gputypes.TextureFormatBGRA8UnormSrgb

	// Packed formats: wgpu-native 25-28 → gputypes 29-32
	case 25:
		return gputypes.TextureFormatRGB10A2Uint
	case 26:
		return gputypes.TextureFormatRGB10A2Unorm
	case 27:
		return gputypes.TextureFormatRG11B10Ufloat
	case 28:
		return gputypes.TextureFormatRGB9E5Ufloat

	// RG32 formats: wgpu-native 29-31 → gputypes 33-35
	case 29:
		return gputypes.TextureFormatRG32Float
	case 30:
		return gputypes.TextureFormatRG32Uint
	case 31:
		return gputypes.TextureFormatRG32Sint

	// RGBA16 formats: wgpu-native 32-34 → gputypes 38-40
	case 32:
		return gputypes.TextureFormatRGBA16Uint
	case 33:
		return gputypes.TextureFormatRGBA16Sint
	case 34:
		return gputypes.TextureFormatRGBA16Float

	// RGBA32 formats: wgpu-native 35-37 → gputypes 41-43
	case 35:
		return gputypes.TextureFormatRGBA32Float
	case 36:
		return gputypes.TextureFormatRGBA32Uint
	case 37:
		return gputypes.TextureFormatRGBA32Sint

	// Depth/Stencil formats: wgpu-native 38-43 → gputypes 44-49
	case 38:
		return gputypes.TextureFormatStencil8
	case 39:
		return gputypes.TextureFormatDepth16Unorm
	case 40:
		return gputypes.TextureFormatDepth24Plus
	case 41:
		return gputypes.TextureFormatDepth24PlusStencil8
	case 42:
		return gputypes.TextureFormatDepth32Float
	case 43:
		return gputypes.TextureFormatDepth32FloatStencil8

	default:
		// For unknown/compressed formats, return as-is
		return gputypes.TextureFormat(f)
	}
}

// =============================================================================
// Types that DON'T need conversion (bitflags or already matching)
// =============================================================================

// ShaderStage - bitflags, same values (but needs widening to uint64!)
// BufferUsage - bitflags, same values
// TextureUsage - bitflags, same values
// ColorWriteMask - bitflags, same values

// =============================================================================
// Simple enums that may need +1 shift (to be verified)
// =============================================================================

// PrimitiveTopology, IndexFormat, FrontFace, CullMode, VertexFormat, VertexStepMode,
// AddressMode, FilterMode, CompareFunction, BlendFactor, BlendOperation, StencilOperation,
// LoadOp, StoreOp, PresentMode, CompositeAlphaMode, TextureDimension
//
// These may or may not have BindingNotUsed/Undefined shift - needs verification.
// For now, we'll add converters as issues are discovered.

func toWGPULoadOp(op gputypes.LoadOp) uint32 {
	// gputypes: Undefined=0, Clear=1, Load=2
	// wgpu-native: Undefined=0, Load=1, Clear=2 (different order!)
	switch op {
	case gputypes.LoadOpClear:
		return 2 // wgpu-native Clear
	case gputypes.LoadOpLoad:
		return 1 // wgpu-native Load
	default:
		return 0 // Undefined
	}
}

func toWGPUStoreOp(op gputypes.StoreOp) uint32 {
	// gputypes: Undefined=0, Store=1, Discard=2
	// wgpu-native: Undefined=0, Store=1, Discard=2 (same!)
	return uint32(op)
}

// =============================================================================
// TextureDimension conversion
// gputypes: Undefined=0, 1D=1, 2D=2, 3D=3
// wgpu-native: Undefined=1, 1D=2, 2D=3, 3D=4
// =============================================================================

func toWGPUTextureDimension(d gputypes.TextureDimension) uint32 {
	// gputypes: Undefined=0, 1D=1, 2D=2, 3D=3
	// wgpu-native: Undefined=0, 1D=1, 2D=2, 3D=3 (SAME values!)
	// TextureDimension does NOT have BindingNotUsed, so no +1 shift needed
	return uint32(d)
}

// =============================================================================
// VertexStepMode conversion
// gputypes: Undefined=0, VertexBufferNotUsed=1, Vertex=2, Instance=3
// wgpu-native: VertexBufferNotUsed=0, Undefined=1, Vertex=2, Instance=3
// Note: Undefined and VertexBufferNotUsed are SWAPPED!
// =============================================================================

func toWGPUVertexStepMode(m gputypes.VertexStepMode) uint32 {
	switch m {
	case gputypes.VertexStepModeUndefined:
		return 1 // wgpu-native Undefined
	case gputypes.VertexStepModeVertexBufferNotUsed:
		return 0 // wgpu-native VertexBufferNotUsed
	default:
		// Vertex=2, Instance=3 are the same
		return uint32(m)
	}
}

// =============================================================================
// VertexFormat conversion
// gputypes has fewer formats (no single-component 8/16-bit).
// wgpu-native has Uint8, Sint8, Unorm8, Snorm8, Uint16, Sint16, Unorm16, Snorm16, Float16
// which shift all subsequent values.
//
// gputypes: Uint8x2=1, Uint8x4=2, Sint8x2=3, Sint8x4=4, Unorm8x2=5, Unorm8x4=6...
// wgpu-native: Uint8=1, Uint8x2=2, Uint8x4=3, Sint8=4, Sint8x2=5, Sint8x4=6...
// =============================================================================

func toWGPUVertexFormat(f gputypes.VertexFormat) uint32 {
	switch f {
	case gputypes.VertexFormatUndefined:
		return 0

	// 8-bit formats: gputypes lacks single-component
	case gputypes.VertexFormatUint8x2:
		return 2 // wgpu Uint8x2
	case gputypes.VertexFormatUint8x4:
		return 3 // wgpu Uint8x4
	case gputypes.VertexFormatSint8x2:
		return 5 // wgpu Sint8x2
	case gputypes.VertexFormatSint8x4:
		return 6 // wgpu Sint8x4
	case gputypes.VertexFormatUnorm8x2:
		return 8 // wgpu Unorm8x2
	case gputypes.VertexFormatUnorm8x4:
		return 9 // wgpu Unorm8x4
	case gputypes.VertexFormatSnorm8x2:
		return 11 // wgpu Snorm8x2
	case gputypes.VertexFormatSnorm8x4:
		return 12 // wgpu Snorm8x4

	// 16-bit formats: gputypes lacks single-component
	case gputypes.VertexFormatUint16x2:
		return 14 // wgpu Uint16x2
	case gputypes.VertexFormatUint16x4:
		return 15 // wgpu Uint16x4
	case gputypes.VertexFormatSint16x2:
		return 17 // wgpu Sint16x2
	case gputypes.VertexFormatSint16x4:
		return 18 // wgpu Sint16x4
	case gputypes.VertexFormatUnorm16x2:
		return 20 // wgpu Unorm16x2
	case gputypes.VertexFormatUnorm16x4:
		return 21 // wgpu Unorm16x4
	case gputypes.VertexFormatSnorm16x2:
		return 23 // wgpu Snorm16x2
	case gputypes.VertexFormatSnorm16x4:
		return 24 // wgpu Snorm16x4

	// Float16 formats: gputypes lacks single-component
	case gputypes.VertexFormatFloat16x2:
		return 26 // wgpu Float16x2
	case gputypes.VertexFormatFloat16x4:
		return 27 // wgpu Float16x4

	// Float32 formats
	case gputypes.VertexFormatFloat32:
		return 28 // wgpu Float32
	case gputypes.VertexFormatFloat32x2:
		return 29 // wgpu Float32x2
	case gputypes.VertexFormatFloat32x3:
		return 30 // wgpu Float32x3
	case gputypes.VertexFormatFloat32x4:
		return 31 // wgpu Float32x4

	// Uint32 formats
	case gputypes.VertexFormatUint32:
		return 32 // wgpu Uint32
	case gputypes.VertexFormatUint32x2:
		return 33 // wgpu Uint32x2
	case gputypes.VertexFormatUint32x3:
		return 34 // wgpu Uint32x3
	case gputypes.VertexFormatUint32x4:
		return 35 // wgpu Uint32x4

	// Sint32 formats
	case gputypes.VertexFormatSint32:
		return 36 // wgpu Sint32
	case gputypes.VertexFormatSint32x2:
		return 37 // wgpu Sint32x2
	case gputypes.VertexFormatSint32x3:
		return 38 // wgpu Sint32x3
	case gputypes.VertexFormatSint32x4:
		return 39 // wgpu Sint32x4

	// Packed format
	case gputypes.VertexFormatUnorm1010102:
		return 40 // wgpu Unorm10_10_10_2

	default:
		return uint32(f)
	}
}
