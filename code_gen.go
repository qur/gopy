// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Code represents objects of the CodeType (or PyCodeType
// in the Python API) type.
type Code struct {
	o C.PyCodeObject
}

var _ Object = (*Code)(nil)

// CodeType is the Type object that represents the Code type.
var CodeType = newType(&C.PyCode_Type)

func codeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.codeCheck(c(obj)) != 0
}

func newCode(obj *C.PyObject) *Code {
	return (*Code)(unsafe.Pointer(obj))
}



func (co *Code) c() *C.PyCodeObject {
	return (*C.PyCodeObject)(unsafe.Pointer(co))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (co *Code) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(co))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (co *Code) Type() *Type {
	return newType(c(co).ob_type)
}

// Decref decrements co's reference count, co may not be nil.
func (co *Code) Decref() {
	C.decref(c(co))
}

// Incref increments co's reference count, co may not be nil.
func (co *Code) Incref() {
	C.incref(c(co))
}

// IsTrue returns true if the value of co is considered to be True. This is
// equivalent to "if co:" in Python.
func (co *Code) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(co))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of co is considered to be False. This is
// equivalent to "if not co:" in Python.
func (co *Code) Not() bool {
	ret := C.PyObject_Not(c(co))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Free deallocates the storage (in Python) for co. After calling this method,
// co should no longer be used.
func (co *Code) Free() {
	free(co)
}

// Repr returns a String representation of "co". This is equivalent to the
// Python "repr(co)".
//
// Return value: New Reference.
func (co *Code) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(co))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of co. The equivalent
// Python is "hash(co)".
func (co *Code) Hash() (int, error) {
	ret := C.PyObject_Hash(c(co))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// HasAttr returns true if "co" has the attribute "name".  This is equivalent
// to the Python "hasattr(co, name)".
func (co *Code) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(co), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "co" with the name "name".  This is
// equivalent to the Python "co.name".
//
// Return value: New Reference.
func (co *Code) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(co), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "co" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "co op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (co *Code) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(co), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "co" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "co op obj",
// where op is the corresponding Python operator for op.
func (co *Code) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(co), c(obj), C.int(op))
	return int2BoolErr(ret)
}



/*
set fields:
  ob_base
  tp_basicsize
  tp_dealloc
  tp_doc
  tp_getattro
  tp_getset
  tp_hash
  tp_itemsize
  tp_members
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_weaklistoffset
*/

