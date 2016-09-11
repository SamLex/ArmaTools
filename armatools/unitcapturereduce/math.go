package unitcapturereduce

import "math"

// Linearly interpolate between two values
// Algorithm from https://en.wikipedia.org/wiki/Linear_interpolation
func lerp(v0, v1, t float64) float64 {
	return (1-t)*v0 + t*v1
}

// Percentage difference
// Algorithm from https://en.wikipedia.org/wiki/Relative_change_and_difference
func percentDifference(theoretical, experimental float64) float64 {
	if theoretical == experimental {
		return 0
	}

	if theoretical == 0 {
		temp := experimental
		experimental = theoretical
		theoretical = temp
	}

	premultiple := math.Abs(experimental-theoretical) / math.Abs(theoretical)
	return premultiple * 100
}
