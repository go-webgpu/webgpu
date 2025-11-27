# Multiple Render Targets (MRT) Example

This example demonstrates **Multiple Render Targets (MRT)** in go-webgpu, rendering to two textures simultaneously from a single draw call. MRT is commonly used in deferred rendering pipelines for G-buffers.

## Features

- **Two Render Targets**: Renders to two textures at once
  - Target 0: Color output (BGRA8Unorm) - displayed on screen
  - Target 1: Position data (RGBA8Unorm) - offscreen texture
- **Single Draw Call**: Fragment shader outputs to both targets simultaneously
- **Rotating Triangle**: Reuses uniform buffer animation from rotating-triangle example
- **WGSL MRT Syntax**: Uses `FragmentOutput` struct with `@location(0)` and `@location(1)`

## Use Cases

MRT is commonly used for:

1. **Deferred Rendering**: Render geometry data to G-buffer (albedo, normals, depth, etc.)
2. **Post-Processing**: Generate multiple intermediate textures for effects
3. **Debug Visualization**: Output normal/position data alongside color
4. **Compute Pipelines**: Pre-process data in multiple formats

## Key Concepts

### Fragment Shader with MRT

```wgsl
struct FragmentOutput {
    @location(0) color: vec4f,
    @location(1) extra: vec4f,
}

@fragment
fn fs_main(in: VertexOutput) -> FragmentOutput {
    var out: FragmentOutput;
    // Target 0: regular color
    out.color = vec4f(in.color, 1.0);
    // Target 1: position encoded as color
    out.extra = vec4f(in.original_pos * 0.5 + 0.5, 0.0, 1.0);
    return out;
}
```

- **`@location(0)`**: First render target (screen/swapchain)
- **`@location(1)`**: Second render target (offscreen texture)
- **Multiple outputs**: Fragment shader returns struct with multiple fields

### Creating the Extra Texture

```go
extraTexture := device.CreateTexture(&wgpu.TextureDescriptor{
    Usage: wgpu.TextureUsageRenderAttachment | wgpu.TextureUsageTextureBinding,
    Dimension: wgpu.TextureDimension2D,
    Size: wgpu.Extent3D{
        Width:              width,
        Height:             height,
        DepthOrArrayLayers: 1,
    },
    Format:        wgpu.TextureFormatRGBA8Unorm,
    MipLevelCount: 1,
    SampleCount:   1,
})
```

- **Usage flags**:
  - `TextureUsageRenderAttachment`: Can be used as render target
  - `TextureUsageTextureBinding`: Can be sampled in future shaders (not shown in this example)
- **Format**: `RGBA8Unorm` for simplicity (can use float formats like `RGBA16Float` for HDR)

### Pipeline with Multiple Targets

```go
Fragment: &wgpu.FragmentState{
    Module:     shader,
    EntryPoint: "fs_main",
    Targets: []wgpu.ColorTargetState{
        {Format: wgpu.TextureFormatBGRA8Unorm, WriteMask: wgpu.ColorWriteMaskAll},
        {Format: wgpu.TextureFormatRGBA8Unorm, WriteMask: wgpu.ColorWriteMaskAll},
    },
}
```

- **Targets array**: Each element corresponds to `@location(N)` in fragment shader
- **Format matching**: Must match texture formats used in render pass
- **Write masks**: Control which channels are written (RGBA, RGB only, etc.)

### Render Pass with MRT

```go
ColorAttachments: []wgpu.RenderPassColorAttachment{
    {
        View:    surfaceTextureView,  // @location(0)
        LoadOp:  wgpu.LoadOpClear,
        StoreOp: wgpu.StoreOpStore,
        ClearValue: wgpu.Color{R: 0.1, G: 0.2, B: 0.3, A: 1.0},
    },
    {
        View:    extraTextureView,     // @location(1)
        LoadOp:  wgpu.LoadOpClear,
        StoreOp: wgpu.StoreOpStore,
        ClearValue: wgpu.Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
    },
}
```

- **Attachment order**: Index in array matches `@location(N)` in shader
- **View binding**: Each attachment points to a different texture view
- **Independent clear values**: Each target can have different clear color

## Building

```bash
cd examples/mrt
go build
```

