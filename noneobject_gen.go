// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// NoneObject is the type of the None value.  The only value of this type is
// None.
type NoneObject struct {
	o C.PyObject
}

var _ Object = (*NoneObject)(nil)



func newNoneObject(obj *C.PyObject) *NoneObject {
	return (*NoneObject)(unsafe.Pointer(obj))
}



func (n *NoneObject) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(n))
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
	obj := (*C.PyObject)(unsafe.Pointer(n))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments n's reference count, n may not be nil.
func (n *NoneObject) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(n)).anon0[0]))
	*refcnt++
}



/*
set fields:
*/

