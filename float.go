// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type Float struct {
	AbstractObject
	NumberProtocol
	o C.PyFloatObject
}

// FloatType is the Type object that represents the Float type.
var FloatType = (*Type)(unsafe.Pointer(&C.PyFloat_Type))

func floatCheck(obj Object) bool {
	return C.floatCheck(c(obj)) != 0
}

func newFloat(obj *C.PyObject) *Float {
	return (*Float)(unsafe.Pointer(obj))
}

func NewFloat(v float64) (*Float, error) {
	ret := C.PyFloat_FromDouble(C.double(v))
	if ret == nil {
		return nil, exception()
	}
	return newFloat(ret), nil
}

func NewFloatString(v string) (*Float, error) {
	s, err := NewString(v)
	if err != nil {
		return nil, err
	}
	defer s.Decref()
	ret := C.PyFloat_FromString(c(s), nil)
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