## Running

```bash
./mrt
# or on Windows:
mrt.exe
```

You should see a window with a rotating triangle displaying colors (red/green/blue vertices). The second render target (position data) is written to an offscreen texture and not displayed in this simple example.

## Code Structure

- `createExtraTexture()`: Creates the second render target texture
- `createPipeline()`: Configures pipeline with two color targets
- `renderTriangle()`: Sets up render pass with two color attachments
- Shader: Uses `FragmentOutput` struct to output to both targets

## What's Being Rendered

### Target 0 (Screen)
- Red vertex at top
- Green vertex at bottom-left
- Blue vertex at bottom-right
- Rotates continuously

### Target 1 (Offscreen)
- Vertex positions remapped to [0, 1] range
- Top vertex: (0.5, 1.0, 0.0) → cyan-ish
- Bottom vertices: varying colors based on position

## Platform Support

- **Windows**: Full support (uses Win32 API for windowing)
- **Linux/macOS**: Windowing code needs adaptation

## Next Steps

To extend this example:

1. **Visualize Extra Target**: Display the second texture in a separate window or quad
2. **More Targets**: Add normal vectors, depth, or other G-buffer data
3. **Float Formats**: Use `TextureFormatRGBA16Float` or `RGBA32Float` for HDR
4. **Deferred Shading**: Implement a full deferred rendering pipeline
5. **Post-Processing**: Use MRT output as input to compute shaders

## Technical Notes

### Format Considerations

| Format | Use Case | Size per Pixel |
|--------|----------|----------------|
| `RGBA8Unorm` | Color, LDR data | 4 bytes |
| `RGBA16Float` | HDR color, normals | 8 bytes |
| `RGBA32Float` | High-precision data | 16 bytes |

### Performance

- **Bandwidth**: Writing to multiple targets increases memory bandwidth
- **Cache**: GPU can optimize writes to multiple targets in same pass
- **Best Practice**: Use smallest format that provides required precision

### WebGPU Limits

Most GPUs support **at least 4 color attachments**. Check device limits for maximum:

```go
// Query max color attachments (not shown in example)
limits := adapter.GetLimits()
maxColorAttachments := limits.MaxColorAttachments // typically 4-8
```

### Position Encoding

The example encodes 2D positions as colors:

```wgsl
// Remap [-1, 1] → [0, 1]
out.extra = vec4f(in.original_pos * 0.5 + 0.5, 0.0, 1.0);
```

- Top vertex (0.0, 0.5) → (0.5, 0.75, 0, 1) → cyan-white
- Left vertex (-0.5, -0.5) → (0.25, 0.25, 0, 1) → dark gray
- Right vertex (0.5, -0.5) → (0.75, 0.25, 0, 1) → reddish-gray

## Common Issues

### Format Mismatch
**Error**: "Render pass attachment format doesn't match pipeline"

**Solution**: Ensure formats in `FragmentState.Targets` match render pass attachment textures:

```go
// Pipeline
Targets: []wgpu.ColorTargetState{
    {Format: wgpu.TextureFormatBGRA8Unorm, ...},  // Must match surface format
    {Format: wgpu.TextureFormatRGBA8Unorm, ...},  // Must match extra texture format
}
```

### Wrong Number of Attachments
**Error**: "Number of color attachments doesn't match pipeline"

**Solution**: Render pass must have same number of attachments as pipeline targets.

### Missing Shader Outputs
**Error**: "Fragment shader doesn't write to all render targets"

**Solution**: Shader must output to all `@location(N)` values declared in pipeline.

## See Also

- [rotating-triangle](../rotating-triangle/) - Uniform buffers and animation (base for this example)
- [colored-triangle](../colored-triangle/) - Basic vertex colors
- [textured-quad](../textured-quad/) - Texture sampling (can be combined with MRT)

## References

- [WebGPU Spec: MRT](https://www.w3.org/TR/webgpu/#color-attachments)
- [WGSL Spec: Fragment Outputs](https://www.w3.org/TR/WGSL/#fragment-outputs)
- Deferred Rendering: [Learn OpenGL - Deferred Shading](https://learnopengl.com/Advanced-Lighting/Deferred-Shading)
