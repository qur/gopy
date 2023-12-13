// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Dict represents objects of the DictType (or PyDict_Type
// in the Python API) type.
//
// This type implements the Mapping protocol.
type Dict struct {
	o C.PyDictObject
}

var _ Object = (*Dict)(nil)

// DictType is the Type object that represents the Dict type.
var DictType = newType(&C.PyDict_Type)

func dictCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.dictCheck(c(obj)) != 0
}

// AsDict casts the given obj to a Dict (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Dict, then nil is returned.
//
// Return value: Borrowed Reference.
func AsDict(obj Object) *Dict {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.dictCheck(o) == 0 {
		return nil
	}
	return (*Dict)(unsafe.Pointer(o))
}

func newDict(obj *C.PyObject) *Dict {
	return (*Dict)(unsafe.Pointer(obj))
}



func (d *Dict) c() *C.PyDictObject {
	return (*C.PyDictObject)(unsafe.Pointer(d))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (d *Dict) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(d))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (d *Dict) Type() *Type {
	return newType(c(d).ob_type)
}

// Decref decrements d's reference count, d may not be nil.
func (d *Dict) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(d))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments d's reference count, d may not be nil.
func (d *Dict) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(d)).anon0[0]))
	*refcnt++
}

// Repr returns a String representation of "d". This is equivalent to the
// Python "repr(d)".
//
// Return value: New Reference.
func (d *Dict) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(d))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of d. The equivalent
// Python is "hash(d)".
func (d *Dict) Hash() (int, error) {
	ret := C.PyObject_Hash(c(d))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// HasAttr returns true if "d" has the attribute "name".  This is equivalent
// to the Python "hasattr(d, name)".
func (d *Dict) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(d), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "d" with the name "name".  This is
// equivalent to the Python "d.name".
//
// Return value: New Reference.
func (d *Dict) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(d), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "d" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "d op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (d *Dict) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(d), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "d" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "d op obj",
// where op is the corresponding Python operator for op.
func (d *Dict) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(d), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// Iter returns an Iterator that will iterate over the members of d.
func (d *Dict) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c(d))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

// Size returns the size of d. The equivalent Python is "len(d)".
func (d *Dict) Size() int {
	ret := C.PyObject_Size(c(d))
	if ret < 0 {
		clearErr();
		return 0
	}
	return int(ret)
}

// AsMappingMethods returns a MappingMethods instance that refers to the same
// underlying Python object as d.
//
// This method also means that Dict implements the MappingProtocol interface.
func (d *Dict) AsMappingMethods() *MappingMethods {
	return (*MappingMethods)(unsafe.Pointer(d.Base()))
}

func (d *Dict) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(c(d), cKey)
	return obj2ObjErr(ret)
}

func (d *Dict) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(c(d), c(key))
	clearErr();
	return ret > 0
}

func (d *Dict) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(c(d), cKey)
	clearErr();
	return ret > 0
}

func (d *Dict) DelItem(key Object) error {
	ret := C.PyObject_DelItem(c(d), c(key))
	return int2Err(ret)
}

func (d *Dict) DelItemString(key string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyObject_DelItemString(c(d), cKey)
	return int2Err(ret)
}

func (d *Dict) SetItemString(key string, v Object) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_SetItemString(c(d), cKey, c(v))
	return int2Err(ret)
}

func (d *Dict) Contains(obj Object) (bool, error) {
	ret := C.PySequence_Contains(c(d), c(obj))
	return int2BoolErr(ret)
}

// Or returns the bitwise or of d and obj. The equivalent Python is
// "d | obj".
//
// Return value: New Reference.
func (d *Dict) Or(obj Object) (Object, error) {
	ret := C.PyNumber_Or(c(d), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceOr returns the bitwise or of d and obj. This is done in place. The
// equivalent Python is "d |= obj".
//
// Return value: New Reference.
func (d *Dict) InPlaceOr(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceOr(c(d), c(obj))
	return obj2ObjErr(ret)
}



/*
set fields:
  mp_ass_subscript
  mp_length
  mp_subscript
  nb_inplace_or
  nb_or
  ob_base
  sq_contains
  tp_alloc
  tp_as_mapping
  tp_as_number
  tp_as_sequence
  tp_basicsize
  tp_clear
  tp_dealloc
  tp_doc
  tp_flags
  tp_free
  tp_getattro
  tp_hash
  tp_init
  tp_iter
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_traverse
  tp_vectorcall
*/

