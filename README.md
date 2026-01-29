# go-webgpu

> **Zero-CGO WebGPU bindings for Go — GPU-accelerated graphics and compute in pure Go**

[![GitHub Release](https://img.shields.io/github/v/release/go-webgpu/webgpu?include_prereleases&style=flat-square&logo=github&color=blue)](https://github.com/go-webgpu/webgpu/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/go-webgpu/webgpu?style=flat-square&logo=go)](https://go.dev/dl/)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-webgpu/webgpu.svg)](https://pkg.go.dev/github.com/go-webgpu/webgpu)
[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/go-webgpu/webgpu/test.yml?branch=main&style=flat-square&logo=github-actions&label=CI)](https://github.com/go-webgpu/webgpu/actions)
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
| Textures, Samplers, Storage Textures | ✅ |
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
- wgpu-native v24.0.3.1 ([download](https://github.com/gfx-rs/wgpu-native/releases))

## Installation

```bash
go get github.com/go-webgpu/webgpu
```

Download wgpu-native and place `wgpu_native.dll` (Windows) or `libwgpu_native.so` (Linux) in your project directory or system PATH.

## Type System

This library uses [gputypes](https://github.com/gogpu/gputypes) for WebGPU type definitions, ensuring compatibility with the [gogpu ecosystem](https://github.com/gogpu) and webgpu.h specification.

```go
import (
    "github.com/go-webgpu/webgpu/wgpu"
    "github.com/gogpu/gputypes"
)

// Use gputypes for WebGPU enums
config.Format = gputypes.TextureFormatBGRA8Unorm
buffer.Usage = gputypes.BufferUsageVertex | gputypes.BufferUsageCopyDst
```

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
├─────────────────────────────────────────┤
│     go-webgpu (this package)            │
│     - Zero CGO                          │
│     - Pure Go FFI via goffi             │
├─────────────────────────────────────────┤
│     wgpu-native (Rust WebGPU)           │
├─────────────────────────────────────────┤
│     Vulkan / Metal / DX12 / OpenGL      │
└─────────────────────────────────────────┘
```

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

- [goffi](https://github.com/go-webgpu/goffi) — Pure Go FFI for callbacks
- [wgpu-native](https://github.com/gfx-rs/wgpu-native) — Rust WebGPU implementation
- [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) — Platform-specific syscalls

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.
