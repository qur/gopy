// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// BufferMethods represents an arbitrary Python type that implements the Buffer
// Protocol.
type BufferMethods struct {
	o C.PyObject
}

var _ Object = (*BufferMethods)(nil)





func (b *BufferMethods) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(b))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (b *BufferMethods) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(b))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (b *BufferMethods) Type() *Type {
	return newType(c(b).ob_type)
}

// Decref decrements b's reference count, b may not be nil.
func (b *BufferMethods) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(b))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments b's reference count, b may not be nil.
func (b *BufferMethods) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(b)).anon0[0]))
	*refcnt++
}



/*
set fields:
*/

