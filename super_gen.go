// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Super represents objects of the SuperType (or PySuper_Type
// in the Python API) type.
type Super struct {
	o C.PyObject
}

var _ Object = (*Super)(nil)

// SuperType is the Type object that represents the Super type.
var SuperType = newType(&C.PySuper_Type)

func newSuper(obj *C.PyObject) *Super {
	return (*Super)(unsafe.Pointer(obj))
}



func (s *Super) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(s))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (s *Super) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(s))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (s *Super) Type() *Type {
	return newType(c(s).ob_type)
}

// Decref decrements s's reference count, s may not be nil.
func (s *Super) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(s))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments s's reference count, s may not be nil.
func (s *Super) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(s)).anon0[0]))
	*refcnt++
}



/*
set fields:
*/

