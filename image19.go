package main

import (
	"io"
	"math"
)

func Image19(w io.Writer) {
	world := NewHittableList()

	R := math.Cos(math.Pi / 4)

	materialLeft := NewLambertianMaterial(NewColor(0, 0, 1))
	materialRight := NewLambertianMaterial(NewColor(1, 0, 0))

	world.Add(NewSphereWithMaterial(NewPoint3(-R, 0, -1), R, materialLeft))
	world.Add(NewSphereWithMaterial(NewPoint3(R, 0, -1), R, materialRight))

	cam := NewPositionableCamera()
	cam.SetVerticalFieldOfView(90)

	cam.Render(w, world, 100, 50)
}
