// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int typeCheck(PyObject *o) { return PyType_Check(o); }
// static inline int typeCheckE(PyObject *o) { return PyType_CheckExact(o); }
// static inline int typeInit(PyObject *t, PyObject *o, PyObject *a, PyObject *k) {
//    return ((PyTypeObject *)t)->tp_init(o, a, k);
// }
import "C"

import (
	"os"
	"unsafe"
)

type Type struct {
	AbstractObject
	o C.PyTypeObject
}

var TypeType = (*Type)(unsafe.Pointer(&C.PyType_Type))

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

func (t *Type) Init(obj Object, args *Tuple, kw *Dict) os.Error {
	ret := C.typeInit(c(t), c(obj), c(args), c(kw))
	if ret < 0 {
		return exception()
	}
	return nil
}

func (t *Type) String() string {
	pyS := C.PyObject_Str(c(t))
	if pyS == nil {
		return "<unknown type>"
	}
	return C.GoString(C.PyString_AsString(pyS))
}
