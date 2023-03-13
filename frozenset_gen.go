// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// FrozenSet represents objects of the FrozenSetType (or PyFrozenSet_Type
// in the Python API) type.
type FrozenSet struct {
	o C.PyObject
}

var _ Object = (*FrozenSet)(nil)

// FrozenSetType is the Type object that represents the FrozenSet type.
var FrozenSetType = newType(&C.PyFrozenSet_Type)

func frozenSetCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.frozenSetCheck(c(obj)) != 0
}

// AsFrozenSet casts the given obj to a FrozenSet (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a FrozenSet, then nil is returned.
//
// Return value: Borrowed Reference.
func AsFrozenSet(obj Object) *FrozenSet {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.frozenSetCheck(o) == 0 {
		return nil
	}
	return (*FrozenSet)(unsafe.Pointer(o))
}

func newFrozenSet(obj *C.PyObject) *FrozenSet {
	return (*FrozenSet)(unsafe.Pointer(obj))
}



func (f *FrozenSet) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(f))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (f *FrozenSet) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(f))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (f *FrozenSet) Type() *Type {
	return newType(c(f).ob_type)
}

// Decref decrements f's reference count, f may not be nil.
func (f *FrozenSet) Decref() {
	C.decref(c(f))
}

// Incref increments f's reference count, f may not be nil.
func (f *FrozenSet) Incref() {
	C.incref(c(f))
}

func (f *FrozenSet) raw() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(f))
}



/*
set fields:
*/

