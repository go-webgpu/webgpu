// math.go contains 3D math helpers optimized for WebGPU/WGSL compatibility.

package wgpu

import "math"

// Mat4 represents a 4x4 matrix in column-major order, compatible with WGSL mat4x4<f32>.
// Layout: [col0, col1, col2, col3] where each column is [x, y, z, w].
// Element at column c, row r is at index c*4+r.
// This matches WebGPU/WGSL/OpenGL convention (column-major).
type Mat4 [16]float32

// Vec3 represents a 3D vector with X, Y, Z components.
type Vec3 struct {
	X, Y, Z float32
}

// Vec4 represents a 4D vector with X, Y, Z, W components.
// Compatible with WGSL vec4<f32>.
type Vec4 struct {
	X, Y, Z, W float32
}

// Mat4Identity returns a 4x4 identity matrix.
// The identity matrix has 1s on the diagonal and 0s elsewhere.
func Mat4Identity() Mat4 {
	return Mat4{
		1, 0, 0, 0, // column 0
		0, 1, 0, 0, // column 1
		0, 0, 1, 0, // column 2
		0, 0, 0, 1, // column 3
	}
}

// Mat4Translate returns a translation matrix for the given offset.
// Translates points by (x, y, z) in 3D space.
func Mat4Translate(x, y, z float32) Mat4 {
	return Mat4{
		1, 0, 0, 0, // column 0
		0, 1, 0, 0, // column 1
		0, 0, 1, 0, // column 2
		x, y, z, 1, // column 3 (translation)
	}
}

// Mat4Scale returns a scaling matrix for the given factors.
// Scales objects by (x, y, z) along each axis.
func Mat4Scale(x, y, z float32) Mat4 {
	return Mat4{
		x, 0, 0, 0, // column 0
		0, y, 0, 0, // column 1
		0, 0, z, 0, // column 2
		0, 0, 0, 1, // column 3
	}
}

// Mat4RotateX returns a rotation matrix around the X axis.
// Angle is in radians. Positive rotation follows right-hand rule.
func Mat4RotateX(radians float32) Mat4 {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))

	return Mat4{
		1, 0, 0, 0, // column 0
		0, c, s, 0, // column 1
		0, -s, c, 0, // column 2
		0, 0, 0, 1, // column 3
	}
}

// Mat4RotateY returns a rotation matrix around the Y axis.
// Angle is in radians. Positive rotation follows right-hand rule.
func Mat4RotateY(radians float32) Mat4 {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))

	return Mat4{
		c, 0, -s, 0, // column 0
		0, 1, 0, 0, // column 1
		s, 0, c, 0, // column 2
		0, 0, 0, 1, // column 3
	}
}

// Mat4RotateZ returns a rotation matrix around the Z axis.
// Angle is in radians. Positive rotation follows right-hand rule.
func Mat4RotateZ(radians float32) Mat4 {
	c := float32(math.Cos(float64(radians)))
	s := float32(math.Sin(float64(radians)))

	return Mat4{
		c, s, 0, 0, // column 0
		-s, c, 0, 0, // column 1
		0, 0, 1, 0, // column 2
		0, 0, 0, 1, // column 3
	}
}

// Mat4Perspective returns a perspective projection matrix.
// fovY: vertical field of view in radians
// aspect: aspect ratio (width/height)
// near: near clipping plane distance (must be > 0)
// far: far clipping plane distance (must be > near)
//
// This uses right-handed coordinate system with Z in [-1, 1] (OpenGL/Vulkan style).
// For WebGPU with Z in [0, 1], post-multiply with depth range transform.
func Mat4Perspective(fovY, aspect, near, far float32) Mat4 {
	tanHalfFovy := float32(math.Tan(float64(fovY) / 2.0))
	f := 1.0 / tanHalfFovy

	return Mat4{
		f / aspect, 0, 0, 0, // column 0
		0, f, 0, 0, // column 1
		0, 0, -(far + near) / (far - near), -1, // column 2
		0, 0, -(2 * far * near) / (far - near), 0, // column 3
	}
}

// Mat4LookAt returns a view matrix that looks from eye position towards center.
// eye: camera position
// center: point the camera is looking at
// up: up direction vector (typically (0, 1, 0))
//
// This creates a right-handed coordinate system view matrix.
func Mat4LookAt(eye, center, up Vec3) Mat4 {
	// Forward direction (z axis)
	f := center.Sub(eye).Normalize()

	// Right direction (x axis)
	s := f.Cross(up).Normalize()

	// Recalculated up direction (y axis)
	u := s.Cross(f)

	// Build view matrix (rotation + translation)
	return Mat4{
		s.X, u.X, -f.X, 0, // column 0
		s.Y, u.Y, -f.Y, 0, // column 1
		s.Z, u.Z, -f.Z, 0, // column 2
		-s.Dot(eye), -u.Dot(eye), f.Dot(eye), 1, // column 3
	}
}

// Mul multiplies this matrix by another matrix (column-major order).
// Returns result = m * other (apply m first, then other).
func (m Mat4) Mul(other Mat4) Mat4 {
	var result Mat4

	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			sum := float32(0)
			for k := 0; k < 4; k++ {
				// m[k][row] * other[col][k]
				sum += m[k*4+row] * other[col*4+k]
			}
			result[col*4+row] = sum
		}
	}

	return result
}

// MulVec4 multiplies this matrix by a 4D vector.
// Returns result = m * v (transforms vector by matrix).
func (m Mat4) MulVec4(v Vec4) Vec4 {
	return Vec4{
		X: m[0]*v.X + m[4]*v.Y + m[8]*v.Z + m[12]*v.W,
		Y: m[1]*v.X + m[5]*v.Y + m[9]*v.Z + m[13]*v.W,
		Z: m[2]*v.X + m[6]*v.Y + m[10]*v.Z + m[14]*v.W,
		W: m[3]*v.X + m[7]*v.Y + m[11]*v.Z + m[15]*v.W,
	}
}

// Sub subtracts another vector from this vector.
// Returns v - other.
func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

// Cross computes the cross product of this vector with another.
// Returns v × other (perpendicular to both vectors).
// Result follows right-hand rule.
func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Normalize returns a unit vector in the same direction as v.
// If v has zero length, returns zero vector.
func (v Vec3) Normalize() Vec3 {
	length := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
	if length == 0 {
		return Vec3{0, 0, 0}
	}
	invLength := 1.0 / length
	return Vec3{
		X: v.X * invLength,
		Y: v.Y * invLength,
		Z: v.Z * invLength,
	}
}

// Dot computes the dot product of this vector with another.
// Returns v · other (scalar projection).
func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}
