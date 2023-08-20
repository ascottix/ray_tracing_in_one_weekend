package main

import (
	"fmt"
	"io"
)

func i2_rayColor(ray Ray) Color {
	unitDirection := ray.Direction().UnitVector() // All vector components are now between -1 and +1
	alpha := (unitDirection.Y + 1) / 2            // Alpha is a number between 0 and 1 proportional to the Y component of the ray direction
	return NewColor(1, 1, 1).Mul(1 - alpha).Add(NewColor(0.5, 0.7, 1.0).Mul(alpha))
}

func Image2(w io.Writer) {
	// The camera is centered at (0,0,0) and oriented so that the Y-axis goes up, the X-axis goes right and the negative Z-axis points in the view direction
	cameraCenter := NewPoint3(0, 0, 0)

	// The viewport U and V vectors have the same alignment as the image we want to produce, which has the (0,0) pixel at the top left
	viewport_U := NewVec3(ViewportWidth, 0, 0)   // Vector from left to right edge of viewport
	viewport_V := NewVec3(0, -ViewportHeight, 0) // Vector from top to bottom edge of viewport

	// The pixel delta vectors represent the distance between adjacent pixels in the viewport
	pixelDelta_U := viewport_U.Div(float64(ImageWidth))
	pixelDelta_V := viewport_V.Div(float64(ImageHeight))

	// The viewport is positioned along the negative Z-axis, at the "focal length" distance from the camera
	viewportUpperLeft := cameraCenter.Sub(NewVec3(0, 0, FocalLength)).Sub(viewport_U.Div(2)).Sub(viewport_V.Div(2))

	// Pixels are just in the middle of each viewport grid cell
	pixelUpperLeft := viewportUpperLeft.Add(pixelDelta_U.Mul(0.5)).Add(pixelDelta_V.Mul(0.5))

	// Render
	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", ImageWidth, ImageHeight)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < ImageHeight; y++ {
		for x := 0; x < ImageWidth; x++ {
			pixelCenter := pixelUpperLeft.Add(pixelDelta_U.Mul(float64(x))).Add(pixelDelta_V.Mul(float64(y)))
			direction := pixelCenter.Sub(cameraCenter) // Note: the direction is not normalized
			ray := NewRay(cameraCenter, direction)
			c := i2_rayColor(ray)

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)

}
