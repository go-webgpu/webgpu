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
// Returns an error if the FFI call fails or the device is nil.
func (d *Device) CreateShaderModuleWGSL(code string) (*ShaderModule, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if d == nil || d.handle == 0 {
		return nil, &WGPUError{Op: "CreateShaderModuleWGSL", Message: "device is nil or released"}
	}
	if code == "" {
		return nil, &WGPUError{Op: "CreateShaderModuleWGSL", Message: "shader source is empty"}
	}

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
		return nil, &WGPUError{Op: "CreateShaderModuleWGSL", Message: "wgpu returned null handle"}
	}
	trackResource(handle, "ShaderModule")
	return &ShaderModule{handle: handle}, nil
}

// CreateShaderModule creates a shader module from a descriptor.
// For WGSL shaders, prefer CreateShaderModuleWGSL or use ShaderDescriptor.
// Returns an error if the FFI call fails or the device/descriptor is nil.
func (d *Device) CreateShaderModule(desc *ShaderModuleDescriptor) (*ShaderModule, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if d == nil || d.handle == 0 {
		return nil, &WGPUError{Op: "CreateShaderModule", Message: "device is nil or released"}
	}
	if desc == nil {
		return nil, &WGPUError{Op: "CreateShaderModule", Message: "descriptor is nil"}
	}
	handle, _, _ := procDeviceCreateShaderModule.Call(
		d.handle,
		uintptr(unsafe.Pointer(desc)),
	)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreateShaderModule", Message: "wgpu returned null handle"}
	}
	trackResource(handle, "ShaderModule")
	return &ShaderModule{handle: handle}, nil
}

// CreateShaderModuleFromDescriptor creates a shader module from a Go-idiomatic ShaderDescriptor.
// If both WGSL and SPIRV are set, WGSL takes precedence.
// Returns an error if the FFI call fails, the device is nil, or both sources are empty.
func (d *Device) CreateShaderModuleFromDescriptor(desc *ShaderDescriptor) (*ShaderModule, error) {
	if desc == nil {
		return nil, &WGPUError{Op: "CreateShaderModule", Message: "descriptor is nil"}
	}
	if desc.WGSL != "" {
		return d.CreateShaderModuleWGSL(desc.WGSL)
	}
	if len(desc.SPIRV) > 0 {
		return d.CreateShaderModuleSPIRV(desc.Label, desc.SPIRV)
	}
	return nil, &WGPUError{Op: "CreateShaderModule", Message: "shader source is empty (set WGSL or SPIRV)"}
}

// CreateShaderModuleSPIRV creates a shader module from SPIR-V bytecode.
// Returns an error if the FFI call fails, the device is nil, or the bytecode is empty.
func (d *Device) CreateShaderModuleSPIRV(label string, spirv []uint32) (*ShaderModule, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if d == nil || d.handle == 0 {
		return nil, &WGPUError{Op: "CreateShaderModuleSPIRV", Message: "device is nil or released"}
	}
	if len(spirv) == 0 {
		return nil, &WGPUError{Op: "CreateShaderModuleSPIRV", Message: "SPIR-V bytecode is empty"}
	}

	spirvSource := struct {
		Chain    ChainedStruct
		Code     uintptr // *uint32
		CodeSize uint32
		_pad     [4]byte
	}{
		Chain: ChainedStruct{
			Next:  0,
			SType: uint32(STypeShaderSourceSPIRV),
		},
		Code:     uintptr(unsafe.Pointer(&spirv[0])),
		CodeSize: uint32(len(spirv)),
	}

	desc := ShaderModuleDescriptor{
		NextInChain: uintptr(unsafe.Pointer(&spirvSource)),
		Label:       stringToStringView(label),
	}

	handle, _, _ := procDeviceCreateShaderModule.Call(
		d.handle,
		uintptr(unsafe.Pointer(&desc)),
	)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreateShaderModuleSPIRV", Message: "wgpu returned null handle"}
	}
	trackResource(handle, "ShaderModule")
	return &ShaderModule{handle: handle}, nil
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
