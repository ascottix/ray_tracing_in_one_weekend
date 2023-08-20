package main

import (
	"io"
)

// Most of the code from Image5() is now moved into the Camera class
func Image6(w io.Writer) {
	world := NewHittableList()

	world.Add(NewSphere(NewPoint3(0, 0, -1), 0.5))
	world.Add(NewSphere(NewPoint3(0, -100.5, -1), 100))

	cam := NewCamera()

	cam.RenderWithMultipleSamples(w, world, 100)
}
