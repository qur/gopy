// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline void incref(PyObject *obj) { Py_INCREF(obj); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
// static inline void typeFree(PyTypeObject *type, PyObject *o) { type->tp_free(o); }
import "C"

import "unsafe"

type AbstractObject struct{}

func newAbstractObject(obj *C.PyObject) *AbstractObject {
	return (*AbstractObject)(unsafe.Pointer(obj))
}

func (obj *AbstractObject) Init(args *Tuple, kw *Dict) error {
	return obj.Type().Init(obj, args, kw)
}

func (obj *AbstractObject) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(obj))
}

func (obj *AbstractObject) Type() *Type {
	o := c(obj).ob_type
	return newType((*C.PyObject)(unsafe.Pointer(o)))
}

func Decref(obj Object) {
	if obj != nil {
		C.decref(c(obj))
	}
}

func (obj *AbstractObject) Decref() {
	C.decref(c(obj))
}

func Incref(obj Object) {
	if obj != nil {
		C.incref(c(obj))
	}
}

func (obj *AbstractObject) Incref() {
	C.incref(c(obj))
}

func (obj *AbstractObject) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(obj))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

func (obj *AbstractObject) Not() bool {
	ret := C.PyObject_Not(c(obj))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

func (obj *AbstractObject) Free() {
	o := c(obj)
	pyType := (*C.PyTypeObject)(unsafe.Pointer(o.ob_type))
	C.typeFree(pyType, o)

	// Make sure this instance isn't registered anymore
	delete(contexts, uintptr(unsafe.Pointer(o)))
}
