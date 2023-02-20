// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// ODict represents objects of the ODictType (or PyODictType
// in the Python API) type.
type ODict struct {
	abstractObject
	o C.PyODictObject
}

// ODictType is the Type object that represents the ODict type.
var ODictType = (*Type)(unsafe.Pointer(&C.PyODict_Type))

func oDictCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.oDictCheck(c(obj)) != 0
}

func newODict(obj *C.PyObject) *ODict {
	return (*ODict)(unsafe.Pointer(obj))
}

// Repr returns a String representation of "o". This is equivalent to the
// Python "repr(o)".
//
// Return value: New Reference.
func (o *ODict) Repr() (Object, error) {
	ret := C.PyObject_Repr(c(o))
	return obj2ObjErr(ret)
}

// RichCompare compares "o" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "o op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (o *ODict) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(o), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "obj" with "obj2" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "obj op obj2",
// where op is the corresponding Python operator for op.
func (o *ODict) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(o), c(obj), C.int(op))
	return int2BoolErr(ret)
}

func (o *ODict) Iter() (Object, error) {
	ret := C.PyObject_GetIter(c(o))
	return obj2ObjErr(ret)
}

func (o *ODict) DelItem(key Object) error {
	ret := C.PyObject_DelItem(c(o), c(key))
	return int2Err(ret)
}

func (o *ODict) DelItemString(key string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyObject_DelItemString(c(o), cKey)
	return int2Err(ret)
}

func (o *ODict) SetItemString(key string, v Object) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_SetItemString(c(o), cKey, c(v))
	return int2Err(ret)
}

// Or returns the bitwise or of o and obj. The equivalent Python is
// "o | obj".
//
// Return value: New Reference.
func (o *ODict) Or(obj Object) (Object, error) {
	ret := C.PyNumber_Or(c(o), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceOr returns the bitwise or of o and obj. This is done in place. The
// equivalent Python is "o |= obj".
//
// Return value: New Reference.
func (o *ODict) InPlaceOr(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceOr(c(o), c(obj))
	return obj2ObjErr(ret)
}



/*
set fields:
  mp_ass_subscript
  nb_inplace_or
  nb_or
  ob_base
  tp_alloc
  tp_as_mapping
  tp_as_number
  tp_base
  tp_basicsize
  tp_clear
  tp_dealloc
  tp_dictoffset
  tp_doc
  tp_flags
  tp_getset
  tp_init
  tp_iter
  tp_methods
  tp_name
  tp_repr
  tp_richcompare
  tp_traverse
  tp_weaklistoffset
*/
