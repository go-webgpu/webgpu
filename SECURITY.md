# Security Policy

## Supported Versions

go-webgpu is currently in initial release (v0.x.x). We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1.0 | :x:                |

Future stable releases (v1.0+) will follow semantic versioning with LTS support.

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability in go-webgpu, please report it responsibly.

### How to Report

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, please report security issues by:

1. **Private Security Advisory** (preferred):
   https://github.com/go-webgpu/webgpu/security/advisories/new

2. **Email** to maintainers:
   Create a private GitHub issue or contact via discussions

### What to Include

Please include the following information in your report:

- **Description** of the vulnerability
- **Steps to reproduce** the issue
- **Affected versions** (which versions are impacted)
- **Potential impact** (memory corruption, crashes, GPU resource leaks, etc.)
- **Suggested fix** (if you have one)
- **Your contact information** (for follow-up questions)

### Response Timeline

- **Initial Response**: Within 48-72 hours
- **Triage & Assessment**: Within 1 week
- **Fix & Disclosure**: Coordinated with reporter

We aim to:
1. Acknowledge receipt within 72 hours
2. Provide an initial assessment within 1 week
3. Work with you on a coordinated disclosure timeline
4. Credit you in the security advisory (unless you prefer to remain anonymous)

## Security Considerations for WebGPU Bindings

go-webgpu provides Go bindings to wgpu-native, which interfaces with GPU hardware. This introduces security considerations that users should be aware of.

### 1. FFI and Unsafe Code

**Risk**: go-webgpu uses FFI (Foreign Function Interface) to call native wgpu library functions.

**Attack Vectors**:
- Incorrect pointer handling in FFI calls
- Memory corruption from mismatched struct layouts
- Use-after-free from incorrect resource lifetime management

**Mitigation**:
- Careful struct layout matching with C headers
- Explicit resource cleanup with Release/Drop methods
- Extensive testing on all platforms
- golangci-lint with FFI-aware configuration

**User Recommendations**:
```go
// Always release GPU resources when done
device := adapter.CreateDevice(nil)
defer device.Release()

buffer := device.CreateBuffer(&wgpu.BufferDescriptor{...})
defer buffer.Release()
```

### 2. GPU Resource Exhaustion

**Risk**: Improper resource management can exhaust GPU memory or cause driver crashes.

**Attack Vectors**:
- Creating buffers/textures without releasing them
- Infinite loops in compute shaders
- Excessive command buffer submissions

**Mitigation**:
- Explicit resource lifetime management
- wgpu-native's built-in validation layer
- Device lost callbacks for error recovery

**User Best Practices**:
```go
// Enable validation in development
instance := wgpu.CreateInstance(nil)

// Use device lost callback
device := adapter.CreateDevice(&wgpu.DeviceDescriptor{
    DeviceLostCallback: func(reason wgpu.DeviceLostReason, message string) {
        log.Printf("Device lost: %v - %s", reason, message)
    },
})

// Always release resources
defer texture.Release()
defer buffer.Release()
defer pipeline.Release()
```

### 3. Shader Security

**Risk**: WGSL shaders execute on GPU and could potentially cause issues.

**Attack Vectors**:
- Malformed WGSL causing driver crashes
- Infinite loops in shaders (GPU hang)
- Out-of-bounds buffer access in shaders

**Mitigation**:
- wgpu-native validates all WGSL shaders
- WebGPU spec mandates bounds checking
- No access to system resources from shaders
- Shader compilation errors are returned to application

**User Recommendations**:
```go
// Always check shader compilation errors
shaderModule, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
    WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{
        Code: wgslSource,
    },
})
if err != nil {
    log.Printf("Shader compilation failed: %v", err)
    return err
}
```

### 4. Buffer Mapping Security

**Risk**: Mapped buffers provide direct memory access.

**Attack Vectors**:
- Reading unmapped buffer memory
- Writing beyond buffer bounds
- Using mapped pointer after unmap

**Mitigation**:
- WebGPU spec enforces mapping state machine
- Bounds checking in GetMappedRange
- wgpu-native validation

**User Best Practices**:
```go
// Wait for map to complete before accessing
buffer.MapAsync(wgpu.MapModeRead, 0, size, func(status wgpu.BufferMapAsyncStatus) {
    if status != wgpu.BufferMapAsyncStatusSuccess {
        return
    }

    // Only access after successful map
    data := buffer.GetMappedRange(0, size)
    // ... use data ...

    buffer.Unmap()
    // Don't use 'data' after Unmap!
})
```

### 5. Surface and Window Handling

