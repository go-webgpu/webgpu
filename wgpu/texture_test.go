package wgpu

import (
	"testing"

	"github.com/gogpu/gputypes"
)

func TestCreateTexture(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	t.Log("Creating 2D texture...")
	texture := device.CreateTexture(&TextureDescriptor{
		Usage:     gputypes.TextureUsageTextureBinding | gputypes.TextureUsageCopyDst,
		Dimension: gputypes.TextureDimension2D,
		Size: gputypes.Extent3D{
			Width:              256,
			Height:             256,
			DepthOrArrayLayers: 1,
		},
		Format:        gputypes.TextureFormatRGBA8Unorm,
		MipLevelCount: 1,
		SampleCount:   1,
	})
	if texture == nil {
		t.Fatal("CreateTexture returned nil")
	}
	defer texture.Release()

	if texture.Handle() == 0 {
		t.Fatal("Texture handle is zero")
	}

	t.Logf("Texture created: handle=%#x", texture.Handle())
}

func TestCreateTextureView(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	texture := device.CreateTexture(&TextureDescriptor{
		Usage:     gputypes.TextureUsageTextureBinding,
		Dimension: gputypes.TextureDimension2D,
		Size: gputypes.Extent3D{
			Width:              128,
			Height:             128,
			DepthOrArrayLayers: 1,
		},
		Format:        gputypes.TextureFormatRGBA8Unorm,
		MipLevelCount: 1,
		SampleCount:   1,
	})
	if texture == nil {
		t.Fatal("CreateTexture returned nil")
	}
	defer texture.Release()

	t.Log("Creating texture view...")
	view := texture.CreateView(nil)
	if view == nil {
		t.Fatal("CreateView returned nil")
	}
	defer view.Release()

	if view.Handle() == 0 {
		t.Fatal("TextureView handle is zero")
	}

	t.Logf("TextureView created: handle=%#x", view.Handle())
}

func TestCreateDepthTexture(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	t.Log("Creating depth texture...")
	depthTexture := device.CreateDepthTexture(800, 600, gputypes.TextureFormatDepth24Plus)
	if depthTexture == nil {
		t.Fatal("CreateDepthTexture returned nil")
	}
	defer depthTexture.Release()

	if depthTexture.Handle() == 0 {
		t.Fatal("Depth texture handle is zero")
	}

	t.Logf("Depth texture created: handle=%#x", depthTexture.Handle())
}

func TestCreateSampler(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	t.Log("Creating sampler...")
	sampler := device.CreateSampler(&SamplerDescriptor{
		AddressModeU:  gputypes.AddressModeRepeat,
		AddressModeV:  gputypes.AddressModeRepeat,
		AddressModeW:  gputypes.AddressModeRepeat,
		MagFilter:     gputypes.FilterModeLinear,
		MinFilter:     gputypes.FilterModeLinear,
		MipmapFilter:  gputypes.MipmapFilterModeLinear,
		MaxAnisotropy: 1,
	})
	if sampler == nil {
		t.Fatal("CreateSampler returned nil")
	}
	defer sampler.Release()

	if sampler.Handle() == 0 {
		t.Fatal("Sampler handle is zero")
	}

	t.Logf("Sampler created: handle=%#x", sampler.Handle())
}

func TestCreateSamplerSimple(t *testing.T) {
	inst, err := CreateInstance(nil)
	if err != nil {
		t.Fatalf("CreateInstance failed: %v", err)
	}
	defer inst.Release()

	adapter, err := inst.RequestAdapter(nil)
	if err != nil {
		t.Fatalf("RequestAdapter failed: %v", err)
	}
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	if err != nil {
		t.Fatalf("RequestDevice failed: %v", err)
	}
	defer device.Release()

	t.Log("Creating sampler with minimal settings...")
	sampler := device.CreateSampler(&SamplerDescriptor{
		MaxAnisotropy: 1, // Required to be >= 1
	})
	if sampler == nil {
		t.Fatal("CreateSampler returned nil")
	}
	defer sampler.Release()

	t.Logf("Simple sampler created: handle=%#x", sampler.Handle())
}

func TestTextureFormats(t *testing.T) {
	// Test common texture format constants
	formats := []struct {
		name   string
		format gputypes.TextureFormat
	}{
		{"RGBA8Unorm", gputypes.TextureFormatRGBA8Unorm},
		{"BGRA8Unorm", gputypes.TextureFormatBGRA8Unorm},
		{"Depth24Plus", gputypes.TextureFormatDepth24Plus},
		{"Depth32Float", gputypes.TextureFormatDepth32Float},
		{"R8Unorm", gputypes.TextureFormatR8Unorm},
		{"RG8Unorm", gputypes.TextureFormatRG8Unorm},
	}

	for _, f := range formats {
		if f.format == 0 {
			t.Errorf("Format %s has zero value", f.name)
		}
		t.Logf("Format %s = %#x", f.name, f.format)
	}
}
