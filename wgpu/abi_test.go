package wgpu

// abi_test.go — ABI validation tests for wgpu-native v29 migration.
//
// These tests verify that Go struct layouts and enum values match
// the C ABI defined in wgpu-native v29 webgpu.h and wgpu.h.
//
// All tests:
//   - Require no GPU
//   - Use only unsafe package for layout queries
//   - Are CI-safe: go test -run "TestStruct|TestEnum|TestGputypes|TestWire" ./wgpu/...
//
// v29 BREAKING CHANGES verified here:
//   - WGPULimits gained nextInChain as first field
//   - MinUniform/StorageBufferOffsetAlignment moved after MaxStorageBufferBindingSize
//   - WGPUStatus Success=0x01 (was 0x00 in v27)
//   - WGPUVertexAttribute gained nextInChain (Go wire struct does NOT have it — known gap)
//   - WGPUBindGroupLayoutEntry gained bindingArraySize (Go wire struct does NOT have it — known gap)
//   - WGPUPassTimestampWrites gained nextInChain

import (
	"testing"
	"unsafe"

	"github.com/gogpu/gputypes"
)

// =============================================================================
// TestABIStructSizes — verify Go struct sizes match expected C ABI sizes on 64-bit.
// Expected values derived from wgpu-native v29 webgpu.h and wgpu.h.
// =============================================================================

func TestABIStructSizes(t *testing.T) {
	// ptr = 8, uint32 = 4, uint64 = 8, float64 = 8, WGPUStringView = 16,
	// WGPUBool = uint32 = 4, WGPUFlags = uint64 = 8, enum = uint32 = 4

	tests := []struct {
		name     string
		got      uintptr
		expected uintptr
	}{
		// Core helper types
		{"StringView", unsafe.Sizeof(StringView{}), 16},
		{"ChainedStruct", unsafe.Sizeof(ChainedStruct{}), 16},
		{"Color", unsafe.Sizeof(Color{}), 32},

		// gputypes pass-through (must match C sizes)
		{"gputypes.Extent3D", unsafe.Sizeof(gputypes.Extent3D{}), 12},
		{"gputypes.Origin3D", unsafe.Sizeof(gputypes.Origin3D{}), 12},

		// Instance-level structs
		// InstanceDescriptor: nextInChain(8)+reqFeatureCount(8)+reqFeatures(8)+reqLimits(8) = 32
		{"InstanceDescriptor", unsafe.Sizeof(InstanceDescriptor{}), 32},
		// InstanceLimits: nextInChain(8)+timedWaitAnyMaxCount(8) = 16
		{"InstanceLimits", unsafe.Sizeof(InstanceLimits{}), 16},

		// Adapter-level structs
		// RequestAdapterOptions: nextInChain(8)+featureLevel(4)+powerPreference(4)+
		//   forceFallbackAdapter(4)+backendType(4)+compatibleSurface(8) = 32
		{"RequestAdapterOptions", unsafe.Sizeof(RequestAdapterOptions{}), 32},

		// Limits: nextInChain(8) + 14×uint32(56) + 2×uint64(16) + 2×uint32(8) +
		//   uint32+uint64+9×uint32 + maxImmediateSize(4) = 152
		// Exact layout: see TestStructFieldOffsets for field-level verification.
		{"Limits", unsafe.Sizeof(Limits{}), 152},

		// AdapterInfo: nextInChain(8)+vendor(16)+architecture(16)+device(16)+
		//   description(16)+backendType(4)+adapterType(4)+vendorID(4)+deviceID(4)+
		//   subgroupMinSize(4)+subgroupMaxSize(4) = 96
		{"AdapterInfo", unsafe.Sizeof(AdapterInfo{}), 96},

		// Buffer
		// BufferDescriptor: nextInChain(8)+label(16)+usage(8)+size(8)+mappedAtCreation(4)+pad(4) = 48
		{"BufferDescriptor", unsafe.Sizeof(BufferDescriptor{}), 48},

		// Device-level structs
		// QueueDescriptor: nextInChain(8)+label(16) = 24
		{"QueueDescriptor", unsafe.Sizeof(QueueDescriptor{}), 24},
		// DeviceLostCallbackInfo: nextInChain(8)+mode(4)+pad(4)+callback(8)+ud1(8)+ud2(8) = 40
		{"DeviceLostCallbackInfo", unsafe.Sizeof(DeviceLostCallbackInfo{}), 40},
		// UncapturedErrorCallbackInfo: nextInChain(8)+callback(8)+ud1(8)+ud2(8) = 32
		{"UncapturedErrorCallbackInfo", unsafe.Sizeof(UncapturedErrorCallbackInfo{}), 32},

		// Texture
		// SamplerDescriptor: nextInChain(8)+label(16)+addressModeU(4)+addressModeV(4)+
		//   addressModeW(4)+magFilter(4)+minFilter(4)+mipmapFilter(4)+
		//   lodMinClamp(4)+lodMaxClamp(4)+compare(4)+maxAnisotropy(2)+pad(2) = 64
		{"SamplerDescriptor", unsafe.Sizeof(SamplerDescriptor{}), 64},

		// Pipeline
		// PipelineLayoutDescriptor: nextInChain(8)+label(16)+bindGroupLayoutCount(8)+
		//   bindGroupLayouts(8)+immediateSize(4)+pad(4) = 48
		{"PipelineLayoutDescriptor", unsafe.Sizeof(PipelineLayoutDescriptor{}), 48},
		// ProgrammableStageDescriptor: nextInChain(8)+module(8)+entryPoint(16)+
		//   constantCount(8)+constants(8) = 48
		{"ProgrammableStageDescriptor", unsafe.Sizeof(ProgrammableStageDescriptor{}), 48},
		// ComputePipelineDescriptor: nextInChain(8)+label(16)+layout(8)+compute(48) = 80
		{"ComputePipelineDescriptor", unsafe.Sizeof(ComputePipelineDescriptor{}), 80},

		// BindGroup types
		// BindGroupEntry: nextInChain(8)+binding(4)+pad(4)+buffer(8)+offset(8)+size(8)+
		//   sampler(8)+textureView(8) = 56
		{"BindGroupEntry", unsafe.Sizeof(BindGroupEntry{}), 56},

		// Render pipeline types
		// BlendComponent: operation(4)+srcFactor(4)+dstFactor(4) = 12
		{"BlendComponent", unsafe.Sizeof(BlendComponent{}), 12},
		// BlendState: color(12)+alpha(12) = 24
		{"BlendState", unsafe.Sizeof(BlendState{}), 24},

		// Native extensions
		// NativeLimits: chain(16)+maxImmediateSize(4)+maxNonSamplerBindings(4)+
		//   maxBindingArrayElementsPerShaderStage(4)+pad(4) = 32
		{"NativeLimits", unsafe.Sizeof(NativeLimits{}), 32},
		// PipelineLayoutExtras: chain(16)+immediateDataSize(4)+pad(4) = 24
		{"PipelineLayoutExtras", unsafe.Sizeof(PipelineLayoutExtras{}), 24},

		// Wire structs (FFI-compatible internal structs)
		// passTimestampWrites: nextInChain(8)+querySet(8)+beginIndex(4)+endIndex(4) = 24
		{"passTimestampWrites (wire)", unsafe.Sizeof(passTimestampWrites{}), 24},
		// renderPassColorAttachment: nextInChain(8)+view(8)+depthSlice(4)+pad(4)+
		//   resolveTarget(8)+loadOp(4)+storeOp(4)+clearValue(32) = 72
		{"renderPassColorAttachment (wire)", unsafe.Sizeof(renderPassColorAttachment{}), 72},
		// vertexBufferLayoutWire: nextInChain(8)+stepMode(4)+pad(4)+arrayStride(8)+
		//   attributeCount(8)+attributes(8) = 40
		{"vertexBufferLayoutWire", unsafe.Sizeof(vertexBufferLayoutWire{}), 40},
		// colorTargetStateWire: nextInChain(8)+format(4)+pad(4)+blend(8)+writeMask(8) = 32
		{"colorTargetStateWire (wire)", unsafe.Sizeof(colorTargetStateWire{}), 32},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("sizeof(%s) = %d, want %d (delta: %+d)",
					tt.name, tt.got, tt.expected, int(tt.got)-int(tt.expected))
			}
		})
	}
}

