package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// SamplerDescriptor describes a sampler to create.
type SamplerDescriptor struct {
	NextInChain   uintptr
	Label         StringView
	AddressModeU  gputypes.AddressMode
	AddressModeV  gputypes.AddressMode
	AddressModeW  gputypes.AddressMode
	MagFilter     gputypes.FilterMode
	MinFilter     gputypes.FilterMode
	MipmapFilter  gputypes.MipmapFilterMode
	LodMinClamp   float32
	LodMaxClamp   float32
	Compare       gputypes.CompareFunction
	MaxAnisotropy uint16
	_pad          [2]byte
}

// CreateSampler creates a sampler with the specified descriptor.
func (d *Device) CreateSampler(desc *SamplerDescriptor) *Sampler {
	mustInit()
	if desc == nil {
		return nil
	}

	// wgpu-native requires MaxAnisotropy >= 1
	descCopy := *desc
	if descCopy.MaxAnisotropy == 0 {
		descCopy.MaxAnisotropy = 1
	}

	handle, _, _ := procDeviceCreateSampler.Call(
		d.handle,
		uintptr(unsafe.Pointer(&descCopy)),
	)
	if handle == 0 {
		return nil
	}
	trackResource(handle, "Sampler")
	return &Sampler{handle: handle}
}

// CreateLinearSampler creates a sampler with linear filtering.
func (d *Device) CreateLinearSampler() *Sampler {
	desc := SamplerDescriptor{
		Label:        EmptyStringView(),
		AddressModeU: gputypes.AddressModeClampToEdge,
		AddressModeV: gputypes.AddressModeClampToEdge,
		AddressModeW: gputypes.AddressModeClampToEdge,
		MagFilter:    gputypes.FilterModeLinear,
		MinFilter:    gputypes.FilterModeLinear,
		MipmapFilter: gputypes.MipmapFilterModeLinear,
		LodMinClamp:  0.0,
		LodMaxClamp:  32.0,
	}
	return d.CreateSampler(&desc)
}

// CreateNearestSampler creates a sampler with nearest filtering.
func (d *Device) CreateNearestSampler() *Sampler {
	desc := SamplerDescriptor{
		Label:        EmptyStringView(),
		AddressModeU: gputypes.AddressModeClampToEdge,
		AddressModeV: gputypes.AddressModeClampToEdge,
		AddressModeW: gputypes.AddressModeClampToEdge,
		MagFilter:    gputypes.FilterModeNearest,
		MinFilter:    gputypes.FilterModeNearest,
		MipmapFilter: gputypes.MipmapFilterModeNearest,
		LodMinClamp:  0.0,
		LodMaxClamp:  1.0,
	}
	return d.CreateSampler(&desc)
}

// Release releases the sampler reference.
func (s *Sampler) Release() {
	if s.handle != 0 {
		untrackResource(s.handle)
		procSamplerRelease.Call(s.handle) //nolint:errcheck
		s.handle = 0
	}
}

// Handle returns the underlying handle. For advanced use only.
func (s *Sampler) Handle() uintptr { return s.handle }
