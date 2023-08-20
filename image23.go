package main

import (
	"io"
)

func Image23(w io.Writer) {
	world := NewHittableList()

	materialGround := NewLambertianMaterial(NewColor(0.5, 0.5, 0.5))
	world.Add(NewSphereWithMaterial(NewPoint3(0.0, -1000, 0), 1000, materialGround))

	ref := NewPoint3(4, 0.2, 0)
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := NewPoint3(float64(a)+0.9*RandomDouble(), 0.2, float64(b)+0.9*RandomDouble())

			if center.Sub(ref).Length() > 0.9 {
				chooseMat := RandomDouble()
				if chooseMat < 0.8 {
					// Diffuse
					albedo := NewRandomVec3().MultiplyComponents(NewRandomVec3())
					mat := NewLambertianMaterial(albedo)
					world.Add(NewSphereWithMaterial(center, 0.2, mat))
				} else if chooseMat < 0.95 {
					// Metal
					albedo := NewRandomInIntervalVec3(0.5, 1)
					fuzz := RandomDoubleInInterval(0, 0.5)
					mat := NewMetalMaterial(albedo, fuzz)
					world.Add(NewSphereWithMaterial(center, 0.2, mat))
				} else {
					// Glass
					mat := NewDielectricMaterial(1.5)
					world.Add(NewSphereWithMaterial(center, 0.2, mat))
				}
			}
		}
	}

	material1 := NewDielectricMaterial(1.5)
	world.Add(NewSphereWithMaterial(NewPoint3(0, 1, 0), 1, material1))
	material2 := NewLambertianMaterial(NewColor(0.4, 0.2, 0.1))
	world.Add(NewSphereWithMaterial(NewPoint3(-4, 1, 0), 1, material2))
	material3 := NewMetalMaterial(NewColor(0.7, 0.6, 0.5), 0)
	world.Add(NewSphereWithMaterial(NewPoint3(4, 1, 0), 1, material3))

	cam := NewPositionableCamera()
	cam.SetLookFrom(NewPoint3(13, 2, 3))
	cam.SetLookAt(NewPoint3(0, 0, 0))
	cam.SetVerticalFieldOfView(20)
	cam.SetFocusDistance(10)
	cam.SetDefocusAngle(0.6)

	samplesPerPixel := 50
	maxRayDepth := 10

	// Renders the final image in the book at full resolution and quality but takes a _very_ long time!
	finalImage := false

	if finalImage {
		samplesPerPixel = 500
		maxRayDepth = 50
		cam.SetImageWidth(1200)
	}

	cam.Render(w, world, samplesPerPixel, maxRayDepth)
}
