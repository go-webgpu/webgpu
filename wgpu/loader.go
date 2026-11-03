// loader.go provides cross-platform library loading abstractions.

package wgpu

// Library represents a dynamically loaded library (DLL/SO/DYLIB).
// Platform-specific implementations handle the actual loading mechanism.
type Library interface {
	// NewProc retrieves a procedure (function) from the library.
	// Returns a Proc that can be used to call the function.
	NewProc(name string) *Proc
}

// procImpl is the platform-specific calling mechanism behind a Proc.
type procImpl interface {
	// Call invokes the procedure with the given arguments.
	// Returns the result value and error (if any).
	// Arguments are passed as uintptr to match C ABI.
	Call(args ...uintptr) (uintptr, uintptr, error)
}

// Proc represents a procedure (function pointer) from a dynamically loaded
// library. It abstracts platform-specific function calling mechanisms.
//
// Proc is a concrete type (not an interface) so that Call can carry
// //go:uintptrescapes. Callers pass Go structs as
// uintptr(unsafe.Pointer(pin(&local))); without the directive those locals stay
// on the goroutine stack and can move (stack growth, GC shrink) between the
// conversion and the native call, so wgpu-native reads stale input or
// writes its output to the old stack. Seen as wgpuSurfaceGetCurrentTexture
// "returning" a zeroed WGPUSurfaceTexture (status 0) while the texture was
// really acquired, which then aborted the process on the next present or
// configure. The directive is only honored on direct calls, hence no
// interface.
type Proc struct {
	impl procImpl
}

func newProc(impl procImpl) *Proc { return &Proc{impl: impl} }

// Call invokes the procedure with the given arguments. Pointer arguments
// converted to uintptr in the call expression are moved to the heap and kept
// alive for the duration of the call.
//
//go:uintptrescapes
func (p *Proc) Call(args ...uintptr) (uintptr, uintptr, error) {
	return p.impl.Call(args...)
}

// float32Proc is implemented by platform loaders for procedures whose native
// return type is float32. Proc.Call intentionally keeps the existing integer
// return contract for the rest of the WebGPU API; this narrow interface lets
// those procedures use the platform's floating-point return ABI instead.
type float32Proc interface {
	CallFloat32(args ...uintptr) (float32, error)
}