// =============================================================================
// TestABIStructFieldOffsets — verify critical field offsets in FFI structs.
// Especially important for v29 structs where field order changed.
// =============================================================================

func TestABIStructFieldOffsets(t *testing.T) {
	t.Run("Limits", func(t *testing.T) {
		// v29 BREAKING: nextInChain added as first field.
		// v29 BREAKING: minUniform/StorageBufferOffsetAlignment moved
		//   after maxStorageBufferBindingSize (were at end in v27).
		// v29 NEW: maxImmediateSize added as last field.
		//
		// Expected layout (all uint32 unless noted):
		// offset 0:   nextInChain (uintptr/8)
		// offset 8:   maxTextureDimension1D (uint32)
		// offset 12:  maxTextureDimension2D (uint32)
		// offset 16:  maxTextureDimension3D (uint32)
		// offset 20:  maxTextureArrayLayers (uint32)
		// offset 24:  maxBindGroups (uint32)
		// offset 28:  maxBindGroupsPlusVertexBuffers (uint32)
		// offset 32:  maxBindingsPerBindGroup (uint32)
		// offset 36:  maxDynamicUniformBuffersPerPipelineLayout (uint32)
		// offset 40:  maxDynamicStorageBuffersPerPipelineLayout (uint32)
		// offset 44:  maxSampledTexturesPerShaderStage (uint32)
		// offset 48:  maxSamplersPerShaderStage (uint32)
		// offset 52:  maxStorageBuffersPerShaderStage (uint32)
		// offset 56:  maxStorageTexturesPerShaderStage (uint32)
		// offset 60:  maxUniformBuffersPerShaderStage (uint32)  [last of 14 uint32s]
		// --- padding to 8-byte align uint64 ---
		// offset 64:  maxUniformBufferBindingSize (uint64)
		// offset 72:  maxStorageBufferBindingSize (uint64)
		// --- v29 MOVED here (were at end in v27): ---
		// offset 80:  minUniformBufferOffsetAlignment (uint32)
		// offset 84:  minStorageBufferOffsetAlignment (uint32)
		// offset 88:  maxVertexBuffers (uint32)
		// --- padding to 8-byte align uint64 ---
		// offset 96:  maxBufferSize (uint64)
		// offset 104: maxVertexAttributes (uint32)
		// offset 108: maxVertexBufferArrayStride (uint32)
		// offset 112: maxInterStageShaderVariables (uint32)
		// offset 116: maxColorAttachments (uint32)
		// offset 120: maxColorAttachmentBytesPerSample (uint32)
		// offset 124: maxComputeWorkgroupStorageSize (uint32)
		// offset 128: maxComputeInvocationsPerWorkgroup (uint32)
		// offset 132: maxComputeWorkgroupSizeX (uint32)
		// offset 136: maxComputeWorkgroupSizeY (uint32)
		// offset 140: maxComputeWorkgroupSizeZ (uint32)
		// offset 144: maxComputeWorkgroupsPerDimension (uint32)
		// offset 148: maxImmediateSize (uint32)  [NEW v29]
		// total: 152 bytes

		var l Limits
		offsets := []struct {
			name     string
			got      uintptr
			expected uintptr
		}{
			{"NextInChain", unsafe.Offsetof(l.NextInChain), 0},
			{"MaxTextureDimension1D", unsafe.Offsetof(l.MaxTextureDimension1D), 8},
			{"MaxTextureDimension2D", unsafe.Offsetof(l.MaxTextureDimension2D), 12},
			{"MaxTextureDimension3D", unsafe.Offsetof(l.MaxTextureDimension3D), 16},
			{"MaxTextureArrayLayers", unsafe.Offsetof(l.MaxTextureArrayLayers), 20},
			{"MaxBindGroups", unsafe.Offsetof(l.MaxBindGroups), 24},
			{"MaxBindGroupsPlusVertexBuffers", unsafe.Offsetof(l.MaxBindGroupsPlusVertexBuffers), 28},
			{"MaxBindingsPerBindGroup", unsafe.Offsetof(l.MaxBindingsPerBindGroup), 32},
			{"MaxDynamicUniformBuffersPerPipelineLayout", unsafe.Offsetof(l.MaxDynamicUniformBuffersPerPipelineLayout), 36},
			{"MaxDynamicStorageBuffersPerPipelineLayout", unsafe.Offsetof(l.MaxDynamicStorageBuffersPerPipelineLayout), 40},
			{"MaxSampledTexturesPerShaderStage", unsafe.Offsetof(l.MaxSampledTexturesPerShaderStage), 44},
			{"MaxSamplersPerShaderStage", unsafe.Offsetof(l.MaxSamplersPerShaderStage), 48},
			{"MaxStorageBuffersPerShaderStage", unsafe.Offsetof(l.MaxStorageBuffersPerShaderStage), 52},
			{"MaxStorageTexturesPerShaderStage", unsafe.Offsetof(l.MaxStorageTexturesPerShaderStage), 56},
			{"MaxUniformBuffersPerShaderStage", unsafe.Offsetof(l.MaxUniformBuffersPerShaderStage), 60},
			// uint64 fields at 8-byte aligned offsets
			{"MaxUniformBufferBindingSize", unsafe.Offsetof(l.MaxUniformBufferBindingSize), 64},
			{"MaxStorageBufferBindingSize", unsafe.Offsetof(l.MaxStorageBufferBindingSize), 72},
			// v29: these two MOVED from end-of-struct to here
			{"MinUniformBufferOffsetAlignment", unsafe.Offsetof(l.MinUniformBufferOffsetAlignment), 80},
			{"MinStorageBufferOffsetAlignment", unsafe.Offsetof(l.MinStorageBufferOffsetAlignment), 84},
			{"MaxVertexBuffers", unsafe.Offsetof(l.MaxVertexBuffers), 88},
			{"MaxBufferSize", unsafe.Offsetof(l.MaxBufferSize), 96},
			{"MaxVertexAttributes", unsafe.Offsetof(l.MaxVertexAttributes), 104},
			{"MaxVertexBufferArrayStride", unsafe.Offsetof(l.MaxVertexBufferArrayStride), 108},
			{"MaxInterStageShaderVariables", unsafe.Offsetof(l.MaxInterStageShaderVariables), 112},
			{"MaxColorAttachments", unsafe.Offsetof(l.MaxColorAttachments), 116},
			{"MaxColorAttachmentBytesPerSample", unsafe.Offsetof(l.MaxColorAttachmentBytesPerSample), 120},
			{"MaxComputeWorkgroupStorageSize", unsafe.Offsetof(l.MaxComputeWorkgroupStorageSize), 124},
			{"MaxComputeInvocationsPerWorkgroup", unsafe.Offsetof(l.MaxComputeInvocationsPerWorkgroup), 128},
			{"MaxComputeWorkgroupSizeX", unsafe.Offsetof(l.MaxComputeWorkgroupSizeX), 132},
			{"MaxComputeWorkgroupSizeY", unsafe.Offsetof(l.MaxComputeWorkgroupSizeY), 136},
			{"MaxComputeWorkgroupSizeZ", unsafe.Offsetof(l.MaxComputeWorkgroupSizeZ), 140},
			{"MaxComputeWorkgroupsPerDimension", unsafe.Offsetof(l.MaxComputeWorkgroupsPerDimension), 144},
			// v29 NEW field
			{"MaxImmediateSize", unsafe.Offsetof(l.MaxImmediateSize), 148},
		}
		for _, o := range offsets {
			o := o
			t.Run(o.name, func(t *testing.T) {
				if o.got != o.expected {
					t.Errorf("offsetof(Limits.%s) = %d, want %d",
						o.name, o.got, o.expected)
				}
			})
		}
	})

	t.Run("InstanceDescriptor", func(t *testing.T) {
		// nextInChain(0)+requiredFeatureCount(8)+requiredFeatures(16)+requiredLimits(24)
		var d InstanceDescriptor
		offsets := []struct {
			name     string
			got      uintptr
			expected uintptr
		}{
			{"NextInChain", unsafe.Offsetof(d.NextInChain), 0},
			{"RequiredFeatureCount", unsafe.Offsetof(d.RequiredFeatureCount), 8},
			{"RequiredFeatures", unsafe.Offsetof(d.RequiredFeatures), 16},
			{"RequiredLimits", unsafe.Offsetof(d.RequiredLimits), 24},
		}
		for _, o := range offsets {
			o := o
			t.Run(o.name, func(t *testing.T) {
				if o.got != o.expected {
					t.Errorf("offsetof(InstanceDescriptor.%s) = %d, want %d",
						o.name, o.got, o.expected)
				}
			})
		}
	})

	t.Run("RequestAdapterOptions", func(t *testing.T) {
		// nextInChain(0)+featureLevel(8)+powerPreference(12)+
		// forceFallbackAdapter(16)+backendType(20)+compatibleSurface(24) = 32
		var o RequestAdapterOptions
		offsets := []struct {
			name     string
			got      uintptr
			expected uintptr
		}{
			{"NextInChain", unsafe.Offsetof(o.NextInChain), 0},
			{"FeatureLevel", unsafe.Offsetof(o.FeatureLevel), 8},
			{"PowerPreference", unsafe.Offsetof(o.PowerPreference), 12},
			{"ForceFallbackAdapter", unsafe.Offsetof(o.ForceFallbackAdapter), 16},
			{"BackendType", unsafe.Offsetof(o.BackendType), 20},
			{"CompatibleSurface", unsafe.Offsetof(o.CompatibleSurface), 24},
		}
		for _, of := range offsets {
			of := of
			t.Run(of.name, func(t *testing.T) {
				if of.got != of.expected {
					t.Errorf("offsetof(RequestAdapterOptions.%s) = %d, want %d",
						of.name, of.got, of.expected)
				}
			})
		}
	})

	t.Run("AdapterInfo", func(t *testing.T) {
		// nextInChain(0)+vendor(8)+architecture(24)+device(40)+description(56)+
		// backendType(72)+adapterType(76)+vendorID(80)+deviceID(84)+
		// subgroupMinSize(88)+subgroupMaxSize(92) = 96
		var ai AdapterInfo
		offsets := []struct {
			name     string
			got      uintptr
			expected uintptr
		}{
			{"NextInChain", unsafe.Offsetof(ai.NextInChain), 0},
			{"Vendor", unsafe.Offsetof(ai.Vendor), 8},
			{"Architecture", unsafe.Offsetof(ai.Architecture), 24},
			{"Device", unsafe.Offsetof(ai.Device), 40},
			{"Description", unsafe.Offsetof(ai.Description), 56},
			{"BackendType", unsafe.Offsetof(ai.BackendType), 72},
			{"AdapterType", unsafe.Offsetof(ai.AdapterType), 76},
			{"VendorID", unsafe.Offsetof(ai.VendorID), 80},
			{"DeviceID", unsafe.Offsetof(ai.DeviceID), 84},
			{"SubgroupMinSize", unsafe.Offsetof(ai.SubgroupMinSize), 88},
			{"SubgroupMaxSize", unsafe.Offsetof(ai.SubgroupMaxSize), 92},
		}
		for _, o := range offsets {
			o := o
			t.Run(o.name, func(t *testing.T) {
				if o.got != o.expected {
					t.Errorf("offsetof(AdapterInfo.%s) = %d, want %d",
						o.name, o.got, o.expected)
				}
			})
		}
	})

	t.Run("BufferDescriptor", func(t *testing.T) {
		// nextInChain(0)+label(8)+usage(24)+size(32)+mappedAtCreation(40)+pad(44) = 48
		var bd BufferDescriptor
		offsets := []struct {
			name     string
			got      uintptr
			expected uintptr
		}{
			{"NextInChain", unsafe.Offsetof(bd.NextInChain), 0},
			{"Label", unsafe.Offsetof(bd.Label), 8},
			{"Usage", unsafe.Offsetof(bd.Usage), 24},
			{"Size", unsafe.Offsetof(bd.Size), 32},
			{"MappedAtCreation", unsafe.Offsetof(bd.MappedAtCreation), 40},
		}
		for _, o := range offsets {
			o := o
			t.Run(o.name, func(t *testing.T) {
				if o.got != o.expected {
					t.Errorf("offsetof(BufferDescriptor.%s) = %d, want %d",
						o.name, o.got, o.expected)
				}
			})
		}
	})

	t.Run("passTimestampWrites_wire", func(t *testing.T) {
		// v29 NEW: nextInChain added as first field in WGPUPassTimestampWrites.
		// nextInChain(0)+querySet(8)+beginningOfPassWriteIndex(16)+endOfPassWriteIndex(20) = 24
		var ptw passTimestampWrites
		offsets := []struct {
			name     string
			got      uintptr
			expected uintptr
		}{
			{"nextInChain", unsafe.Offsetof(ptw.nextInChain), 0},
			{"querySet", unsafe.Offsetof(ptw.querySet), 8},
			{"beginningOfPassWriteIndex", unsafe.Offsetof(ptw.beginningOfPassWriteIndex), 16},
			{"endOfPassWriteIndex", unsafe.Offsetof(ptw.endOfPassWriteIndex), 20},
		}
		for _, o := range offsets {
			o := o
			t.Run(o.name, func(t *testing.T) {
				if o.got != o.expected {
					t.Errorf("offsetof(passTimestampWrites.%s) = %d, want %d",
						o.name, o.got, o.expected)
				}
			})
		}
	})

	t.Run("renderPassColorAttachment_wire", func(t *testing.T) {
		// nextInChain(0)+view(8)+depthSlice(16)+pad(20)+resolveTarget(24)+
		// loadOp(32)+storeOp(36)+clearValue(40) = 72
		var rca renderPassColorAttachment
		offsets := []struct {
			name     string
			got      uintptr
			expected uintptr
		}{
			{"nextInChain", unsafe.Offsetof(rca.nextInChain), 0},
			{"view", unsafe.Offsetof(rca.view), 8},
			{"depthSlice", unsafe.Offsetof(rca.depthSlice), 16},
			{"resolveTarget", unsafe.Offsetof(rca.resolveTarget), 24},
			{"loadOp", unsafe.Offsetof(rca.loadOp), 32},
			{"storeOp", unsafe.Offsetof(rca.storeOp), 36},
			{"clearValue", unsafe.Offsetof(rca.clearValue), 40},
		}
		for _, o := range offsets {
			o := o
			t.Run(o.name, func(t *testing.T) {
				if o.got != o.expected {
					t.Errorf("offsetof(renderPassColorAttachment.%s) = %d, want %d",
						o.name, o.got, o.expected)
				}
			})
		}
	})
}

