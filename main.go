package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const OutputFilename = "out.ppm"

type Renderer func(w io.Writer)

func main() {
	renderers := []Renderer{Image1, Image2, Image3, Image4, Image5, Image6, Image7, Image8, Image9, Image10, Image11, Image12, Image13, Image14, Image15, Image16, Image17, Image18, Image19, Image20, Image21, Image22, Image23}

	imageNo := 23

	if len(os.Args) == 2 {
		imageNo, _ = strconv.Atoi(os.Args[1])
	} else {
		fmt.Fprintln(os.Stderr, "No image number specified, default is", imageNo)
	}

	if imageNo < 1 || imageNo > 23 {
		fmt.Fprintln(os.Stderr, "Image number must be between 1 and 23")
	} else {
		renderer := renderers[imageNo-1]

		fmt.Fprintln(os.Stderr, "Rendering image no.", imageNo, "on file", OutputFilename)

		start := time.Now()

		f, err := os.Create(OutputFilename)

		if err != nil {
			panic(err)
		}

		defer f.Close()

		renderer(f)

		elapsed := time.Since(start)

		fmt.Fprintln(os.Stderr, "Done in", elapsed)
	}
}
