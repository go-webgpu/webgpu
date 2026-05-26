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

## Buffer Mapping Architecture

Buffer mapping in WebGPU is inherently async: the GPU must finish any in-flight work on the buffer before the CPU can access it. go-webgpu provides three access patterns:

### Map (blocking, context-aware)

`Buffer.Map(ctx, mode, offset, size) error` — the recommended path for most applications.

1. Calls `mapAsyncStart` which issues `wgpuBufferMapAsync` with a Go callback registered in `mapRequests` (global map, protected by `sync.Mutex`)
2. The callback (`mapCallbackHandler`) is a C-callable function pointer created once via `ffi.NewCallback`
3. After submitting the request, `Map` kicks an initial `Device.Poll(false)` (for synchronous-complete backends)
4. If not immediately complete, a background goroutine drives `Device.Poll` continuously so the mapping resolves without the caller needing to pump events
5. Blocks on `<-req.done` or `ctx.Done()`, whichever fires first

### MapAsync (non-blocking)

`Buffer.MapAsync(mode, offset, size) (*MapPending, error)` — for callers that want to do other work while waiting.

1. Same `mapAsyncStart` path as Map
2. Returns a `*MapPending` immediately without blocking
3. `MapPending.Status()` performs a non-blocking `select` on `req.done`
4. `MapPending.Wait(ctx)` blocks with context support
5. Caller must drive `Device.Poll` externally until `Status()` returns ready

### MappedRange (type-safe access)

After `Map` or `MapAsync` resolves, `Buffer.MappedRange(offset, size) (*MappedRange, error)` wraps `Buffer.GetMappedRange` and validates the buffer state (must be `BufferMapStateMapped`). The returned `MappedRange.Bytes()` returns a `[]byte` backed by the GPU-mapped memory, valid until `Buffer.Unmap()`.

```
caller
  │ Map(ctx, ...)
  ▼
mapAsyncStart ──► wgpuBufferMapAsync (C FFI)
  │                        │
  │                   C callback ──► mapCallbackHandler (Go)
  │                                        │ closes req.done
  ▼
select req.done / ctx.Done
  │ done
  ▼
MappedRange ──► Bytes() ──► []byte view into GPU memory
  │
  ▼
Unmap ──► wgpuBufferUnmap (C FFI)
```

---

## Limits Caching

`Adapter.Limits()` and `Device.Limits()` return a cached `Limits` value with no error. Limits are fetched once via `wgpuAdapterGetLimits` / `wgpuDeviceGetLimits` at `RequestAdapter` / `RequestDevice` time and stored inside the `Adapter` / `Device` struct.

This design has two benefits:
1. **No error handling at call site** — limits are always available once you have a valid Adapter or Device
2. **No FFI overhead on repeated access** — common in render loops that check `MaxUniformBufferBindingSize` or similar

The cached value is read-only and thread-safe (written once before the struct is returned to the caller).

---

## Wire Struct Pattern Summary

Every public descriptor type has an unexported `*Wire` counterpart used at the FFI boundary. The pattern is always:

1. **Public struct** — Go-idiomatic types (`string`, `bool`, `*T`, `[]T`)
2. **Conversion** — inside the method body, construct the wire struct from the public struct
3. **Wire struct** — C-layout types (`uintptr`, `uint32`, `uint64`, `StringView`, `Bool`)
4. **FFI call** — pass `unsafe.Pointer` to the wire struct

The wire struct is a local variable on the stack; its address is safe to pass to wgpu-native only for the duration of the FFI call (wgpu-native copies all descriptor data synchronously).

```go
// Public descriptor — what the user writes
type BufferDescriptor struct {
    Label            string
    Usage            gputypes.BufferUsage
    Size             uint64
    MappedAtCreation bool
}

// Wire struct — matches WGPUBufferDescriptor byte-for-byte
type bufferDescriptorWire struct {
    NextInChain      uintptr
    Label            StringView           // {Data uintptr, Length uintptr}
    Usage            gputypes.BufferUsage // uint64 on wgpu-native
    Size             uint64
    MappedAtCreation Bool                 // uint32 (WGPUBool)
    _pad             [4]byte
}
```

The 271 ABI tests in `abi_test.go` verify that every wire struct matches the C header at compile time.

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
