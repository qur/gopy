// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// DictKeys represents objects of the DictKeysType (or PyDictKeysType
// in the Python API) type.
type DictKeys struct {
	o C.PyDictKeysObject
}

var _ Object = (*DictKeys)(nil)

// DictKeysType is the Type object that represents the DictKeys type.
var DictKeysType = (*Type)(unsafe.Pointer(&C.PyDictKeys_Type))

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
	return newType((*C.PyObject)(unsafe.Pointer(c(d).ob_type)))
}

// Decref decrements d's reference count, d may not be nil.
func (d *DictKeys) Decref() {
	C.decref(c(d))
}

// Incref increments d's reference count, d may not be nil.
func (d *DictKeys) Incref() {
	C.incref(c(d))
}

// IsTrue returns true if the value of d is considered to be True. This is
// equivalent to "if d:" in Python.
func (d *DictKeys) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(d))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of d is considered to be False. This is
// equivalent to "if not d:" in Python.
func (d *DictKeys) Not() bool {
	ret := C.PyObject_Not(c(d))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Free deallocates the storage (in Python) for d. After calling this method,
// d should no longer be used.
func (d *DictKeys) Free() {
	free(d)
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
	if ret == 1 {
		return true
	}
	return false
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

// RichCompare compares "obj" with "obj2" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "obj op obj2",
// where op is the corresponding Python operator for op.
func (d *DictKeys) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(d), c(obj), C.int(op))
	return int2BoolErr(ret)
}

func (d *DictKeys) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c(d))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

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

