// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Dict represents objects of the DictType (or PyDictType
// in the Python API) type.
//
// This type implements the Mapping protocol.
type Dict struct {
	abstractObject
	o C.PyDictObject
}

// DictType is the Type object that represents the Dict type.
var DictType = (*Type)(unsafe.Pointer(&C.PyDict_Type))

func dictCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.dictCheck(c(obj)) != 0
}

func newDict(obj *C.PyObject) *Dict {
	return (*Dict)(unsafe.Pointer(obj))
}

func (d *Dict) Size() int {
	ret := C.PyObject_Size(c(d))
	if ret < 0 {
		clearErr();
		return 0
	}
	return int(ret)
}

func (d *Dict) AsMapping() *MappingMethods {
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

func (d *Dict) Or(obj Object) (Object, error) {
	ret := C.PyNumber_Or(c(d), c(obj))
	return obj2ObjErr(ret)
}

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

