// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Bytes represents objects of the BytesType (or PyBytesType
// in the Python API) type.
//
// This type implements the Mapping protocol.
//
// This type implements the Sequence protocol.
type Bytes struct {
	abstractObject
	o C.PyBytesObject
}

// BytesType is the Type object that represents the Bytes type.
var BytesType = (*Type)(unsafe.Pointer(&C.PyBytes_Type))

func bytesCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.bytesCheck(c(obj)) != 0
}

func newBytes(obj *C.PyObject) *Bytes {
	return (*Bytes)(unsafe.Pointer(obj))
}

func (b *Bytes) Size() int {
	ret := C.PyObject_Size(c(b))
	if ret < 0 {
		clearErr();
		return 0
	}
	return int(ret)
}

func (b *Bytes) AsMapping() *MappingMethods {
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

func (b *Bytes) AsSequence() *SequenceMethods {
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

func (b *Bytes) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(c(b), c(obj))
	return obj2ObjErr(ret)
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

