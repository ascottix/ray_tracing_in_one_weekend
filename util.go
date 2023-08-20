package main

import (
	"math"
	"math/rand"
)

func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Returns a random number in the interval [0,1)
func RandomDouble() float64 {
	return rand.Float64()
}

// Returns a random number in the interval [min, max)
func RandomDoubleInInterval(min, max float64) float64 {
	return min + (max-min)*RandomDouble()
}

func LinearToGamma(linear float64) float64 {
	return math.Sqrt(linear)
}
