/*
Copyright 2016 Euan James Hunter

math.go: Math utilities for armatools/unitcapturereduce

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package unitcapturereduce

import "math"

// Linearly interpolate between two values
// Algorithm from https://en.wikipedia.org/wiki/Linear_interpolation
func lerp(v0, v1, t float64) float64 {
	return (1-t)*v0 + t*v1
}

// 'Vectorised' form of the probability density function for the normal distribution
// Uses the dot product for the squaring
func vectorisedNormalDistributionPDF(x []float64, mean []float64, stddev float64) float64 {
	diff := subArray(x, mean)
	power := -(dotArray(diff, diff) / (2 * stddev * stddev))
	norm := 1 / math.Sqrt(2*stddev*stddev*math.Pi)

	return norm * math.Pow(math.E, power)
}

// Subtract the contents of two float64 slices
// Both slices are assumed to be the same size
func subArray(a []float64, b []float64) []float64 {
	sub := make([]float64, len(a))

	for i := range a {
		sub[i] = a[i] - b[i]
	}

	return sub
}

// Calculate dot product the two float64 slices (treating them as vectors)
// Both slices are assumed to be the same size
func dotArray(a []float64, b []float64) float64 {
	multiple := mulArray(a, b)

	return sumArray(multiple)
}

// Multiple the contents of two float64 slices
// Both slices are assumed to be the same size
func mulArray(a []float64, b []float64) []float64 {
	mul := make([]float64, len(a))

	for i := range a {
		mul[i] = a[i] * b[i]
	}

	return mul
}

// Sum the contents of a float64 slice
func sumArray(a []float64) float64 {
	sum := 0.0

	for _, f := range a {
		sum += f
	}

	return sum
}
