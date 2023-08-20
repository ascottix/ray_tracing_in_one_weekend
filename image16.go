package main

import (
	"io"
)

func Image16(w io.Writer) {
	world := NewHittableList()

	materialGround := NewLambertianMaterial(NewColor(0.8, 0.8, 0.0))
	materialCenter := NewDielectricAlwaysRefractMaterial(1.5)
	materialLeft := NewDielectricAlwaysRefractMaterial(1.5)
	materialRight := NewMetalMaterial(NewColor(0.8, 0.6, 0.2), 1.0)

	world.Add(NewSphereWithMaterial(NewPoint3(0.0, -100.5, -1), 100, materialGround))
	world.Add(NewSphereWithMaterial(NewPoint3(0, 0, -1), 0.5, materialCenter))
	world.Add(NewSphereWithMaterial(NewPoint3(-1, 0, -1), 0.5, materialLeft))
	world.Add(NewSphereWithMaterial(NewPoint3(1, 0, -1), 0.5, materialRight))

	cam := NewCamera()
	cam.HitRayTmin = 0.001 // Avoid casting rays that start too close to a surface

	cam.RenderWithObjectMaterial(w, world, 100, 50)
}
