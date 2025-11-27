# Textured Quad Example

This example demonstrates texture creation, sampling, and rendering a textured quad using go-webgpu.

## Features

1. **Procedural Texture Generation**
   - Creates a 256x256 RGBA8 checkerboard pattern
   - Alternating light yellow and dark blue squares (8x8 grid)

2. **Texture Upload**
   - Uses `Queue.WriteTexture()` for GPU data transfer
   - Proper BytesPerRow alignment (256-byte aligned)

3. **Sampler Creation**
   - Linear filtering using `Device.CreateLinearSampler()`
   - Clamp-to-edge addressing mode

4. **Bind Group Setup**
   - Group 0, Binding 0: Sampler
   - Group 0, Binding 1: Texture 2D
   - Uses helper functions `SamplerBindingEntry()` and `TextureBindingEntry()`

5. **Vertex Buffer with UV Coordinates**
   - 4 vertices forming a quad
   - Each vertex: position (vec2f) + UV coordinates (vec2f)
   - ArrayStride: 16 bytes (4 floats × 4 bytes)

6. **Index Buffer**
   - 6 indices forming 2 triangles
   - IndexFormat: Uint16

## Building

```bash
cd examples/textured-quad
go build -o textured-quad.exe .
```

## Running

```bash
./textured-quad.exe
```

A window will open showing a textured quad with a checkerboard pattern.

## Key Technical Details

- **Texture Format**: RGBA8Unorm
- **Texture Size**: 256×256
- **Vertex Format**: Float32x2 (position) + Float32x2 (UV)
- **BytesPerRow Alignment**: Must be multiple of 256 bytes for `WriteTexture()`
- **DepthSliceUndefined**: Used for 2D textures (0xFFFFFFFF)

## Code Structure

```go
// 1. Create texture
texture := device.CreateTexture(&TextureDescriptor{...})
queue.WriteTexture(&destInfo, textureData, &layout, &extent)

// 2. Create sampler
sampler := device.CreateLinearSampler()

// 3. Create bind group layout
layoutEntries := []BindGroupLayoutEntry{
    {Binding: 0, Visibility: ShaderStageFragment, Sampler: ...},
    {Binding: 1, Visibility: ShaderStageFragment, Texture: ...},
}
bindGroupLayout := device.CreateBindGroupLayoutSimple(layoutEntries)

// 4. Create bind group
entries := []BindGroupEntry{
    SamplerBindingEntry(0, sampler),
    TextureBindingEntry(1, textureView),
}
bindGroup := device.CreateBindGroupSimple(bindGroupLayout, entries)

// 5. Use in render pass
pass.SetBindGroup(0, bindGroup, nil)
```

## WGSL Shader

```wgsl
@group(0) @binding(0) var texSampler: sampler;
@group(0) @binding(1) var tex: texture_2d<f32>;

@fragment
fn fs_main(@location(0) uv: vec2f) -> @location(0) vec4f {
    return textureSample(tex, texSampler, uv);
}
```

## Platform Support

- ✅ Windows (tested)
- 🔄 Linux (should work, needs testing)
- 🔄 macOS (should work, needs testing)

## See Also

- [colored-triangle](../colored-triangle/) - Basic rendering with vertex colors
- [compute-shader](../compute-shader/) - Compute shader example
