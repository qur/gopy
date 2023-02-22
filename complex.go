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

func (co *Complex) Complex128() complex128 {
	r := float64(C.PyComplex_RealAsDouble(c(co)))
	i := float64(C.PyComplex_ImagAsDouble(c(co)))
	return complex(r, i)
}

func (co *Complex) String() string {
	if co == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", co.Complex128())
}
