package wgpu

import (
	"testing"
	"unsafe"

	"github.com/gogpu/gputypes"
)

// =============================================================================
// Fuzz tests for FFI boundary — enum conversion functions
// These ensure no panics or invalid memory access for arbitrary input values.
// =============================================================================

func FuzzToWGPUBufferBindingType(f *testing.F) {
	f.Add(uint32(0))
	f.Add(uint32(1))
	f.Add(uint32(3))
	f.Add(uint32(255))
	f.Add(uint32(0xFFFFFFFF))
	f.Fuzz(func(t *testing.T, v uint32) {
		result := toWGPUBufferBindingType(gputypes.BufferBindingType(v))
		if v == 0 && result != 0 {
			t.Errorf("zero input must map to zero, got %d", result)
		}
	})
}

func FuzzToWGPUSamplerBindingType(f *testing.F) {
	f.Add(uint32(0))
	f.Add(uint32(1))
	f.Add(uint32(3))
	f.Add(uint32(0xFFFFFFFF))
	f.Fuzz(func(t *testing.T, v uint32) {
		result := toWGPUSamplerBindingType(gputypes.SamplerBindingType(v))
		if v == 0 && result != 0 {
			t.Errorf("zero input must map to zero, got %d", result)
		}
	})
}

func FuzzToWGPUTextureSampleType(f *testing.F) {
	f.Add(uint32(0))
	f.Add(uint32(1))
	f.Add(uint32(5))
	f.Add(uint32(0xFFFFFFFF))
	f.Fuzz(func(t *testing.T, v uint32) {
		result := toWGPUTextureSampleType(gputypes.TextureSampleType(v))
		if v == 0 && result != 0 {
			t.Errorf("zero input must map to zero, got %d", result)
		}
	})
}

func FuzzToWGPUStorageTextureAccess(f *testing.F) {
	f.Add(uint32(0))
	f.Add(uint32(1))
	f.Add(uint32(3))
	f.Add(uint32(0xFFFFFFFF))
	f.Fuzz(func(t *testing.T, v uint32) {
		result := toWGPUStorageTextureAccess(gputypes.StorageTextureAccess(v))
		if v == 0 && result != 0 {
			t.Errorf("zero input must map to zero, got %d", result)
		}
	})
}

// FuzzTextureFormatDirectCast verifies that direct uint32 cast of TextureFormat
// never panics. gputypes v0.3.0 and wgpu-native v29 share identical format values.
func FuzzTextureFormatDirectCast(f *testing.F) {
	for i := uint32(0); i <= 0x65; i++ {
		f.Add(i)
	}
	f.Add(uint32(0xFFFFFFFF))
	f.Fuzz(func(t *testing.T, v uint32) {
		// Direct cast must be identity: gputypes matches v29 exactly.
		gf := gputypes.TextureFormat(v)
		wire := uint32(gf)
		if wire != v {
			t.Errorf("direct cast changed value: input %d, got %d", v, wire)
		}
		// Round-trip from wire back to gputypes must also be identity.
		back := gputypes.TextureFormat(wire)
		if back != gf {
			t.Errorf("round-trip failed for format %d", v)
		}
	})
}

func FuzzToWGPUVertexStepMode(f *testing.F) {
	f.Add(uint32(0))
	f.Add(uint32(1))
	f.Add(uint32(2))
	f.Add(uint32(3))
	f.Add(uint32(0xFFFFFFFF))
	f.Fuzz(func(t *testing.T, v uint32) {
		_ = toWGPUVertexStepMode(gputypes.VertexStepMode(v))
	})
}

func FuzzToWGPUVertexFormat(f *testing.F) {
	for i := uint32(0); i <= 40; i++ {
		f.Add(i)
	}
	f.Add(uint32(0xFFFFFFFF))
	f.Fuzz(func(t *testing.T, v uint32) {
		_ = toWGPUVertexFormat(gputypes.VertexFormat(v))
	})
}

// =============================================================================
// Fuzz tests for struct layout correctness
// Verify struct sizes match C ABI expectations.
// =============================================================================

