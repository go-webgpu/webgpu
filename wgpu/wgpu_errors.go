package wgpu

import "fmt"

// Sentinel errors for programmatic error handling via errors.Is().
//
// Usage:
//
//	if errors.Is(err, wgpu.ErrValidation) {
//	    // handle validation error
//	}
var (
	// ErrValidation matches any WebGPU validation error.
	ErrValidation = &WGPUError{Type: ErrorTypeValidation}
	// ErrOutOfMemory matches any WebGPU out-of-memory error.
	ErrOutOfMemory = &WGPUError{Type: ErrorTypeOutOfMemory}
	// ErrInternal matches any WebGPU internal error.
	ErrInternal = &WGPUError{Type: ErrorTypeInternal}
	// ErrDeviceLost matches device lost errors.
	ErrDeviceLost = &WGPUError{Type: ErrorTypeUnknown, Message: "device lost"}
)

// WGPUError represents a WebGPU operation error with context.
// Supports errors.Is() and errors.As() for programmatic handling.
type WGPUError struct {
	// Op is the operation that failed (e.g., "CreateBuffer", "RequestAdapter").
	Op string
	// Type is the WebGPU error type (Validation, OutOfMemory, Internal).
	Type ErrorType
	// Message is the error message from wgpu-native.
	Message string
}

// Error returns a formatted error string including the operation name and message.
func (e *WGPUError) Error() string {
	if e.Op != "" && e.Message != "" {
		return fmt.Sprintf("wgpu: %s: %s", e.Op, e.Message)
	}
	if e.Op != "" {
		return fmt.Sprintf("wgpu: %s failed", e.Op)
	}
	if e.Message != "" {
		return "wgpu: " + e.Message
	}
	return "wgpu: unknown error"
}

// Is supports errors.Is() matching by error Type.
// This allows: errors.Is(err, wgpu.ErrValidation)
func (e *WGPUError) Is(target error) bool {
	t, ok := target.(*WGPUError)
	if !ok {
		return false
	}
	// Match by Type if target has no specific Op/Message
	if t.Op == "" && t.Message == "" {
		return e.Type == t.Type
	}
	return e.Op == t.Op && e.Type == t.Type && e.Message == t.Message
}
