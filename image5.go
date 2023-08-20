package main

import (
	"fmt"
	"io"
	"math"
)

func i5_rayColor(ray Ray, world HittableList) Color {
	rec := HitRecord{}

	if world.Hit(ray, 0, math.Inf(+1), &rec) {
		return rec.Normal.Add(NewVec3(1, 1, 1)).Div(2)
	}

	return i2_rayColor(ray) // Reuse gradient background from image 2
}

func Image5(w io.Writer) {
	// The camera is centered at (0,0,0) and oriented so that the Y-axis goes up, the X-axis goes right and the negative Z-axis points in the view direction
	cameraCenter := NewPoint3(0, 0, 0)

	// The viewport U and V vectors have the same alignment as the image we want to produce, which has the (0,0) pixel at the top left
	viewport_U := NewVec3(ViewportWidth, 0, 0)   // Vector from left to right edge of viewport
	viewport_V := NewVec3(0, -ViewportHeight, 0) // Vector from top to bottom edge of viewport

	// The pixel delta vectors represent the distance between adjacent pixels in the viewport
	pixelDelta_U := viewport_U.Div(float64(ImageWidth))
	pixelDelta_V := viewport_V.Div(float64(ImageHeight))

	// The viewport is positioned along the negative Z-axis, at the "focal length" distance from the camera, centered with respect to the X and Y axis
	viewportUpperLeft := cameraCenter.Sub(NewVec3(0, 0, FocalLength)).Sub(viewport_U.Div(2)).Sub(viewport_V.Div(2))

	// We want to place pixels in the middle of viewport "grid" cells, so add a half-delta to each coordinate
	pixelUpperLeft := viewportUpperLeft.Add(pixelDelta_U.Mul(0.5)).Add(pixelDelta_V.Mul(0.5))

	// World
	world := NewHittableList()

	world.Add(NewSphere(NewPoint3(0, 0, -1), 0.5))
	world.Add(NewSphere(NewPoint3(0, -100.5, -1), 100))

	// Render
	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", ImageWidth, ImageHeight)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < ImageHeight; y++ {
		for x := 0; x < ImageWidth; x++ {
			pixelCenter := pixelUpperLeft.Add(pixelDelta_U.Mul(float64(x))).Add(pixelDelta_V.Mul(float64(y)))
			direction := pixelCenter.Sub(cameraCenter) // Note: the direction is not normalized
			ray := NewRay(cameraCenter, direction)
			c := i5_rayColor(ray, world)

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}
