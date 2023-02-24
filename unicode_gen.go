// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Unicode represents objects of the UnicodeType (or PyUnicodeType
// in the Python API) type.
//
// This type implements the Mapping protocol.
//
// This type implements the Sequence protocol.
type Unicode struct {
	o C.PyUnicodeObject
}

var _ Object = (*Unicode)(nil)

// UnicodeType is the Type object that represents the Unicode type.
var UnicodeType = (*Type)(unsafe.Pointer(&C.PyUnicode_Type))

func unicodeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.unicodeCheck(c(obj)) != 0
}

func newUnicode(obj *C.PyObject) *Unicode {
	return (*Unicode)(unsafe.Pointer(obj))
}

func (u *Unicode) c() *C.PyUnicodeObject {
	return (*C.PyUnicodeObject)(unsafe.Pointer(u))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (u *Unicode) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(u))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (u *Unicode) Type() *Type {
	return newType((*C.PyObject)(unsafe.Pointer(c(u).ob_type)))
}

// Decref decrements u's reference count, u may not be nil.
func (u *Unicode) Decref() {
	C.decref(c(u))
}

// Incref increments u's reference count, u may not be nil.
func (u *Unicode) Incref() {
	C.incref(c(u))
}

// IsTrue returns true if the value of u is considered to be True. This is
// equivalent to "if u:" in Python.
func (u *Unicode) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(u))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of u is considered to be False. This is
// equivalent to "if not u:" in Python.
func (u *Unicode) Not() bool {
	ret := C.PyObject_Not(c(u))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Free deallocates the storage (in Python) for u. After calling this method,
// u should no longer be used.
func (u *Unicode) Free() {
	free(u)
}

// Repr returns a String representation of "u". This is equivalent to the
// Python "repr(u)".
//
// Return value: New Reference.
func (u *Unicode) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(u))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of u. The equivalent
// Python is "hash(u)".
func (u *Unicode) Hash() (int, error) {
	ret := C.PyObject_Hash(c(u))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// Str returns a String representation of "u". This is equivalent to the
// Python "str(u)".
//
// Return value: New Reference.
func (u *Unicode) Str() (*Unicode, error) {
	ret := C.PyObject_Str(c(u))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// HasAttr returns true if "u" has the attribute "name".  This is equivalent
// to the Python "hasattr(u, name)".
func (u *Unicode) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(u), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "u" with the name "name".  This is
// equivalent to the Python "u.name".
//
// Return value: New Reference.
func (u *Unicode) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(u), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "u" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "u op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (u *Unicode) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(u), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "u" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "u op obj",
// where op is the corresponding Python operator for op.
func (u *Unicode) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(u), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// Iter returns an Iterator that will iterate over the members of u.
func (u *Unicode) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c(u))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

// Size returns the size of u. The equivalent Python is "len(u)".
func (u *Unicode) Size() int {
	ret := C.PyObject_Size(c(u))
	if ret < 0 {
		clearErr();
		return 0
	}
	return int(ret)
}

// AsMappingMethods returns a MappingMethods instance that refers to the same
// underlying Python object as u.
//
// This method also means that Unicode implements the MappingProtocol interface.
func (u *Unicode) AsMappingMethods() *MappingMethods {
	return (*MappingMethods)(unsafe.Pointer(u.Base()))
}

func (u *Unicode) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(c(u), cKey)
	return obj2ObjErr(ret)
}

func (u *Unicode) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(c(u), c(key))
	clearErr();
	return ret > 0
}

func (u *Unicode) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(c(u), cKey)
	clearErr();
	return ret > 0
}

// AsSequenceMethods returns a SequenceMethods instance that refers to the same
// underlying Python object as u.
//
// This method also means that Unicode implements the SequenceProtocol
// interface.
func (u *Unicode) AsSequenceMethods() *SequenceMethods {
	return (*SequenceMethods)(unsafe.Pointer(u.Base()))
}

func (u *Unicode) GetIndex(idx int) (Object, error) {
	ret := C.PySequence_GetItem(c(u), C.Py_ssize_t(idx))
	return obj2ObjErr(ret)
}

func (u *Unicode) Count(obj Object) (int, error) {
	ret := C.PySequence_Count(c(u), c(obj))
	return ssize_t2IntErr(ret)
}

func (u *Unicode) Index(obj Object) (int, error) {
	ret := C.PySequence_Index(c(u), c(obj))
	return ssize_t2IntErr(ret)
}

func (u *Unicode) List() (*List, error) {
	ret := C.PySequence_List(c(u))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

func (u *Unicode) Tuple() (*Tuple, error) {
	ret := C.PySequence_Tuple(c(u))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

func (u *Unicode) Concat(obj Object) (Object, error) {
	ret := C.PySequence_Concat(c(u), c(obj))
	return obj2ObjErr(ret)
}

func (u *Unicode) Repeat(count int) (Object, error) {
	ret := C.PySequence_Repeat(c(u), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (u *Unicode) Contains(obj Object) (bool, error) {
	ret := C.PySequence_Contains(c(u), c(obj))
	return int2BoolErr(ret)
}

func (u *Unicode) GetSlice(start, end int) (Object, error) {
	ret := C.PySequence_GetSlice(c(u), C.Py_ssize_t(start), C.Py_ssize_t(end))
	return obj2ObjErr(ret)
}

// Remainder returns the remainder of dividing u by obj. The equivalent Python
// is "u % obj".
//
// Return value: New Reference.
func (u *Unicode) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(c(u), c(obj))
	return obj2ObjErr(ret)
}



/*
set fields:
  mp_length
  mp_subscript
  nb_remainder
  ob_base
  sq_concat
  sq_contains
  sq_item
  sq_length
  sq_repeat
  tp_as_mapping
  tp_as_number
  tp_as_sequence
  tp_basicsize
  tp_dealloc
  tp_doc
  tp_flags
  tp_free
  tp_getattro
  tp_hash
  tp_iter
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_str
*/

