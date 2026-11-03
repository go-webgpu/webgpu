package wgpu

// pinSink and pinSinkOn force the escape: the compiler cannot prove that
// pinSinkOn is always false, so anything passed through pin is heap-allocated.
var (
	pinSink   any
	pinSinkOn bool
)

// pin returns p unchanged but forces the pointed-to object onto the heap.
//
// //go:uintptrescapes on Proc.Call only covers pointers converted to uintptr
// directly in the call's argument list. Descriptors that wgpu-native receives
// nested (a pointer stored as uintptr inside another struct, or converted
// before the call) are not covered: if they live on the goroutine stack, a
// stack growth or shrink between the conversion and the native call leaves
// the uintptr pointing at stale memory. The Go heap does not move objects, so
// the pointer stays valid.
func pin[T any](p *T) *T {
	if pinSinkOn {
		pinSink = p
	}
	return p
}
