# Contributing to go-webgpu

Thank you for considering contributing to go-webgpu! This document outlines the development workflow and guidelines.

## Git Workflow (Git-Flow)

This project uses Git-Flow branching model for development.

### Branch Structure

```
main                 # Production-ready code (tagged releases)
  └─ develop         # Integration branch for next release
       ├─ feature/*  # New features
       ├─ bugfix/*   # Bug fixes
       └─ hotfix/*   # Critical fixes from main
```

### Branch Purposes

- **main**: Production-ready code. Only releases are merged here.
- **develop**: Active development branch. All features merge here first.
- **feature/\***: New features. Branch from `develop`, merge back to `develop`.
- **bugfix/\***: Bug fixes. Branch from `develop`, merge back to `develop`.
- **hotfix/\***: Critical production fixes. Branch from `main`, merge to both `main` and `develop`.

### Workflow Commands

#### Starting a New Feature

```bash
# Create feature branch from develop
git checkout develop
git pull origin develop
git checkout -b feature/my-new-feature

# Work on your feature...
git add .
git commit -m "feat: add my new feature"

# When done, merge back to develop
git checkout develop
git merge --squash feature/my-new-feature
git commit -m "feat: my new feature"
git branch -d feature/my-new-feature
git push origin develop
```

#### Fixing a Bug

```bash
# Create bugfix branch from develop
git checkout develop
git pull origin develop
git checkout -b bugfix/fix-issue-123

# Fix the bug...
git add .
git commit -m "fix: resolve issue #123"

# Merge back to develop
git checkout develop
git merge --squash bugfix/fix-issue-123
git commit -m "fix: resolve issue #123"
git branch -d bugfix/fix-issue-123
git push origin develop
```

#### Creating a Release

```bash
# Create release branch from develop
git checkout develop
git pull origin develop
git checkout -b release/v0.2.0

# Update version numbers, CHANGELOG, etc.
git add .
git commit -m "chore: prepare release v0.2.0"

# Merge to main and tag
git checkout main
git merge --no-ff release/v0.2.0
git tag -a v0.2.0 -m "Release v0.2.0"

# Merge back to develop
git checkout develop
git merge --no-ff release/v0.2.0

# Delete release branch
git branch -d release/v0.2.0

# Push everything
git push origin main develop --tags
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
- **style**: Code style changes (formatting, etc.)
- **refactor**: Code refactoring
- **test**: Adding or updating tests
- **chore**: Maintenance tasks (build, dependencies, etc.)
- **perf**: Performance improvements

### Examples

```bash
feat: add Texture3D support
fix: correct buffer mapping alignment
docs: update README with compute shader example
refactor: simplify device creation flow
test: add render pipeline tests
chore: update wgpu-native to v24.0.1
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

- Go 1.21 or later
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
│   └── workflows/       # CI/CD pipelines
├── wgpu/                 # WebGPU bindings (PUBLIC)
│   ├── adapter.go       # GPU adapter
│   ├── buffer.go        # Buffer management
│   ├── command.go       # Command encoder/buffer
│   ├── device.go        # GPU device
│   ├── instance.go      # WebGPU instance
│   ├── loader_*.go      # Platform-specific loaders
│   ├── pipeline.go      # Compute/render pipelines
│   ├── render.go        # Render pass
│   ├── surface_*.go     # Platform-specific surfaces
│   ├── texture.go       # Texture management
│   └── *_test.go        # Tests
├── examples/             # Usage examples
│   ├── triangle/        # Basic rendering
│   ├── compute/         # Compute shaders
│   ├── cube/            # 3D with depth buffer
│   └── ...              # More examples
├── include/              # WebGPU C headers (reference)
├── scripts/              # Development scripts
│   ├── download-wgpu-native.sh
│   └── pre-release-check.sh
├── docs/                 # Documentation
├── CHANGELOG.md          # Version history
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

- Always check and handle errors
- Use descriptive error messages with context
- Return errors immediately, don't wrap unnecessarily
- Validate inputs before FFI calls

### Testing

- Use table-driven tests when appropriate
- Test both success and error cases
- Test on all supported platforms
- Add examples for new features
- Compare with wgpu-native C examples for correctness

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
- Currently x86_64 only (goffi ARM64 support pending)
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
