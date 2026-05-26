# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v0.5.0 (Unreleased)

### Breaking Changes
- **wgpu-native v29.0.0.0**: Migrated from v27.0.4.0 to v29.0.0.0 with stable webgpu-headers
- **API redesign**: All `Create*` methods now return `(*T, error)` instead of `*T`
- **Method renames**: `GetQueue()` → `Queue()`, `GetSize()` → `Size()`, `GetLimits()` → `Limits()`, etc.
- **Struct layout changes**: `Limits` field order fixed (ABI-breaking), `VertexBufferLayout` gains `nextInChain`
- **Removed types**: `SupportedLimits`, `ChainedStructOut`, `InstanceCapabilities`
- **Enum changes**: `SurfaceGetCurrentTextureStatus` simplified, `InstanceFlag_Default` semantic change
- **gputypes aliases**: Types re-exported for single-import ergonomics

### Added
- 271 enterprise ABI verification tests (`TestABI*`)
- gputypes type aliases and constant re-exports in `wgpu` package
- New v29 API functions: `GetFeatures`, `GetInstanceFeatures`, `BufferReadMappedRange`, etc.
- New enums: `InstanceFeatureName`, `ComponentSwizzle`, `PredefinedColorSpace`, `ToneMappingMode`
- New instance flags: `GPUBasedValidation`, `Debugging`, `AdvancedDebugging`, `WithEnv`

### Changed
- `convert.go`: Removed 6 identity converters (TextureFormat now matches v29 natively)
- `wgpuAdapterEnumerateFeatures` → `wgpuAdapterGetFeatures` (single-call pattern)
- PushConstants → Immediates rename throughout

### Removed
- `SupportedLimits` wrapper struct
- `ChainedStructOut` type (aliased to `ChainedStruct`)
- `InstanceCapabilities` struct
- DX11 backend support (`InstanceBackendDX11`)
- `SurfaceGetCurrentTextureStatusOutOfMemory`, `SurfaceGetCurrentTextureStatusDeviceLost`

### Dependencies
- wgpu-native: v27.0.4.0 → v29.0.0.0
- goffi: v0.5.0 (unchanged)
- gputypes: v0.3.0 (unchanged)

---

## [0.4.3] - 2026-03-29

### Changed

- Update goffi v0.4.2 → v0.5.0 — Windows ARM64 (Snapdragon X) and FreeBSD amd64 support
- Update gputypes v0.2.0 → v0.3.0 — `TextureUsage.ContainsUnknownBits()` method
- Update golang.org/x/sys v0.41.0 → v0.42.0

---

## [0.4.2] - 2026-03-04

### Changed

- Update goffi v0.4.1 → v0.4.2 — purego compatibility fix (`nofakecgo` build tag for `_cgo_init` linker collision)
- Update golang.org/x/sys v0.40.0 → v0.41.0

---

## [0.4.1] - 2026-03-02

### Changed

- Update goffi v0.4.0 → v0.4.1 — ABI compliance hotfix (float32 encoding, stack spill for 7+ args, struct return, runtime.KeepAlive)

---

## [0.4.0] - 2026-02-27

### Added

- **Null handle guards** on all public FFI methods — prevents SIGSEGV when passing nil/released objects
- **85 null guard tests** (`TestNullGuard_*`) — CI-safe, no GPU required
- **`WGPU_NATIVE_PATH` env var** — override library path for custom wgpu-native locations
- **`ptrFromUintptr` helper** — eliminates all `go vet` unsafe.Pointer warnings in FFI code

### Changed

- `loadLibrary` now returns `(Library, error)` — proper error propagation on init failure
- Windows: eager DLL loading via `dll.Load()` — errors at `Init()` instead of first FFI call
- `Init()` returns descriptive error messages with library path and override hint
- CI: wgpu-native binary downloaded in all workflows — tests run against real library, no skips
- CI: removed `-unsafeptr=false` go vet workaround — all warnings properly fixed

### Fixed

- **15 `go vet` warnings** — all `possible misuse of unsafe.Pointer` eliminated via `ptrFromUintptr`
- Silent library loading failures — `Init()` now properly reports missing DLL/so/dylib

---

## [0.3.2] - 2026-02-27

### Changed

- **goffi:** v0.3.9 → v0.4.0 (crosscall2 integration for C-thread callbacks, fakecgo trampoline fixes)

---

## [0.3.1] - 2026-02-18

### Changed

- **goffi:** v0.3.8 → v0.3.9 (ARM64 callback trampoline fix, symbol rename to avoid linker collision)

---

## [0.3.0] - 2026-02-09

### Added

- **Surface.GetCapabilities()** — query supported formats, present modes, alpha modes
- **Device.GetFeatures()** — enumerate all features enabled on the device
- **Device.HasFeature()** — check if a specific feature is enabled
- **Device.GetLimits()** — retrieve device limits (experimental, may return error)
- **Typed error system** — `WGPUError` with `errors.Is()`/`errors.As()` support
- **Sentinel errors** — `ErrValidation`, `ErrOutOfMemory`, `ErrInternal`, `ErrDeviceLost`
- **Resource leak detection** — `SetDebugMode(true)`, `ReportLeaks()`, zero overhead when disabled
- **Thread safety documentation** — `doc.go` with threading model, safe/unsafe operations
- **Fuzz testing** — 14 fuzz targets for FFI boundary (enum conversions, struct sizes, math)
- **API stability policy** — `STABILITY.md` with stable/experimental classification
- **Comprehensive godoc** — all exported symbols documented for pkg.go.dev
- **Release automation** — GitHub Actions workflow for automated release on tag push

