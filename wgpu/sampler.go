package wgpu

import (
	"unsafe"

	"github.com/gogpu/gputypes"
)

// SamplerDescriptor describes a sampler to create.
type SamplerDescriptor struct {
	Label         string
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
}

// samplerDescriptorWire is the FFI-compatible C-layout struct for wgpu-native.
// CRITICAL: layout must match WGPUSamplerDescriptor exactly.
// nextInChain(8)+label(16)+addressModeU(4)+addressModeV(4)+addressModeW(4)+
// magFilter(4)+minFilter(4)+mipmapFilter(4)+lodMinClamp(4)+lodMaxClamp(4)+
// compare(4)+maxAnisotropy(2)+pad(2) = 64 bytes.
type samplerDescriptorWire struct {
	NextInChain   uintptr              // 8 bytes
	Label         StringView           // 16 bytes
	AddressModeU  gputypes.AddressMode // 4 bytes
	AddressModeV  gputypes.AddressMode // 4 bytes
	AddressModeW  gputypes.AddressMode // 4 bytes
	MagFilter     gputypes.FilterMode  // 4 bytes
	MinFilter     gputypes.FilterMode  // 4 bytes
	MipmapFilter  gputypes.MipmapFilterMode // 4 bytes
	LodMinClamp   float32              // 4 bytes
	LodMaxClamp   float32              // 4 bytes
	Compare       gputypes.CompareFunction // 4 bytes
	MaxAnisotropy uint16               // 2 bytes
	_pad          [2]byte              //nolint:unused // padding to align to 4 bytes
}

// CreateSampler creates a sampler with the specified descriptor.
func (d *Device) CreateSampler(desc *SamplerDescriptor) (*Sampler, error) {
	if err := checkInit(); err != nil {
		return nil, err
	}
	if d == nil || d.handle == 0 {
		return nil, &WGPUError{Op: "CreateSampler", Message: "device is nil or released"}
	}
	if desc == nil {
		return nil, &WGPUError{Op: "CreateSampler", Message: "descriptor is nil"}
	}

	// wgpu-native requires MaxAnisotropy >= 1
	maxAnisotropy := desc.MaxAnisotropy
	if maxAnisotropy == 0 {
		maxAnisotropy = 1
	}

	wire := samplerDescriptorWire{
		Label:         stringToStringView(desc.Label),
		AddressModeU:  desc.AddressModeU,
		AddressModeV:  desc.AddressModeV,
		AddressModeW:  desc.AddressModeW,
		MagFilter:     desc.MagFilter,
		MinFilter:     desc.MinFilter,
		MipmapFilter:  desc.MipmapFilter,
		LodMinClamp:   desc.LodMinClamp,
		LodMaxClamp:   desc.LodMaxClamp,
		Compare:       desc.Compare,
		MaxAnisotropy: maxAnisotropy,
	}

	handle, _, _ := procDeviceCreateSampler.Call(
		d.handle,
		uintptr(unsafe.Pointer(&wire)),
	)
	if handle == 0 {
		return nil, &WGPUError{Op: "CreateSampler", Message: "wgpu returned null handle"}
	}
	trackResource(handle, "Sampler")
	return &Sampler{handle: handle}, nil
}

// CreateLinearSampler creates a sampler with linear filtering.
func (d *Device) CreateLinearSampler() (*Sampler, error) {
	return d.CreateSampler(&SamplerDescriptor{
		AddressModeU: gputypes.AddressModeClampToEdge,
		AddressModeV: gputypes.AddressModeClampToEdge,
		AddressModeW: gputypes.AddressModeClampToEdge,
		MagFilter:    gputypes.FilterModeLinear,
		MinFilter:    gputypes.FilterModeLinear,
		MipmapFilter: gputypes.MipmapFilterModeLinear,
		LodMinClamp:  0.0,
		LodMaxClamp:  32.0,
	})
}

// CreateNearestSampler creates a sampler with nearest filtering.
func (d *Device) CreateNearestSampler() (*Sampler, error) {
	return d.CreateSampler(&SamplerDescriptor{
		AddressModeU: gputypes.AddressModeClampToEdge,
		AddressModeV: gputypes.AddressModeClampToEdge,
		AddressModeW: gputypes.AddressModeClampToEdge,
		MagFilter:    gputypes.FilterModeNearest,
		MinFilter:    gputypes.FilterModeNearest,
		MipmapFilter: gputypes.MipmapFilterModeNearest,
		LodMinClamp:  0.0,
		LodMaxClamp:  1.0,
	})
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
