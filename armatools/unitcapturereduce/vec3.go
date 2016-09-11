/*
Copyright 2016 Euan James Hunter

vec3.go: 3d Vector utilities for armatools/unitcapturereduce

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

import "fmt"

// Basic immutable 3 component vector
type vec3 struct {
	A float64
	B float64
	C float64
}

// Element-wise linear interpolation been this vector and another
func (v vec3) lerp(vOther vec3, t float64) vec3 {
	return vec3{
		A: lerp(v.A, vOther.A, t),
		B: lerp(v.B, vOther.B, t),
		C: lerp(v.C, vOther.C, t),
	}
}

// Represents a vec3 as a SQF array with 3 elements
func (v vec3) SQFString() string {
	return fmt.Sprintf("[%v,%v,%v]", v.A, v.B, v.C)
}
