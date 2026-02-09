# Contributing to go-webgpu

Thank you for considering contributing to go-webgpu! This document outlines the development workflow and guidelines.

## Git Workflow (Pull Request)

All changes to `main` branch **must** go through Pull Requests. The `main` branch is protected.

### Branch Structure

```
main                 # Protected. Production-ready code (tagged releases)
  ├─ feat/*          # New features
  ├─ fix/*           # Bug fixes
  ├─ deps/*          # Dependency updates
  ├─ docs/*          # Documentation
  └─ hotfix/*        # Critical fixes
```

### Branch Protection

- **main** is protected — no direct pushes allowed
- All changes require a Pull Request
- Admins can bypass protection for emergency fixes

### Workflow Commands

#### Starting a New Feature

```bash
# Create feature branch from main
git checkout main
git pull origin main
git checkout -b feat/my-new-feature

# Work on your feature...
git add .
git commit -m "feat: add my new feature"

# Push branch and create PR
git push -u origin feat/my-new-feature
gh pr create --title "feat: add my new feature" --body "Description..."

# After PR is merged, clean up
git checkout main
git pull origin main
git branch -d feat/my-new-feature
```

#### Fixing a Bug

```bash
# Create fix branch from main
git checkout main
git pull origin main
git checkout -b fix/issue-123

# Fix the bug...
git add .
git commit -m "fix: resolve issue #123"

# Push and create PR
git push -u origin fix/issue-123
gh pr create --title "fix: resolve issue #123" --body "Closes #123"
```

#### Creating a Release

```bash
# After PR is merged, create release from main
git checkout main
git pull origin main

# Create tag and release
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
gh release create v0.2.0 --title "v0.2.0" --notes "Release notes..."
```

## Commit Message Guidelines

Follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **deps**: Dependency updates
- **style**: Code style changes (formatting, etc.)
- **refactor**: Code refactoring
- **test**: Adding or updating tests
- **chore**: Maintenance tasks (build, CI, etc.)
- **perf**: Performance improvements

### Examples

```bash
feat: add Texture3D support
fix: correct buffer mapping alignment
docs: update README with compute shader example
deps: update goffi v0.3.1 → v0.3.3
refactor: simplify device creation flow
test: add render pipeline tests
chore: update CI workflow
perf: optimize command encoder batch submission
```

## Code Quality Standards

### Before Committing

1. **Format code**:
   ```bash
   go fmt ./...
   ```

2. **Run linter**:
   ```bash
   golangci-lint run
   ```

3. **Run tests**:
   ```bash
   go test ./wgpu/...
   ```

4. **All-in-one** (use pre-release script):
   ```bash
   bash scripts/pre-release-check.sh
   ```

### Pull Request Requirements

- [ ] Code is formatted (`go fmt ./...`)
- [ ] Linter passes (`golangci-lint run` - 0 issues)
- [ ] All tests pass (`go test ./wgpu/...`)
- [ ] New code has tests (minimum 70% coverage)
- [ ] Documentation updated (if applicable)
- [ ] Commit messages follow conventions
- [ ] No sensitive data (credentials, tokens, etc.)
- [ ] Examples updated for new features

## Development Setup

### Prerequisites

- Go 1.25 or later
- golangci-lint v2
- wgpu-native shared libraries (download script provided)

### Platform Requirements

| Platform | Requirements |
|----------|-------------|
| Windows | Go + golangci-lint |
| Linux | Go + golangci-lint (CGO_ENABLED=0) |
| macOS | Go + golangci-lint (CGO_ENABLED=0, Intel x86_64 only) |

### Install Dependencies

```bash
# Clone repository
git clone https://github.com/go-webgpu/webgpu.git
cd webgpu

# Download wgpu-native libraries
bash scripts/download-wgpu-native.sh

# Download dependencies
go mod download

# Install golangci-lint
# See: https://golangci-lint.run/welcome/install/
```

### Running Tests

```bash
# Run all tests
go test ./wgpu/...

# Run with coverage
go test -cover ./wgpu/...

# Run specific test
go test -v ./wgpu/... -run "TestDeviceCreation"

# Run benchmarks
go test -bench=. -benchmem ./wgpu/...
```

### Running Linter

```bash
# Run linter
golangci-lint run

# Run with verbose output
golangci-lint run -v

# Verify config
golangci-lint config verify
```

## Project Structure

