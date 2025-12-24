# go-webgpu - Development Roadmap

> **Strategic Focus**: Production-grade Zero-CGO WebGPU bindings for Go

**Last Updated**: 2024-12-24 | **Current Version**: v0.1.1 | **Target**: v1.0.0 stable

---

## Vision

Build **production-ready, cross-platform WebGPU bindings** for Go with zero CGO dependency, enabling GPU-accelerated graphics and compute in pure Go applications.

### Current State vs Target

| Metric | Current (v0.1.1) | Target (v1.0.0) |
|--------|------------------|-----------------|
| Platforms | Windows, Linux, macOS (x64, arm64) | All major platforms |
| CGO Required | No (Zero-CGO) | No |
| API Coverage | ~80% WebGPU | 100% WebGPU |
| wgpu-native | v24.0.3.1 | Latest stable |
| Test Coverage | ~70% | 90%+ |
| Examples | 11 | 20+ |

---

## Release Strategy

```
v0.1.1 (Current) -> Hotfix: goffi PointerType bug + PR workflow
         |
v0.2.0 (Next) -> API improvements, builder patterns
         |
v0.3.0 -> Advanced features (storage textures, texture arrays)
         |
v0.4.0 -> Performance optimizations
         |
v0.5.0 -> Extended examples and documentation
         |
v1.0.0-rc -> Feature freeze, API locked
         |
v1.0.0 STABLE -> Production release with API stability guarantee
```

---

## v0.2.0 - API Improvements (NEXT)

**Goal**: Improve API ergonomics and developer experience

| ID | Feature | Impact | Status |
|----|---------|--------|--------|
| API-001 | Builder pattern for descriptors | Better ergonomics | Planned |
| API-002 | Error wrapping with context | Better debugging | Planned |
| API-003 | Resource tracking helpers | Memory management | Planned |

**Target**: Q1 2025

---

## v0.3.0 - Advanced Features (MEDIUM PRIORITY)

**Goal**: Complete WebGPU API coverage

| ID | Feature | Impact | Status |
|----|---------|--------|--------|
| FEAT-001 | Storage textures | Compute image processing | Planned |
| FEAT-002 | Texture arrays | Sprite sheets, cubemaps | Planned |
| FEAT-003 | Occlusion queries | Visibility testing | Planned |
| FEAT-004 | Pipeline statistics | Performance profiling | Planned |
| FEAT-005 | Multi-draw indirect | Batch rendering | Planned |

**Target**: Q2 2025

---

## v0.4.0 - Performance (MEDIUM PRIORITY)

**Goal**: Optimize hot paths and reduce allocations

| ID | Feature | Impact | Status |
|----|---------|--------|--------|
| PERF-001 | Command buffer pooling | Reduce allocations | Planned |
| PERF-002 | Descriptor caching | Faster pipeline creation | Planned |
| PERF-003 | Batch resource creation | Startup optimization | Planned |
| PERF-004 | Memory-mapped staging | Faster uploads | Planned |

**Target**: Q2 2025

---

## v0.5.0 - Examples & Documentation (MEDIUM PRIORITY)

**Goal**: Comprehensive learning resources

### New Examples

| ID | Example | Demonstrates |
|----|---------|--------------|
| EX-001 | PBR Renderer | Material system, lighting |
| EX-002 | Shadow Mapping | Depth textures, multi-pass |
| EX-003 | Post-processing | Framebuffers, effects |
| EX-004 | Particle System | Compute + render integration |
| EX-005 | Text Rendering | Texture atlases, SDF fonts |
| EX-006 | Deferred Shading | G-buffer, MRT |
| EX-007 | Ray Marching | Compute shaders |
| EX-008 | Image Processing | Compute filters |
| EX-009 | Physics Simulation | GPU compute |

### Documentation

| ID | Document | Content |
|----|----------|---------|
| DOC-001 | API Reference | Complete godoc |
| DOC-002 | Migration Guide | From other GPU libs |
| DOC-003 | Performance Guide | Optimization tips |
| DOC-004 | Troubleshooting | Common issues |

**Target**: Q3 2025

---

## v1.0.0 - Production Ready

**Requirements**:
- [ ] All v0.2.0-v0.5.0 features complete
- [ ] API stability guarantee
- [ ] Comprehensive documentation
- [ ] 90%+ test coverage
- [ ] Performance benchmarks
- [ ] Security review

**Guarantees**:
- API stability (no breaking changes in v1.x.x)
- Semantic versioning
- Long-term support

**Target**: Q4 2025

---

## Feature Comparison Matrix

| Feature | wgpu-rs | Dawn | go-webgpu v0.1 | go-webgpu v1.0 |
|---------|---------|------|----------------|----------------|
| Zero-CGO | N/A | N/A | Yes | Yes |
| Windows x64 | Yes | Yes | Yes | Yes |
| Linux x64 | Yes | Yes | Yes | Yes |
| Linux ARM64 | Yes | Yes | Yes | Yes |
| macOS x64 | Yes | Yes | Yes | Yes |
| macOS ARM64 | Yes | Yes | Yes | Yes |
| Buffer mapping | Yes | Yes | Yes | Yes |
| Compute shaders | Yes | Yes | Yes | Yes |
| Render bundles | Yes | Yes | Yes | Yes |
| Timestamp queries | Yes | Yes | Yes | Yes |
| Storage textures | Yes | Yes | No | Yes |
| Texture arrays | Yes | Yes | No | Yes |

---

## Current Examples (v0.1.x)

| Example | Features Demonstrated |
|---------|----------------------|
| Triangle | Basic rendering, shaders |
| Colored Triangle | Vertex attributes |
| Rotating Triangle | Uniform buffers, animation |
| Textured Quad | Texture sampling, UV coords |
| 3D Cube | Depth buffer, transforms, MVP |
| MRT | Multiple render targets |
| Compute | Compute shaders, storage buffers |
| Instanced | Instance rendering, vertex step mode |
| RenderBundle | Pre-recorded commands |
| Timestamp Query | GPU timing |
| Error Handling | Error scopes |

---

## Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| wgpu-native | v24.0.3.1 | WebGPU implementation |
| goffi | v0.3.3 | Pure-Go FFI (x64 + ARM64) |
| Go | 1.25+ | Language runtime |

### Upstream Tracking

- **wgpu-native**: Track releases for new features and security fixes
- **goffi**: Track for performance improvements and new platforms

---

## Out of Scope

**Not planned**:
- WebGL fallback (WebGPU only)
- DirectX 11 backend (wgpu-native uses D3D12)
- OpenGL backend (wgpu-native uses Vulkan/Metal)
- Custom shader language (WGSL only)

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to contribute to the roadmap.

Priority features are marked in GitHub Issues with labels:
- `priority:high` - Next release
- `priority:medium` - Future release
- `help-wanted` - Community contributions welcome

---

## Release History

| Version | Date | Type | Key Changes |
|---------|------|------|-------------|
| v0.1.1 | 2024-12-24 | Hotfix | goffi v0.3.3 (PointerType fix), PR workflow |
| v0.1.0 | 2024-11-28 | Initial | Core API, 11 examples, 5 platforms (x64 + ARM64) |

---

*Current: v0.1.1 | Next: v0.2.0 (API Improvements) | Target: v1.0.0 (Q4 2025)*
