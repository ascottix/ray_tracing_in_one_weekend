package main

import (
	"fmt"
	"io"
	"math"
)

const (
	// Viewport
	AspectRatio    = 16.0 / 9.0
	ImageWidth     = 400
	ImageHeight    = int(ImageWidth / AspectRatio)
	ViewportHeight = 2.0
	ViewportWidth  = ViewportHeight * float64(ImageWidth) / float64(ImageHeight)

	// Camera
	FocalLength = 1.0 // Distance between the camera center and the viewport
)

type Camera struct {
	aspectRatio    float64
	imageWidth     int
	imageHeight    int
	center         Point3
	pixelDelta_U   Vec3
	pixelDelta_V   Vec3
	pixelUpperLeft Point3
	// Public members
	HitRayTmin float64
}

func NewCamera() Camera {
	return Camera{imageWidth: ImageWidth, aspectRatio: AspectRatio, HitRayTmin: 0}
}

func (camera *Camera) Initialize() {
	camera.center = NewPoint3(0, 0, 0)

	camera.imageHeight = int(float64(camera.imageWidth) / camera.aspectRatio)

	viewportWidth := ViewportHeight * float64(camera.imageWidth) / float64(camera.imageHeight)

	// The viewport U and V vectors have the same alignment as the image we want to produce, which has the (0,0) pixel at the top left
	viewport_U := NewVec3(viewportWidth, 0, 0)   // Vector from left to right edge of viewport
	viewport_V := NewVec3(0, -ViewportHeight, 0) // Vector from top to bottom edge of viewport

	// The pixel delta vectors represent the distance between adjacent pixels in the viewport
	camera.pixelDelta_U = viewport_U.Div(float64(camera.imageWidth))
	camera.pixelDelta_V = viewport_V.Div(float64(camera.imageHeight))

	// The viewport is positioned along the negative Z-axis, at the "focal length" distance from the camera, centered with respect to the X and Y axis
	viewportUpperLeft := camera.center.Sub(NewVec3(0, 0, FocalLength)).Sub(viewport_U.Div(2)).Sub(viewport_V.Div(2))

	// We want to place pixels in the middle of viewport "grid" cells, so add a half-delta to each coordinate
	camera.pixelUpperLeft = viewportUpperLeft.Add(camera.pixelDelta_U.Mul(0.5)).Add(camera.pixelDelta_V.Mul(0.5))
}

func (camera Camera) RayColor(ray Ray, world Hittable) Color {
	rec := HitRecord{}

	if world.Hit(ray, 0, math.Inf(+1), &rec) {
		return rec.Normal.Add(NewVec3(1, 1, 1)).Div(2)
	}

	return i2_rayColor(ray) // Reuse gradient background from image 2
}

