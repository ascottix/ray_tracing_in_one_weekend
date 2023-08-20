package main

type HitRecord struct {
	P         Point3
	Normal    Vec3
	T         float64
	FrontFace bool
	Mat       Material // Used starting from image 13
}

type Hittable interface {
	Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool
}

// A normal to an object surface may point outwards or inwards... how do we choose?
// There are two main conventions:
// 1. the normal always points outwards
// 2. the normal always points against the ray
// In the first case we can use the dot product between ray and normal to determine whether
// the ray is outside the object (dot product is negative) or inside the object (dot product is positive).
// In the second case the dot product will always be negative so we need to determine that information
// first, and store it for later.
// This book takes the second approach.
// Note: outwardNormal is assumed to have unit length
func (h *HitRecord) SetFaceNormal(ray Ray, outwardNormal Vec3) {
	h.FrontFace = ray.Direction().Dot(outwardNormal) < 0
	if h.FrontFace {
		// Ray is outside the sphere
		h.Normal = outwardNormal
	} else {
		// Ray is inside the sphere, flip the normal
		h.Normal = outwardNormal.Negate()
	}
}
