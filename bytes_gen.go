// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Bytes represents objects of the BytesType (or PyBytes_Type
// in the Python API) type.
//
// This type implements the Mapping protocol.
//
// This type implements the Sequence protocol.
type Bytes struct {
	o C.PyBytesObject
}

var _ Object = (*Bytes)(nil)

// BytesType is the Type object that represents the Bytes type.
var BytesType = newType(&C.PyBytes_Type)

func bytesCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.bytesCheck(c(obj)) != 0
}

// AsBytes casts the given obj to a Bytes (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Bytes, then nil is returned.
//
// Return value: Borrowed Reference.
func AsBytes(obj Object) *Bytes {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.bytesCheck(o) == 0 {
		return nil
	}
	return (*Bytes)(unsafe.Pointer(o))
}

func newBytes(obj *C.PyObject) *Bytes {
	return (*Bytes)(unsafe.Pointer(obj))
}



func (b *Bytes) c() *C.PyBytesObject {
	return (*C.PyBytesObject)(unsafe.Pointer(b))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (b *Bytes) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(b))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (b *Bytes) Type() *Type {
	return newType(c(b).ob_type)
}


// Decref decrements b's reference count, b may not be nil.
func (b *Bytes) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(b))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments b's reference count, b may not be nil.
func (b *Bytes) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(b)).anon0[0]))
	*refcnt++
}

// Repr returns a String representation of "b". This is equivalent to the
// Python "repr(b)".
//
// Return value: New Reference.
func (b *Bytes) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of b. The equivalent
// Python is "hash(b)".
func (b *Bytes) Hash() (int, error) {
	ret := C.PyObject_Hash(c(b))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// Str returns a String representation of "b". This is equivalent to the
// Python "str(b)".
//
// Return value: New Reference.
func (b *Bytes) Str() (*Unicode, error) {
	ret := C.PyObject_Str(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// HasAttr returns true if "b" has the attribute "name".  This is equivalent
// to the Python "hasattr(b, name)".
func (b *Bytes) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(b), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "b" with the name "name".  This is
// equivalent to the Python "b.name".
//
// Return value: New Reference.
func (b *Bytes) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(b), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "b" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "b op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (b *Bytes) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(b), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "b" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "b op obj",
// where op is the corresponding Python operator for op.
func (b *Bytes) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(b), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// Iter returns an Iterator that will iterate over the members of b.
func (b *Bytes) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

// Size returns the size of b. The equivalent Python is "len(b)".
func (b *Bytes) Size() int {
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
// This method also means that Bytes implements the MappingProtocol interface.
func (b *Bytes) AsMappingMethods() *MappingMethods {
	return (*MappingMethods)(unsafe.Pointer(b.Base()))
}

func (b *Bytes) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(c(b), cKey)
	return obj2ObjErr(ret)
}

func (b *Bytes) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(c(b), c(key))
	clearErr();
	return ret > 0
}

func (b *Bytes) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(c(b), cKey)
	clearErr();
	return ret > 0
}

// AsSequenceMethods returns a SequenceMethods instance that refers to the same
// underlying Python object as b.
//
// This method also means that Bytes implements the SequenceProtocol
// interface.
func (b *Bytes) AsSequenceMethods() *SequenceMethods {
	return (*SequenceMethods)(unsafe.Pointer(b.Base()))
}

func (b *Bytes) GetIndex(idx int) (Object, error) {
	ret := C.PySequence_GetItem(c(b), C.Py_ssize_t(idx))
	return obj2ObjErr(ret)
}

func (b *Bytes) Count(obj Object) (int, error) {
	ret := C.PySequence_Count(c(b), c(obj))
	return ssize_t2IntErr(ret)
}

func (b *Bytes) Index(obj Object) (int, error) {
	ret := C.PySequence_Index(c(b), c(obj))
	return ssize_t2IntErr(ret)
}

func (b *Bytes) List() (*List, error) {
	ret := C.PySequence_List(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

func (b *Bytes) Tuple() (*Tuple, error) {
	ret := C.PySequence_Tuple(c(b))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

func (b *Bytes) Concat(obj Object) (Object, error) {
	ret := C.PySequence_Concat(c(b), c(obj))
	return obj2ObjErr(ret)
}

func (b *Bytes) Repeat(count int) (Object, error) {
	ret := C.PySequence_Repeat(c(b), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (b *Bytes) Contains(obj Object) (bool, error) {
	ret := C.PySequence_Contains(c(b), c(obj))
	return int2BoolErr(ret)
}

func (b *Bytes) GetSlice(start, end int) (Object, error) {
	ret := C.PySequence_GetSlice(c(b), C.Py_ssize_t(start), C.Py_ssize_t(end))
	return obj2ObjErr(ret)
}

// Remainder returns the remainder of dividing b by obj. The equivalent Python
// is "b % obj".
//
// Return value: New Reference.
func (b *Bytes) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(c(b), c(obj))
	return obj2ObjErr(ret)
}

// AsBufferMethods returns a BufferMethods instance that refers to the same
// underlying Python object as b.
//
// This method also means that Bytes implements the BufferProtocol
// interface.
func (b *Bytes) AsBufferMethods() *BufferMethods {
	return (*BufferMethods)(unsafe.Pointer(b.Base()))
}

func(b *Bytes) GetBuffer(flags BufferFlags) (*Buffer, error) {
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
  mp_length
  mp_subscript
  nb_remainder
  ob_base
  sq_concat
  sq_contains
  sq_item
  sq_length
  sq_repeat
  tp_alloc
  tp_as_buffer
  tp_as_mapping
  tp_as_number
  tp_as_sequence
  tp_basicsize
  tp_doc
  tp_flags
  tp_free
  tp_getattro
  tp_hash
  tp_itemsize
  tp_iter
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_str
*/

