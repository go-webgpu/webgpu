# Rotating Triangle Example

This example demonstrates how to use **uniform buffers** in go-webgpu to create animated graphics. It renders a continuously rotating colored triangle by updating a transformation matrix every frame.

## Features

- **Uniform Buffer**: Creates a 64-byte buffer for a 4x4 transformation matrix
- **Bind Groups**: Uses bind groups to bind the uniform buffer to the vertex shader
- **Animation Loop**: Updates the rotation matrix each frame via `queue.WriteBufferRaw()`
- **2D Rotation**: Implements Z-axis rotation using trigonometric functions
- **Vertex Colors**: Triangle vertices have red, green, and blue colors

## Key Concepts

### Uniform Buffer Creation

```go
// Size: mat4x4f = 16 floats * 4 bytes = 64 bytes
const uniformBufferSize = 64

uniformBuffer := device.CreateBuffer(&wgpu.BufferDescriptor{
    Usage:            wgpu.BufferUsageUniform | wgpu.BufferUsageCopyDst,
    Size:             uniformBufferSize,
    MappedAtCreation: wgpu.False,
})
```

- `BufferUsageUniform`: Marks buffer as uniform buffer
- `BufferUsageCopyDst`: Allows updating via `WriteBuffer` operations

### Bind Group Layout

```go
entries := []wgpu.BindGroupLayoutEntry{
    {
        Binding:    0,
        Visibility: wgpu.ShaderStageVertex,
        Buffer: wgpu.BufferBindingLayout{
            Type:             wgpu.BufferBindingTypeUniform,
            HasDynamicOffset: wgpu.False,
            MinBindingSize:   64, // mat4x4f = 64 bytes
        },
    },
}
```

- `Binding: 0`: Corresponds to `@binding(0)` in WGSL shader
- `Visibility: ShaderStageVertex`: Only accessible in vertex shader
- `BufferBindingTypeUniform`: Specifies uniform buffer binding

### Bind Group Creation

```go
entries := []wgpu.BindGroupEntry{
    wgpu.BufferBindingEntry(0, uniformBuffer, 0, 64),
}

bindGroup := device.CreateBindGroupSimple(bindGroupLayout, entries)
```

- Helper function `BufferBindingEntry` creates the entry
- Binds the entire 64-byte uniform buffer

### Matrix Update (Column-Major)

```go
// 2D rotation matrix in column-major layout
cos := float32(math.Cos(float64(angle)))
sin := float32(math.Sin(float64(angle)))

matrix := [16]float32{
    cos, sin, 0.0, 0.0,  // column 0
    -sin, cos, 0.0, 0.0, // column 1
    0.0, 0.0, 1.0, 0.0,  // column 2
    0.0, 0.0, 0.0, 1.0,  // column 3
}

queue.WriteBufferRaw(uniformBuffer, 0, unsafe.Pointer(&matrix[0]), 64)
```

- **Column-major layout**: Standard for graphics (matches WGSL/GLSL)
- **Rotation formula**:
  - `cos(θ)` and `sin(θ)` for rotation
  - Upper-left 2x2 submatrix rotates 2D vectors
  - Rest is identity matrix

### WGSL Shader

```wgsl
struct Uniforms {
    transform: mat4x4f,
}

@group(0) @binding(0) var<uniform> uniforms: Uniforms;

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.position = uniforms.transform * vec4f(in.position, 0.0, 1.0);
    out.color = in.color;
    return out;
}
```

- `@group(0) @binding(0)`: Matches bind group 0, binding 0
- `var<uniform>`: Declares uniform variable
- Matrix-vector multiplication rotates vertex position

## Building

```bash
cd examples/rotating-triangle
go build
```

## Running

```bash
./rotating-triangle
# or on Windows:
rotating-triangle.exe
```

You should see a window with a rotating triangle. The triangle has:
- Red vertex at the top
- Green vertex at bottom-left
- Blue vertex at bottom-right

The triangle rotates approximately 1 radian per second (about 57 degrees/second).

## Code Structure

- `createUniformBuffer()`: Creates 64-byte uniform buffer
- `createBindGroupLayout()`: Defines layout for uniform binding
- `createBindGroup()`: Binds uniform buffer to shader
- `updateUniformBuffer()`: Calculates rotation matrix and uploads to GPU
- `renderTriangle()`: Sets bind group and draws triangle

## Platform Support

- **Windows**: Full support (uses Win32 API for windowing)
- **Linux/macOS**: Windowing code needs to be adapted

## Next Steps

To extend this example:

1. **Multiple Uniforms**: Add scale, translation, or other transformations
2. **3D Rotation**: Extend to full 3D transformations (pitch, yaw, roll)
3. **Camera Matrix**: Add view and projection matrices for 3D rendering
4. **Dynamic Speed**: Allow user to control rotation speed with keyboard input

## Technical Notes

### Matrix Layout

WebGPU/WGSL uses **column-major** matrix layout. When writing a 4x4 matrix:

```
float32 array:  [m00, m10, m20, m30,  m01, m11, m21, m31, ...]
                  ↑              ↑      ↑              ↑
                column 0      column 1
```

For 2D rotation around Z-axis:
```
Math notation:           Memory layout (column-major):
| cos  -sin  0  0 |     [cos, sin, 0, 0,  -sin, cos, 0, 0, ...]
| sin   cos  0  0 |
|  0     0   1  0 |
|  0     0   0  1 |
```

### Performance Considerations

- Uniform buffers are updated every frame (~60 times/second)
- `WriteBufferRaw` is efficient for small updates (64 bytes)
- For larger data, consider double-buffering strategies

### Rotation Speed

Current rotation speed: `angle = elapsed_seconds` radians

- 1 radian ≈ 57.3 degrees
- Full rotation (2π radians) takes ~6.28 seconds

To change speed: `angle = elapsed * speed_factor`

## See Also

- [colored-triangle](../colored-triangle/) - Basic vertex colors without uniforms
- [triangle](../triangle/) - Minimal triangle without colors
- [textured-quad](../textured-quad/) - Texture mapping example
