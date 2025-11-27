# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2024-11-28

### Added

#### Core API
- WebGPU Instance, Adapter, Device creation
- Buffer creation and management (MapAsync, GetMappedRange, Unmap)
- Texture creation and management (2D, depth, render targets)
- Sampler API with filtering and address modes
- Shader module compilation (WGSL)

#### Pipelines
- Compute Pipeline with workgroups
- Render Pipeline with vertex/fragment stages
- Pipeline Layout and Bind Group Layout
- Bind Groups for resource binding

#### Rendering
- Command Encoder and Queue submission
- Render Pass with color and depth attachments
- Vertex buffers with custom layouts
- Index buffers (Uint16, Uint32)
- Depth buffer support (DepthStencilState)
- MRT (Multiple Render Targets) support
- RenderBundle API for pre-recording render commands

#### Advanced Features
- Instanced rendering (VertexStepModeInstance)
- Indirect drawing (DrawIndirect, DrawIndexedIndirect)
- Indirect compute dispatch (DispatchWorkgroupsIndirect)
- QuerySet API for GPU timestamps
- Error Handling API (PushErrorScope, PopErrorScope)

#### Cross-Platform
- Windows support via syscall.LazyDLL
- Linux support via goffi (CGO_ENABLED=0)
- macOS support via goffi (CGO_ENABLED=0)
- Platform-specific Surface creation (HWND, X11, Wayland, Metal)

#### Math Helpers
- Mat4 (4x4 matrix) with Identity, Translate, Scale, Rotate, Perspective, LookAt
- Vec3, Vec4 vector types with common operations

#### Examples
- Triangle (basic rendering)
- Colored triangle (vertex colors)
- Rotating triangle (uniform buffers)
- Textured quad (texture sampling)
- 3D cube (depth buffer, transforms)
- MRT (multiple render targets)
- Compute shader (parallel processing)
- Instanced rendering (25 objects in 1 draw call)
- RenderBundle (pre-recorded commands)
- Timestamp query (GPU timing)
- Error handling (error scopes)

#### Infrastructure
- Zero-CGO architecture using goffi for FFI
- GitHub Actions CI/CD (Linux, macOS, Windows)
- golangci-lint configuration for FFI code
- Pre-release validation script
- Comprehensive test suite (76+ tests)

### Dependencies
- github.com/go-webgpu/goffi v0.3.1
- wgpu-native v24.0.0.2
