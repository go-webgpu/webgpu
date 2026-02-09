# API Stability Policy

This document describes the API stability guarantees for go-webgpu.

## Versioning

This project follows [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html)
and [Go module versioning](https://go.dev/doc/modules/version-numbers).

- **v0.x.y** — Pre-stable. API may change between minor versions.
- **v1.0.0** — First stable release. Breaking changes require a new major version.

## Current Status: Pre-stable (v0.x)

While in v0.x, the API is evolving. However, we follow these principles:

- **Patch versions (v0.x.Y)** — Bug fixes, dependency updates. No API changes.
- **Minor versions (v0.X.0)** — New features, possible breaking changes (documented in CHANGELOG).

## API Surface Classification

### Stable API

These APIs are unlikely to change and will follow the deprecation policy below:

| Category | Examples |
|----------|---------|
| **Instance lifecycle** | `CreateInstance`, `Instance.Release`, `Instance.ProcessEvents` |
| **Adapter** | `Instance.RequestAdapter`, `Adapter.GetLimits`, `Adapter.GetInfo`, `Adapter.Release` |
| **Device** | `Adapter.RequestDevice`, `Device.GetQueue`, `Device.Release` |
| **Buffer** | `Device.CreateBuffer`, `Buffer.MapAsync`, `Buffer.GetMappedRange`, `Buffer.Unmap`, `Buffer.Release` |
| **Texture** | `Device.CreateTexture`, `Texture.CreateView`, `Texture.Release`, `TextureView.Release` |
| **Shader** | `Device.CreateShaderModuleWGSL`, `ShaderModule.Release` |
| **Pipeline** | `Device.CreateRenderPipeline`, `Device.CreateComputePipeline`, `*Pipeline.Release` |
| **Bind Group** | `Device.CreateBindGroup`, `Device.CreateBindGroupLayout`, `*.Release` |
| **Command** | `Device.CreateCommandEncoder`, `CommandEncoder.Finish`, `Queue.Submit` |
| **Render Pass** | `CommandEncoder.BeginRenderPass`, `RenderPassEncoder.*`, `RenderPassEncoder.End` |
| **Compute Pass** | `CommandEncoder.BeginComputePass`, `ComputePassEncoder.*`, `ComputePassEncoder.End` |
| **Surface** | `Surface.Configure`, `Surface.GetCurrentTexture`, `Surface.Present`, `Surface.Release` |
| **Error handling** | `WGPUError`, `ErrValidation`, `ErrOutOfMemory`, `ErrInternal`, `ErrDeviceLost` |
| **Debug** | `SetDebugMode`, `ReportLeaks`, `ResetLeakTracker` |
| **Types** | All types in `gputypes` package (external dependency) |

### Experimental API

These APIs may change in minor versions:

| API | Reason | Stability target |
|-----|--------|-----------------|
| `Device.GetLimits` | Returns error on some wgpu-native versions | v1.0 |
| `Device.GetFeatures` | Uses newer wgpuDeviceGetFeatures API | v1.0 |
| `Surface.GetCapabilities` | New in v0.3 | v1.0 |
| `*Simple` convenience methods | May be renamed or adjusted | v1.0 |
| Math helpers (`Mat4`, `Vec3`) | May move to separate package | v1.0 |

### Internal API (not for external use)

| API | Purpose |
|-----|---------|
| `Handle()` methods | Raw FFI handle access — no stability guarantee |
| Wire structs (`*wire`, `*Wire`) | FFI layout structs — change with wgpu-native |
| `Init()` / `mustInit()` / `checkInit()` | Library initialization internals |

## Deprecation Policy

When an API needs to change:

1. **Announce** — The function is marked with `// Deprecated: Use X instead.`
   This follows the [Go deprecation convention](https://go.dev/wiki/Deprecated).

2. **Maintain** — The deprecated function continues to work for at least one minor version.

3. **Remove** — In the next major version (or next minor for v0.x), the function is removed.

### Currently Deprecated

| Function | Deprecated in | Replacement | Removal target |
|----------|--------------|-------------|----------------|
| `Device.PopErrorScope` | v0.3.0 | `Device.PopErrorScopeAsync` | v1.0.0 |

## Breaking Change Policy

For v0.x releases, breaking changes are documented in [CHANGELOG.md](CHANGELOG.md)
with migration guides. We minimize breaking changes and batch them into minor releases.

For v1.0+, breaking changes will only occur in major version bumps.

## Compatibility with wgpu-native

This library tracks wgpu-native releases. When wgpu-native makes breaking changes:

1. We absorb the changes in our conversion layer (`convert.go`) when possible.
2. If user-facing API must change, it follows the deprecation policy above.
3. The supported wgpu-native version is documented in `go.mod` and README.

Current: **wgpu-native v24.0.3.1** (v27 API)
