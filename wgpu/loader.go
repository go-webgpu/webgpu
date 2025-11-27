// Package wgpu provides cross-platform library loading abstractions.
package wgpu

// Library represents a dynamically loaded library (DLL/SO/DYLIB).
// Platform-specific implementations handle the actual loading mechanism.
type Library interface {
	// NewProc retrieves a procedure (function) from the library.
	// Returns a Proc interface that can be used to call the function.
	NewProc(name string) Proc
}

// Proc represents a procedure (function pointer) from a dynamically loaded library.
// It abstracts platform-specific function calling mechanisms.
type Proc interface {
	// Call invokes the procedure with the given arguments.
	// Returns the result value and error (if any).
	// Arguments are passed as uintptr to match C ABI.
	Call(args ...uintptr) (uintptr, uintptr, error)
}
