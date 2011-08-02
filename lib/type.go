// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include <Python.h>
// #include "utils.h"
// static inline int typeCheck(PyObject *o) { return PyType_Check(o); }
// static inline int typeCheckE(PyObject *o) { return PyType_CheckExact(o); }
import "C"

import (
	"os"
	"unsafe"
)

type Type struct {
	BaseObject
}

func typeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.typeCheck(c(obj)) != 0
}

func newType(obj *C.PyObject) *Type {
	return (*Type)(unsafe.Pointer(obj))
}

func (t *Type) Alloc(n int64) (Object, os.Error) {
	ct := (*C.PyTypeObject)(unsafe.Pointer(c(t)))
	ret := C.PyType_GenericAlloc(ct, C.Py_ssize_t(n))
	return obj2ObjErr(ret)
}
