// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int floatCheck(PyObject *o) { return PyFloat_Check(o); }
import "C"

import (
	"fmt"
	"unsafe"
)

type Float struct {
	AbstractObject
	NumberProtocol
	o *C.PyFloatObject
}

var floatObjMap = make(map[*C.PyObject]*Float)

// FloatType is the Type object that represents the Float type.
var FloatType = newType((*C.PyObject)(unsafe.Pointer(C.getBasePyType(C.GoPyFloat_Type))))

func floatCheck(obj Object) bool {
	return C.floatCheck(c(obj)) != 0
}

func newFloat(obj *C.PyObject) *Float {
	if f, ok := floatObjMap[obj]; ok {
		return f
	}
	f := &Float{o: (*C.PyFloatObject)(unsafe.Pointer(obj))}
	floatObjMap[obj] = f
	return f
}

func NewFloat(v float64) (*Float, error) {
	ret := C.PyFloat_FromDouble(C.double(v))
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
