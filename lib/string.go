// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int stringCheck(PyObject *o) { return PyString_Check(o); }
import "C"

import (
	"os"
	"unsafe"
)

type String struct {
	BaseObject
}

func stringCheck(obj Object) bool {
	return C.stringCheck(c(obj)) != 0
}

func newString(obj *C.PyObject) *String {
	return (*String)(unsafe.Pointer(obj))
}

func NewString(s string) (*String, os.Error) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	ret := C.PyString_FromString(cs)
	if ret == nil {
		return nil, exception()
	}
	return newString(ret), nil
}

func (s *String) String() string {
	if s == nil {
		return "<nil>"
	}
	ret := C.PyString_AsString(c(s))
	if ret == nil {
		panic(exception())
	}
	return C.GoString(ret)
}

func (s *String) Format(args *Tuple) (*String, os.Error) {
	ret := C.PyString_Format(c(s), c(args))
	if ret == nil {
		return nil, exception()
	}
	return newString(ret), nil
}
