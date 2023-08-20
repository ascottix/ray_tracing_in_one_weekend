package main

import (
	"io"
)

// This image is my attempt to simulate the bug mentioned by the book author
func Image15(w io.Writer) {
	world := NewHittableList()

	materialGround := NewLambertianMaterial(NewColor(0.8, 0.8, 0.0))
	materialCenter := NewBuggyDielectricMaterial(1)
	materialLeft := NewBuggyDielectricMaterial(1)
	materialRight := NewMetalMaterial(NewColor(0.8, 0.6, 0.2), 0)

	world.Add(NewSphereWithMaterial(NewPoint3(0.0, -100.5, -1), 100, materialGround))
	world.Add(NewSphereWithMaterial(NewPoint3(0, 0, -1), 0.5, materialCenter))
	world.Add(NewSphereWithMaterial(NewPoint3(-1, 0, -1), 0.5, materialLeft))
	world.Add(NewSphereWithMaterial(NewPoint3(1, 0, -1), 0.5, materialRight))

	cam := NewCamera()
	cam.HitRayTmin = 0.001 // Avoid casting rays that start too close to a surface

	cam.RenderWithObjectMaterial(w, world, 100, 10)
}
