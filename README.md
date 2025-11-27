# go-webgpu

Zero-CGO WebGPU bindings for Go using [goffi](https://github.com/go-webgpu/goffi) + [wgpu-native](https://github.com/gfx-rs/wgpu-native).

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

## Dependencies

- [goffi](https://github.com/go-webgpu/goffi) — Pure Go FFI for callbacks
- [wgpu-native](https://github.com/gfx-rs/wgpu-native) — Rust WebGPU implementation
- [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) — Platform-specific syscalls

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.
