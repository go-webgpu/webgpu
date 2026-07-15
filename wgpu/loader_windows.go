//go:build windows

package wgpu

import (
	"syscall"
	"unsafe"

	"github.com/go-webgpu/goffi/ffi"
	"github.com/go-webgpu/goffi/types"
)

// windowsLibrary wraps syscall.LazyDLL to implement the Library interface.
type windowsLibrary struct {
	dll *syscall.LazyDLL
}

// windowsProc wraps syscall.LazyProc to implement the Proc interface.
type windowsProc struct {
	proc *syscall.LazyProc
}

// loadLibrary loads a DLL using Windows syscall.NewLazyDLL.
// Returns a Library interface that can be used to get procedures.
// The DLL is eagerly loaded to report errors immediately.
func loadLibrary(name string) (Library, error) {
	dll := syscall.NewLazyDLL(name)
	// Force eager load to detect missing DLL immediately.
	if err := dll.Load(); err != nil {
		return nil, err
	}
	return &windowsLibrary{dll: dll}, nil
}

// NewProc retrieves a procedure from the Windows DLL.
func (w *windowsLibrary) NewProc(name string) Proc {
	return &windowsProc{
		proc: w.dll.NewProc(name),
	}
}

// Call invokes the Windows procedure with the given arguments.
// This directly delegates to syscall.LazyProc.Call().
func (w *windowsProc) Call(args ...uintptr) (uintptr, uintptr, error) {
	return w.proc.Call(args...)
}

// CallFloat32 invokes a float32-returning procedure through goffi so the
// Windows x64 ABI reads XMM0. syscall.LazyProc.Call only exposes integer
// return registers and therefore cannot safely call this signature.
func (w *windowsProc) CallFloat32(args ...uintptr) (float32, error) {
	if err := w.proc.Find(); err != nil {
		return 0, err
	}

	argTypes := make([]*types.TypeDescriptor, len(args))
	for i := range argTypes {
		argTypes[i] = types.PointerTypeDescriptor
	}
	var cif types.CallInterface
	if err := ffi.PrepareCallInterface(
		&cif,
		types.WindowsCallingConvention,
		types.FloatTypeDescriptor,
		argTypes,
	); err != nil {
		return 0, err
	}

	argPtrs := make([]unsafe.Pointer, len(args))
	for i := range args {
		argPtrs[i] = unsafe.Pointer(&args[i])
	}
	var result float32
	_, err := ffi.CallFunction(&cif, unsafe.Pointer(w.proc.Addr()), unsafe.Pointer(&result), argPtrs)
	return result, err
}
