package py

// #include "utils.h"
import "C"

import (
	"fmt"
)

func NewFloat[T ~float32 | ~float64](f T) (*Float, error) {
	ret := C.PyFloat_FromDouble(C.double(float64(f)))
	if ret == nil {
		return nil, exception()
	}
	return newFloat(ret), nil
}

func NewFloatString(v string) (*Float, error) {
	s, err := NewUnicode(v)
	if err != nil {
		return nil, err
	}
	defer s.Decref()
	ret := C.PyFloat_FromString(c(s))
	if ret == nil {
		return nil, exception()
	}
	return newFloat(ret), nil
}

func (f *Float) Float64() float64 {
	return float64(C.PyFloat_AsDouble(c(f)))
}

func (f *Float) String() string {
	if f == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", f.Float64())
}
