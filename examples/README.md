# WebGPU Go Examples

This directory contains examples demonstrating the go-webgpu API.

## Prerequisites

1. **wgpu_native.dll** must be in your PATH or the same directory as the executable
2. **Go 1.25+**
3. **GPU with WebGPU support** (most modern GPUs)

## Building Examples

```bash
# Build all examples
cd examples
go build ./...

# Build specific example
cd adapter_info
go build
```

## Running Examples

Make sure `wgpu_native.dll` is accessible, then:

```bash
cd adapter_info
./adapter_info.exe    # Windows
./adapter_info        # Linux/macOS
```

---

## Example Categories

### 🔍 Introspection & Info

#### adapter_info
**Purpose:** Query GPU capabilities and adapter information

**What it demonstrates:**
- Getting adapter limits (max texture size, buffer size, etc.)
- Querying GPU information (vendor, device name, backend type)
- Enumerating supported features
- Checking for specific features

**Output:**
```
=== Adapter Information ===
Vendor:       NVIDIA Corporation
Device:       NVIDIA GeForce RTX 3080
Backend Type: Vulkan
Adapter Type: Discrete GPU

=== Key Adapter Limits ===
Max Texture 2D:              16384 x 16384
Max Buffer Size:             4294967296 bytes (4.00 GB)
Max Compute Workgroup Size:  1024 x 1024 x 64
```

**Use cases:**
- Check GPU capabilities before creating resources
- Display system information to users
- Debug platform-specific issues
- Validate hardware requirements

---

#### buffer_introspection
**Purpose:** Demonstrate buffer state introspection at runtime

**What it demonstrates:**
- Querying buffer size
- Checking buffer usage flags
- Getting buffer map state
- Buffer lifecycle management

**Output:**
```
=== Buffer Introspection ===
Buffer size: 1048576 bytes (1.00 MB)
Buffer usage: Storage | CopySrc | CopyDst
Buffer map state: Unmapped

=== Mappable Buffer Example ===
Initial map state (MappedAtCreation): Mapped
Map state after Unmap(): Unmapped
Map state after MapAsync(): Mapped
```

**Use cases:**
- Debug buffer mapping issues
- Validate buffer state before operations
- GPU memory profiling
- Resource lifecycle tracking

---

### 🐛 Debugging & Profiling

#### render_debug_markers
**Purpose:** GPU debugging with hierarchical markers

**What it demonstrates:**
- Inserting single debug markers
- Pushing/popping nested debug groups
- Creating hierarchical GPU timelines
- Integration with GPU profiling tools

**Output:**
```
=== Demonstrating RenderPass Debug Markers ===
✓ Inserted debug marker: 'Frame Start'
✓ Pushed debug group: 'Scene Rendering'
  ✓ Pushed nested group: 'Geometry Pass'
    ✓ Inserted marker: 'Draw Opaque Objects'
  ✓ Popped debug group: 'Geometry Pass'

=== Debug Marker Hierarchy ===
Frame Start [marker]
└─ Scene Rendering [group]
   ├─ Geometry Pass [group]
   │  ├─ Draw Opaque Objects [marker]
   │  └─ Draw Alpha-Tested Objects [marker]
   └─ Lighting Pass [group]
```

**GPU Tools Support:**
- RenderDoc - Event browser
- PIX (Windows) - Timeline view
- Xcode GPU Debugger (macOS) - Capture hierarchy
- Chrome DevTools - WebGPU profiler
- NVIDIA Nsight - Performance timeline

**Use cases:**
- Identify GPU bottlenecks
- Debug rendering issues
- Optimize draw call order
- Profile complex render passes
- Create professional GPU captures

---

## Advanced Examples

### compute_shader
Demonstrates compute shaders, storage buffers, and GPU computation.

### triangle
Basic rendering pipeline with vertex buffers and render passes.

### texture
Texture creation, sampling, and rendering to texture.

### surface
Window surface creation and presentation (requires windowing library).

---

## Example Structure

All examples follow this pattern:

```go
package main

import (
    "github.com/go-webgpu/webgpu/wgpu"
)

func main() {
    // 1. Create instance
    instance, _ := wgpu.CreateInstance(nil)
    defer instance.Release()

    // 2. Request adapter
    adapter, _ := instance.RequestAdapter(nil)
    defer adapter.Release()

    // 3. Request device
    device, _ := adapter.RequestDevice(nil)
    defer device.Release()

    // 4. Create resources and render
    // ...
}
```

### Important Notes

1. **Always defer Release()** - Prevents resource leaks
2. **Check errors** - Production code should handle all errors
3. **Init order** - Instance → Adapter → Device → Resources
4. **GPU polling** - Some operations require `device.Poll()` or `instance.ProcessEvents()`

---

## Troubleshooting

### "Failed to load library"
- Ensure `wgpu_native.dll` is in PATH or current directory
- Check DLL matches your architecture (x64/x86)

### "Failed to create instance"
- Update GPU drivers
- Check system has WebGPU-compatible GPU

### "Failed to request adapter"
- Try `PowerPreference: wgpu.PowerPreferenceLowPower`
- Check GPU supports required features

### Example crashes or hangs
- Update to latest wgpu-native build
- Check GPU drivers are up to date
- Try different backend (set via environment variables)

---

## Environment Variables

```bash
# Force specific backend (Windows)
set WGPU_BACKEND=dx12    # or vulkan, dx11

# Force specific backend (Linux/macOS)
export WGPU_BACKEND=vulkan  # or metal, gl

# Enable validation (helpful for debugging)
export WGPU_VALIDATION=1

# Enable debug output
export RUST_LOG=wgpu=debug
```

---

## Performance Tips

1. **Reuse resources** - Don't create/destroy buffers every frame
2. **Batch draw calls** - Group by pipeline/bind group
3. **Use staging buffers** - For frequent CPU→GPU uploads
4. **Profile first** - Use debug markers to identify bottlenecks
5. **Cache limits** - Call `GetLimits()` once, not per frame

---

## Contributing Examples

To add a new example:

1. Create directory: `examples/my_example/`
2. Add `main.go` with clear documentation
3. Update this README with description
4. Test on Windows, Linux, and macOS if possible
5. Follow existing code style (see LINTER_RULES.md)

---

## References

- [WebGPU Specification](https://www.w3.org/TR/webgpu/)
- [wgpu-native Documentation](https://wgpu-native.github.io/wgpu-native/)
- [Project Documentation](../docs/)
- [API Implementation Details](../docs/NEW_API_IMPLEMENTATION.md)

---

**Last Updated:** 2026-01-29
**go-webgpu Version:** 0.2.0

**Note:** All examples use [gputypes](https://github.com/gogpu/gputypes) for WebGPU type definitions.
