# WebGPU Error Handling API

WebGPU error handling API for go-webgpu provides error scopes to catch and handle GPU errors during operations.

## Overview

Error scopes allow you to catch specific types of GPU errors:
- **Validation errors**: API misuse, invalid parameters
- **Out-of-memory errors**: Insufficient GPU memory
- **Internal errors**: Internal implementation errors

## API

### Types

```go
// ErrorFilter - filters which errors to catch
type ErrorFilter uint32
const (
    ErrorFilterValidation  ErrorFilter = 0x00000001
    ErrorFilterOutOfMemory ErrorFilter = 0x00000002
    ErrorFilterInternal    ErrorFilter = 0x00000003
)

// ErrorType - describes the error that occurred
type ErrorType uint32
const (
    ErrorTypeNoError     ErrorType = 0x00000001
    ErrorTypeValidation  ErrorType = 0x00000002
    ErrorTypeOutOfMemory ErrorType = 0x00000003
    ErrorTypeInternal    ErrorType = 0x00000004
    ErrorTypeUnknown     ErrorType = 0x00000005
)
```

### Functions

#### PushErrorScope

```go
func (d *Device) PushErrorScope(filter ErrorFilter)
```

Pushes an error scope to catch errors of the specified type.

**Parameters:**
- `filter`: Type of errors to catch

**Example:**
```go
device.PushErrorScope(wgpu.ErrorFilterValidation)
```

#### PopErrorScope

```go
func (d *Device) PopErrorScope(instance *Instance) (ErrorType, string)
```

Pops the current error scope and returns the first error caught (if any).

**Parameters:**
- `instance`: WebGPU instance (required for processing events)

**Returns:**
- `ErrorType`: Type of error (ErrorTypeNoError if no error)
- `string`: Error message (empty if no error)

**Panics** if:
- Error scope stack is empty
- Instance is nil
- Operation fails

**Example:**
```go
errType, message := device.PopErrorScope(instance)
if errType != wgpu.ErrorTypeNoError {
    log.Printf("GPU error: %s", message)
}
```

#### PopErrorScopeAsync

```go
func (d *Device) PopErrorScopeAsync(instance *Instance) (ErrorType, string, error)
```

Similar to PopErrorScope but returns an error instead of panicking.

**Returns:**
- `ErrorType`: Type of error (ErrorTypeNoError if no error)
- `string`: Error message
- `error`: Go error if operation failed (nil on success)

**Example:**
```go
errType, message, err := device.PopErrorScopeAsync(instance)
if err != nil {
    log.Printf("PopErrorScope failed: %v", err)
    return
}
if errType != wgpu.ErrorTypeNoError {
    log.Printf("GPU error: %s", message)
}
```

## Usage

### Basic Error Scope

```go
// Push error scope
device.PushErrorScope(wgpu.ErrorFilterValidation)

// Perform GPU operations
buffer := device.CreateBuffer(&desc)
// ... more operations ...

// Pop scope and check for errors
errType, message := device.PopErrorScope(instance)
if errType != wgpu.ErrorTypeNoError {
    log.Printf("Validation error: %s", message)
}
```

### Nested Error Scopes

Error scopes are LIFO (stack-based):

```go
device.PushErrorScope(wgpu.ErrorFilterValidation)  // Outer
device.PushErrorScope(wgpu.ErrorFilterValidation)  // Inner

// GPU operations...

// Pop in reverse order
errType, msg := device.PopErrorScope(instance)  // Inner first
errType, msg = device.PopErrorScope(instance)   // Then outer
```

### Different Error Filters

```go
// Monitor specific error types
device.PushErrorScope(wgpu.ErrorFilterOutOfMemory)
largeBuffer := device.CreateBuffer(&largeDesc)
errType, message := device.PopErrorScope(instance)
if errType == wgpu.ErrorTypeOutOfMemory {
    // Handle OOM error
    log.Printf("Out of memory: %s", message)
}
```

## Important Notes

### Stack Management

⚠️ **CRITICAL**: Always pop every pushed error scope!

```go
// ❌ BAD - unbalanced push/pop
device.PushErrorScope(wgpu.ErrorFilterValidation)
// ... forgot to pop ...

// ✅ GOOD - balanced push/pop
device.PushErrorScope(wgpu.ErrorFilterValidation)
defer device.PopErrorScope(instance) // Ensures pop happens
```

### Stack Underflow

⚠️ **Known Limitation**: Popping an empty error scope stack will cause a **panic** in wgpu-native.

```go
// ❌ BAD - will panic!
device.PopErrorScope(instance)  // No push before this

// ✅ GOOD - track push/pop manually
device.PushErrorScope(wgpu.ErrorFilterValidation)
device.PopErrorScope(instance)  // Balanced
```

If you need graceful handling of empty stack:
```go
errType, message, err := device.PopErrorScopeAsync(instance)
if err != nil {
    // Handle error (e.g., empty stack)
    log.Printf("Error: %v", err)
    return
}
```

### Instance Requirement

PopErrorScope requires an `Instance` to process events:

```go
// Keep instance alive
instance, _ := wgpu.CreateInstance(nil)
defer instance.Release()

// Use instance for PopErrorScope
errType, message := device.PopErrorScope(instance)
```

## Error Types

| ErrorType | Description |
|-----------|-------------|
| `ErrorTypeNoError` | No error occurred |
| `ErrorTypeValidation` | API misuse or invalid parameters |
| `ErrorTypeOutOfMemory` | Insufficient GPU memory |
| `ErrorTypeInternal` | Internal implementation error |
| `ErrorTypeUnknown` | Unknown error |

## Example

See [examples/error_handling](../../examples/error_handling/main.go) for a complete example.

## References

- [WebGPU Spec - Error Scopes](https://gpuweb.github.io/gpuweb/#error-scopes)
- [WebGPU Error Handling Best Practices](https://toji.dev/webgpu-best-practices/error-handling.html)
