package main

// A ray (i.e. a line) is defined by a point (its origin) and a vector (its direction): ray(t) = origin + t*direction
//
// The t parameter allows access to every point on the ray. Usually the term "ray" implies that t > 0, which generates a half-line.
type Ray struct {
	orig Point3
	dir  Vec3
}

func NewRay(origin Point3, direction Vec3) Ray {
	return Ray{orig: origin, dir: direction}
}

func (r Ray) At(t float64) Point3 {
	return r.orig.Add(r.dir.Mul(t))
}

func (r Ray) Origin() Point3 {
	return r.orig
}

func (r Ray) Direction() Vec3 {
	return r.dir
}
