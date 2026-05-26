package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// =============================================================================
// Helper functions for public → wire conversion
// =============================================================================

// stringToStringView converts a Go string to a wgpu-native StringView.
// The returned StringView points to the string's backing data — the string
// must remain alive for the duration of any FFI call using the result.
func stringToStringView(s string) StringView {
	if len(s) == 0 {
		return EmptyStringView()
	}
	b := []byte(s)
	return StringView{
		Data:   uintptr(unsafe.Pointer(&b[0])),
		Length: uintptr(len(b)),
	}
}

// boolToWGPU converts a Go bool to a WebGPU WGPUBool (uint32).
func boolToWGPU(b bool) Bool {
	if b {
		return True
	}
	return False
}


// =============================================================================
// Shader descriptor
// =============================================================================

// ShaderDescriptor is the Go-idiomatic descriptor for creating a shader module.
// Use WGSL to specify WGSL source, or SPIRV for SPIR-V bytecode.
// If both are set, WGSL takes precedence.
// Use Device.CreateShaderModule or Device.CreateShaderModuleFromDescriptor to create from this.
type ShaderDescriptor struct {
	Label string
	WGSL  string   // WGSL source code
	SPIRV []uint32 // SPIR-V bytecode (alternative to WGSL)
}

// ShaderModuleDescriptorGo is an alias for ShaderDescriptor.
// Use this name when matching the gogpu/wgpu naming convention.
type ShaderModuleDescriptorGo = ShaderDescriptor

// =============================================================================
// Render bundle descriptor
// =============================================================================

// RenderBundleEncoderDescriptorGo is kept for backward compatibility.
// Deprecated: Use RenderBundleEncoderDescriptor directly (now Go-idiomatic).
type RenderBundleEncoderDescriptorGo = RenderBundleEncoderDescriptor

// =============================================================================
// Legacy types (kept for backward compatibility, may be removed in v0.6.0)
// =============================================================================

// BindGroupEntryGo is kept for backward compatibility.
// Deprecated: Use BindGroupEntry directly (now Go-idiomatic).
type BindGroupEntryGo = BindGroupEntry

// BindGroupDescriptorGo is kept for backward compatibility.
// Deprecated: Use BindGroupDescriptor directly (now Go-idiomatic).
type BindGroupDescriptorGo = BindGroupDescriptor

// BindGroupLayoutDescriptorGo is kept for backward compatibility.
// Deprecated: Use BindGroupLayoutDescriptor directly (now Go-idiomatic).
type BindGroupLayoutDescriptorGo = BindGroupLayoutDescriptor

// PipelineLayoutDescriptorGo is kept for backward compatibility.
// Deprecated: Use PipelineLayoutDescriptor directly (now Go-idiomatic).
type PipelineLayoutDescriptorGo = PipelineLayoutDescriptor

// ComputePipelineDescriptorGo is kept for backward compatibility.
// Deprecated: Use ComputePipelineDescriptor directly (now Go-idiomatic).
type ComputePipelineDescriptorGo = ComputePipelineDescriptor

// CommandEncoderDescriptorGo is kept for backward compatibility.
// Deprecated: Use CommandEncoderDescriptor directly (now Go-idiomatic).
type CommandEncoderDescriptorGo = CommandEncoderDescriptor

// =============================================================================
// Suppressing unused import warning
// =============================================================================

var _ = gputypes.TextureFormatUndefined // ensure gputypes import is used
