# go-webgpu Roadmap

> **Mission**: Production-grade Zero-CGO WebGPU bindings for Go

[![GitHub Project](https://img.shields.io/badge/GitHub-Project%20Board-blue?style=flat-square&logo=github)](https://github.com/go-webgpu/webgpu/projects)
[![GitHub Issues](https://img.shields.io/github/issues/go-webgpu/webgpu?style=flat-square&logo=github)](https://github.com/go-webgpu/webgpu/issues)

---

## Disclaimer

> **This roadmap represents our current plans and priorities, not commitments.**
> Features, timelines, and priorities may change based on community feedback, technical constraints, and ecosystem developments. For the most current status, see our [GitHub Issues](https://github.com/go-webgpu/webgpu/issues) and [Project Board](https://github.com/go-webgpu/webgpu/projects).

---

## Vision

Enable **GPU-accelerated graphics and compute in pure Go** — no CGO, no complexity, just Go.

### Why go-webgpu?

| Challenge | Our Solution |
|-----------|--------------|
| CGO complexity | Zero-CGO via [goffi](https://github.com/go-webgpu/goffi) FFI |
| Cross-compilation pain | Pure Go builds for all platforms |
| WebGPU fragmentation | Unified API via [gputypes](https://github.com/gogpu/gputypes) |
| Vendor lock-in | Open source, part of [gogpu ecosystem](https://github.com/gogpu) |

---

## Current Status

| Metric | Status |
|--------|--------|
| **Latest Release** | ![GitHub Release](https://img.shields.io/github/v/release/go-webgpu/webgpu?style=flat-square) |
| **Platforms** | Windows, Linux, macOS (x64, arm64) |
| **API Coverage** | ~80% WebGPU |
| **Examples** | 11 working demos |
| **Test Coverage** | ~70% |

### Technology Stack

| Component | Version | Role |
|-----------|---------|------|
| [wgpu-native](https://github.com/gfx-rs/wgpu-native) | v27.0.4.0 | WebGPU implementation (Rust) |
| [goffi](https://github.com/go-webgpu/goffi) | v0.3.7 | Zero-CGO FFI layer |
| [gputypes](https://github.com/gogpu/gputypes) | latest | WebGPU type definitions |

---

## Roadmap Phases

We use GitHub labels to track feature progress:

| Label | Meaning |
|-------|---------|
| `phase:exploring` | Under consideration, gathering feedback |
| `phase:design` | Actively designing solution |
| `phase:development` | Implementation in progress |
| `phase:preview` | Available for testing |
| `phase:stable` | Production ready |

---

## Now: Stability & Ecosystem

**Focus**: Ensure rock-solid foundation for production use.

| Feature | Status | Issue |
|---------|--------|-------|
| gputypes integration | `stable` | — |
| wgpu-native v27 compatibility | `stable` | — |
| All 11 examples working | `stable` | — |
| Enum conversion layer | `stable` | — |

---

## Next: Advanced Features

**Focus**: Complete WebGPU API coverage.

| Feature | Status | Issue |
|---------|--------|-------|
| Storage textures | `exploring` | [#TBD](https://github.com/go-webgpu/webgpu/issues) |
| Texture arrays | `exploring` | [#TBD](https://github.com/go-webgpu/webgpu/issues) |
| Occlusion queries | `exploring` | [#TBD](https://github.com/go-webgpu/webgpu/issues) |
| Pipeline statistics | `exploring` | [#TBD](https://github.com/go-webgpu/webgpu/issues) |
| Multi-draw indirect | `exploring` | [#TBD](https://github.com/go-webgpu/webgpu/issues) |

---

## Later: Performance & DX

**Focus**: Optimize performance and developer experience.

| Feature | Status | Issue |
|---------|--------|-------|
| Builder pattern for descriptors | `exploring` | — |
| Command buffer pooling | `exploring` | — |
| Descriptor caching | `exploring` | — |
| Memory-mapped staging | `exploring` | — |
| Error wrapping with context | `exploring` | — |

---

## Future: Extended Examples

| Example | Demonstrates | Status |
|---------|--------------|--------|
| PBR Renderer | Material system, lighting | `exploring` |
| Shadow Mapping | Depth textures, multi-pass | `exploring` |
| Post-processing | Framebuffers, effects | `exploring` |
| Particle System | Compute + render | `exploring` |
| Text Rendering | SDF fonts, atlases | `exploring` |
| Deferred Shading | G-buffer, MRT | `exploring` |

---

## v1.0 Requirements

Before we tag v1.0.0 stable:

- [ ] 100% WebGPU API coverage
- [ ] 90%+ test coverage
- [ ] Comprehensive documentation
- [ ] Performance benchmarks
- [ ] Security review
- [ ] API stability guarantee

**v1.0 Guarantees**:
- No breaking changes in v1.x.x
- Semantic versioning
- Long-term support commitment

---

## Out of Scope

Features we **do not plan** to implement:

| Feature | Reason |
|---------|--------|
| WebGL fallback | WebGPU-only library |
| DirectX 11 backend | wgpu-native uses D3D12 |
| OpenGL backend | wgpu-native uses Vulkan/Metal |
| Custom shader language | WGSL standard only |
| Browser support | Native applications only |

---

## How to Contribute

We welcome contributions! Here's how to get involved:

### 1. Find an Issue

- [`good-first-issue`](https://github.com/go-webgpu/webgpu/labels/good-first-issue) — Great for newcomers
- [`help-wanted`](https://github.com/go-webgpu/webgpu/labels/help-wanted) — Community contributions welcome
- [`priority:high`](https://github.com/go-webgpu/webgpu/labels/priority%3Ahigh) — Most impactful work

### 2. Propose Features

Open a [Feature Request](https://github.com/go-webgpu/webgpu/issues/new?template=feature_request.md) to discuss before implementing.

### 3. Submit PRs

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### 4. Join Discussion

- [GitHub Discussions](https://github.com/go-webgpu/webgpu/discussions) — Questions, ideas, showcase
- [Issues](https://github.com/go-webgpu/webgpu/issues) — Bug reports, feature requests

---

## Upstream Dependencies

We track these projects for updates:

| Project | What We Track | Our Issue |
|---------|---------------|-----------|
| [wgpu-native](https://github.com/gfx-rs/wgpu-native) | Releases, security fixes | [#3](https://github.com/go-webgpu/webgpu/issues/3) |
| [webgpu-headers](https://github.com/webgpu-native/webgpu-headers) | Spec changes | [#3](https://github.com/go-webgpu/webgpu/issues/3) |
| [goffi](https://github.com/go-webgpu/goffi) | Performance, platforms | — |
| [gputypes](https://github.com/gogpu/gputypes) | Type definitions | — |

---

## Release History

| Version | Date | Highlights |
|---------|------|------------|
| **v0.2.0** | 2026-01-29 | gputypes integration, wgpu-native v27, all examples fixed |
| v0.1.4 | 2026-01-03 | goffi v0.3.7 (ARM64 Darwin) |
| v0.1.3 | 2025-12-29 | goffi v0.3.6 (ARM64 HFA fix) |
| v0.1.2 | 2025-12-27 | goffi v0.3.5 |
| v0.1.1 | 2024-12-24 | goffi hotfix, PR workflow |
| v0.1.0 | 2024-11-28 | Initial release, 11 examples, 5 platforms |

See [CHANGELOG.md](CHANGELOG.md) for detailed release notes.

---

## Related Projects

| Project | Description |
|---------|-------------|
| [gogpu](https://github.com/gogpu) | Pure Go WebGPU ecosystem |
| [gputypes](https://github.com/gogpu/gputypes) | Shared WebGPU type definitions |
| [goffi](https://github.com/go-webgpu/goffi) | Zero-CGO FFI for Go |

---

<p align="center">
  <sub>This roadmap is inspired by <a href="https://github.com/github/roadmap">GitHub's public roadmap</a> practices.</sub>
</p>
