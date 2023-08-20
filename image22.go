package main

import (
	"io"
)

func Image22(w io.Writer) {
	world := NewHittableList()

	materialGround := NewLambertianMaterial(NewColor(0.8, 0.8, 0.0))
	materialCenter := NewLambertianMaterial(NewColor(0.1, 0.2, 0.5))
	materialLeft := NewDielectricMaterial(1.5)
	materialRight := NewMetalMaterial(NewColor(0.8, 0.6, 0.2), 0)

	world.Add(NewSphereWithMaterial(NewPoint3(0.0, -100.5, -1), 100, materialGround))
	world.Add(NewSphereWithMaterial(NewPoint3(0, 0, -1), 0.5, materialCenter))
	world.Add(NewSphereWithMaterial(NewPoint3(-1, 0, -1), 0.5, materialLeft))
	world.Add(NewSphereWithMaterial(NewPoint3(-1, 0, -1), -0.4, materialLeft)) // A negative radius creates a hollow sphere
	world.Add(NewSphereWithMaterial(NewPoint3(1, 0, -1), 0.5, materialRight))

	cam := NewPositionableCamera()
	cam.SetLookFrom(NewPoint3(-2, 2, 1))
	cam.SetVerticalFieldOfView(20)
	cam.SetFocusDistance(3.4)
	cam.SetDefocusAngle(10)

	cam.Render(w, world, 100, 50)
}
