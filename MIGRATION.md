# Migration Guide: v0.4.x → v0.5.0

v0.5.0 is a major breaking release that upgrades wgpu-native from v27 to v29 (stable webgpu-headers) and redesigns the public API for idiomatic Go usage. This guide covers every breaking change with before/after examples.

See [CHANGELOG.md](CHANGELOG.md) for the full list of changes.

---

## Table of Contents

- [wgpu-native Binary](#wgpu-native-binary)
- [Error Returns on Create Methods](#error-returns-on-create-methods)
- [Method Renames](#method-renames)
- [Single Import: gputypes Aliases](#single-import-gputypes-aliases)
- [Removed Types](#removed-types)
- [Enum Changes](#enum-changes)
- [Quick Checklist](#quick-checklist)

---

## wgpu-native Binary

Download the new binary from the [wgpu-native releases page](https://github.com/gfx-rs/wgpu-native/releases).

| Platform | Filename |
|----------|---------|
| Windows x64 | `wgpu-windows-x86_64-msvc-release.zip` → `wgpu_native.dll` |
| Linux x64 | `wgpu-linux-x86_64-release.zip` → `libwgpu_native.so` |
| macOS ARM64 | `wgpu-macos-aarch64-release.zip` → `libwgpu_native.dylib` |

Required version: **v29.0.0.0**. The v0.5.0 Go bindings are not compatible with v27.

---

## Error Returns on Create Methods

All `Create*` methods now return `(*T, error)` instead of `*T`. This applies to every resource creation function.

```go
// Before (v0.4.x) — nil on failure
buffer := device.CreateBuffer(&desc)
if buffer == nil {
    log.Fatal("failed to create buffer")
}
defer buffer.Release()

// After (v0.5.0) — idiomatic error return
buffer, err := device.CreateBuffer(&desc)
if err != nil {
    log.Fatal(err)
}
defer buffer.Release()
```

Affected methods (all on `*Device` unless noted):

| Method | v0.4.x return | v0.5.0 return |
|--------|--------------|--------------|
| `CreateBuffer` | `*Buffer` | `(*Buffer, error)` |
| `CreateTexture` | `*Texture` | `(*Texture, error)` |
| `CreateShaderModuleWGSL` | `*ShaderModule` | `(*ShaderModule, error)` |
| `CreateRenderPipeline` | `*RenderPipeline` | `(*RenderPipeline, error)` |
| `CreateComputePipeline` | `*ComputePipeline` | `(*ComputePipeline, error)` |
| `CreateBindGroup` | `*BindGroup` | `(*BindGroup, error)` |
| `CreateBindGroupLayout` | `*BindGroupLayout` | `(*BindGroupLayout, error)` |
| `CreatePipelineLayout` | `*PipelineLayout` | `(*PipelineLayout, error)` |
| `CreateCommandEncoder` | `*CommandEncoder` | `(*CommandEncoder, error)` |
| `CreateSampler` | `*Sampler` | `(*Sampler, error)` |
| `CreateQuerySet` | `*QuerySet` | `(*QuerySet, error)` |
| `Texture.CreateView` | `*TextureView` | `(*TextureView, error)` |
| `CommandEncoder.Finish` | `*CommandBuffer` | `(*CommandBuffer, error)` |

---

## Method Renames

`Get` prefix removed from accessor methods to follow Go naming conventions:

```go
// Before (v0.4.x)          // After (v0.5.0)
device.GetQueue()           → device.Queue()
buffer.GetSize()            → buffer.Size()
adapter.GetLimits()         → adapter.Limits()
adapter.GetInfo()           → adapter.Info()
device.GetLimits()          → device.Limits()
device.GetFeatures()        → device.Features()
texture.GetWidth()          → texture.Width()
texture.GetHeight()         → texture.Height()
texture.GetDepthOrArrayLayers() → texture.DepthOrArrayLayers()
texture.GetFormat()         → texture.Format()
texture.GetDimension()      → texture.Dimension()
texture.GetUsage()          → texture.Usage()
```

Additionally, `Adapter.EnumerateFeatures()` is now `Adapter.Features()`. The old name is kept as a deprecated alias until v1.0.

---

## Single Import: gputypes Aliases

v0.5.0 re-exports gputypes types and constants as aliases in the `wgpu` package. A separate `gputypes` import is no longer required.

```go
// Before (v0.4.x) — two imports required
import (
    "github.com/go-webgpu/webgpu/wgpu"
    "github.com/gogpu/gputypes"
)

surfaceConfig.Format = gputypes.TextureFormatBGRA8Unorm
bufferDesc.Usage     = gputypes.BufferUsageVertex | gputypes.BufferUsageCopyDst
```

```go
// After (v0.5.0) — single import
import "github.com/go-webgpu/webgpu/wgpu"

surfaceConfig.Format = wgpu.TextureFormatBGRA8Unorm
bufferDesc.Usage     = wgpu.BufferUsageVertex | wgpu.BufferUsageCopyDst
```

The direct `gputypes` import still works and produces the same types (they are aliases, not copies). Use the direct import when sharing types with other gogpu ecosystem packages.

---

## Removed Types

### SupportedLimits

`SupportedLimits` wrapper struct is removed. `Limits` is now returned directly.

```go
// Before (v0.4.x)
supported, err := adapter.GetLimits()
maxBuffers := supported.Limits.MaxVertexBuffers

// After (v0.5.0)
limits, err := adapter.Limits()
maxBuffers := limits.MaxVertexBuffers
```

### ChainedStructOut

`ChainedStructOut` is aliased to `ChainedStruct` (v29 unified them). If your code references `ChainedStructOut`, replace with `ChainedStruct`.

### InstanceCapabilities

`InstanceCapabilities` struct is removed. Use `GetInstanceFeatures()` if you need instance-level feature queries.

---

## Enum Changes

### SurfaceGetCurrentTextureStatus

Status values simplified. `SurfaceGetCurrentTextureStatusOutOfMemory` and `SurfaceGetCurrentTextureStatusDeviceLost` are removed. Check `WGPUError` from the error return of `CreateTexture` instead.

### InstanceFlag

`InstanceFlag_Default` semantic changed in v29. Use explicit flags:

```go
// Before (v0.4.x)
desc := wgpu.InstanceDescriptor{Flags: wgpu.InstanceFlagDefault}

// After (v0.5.0)
desc := wgpu.InstanceDescriptor{Flags: wgpu.InstanceFlagNone}
// or with validation:
desc := wgpu.InstanceDescriptor{Flags: wgpu.InstanceFlagDebugging}
```

### DX11 Backend Removed

`InstanceBackendDX11` is removed from `InstanceBackend` flags. wgpu-native v29 uses D3D12 on Windows.

---

## Quick Checklist

After updating to v0.5.0:

- [ ] Download wgpu-native v29.0.0.0 for your platform
- [ ] Update Go module: `go get github.com/go-webgpu/webgpu@v0.5.0`
- [ ] Add `err` return handling to all `Create*` calls
- [ ] Replace `GetQueue()` → `Queue()`, `GetLimits()` → `Limits()`, `GetInfo()` → `Info()`
- [ ] Replace `GetSize()` → `Size()` on Buffer
- [ ] Replace `GetWidth()`/`GetHeight()` → `Width()`/`Height()` on Texture
- [ ] Remove separate `gputypes` import if no longer needed
- [ ] Replace `supported.Limits.X` → `limits.X` (SupportedLimits removed)
- [ ] Replace `ChainedStructOut` → `ChainedStruct`
- [ ] Remove any reference to `InstanceBackendDX11`
- [ ] Run `go build ./...` and fix remaining compilation errors
