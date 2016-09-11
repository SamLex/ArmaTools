package unitcapturereduce

import "math"

type finalFloat64 struct {
	value float64
}

func newFinalFloat64(value float64) finalFloat64 {
	return finalFloat64{
		value: value,
	}
}

func (ff finalFloat64) Value() float64 {
	return ff.value
}

type finalInt struct {
	value int
}

func newFinalInt(value int) finalInt {
	return finalInt{
		value: value,
	}
}

func (fi finalInt) Value() int {
	return fi.value
}

// Linearly interpolate between two values
// Algorithm from https://en.wikipedia.org/wiki/Linear_interpolation
func lerp(v0, v1, t float64) float64 {
	return (1-t)*v0 + t*v1
}

// 'Vectorised' form of the probability density function for the normal distribution
func vectorisedNormalDistributionPDF(x []float64, mean []float64, stddev float64) float64 {
	diff := subArray(x, mean)
	power := -(dotArray(diff, diff) / (2 * stddev * stddev))
	norm := 1 / math.Sqrt(2*stddev*stddev*math.Pi)

	return norm * math.Pow(math.E, power)
}

func subArray(a []float64, b []float64) []float64 {
	sub := make([]float64, len(a))

	for i := range a {
		sub[i] = a[i] - b[i]
	}

	return sub
}

func dotArray(a []float64, b []float64) float64 {
	multiple := mulArray(a, b)

	return sumArray(multiple)
}

func mulArray(a []float64, b []float64) []float64 {
	mul := make([]float64, len(a))

	for i := range a {
		mul[i] = a[i] * b[i]
	}

	return mul
}

func sumArray(a []float64) float64 {
	sum := 0.0

	for _, f := range a {
		sum += f
	}

	return sum
}