**Risk**: Surface creation requires platform-specific window handles.

**Attack Vectors**:
- Invalid window handle causing crashes
- Surface use after window destruction
- Cross-platform handle confusion

**Mitigation**:
- Platform-specific surface creation functions
- Validation of window handles where possible
- Clear documentation of lifetime requirements

**User Recommendations**:
```go
// Ensure window is valid before creating surface
surface, err := instance.CreateSurface(&wgpu.SurfaceDescriptor{
    WindowsHWND: &wgpu.SurfaceDescriptorFromWindowsHWND{
        Hwnd: hwnd,  // Must be valid HWND
    },
})
if err != nil {
    return err
}

// Release surface before destroying window
surface.Release()
// Then destroy window
```

### 6. Dynamic Library Loading

**Risk**: go-webgpu loads wgpu-native as a dynamic library.

**Attack Vectors**:
- DLL hijacking (malicious library in search path)
- Library version mismatch
- Missing library dependencies

**Mitigation**:
- Explicit library paths recommended
- Version checking at load time
- Clear error messages for missing libraries

**User Best Practices**:
```go
// Set explicit library path if needed
os.Setenv("WGPU_NATIVE_PATH", "/path/to/wgpu_native.dll")

// Or place library in application directory
```

## Security Best Practices for Users

### Resource Management

Always release GPU resources:

```go
device := adapter.CreateDevice(nil)
defer device.Release()

// Create resources
buffer := device.CreateBuffer(...)
texture := device.CreateTexture(...)
pipeline := device.CreateRenderPipeline(...)

// Release in reverse order
defer pipeline.Release()
defer texture.Release()
defer buffer.Release()
```

### Error Handling

Always check errors:

```go
// Check adapter request
adapter, err := instance.RequestAdapter(...)
if err != nil {
    return fmt.Errorf("failed to get adapter: %w", err)
}

// Check device creation
device, err := adapter.CreateDevice(...)
if err != nil {
    return fmt.Errorf("failed to create device: %w", err)
}

// Check shader compilation
module, err := device.CreateShaderModule(...)
if err != nil {
    return fmt.Errorf("shader error: %w", err)
}
```

### Error Scopes

Use error scopes for detailed GPU error information:

```go
device.PushErrorScope(wgpu.ErrorFilterValidation)

// ... GPU operations ...

device.PopErrorScope(func(err wgpu.Error) {
    if err != nil {
        log.Printf("GPU validation error: %v", err)
    }
})
```

## Known Security Considerations

### 1. Relies on wgpu-native Security

**Status**: Dependency on upstream security.

**Risk Level**: Low

**Description**: go-webgpu's security depends on wgpu-native. We track wgpu-native releases and update accordingly.

### 2. FFI Boundary

**Status**: Careful implementation with testing.

**Risk Level**: Medium

**Description**: FFI calls require careful pointer and memory management. We maintain extensive tests and linter checks.

### 3. Platform Differences

**Status**: Platform-specific code paths tested.

**Risk Level**: Low

**Description**: Different platforms (Windows, Linux, macOS) have different loader and surface implementations.

## Dependency Security

go-webgpu dependencies:

| Dependency | Purpose | Security Notes |
|------------|---------|----------------|
| github.com/go-webgpu/goffi | Pure-Go FFI | Minimal, no CGO |
| wgpu-native (external) | WebGPU implementation | Mozilla/gfx-rs maintained |

**Monitoring**:
- Dependabot enabled for Go dependencies
- Track wgpu-native releases for security updates

## Security Testing

### Current Testing

- Unit tests for all public APIs
- Platform-specific tests (Windows, Linux, macOS)
- Memory leak detection in examples
- golangci-lint with security linters

### Planned

- Fuzz testing for FFI boundaries
- Integration tests with various GPU drivers
- Security-focused code review

## Security Disclosure History

### v0.1.0 (2024-11-28)

**Initial release** - No security issues reported yet.

## Security Contact

- **GitHub Security Advisory**: https://github.com/go-webgpu/webgpu/security/advisories/new
- **Public Issues** (for non-sensitive bugs): https://github.com/go-webgpu/webgpu/issues
- **Discussions**: https://github.com/go-webgpu/webgpu/discussions

## Bug Bounty Program

go-webgpu does not currently have a bug bounty program. We rely on responsible disclosure from the security community.

If you report a valid security vulnerability:
- Public credit in security advisory (if desired)
- Acknowledgment in CHANGELOG
- Our gratitude and recognition in README

---

**Thank you for helping keep go-webgpu secure!**

*Security is a journey, not a destination. We continuously improve our security posture with each release.*
