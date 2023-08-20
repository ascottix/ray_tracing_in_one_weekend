package main

import (
	"fmt"
	"io"
	"math"
)

// To see if a ray hits a sphere we first need to consider the equation that describes all points on a sphere:
// (x - center_x)^2 + (y - center_y)^2 + (z - center_z)^2 = radius^2
// If P is a generic point we can rewrite above using the dot product:
// (P - center) • (P - center) = radius^2
// In our case we are interested in points generated by the ray, which have the form:
// Origin + t*Direction
// Plugging this definition in place of P above and doing the math, all vectors disappear because of dot products
// and we get a standard 2nd degree equation in t.
//
// We return the actual t value of the intersection, or -1 if the ray misses the sphere.
func i4_hitSphere(center Point3, radius float64, ray Ray) float64 {
	oc := ray.Origin().Sub(center)
	a := ray.Direction().Dot(ray.Direction())
	b := 2 * oc.Dot(ray.Direction())
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return -1 // No intersection (any point where t < 0 is behind the camera)
	} else {
		return (-b - math.Sqrt(discriminant)) / (2 * a) // Return the first intersection (the one closest to the camera)
	}
}

func i4_rayColor(ray Ray) Color {
	c := Point3{0, 0, -1} // Sphere center
	t := i4_hitSphere(c, 0.5, ray)
	if t >= 0 {
		// Compute the normal vector, which is a unit vector perpendicular to the sphere surface on the point of intersection, pointing outwards
		n := ray.At(t).Sub(c).UnitVector() // Just take the point on the surface and subtract the sphere center from it, then make it a unit vector
		return Color{n.X + 1, n.Y + 1, n.Z + 1}.Div(2)
	}

	return i2_rayColor(ray) // Reuse gradient background from image 2
}

func Image4(w io.Writer) {
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
			c := i4_rayColor(ray)

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)

}