package py

// #include "utils.h"
import "C"

import (
	"fmt"
)

func NewComplex(v complex128) (*Complex, error) {
	ret := C.PyComplex_FromDoubles(C.double(real(v)), C.double(imag(v)))
	if ret == nil {
		return nil, exception()
	}
	return newComplex(ret), nil
}

func (o *Complex) Complex128() complex128 {
	r := float64(C.PyComplex_RealAsDouble(c(o)))
	i := float64(C.PyComplex_ImagAsDouble(c(o)))
	return complex(r, i)
}

func (c *Complex) String() string {
	if c == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", c.Complex128())
}
