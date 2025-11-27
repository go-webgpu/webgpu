# Triangle Example

Simple triangle rendering using go-webgpu bindings.

## Description

This example demonstrates the basic WebGPU rendering pipeline:
- Creating a window using Win32 API (Windows only)
- Initializing WebGPU (Instance, Adapter, Device, Queue)
- Creating a Surface from HWND
- Compiling WGSL shaders
- Creating a RenderPipeline
- Rendering a red triangle with clear color

## Building

```bash
# From webgpu root directory
go build -o examples/triangle/triangle.exe ./examples/triangle/

# Or directly
cd examples/triangle
go build
```

## Running

```bash
# From examples/triangle directory
./triangle.exe

# Or from webgpu root
./examples/triangle/triangle.exe
```

## What You Should See

A window with a cornflower blue background (RGB: 0.1, 0.2, 0.3) and a red triangle in the center.

The triangle:
- **Top vertex**: (0.0, 0.5) - center top
- **Bottom-left vertex**: (-0.5, -0.5)
- **Bottom-right vertex**: (0.5, -0.5)
- **Color**: Red (1.0, 0.0, 0.0, 1.0)

## Platform Support

Currently Windows-only (Win32 API). Linux/macOS support coming in future milestones.

## Code Structure

- **Window creation**: Win32 FFI using `golang.org/x/sys/windows`
- **Message loop**: `PeekMessage` for non-blocking event processing
- **WebGPU initialization**: Standard flow (Instance → Adapter → Device)
- **Surface configuration**: BGRA8Unorm format, VSync enabled (Fifo)
- **Render pipeline**: Simple vertex + fragment shaders (WGSL)
- **Render loop**: Acquire texture → Encode commands → Submit → Present

## Troubleshooting

**Problem**: Window doesn't appear or crashes immediately
- **Solution**: Make sure `wgpu_native.dll` is in the same directory as the executable or in PATH

**Problem**: Black screen
- **Solution**: Check if your GPU supports WebGPU. Run with validation layers to see errors.

**Problem**: Build fails with "missing wgpu package"
- **Solution**: Run `go mod tidy` from webgpu root directory

## Next Steps

After running this example, try:
1. Modifying triangle vertices in the shader
2. Changing clear color
3. Adding multiple triangles
4. Experimenting with different colors per vertex

See other examples for more advanced features (textures, buffers, compute shaders).
