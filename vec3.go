package main

import "math"

// Vec3 is used to represent points, vectors and RGB colors
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func NewVec3(x, y, z float64) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}

func (v Vec3) Negate() Vec3 {
	return Vec3{X: -v.X, Y: -v.Y, Z: -v.Z}
}

func (v Vec3) Add(w Vec3) Vec3 {
	return Vec3{X: v.X + w.X, Y: v.Y + w.Y, Z: v.Z + w.Z}
}

func (v Vec3) Sub(w Vec3) Vec3 {
	return Vec3{X: v.X - w.X, Y: v.Y - w.Y, Z: v.Z - w.Z} // Slight optimization over v.Add(w.Negate())
}

func (v Vec3) Mul(t float64) Vec3 {
	return Vec3{X: v.X * t, Y: v.Y * t, Z: v.Z * t}
}

func (v Vec3) Div(t float64) Vec3 {
	return v.Mul(1 / t)
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) Dot(w Vec3) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

func (v Vec3) Cross(w Vec3) Vec3 {
	return Vec3{
		X: v.Y*w.Z - v.Z*w.Y,
		Y: v.Z*w.X - v.X*w.Z,
		Z: v.X*w.Y - v.Y*w.X}
}

func (v Vec3) MultiplyComponents(w Vec3) Vec3 {
	return Vec3{
		X: v.X * w.X,
		Y: v.Y * w.Y,
		Z: v.Z * w.Z}
}

func (v Vec3) UnitVector() Vec3 {
	return v.Div(v.Length())
}

// These functions create a random vector with various constraints, they are used to simulate diffuse reflection
func NewRandomVec3() Vec3 {
	return NewVec3(RandomDouble(), RandomDouble(), RandomDouble())
}

func NewRandomInIntervalVec3(min, max float64) Vec3 {
	return NewVec3(RandomDoubleInInterval(min, max), RandomDoubleInInterval(min, max), RandomDoubleInInterval(min, max))
}

func NewRandomInUnitSphereVec3() Vec3 {
	for {
		p := NewRandomInIntervalVec3(-1, 1) // Create a random vector inside a cube
		if p.LengthSquared() <= 1 {         // If the length of the vector is less than 1 then the vector is inside a sphere (centered at the origin)
			return p
		}
	}
}

func NewRandomUnitVec3() Vec3 {
	return NewRandomInUnitSphereVec3().UnitVector()
}

func NewRandomUnitInHemisphereVec3(normal Vec3) Vec3 {
	vecOnUnitSphere := NewRandomInUnitSphereVec3().UnitVector()
	if vecOnUnitSphere.Dot(normal) > 0 {
		return vecOnUnitSphere
	} else {
		return vecOnUnitSphere.Negate()
	}
}

// Checks whether the vector is close to zero
func (v Vec3) NearZero() bool {
	s := 1e-8
	return math.Abs(v.X) < s && math.Abs(v.Y) < s && math.Abs(v.Z) < s
}

// Point3 is an alias for Vec3, it is just used for clarity
type Point3 = Vec3

func NewPoint3(x, y, z float64) Point3 {
	return Point3{X: x, Y: y, Z: z}
}

// Color is also an alias for Vec3
type Color = Vec3

func NewColor(x, y, z float64) Color {
	return Color{X: x, Y: y, Z: z}
}
