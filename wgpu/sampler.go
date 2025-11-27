package wgpu

import "unsafe"

// SamplerDescriptor describes a sampler to create.
type SamplerDescriptor struct {
	NextInChain   uintptr
	Label         StringView
	AddressModeU  AddressMode
	AddressModeV  AddressMode
	AddressModeW  AddressMode
	MagFilter     FilterMode
	MinFilter     FilterMode
	MipmapFilter  MipmapFilterMode
	LodMinClamp   float32
	LodMaxClamp   float32
	Compare       CompareFunction
	MaxAnisotropy uint16
	_pad          [2]byte
}

// CreateSampler creates a sampler with the specified descriptor.
func (d *Device) CreateSampler(desc *SamplerDescriptor) *Sampler {
	mustInit()
	if desc == nil {
		return nil
	}
	handle, _, _ := procDeviceCreateSampler.Call(
		d.handle,
		uintptr(unsafe.Pointer(desc)),
	)
	if handle == 0 {
		return nil
	}
	return &Sampler{handle: handle}
}

// CreateLinearSampler creates a sampler with linear filtering.
func (d *Device) CreateLinearSampler() *Sampler {
	desc := SamplerDescriptor{
		Label:        EmptyStringView(),
		AddressModeU: AddressModeClampToEdge,
		AddressModeV: AddressModeClampToEdge,
		AddressModeW: AddressModeClampToEdge,
		MagFilter:    FilterModeLinear,
		MinFilter:    FilterModeLinear,
		MipmapFilter: MipmapFilterModeLinear,
		LodMinClamp:  0.0,
		LodMaxClamp:  32.0,
	}
	return d.CreateSampler(&desc)
}

// CreateNearestSampler creates a sampler with nearest filtering.
func (d *Device) CreateNearestSampler() *Sampler {
	desc := SamplerDescriptor{
		Label:        EmptyStringView(),
		AddressModeU: AddressModeClampToEdge,
		AddressModeV: AddressModeClampToEdge,
		AddressModeW: AddressModeClampToEdge,
		MagFilter:    FilterModeNearest,
		MinFilter:    FilterModeNearest,
		MipmapFilter: MipmapFilterModeNearest,
		LodMinClamp:  0.0,
		LodMaxClamp:  1.0,
	}
	return d.CreateSampler(&desc)
}

// Release releases the sampler reference.
func (s *Sampler) Release() {
	if s.handle != 0 {
		procSamplerRelease.Call(s.handle) //nolint:errcheck
		s.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (s *Sampler) Handle() uintptr { return s.handle }
