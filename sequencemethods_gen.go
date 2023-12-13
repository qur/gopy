// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// SequenceMethods is a concrete realisation of the full set of Sequence
// Protocol methods. A type that implements the "Sequence Protocol" can be
// turned into a SequenceMethods instance using AsSequence.
//
// Note that the methods not already implemented on the type itself may return
// an error, as not all methods are implemented by all types that support the
// protocol.
type SequenceMethods struct {
	o C.PyObject
}

var _ Object = (*SequenceMethods)(nil)





func (s *SequenceMethods) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(s))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (s *SequenceMethods) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(s))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (s *SequenceMethods) Type() *Type {
	return newType(c(s).ob_type)
}


// Decref decrements s's reference count, s may not be nil.
func (s *SequenceMethods) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(s))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments s's reference count, s may not be nil.
func (s *SequenceMethods) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(s)).anon0[0]))
	*refcnt++
}



/*
set fields:
*/

