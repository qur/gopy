// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Slice represents objects of the SliceType (or PySlice_Type
// in the Python API) type.
type Slice struct {
	o C.PySliceObject
}

var _ Object = (*Slice)(nil)

// SliceType is the Type object that represents the Slice type.
var SliceType = newType(&C.PySlice_Type)

func sliceCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.sliceCheck(c(obj)) != 0
}

// AsSlice casts the given obj to a Slice (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Slice, then nil is returned.
//
// Return value: Borrowed Reference.
func AsSlice(obj Object) *Slice {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.sliceCheck(o) == 0 {
		return nil
	}
	return (*Slice)(unsafe.Pointer(o))
}

func newSlice(obj *C.PyObject) *Slice {
	return (*Slice)(unsafe.Pointer(obj))
}



func (s *Slice) c() *C.PySliceObject {
	return (*C.PySliceObject)(unsafe.Pointer(s))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (s *Slice) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(s))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (s *Slice) Type() *Type {
	return newType(c(s).ob_type)
}

// Decref decrements s's reference count, s may not be nil.
func (s *Slice) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(s))
	obj.ob_refcnt--
	if obj.ob_refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments s's reference count, s may not be nil.
func (s *Slice) Incref() {
	obj := (*C.PyObject)(unsafe.Pointer(s))
	obj.ob_refcnt++
}

// Repr returns a String representation of "s". This is equivalent to the
// Python "repr(s)".
//
// Return value: New Reference.
func (s *Slice) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(s))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of s. The equivalent
// Python is "hash(s)".
func (s *Slice) Hash() (int, error) {
	ret := C.PyObject_Hash(c(s))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// HasAttr returns true if "s" has the attribute "name".  This is equivalent
// to the Python "hasattr(s, name)".
func (s *Slice) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(s), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "s" with the name "name".  This is
// equivalent to the Python "s.name".
//
// Return value: New Reference.
func (s *Slice) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(s), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "s" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "s op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (s *Slice) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(s), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "s" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "s op obj",
// where op is the corresponding Python operator for op.
func (s *Slice) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(s), c(obj), C.int(op))
	return int2BoolErr(ret)
}



/*
set fields:
  ob_base
  tp_basicsize
  tp_dealloc
  tp_doc
  tp_flags
  tp_getattro
  tp_hash
  tp_members
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_traverse
*/

