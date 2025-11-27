//go:build linux || darwin

package wgpu

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/go-webgpu/goffi/ffi"
	"github.com/go-webgpu/goffi/types"
)

// unixLibrary wraps goffi library handle to implement the Library interface.
type unixLibrary struct {
	handle unsafe.Pointer
	name   string
}

// unixProc wraps a goffi function pointer and prepared CIF.
type unixProc struct {
	lib      *unixLibrary
	name     string
	fnPtr    unsafe.Pointer
	cif      types.CallInterface
	cifMu    sync.Mutex
	prepared bool
}

// loadLibrary loads a shared library using goffi.LoadLibrary.
// Returns a Library interface that can be used to get procedures.
func loadLibrary(name string) Library {
	handle, err := ffi.LoadLibrary(name)
	if err != nil {
		// Return a library that will fail on any NewProc call
		// This maintains compatibility with Windows LazyDLL behavior
		return &unixLibrary{
			handle: nil,
			name:   name,
		}
	}

	return &unixLibrary{
		handle: handle,
		name:   name,
	}
}

// NewProc retrieves a procedure from the Unix shared library.
func (u *unixLibrary) NewProc(name string) Proc {
	if u.handle == nil {
		// Return a proc that will fail on Call
		return &unixProc{
			lib:      u,
			name:     name,
			fnPtr:    nil,
			prepared: false,
		}
	}

	fnPtr, err := ffi.GetSymbol(u.handle, name)
	if err != nil {
		// Return a proc that will fail on Call
		return &unixProc{
			lib:      u,
			name:     name,
			fnPtr:    nil,
			prepared: false,
		}
	}

	return &unixProc{
		lib:      u,
		name:     name,
		fnPtr:    fnPtr,
		prepared: false,
	}
}

// Call invokes the Unix procedure with the given arguments.
// This uses goffi's CallFunction with lazy CIF preparation.
//
// Note: WebGPU functions have varying signatures, so we use a conservative
// approach: prepare CIF on first call with actual argument count.
// Most WebGPU functions return uintptr (handles) or void.
func (u *unixProc) Call(args ...uintptr) (uintptr, uintptr, error) {
	if u.fnPtr == nil {
		return 0, 0, fmt.Errorf("wgpu: failed to get symbol %s from %s", u.name, u.lib.name)
	}

	// Lazy CIF preparation on first call
	u.cifMu.Lock()
	if !u.prepared {
		argCount := len(args)
		argTypes := make([]*types.TypeDescriptor, argCount)
		for i := 0; i < argCount; i++ {
			argTypes[i] = types.PointerTypeDescriptor // Conservative: treat all args as uintptr
		}

		// Use platform-specific calling convention
		// Linux/macOS use System V AMD64 ABI (UnixCallingConvention)
		err := ffi.PrepareCallInterface(
			&u.cif,
			types.UnixCallingConvention,
			types.PointerTypeDescriptor, // Most WebGPU functions return uintptr handle
			argTypes,
		)
		if err != nil {
			u.cifMu.Unlock()
			return 0, 0, fmt.Errorf("wgpu: failed to prepare CIF for %s: %w", u.name, err)
		}
		u.prepared = true
	}
	u.cifMu.Unlock()

	// Prepare argument pointers
	argPtrs := make([]unsafe.Pointer, len(args))
	for i := range args {
		argPtrs[i] = unsafe.Pointer(&args[i])
	}

	// Call the function
	var result uintptr
	err := ffi.CallFunction(&u.cif, u.fnPtr, unsafe.Pointer(&result), argPtrs)
	if err != nil {
		return 0, 0, fmt.Errorf("wgpu: call to %s failed: %w", u.name, err)
	}

	// Return result, 0 (no secondary return), nil error
	// This matches Windows syscall.LazyProc.Call signature
	return result, 0, nil
}
