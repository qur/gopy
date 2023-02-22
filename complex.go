package py

// #include "utils.h"
import "C"

import (
	"fmt"
)

func NewComplex[T ~complex64 | ~complex128](v T) (*Complex, error) {
	c := complex128(v)
	ret := C.PyComplex_FromDoubles(C.double(real(c)), C.double(imag(c)))
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
