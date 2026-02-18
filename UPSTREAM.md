# Upstream Dependencies

This document tracks upstream dependencies, pinned versions, and compatibility for go-webgpu.

---

## Pinned Versions

| Dependency | Version | Commit | Date |
|------------|---------|--------|------|
| **wgpu-native** | [v27.0.4.0](https://github.com/gfx-rs/wgpu-native/releases/tag/v27.0.4.0) | [`768f15f`](https://github.com/gfx-rs/wgpu-native/commit/768f15f6ace8e4ec8e8720d5732b29e0b34250a8) | 2025-12-23 |
| **webgpu.h** | wgpu-native bundled | same as above | — |
| **goffi** | [v0.3.9](https://github.com/go-webgpu/goffi/releases/tag/v0.3.9) | [`aa78271`](https://github.com/go-webgpu/goffi/commit/aa782710c349c09cebe2e5b9f76df859512884ef) | 2026-02-18 |
| **gputypes** | [v0.2.0](https://github.com/gogpu/gputypes/releases/tag/v0.2.0) | [`146b8b2`](https://github.com/gogpu/gputypes/commit/146b8b253ad16fe23db83cc593601081d009e3a6) | 2026-01-29 |

## Compatibility Matrix

| go-webgpu | wgpu-native | goffi | gputypes | Go |
|-----------|-------------|-------|----------|----|
| v0.3.1 | v27.0.4.0 | v0.3.9 | v0.2.0 | 1.25+ |
| v0.3.0 | v27.0.4.0 | v0.3.8 | v0.2.0 | 1.25+ |
| v0.2.1 | v27.0.4.0 | v0.3.8 | v0.2.0 | 1.25+ |
| v0.2.0 | v27.0.4.0 | v0.3.7 | v0.2.0 | 1.25+ |
| v0.1.1 | v24.0.3.1 | v0.3.3 | — | 1.23+ |
| v0.1.0 | v24.0.3.1 | v0.3.1 | — | 1.23+ |

## Binary Dependencies

wgpu-native is a **binary dependency** (not a Go module). Users must have the shared library available at runtime:

| Platform | Library | Location |
|----------|---------|----------|
| Windows | `wgpu_native.dll` | Same directory as executable or `%PATH%` |
| Linux | `libwgpu_native.so` | `LD_LIBRARY_PATH` or `/usr/local/lib` |
| macOS | `libwgpu_native.dylib` | `DYLD_LIBRARY_PATH` or `/usr/local/lib` |

Download pre-built binaries: [wgpu-native releases](https://github.com/gfx-rs/wgpu-native/releases)

## Upstream Repositories

| Project | Role | URL |
|---------|------|-----|
| **wgpu-native** | C API to wgpu (Rust) | https://github.com/gfx-rs/wgpu-native |
| **wgpu** | Rust WebGPU implementation | https://github.com/gfx-rs/wgpu |
| **webgpu-headers** | WebGPU C header spec | https://github.com/webgpu-native/webgpu-headers |
| **goffi** | Pure Go FFI library | https://github.com/go-webgpu/goffi |
| **gputypes** | Shared WebGPU type defs | https://github.com/gogpu/gputypes |

## Our Contributions to Upstream

| PR / Issue | Project | Status | Description |
|------------|---------|--------|-------------|
| [PR #543](https://github.com/gfx-rs/wgpu-native/pull/543) | wgpu-native | **Merged** | docs: update Go binding to go-webgpu/webgpu |
| [Issue #546 comment](https://github.com/gfx-rs/wgpu-native/issues/546) | wgpu-native | Open | webgpu.h / wgpu.h header inconsistencies |

## Breaking Changes Watch

### wgpu-native

wgpu follows a ~12-week release cadence. Major releases regularly include breaking changes. The project has **no SemVer stability guarantee** yet.

**Key risks for go-webgpu:**

- **webgpu.h header updates** — enum values, struct layouts, function signatures may change
- **wgpu.h extensions** — wgpu-native specific API beyond the standard header
- **Enum value remapping** — our `convert.go` maps between gputypes (spec) and wgpu-native (internal)

**Tracking:**
- Watch [wgpu-native releases](https://github.com/gfx-rs/wgpu-native/releases)
- Watch [wgpu CHANGELOG](https://github.com/gfx-rs/wgpu/blob/trunk/CHANGELOG.md)
- Watch [webgpu-headers](https://github.com/webgpu-native/webgpu-headers) for spec changes
- Issue [#3](https://github.com/go-webgpu/webgpu/issues/3) tracks webgpu-headers upgrade in wgpu-native

### gputypes

Enum values in gputypes follow the webgpu.h specification. When gputypes updates to match a new spec revision, our `convert.go` must be updated to maintain the mapping.

## Update Checklist

### Updating wgpu-native

1. Download new release binaries from [wgpu-native releases](https://github.com/gfx-rs/wgpu-native/releases)
2. Compare `ffi/webgpu.h` and `ffi/wgpu.h` with previous version for breaking changes
3. Update struct layouts in `wgpu/types.go` if struct sizes changed
4. Update enum mappings in `wgpu/convert.go` if enum values changed
5. Update function signatures in binding files if API changed
6. Run `go test ./wgpu/...` locally (with GPU) to verify
7. Update this file with new pinned version and commit
8. Update `CHANGELOG.md`

### Updating goffi

1. `go get github.com/go-webgpu/goffi@vX.Y.Z`
2. `go mod tidy`
3. Run tests
4. Update this file

### Updating gputypes

1. `go get github.com/gogpu/gputypes@vX.Y.Z`
2. Check if new enum values were added — update `convert.go`
3. `go mod tidy`
4. Run tests
5. Update this file

---

*Last updated: 2026-02-18 (v0.3.1)*
