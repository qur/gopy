// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

// AbstractObject is an 0-sized type that can be embedded as the first item in
// concrete types to provide the Object interface functions.
type AbstractObject struct{}

func newAbstractObject(obj *C.PyObject) *AbstractObject {
	return (*AbstractObject)(unsafe.Pointer(obj))
}

// Init initialises obj.  It is equivalent to "obj.__init__(*args, **kw)" in
// Python.
func (obj *AbstractObject) Init(args *Tuple, kw *Dict) error {
	return obj.Type().Init(obj, args, kw)
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (obj *AbstractObject) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(obj))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (obj *AbstractObject) Type() *Type {
	o := c(obj).ob_type
	return newType((*C.PyObject)(unsafe.Pointer(o)))
}

// Decref decrements obj's reference count, obj may be nil.
func Decref(obj Object) {
	if obj != nil {
		C.decref(c(obj))
	}
}

// Decref decrements obj's reference count, obj may not be nil.
func (obj *AbstractObject) Decref() {
	C.decref(c(obj))
}

// Incref increments obj's reference count, obj may be nil.
func Incref(obj Object) {
	if obj != nil {
		C.incref(c(obj))
	}
}

// Incref increments obj's reference count, obj may not be nil.
func (obj *AbstractObject) Incref() {
	C.incref(c(obj))
}

// IsTrue returns true if the value of obj is considered to be True.  This is
// equivalent to "if obj:" in Python.
func (obj *AbstractObject) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(obj))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of obj is considered to be False.  This is
// equivalent to "if not obj:" in Python.
func (obj *AbstractObject) Not() bool {
	ret := C.PyObject_Not(c(obj))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Free deallocates the storage (in Python) for obj.  After calling this method,
// obj should no longer be used.
func (obj *AbstractObject) Free() {
	o := c(obj)

	// Make sure this instance isn't registered anymore
	clearClassContext(unsafe.Pointer(o))

	// Call Python free function
	pyType := (*C.PyTypeObject)(unsafe.Pointer(o.ob_type))
	C.typeFree(pyType, o)
}
