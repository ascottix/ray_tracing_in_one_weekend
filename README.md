# Ray Tracing in One Weekend, in Go

This is a Go implementation of the raytracer described in the book [_Ray Tracing in One Weekend_](https://raytracing.github.io/books/RayTracingInOneWeekend.html) by Peter Shirley, Trevor David Black, Steve Hollasch.

It can render all images described in the book, that's why you'll find some duplicated code.

Here's the final image rendered with a generous amount of rays per pixel:

![Final image](https://ascottix.github.io/rtiow/rtiow_final_large.png)

To generate an image run:

> go run . [image_number]

where __image_number__ is a number between 1 and 23.

Output is a file named `out.ppm` in PPM format.

All images are rendered with default parameter values. Different values can only be set by editing the source code.
