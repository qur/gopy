// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// AsyncIteratorMethods represents an arbitrary async Python iterator.
type AsyncIteratorMethods struct {
	o C.PyObject
}

var _ Object = (*AsyncIteratorMethods)(nil)





func (a *AsyncIteratorMethods) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(a))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (a *AsyncIteratorMethods) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(a))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (a *AsyncIteratorMethods) Type() *Type {
	return newType(c(a).ob_type)
}

// Decref decrements a's reference count, a may not be nil.
func (a *AsyncIteratorMethods) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(a))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments a's reference count, a may not be nil.
func (a *AsyncIteratorMethods) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(a)).anon0[0]))
	*refcnt++
}



/*
set fields:
*/

