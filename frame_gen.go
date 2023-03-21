// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Frame represents objects of the FrameType (or PyFrame_Type
// in the Python API) type.
type Frame struct {
	o C.PyFrameObject
}

var _ Object = (*Frame)(nil)

// FrameType is the Type object that represents the Frame type.
var FrameType = newType(&C.PyFrame_Type)

func frameCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.frameCheck(c(obj)) != 0
}

// AsFrame casts the given obj to a Frame (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Frame, then nil is returned.
//
// Return value: Borrowed Reference.
func AsFrame(obj Object) *Frame {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.frameCheck(o) == 0 {
		return nil
	}
	return (*Frame)(unsafe.Pointer(o))
}

func newFrame(obj *C.PyObject) *Frame {
	return (*Frame)(unsafe.Pointer(obj))
}



func (f *Frame) c() *C.PyFrameObject {
	return (*C.PyFrameObject)(unsafe.Pointer(f))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (f *Frame) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(f))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (f *Frame) Type() *Type {
	return newType(c(f).ob_type)
}

// Decref decrements f's reference count, f may not be nil.
func (f *Frame) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(f))
	obj.ob_refcnt--
	if obj.ob_refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments f's reference count, f may not be nil.
func (f *Frame) Incref() {
	obj := (*C.PyObject)(unsafe.Pointer(f))
	obj.ob_refcnt++
}

// Repr returns a String representation of "f". This is equivalent to the
// Python "repr(f)".
//
// Return value: New Reference.
func (f *Frame) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(f))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// HasAttr returns true if "f" has the attribute "name".  This is equivalent
// to the Python "hasattr(f, name)".
func (f *Frame) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(f), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "f" with the name "name".  This is
// equivalent to the Python "f.name".
//
// Return value: New Reference.
func (f *Frame) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(f), c(name))
	return obj2ObjErr(ret)
}

// SetAttr sets the attribute of "f" with the name "name" to "value".  This is
// equivalent to the Python "f.name = value".
func (f *Frame) SetAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c(f), c(name), c(value))
	return int2Err(ret)
}

// DelAttr deletes the attribute with the name "name" from "f".  This is
// equivalent to the Python "del f.name".
func (f *Frame) DelAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c(f), c(name), nil)
	return int2Err(ret)
}



/*
set fields:
  ob_base
  tp_basicsize
  tp_clear
  tp_dealloc
  tp_flags
  tp_getattro
  tp_getset
  tp_itemsize
  tp_members
  tp_methods
  tp_name
  tp_repr
  tp_setattro
  tp_traverse
*/

