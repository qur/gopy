// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// ByteArray represents objects of the ByteArrayType (or PyByteArray_Type
// in the Python API) type.
//
// This type implements the Mapping protocol.
//
// This type implements the Sequence protocol.
type ByteArray struct {
	o C.PyByteArrayObject
}

var _ Object = (*ByteArray)(nil)

// ByteArrayType is the Type object that represents the ByteArray type.
var ByteArrayType = newType(&C.PyByteArray_Type)

func byteArrayCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.byteArrayCheck(c(obj)) != 0
}

// AsByteArray casts the given obj to a ByteArray (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a ByteArray, then nil is returned.
//
// Return value: Borrowed Reference.
func AsByteArray(obj Object) *ByteArray {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.byteArrayCheck(o) == 0 {
		return nil
	}
	return (*ByteArray)(unsafe.Pointer(o))
}

func newByteArray(obj *C.PyObject) *ByteArray {
	return (*ByteArray)(unsafe.Pointer(obj))
}



func (b *ByteArray) c() *C.PyByteArrayObject {
	return (*C.PyByteArrayObject)(unsafe.Pointer(b))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (b *ByteArray) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(b))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (b *ByteArray) Type() *Type {
	return newType(c(b).ob_type)
}


// Decref decrements b's reference count, b may not be nil.
func (b *ByteArray) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(b))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments b's reference count, b may not be nil.
func (b *ByteArray) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(b)).anon0[0]))
	*refcnt++
}

// Repr returns a String representation of "b". This is equivalent to the
// Python "repr(b)".
//
// Return value: New Reference.
func (b *ByteArray) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Str returns a String representation of "b". This is equivalent to the
// Python "str(b)".
//
// Return value: New Reference.
func (b *ByteArray) Str() (*Unicode, error) {
	ret := C.PyObject_Str(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// HasAttr returns true if "b" has the attribute "name".  This is equivalent
// to the Python "hasattr(b, name)".
func (b *ByteArray) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(b), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "b" with the name "name".  This is
// equivalent to the Python "b.name".
//
// Return value: New Reference.
func (b *ByteArray) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(b), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "b" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "b op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (b *ByteArray) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(b), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "b" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "b op obj",
// where op is the corresponding Python operator for op.
func (b *ByteArray) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(b), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// Iter returns an Iterator that will iterate over the members of b.
func (b *ByteArray) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

// Size returns the size of b. The equivalent Python is "len(b)".
func (b *ByteArray) Size() int {
	ret := C.PyObject_Size(c(b))
	if ret < 0 {
		clearErr();
		return 0
	}
	return int(ret)
}

// AsMappingMethods returns a MappingMethods instance that refers to the same
// underlying Python object as b.
//
// This method also means that ByteArray implements the MappingProtocol interface.
func (b *ByteArray) AsMappingMethods() *MappingMethods {
	return (*MappingMethods)(unsafe.Pointer(b.Base()))
}

func (b *ByteArray) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(c(b), cKey)
	return obj2ObjErr(ret)
}

func (b *ByteArray) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(c(b), c(key))
	clearErr();
	return ret > 0
}

func (b *ByteArray) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(c(b), cKey)
	clearErr();
	return ret > 0
}

func (b *ByteArray) DelItem(key Object) error {
	ret := C.PyObject_DelItem(c(b), c(key))
	return int2Err(ret)
}

func (b *ByteArray) DelItemString(key string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyObject_DelItemString(c(b), cKey)
	return int2Err(ret)
}

func (b *ByteArray) SetItemString(key string, v Object) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_SetItemString(c(b), cKey, c(v))
	return int2Err(ret)
}

// AsSequenceMethods returns a SequenceMethods instance that refers to the same
// underlying Python object as b.
//
// This method also means that ByteArray implements the SequenceProtocol
// interface.
func (b *ByteArray) AsSequenceMethods() *SequenceMethods {
	return (*SequenceMethods)(unsafe.Pointer(b.Base()))
}

