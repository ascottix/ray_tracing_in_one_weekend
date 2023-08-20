package main

import "math"

type Interval struct {
	Min float64
	Max float64
}

func NewInterval(min, max float64) Interval {
	return Interval{Min: min, Max: max}
}

func NewEmptyInterval() Interval {
	return Interval{Min: math.Inf(+1), Max: math.Inf(-1)}
}

func (i Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

func (i Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}

func (i Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}

	if x > i.Max {
		return i.Max
	}

	return x
}

var Empty = NewEmptyInterval()
var Universe = NewInterval(math.Inf(-1), math.Inf(+1))
