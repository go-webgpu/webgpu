package wgpu

import "fmt"

// Sentinel errors for programmatic error handling via errors.Is().
var (
	ErrValidation  = &WGPUError{Type: ErrorTypeValidation}
	ErrOutOfMemory = &WGPUError{Type: ErrorTypeOutOfMemory}
	ErrInternal    = &WGPUError{Type: ErrorTypeInternal}
	ErrDeviceLost  = &WGPUError{Type: ErrorTypeUnknown, Message: "device lost"}
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
