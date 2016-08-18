// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int complexCheck(PyObject *o) { return PyComplex_Check(o); }
import "C"

import (
	"fmt"
	"unsafe"
)

type Complex struct {
	AbstractObject
	NumberProtocol
	o *C.PyComplexObject
}

var complexObjMap = make(map[*C.PyObject]*Complex)

// ComplexType is the Type object that represents the Complex type.
var ComplexType = newType((*C.PyObject)(unsafe.Pointer(C.getBasePyType(C.GoPyComplex_Type))))

func complexCheck(obj Object) bool {
	return C.complexCheck(c(obj)) != 0
}

func newComplex(obj *C.PyObject) *Complex {
	if c, ok := complexObjMap[obj]; ok {
		return c
	}
	c := &Complex{o: (*C.PyComplexObject)(unsafe.Pointer(obj))}
	complexObjMap[obj] = c
	return c
}

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
