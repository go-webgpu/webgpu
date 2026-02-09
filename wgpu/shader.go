package wgpu

import (
	"unsafe"
)

// ShaderModuleDescriptor describes a shader module to create.
type ShaderModuleDescriptor struct {
	NextInChain uintptr    // *ChainedStruct
	Label       StringView // Shader module label for debugging
}

// ShaderSourceWGSL provides WGSL source code for shader creation.
type ShaderSourceWGSL struct {
	Chain ChainedStruct
	Code  StringView
}

// CreateShaderModuleWGSL creates a shader module from WGSL source code.
func (d *Device) CreateShaderModuleWGSL(code string) *ShaderModule {
	mustInit()

	// Create WGSL source with embedded string data
	codeBytes := []byte(code)

	wgslSource := ShaderSourceWGSL{
		Chain: ChainedStruct{
			Next:  0,
			SType: uint32(STypeShaderSourceWGSL),
		},
		Code: StringView{
			Data:   uintptr(unsafe.Pointer(&codeBytes[0])),
			Length: uintptr(len(codeBytes)),
		},
	}

	desc := ShaderModuleDescriptor{
		NextInChain: uintptr(unsafe.Pointer(&wgslSource)),
		Label:       EmptyStringView(),
	}

	handle, _, _ := procDeviceCreateShaderModule.Call(
		d.handle,
		uintptr(unsafe.Pointer(&desc)),
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "ShaderModule")
	return &ShaderModule{handle: handle}
}

// CreateShaderModule creates a shader module from a descriptor.
// For WGSL shaders, use CreateShaderModuleWGSL instead.
func (d *Device) CreateShaderModule(desc *ShaderModuleDescriptor) *ShaderModule {
	mustInit()
	if desc == nil {
		return nil
	}
	handle, _, _ := procDeviceCreateShaderModule.Call(
		d.handle,
		uintptr(unsafe.Pointer(desc)),
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "ShaderModule")
	return &ShaderModule{handle: handle}
}

// Release releases the shader module resources.
func (s *ShaderModule) Release() {
	if s.handle != 0 {
		untrackResource(s.handle)
		procShaderModuleRelease.Call(s.handle) //nolint:errcheck
		s.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (s *ShaderModule) Handle() uintptr { return s.handle }
