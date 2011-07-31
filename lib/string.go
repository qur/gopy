// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include <Python.h>
// static inline int stringCheck(PyObject *o) { return PyString_Check(o); }
import "C"

import (
//	"fmt"
//	"os"
//	"unsafe"
)

type String struct {
	BaseObject
}

func stringCheck(obj Object) bool {
	return C.stringCheck(c(obj)) != 0
}

func newString(obj *C.PyObject) *String {
	if obj == nil {
		return nil
	}
	return &String{BaseObject{obj}}
}

func (s *String) String() string {
	ret := C.PyString_AsString(c(s))
	if ret == nil {
		panic(exception())
	}
	return C.GoString(ret)
}