// =============================================================================
// TestABIEnumValues — verify Go enum constants match v29 webgpu.h / wgpu.h values.
// =============================================================================

func TestABIEnumValues(t *testing.T) {
	t.Run("FeatureLevel", func(t *testing.T) {
		// v29 added FeatureLevel; Undefined=0x00, Compatibility=0x01, Core=0x02
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"Undefined", uint32(FeatureLevelUndefined), 0x00000000},
			{"Compatibility", uint32(FeatureLevelCompatibility), 0x00000001},
			{"Core", uint32(FeatureLevelCore), 0x00000002},
		}
		runEnumTests(t, tests)
	})

	t.Run("FeatureName", func(t *testing.T) {
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"CoreFeaturesAndLimits", uint32(FeatureNameCoreFeaturesAndLimits), 0x00000001},
			{"DepthClipControl", uint32(FeatureNameDepthClipControl), 0x00000002},
			{"Depth32FloatStencil8", uint32(FeatureNameDepth32FloatStencil8), 0x00000003},
			{"TextureCompressionBC", uint32(FeatureNameTextureCompressionBC), 0x00000004},
			{"TextureCompressionBCSliced3D", uint32(FeatureNameTextureCompressionBCSliced3D), 0x00000005},
			{"TextureCompressionETC2", uint32(FeatureNameTextureCompressionETC2), 0x00000006},
			{"TextureCompressionASTC", uint32(FeatureNameTextureCompressionASTC), 0x00000007},
			{"TextureCompressionASTCSliced3D", uint32(FeatureNameTextureCompressionASTCSliced3D), 0x00000008},
			{"TimestampQuery", uint32(FeatureNameTimestampQuery), 0x00000009},
			{"IndirectFirstInstance", uint32(FeatureNameIndirectFirstInstance), 0x0000000A},
			{"ShaderF16", uint32(FeatureNameShaderF16), 0x0000000B},
			{"RG11B10UfloatRenderable", uint32(FeatureNameRG11B10UfloatRenderable), 0x0000000C},
			{"BGRA8UnormStorage", uint32(FeatureNameBGRA8UnormStorage), 0x0000000D},
			{"Float32Filterable", uint32(FeatureNameFloat32Filterable), 0x0000000E},
			{"Float32Blendable", uint32(FeatureNameFloat32Blendable), 0x0000000F},
			{"ClipDistances", uint32(FeatureNameClipDistances), 0x00000010},
			{"DualSourceBlending", uint32(FeatureNameDualSourceBlending), 0x00000011},
			{"Subgroups", uint32(FeatureNameSubgroups), 0x00000012},
			{"TextureFormatsTier1", uint32(FeatureNameTextureFormatsTier1), 0x00000013},
			{"TextureFormatsTier2", uint32(FeatureNameTextureFormatsTier2), 0x00000014},
			{"PrimitiveIndex", uint32(FeatureNamePrimitiveIndex), 0x00000015},
			{"TextureComponentSwizzle", uint32(FeatureNameTextureComponentSwizzle), 0x00000016},
		}
		runEnumTests(t, tests)
	})

	t.Run("BackendType", func(t *testing.T) {
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"Undefined", uint32(BackendTypeUndefined), 0x00000000},
			{"Null", uint32(BackendTypeNull), 0x00000001},
			{"WebGPU", uint32(BackendTypeWebGPU), 0x00000002},
			{"D3D11", uint32(BackendTypeD3D11), 0x00000003},
			{"D3D12", uint32(BackendTypeD3D12), 0x00000004},
			{"Metal", uint32(BackendTypeMetal), 0x00000005},
			{"Vulkan", uint32(BackendTypeVulkan), 0x00000006},
			{"OpenGL", uint32(BackendTypeOpenGL), 0x00000007},
			{"OpenGLES", uint32(BackendTypeOpenGLES), 0x00000008},
		}
		runEnumTests(t, tests)
	})

	t.Run("AdapterType", func(t *testing.T) {
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"DiscreteGPU", uint32(AdapterTypeDiscreteGPU), 0x00000001},
			{"IntegratedGPU", uint32(AdapterTypeIntegratedGPU), 0x00000002},
			{"CPU", uint32(AdapterTypeCPU), 0x00000003},
			{"Unknown", uint32(AdapterTypeUnknown), 0x00000004},
		}
		runEnumTests(t, tests)
	})

	t.Run("SurfaceGetCurrentTextureStatus", func(t *testing.T) {
		// v29 BREAKING: OutOfMemory(0x06) and DeviceLost(0x07) removed;
		// collapsed to single Error(0x06). Occluded is native extension (0x00030001).
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"SuccessOptimal", uint32(SurfaceGetCurrentTextureStatusSuccessOptimal), 0x00000001},
			{"SuccessSuboptimal", uint32(SurfaceGetCurrentTextureStatusSuccessSuboptimal), 0x00000002},
			{"Timeout", uint32(SurfaceGetCurrentTextureStatusTimeout), 0x00000003},
			{"Outdated", uint32(SurfaceGetCurrentTextureStatusOutdated), 0x00000004},
			{"Lost", uint32(SurfaceGetCurrentTextureStatusLost), 0x00000005},
			{"Error", uint32(SurfaceGetCurrentTextureStatusError), 0x00000006},
			{"Occluded (native)", uint32(NativeSurfaceGetCurrentTextureStatusOccluded), 0x00030001},
		}
		runEnumTests(t, tests)
	})

	t.Run("RequestAdapterStatus", func(t *testing.T) {
		// v29: InstanceDropped renamed to CallbackCancelled
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"Success", uint32(RequestAdapterStatusSuccess), 0x00000001},
			{"CallbackCancelled", uint32(RequestAdapterStatusCallbackCancelled), 0x00000002},
			{"Unavailable", uint32(RequestAdapterStatusUnavailable), 0x00000003},
			{"Error", uint32(RequestAdapterStatusError), 0x00000004},
		}
		runEnumTests(t, tests)
	})

	t.Run("RequestDeviceStatus", func(t *testing.T) {
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"Success", uint32(RequestDeviceStatusSuccess), 0x00000001},
			{"CallbackCancelled", uint32(RequestDeviceStatusCallbackCancelled), 0x00000002},
			{"Error", uint32(RequestDeviceStatusError), 0x00000003},
		}
		runEnumTests(t, tests)
	})

	t.Run("BufferMapState", func(t *testing.T) {
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"Unmapped", uint32(BufferMapStateUnmapped), 0x00000001},
			{"Pending", uint32(BufferMapStatePending), 0x00000002},
			{"Mapped", uint32(BufferMapStateMapped), 0x00000003},
		}
		runEnumTests(t, tests)
	})

	t.Run("WGPUStatus", func(t *testing.T) {
		// v29 BREAKING: Success changed from 0x00 to 0x01; Error from 0x01 to 0x02.
		// Any code comparing against 0 for success is broken after v29.
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"Success", uint32(WGPUStatusSuccess), 0x00000001},
			{"Error", uint32(WGPUStatusError), 0x00000002},
		}
		runEnumTests(t, tests)
	})

	t.Run("SType_standard", func(t *testing.T) {
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"ShaderSourceSPIRV", uint32(STypeShaderSourceSPIRV), 0x00000001},
			{"ShaderSourceWGSL", uint32(STypeShaderSourceWGSL), 0x00000002},
			{"RenderPassMaxDrawCount", uint32(STypeRenderPassMaxDrawCount), 0x00000003},
			{"SurfaceSourceMetalLayer", uint32(STypeSurfaceSourceMetalLayer), 0x00000004},
			{"SurfaceSourceWindowsHWND", uint32(STypeSurfaceSourceWindowsHWND), 0x00000005},
			{"SurfaceSourceXlibWindow", uint32(STypeSurfaceSourceXlibWindow), 0x00000006},
			{"SurfaceSourceWaylandSurface", uint32(STypeSurfaceSourceWaylandSurface), 0x00000007},
			{"SurfaceSourceAndroidNativeWindow", uint32(STypeSurfaceSourceAndroidNativeWindow), 0x00000008},
			{"SurfaceSourceXCBWindow", uint32(STypeSurfaceSourceXCBWindow), 0x00000009},
			{"SurfaceColorManagement", uint32(STypeSurfaceColorManagement), 0x0000000A},
		}
		runEnumTests(t, tests)
	})

	t.Run("SType_native", func(t *testing.T) {
		// Native wgpu-native extension STypes in 0x0003XXXX range
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"DeviceExtras", uint32(STypeDeviceExtras), 0x00030001},
			{"NativeLimits", uint32(STypeNativeLimits), 0x00030002},
			{"PipelineLayoutExtras", uint32(STypePipelineLayoutExtras), 0x00030003},
			{"ShaderSourceGLSL", uint32(STypeShaderSourceGLSL), 0x00030004},
			{"InstanceExtras", uint32(STypeInstanceExtras), 0x00030006},
			{"BindGroupEntryExtras", uint32(STypeBindGroupEntryExtras), 0x00030007},
			{"BindGroupLayoutEntryExtras", uint32(STypeBindGroupLayoutEntryExtras), 0x00030008},
		}
		runEnumTests(t, tests)
	})

	t.Run("NativeFeature", func(t *testing.T) {
		// Selected subset of NativeFeature values; all in 0x0003XXXX range
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"Immediates", uint32(NativeFeatureImmediates), 0x00030001},
			{"TextureAdapterSpecificFormatFeatures", uint32(NativeFeatureTextureAdapterSpecificFormatFeatures), 0x00030002},
			{"MultiDrawIndirectCount", uint32(NativeFeatureMultiDrawIndirectCount), 0x00030004},
			{"VertexWritableStorage", uint32(NativeFeatureVertexWritableStorage), 0x00030005},
			{"TextureBindingArray", uint32(NativeFeatureTextureBindingArray), 0x00030006},
		}
		runEnumTests(t, tests)
	})

	t.Run("InstanceBackend_bitflags", func(t *testing.T) {
		// InstanceBackend is uint64 bitflags
		tests := []struct {
			name     string
			got      uint64
			expected uint64
		}{
			{"All (zero)", uint64(InstanceBackendAll), 0x00000000},
			{"Vulkan", uint64(InstanceBackendVulkan), 1 << 0},
			{"GL", uint64(InstanceBackendGL), 1 << 1},
			{"Metal", uint64(InstanceBackendMetal), 1 << 2},
			{"DX12", uint64(InstanceBackendDX12), 1 << 3},
			{"BrowserWebGPU", uint64(InstanceBackendBrowserWebGPU), 1 << 5},
			// BREAKING v29: Secondary = GL only (v27 had GL|DX11; DX11 removed)
			{"Secondary (GL only)", uint64(InstanceBackendSecondary), 1 << 1},
			{"Primary", uint64(InstanceBackendPrimary),
				(1 << 0) | (1 << 2) | (1 << 3) | (1 << 5)},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				if tt.got != tt.expected {
					t.Errorf("InstanceBackend.%s = %#x, want %#x",
						tt.name, tt.got, tt.expected)
				}
			})
		}
	})

	t.Run("InstanceFlag_bitflags", func(t *testing.T) {
		// InstanceFlag is uint64 bitflags
		// BREAKING v29: Default moved from 0 to 1<<24
		tests := []struct {
			name     string
			got      uint64
			expected uint64
		}{
			{"Empty (zero)", uint64(InstanceFlagEmpty), 0x00000000},
			{"Debug", uint64(InstanceFlagDebug), 1 << 0},
			{"Validation", uint64(InstanceFlagValidation), 1 << 1},
			{"DiscardHalLabels", uint64(InstanceFlagDiscardHalLabels), 1 << 2},
			{"AllowUnderlyingNoncompliantAdapter", uint64(InstanceFlagAllowUnderlyingNoncompliantAdapter), 1 << 3},
			{"GPUBasedValidation", uint64(InstanceFlagGPUBasedValidation), 1 << 4},
			{"ValidationIndirectCall", uint64(InstanceFlagValidationIndirectCall), 1 << 5},
			{"AutomaticTimestampNormalization", uint64(InstanceFlagAutomaticTimestampNormalization), 1 << 6},
			// BREAKING: Default=1<<24 in v29 (was 0 in v27)
			{"Default", uint64(InstanceFlagDefault), 1 << 24},
			{"Debugging", uint64(InstanceFlagDebugging), 1 << 25},
			{"AdvancedDebugging", uint64(InstanceFlagAdvancedDebugging), 1 << 26},
			{"WithEnv", uint64(InstanceFlagWithEnv), 1 << 27},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				if tt.got != tt.expected {
					t.Errorf("InstanceFlag.%s = %#x, want %#x",
						tt.name, tt.got, tt.expected)
				}
			})
		}
	})
}

