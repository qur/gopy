// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// MemoryView represents objects of the MemoryViewType (or PyMemoryViewType
// in the Python API) type.
//
// This type implements the Mapping protocol.
//
// This type implements the Sequence protocol.
type MemoryView struct {
	abstractObject
	o C.PyMemoryViewObject
}

// MemoryViewType is the Type object that represents the MemoryView type.
var MemoryViewType = (*Type)(unsafe.Pointer(&C.PyMemoryView_Type))

func memoryViewCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.memoryViewCheck(c(obj)) != 0
}

func newMemoryView(obj *C.PyObject) *MemoryView {
	return (*MemoryView)(unsafe.Pointer(obj))
}

func (m *MemoryView) Size() int {
	ret := C.PyObject_Size(c(m))
	if ret < 0 {
		clearErr();
		return 0
	}
	return int(ret)
}

func (m *MemoryView) AsMapping() *MappingMethods {
	return (*MappingMethods)(unsafe.Pointer(m.Base()))
}

func (m *MemoryView) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(c(m), cKey)
	return obj2ObjErr(ret)
}

func (m *MemoryView) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(c(m), c(key))
	clearErr();
	return ret > 0
}

func (m *MemoryView) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(c(m), cKey)
	clearErr();
	return ret > 0
}

func (m *MemoryView) DelItem(key Object) error {
	ret := C.PyObject_DelItem(c(m), c(key))
	return int2Err(ret)
}

func (m *MemoryView) DelItemString(key string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyObject_DelItemString(c(m), cKey)
	return int2Err(ret)
}

func (m *MemoryView) SetItemString(key string, v Object) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_SetItemString(c(m), cKey, c(v))
	return int2Err(ret)
}

func (m *MemoryView) AsSequence() *SequenceMethods {
	return (*SequenceMethods)(unsafe.Pointer(m.Base()))
}

func (m *MemoryView) GetIndex(idx int) (Object, error) {
	ret := C.PySequence_GetItem(c(m), C.Py_ssize_t(idx))
	return obj2ObjErr(ret)
}



/*
set fields:
  bf_getbuffer
  bf_releasebuffer
  mp_ass_subscript
  mp_length
  mp_subscript
  ob_base
  sq_item
  sq_length
  tp_as_buffer
  tp_as_mapping
  tp_as_sequence
  tp_basicsize
  tp_clear
  tp_dealloc
  tp_doc
  tp_flags
  tp_getattro
  tp_getset
  tp_hash
  tp_itemsize
  tp_iter
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_traverse
  tp_weaklistoffset
*/

