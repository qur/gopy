package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// None is the Python equivalent to nil.
var None = newNone(&C._Py_NoneStruct)

// NoneObject is the type of the None value.  The only value of this type is
// None.
type NoneObject struct {
	c C.PyObject
}

var _ Object = (*NoneObject)(nil)

func newNone(obj *C.PyObject) *NoneObject {
	return (*NoneObject)(unsafe.Pointer(obj))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (n *NoneObject) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(n))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (n *NoneObject) Type() *Type {
	return newType(c(n).ob_type)
}

// Decref decrements n's reference count, n may not be nil.
func (n *NoneObject) Decref() {
	C.decref(c(n))
}

// Incref increments n's reference count, n may not be nil.
func (n *NoneObject) Incref() {
	C.incref(c(n))
}

func (n *NoneObject) raw() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(n))
}

// IsTrue returns true if the value of n is considered to be True. This is
// equivalent to "if n:" in Python.
func (n *NoneObject) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(n))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of n is considered to be False. This is
// equivalent to "if not n:" in Python.
func (n *NoneObject) Not() bool {
	ret := C.PyObject_Not(c(n))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Free deallocates the storage (in Python) for n. After calling this method,
// n should no longer be used.
func (n *NoneObject) Free() {
	free(n)
}

// String returns a string representation of the NoneObject type. This function
// always returns "None".
func (*NoneObject) String() string {
	return "None"
}

// ReturnNone returns None (after incrementing the reference count).
// This is a convenience function for returning None from a function.
//
// Return value: New Reference.
func ReturnNone() Object {
	None.Incref()
	return None
}