// runEnumTests is a table-driven helper for uint32 enum value checks.
func runEnumTests(t *testing.T, tests []struct {
	name     string
	got      uint32
	expected uint32
}) {
	t.Helper()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("= %#08x, want %#08x", tt.got, tt.expected)
			}
		})
	}
}

// =============================================================================
// TestABIGputypesEnumAlignment — verify gputypes enum values pass directly through
// FFI without conversion and match v29 webgpu.h constants.
// =============================================================================

func TestABIGputypesEnumAlignment(t *testing.T) {
	t.Run("TextureFormat_common", func(t *testing.T) {
		// gputypes.TextureFormat values used in textureDescriptorWire (passed directly, no conversion).
		// gputypes v0.3.0 uses sequential numbering from the webgpu.h spec where
		// R8Unorm=0x01 through the format list; BGRA8Unorm comes later in the sequence.
		// These expected values are from gputypes v0.3.0 (github.com/gogpu/gputypes).
		// They must remain stable across gputypes versions to avoid silent ABI breaks.
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"Undefined", uint32(gputypes.TextureFormatUndefined), 0x00000000},
			{"R8Unorm", uint32(gputypes.TextureFormatR8Unorm), 0x00000001},
			{"R8Snorm", uint32(gputypes.TextureFormatR8Snorm), 0x00000002},
			{"R8Uint", uint32(gputypes.TextureFormatR8Uint), 0x00000003},
			{"R8Sint", uint32(gputypes.TextureFormatR8Sint), 0x00000004},
			// gputypes v0.3.0 sequential values (after R-only, RG, RGBA formats):
			{"BGRA8Unorm", uint32(gputypes.TextureFormatBGRA8Unorm), 0x0000001b},
			{"BGRA8UnormSrgb", uint32(gputypes.TextureFormatBGRA8UnormSrgb), 0x0000001c},
			{"Depth16Unorm", uint32(gputypes.TextureFormatDepth16Unorm), 0x0000002d},
			{"Depth24Plus", uint32(gputypes.TextureFormatDepth24Plus), 0x0000002e},
			{"Depth32Float", uint32(gputypes.TextureFormatDepth32Float), 0x00000030},
		}
		runEnumTests(t, tests)
	})

	t.Run("TextureUsage_bitflags", func(t *testing.T) {
		// gputypes.TextureUsage bitflags passed directly as uint64 in wire structs.
		// Must match WGPUTextureUsageFlags in webgpu.h v29.
		tests := []struct {
			name     string
			got      uint64
			expected uint64
		}{
			{"None", uint64(gputypes.TextureUsageNone), 0x0000000000000000},
			{"CopySrc", uint64(gputypes.TextureUsageCopySrc), 0x0000000000000001},
			{"CopyDst", uint64(gputypes.TextureUsageCopyDst), 0x0000000000000002},
			{"TextureBinding", uint64(gputypes.TextureUsageTextureBinding), 0x0000000000000004},
			{"StorageBinding", uint64(gputypes.TextureUsageStorageBinding), 0x0000000000000008},
			{"RenderAttachment", uint64(gputypes.TextureUsageRenderAttachment), 0x0000000000000010},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				if tt.got != tt.expected {
					t.Errorf("TextureUsage.%s = %#x, want %#x", tt.name, tt.got, tt.expected)
				}
			})
		}
	})

	t.Run("BufferUsage_bitflags", func(t *testing.T) {
		// gputypes.BufferUsage bitflags.
		// Must match WGPUBufferUsageFlags in webgpu.h v29.
		tests := []struct {
			name     string
			got      uint64
			expected uint64
		}{
			{"None", uint64(gputypes.BufferUsageNone), 0x0000000000000000},
			{"MapRead", uint64(gputypes.BufferUsageMapRead), 0x0000000000000001},
			{"MapWrite", uint64(gputypes.BufferUsageMapWrite), 0x0000000000000002},
			{"CopySrc", uint64(gputypes.BufferUsageCopySrc), 0x0000000000000004},
			{"CopyDst", uint64(gputypes.BufferUsageCopyDst), 0x0000000000000008},
			{"Index", uint64(gputypes.BufferUsageIndex), 0x0000000000000010},
			{"Vertex", uint64(gputypes.BufferUsageVertex), 0x0000000000000020},
			{"Uniform", uint64(gputypes.BufferUsageUniform), 0x0000000000000040},
			{"Storage", uint64(gputypes.BufferUsageStorage), 0x0000000000000080},
			{"Indirect", uint64(gputypes.BufferUsageIndirect), 0x0000000000000100},
			{"QueryResolve", uint64(gputypes.BufferUsageQueryResolve), 0x0000000000000200},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				if tt.got != tt.expected {
					t.Errorf("BufferUsage.%s = %#x, want %#x", tt.name, tt.got, tt.expected)
				}
			})
		}
	})

	t.Run("ShaderStage_bitflags", func(t *testing.T) {
		// gputypes.ShaderStage bitflags.
		// CRITICAL: wgpu-native uses WGPUShaderStageFlags = uint64 (WGPUFlags).
		// These values are widened to uint64 when placed in bindGroupLayoutEntryWire.
		tests := []struct {
			name     string
			got      uint32
			expected uint32
		}{
			{"None", uint32(gputypes.ShaderStageNone), 0x00000000},
			{"Vertex", uint32(gputypes.ShaderStageVertex), 0x00000001},
			{"Fragment", uint32(gputypes.ShaderStageFragment), 0x00000002},
			{"Compute", uint32(gputypes.ShaderStageCompute), 0x00000004},
		}
		runEnumTests(t, tests)
	})
}

