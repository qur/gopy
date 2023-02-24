// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Frame represents objects of the FrameType (or PyFrameType
// in the Python API) type.
type Frame struct {
	o C.PyFrameObject
}

var _ Object = (*Frame)(nil)

// FrameType is the Type object that represents the Frame type.
var FrameType = (*Type)(unsafe.Pointer(&C.PyFrame_Type))

func frameCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.frameCheck(c(obj)) != 0
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
	return newType((*C.PyObject)(unsafe.Pointer(c(f).ob_type)))
}

// Decref decrements f's reference count, f may not be nil.
func (f *Frame) Decref() {
	C.decref(c(f))
}

// Incref increments f's reference count, f may not be nil.
func (f *Frame) Incref() {
	C.incref(c(f))
}

// IsTrue returns true if the value of f is considered to be True. This is
// equivalent to "if f:" in Python.
func (f *Frame) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(f))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of f is considered to be False. This is
// equivalent to "if not f:" in Python.
func (f *Frame) Not() bool {
	ret := C.PyObject_Not(c(f))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Free deallocates the storage (in Python) for f. After calling this method,
// f should no longer be used.
func (f *Frame) Free() {
	free(f)
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
	if ret == 1 {
		return true
	}
	return false
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

