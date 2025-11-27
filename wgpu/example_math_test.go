package wgpu_test

import (
	"fmt"

	"github.com/go-webgpu/webgpu/wgpu"
)

// ExampleMat4_basic demonstrates basic matrix operations
func ExampleMat4_basic() {
	// Create transformation matrices
	translate := wgpu.Mat4Translate(10, 0, 0)
	scale := wgpu.Mat4Scale(2, 2, 2)

	// Combine transformations (order matters: scale first, then translate)
	combined := translate.Mul(scale)

	// Transform a point
	point := wgpu.Vec4{X: 1, Y: 0, Z: 0, W: 1}
	result := combined.MulVec4(point)

	fmt.Printf("Transformed point: (%.0f, %.0f, %.0f)\n", result.X, result.Y, result.Z)
	// Output: Transformed point: (12, 0, 0)
}

// ExampleMat4Perspective demonstrates perspective projection setup
func ExampleMat4Perspective() {
	// Common 3D camera setup
	fov := float32(45.0 * 3.14159 / 180.0) // 45 degrees in radians
	aspect := float32(16.0 / 9.0)          // 16:9 aspect ratio
	near := float32(0.1)
	far := float32(100.0)

	projection := wgpu.Mat4Perspective(fov, aspect, near, far)

	// Create view matrix (camera looking from (0, 0, 5) towards origin)
	eye := wgpu.Vec3{X: 0, Y: 0, Z: 5}
	center := wgpu.Vec3{X: 0, Y: 0, Z: 0}
	up := wgpu.Vec3{X: 0, Y: 1, Z: 0}
	view := wgpu.Mat4LookAt(eye, center, up)

	// Combine projection and view
	viewProj := projection.Mul(view)

	fmt.Printf("View-Projection matrix ready: %v elements\n", len(viewProj))
	// Output: View-Projection matrix ready: 16 elements
}

// ExampleVec3_cross demonstrates cross product for normal calculation
func ExampleVec3_cross() {
	// Calculate normal for a triangle (right-hand rule)
	// Triangle vertices: v0, v1, v2
	v0 := wgpu.Vec3{X: 0, Y: 0, Z: 0}
	v1 := wgpu.Vec3{X: 1, Y: 0, Z: 0}
	v2 := wgpu.Vec3{X: 0, Y: 1, Z: 0}

	// Edge vectors
	edge1 := v1.Sub(v0)
	edge2 := v2.Sub(v0)

	// Normal = edge1 × edge2
	normal := edge1.Cross(edge2).Normalize()

	fmt.Printf("Triangle normal: (%.0f, %.0f, %.0f)\n", normal.X, normal.Y, normal.Z)
	// Output: Triangle normal: (0, 0, 1)
}

// ExampleMat4_rotation demonstrates rotation matrices
func ExampleMat4_rotation() {
	// Rotate 90 degrees around Y axis
	rotation := wgpu.Mat4RotateY(3.14159 / 2) // 90 degrees in radians

	// Apply rotation to a point on X axis
	point := wgpu.Vec4{X: 1, Y: 0, Z: 0, W: 1}
	rotated := rotation.MulVec4(point)

	// Right-hand rule: rotating +X around +Y by 90° gives -Z direction
	fmt.Printf("Rotated point: X≈%.0f, Z≈%.0f\n", rotated.X, rotated.Z)
	// Output: Rotated point: X≈0, Z≈-1
}
