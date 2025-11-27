# Rotating 3D Cube Example

This example demonstrates a **3D rotating cube with depth buffer** using go-webgpu.

## Features

- **3D Cube Geometry**: 36 vertices (6 faces × 2 triangles × 3 vertices per triangle)
- **Colored Faces**: Each face has a different color (Red, Green, Blue, Yellow, Cyan, Magenta)
- **MVP Matrices**: Model-View-Projection transformation using `wgpu.Mat4*` helpers
  - **Model**: Rotation around Y and X axes (animated)
  - **View**: Camera positioned at (0, 0, 3) looking at origin
  - **Projection**: Perspective projection with 45° FOV
- **Depth Buffer**: Proper depth testing using `TextureFormatDepth24Plus`
- **Back-face Culling**: Only front faces are rendered for performance
- **Uniform Buffer**: MVP matrix updated each frame via `queue.WriteBufferRaw()`

## Vertex Layout

Each vertex has 6 floats (24 bytes):
```
Position (vec3f): 12 bytes
Color (vec3f):    12 bytes
Total:            24 bytes (ArrayStride)
```

## Shaders (WGSL)

```wgsl
struct Uniforms {
    mvp: mat4x4f,
}
@group(0) @binding(0) var<uniform> uniforms: Uniforms;

@vertex
fn vs_main(@location(0) pos: vec3f, @location(1) color: vec3f) -> VertexOutput {
    var out: VertexOutput;
    out.position = uniforms.mvp * vec4f(pos, 1.0);
    out.color = color;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4f {
    return vec4f(in.color, 1.0);
}
```

## Depth Testing Configuration

```go
DepthStencil: &wgpu.DepthStencilState{
    Format:            wgpu.TextureFormatDepth24Plus,
    DepthWriteEnabled: true,
    DepthCompare:      wgpu.CompareFunctionLess,
    // ... stencil settings
}
```

## Render Pass with Depth Attachment

```go
DepthStencilAttachment: &wgpu.RenderPassDepthStencilAttachment{
    View:              depthTextureView,
    DepthLoadOp:       wgpu.LoadOpClear,
    DepthStoreOp:      wgpu.StoreOpStore,
    DepthClearValue:   1.0,
    DepthReadOnly:     wgpu.False,
    // ... stencil settings
}
```

## Building and Running

### Windows
```bash
cd examples/cube
go build -o cube.exe
./cube.exe
```

### Linux/macOS
```bash
cd examples/cube
go build -o cube
./cube
```

## Controls

- **Close Window**: Press the X button or Alt+F4

## Expected Output

A window displaying a **rotating 3D cube** with each face having a different color:
- **Front** (Z+): Red
- **Back** (Z-): Green
- **Top** (Y+): Blue
- **Bottom** (Y-): Yellow
- **Right** (X+): Cyan
- **Left** (X-): Magenta

The cube rotates around both Y and X axes, with proper depth testing ensuring correct face visibility.

## Technical Details

### MVP Matrix Calculation

```go
// Model: rotation around Y and X
modelY := wgpu.Mat4RotateY(angleY)
modelX := wgpu.Mat4RotateX(angleX)
model := modelY.Mul(modelX)

// View: camera at (0, 0, 3) looking at origin
view := wgpu.Mat4LookAt(
    wgpu.Vec3{X: 0, Y: 0, Z: 3},
    wgpu.Vec3{X: 0, Y: 0, Z: 0},
    wgpu.Vec3{X: 0, Y: 1, Z: 0},
)

// Projection: perspective 45° FOV
aspect := float32(width) / float32(height)
projection := wgpu.Mat4Perspective(45.0*π/180.0, aspect, 0.1, 100.0)

// Combine: MVP = P * V * M
mvp := projection.Mul(view).Mul(model)
```

### Depth Texture Creation

Uses the helper function from `wgpu.Device`:
```go
depthTexture := device.CreateDepthTexture(width, height, wgpu.TextureFormatDepth24Plus)
depthTextureView := depthTexture.CreateView(nil)
```

## Related Examples

- **rotating-triangle**: 2D rotation with uniform buffer (no depth)
- **textured-quad**: Texture mapping (no 3D)
- **triangle**: Basic triangle rendering

## References

- **WebGPU Depth Testing**: https://www.w3.org/TR/webgpu/#depth-stencil-state
- **WGSL Matrix Types**: https://www.w3.org/TR/WGSL/#matrix-types
- **go-webgpu Math Helpers**: `wgpu/math.go`
