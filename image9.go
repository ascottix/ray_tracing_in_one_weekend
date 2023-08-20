package main

import (
	"io"
)

func Image9(w io.Writer) {
	world := NewHittableList()

	world.Add(NewSphere(NewPoint3(0, 0, -1), 0.5))
	world.Add(NewSphere(NewPoint3(0, -100.5, -1), 100))

	cam := NewCamera()
	cam.HitRayTmin = 0.001 // Avoid casting rays that start too close to a surface

	cam.RenderWithDiffuseMaterial(w, world, 100, 50)
}