func TestFFIStructSizes(t *testing.T) {
	ptrSize := unsafe.Sizeof(uintptr(0))

	tests := []struct {
		name     string
		got      uintptr
		expected uintptr
	}{
		// StringView: uintptr (Data) + uintptr (Length)
		{"StringView", unsafe.Sizeof(StringView{}), 2 * ptrSize},
		// Future: uint64
		{"Future", unsafe.Sizeof(Future{}), 8},
		// DrawIndirectArgs: 4 x uint32
		{"DrawIndirectArgs", unsafe.Sizeof(DrawIndirectArgs{}), 16},
		// DrawIndexedIndirectArgs: 5 x uint32/int32
		{"DrawIndexedIndirectArgs", unsafe.Sizeof(DrawIndexedIndirectArgs{}), 20},
		// DispatchIndirectArgs: 3 x uint32
		{"DispatchIndirectArgs", unsafe.Sizeof(DispatchIndirectArgs{}), 12},
		// Color: 4 x float64
		{"Color", unsafe.Sizeof(Color{}), 32},
		// ChainedStruct: uintptr + uint32 (+ padding)
		{"ChainedStruct", unsafe.Sizeof(ChainedStruct{}), ptrSize + ptrSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s size = %d, want %d", tt.name, tt.got, tt.expected)
			}
		})
	}
}

// =============================================================================
// Fuzz tests for LeakReport.String() — must not panic
// =============================================================================

func FuzzLeakReportString(f *testing.F) {
	f.Add(0, "Buffer", 3)
	f.Add(1, "Texture", 1)
	f.Add(100, "", 0)
	f.Fuzz(func(t *testing.T, count int, typeName string, typeCount int) {
		report := &LeakReport{
			Count: count,
			Types: map[string]int{typeName: typeCount},
		}
		// Must not panic
		_ = report.String()
	})
}

// =============================================================================
// Fuzz tests for WGPUError.Error() — must not panic
// =============================================================================

func FuzzWGPUErrorString(f *testing.F) {
	f.Add("CreateBuffer", uint32(2), "invalid size")
	f.Add("", uint32(0), "")
	f.Add("Op", uint32(5), "msg")
	f.Fuzz(func(t *testing.T, op string, errType uint32, msg string) {
		e := &WGPUError{
			Op:      op,
			Type:    ErrorType(errType),
			Message: msg,
		}
		// Must not panic
		s := e.Error()
		if s == "" {
			t.Error("Error() must never return empty string")
		}
	})
}

// =============================================================================
// Fuzz tests for WGPUError.Is() — must not panic
// =============================================================================

func FuzzWGPUErrorIs(f *testing.F) {
	f.Add("Op", uint32(2), "msg", uint32(2))
	f.Add("", uint32(0), "", uint32(0))
	f.Fuzz(func(t *testing.T, op string, errType uint32, msg string, targetType uint32) {
		e := &WGPUError{Op: op, Type: ErrorType(errType), Message: msg}
		target := &WGPUError{Type: ErrorType(targetType)}
		// Must not panic
		_ = e.Is(target)
	})
}

// =============================================================================
// Math helpers — no panics on edge values
// =============================================================================

func FuzzMat4Perspective(f *testing.F) {
	f.Add(float64(1.0), float64(1.0), float64(0.1), float64(100.0))
	f.Add(float64(0.0), float64(0.0), float64(0.0), float64(0.0))
	f.Add(float64(-1.0), float64(1e10), float64(1e-10), float64(1e20))
	f.Fuzz(func(t *testing.T, fovY, aspect, near, far float64) {
		// Must not panic (even with degenerate values)
		_ = Mat4Perspective(float32(fovY), float32(aspect), float32(near), float32(far))
	})
}

func FuzzVec3Normalize(f *testing.F) {
	f.Add(float64(1.0), float64(0.0), float64(0.0))
	f.Add(float64(0.0), float64(0.0), float64(0.0))
	f.Add(float64(1e30), float64(1e30), float64(1e30))
	f.Fuzz(func(t *testing.T, x, y, z float64) {
		v := Vec3{float32(x), float32(y), float32(z)}
		// Must not panic
		_ = v.Normalize()
	})
}
