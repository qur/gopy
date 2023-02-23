// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Bool is the representation of the Python bool type. There are only two
// possible values for a Bool, True and False. Every True value refers to the
// same instance, and every False value refers to the same value.
type Bool struct {
	o C.PyObject
}

var _ Object = (*Bool)(nil)

// BoolType is the Type object that represents the Bool type.
var BoolType = (*Type)(unsafe.Pointer(&C.PyBool_Type))

func boolCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.boolCheck(c(obj)) != 0
}



func (b *Bool) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(b))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (b *Bool) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(b))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (b *Bool) Type() *Type {
	return newType((*C.PyObject)(unsafe.Pointer(c(b).ob_type)))
}

// Decref decrements b's reference count, b may not be nil.
func (b *Bool) Decref() {
	C.decref(c(b))
}

// Incref increments b's reference count, b may not be nil.
func (b *Bool) Incref() {
	C.incref(c(b))
}

// IsTrue returns true if the value of b is considered to be True. This is
// equivalent to "if b:" in Python.
func (b *Bool) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(b))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of b is considered to be False. This is
// equivalent to "if not b:" in Python.
func (b *Bool) Not() bool {
	ret := C.PyObject_Not(c(b))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Free deallocates the storage (in Python) for b. After calling this method,
// b should no longer be used.
func (b *Bool) Free() {
	free(b)
}



/*
set fields:
*/
