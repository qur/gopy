// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// List represents objects of the ListType (or PyList_Type
// in the Python API) type.
//
// This type implements the Mapping protocol.
//
// This type implements the Sequence protocol.
type List struct {
	o C.PyListObject
}

var _ Object = (*List)(nil)

// ListType is the Type object that represents the List type.
var ListType = newType(&C.PyList_Type)

func listCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.listCheck(c(obj)) != 0
}

// AsList casts the given obj to a List (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a List, then nil is returned.
//
// Return value: Borrowed Reference.
func AsList(obj Object) *List {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.listCheck(o) == 0 {
		return nil
	}
	return (*List)(unsafe.Pointer(o))
}

func newList(obj *C.PyObject) *List {
	return (*List)(unsafe.Pointer(obj))
}



func (l *List) c() *C.PyListObject {
	return (*C.PyListObject)(unsafe.Pointer(l))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (l *List) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(l))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (l *List) Type() *Type {
	return newType(c(l).ob_type)
}

// Decref decrements l's reference count, l may not be nil.
func (l *List) Decref() {
	C.decref(c(l))
}

// Incref increments l's reference count, l may not be nil.
func (l *List) Incref() {
	C.incref(c(l))
}

func (l *List) raw() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(l))
}

// Repr returns a String representation of "l". This is equivalent to the
// Python "repr(l)".
//
// Return value: New Reference.
func (l *List) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(l))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of l. The equivalent
// Python is "hash(l)".
func (l *List) Hash() (int, error) {
	ret := C.PyObject_Hash(c(l))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// HasAttr returns true if "l" has the attribute "name".  This is equivalent
// to the Python "hasattr(l, name)".
func (l *List) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(l), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "l" with the name "name".  This is
// equivalent to the Python "l.name".
//
// Return value: New Reference.
func (l *List) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(l), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "l" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "l op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (l *List) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(l), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "l" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "l op obj",
// where op is the corresponding Python operator for op.
func (l *List) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(l), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// Iter returns an Iterator that will iterate over the members of l.
func (l *List) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c(l))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

// Size returns the size of l. The equivalent Python is "len(l)".
func (l *List) Size() int {
	ret := C.PyObject_Size(c(l))
	if ret < 0 {
		clearErr();
		return 0
	}
	return int(ret)
}

// AsMappingMethods returns a MappingMethods instance that refers to the same
// underlying Python object as l.
//
// This method also means that List implements the MappingProtocol interface.
func (l *List) AsMappingMethods() *MappingMethods {
	return (*MappingMethods)(unsafe.Pointer(l.Base()))
}

func (l *List) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(c(l), cKey)
	return obj2ObjErr(ret)
}

func (l *List) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(c(l), c(key))
	clearErr();
	return ret > 0
}

func (l *List) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(c(l), cKey)
	clearErr();
	return ret > 0
}

func (l *List) DelItem(key Object) error {
	ret := C.PyObject_DelItem(c(l), c(key))
	return int2Err(ret)
}

func (l *List) DelItemString(key string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyObject_DelItemString(c(l), cKey)
	return int2Err(ret)
}

func (l *List) SetItemString(key string, v Object) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_SetItemString(c(l), cKey, c(v))
	return int2Err(ret)
}

// AsSequenceMethods returns a SequenceMethods instance that refers to the same
// underlying Python object as l.
//
// This method also means that List implements the SequenceProtocol
// interface.
func (l *List) AsSequenceMethods() *SequenceMethods {
	return (*SequenceMethods)(unsafe.Pointer(l.Base()))
}

func (l *List) GetIndex(idx int) (Object, error) {
	ret := C.PySequence_GetItem(c(l), C.Py_ssize_t(idx))
	return obj2ObjErr(ret)
}

func (l *List) Count(obj Object) (int, error) {
	ret := C.PySequence_Count(c(l), c(obj))
	return ssize_t2IntErr(ret)
}

func (l *List) Index(obj Object) (int, error) {
	ret := C.PySequence_Index(c(l), c(obj))
	return ssize_t2IntErr(ret)
}

func (l *List) Tuple() (*Tuple, error) {
	ret := C.PySequence_Tuple(c(l))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

func (l *List) SetIndex(idx int, obj Object) error {
	ret := C.PySequence_SetItem(c(l), C.Py_ssize_t(idx), c(obj))
	return int2Err(ret)
}

func (l *List) DelIndex(idx int) error {
	ret := C.PySequence_DelItem(c(l), C.Py_ssize_t(idx))
	return int2Err(ret)
}

func (l *List) Concat(obj Object) (Object, error) {
	ret := C.PySequence_Concat(c(l), c(obj))
	return obj2ObjErr(ret)
}

func (l *List) InPlaceConcat(obj Object) (Object, error) {
	ret := C.PySequence_InPlaceConcat(c(l), c(obj))
	return obj2ObjErr(ret)
}

func (l *List) Repeat(count int) (Object, error) {
	ret := C.PySequence_Repeat(c(l), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (l *List) InPlaceRepeat(count int) (Object, error) {
	ret := C.PySequence_InPlaceRepeat(c(l), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (l *List) Contains(obj Object) (bool, error) {
	ret := C.PySequence_Contains(c(l), c(obj))
	return int2BoolErr(ret)
}

func (l *List) GetSlice(start, end int) (Object, error) {
	ret := C.PySequence_GetSlice(c(l), C.Py_ssize_t(start), C.Py_ssize_t(end))
	return obj2ObjErr(ret)
}

func (l *List) SetSlice(start, end int, obj Object) error {
	ret := C.PySequence_SetSlice(c(l), C.Py_ssize_t(start), C.Py_ssize_t(end), c(obj))
	return int2Err(ret)
}

func (l *List) DelSlice(start, end int) error {
	ret := C.PySequence_DelSlice(c(l), C.Py_ssize_t(start), C.Py_ssize_t(end))
	return int2Err(ret)
}



/*
set fields:
  mp_ass_subscript
  mp_length
  mp_subscript
  ob_base
  sq_ass_item
  sq_concat
  sq_contains
  sq_inplace_concat
  sq_inplace_repeat
  sq_item
  sq_length
  sq_repeat
  tp_alloc
  tp_as_mapping
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

