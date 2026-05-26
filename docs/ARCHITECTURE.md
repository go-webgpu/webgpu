# Architecture

This document explains the internal design of go-webgpu for contributors. It covers the layered architecture, the FFI contract, the enum translation layer, and the testing strategy.

---

## Table of Contents

- [High-Level Overview](#high-level-overview)
- [FFI Layer Architecture](#ffi-layer-architecture)
- [Struct Layout Contract](#struct-layout-contract)
- [convert.go — The Enum Translation Layer](#convertgo--the-enum-translation-layer)
- [gputypes Relationship](#gputypes-relationship)
- [Async Callbacks](#async-callbacks)
- [Testing Strategy](#testing-strategy)
- [Ecosystem Context](#ecosystem-context)

---

## High-Level Overview

```
Your Go Application
        │ uses public API (Go strings, typed values)
        ▼
go-webgpu/wgpu (public package)
   ├─ Public structs   — Go-idiomatic types, gputypes aliases
   ├─ Method receivers — error-returning wrappers
   └─ gputypes_aliases.go — re-exports for single-import UX
        │ calls
        ▼
Internal FFI layer
   ├─ Wire structs     — must match C ABI exactly (field order, padding)
   ├─ convert.go       — translates enum values between gputypes and wgpu-native
   ├─ loader*.go       — cross-platform dynamic library loading
   └─ wgpu.go / procs  — syscall procedure handles
        │ calls
        ▼
wgpu-native v29 (Rust) — dynamically loaded .dll / .so / .dylib
        │ calls
        ▼
Vulkan / Metal / D3D12 / OpenGL
```

**Design principle**: the public API is Go-idiomatic. The FFI-level ABI complexity is entirely contained in the internal layer and never surfaces to users.

---

## FFI Layer Architecture

### Public structs vs Wire structs

Public structs use Go-friendly types. Wire structs must match C memory layout exactly.

```go
// Public struct — user creates this
type TextureDescriptor struct {
    Label           StringView
    Usage           gputypes.TextureUsage  // Go-typed
    Dimension       gputypes.TextureDimension
    Format          gputypes.TextureFormat
    // ...
}

// Wire struct — passed directly to wgpu-native
type textureDescriptorWire struct {
    NextInChain     uintptr
    Label           StringView
    Usage           uint64   // CRITICAL: wgpu-native uses uint64 for flags
    Dimension       uint32   // converted: gputypes value → wgpu-native value
    Format          uint32   // converted via map in convert.go
    // ...
}
```

Conversion happens **inside each method** before the FFI call:

```go
func (d *Device) CreateTexture(desc *TextureDescriptor) (*Texture, error) {
    wire := textureDescriptorWire{
        Usage:     uint64(desc.Usage),                   // bitflag: direct cast
        Dimension: toWGPUTextureDimension(desc.Dimension), // +1 shift
        Format:    toWGPUTextureFormat(desc.Format),      // lookup table
        // ...
    }
    // Call wgpu-native with &wire, not &desc
}
```

The user never sees wire structs. They are unexported and created only at the FFI call site.

### Procedure handles

All wgpu-native functions are loaded once at `Init()` time via platform-specific loaders:

- `loader_windows.go` — `syscall.LazyDLL` / `syscall.LazyProc`
- `loader_unix.go` — goffi `ffi.LoadLibrary` / `ffi.GetSymbol`

Procedures are package-level variables (`procCreateInstance`, `procDeviceGetQueue`, etc.) used directly in method bodies.

---

## Struct Layout Contract

**Every wire struct field order, type, and padding must exactly match the C struct in webgpu.h.**

This is verified at compile time by `abi_test.go` (271 tests) using `unsafe.Sizeof` and `unsafe.Offsetof` assertions.

Rules:
1. Field order must follow the C struct definition (copy from webgpu.h comments)
2. Enum fields that are `uint32` in C must be `uint32` in Go (not Go enum types)
3. Flag fields that are `uint64` in C must be `uint64` in Go (wgpu-native extends several flag types to 64-bit)
4. Pointer fields are `uintptr` (8 bytes on 64-bit, covers `WGPUSomething*` and `void*`)
5. Booleans are `uint32` (`WGPUBool` = C `uint32_t`)
6. `StringView` is `{Data uintptr, Length uintptr}` — two pointers (16 bytes on 64-bit)

**Any struct change requires updating both the Go struct and `abi_test.go`.**

Example of a v29 ABI-breaking change: `Limits` gained `NextInChain uintptr` as its first field and two fields changed position. This would silently corrupt all limit queries if the wire struct was not updated.

---

## convert.go — The Enum Translation Layer

`convert.go` bridges two independent numbering systems: **gputypes** (Go-idiomatic, WebGPU JS spec) and **wgpu-native v29** (C spec, with some structural differences).

### Why conversions are needed

The wgpu-native v29 C header introduces `BindingNotUsed=0` as a sentinel for binding-related enums. gputypes does not have this sentinel and starts from `Undefined=0`, shifting all subsequent values by one position.

### Enums that need explicit conversion

| Enum | Reason | Conversion |
|------|--------|-----------|
| `BufferBindingType` | v29 adds `BindingNotUsed=0` | `Undefined=0 → 0`, others `+1` |
| `SamplerBindingType` | same | same |
| `TextureSampleType` | same | same |
| `StorageTextureAccess` | same | same |
| `VertexFormat` | gputypes omits single-component 8/16-bit variants added in v29 | full lookup table |
| `VertexStepMode` | gputypes has `VertexBufferNotUsed=1` removed in v29; values shifted | lookup table |

Conversion functions follow the naming convention `toWGPU<EnumName>` and return `uint32`.

### Enums that match exactly — no converter needed

All bitflag enums match between gputypes and wgpu-native v29 because they use power-of-2 values defined by the WebGPU spec:

- `TextureFormat` — gputypes v0.3.0 matches v29 exactly, including `R16*/RG16*` Unorm/Snorm variants
- All flag types: `BufferUsage`, `TextureUsage`, `ShaderStage`, `ColorWriteMask`, `MapMode`
- `LoadOp`, `StoreOp`, `BlendFactor`, `BlendOperation`
- `PrimitiveTopology`, `FrontFace`, `CullMode`, `IndexFormat`
- `FilterMode`, `MipmapFilterMode`, `AddressMode`, `CompareFunction`, `StencilOperation`
- `TextureViewDimension`, `TextureDimension`, `TextureAspect`
- `PresentMode`, `CompositeAlphaMode`, `PowerPreference`

For these enums, a direct `uint32(value)` cast is safe.

### Adding a new conversion

When adding support for a new enum, check whether it has a `BindingNotUsed=0` or structural gap. Compare the gputypes values against the webgpu.h C header values before deciding whether a converter is needed.

---

## gputypes Relationship

`gogpu/gputypes` is the shared type registry for the entire gogpu ecosystem. It defines WebGPU enums and struct types following the WebGPU JavaScript specification numbering.

go-webgpu re-exports gputypes types and constants as type aliases in `gputypes_aliases.go`:

```go
// gputypes_aliases.go
type TextureFormat = gputypes.TextureFormat
const TextureFormatBGRA8Unorm = gputypes.TextureFormatBGRA8Unorm
// ...
```

This means users can write `wgpu.TextureFormatBGRA8Unorm` with a single import. There is no wrapping or copying — `wgpu.TextureFormat` and `gputypes.TextureFormat` are the same type at the Go type system level.

**Important**: gputypes values are NOT the same numbers as wgpu-native C values for the enums listed in the conversion table above. `convert.go` is the bridge between these two numbering systems. Never pass a gputypes enum value directly to an FFI function without checking whether a converter exists for that type.

---

## Async Callbacks

wgpu-native uses callbacks for `RequestAdapter`, `RequestDevice`, `MapAsync`, and error scopes. The Go implementation converts these to synchronous calls using channels:

1. A `*Request` state struct with a `done chan struct{}` is allocated and registered in a global map
2. The request ID is passed as `Userdata1` to the C callback
3. The callback (registered via `ffi.NewCallback`) looks up the request by ID, writes results, and closes `done`
4. The calling goroutine blocks on `<-req.done` in a select loop that also calls `ProcessEvents()`

Callback function pointers are created once via `sync.Once` and reused across calls. The global registry is protected by a `sync.Mutex`.

**Thread safety**: each callback modifies only its own request struct, and deletes the map entry atomically under the lock before writing to the struct, so there is no race between the callback and the waiting goroutine.

---

## Testing Strategy

Tests are organized into three tiers by GPU requirement:

### ABI tests (no GPU, always run in CI)

`abi_test.go` — 271 assertions verifying:
- `unsafe.Sizeof(wireStruct)` matches `C sizeof(WGPUStruct)`
- `unsafe.Offsetof(wireStruct.Field)` matches C field offsets
- Enum constant values match webgpu.h integers
- gputypes type alignment with wgpu-native values for pass-through enums

These tests have zero external dependencies and catch ABI regressions immediately.

### Safe tests (no GPU, CI-safe)

Tests that exercise logic without making FFI calls to wgpu-native:
- `TestMat4*`, `TestVec3*` — math helpers
- `TestStructSizes*` — additional struct size checks
- `TestWGPUError*` — error type assertions
- `TestNullGuard*` — nil/released handle defensive checks
- `Fuzz*` — FFI boundary fuzz targets (seed corpus only in CI)

Filter: `-run "Mat4|Vec3|StructSizes|CheckInit|WGPUError|Fuzz|NullGuard"` in `.github/workflows/`.

### GPU tests (local only, skipped in CI)

Tests that require a real GPU and wgpu-native loaded:
- `TestAdapter*`, `TestDevice*`, `TestBuffer*`, `TestSurface*`
- `TestLeak*`, `TestErrorScope*`

These require `WGPU_NATIVE_PATH` pointing to a valid wgpu-native binary. GitHub Actions runners have no GPU, so these tests are excluded from CI filter patterns.

---

## Ecosystem Context

```
born-ml/born (ML framework)
        │ uses gogpu for GPU compute
        ▼
gogpu/gogpu (graphics framework)
        │ backend selection
       ┌┴──────────────────────────┐
       ▼                           ▼
go-webgpu/webgpu            gogpu/wgpu
(FFI → wgpu-native)         (Pure Go WebGPU)
       │
       ▼
go-webgpu/goffi (Pure Go FFI)

Shared: gogpu/gputypes ← WebGPU type definitions used by all
```

| Project | Approach | Runtime requirement |
|---------|----------|---------------------|
| go-webgpu/webgpu | FFI to wgpu-native (Rust) | `wgpu_native.dll` / `.so` / `.dylib` |
| gogpu/wgpu | Pure Go WebGPU implementation | None |

The two implementations share `gputypes` as their type contract. This is what makes gogpu backend switching possible: the same `gputypes.TextureFormat` constant works with both backends, and only the FFI translation (our `convert.go`) changes.
