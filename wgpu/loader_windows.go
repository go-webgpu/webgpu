//go:build windows

package wgpu

import (
	"syscall"
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
func loadLibrary(name string) Library {
	return &windowsLibrary{
		dll: syscall.NewLazyDLL(name),
	}
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
