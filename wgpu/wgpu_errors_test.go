package wgpu

import (
	"errors"
	"testing"
)

func TestWGPUErrorIs(t *testing.T) {
	err := &WGPUError{Op: "CreateBuffer", Type: ErrorTypeValidation, Message: "invalid size"}

	if !errors.Is(err, ErrValidation) {
		t.Error("expected errors.Is(err, ErrValidation) to be true")
	}
	if errors.Is(err, ErrOutOfMemory) {
		t.Error("expected errors.Is(err, ErrOutOfMemory) to be false")
	}
}

func TestWGPUErrorAs(t *testing.T) {
	var err error = &WGPUError{Op: "RequestAdapter", Message: "no GPU found"}

	var wgpuErr *WGPUError
	if !errors.As(err, &wgpuErr) {
		t.Fatal("expected errors.As to succeed")
	}
	if wgpuErr.Op != "RequestAdapter" {
		t.Errorf("expected Op=RequestAdapter, got %q", wgpuErr.Op)
	}
}

func TestWGPUErrorString(t *testing.T) {
	tests := []struct {
		err  *WGPUError
		want string
	}{
		{&WGPUError{Op: "CreateBuffer", Message: "too large"}, "wgpu: CreateBuffer: too large"},
		{&WGPUError{Op: "CreateBuffer"}, "wgpu: CreateBuffer failed"},
		{&WGPUError{Message: "something"}, "wgpu: something"},
		{&WGPUError{}, "wgpu: unknown error"},
	}
	for _, tt := range tests {
		if got := tt.err.Error(); got != tt.want {
			t.Errorf("Error() = %q, want %q", got, tt.want)
		}
	}
}
