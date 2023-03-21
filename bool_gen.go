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
var BoolType = newType(&C.PyBool_Type)

func boolCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.boolCheck(c(obj)) != 0
}

// AsBool casts the given obj to a Bool (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Bool, then nil is returned.
//
// Return value: Borrowed Reference.
func AsBool(obj Object) *Bool {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.boolCheck(o) == 0 {
		return nil
	}
	return (*Bool)(unsafe.Pointer(o))
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
	return newType(c(b).ob_type)
}

// Decref decrements b's reference count, b may not be nil.
func (b *Bool) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(b))
	obj.ob_refcnt--
	if obj.ob_refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments b's reference count, b may not be nil.
func (b *Bool) Incref() {
	obj := (*C.PyObject)(unsafe.Pointer(b))
	obj.ob_refcnt++
}



/*
set fields:
*/