```
webgpu/
├── .github/              # GitHub workflows and templates
│   ├── CODEOWNERS       # Code ownership
│   └── workflows/       # CI/CD pipelines (test.yml, release.yml)
├── wgpu/                 # WebGPU bindings (PUBLIC API)
│   ├── doc.go           # Package-level documentation (godoc)
│   ├── types.go         # Core WebGPU handle types (Instance, Device, Buffer, ...)
│   ├── enums.go         # WebGPU enum types and constants
│   ├── wgpu.go          # Library initialization (Init, mustInit, checkInit)
│   ├── wgpu_errors.go   # Typed error system (WGPUError, sentinel errors)
│   ├── convert.go       # gputypes ↔ wgpu-native enum conversion
│   ├── loader.go        # Cross-platform library loading abstraction
│   ├── loader_*.go      # Platform-specific loaders (Windows/Unix)
│   ├── debug.go         # Resource leak detection (SetDebugMode, ReportLeaks)
│   ├── math.go          # 3D math helpers (Mat4, Vec3)
│   ├── instance.go      # WebGPU instance
│   ├── adapter.go       # GPU adapter (StringView, Limits, SupportedLimits)
│   ├── device.go        # GPU device (RequestDevice, GetQueue, GetLimits)
│   ├── buffer.go        # Buffer management (MapAsync, WriteBuffer)
│   ├── texture.go       # Texture management (CreateTexture, WriteTexture)
│   ├── shader.go        # Shader module (CreateShaderModuleWGSL)
│   ├── pipeline.go      # Compute pipeline, pipeline layout
│   ├── render_pipeline.go # Render pipeline
│   ├── bindgroup.go     # Bind groups and layouts
│   ├── command.go       # Command encoder/buffer, compute pass
│   ├── render.go        # Render pass encoder
│   ├── render_bundle.go # Render bundle encoder
│   ├── sampler.go       # Texture sampler
│   ├── queryset.go      # Query set (occlusion, timestamp)
│   ├── surface.go       # Surface (Configure, GetCurrentTexture, Present)
│   ├── surface_*.go     # Platform-specific surfaces (Windows/Linux/macOS)
│   ├── errors.go        # Error scope API (PushErrorScope, PopErrorScope)
│   ├── fuzz_test.go     # Fuzz tests for FFI boundary
│   └── *_test.go        # Tests
├── examples/             # Usage examples
│   ├── triangle/        # Basic rendering
│   ├── compute/         # Compute shaders
│   ├── cube/            # 3D with depth buffer
│   └── ...              # More examples
├── scripts/              # Development scripts
│   ├── download-wgpu-native.sh
│   └── pre-release-check.sh
├── CHANGELOG.md          # Version history
├── STABILITY.md          # API stability & deprecation policy
├── LICENSE               # MIT License
└── README.md             # Main documentation
```

## Adding New Features

1. Check if issue exists, if not create one
2. Discuss approach in the issue
3. Create feature branch from `develop`
4. Implement feature with tests
5. Update documentation and examples
6. Run quality checks (`bash scripts/pre-release-check.sh`)
7. Create pull request to `develop`
8. Wait for code review
9. Address feedback
10. Merge when approved

## Code Style Guidelines

### General Principles

- Follow Go conventions and idioms
- Write self-documenting code
- Add comments for complex FFI logic
- Keep functions small and focused
- Use meaningful variable names

### Naming Conventions

- **Public types/functions**: `PascalCase` (e.g., `CreateDevice`, `Buffer`)
- **Private types/functions**: `camelCase` (e.g., `loadLibrary`, `callProc`)
- **Constants**: `PascalCase` with context prefix (e.g., `BufferUsageVertex`, `TextureFormatRGBA8Unorm`)
- **Test functions**: `Test*` (e.g., `TestBufferCreation`)
- **Benchmark functions**: `Benchmark*` (e.g., `BenchmarkTextureUpload`)

### FFI-Specific Guidelines

```go
// Always check errors from FFI calls where applicable
result := proc.Call(args...)
if result == 0 {
    return nil, errors.New("FFI call failed")
}

// Use //nolint:errcheck for .Call() methods that don't return errors
proc.Call(ptr) //nolint:errcheck

// Document unsafe pointer usage
// SAFETY: ptr is valid for the lifetime of the buffer
unsafe.Pointer(ptr)

// Match C struct layouts exactly (no fieldalignment optimization)
type CStruct struct {
    Field1 uint32
    Field2 uint64  // 4-byte padding implied
}
```

### Error Handling

- Use `WGPUError` for all WebGPU errors (supports `errors.Is()`/`errors.As()`)
- Match against sentinel errors: `ErrValidation`, `ErrOutOfMemory`, `ErrInternal`
- Use `checkInit()` (returns error) for functions that return errors
- Use `mustInit()` (panics) only for void functions
- Validate inputs before FFI calls

### Resource Lifecycle

Every `Create*()` function must:
1. Call `trackResource(handle, "TypeName")` after successful creation
2. Have a corresponding `Release()` method that calls `untrackResource(handle)`

Use `SetDebugMode(true)` + `ReportLeaks()` during development to catch leaks.

### Testing

- Use table-driven tests when appropriate
- Test both success and error cases
- Test on all supported platforms
- Add examples for new features
- Compare with wgpu-native C examples for correctness

### Fuzz Testing

The FFI boundary is fuzz-tested via Go native fuzzing. When adding new conversion
functions, add corresponding fuzz targets in `fuzz_test.go`:

```bash
# Run fuzz tests on seed corpus (fast, part of normal test suite)
go test ./wgpu/... -run Fuzz

# Run actual fuzzing for a specific target
go test ./wgpu/ -fuzz=FuzzToWGPUTextureFormat -fuzztime=60s
```

### API Stability

See [STABILITY.md](STABILITY.md) for the API stability policy. When deprecating a function:

```go
// Deprecated: Use NewFunction instead.
func OldFunction() { ... }
```

## Platform-Specific Notes

### Windows
- Uses `syscall.LazyDLL` for dynamic loading
- Surface requires HWND from win32 window

### Linux
- Uses goffi for pure-Go FFI (CGO_ENABLED=0)
- Supports X11 and Wayland surfaces
- Requires wgpu-native .so in LD_LIBRARY_PATH

### macOS
- Uses goffi for pure-Go FFI (CGO_ENABLED=0)
- Supports both x86_64 and ARM64 (Apple Silicon)
- Surface requires Metal layer

## Getting Help

- Check existing issues and discussions
- Read ROADMAP.md for project direction
- Ask questions in GitHub Issues
- Reference wgpu-native documentation

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to go-webgpu!**