// =============================================================================
// TestABIWireStructAlignment — verify wire struct (FFI-compatible) field offsets.
// Wire structs are internal types; their layout must exactly match C ABI.
// =============================================================================

func TestABIWireStructAlignment(t *testing.T) {
	t.Run("vertexBufferLayoutWire", func(t *testing.T) {
		// v29 BREAKING: nextInChain added as FIRST field in WGPUVertexBufferLayout.
		// nextInChain(0)+stepMode(8)+pad(12)+arrayStride(16)+attributeCount(24)+attributes(32) = 40
		var w vertexBufferLayoutWire
		offsets := []struct {
			name     string
			got      uintptr
			expected uintptr
		}{
			{"NextInChain", unsafe.Offsetof(w.NextInChain), 0},
			{"StepMode", unsafe.Offsetof(w.StepMode), 8},
			{"ArrayStride", unsafe.Offsetof(w.ArrayStride), 16},
			{"AttributeCount", unsafe.Offsetof(w.AttributeCount), 24},
			{"Attributes", unsafe.Offsetof(w.Attributes), 32},
		}
		for _, o := range offsets {
			o := o
			t.Run(o.name, func(t *testing.T) {
				if o.got != o.expected {
					t.Errorf("offsetof(vertexBufferLayoutWire.%s) = %d, want %d",
						o.name, o.got, o.expected)
				}
			})
		}
	})

	t.Run("vertexAttributeWire_size", func(t *testing.T) {
		// v29 STATUS: WGPUVertexAttribute in C v29 has nextInChain as first field (32 bytes).
		// Our vertexAttributeWire does NOT have nextInChain (24 bytes).
		//
		// This is a KNOWN MIGRATION GAP:
		//   C v29 WGPUVertexAttribute:
		//     nextInChain(8)+format(4)+pad(4)+offset(8)+shaderLocation(4)+pad(4) = 32 bytes
		//   Go vertexAttributeWire (current):
		//     format(4)+pad(4)+offset(8)+shaderLocation(4)+pad(4) = 24 bytes  [MISSING nextInChain]
		//
		// TODO(v29-migration): Add nextInChain to vertexAttributeWire when upgrading to wgpu-native v29.
		// Tracked in: docs/dev/kanban/blocked/0010-webgpu-headers-upgrade.md
		const gotSize = unsafe.Sizeof(vertexAttributeWire{})
		const expectedCurrent = uintptr(24) // current Go wire (no nextInChain)
		const expectedV29C = uintptr(32)    // C v29 target (has nextInChain)

		if gotSize != expectedCurrent {
			t.Errorf("sizeof(vertexAttributeWire) = %d, want %d (current Go layout)",
				gotSize, expectedCurrent)
		}
		// Document the gap: once v29 migration is complete, this must be 32.
		if gotSize == expectedV29C {
			t.Log("vertexAttributeWire already matches C v29 size (32 bytes) — remove migration TODO")
		} else {
			t.Logf("MIGRATION GAP: vertexAttributeWire is %d bytes, C v29 target is %d bytes (missing nextInChain)",
				gotSize, expectedV29C)
		}
	})

	t.Run("bindGroupLayoutEntryWire_knownGap", func(t *testing.T) {
		// v29 STATUS: WGPUBindGroupLayoutEntry in C v29 has bindingArraySize (uint32)
		// between visibility (uint64) and buffer (bufferBindingLayoutWire).
		//
		// This is a KNOWN MIGRATION GAP:
		//   C v29 layout after visibility:
		//     bindingArraySize(4)+pad(4)+buffer(...)+sampler(...)+...
		//   Go bindGroupLayoutEntryWire (current):
		//     NO bindingArraySize field between visibility and buffer
		//
		// Impact: buffer, sampler, texture, storageTexture offsets are all shifted
		// by -8 relative to C v29. This will cause incorrect binding when binding arrays
		// are used (NativeFeatureTextureBindingArray).
		//
		// TODO(v29-migration): Add bindingArraySize uint32 + padding after Visibility
		// in bindGroupLayoutEntryWire when upgrading to wgpu-native v29.
		// Tracked in: docs/dev/kanban/blocked/0010-webgpu-headers-upgrade.md

		var e bindGroupLayoutEntryWire
		// Verify current layout is self-consistent (no accidental regressions)
		visibilityOffset := unsafe.Offsetof(e.Visibility)
		bufferOffset := uintptr(unsafe.Pointer(&e.Buffer)) - uintptr(unsafe.Pointer(&e))

		// Current: visibility at some offset, buffer directly after (no bindingArraySize gap)
		// In C v29: buffer should be at visibility+8+8 = visibility+16 (bindingArraySize+pad)
		// Currently buffer is at visibility+8 (just uint64 visibility, no bindingArraySize)
		expectedCurrentGap := uintptr(8) // sizeof(Visibility uint64) = 8, buffer follows directly
		actualGap := bufferOffset - visibilityOffset
		if actualGap != expectedCurrentGap {
			t.Errorf("gap(Visibility→Buffer) = %d bytes, want %d (current layout without bindingArraySize)",
				actualGap, expectedCurrentGap)
		}
		t.Logf("MIGRATION GAP: C v29 expects gap(Visibility→Buffer)=16 bytes (bindingArraySize+pad), current Go has %d bytes",
			actualGap)
	})

	t.Run("colorTargetStateWire", func(t *testing.T) {
		// nextInChain(0)+format(8)+pad(12)+blend(16)+writeMask(24) = 32
		// CRITICAL: writeMask is uint64 (WGPUColorWriteMaskFlags = WGPUFlags = uint64)
		var w colorTargetStateWire
		offsets := []struct {
			name     string
			got      uintptr
			expected uintptr
		}{
			{"nextInChain", unsafe.Offsetof(w.nextInChain), 0},
			{"format", unsafe.Offsetof(w.format), 8},
			{"blend", unsafe.Offsetof(w.blend), 16},
			{"writeMask", unsafe.Offsetof(w.writeMask), 24},
		}
		for _, o := range offsets {
			o := o
			t.Run(o.name, func(t *testing.T) {
				if o.got != o.expected {
					t.Errorf("offsetof(colorTargetStateWire.%s) = %d, want %d",
						o.name, o.got, o.expected)
				}
			})
		}
	})

	t.Run("bindGroupLayoutEntryWire_visibility_uint64", func(t *testing.T) {
		// CRITICAL: Visibility must be uint64 (WGPUShaderStageFlags = WGPUFlags = uint64 in wgpu-native).
		// This is NOT uint32 as in the webgpu.h spec — wgpu-native uses WGPUFlags typedef.
		// Verify the Visibility field size via its offset and the next field offset.
		var e bindGroupLayoutEntryWire
		visibilityOffset := unsafe.Offsetof(e.Visibility)
		bufferOffset := uintptr(unsafe.Pointer(&e.Buffer)) - uintptr(unsafe.Pointer(&e))
		visibilitySize := bufferOffset - visibilityOffset
		const expectedVisibilitySize = uintptr(8) // must be uint64 = 8 bytes
		if visibilitySize != expectedVisibilitySize {
			t.Errorf("sizeof(Visibility in bindGroupLayoutEntryWire) = %d, want %d (must be uint64)",
				visibilitySize, expectedVisibilitySize)
		}
	})
}
