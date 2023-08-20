package main

import (
	"io"
	"math"
)

func Image7(w io.Writer) {
	world := NewHittableList()

	world.Add(NewSphere(NewPoint3(0, 0, -1), 0.5))
	world.Add(NewSphere(NewPoint3(0, -100.5, -1), 100))

	cam := NewCamera()

	// To avoid too much refactoring, we use the same code of Image8 but set the maxDepth to math.MaxInt (i.e. practically limitless)
	cam.RenderWithDiffuseMaterial(w, world, 100, math.MaxInt)
}
