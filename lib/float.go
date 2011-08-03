// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int floatCheck(PyObject *o) { return PyFloat_Check(o); }
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type Float struct {
	BaseObject
}

func floatCheck(obj Object) bool {
	return C.floatCheck(c(obj)) != 0
}

func newFloat(obj *C.PyObject) *Float {
	return (*Float)(unsafe.Pointer(obj))
}

func Float_FromFloat64(v float64) (*Float, os.Error) {
	ret := C.PyFloat_FromDouble(C.double(v))
	if ret == nil {
		return nil, exception()
	}
	return newFloat(ret), nil
}

func Float_FromString(v string) (*Float, os.Error) {
	s, err := String_FromString(v)
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
	return fmt.Sprintf("%v", f.Float64())
}
