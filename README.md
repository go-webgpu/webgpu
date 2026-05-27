# go-webgpu

> **Zero-CGO WebGPU bindings for Go — GPU-accelerated graphics and compute in pure Go**

[![GitHub Release](https://img.shields.io/github/v/release/go-webgpu/webgpu?include_prereleases&style=flat-square&logo=github&color=blue)](https://github.com/go-webgpu/webgpu/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/go-webgpu/webgpu?style=flat-square&logo=go)](https://go.dev/dl/)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-webgpu/webgpu.svg)](https://pkg.go.dev/github.com/go-webgpu/webgpu)
[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/go-webgpu/webgpu/test.yml?branch=main&style=flat-square&logo=github-actions&label=CI)](https://github.com/go-webgpu/webgpu/actions)
[![Codecov](https://img.shields.io/codecov/c/github/go-webgpu/webgpu?style=flat-square&logo=codecov)](https://app.codecov.io/gh/go-webgpu/webgpu)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-webgpu/webgpu?style=flat-square)](https://goreportcard.com/report/github.com/go-webgpu/webgpu)
[![License](https://img.shields.io/github/license/go-webgpu/webgpu?style=flat-square)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/go-webgpu/webgpu?style=flat-square&logo=github)](https://github.com/go-webgpu/webgpu/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/go-webgpu/webgpu?style=flat-square&logo=github)](https://github.com/go-webgpu/webgpu/issues)

Pure Go WebGPU bindings using [goffi](https://github.com/go-webgpu/goffi) + [wgpu-native](https://github.com/gfx-rs/wgpu-native). No CGO required.

## Status

**Beta** — Comprehensive API ready for testing and feedback.

| Feature | Status |
|---------|--------|
| Instance, Adapter, Device | ✅ |
| Buffers (vertex, index, uniform, storage) | ✅ |
| Buffer Mapping (Map with context, async MapPending, type-safe MappedRange) | ✅ |
| Queue Submission Index Tracking | ✅ |
| Textures, Samplers, Storage Textures | ✅ |
| Region-Based Copy Operations (CopyTextureToBuffer) | ✅ |
| Render Pipelines | ✅ |
| Compute Pipelines | ✅ |
| Depth Buffer | ✅ |
| MRT (Multiple Render Targets) | ✅ |
| Instanced Rendering | ✅ |
| Indirect Drawing (GPU-driven) | ✅ |
| RenderBundle (pre-recorded commands) | ✅ |
| Cross-Platform Surface (Win/Linux/macOS) | ✅ |
| Error Handling (error scopes) | ✅ |
| QuerySet (GPU timestamps) | ✅ |
| BindGroupLayout (explicit) | ✅ |
| 3D Math (Mat4, Vec3) | ✅ |

## Requirements

- Go 1.25+
- wgpu-native v29.0.0.0 ([download](https://github.com/gfx-rs/wgpu-native/releases))

## Installation

```bash
go get github.com/go-webgpu/webgpu
```

### wgpu-native Setup (auto)

```bash
go run github.com/go-webgpu/webgpu/cmd/setup@latest
```

This downloads the correct wgpu-native v29 binary for your platform (Windows/macOS/Linux, amd64/arm64) into `./lib/`.

The library is found automatically from `./lib/` — no environment variable needed. To override the search path explicitly:

```bash
# Linux
export WGPU_NATIVE_PATH=./lib/libwgpu_native.so

# macOS
export WGPU_NATIVE_PATH=./lib/libwgpu_native.dylib

# Windows (PowerShell)
$env:WGPU_NATIVE_PATH = "lib\wgpu_native.dll"

# Windows (cmd)
set WGPU_NATIVE_PATH=lib\wgpu_native.dll
```

Or copy the library to your project root — it will also be found automatically.

### wgpu-native Setup (manual)

Download from [gfx-rs/wgpu-native releases](https://github.com/gfx-rs/wgpu-native/releases/tag/v29.0.0.0) and place in your project directory or system PATH.

Custom library location:
```bash
export WGPU_NATIVE_PATH=/path/to/libwgpu_native.so      # Linux/macOS
set WGPU_NATIVE_PATH=C:\path\to\wgpu_native.dll          # Windows
```

## gogpu/wgpu Integration

This library is the **Rust FFI backend** for [gogpu/wgpu](https://github.com/gogpu/wgpu) — the unified Go WebGPU package. Build with `-tags rust` to use wgpu-native instead of the Pure Go implementation:

```bash
go build -tags rust ./myapp
```

Same API, same types, same user code — build tag selects the implementation.

## Type System

WebGPU types from [gputypes](https://github.com/gogpu/gputypes) are re-exported as type aliases in the `wgpu` package. A single import is sufficient:

```go
import "github.com/go-webgpu/webgpu/wgpu"

// gputypes constants available directly via wgpu package
config.Format = wgpu.TextureFormatBGRA8Unorm
buffer.Usage = wgpu.BufferUsageVertex | wgpu.BufferUsageCopyDst
```

If you use multiple gogpu ecosystem packages, importing `gputypes` directly also works and is fully compatible:

```go
import (
    "github.com/go-webgpu/webgpu/wgpu"
    "github.com/gogpu/gputypes" // optional — for ecosystem interop
)

config.Format = gputypes.TextureFormatBGRA8Unorm // same underlying type
```

go-webgpu API is designed to be compatible with gogpu/wgpu, enabling future backend switching within the gogpu ecosystem.

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/go-webgpu/webgpu/wgpu"
)

func main() {
    // Initialize library
    if err := wgpu.Init(); err != nil {
        log.Fatal(err)
    }

    // Create WebGPU instance
    instance, err := wgpu.CreateInstance(nil)
    if err != nil {
        log.Fatal(err)
    }
    defer instance.Release()

    // Request GPU adapter
    adapter, err := instance.RequestAdapter(nil)
    if err != nil {
        log.Fatal(err)
    }
    defer adapter.Release()

    fmt.Printf("Adapter: %#x\n", adapter.Handle())
}
```

## Examples

| Example | Description |
|---------|-------------|
| [triangle](examples/triangle) | Basic triangle rendering |
| [colored-triangle](examples/colored-triangle) | Vertex colors and buffers |
| [textured-quad](examples/textured-quad) | Textures, samplers, index buffers |
| [rotating-triangle](examples/rotating-triangle) | Uniform buffers, animation |
| [cube](examples/cube) | 3D rendering with depth buffer |
| [instanced](examples/instanced) | Instanced rendering (25 objects in 1 draw call) |
| [compute](examples/compute) | Compute shader parallel processing |
| [indirect](examples/indirect) | GPU-driven rendering (DrawIndirect) |
| [render_bundle](examples/render_bundle) | Pre-recorded draw commands |
| [timestamp_query](examples/timestamp_query) | GPU profiling with timestamps |
| [mrt](examples/mrt) | Multiple Render Targets |
| [error_handling](examples/error_handling) | Error scopes API |

Run examples:
```bash
cd examples/triangle && go run .
```

## Architecture

```
┌─────────────────────────────────────────┐
│            Your Go Application          │
│    (uses wgpu.TextureFormatBGRA8Unorm)  │
├─────────────────────────────────────────┤
│     go-webgpu/wgpu (this package)       │
│     - Go-idiomatic API                  │
│     - gputypes type aliases             │
│     - Error returns on all operations   │
├─────────────────────────────────────────┤
│     Internal FFI layer                  │
│     - Wire structs (C-layout)           │
│     - convert.go (enum translation)     │
│     - goffi (Pure Go FFI)               │
├─────────────────────────────────────────┤
│     wgpu-native v29 (Rust WebGPU)       │
├─────────────────────────────────────────┤
│     Vulkan / Metal / DX12 / OpenGL      │
└─────────────────────────────────────────┘
```

For a detailed explanation of the architecture, including why convert.go exists and the FFI wire struct contract, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Looking for Pure Go WebGPU?

This project uses FFI bindings to wgpu-native. If you're looking for a **100% Pure Go** WebGPU implementation (no native dependencies), check out:

👉 **[github.com/gogpu](https://github.com/gogpu)** — Pure Go GPU ecosystem

| Project | Description |
|---------|-------------|
| [gogpu/wgpu](https://github.com/gogpu/wgpu) | Pure Go WebGPU implementation |
| [gogpu/naga](https://github.com/gogpu/naga) | Pure Go shader compiler (WGSL/SPIR-V) |
| [gogpu/gogpu](https://github.com/gogpu/gogpu) | High-level GPU compute framework |
| [gogpu/gg](https://github.com/gogpu/gg) | Pure Go graphics library |

## Dependencies

- [goffi](https://github.com/go-webgpu/goffi) — Pure Go FFI (callbacks, cross-platform loading)
- [gputypes](https://github.com/gogpu/gputypes) — Shared WebGPU type definitions for the gogpu ecosystem
- [wgpu-native](https://github.com/gfx-rs/wgpu-native) — Rust WebGPU implementation (runtime binary, not a Go dependency)
- [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) — Platform-specific syscalls

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.
