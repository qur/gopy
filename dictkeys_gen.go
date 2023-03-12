// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// DictKeys represents objects of the DictKeysType (or PyDictKeys_Type
// in the Python API) type.
type DictKeys struct {
	o C.PyDictKeysObject
}

var _ Object = (*DictKeys)(nil)

// DictKeysType is the Type object that represents the DictKeys type.
var DictKeysType = newType(&C.PyDictKeys_Type)

func dictKeysCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.dictKeysCheck(c(obj)) != 0
}

func newDictKeys(obj *C.PyObject) *DictKeys {
	return (*DictKeys)(unsafe.Pointer(obj))
}



func (d *DictKeys) c() *C.PyDictKeysObject {
	return (*C.PyDictKeysObject)(unsafe.Pointer(d))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (d *DictKeys) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(d))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (d *DictKeys) Type() *Type {
	return newType(c(d).ob_type)
}

// Decref decrements d's reference count, d may not be nil.
func (d *DictKeys) Decref() {
	C.decref(c(d))
}

// Incref increments d's reference count, d may not be nil.
func (d *DictKeys) Incref() {
	C.incref(c(d))
}

func (d *DictKeys) raw() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(d))
}

// Repr returns a String representation of "d". This is equivalent to the
// Python "repr(d)".
//
// Return value: New Reference.
func (d *DictKeys) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(d))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// HasAttr returns true if "d" has the attribute "name".  This is equivalent
// to the Python "hasattr(d, name)".
func (d *DictKeys) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(d), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "d" with the name "name".  This is
// equivalent to the Python "d.name".
//
// Return value: New Reference.
func (d *DictKeys) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(d), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "d" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "d op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (d *DictKeys) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(d), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "d" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "d op obj",
// where op is the corresponding Python operator for op.
func (d *DictKeys) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(d), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// Iter returns an Iterator that will iterate over the members of d.
func (d *DictKeys) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c(d))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

// Size returns the size of d. The equivalent Python is "len(d)".
func (d *DictKeys) Size() int {
	ret := C.PyObject_Size(c(d))
	if ret < 0 {
		clearErr();
		return 0
	}
	return int(ret)
}

func (d *DictKeys) Contains(obj Object) (bool, error) {
	ret := C.PySequence_Contains(c(d), c(obj))
	return int2BoolErr(ret)
}

// Subtract returns the result of subtracting obj from d. The equivalent Python
// is "d - obj".
//
// Return value: New Reference.
func (d *DictKeys) Subtract(obj Object) (Object, error) {
	ret := C.PyNumber_Subtract(c(d), c(obj))
	return obj2ObjErr(ret)
}

// And returns the bitwise and of d and obj. The equivalent Python is
// "d & obj".
//
// Return value: New Reference.
func (d *DictKeys) And(obj Object) (Object, error) {
	ret := C.PyNumber_And(c(d), c(obj))
	return obj2ObjErr(ret)
}

// Xor returns the bitwise xor of d and obj. The equivalent Python is
// "d ^ obj".
//
// Return value: New Reference.
func (d *DictKeys) Xor(obj Object) (Object, error) {
	ret := C.PyNumber_Xor(c(d), c(obj))
	return obj2ObjErr(ret)
}

// Or returns the bitwise or of d and obj. The equivalent Python is
// "d | obj".
//
// Return value: New Reference.
func (d *DictKeys) Or(obj Object) (Object, error) {
	ret := C.PyNumber_Or(c(d), c(obj))
	return obj2ObjErr(ret)
}



/*
set fields:
  nb_and
  nb_or
  nb_subtract
  nb_xor
  ob_base
  sq_contains
  sq_length
  tp_as_number
  tp_as_sequence
  tp_basicsize
  tp_dealloc
  tp_flags
  tp_getattro
  tp_getset
  tp_iter
  tp_methods
  tp_name
  tp_repr
  tp_richcompare
  tp_traverse
*/

