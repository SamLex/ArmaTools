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
