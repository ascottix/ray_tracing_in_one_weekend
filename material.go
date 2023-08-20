package main

import "math"

type Material interface {
	// Returns true if the surface scattered (reflected) the incoming ray, or false if it has absorbed it.
	// If the ray has been scattered, also returns the scattered ray and the attenuation color (which depends on the material).
	Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool
}

type LambertianMaterial struct {
	albedo Color
}

type MetalMaterial struct {
	albedo Color
	fuzz   float64 // If fuzz is 0 the material is perfectly smooth, setting 0 < fuzz <= 1 adds roughness to the surface
}

// Lambertian material
func NewLambertianMaterial(a Color) LambertianMaterial {
	return LambertianMaterial{albedo: a}
}

func (m LambertianMaterial) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	scatterDirection := rec.Normal.Add(NewRandomUnitVec3())

	// Catch an edge case where the random unit vector is exactly opposite to the surface normal and nullifies the scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}

	*scattered = NewRay(rec.P, scatterDirection)
	*attenuation = m.albedo

	return true
}

// Metal material
func NewMetalMaterial(a Color, fuzz float64) MetalMaterial {
	return MetalMaterial{albedo: a, fuzz: math.Min(fuzz, 1)}
}

func Reflect(v, n Vec3) Vec3 {
	return v.Add(n.Mul(-2 * n.Dot(v)))
}

func Refract(uv, n Vec3, etaiOverEtat float64) Vec3 {
	cosTheta := n.Dot(uv.Negate()) // It's math.Min(n.Dot(uv.Negate()), 1.0) in the original source
	rOutPerp := uv.Add(n.Mul(cosTheta)).Mul(etaiOverEtat)
	rOutParallel := n.Mul(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

func (m MetalMaterial) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	reflected := Reflect(ray.Direction().UnitVector(), rec.Normal)

	*scattered = NewRay(rec.P, reflected.Add(NewRandomUnitVec3().Mul(m.fuzz)))
	*attenuation = m.albedo

	// We should just return true here, but because of the fuzziness it may happen that a ray is scattered below the surface.
	// If that happens, just pretend the surface has absorbed it and don't scatter.
	return rec.Normal.Dot(scattered.Direction()) > 0
}

// Buggy dielectric material that always refracts
type BuggyDielectricMaterial struct {
	ir float64
}

func NewBuggyDielectricMaterial(indexOfRefraction float64) BuggyDielectricMaterial {
	return BuggyDielectricMaterial{ir: indexOfRefraction}
}

func (m BuggyDielectricMaterial) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	refractionRatio := m.ir
	if rec.FrontFace {
		refractionRatio = 1 / refractionRatio
	}

	unitDirection := ray.Direction().UnitVector()

	// Inline the Refract() function so we can try to bug it!
	cosTheta := rec.Normal.Dot(unitDirection.Negate())
	rOutPerp := unitDirection.Add(rec.Normal.Mul(cosTheta)).Mul(refractionRatio)
	rOutParallel := rec.Normal.Mul(math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared()))) // Missing a minus sign before math.Sqrt
	refracted := rOutPerp.Add(rOutParallel)

	*attenuation = Color{1, 1, 1}
	*scattered = NewRay(rec.P, refracted)

	// This adds the "weird black stuff" to the image, not sure where the original bug could have come from
	if cosTheta < 0.3 {
		*attenuation = Color{0, 0, 0}
	}

	return true
}

// Dielectric material that always refracts
type DielectricAlwaysRefractMaterial struct {
	ir float64
}

func NewDielectricAlwaysRefractMaterial(indexOfRefraction float64) DielectricAlwaysRefractMaterial {
	return DielectricAlwaysRefractMaterial{ir: indexOfRefraction}
}

func (m DielectricAlwaysRefractMaterial) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	refractionRatio := m.ir
	if rec.FrontFace {
		refractionRatio = 1 / refractionRatio
	}

	unitDirection := ray.Direction().UnitVector()
	refracted := Refract(unitDirection, rec.Normal, refractionRatio)

	*attenuation = Color{1, 1, 1}
	*scattered = NewRay(rec.P, refracted)

	return true
}

// Dielectric material with correct behavior
type DielectricMaterial struct {
	ir             float64
	useReflectance bool
}

func NewDielectricMaterial(indexOfRefraction float64) DielectricMaterial {
	return DielectricMaterial{ir: indexOfRefraction, useReflectance: true}
}

// Schlick's approximation for reflectance,
// it is used to handle total internal reflection
func SchlickReflectance(cosine, refIdx float64) float64 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

func (m *DielectricMaterial) DisableReflectance() {
	m.useReflectance = false
}

func (m DielectricMaterial) Scatter(ray Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
	refractionRatio := m.ir
	if rec.FrontFace {
		refractionRatio = 1 / refractionRatio
	}

	unitDirection := ray.Direction().UnitVector()

	cosTheta := rec.Normal.Dot(unitDirection.Negate())
	sinTheta := math.Sqrt(1 - cosTheta*cosTheta)

	cannotRefract := (refractionRatio*sinTheta > 1) || (m.useReflectance && SchlickReflectance(cosTheta, refractionRatio) >= RandomDouble())

	if cannotRefract {
		reflected := Reflect(unitDirection, rec.Normal)
		*scattered = NewRay(rec.P, reflected)
	} else {
		refracted := Refract(unitDirection, rec.Normal, refractionRatio)
		*scattered = NewRay(rec.P, refracted)
	}

	*attenuation = Color{1, 1, 1}

	return true
}
