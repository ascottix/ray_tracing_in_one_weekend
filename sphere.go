package main

import "math"

type Sphere struct {
	center Point3
	radius float64
	mat    Material // Used from image 13
}

func NewSphere(center Point3, radius float64) Sphere {
	return Sphere{center: center, radius: radius}
}

func NewSphereWithMaterial(center Point3, radius float64, mat Material) Sphere {
	return Sphere{center: center, radius: radius, mat: mat}
}

// Implement the Hittable interface
func (s Sphere) Hit(ray Ray, rayTmin, rayTmax float64, rec *HitRecord) bool {
	oc := ray.Origin().Sub(s.center)
	a := ray.Direction().Dot(ray.Direction())
	half_b := oc.Dot(ray.Direction()) // With respect to the initial versions we have now removed the 2 factor from b and simplified the rest accordingly
	c := oc.Dot(oc) - s.radius*s.radius
	discriminant := half_b*half_b - a*c

	if discriminant < 0 {
		return false // No intersection (any point where t < 0 is behind the camera)
	}

	// Find the nearest root that lies in the allowed range
	sqrtd := math.Sqrt(discriminant)
	root := (-half_b - sqrtd) / a // First intersection (closest to the camera)
	if root <= rayTmin || root >= rayTmax {
		root = (-half_b + sqrtd) / a // Second intersection
		if root <= rayTmin || root >= rayTmax {
			return false
		}
	}

	rec.T = root
	rec.P = ray.At(rec.T)
	outwardNormal := rec.P.Sub(s.center).Div(s.radius) // Divide by the sphere radius as it's cheaper that calling UnitVector() and gets the same result here
	rec.SetFaceNormal(ray, outwardNormal)
	rec.Mat = s.mat

	return true
}