### Changed

- Error-returning functions now use `checkInit()` instead of `mustInit()` panic
- `PopErrorScope` deprecated in favor of `PopErrorScopeAsync`
- Package doc consolidated into single `doc.go` (no more duplicate package comments)
- `CONTRIBUTING.md` expanded with architecture, error handling, fuzz testing, stability sections

### Fixed

- Struct size assertions for C ABI compatibility (DrawIndirectArgs, StringView, etc.)

---

## [0.2.1] - 2026-01-29

### Changed

- **goffi:** v0.3.7 → v0.3.8
- **golang.org/x/sys:** v0.39.0 → v0.40.0

---

## [0.2.0] - 2026-01-29

### Changed

- **BREAKING:** All WebGPU types now use `github.com/gogpu/gputypes` directly
  - `TextureFormat`, `BufferUsage`, `ShaderStage`, etc. are now from gputypes
  - Enum values now match webgpu.h specification (fixes compatibility issues)
  - Example: `wgpu.TextureFormatBGRA8Unorm` → `gputypes.TextureFormatBGRA8Unorm`
- **goffi:** v0.3.7 (ARM64 Darwin improvements)

### Added

- Integration with [gogpu ecosystem](https://github.com/gogpu) via gputypes
- Full webgpu.h spec compliance for enum values
- Comprehensive conversion layer (`wgpu/convert.go`) for wgpu-native v27 compatibility
  - TextureFormat (~45 formats), VertexFormat (~30 formats)
  - VertexStepMode, TextureSampleType, TextureViewDimension, StorageTextureAccess
  - Wire structs with correct FFI padding (uint64 flags)

### Fixed

- TextureFormat enum values mismatch (BGRA8Unorm was 0x17, now correct 0x1B)
- Compatibility with gogpu Rust backend
- Struct padding in BindGroupLayout wire structs (sampler, texture, storage)
- PipelineLayout creation in examples (use CreatePipelineLayoutSimple)
- GetModuleHandleW: kernel32.dll instead of user32.dll (all Windows examples)
- Sampler MaxAnisotropy default (wgpu-native requires >= 1)
- Texture SampleCount/MipLevelCount defaults (wgpu-native requires >= 1)
- render_bundle shader: fallback without primitive_index (works on all GPUs)

### Migration Guide

Update imports in your code:
```go
import (
    "github.com/go-webgpu/webgpu/wgpu"
    "github.com/gogpu/gputypes"  // NEW
)

// Before:
config.Format = wgpu.TextureFormatBGRA8Unorm

// After:
config.Format = gputypes.TextureFormatBGRA8Unorm
```

---

## [0.1.1] - 2024-12-24

### Changed

- **goffi:** v0.3.1 → v0.3.3 (PointerType argument passing hotfix)
- **golang.org/x/sys:** v0.38.0 → v0.39.0

### Fixed

- Critical bug in PointerType argument passing ([goffi#4](https://github.com/go-webgpu/goffi/issues/4))

### Infrastructure

- Branch protection enabled for `main`
- All changes now require Pull Requests
- Updated CONTRIBUTING.md with PR workflow

---

## [0.1.0] - 2024-11-28

### Added

#### Core API
- WebGPU Instance, Adapter, Device creation
- Buffer creation and management (MapAsync, GetMappedRange, Unmap)
- Texture creation and management (2D, depth, render targets)
- Sampler API with filtering and address modes
- Shader module compilation (WGSL)

#### Pipelines
- Compute Pipeline with workgroups
- Render Pipeline with vertex/fragment stages
- Pipeline Layout and Bind Group Layout
- Bind Groups for resource binding

#### Rendering
- Command Encoder and Queue submission
- Render Pass with color and depth attachments
- Vertex buffers with custom layouts
- Index buffers (Uint16, Uint32)
- Depth buffer support (DepthStencilState)
- MRT (Multiple Render Targets) support
- RenderBundle API for pre-recording render commands

#### Advanced Features
- Instanced rendering (VertexStepModeInstance)
- Indirect drawing (DrawIndirect, DrawIndexedIndirect)
- Indirect compute dispatch (DispatchWorkgroupsIndirect)
- QuerySet API for GPU timestamps
- Error Handling API (PushErrorScope, PopErrorScope)

#### Cross-Platform
- Windows support via syscall.LazyDLL
- Linux support via goffi (CGO_ENABLED=0)
- macOS support via goffi (CGO_ENABLED=0)
- Platform-specific Surface creation (HWND, X11, Wayland, Metal)

#### Math Helpers
- Mat4 (4x4 matrix) with Identity, Translate, Scale, Rotate, Perspective, LookAt
- Vec3, Vec4 vector types with common operations

#### Examples
- Triangle (basic rendering)
- Colored triangle (vertex colors)
- Rotating triangle (uniform buffers)
- Textured quad (texture sampling)
- 3D cube (depth buffer, transforms)
- MRT (multiple render targets)
- Compute shader (parallel processing)
- Instanced rendering (25 objects in 1 draw call)
- RenderBundle (pre-recorded commands)
- Timestamp query (GPU timing)
- Error handling (error scopes)

#### Infrastructure
- Zero-CGO architecture using goffi for FFI
- GitHub Actions CI/CD (Linux, macOS, Windows)
- golangci-lint configuration for FFI code
- Pre-release validation script
- Comprehensive test suite (76+ tests)

### Dependencies
- github.com/go-webgpu/goffi v0.3.1
- wgpu-native v24.0.3.1
