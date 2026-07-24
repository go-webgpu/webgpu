package wgpu

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/go-webgpu/goffi/ffi"
	"github.com/go-webgpu/goffi/types"
)

type float32CallOps struct {
	prepare func(
		*types.CallInterface,
		types.CallingConvention,
		*types.TypeDescriptor,
		[]*types.TypeDescriptor,
	) error
	call func(
		*types.CallInterface,
		unsafe.Pointer,
		unsafe.Pointer,
		[]unsafe.Pointer,
	) (syscall.Errno, error)
}

var nativeFloat32CallOps = float32CallOps{
	prepare: ffi.PrepareCallInterface,
	call:    ffi.CallFunction,
}

// callFloat32 invokes a native function using the platform's scalar
// floating-point return convention.
func callFloat32(
	ops float32CallOps,
	name string,
	convention types.CallingConvention,
	fn unsafe.Pointer,
	args ...uintptr,
) (float32, error) {
	// TODO: cache float-return CIFs once each procedure's call shape is stable.
	argTypes := make([]*types.TypeDescriptor, len(args))
	for i := range argTypes {
		argTypes[i] = types.PointerTypeDescriptor
	}

	var cif types.CallInterface
	if err := ops.prepare(&cif, convention, types.FloatTypeDescriptor, argTypes); err != nil {
		return 0, fmt.Errorf("wgpu: failed to prepare CIF for %s: %w", name, err)
	}

	argPtrs := make([]unsafe.Pointer, len(args))
	for i := range args {
		argPtrs[i] = unsafe.Pointer(&args[i])
	}

	var result float32
	if _, err := ops.call(&cif, fn, unsafe.Pointer(&result), argPtrs); err != nil {
		return 0, fmt.Errorf("wgpu: call to %s failed: %w", name, err)
	}
	return result, nil
}
