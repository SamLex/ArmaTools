package unitcapturereduce

import (
	"fmt"
	"math"
)

// Basic 3 component vector
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

// Maximum element-wise percentage difference
func (v vec3) maxPercentDifference(vOther vec3) float64 {
	aDifference := percentDifference(v.A, vOther.A)
	bDifference := percentDifference(v.B, vOther.B)
	cDifference := percentDifference(v.C, vOther.C)

	return math.Max(aDifference, math.Max(bDifference, cDifference))
}

func (v vec3) String() string {
	return fmt.Sprintf("[%v,%v,%v]", v.A, v.B, v.C)
}