func (b *ByteArray) GetIndex(idx int) (Object, error) {
	ret := C.PySequence_GetItem(c(b), C.Py_ssize_t(idx))
	return obj2ObjErr(ret)
}

func (b *ByteArray) Count(obj Object) (int, error) {
	ret := C.PySequence_Count(c(b), c(obj))
	return ssize_t2IntErr(ret)
}

func (b *ByteArray) Index(obj Object) (int, error) {
	ret := C.PySequence_Index(c(b), c(obj))
	return ssize_t2IntErr(ret)
}

func (b *ByteArray) List() (*List, error) {
	ret := C.PySequence_List(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

func (b *ByteArray) Tuple() (*Tuple, error) {
	ret := C.PySequence_Tuple(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

func (b *ByteArray) SetIndex(idx int, obj Object) error {
	ret := C.PySequence_SetItem(c(b), C.Py_ssize_t(idx), c(obj))
	return int2Err(ret)
}

func (b *ByteArray) DelIndex(idx int) error {
	ret := C.PySequence_DelItem(c(b), C.Py_ssize_t(idx))
	return int2Err(ret)
}

func (b *ByteArray) Concat(obj Object) (Object, error) {
	ret := C.PySequence_Concat(c(b), c(obj))
	return obj2ObjErr(ret)
}

func (b *ByteArray) InPlaceConcat(obj Object) (Object, error) {
	ret := C.PySequence_InPlaceConcat(c(b), c(obj))
	return obj2ObjErr(ret)
}

func (b *ByteArray) Repeat(count int) (Object, error) {
	ret := C.PySequence_Repeat(c(b), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (b *ByteArray) InPlaceRepeat(count int) (Object, error) {
	ret := C.PySequence_InPlaceRepeat(c(b), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (b *ByteArray) Contains(obj Object) (bool, error) {
	ret := C.PySequence_Contains(c(b), c(obj))
	return int2BoolErr(ret)
}

func (b *ByteArray) GetSlice(start, end int) (Object, error) {
	ret := C.PySequence_GetSlice(c(b), C.Py_ssize_t(start), C.Py_ssize_t(end))
	return obj2ObjErr(ret)
}

func (b *ByteArray) SetSlice(start, end int, obj Object) error {
	ret := C.PySequence_SetSlice(c(b), C.Py_ssize_t(start), C.Py_ssize_t(end), c(obj))
	return int2Err(ret)
}

func (b *ByteArray) DelSlice(start, end int) error {
	ret := C.PySequence_DelSlice(c(b), C.Py_ssize_t(start), C.Py_ssize_t(end))
	return int2Err(ret)
}

// Remainder returns the remainder of dividing b by obj. The equivalent Python
// is "b % obj".
//
// Return value: New Reference.
func (b *ByteArray) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(c(b), c(obj))
	return obj2ObjErr(ret)
}

// AsBufferMethods returns a BufferMethods instance that refers to the same
// underlying Python object as b.
//
// This method also means that ByteArray implements the BufferProtocol
// interface.
func (b *ByteArray) AsBufferMethods() *BufferMethods {
	return (*BufferMethods)(unsafe.Pointer(b.Base()))
}

func(b *ByteArray) GetBuffer(flags BufferFlags) (*Buffer, error) {
	buf := newBuffer()
	ret := C.PyObject_GetBuffer(c(b), buf.c(), C.int(flags))
	if ret < 0 {
		return nil, exception()
	}
	return buf, nil
}



/*
set fields:
  bf_getbuffer
  bf_releasebuffer
  mp_ass_subscript
  mp_length
  mp_subscript
  nb_remainder
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
  tp_as_buffer
  tp_as_mapping
  tp_as_number
  tp_as_sequence
  tp_basicsize
  tp_dealloc
  tp_doc
  tp_flags
  tp_free
  tp_getattro
  tp_init
  tp_iter
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_str
*/

