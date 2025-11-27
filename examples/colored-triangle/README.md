# Colored Triangle Example

Демонстрация рендеринга разноцветного треугольника с использованием vertex buffers.

## Описание

Этот пример показывает:
- Создание и загрузку vertex buffer с позициями и цветами вершин
- Использование `VertexBufferLayout` и `VertexAttribute` для описания структуры вершин
- Интерполяцию цветов между вершинами (GPU автоматически интерполирует varying переменные)
- Передачу данных из vertex shader в fragment shader через `@location`

## Vertex Buffer Layout

```
ArrayStride: 20 bytes (5 floats * 4 bytes)
┌─────────────────────────────────────────────┐
│ Vertex 0                                    │
│ ┌────────────┬──────────────────────────┐  │
│ │ Position   │ Color                    │  │
│ │ (x, y)     │ (r, g, b)                │  │
│ │ 2 floats   │ 3 floats                 │  │
│ │ offset: 0  │ offset: 8                │  │
│ └────────────┴──────────────────────────┘  │
└─────────────────────────────────────────────┘

Vertex 0: (0.0, 0.5)   - Red   (1.0, 0.0, 0.0)
Vertex 1: (-0.5, -0.5) - Green (0.0, 1.0, 0.0)
Vertex 2: (0.5, -0.5)  - Blue  (0.0, 0.0, 1.0)
```

## Shader Pipeline

### Vertex Shader Input
```wgsl
@location(0) position: vec2f  // from attribute 0
@location(1) color: vec3f     // from attribute 1
```

### Vertex to Fragment (varying)
```wgsl
@location(0) color: vec3f     // interpolated automatically
```

### Fragment Shader Output
```wgsl
@location(0) vec4f            // RGBA color to render target
```

## Сборка и запуск

```bash
# Компиляция
go build

# Запуск
./colored-triangle.exe  # Windows
./colored-triangle      # Linux/macOS
```

## Требования

- **wgpu-native.dll** в PATH или рядом с исполняемым файлом
- Поддержка WebGPU-совместимого GPU (DX12, Vulkan, Metal)
- Windows: DirectX 12
- Linux: Vulkan
- macOS: Metal

## Технические детали

### Vertex Buffer Creation
```go
vertices := []float32{
    0.0,  0.5,  1.0, 0.0, 0.0,  // Top: red
   -0.5, -0.5,  0.0, 1.0, 0.0,  // Left: green
    0.5, -0.5,  0.0, 0.0, 1.0,  // Right: blue
}

buffer := device.CreateBuffer(&wgpu.BufferDescriptor{
    Usage:            wgpu.BufferUsageVertex | wgpu.BufferUsageCopyDst,
    Size:             60, // 15 floats * 4 bytes
    MappedAtCreation: wgpu.True,
})
```

### Pipeline Configuration
```go
Vertex: wgpu.VertexState{
    Buffers: []wgpu.VertexBufferLayout{{
        ArrayStride: 20, // bytes per vertex
        StepMode:    wgpu.VertexStepModeVertex,
        Attributes: []wgpu.VertexAttribute{
            {Format: wgpu.VertexFormatFloat32x2, Offset: 0, ShaderLocation: 0},
            {Format: wgpu.VertexFormatFloat32x3, Offset: 8, ShaderLocation: 1},
        },
    }},
}
```

### Rendering
```go
pass.SetPipeline(pipeline)
pass.SetVertexBuffer(0, vertexBuffer, 0, 60)
pass.Draw(3, 1, 0, 0)  // 3 vertices, 1 instance
```

## Следующие шаги

После освоения этого примера смотрите:
- `rotating-triangle` - добавление uniform buffers и анимации
- `textured-cube` - 3D трансформации и текстуры
- `indexed-geometry` - использование index buffers

## Отладка

### Треугольник не отображается
- Проверьте порядок вершин (CCW by default)
- Убедитесь что ArrayStride соответствует размеру данных
- Проверьте Offset для каждого атрибута

### Неправильные цвета
- Проверьте порядок компонент (r, g, b)
- Убедитесь что ShaderLocation соответствуют @location в шейдере
- Проверьте Format атрибутов (Float32x2 vs Float32x3)

### Ошибки валидации
Включите debug mode:
```go
instance, _ := wgpu.CreateInstance(&wgpu.InstanceDescriptor{
    // Validation будет включен автоматически в debug builds
})
```
