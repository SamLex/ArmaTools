package unitcapturereduce

import "fmt"

// Basic immutable 3 component vector
type vec3 struct {
	A finalFloat64
	B finalFloat64
	C finalFloat64
}

// Element-wise linear interpolation been this vector and another
func (v vec3) lerp(vOther vec3, t float64) vec3 {
	return vec3{
		A: newFinalFloat64(lerp(v.A.Value(), vOther.A.Value(), t)),
		B: newFinalFloat64(lerp(v.B.Value(), vOther.B.Value(), t)),
		C: newFinalFloat64(lerp(v.C.Value(), vOther.C.Value(), t)),
	}
}

func (v vec3) SQFString() string {
	return fmt.Sprintf("[%v,%v,%v]", v.A.Value(), v.B.Value(), v.C.Value())
}
