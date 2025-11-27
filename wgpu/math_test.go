package wgpu

import (
	"math"
	"testing"
)

// epsilon for floating point comparisons
const epsilon = 1e-6

// almostEqual checks if two float32 values are approximately equal
func almostEqual(a, b float32) bool {
	return math.Abs(float64(a-b)) < epsilon
}

// mat4AlmostEqual checks if two matrices are approximately equal
func mat4AlmostEqual(a, b Mat4) bool {
	for i := 0; i < 16; i++ {
		if !almostEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// vec3AlmostEqual checks if two Vec3 are approximately equal
func vec3AlmostEqual(a, b Vec3) bool {
	return almostEqual(a.X, b.X) && almostEqual(a.Y, b.Y) && almostEqual(a.Z, b.Z)
}

// vec4AlmostEqual checks if two Vec4 are approximately equal
func vec4AlmostEqual(a, b Vec4) bool {
	return almostEqual(a.X, b.X) && almostEqual(a.Y, b.Y) &&
		almostEqual(a.Z, b.Z) && almostEqual(a.W, b.W)
}

func TestMat4Identity(t *testing.T) {
	identity := Mat4Identity()

	expected := Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	if !mat4AlmostEqual(identity, expected) {
		t.Errorf("Mat4Identity() = %v, want %v", identity, expected)
	}
}

func TestMat4Translate(t *testing.T) {
	translate := Mat4Translate(10, 20, 30)

	expected := Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		10, 20, 30, 1,
	}

	if !mat4AlmostEqual(translate, expected) {
		t.Errorf("Mat4Translate(10, 20, 30) = %v, want %v", translate, expected)
	}
}

func TestMat4Scale(t *testing.T) {
	scale := Mat4Scale(2, 3, 4)

	expected := Mat4{
		2, 0, 0, 0,
		0, 3, 0, 0,
		0, 0, 4, 0,
		0, 0, 0, 1,
	}

	if !mat4AlmostEqual(scale, expected) {
		t.Errorf("Mat4Scale(2, 3, 4) = %v, want %v", scale, expected)
	}
}

func TestMat4RotateX(t *testing.T) {
	// Rotate 90 degrees around X
	rot := Mat4RotateX(math.Pi / 2)

	// Expected: Y becomes Z, Z becomes -Y
	expected := Mat4{
		1, 0, 0, 0,
		0, 0, 1, 0,
		0, -1, 0, 0,
		0, 0, 0, 1,
	}

	if !mat4AlmostEqual(rot, expected) {
		t.Errorf("Mat4RotateX(π/2) = %v, want %v", rot, expected)
	}
}

func TestMat4RotateY(t *testing.T) {
	// Rotate 90 degrees around Y
	rot := Mat4RotateY(math.Pi / 2)

	// Expected: Z becomes X, X becomes -Z
	expected := Mat4{
		0, 0, -1, 0,
		0, 1, 0, 0,
		1, 0, 0, 0,
		0, 0, 0, 1,
	}

	if !mat4AlmostEqual(rot, expected) {
		t.Errorf("Mat4RotateY(π/2) = %v, want %v", rot, expected)
	}
}

func TestMat4RotateZ(t *testing.T) {
	// Rotate 90 degrees around Z
	rot := Mat4RotateZ(math.Pi / 2)

	// Expected: X becomes Y, Y becomes -X
	expected := Mat4{
		0, 1, 0, 0,
		-1, 0, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	if !mat4AlmostEqual(rot, expected) {
		t.Errorf("Mat4RotateZ(π/2) = %v, want %v", rot, expected)
	}
}

func TestMat4Perspective(t *testing.T) {
	// FOV 45 degrees, aspect 16:9, near 0.1, far 100
	fov := float32(math.Pi / 4)
	aspect := float32(16.0 / 9.0)
	near := float32(0.1)
	far := float32(100.0)

	persp := Mat4Perspective(fov, aspect, near, far)

	// Verify basic properties:
	// - Non-zero diagonal elements for X and Y scaling
	// - Perspective division in column 2
	if persp[0] == 0 || persp[5] == 0 {
		t.Errorf("Mat4Perspective() has zero scaling: %v", persp)
	}

	if persp[11] != -1 {
		t.Errorf("Mat4Perspective()[11] = %v, want -1 (perspective division)", persp[11])
	}
}

func TestMat4LookAt(t *testing.T) {
	eye := Vec3{0, 0, 5}
	center := Vec3{0, 0, 0}
	up := Vec3{0, 1, 0}

	view := Mat4LookAt(eye, center, up)

	// Transform eye position - should result in origin
	eyePos := view.MulVec4(Vec4{eye.X, eye.Y, eye.Z, 1})

	// Due to view matrix, eye should be at origin in view space
	expected := Vec4{0, 0, 0, 1}

	if !vec4AlmostEqual(eyePos, expected) {
		t.Errorf("Mat4LookAt() eye transform = %v, want %v", eyePos, expected)
	}
}

func TestMat4Mul(t *testing.T) {
	// Test: Translation * Scale = Combined transform
	translate := Mat4Translate(10, 0, 0)
	scale := Mat4Scale(2, 2, 2)

	combined := translate.Mul(scale)

	// Apply to point (1, 0, 0)
	point := Vec4{1, 0, 0, 1}
	result := combined.MulVec4(point)

	// Expected: scale first (2, 0, 0), then translate (12, 0, 0)
	expected := Vec4{12, 0, 0, 1}

	if !vec4AlmostEqual(result, expected) {
		t.Errorf("Mat4.Mul() transform result = %v, want %v", result, expected)
	}
}

func TestMat4MulIdentity(t *testing.T) {
	identity := Mat4Identity()
	translate := Mat4Translate(5, 10, 15)

	// Identity * M = M
	result := identity.Mul(translate)

	if !mat4AlmostEqual(result, translate) {
		t.Errorf("Identity.Mul(M) = %v, want %v", result, translate)
	}

	// M * Identity = M
	result = translate.Mul(identity)

	if !mat4AlmostEqual(result, translate) {
		t.Errorf("M.Mul(Identity) = %v, want %v", result, translate)
	}
}

func TestMat4MulVec4(t *testing.T) {
	// Test translation
	translate := Mat4Translate(10, 20, 30)
	point := Vec4{1, 2, 3, 1}

	result := translate.MulVec4(point)
	expected := Vec4{11, 22, 33, 1}

	if !vec4AlmostEqual(result, expected) {
		t.Errorf("Translate.MulVec4() = %v, want %v", result, expected)
	}
}

func TestVec3Sub(t *testing.T) {
	a := Vec3{10, 20, 30}
	b := Vec3{1, 2, 3}

	result := a.Sub(b)
	expected := Vec3{9, 18, 27}

	if !vec3AlmostEqual(result, expected) {
		t.Errorf("Vec3.Sub() = %v, want %v", result, expected)
	}
}

func TestVec3Cross(t *testing.T) {
	// Standard basis vectors
	x := Vec3{1, 0, 0}
	y := Vec3{0, 1, 0}

	// X × Y = Z
	result := x.Cross(y)
	expected := Vec3{0, 0, 1}

	if !vec3AlmostEqual(result, expected) {
		t.Errorf("Vec3.Cross(X, Y) = %v, want Z %v", result, expected)
	}

	// Y × X = -Z (anti-commutative)
	result = y.Cross(x)
	expected = Vec3{0, 0, -1}

	if !vec3AlmostEqual(result, expected) {
		t.Errorf("Vec3.Cross(Y, X) = %v, want -Z %v", result, expected)
	}
}

func TestVec3Normalize(t *testing.T) {
	v := Vec3{3, 4, 0}
	result := v.Normalize()

	// Length should be 1
	length := float32(math.Sqrt(float64(result.X*result.X + result.Y*result.Y + result.Z*result.Z)))

	if !almostEqual(length, 1.0) {
		t.Errorf("Vec3.Normalize() length = %v, want 1.0", length)
	}

	// Direction should be preserved
	expected := Vec3{0.6, 0.8, 0}
	if !vec3AlmostEqual(result, expected) {
		t.Errorf("Vec3.Normalize() = %v, want %v", result, expected)
	}
}

func TestVec3NormalizeZero(t *testing.T) {
	v := Vec3{0, 0, 0}
	result := v.Normalize()

	expected := Vec3{0, 0, 0}
	if !vec3AlmostEqual(result, expected) {
		t.Errorf("Vec3.Normalize(zero) = %v, want %v", result, expected)
	}
}

func TestVec3Dot(t *testing.T) {
	a := Vec3{1, 2, 3}
	b := Vec3{4, 5, 6}

	result := a.Dot(b)
	expected := float32(1*4 + 2*5 + 3*6) // 32

	if !almostEqual(result, expected) {
		t.Errorf("Vec3.Dot() = %v, want %v", result, expected)
	}
}

func TestVec3DotOrthogonal(t *testing.T) {
	// Orthogonal vectors have dot product = 0
	x := Vec3{1, 0, 0}
	y := Vec3{0, 1, 0}

	result := x.Dot(y)

	if !almostEqual(result, 0) {
		t.Errorf("Vec3.Dot(orthogonal) = %v, want 0", result)
	}
}

// Benchmark tests

func BenchmarkMat4Identity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Mat4Identity()
	}
}

func BenchmarkMat4Mul(b *testing.B) {
	m1 := Mat4Translate(10, 20, 30)
	m2 := Mat4Scale(2, 3, 4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m1.Mul(m2)
	}
}

func BenchmarkMat4MulVec4(b *testing.B) {
	m := Mat4Translate(10, 20, 30)
	v := Vec4{1, 2, 3, 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.MulVec4(v)
	}
}

func BenchmarkVec3Normalize(b *testing.B) {
	v := Vec3{3, 4, 5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Normalize()
	}
}

func BenchmarkMat4Perspective(b *testing.B) {
	fov := float32(math.Pi / 4)
	aspect := float32(16.0 / 9.0)
	near := float32(0.1)
	far := float32(100.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Mat4Perspective(fov, aspect, near, far)
	}
}

func BenchmarkMat4LookAt(b *testing.B) {
	eye := Vec3{0, 0, 5}
	center := Vec3{0, 0, 0}
	up := Vec3{0, 1, 0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Mat4LookAt(eye, center, up)
	}
}
