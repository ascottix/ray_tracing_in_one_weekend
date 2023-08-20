package main

import (
	"fmt"
	"io"
)

func Image1(w io.Writer) {
	width := 256
	height := 256

	fmt.Fprintf(w, "P3\n") // Magic
	fmt.Fprintf(w, "%d %d\n", width, height)
	fmt.Fprintf(w, "255\n") // Maximum value of a color component

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := float64(x) / float64(width-1)
			g := float64(y) / float64(height-1)
			b := 0.0

			ir := int(255.999 * r)
			ig := int(255.999 * g)
			ib := int(255.999 * b)

			fmt.Fprintf(w, "%d %d %d\n", ir, ig, ib)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
}
