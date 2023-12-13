// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// IterableMethods represents an arbitrary Python iterable.
type IterableMethods struct {
	o C.PyObject
}

var _ Object = (*IterableMethods)(nil)





func (i *IterableMethods) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(i))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (i *IterableMethods) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(i))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (i *IterableMethods) Type() *Type {
	return newType(c(i).ob_type)
}


// Decref decrements i's reference count, i may not be nil.
func (i *IterableMethods) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(i))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	if *refcnt == C._Py_IMMORTAL_REFCNT {
		return
	}
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments i's reference count, i may not be nil.
func (i *IterableMethods) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(i)).anon0[0]))
	if *refcnt == C._Py_IMMORTAL_REFCNT {
		return
	}
	*refcnt++
}



/*
set fields:
*/

