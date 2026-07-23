package wgpu

import (
	"errors"
	"syscall"
	"testing"
	"unsafe"

	"github.com/go-webgpu/goffi/types"
)

func TestABIFloat32Call(t *testing.T) {
	const argument = uintptr(0x1234)
	functionToken := byte(1)
	function := unsafe.Pointer(&functionToken)

	t.Run("success", func(t *testing.T) {
		ops := float32CallOps{
			prepare: func(
				_ *types.CallInterface,
				convention types.CallingConvention,
				returnType *types.TypeDescriptor,
				argTypes []*types.TypeDescriptor,
			) error {
				if convention != types.UnixCallingConvention {
					t.Fatalf("calling convention = %v, want Unix", convention)
				}
				if returnType != types.FloatTypeDescriptor {
					t.Fatalf("return type = %v, want float32", returnType)
				}
				if len(argTypes) != 1 || argTypes[0] != types.PointerTypeDescriptor {
					t.Fatalf("argument types = %v, want one pointer", argTypes)
				}
				return nil
			},
			call: func(
				_ *types.CallInterface,
				gotFunction unsafe.Pointer,
				result unsafe.Pointer,
				args []unsafe.Pointer,
			) (syscall.Errno, error) {
				if gotFunction != function {
					t.Fatalf("function = %p, want %p", gotFunction, function)
				}
				if len(args) != 1 || *(*uintptr)(args[0]) != argument {
					t.Fatalf("arguments do not preserve %#x", argument)
				}
				*(*float32)(result) = 0.125
				return 0, nil
			},
		}

		got, err := callFloat32(ops, "testFloat32", types.UnixCallingConvention, function, argument)
		if err != nil {
			t.Fatal(err)
		}
		if got != 0.125 {
			t.Fatalf("result = %v, want 0.125", got)
		}
	})

	t.Run("prepare error", func(t *testing.T) {
		wantErr := errors.New("prepare failed")
		ops := float32CallOps{
			prepare: func(
				*types.CallInterface,
				types.CallingConvention,
				*types.TypeDescriptor,
				[]*types.TypeDescriptor,
			) error {
				return wantErr
			},
		}

		if _, err := callFloat32(ops, "testFloat32", types.UnixCallingConvention, function, argument); !errors.Is(err, wantErr) {
			t.Fatalf("error = %v, want wrapped %v", err, wantErr)
		}
	})

	t.Run("call error", func(t *testing.T) {
		wantErr := errors.New("call failed")
		ops := float32CallOps{
			prepare: func(
				*types.CallInterface,
				types.CallingConvention,
				*types.TypeDescriptor,
				[]*types.TypeDescriptor,
			) error {
				return nil
			},
			call: func(
				*types.CallInterface,
				unsafe.Pointer,
				unsafe.Pointer,
				[]unsafe.Pointer,
			) (syscall.Errno, error) {
				return 0, wantErr
			},
		}

		if _, err := callFloat32(ops, "testFloat32", types.UnixCallingConvention, function, argument); !errors.Is(err, wantErr) {
			t.Fatalf("error = %v, want wrapped %v", err, wantErr)
		}
	})
}
