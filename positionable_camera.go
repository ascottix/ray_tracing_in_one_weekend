package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

type PositionableCamera struct {
	aspectRatio    float64
	imageWidth     int
	imageHeight    int
	vfov           float64 // Vertical field of view angle in degrees
	lookFrom       Point3  // Where the camera "eye" is positioned
	lookAt         Point3  // Where the camera is looking at
	vUp            Vec3    // Up direction relative to the camera
	focusDistance  float64
	defocusAngle   float64
	pixelDelta_U   Vec3
	pixelDelta_V   Vec3
	pixelUpperLeft Point3
	defocusDisk_U  Vec3
	defocusDisk_V  Vec3
}

func NewPositionableCamera() PositionableCamera {
	return PositionableCamera{imageWidth: ImageWidth, aspectRatio: AspectRatio, vfov: 90, lookFrom: NewPoint3(0, 0, 0), lookAt: NewPoint3(0, 0, -1), vUp: NewVec3(0, 1, 0), focusDistance: 0, defocusAngle: 0}
}

func (camera *PositionableCamera) SetAspectRatio(ratio float64) {
	camera.aspectRatio = ratio
}

func (camera *PositionableCamera) SetDefocusAngle(angle float64) {
	camera.defocusAngle = angle
}

func (camera *PositionableCamera) SetFocusDistance(distance float64) {
	camera.focusDistance = distance
}

func (camera *PositionableCamera) SetImageWidth(width int) {
	camera.imageWidth = width
}

func (camera *PositionableCamera) SetLookAt(p Point3) {
	camera.lookAt = p
}

func (camera *PositionableCamera) SetLookFrom(p Point3) {
	camera.lookFrom = p
}

func (camera *PositionableCamera) SetVerticalFieldOfView(vfov float64) {
	camera.vfov = vfov
}

func (camera *PositionableCamera) Initialize() {
	camera.imageHeight = int(float64(camera.imageWidth) / camera.aspectRatio)

	// Determine the viewport dimentions
	focusDistance := camera.focusDistance
	if focusDistance == 0 { // If focus distance is unassigned, use the distance between the camera center and the "look at" point
		focusDistance = camera.lookAt.Sub(camera.lookFrom).Length()
	}
	theta := DegreesToRadians(camera.vfov)
	h := math.Tan(theta / 2)
	viewportHeight := h * 2 * focusDistance
	viewportWidth := viewportHeight * float64(camera.imageWidth) / float64(camera.imageHeight)

	// Calculate the u,v,w unit basis vectors for the camera coordinate frame.
	w := camera.lookFrom.Sub(camera.lookAt).UnitVector()
	u := camera.vUp.Cross(w).UnitVector()
	v := w.Cross(u)

	// The viewport U and V vectors have the same alignment as the image we want to produce, which has the (0,0) pixel at the top left
	viewport_U := u.Mul(viewportWidth)   // Vector from left to right edge of viewport
	viewport_V := v.Mul(-viewportHeight) // Vector from top to bottom edge of viewport

	// The pixel delta vectors represent the distance between adjacent pixels in the viewport
	camera.pixelDelta_U = viewport_U.Div(float64(camera.imageWidth))
	camera.pixelDelta_V = viewport_V.Div(float64(camera.imageHeight))

	// The viewport is positioned along the negative Z-axis, at the "focal length" distance from the camera, centered with respect to the X and Y axis
	viewportUpperLeft := camera.lookFrom.Sub(w.Mul(focusDistance)).Sub(viewport_U.Div(2)).Sub(viewport_V.Div(2))

	// We want to place pixels in the middle of viewport "grid" cells, so add a half-delta to each coordinate
	camera.pixelUpperLeft = viewportUpperLeft.Add(camera.pixelDelta_U.Mul(0.5)).Add(camera.pixelDelta_V.Mul(0.5))

	// Calculate the camera defocus disk basis vectors
	defocusRadius := camera.focusDistance * math.Tan(DegreesToRadians(camera.defocusAngle/2))
	camera.defocusDisk_U = u.Mul(defocusRadius)
	camera.defocusDisk_V = v.Mul(defocusRadius)
}

// Returns a random point in the square surrounding a pixel at the origin
func (camera PositionableCamera) pixelSampleSquare() Vec3 {
	// Get a random point position, each coordinate is in the [-0.5, 0.5) interval
	// (remember that pixelUpperLeft starts at x=0.5, y=0.5)
	px := -0.5 + RandomDouble()
	py := -0.5 + RandomDouble()

	// Return the vector that leads the ray into the above randomized point of the viewport
	return camera.pixelDelta_U.Mul(px).Add(camera.pixelDelta_V.Mul(py))
}

func (camera PositionableCamera) getRandomPointInDefocusDisk() Point3 {
	// Get a random point in the unit disk
	var x, y float64

	for {
		x = RandomDoubleInInterval(-1, 1)
		y = RandomDoubleInInterval(-1, 1)
		if x*x+y*y <= 1 {
			break
		}
	}

	// Return the corresponding point in the defocus disk
	return camera.lookFrom.Add(camera.defocusDisk_U.Mul(x)).Add(camera.defocusDisk_V.Mul(y))
}

// Get a randomly sampled camera ray for the pixel at location i, j
func (camera PositionableCamera) getRay(i, j int) Ray {
	pixelCenter := camera.pixelUpperLeft.Add(camera.pixelDelta_U.Mul(float64(i))).Add(camera.pixelDelta_V.Mul(float64(j)))
	pixelSample := pixelCenter.Add(camera.pixelSampleSquare())
	origin := camera.lookFrom
	if camera.defocusAngle > 0 {
		origin = camera.getRandomPointInDefocusDisk()
	}
	direction := pixelSample.Sub(origin) // Note: the direction is not normalized

	return NewRay(origin, direction)
}

// The following function uses the properties of the object material to properly compute the ray color
func (camera PositionableCamera) RayColor(ray Ray, world Hittable, depth int) Color {
	rec := HitRecord{}

	if depth <= 0 {
		return Color{0, 0, 0}
	}

	if world.Hit(ray, 0.001, math.Inf(+1), &rec) {
		scattered := Ray{}
		attenuation := Color{}

		if rec.Mat.Scatter(ray, &rec, &attenuation, &scattered) {
			c := camera.RayColor(scattered, world, depth-1)

			return Color{c.X * attenuation.X, c.Y * attenuation.Y, c.Z * attenuation.Z}
		}

		return Color{0, 0, 0}
	}

	return i2_rayColor(ray) // Reuse gradient background from image 2
}

func (camera *PositionableCamera) Render(w io.Writer, world Hittable, samplesPerPixel, maxRayDepth int) {
	camera.Initialize()

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", camera.imageWidth, camera.imageHeight)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < camera.imageHeight; y++ {
		fmt.Fprintf(os.Stderr, "Rendering scanline %d of %d (%d%%)\n", y+1, camera.imageHeight, (y+1)*100/camera.imageHeight)

		for x := 0; x < camera.imageWidth; x++ {
			c := NewColor(0, 0, 0) // Start with black

			// Accumulate all samples into one color, this may bring the color components out of their nominal [0,1] range
			for sample := 0; sample < samplesPerPixel; sample++ {
				ray := camera.getRay(x, y)
				c = c.Add(camera.RayColor(ray, world, maxRayDepth))
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