func (camera *Camera) Render(w io.Writer, world Hittable) {
	camera.Initialize()

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", camera.imageWidth, camera.imageHeight)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < camera.imageHeight; y++ {
		for x := 0; x < camera.imageWidth; x++ {
			pixelCenter := camera.pixelUpperLeft.Add(camera.pixelDelta_U.Mul(float64(x))).Add(camera.pixelDelta_V.Mul(float64(y)))
			direction := pixelCenter.Sub(camera.center) // Note: the direction is not normalized
			ray := NewRay(camera.center, direction)
			c := camera.RayColor(ray, world)

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}

// Returns a random point in the square surrounding a pixel at the origin
func (camera Camera) pixelSampleSquare() Vec3 {
	// Get a random point position, each coordinate is in the [-0.5, 0.5) interval
	// (remember that pixelUpperLeft starts at x=0.5, y=0.5)
	px := -0.5 + RandomDouble()
	py := -0.5 + RandomDouble()

	// Return the vector that leads the ray into the above randomized point of the viewport
	return camera.pixelDelta_U.Mul(px).Add(camera.pixelDelta_V.Mul(py))
}

// Get a randomly sampled camera ray for the pixel at location i, j
func (camera Camera) getRay(i, j int) Ray {
	pixelCenter := camera.pixelUpperLeft.Add(camera.pixelDelta_U.Mul(float64(i))).Add(camera.pixelDelta_V.Mul(float64(j)))
	pixelSample := pixelCenter.Add(camera.pixelSampleSquare())
	direction := pixelSample.Sub(camera.center) // Note: the direction is not normalized

	return NewRay(camera.center, direction)
}

// Casts multiple rays per pixel in order to get a higher quality, antialiased image
func (camera *Camera) RenderWithMultipleSamples(w io.Writer, world Hittable, samplesPerPixel int) {
	camera.Initialize()

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", camera.imageWidth, camera.imageHeight)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < camera.imageHeight; y++ {
		for x := 0; x < camera.imageWidth; x++ {
			c := NewColor(0, 0, 0) // Start with black

			// Accumulate all samples into one color, this may bring the color components out of their nominal [0,1] range
			for sample := 0; sample < samplesPerPixel; sample++ {
				ray := camera.getRay(x, y)
				c = c.Add(camera.RayColor(ray, world))
			}

			c = c.Div(float64(samplesPerPixel)) // Bring the color components back to the [0,1] range

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}

// Simulation of a diffuse (matte) material
func (camera Camera) RayColorOfDiffuseMaterial(ray Ray, world Hittable, depth int) Color {
	rec := HitRecord{}

	if depth <= 0 {
		return Color{0, 0, 0}
	}

	if world.Hit(ray, camera.HitRayTmin, math.Inf(+1), &rec) {
		// Here we want to simulate the behavior of a matte material, i.e. diffuse reflection. It works like this:
		// 1. the direction of the reflected ray is random, so we generate a random vector on the surface hemisphere
		// 2. we create a new ray that goes from the point of surface intersection towards the random direction above
		// 2. we get the color of this reflected ray
		// 3. we arbitrarily weight that color by 50%
		direction := NewRandomUnitInHemisphereVec3(rec.Normal)
		return camera.RayColorOfDiffuseMaterial(NewRay(rec.P, direction), world, depth-1).Mul(0.5)
	}

	return i2_rayColor(ray) // Reuse gradient background from image 2
}

func (camera *Camera) RenderWithDiffuseMaterial(w io.Writer, world Hittable, samplesPerPixel, maxRayDepth int) {
	camera.Initialize()

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", camera.imageWidth, camera.imageHeight)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < camera.imageHeight; y++ {
		for x := 0; x < camera.imageWidth; x++ {
			c := NewColor(0, 0, 0) // Start with black

			// Accumulate all samples into one color, this may bring the color components out of their nominal [0,1] range
			for sample := 0; sample < samplesPerPixel; sample++ {
				ray := camera.getRay(x, y)
				c = c.Add(camera.RayColorOfDiffuseMaterial(ray, world, maxRayDepth))
			}

			c = c.Div(float64(samplesPerPixel)) // Bring the color components back to the [0,1] range

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}

// Simulation of a diffuse (matte) material with the Lambert model
func (camera Camera) RayColorOfLambertianMaterial(ray Ray, world Hittable, depth int, reflectance float64) Color {
	rec := HitRecord{}

	if depth <= 0 {
		return Color{0, 0, 0}
	}

	if world.Hit(ray, camera.HitRayTmin, math.Inf(+1), &rec) {
		// Here we want to simulate the behavior of a matte material, i.e. diffuse reflection, with the Lambert model.
		// In practice, as explained in chapter 9.4, we want to increase the probability that the reflected ray is
		// close to the surface normal.
		// We can do this by simply adding to the surface normal a random vector that goes from the point (rec.P + rec.Normal)
		// to a random point on a unit sphere centered in (rec.P + rec.Normal). I.e. we add a random unit sphere to the normal.
		direction := rec.Normal.Add(NewRandomUnitVec3())
		return camera.RayColorOfDiffuseMaterial(NewRay(rec.P, direction), world, depth-1).Mul(reflectance)
	}

	return i2_rayColor(ray) // Reuse gradient background from image 2
}

func (camera *Camera) RenderWithLambertianMaterial(w io.Writer, world Hittable, samplesPerPixel, maxRayDepth int) {
	camera.Initialize()

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", camera.imageWidth, camera.imageHeight)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < camera.imageHeight; y++ {
		for x := 0; x < camera.imageWidth; x++ {
			c := NewColor(0, 0, 0) // Start with black

			// Accumulate all samples into one color, this may bring the color components out of their nominal [0,1] range
			for sample := 0; sample < samplesPerPixel; sample++ {
				ray := camera.getRay(x, y)
				c = c.Add(camera.RayColorOfLambertianMaterial(ray, world, maxRayDepth, 0.5))
			}

			c = c.Div(float64(samplesPerPixel)) // Bring the color components back to the [0,1] range

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}

func (camera *Camera) RenderGamut(w io.Writer, world Hittable, samplesPerPixel, maxRayDepth int, gammaCorrection bool) {
	camera.Initialize()

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", camera.imageWidth, camera.imageHeight)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < camera.imageHeight; y++ {
		for x := 0; x < camera.imageWidth; x++ {
			c := NewColor(0, 0, 0) // Start with black

			// Accumulate all samples into one color, this may bring the color components out of their nominal [0,1] range
			for sample := 0; sample < samplesPerPixel; sample++ {
				ray := camera.getRay(x, y)
				// Test image is divided into 5 vertical bands, with reflectance going from 10% to 90%
				band := x * 5 / camera.imageWidth // Same as x / (camera.imageWidth / 5)
				reflectance := 0.1 + 0.2*float64(band)
				c = c.Add(camera.RayColorOfLambertianMaterial(ray, world, maxRayDepth, reflectance))
			}

			c = c.Div(float64(samplesPerPixel)) // Bring the color components back to the [0,1] range

			if gammaCorrection {
				c.X = LinearToGamma(c.X)
				c.Y = LinearToGamma(c.Y)
				c.Z = LinearToGamma(c.Z)
			}

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}

// The following function uses the properties of the object material to properly compute the ray color
func (camera Camera) RayColorOfObjectMaterial(ray Ray, world Hittable, depth int) Color {
	rec := HitRecord{}

	if depth <= 0 {
		return Color{0, 0, 0}
	}

	if world.Hit(ray, camera.HitRayTmin, math.Inf(+1), &rec) {
		scattered := Ray{}
		attenuation := Color{}

		if rec.Mat.Scatter(ray, &rec, &attenuation, &scattered) {
			c := camera.RayColorOfObjectMaterial(scattered, world, depth-1)

			return Color{c.X * attenuation.X, c.Y * attenuation.Y, c.Z * attenuation.Z}
		}

		return Color{0, 0, 0}
	}

	return i2_rayColor(ray) // Reuse gradient background from image 2
}

func (camera *Camera) RenderWithObjectMaterial(w io.Writer, world Hittable, samplesPerPixel, maxRayDepth int) {
	camera.Initialize()

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", camera.imageWidth, camera.imageHeight)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < camera.imageHeight; y++ {
		for x := 0; x < camera.imageWidth; x++ {
			c := NewColor(0, 0, 0) // Start with black

			// Accumulate all samples into one color, this may bring the color components out of their nominal [0,1] range
			for sample := 0; sample < samplesPerPixel; sample++ {
				ray := camera.getRay(x, y)
				c = c.Add(camera.RayColorOfObjectMaterial(ray, world, maxRayDepth))
			}

			c = c.Div(float64(samplesPerPixel)) // Bring the color components back to the [0,1] range

			// Apply gamma correction
			c.X = LinearToGamma(c.X)
			c.Y = LinearToGamma(c.Y)
			c.Z = LinearToGamma(c.Z)

			ir := int(255.999 * c.X)
			ig := int(255.999 * c.Y)
			ib := int(255.999 * c.Z)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}
