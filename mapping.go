package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Mapping is an interface that is implemented by types that implement the
// Python "Mapping Protocol".
type Mapping interface {
	Object
	AsMapping() *MappingMethods
}

// AsMapping returns a MappingMethods instance that refers to the same
// underlying Python object as obj. If obj doesn't implement the "Mapping
// Protocol" (i.e. the Mapping interface), then nil is returned.
func AsMapping(obj Object) *MappingMethods {
	if n, ok := obj.(Mapping); ok {
		return n.AsMapping()
	}
	if C.mappingCheck(c(obj)) > 0 {
		return (*MappingMethods)(unsafe.Pointer(obj.Base()))
	}
	return nil
}

func (m *MappingMethods) Size() (int, error) {
	ret := C.PyMapping_Size(c(m))
	return ssize_t2IntErr(ret)
}

func (m *MappingMethods) DelItemString(key string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyObject_DelItemString(c(m), cKey)
	return int2Err(ret)
}

func (m *MappingMethods) DelItem(key Object) error {
	ret := C.PyObject_DelItem(c(m), c(key))
	return int2Err(ret)
}

func (m *MappingMethods) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(c(m), cKey)
	return ret > 0
}

func (m *MappingMethods) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(c(m), c(key))
	return ret > 0
}

func (m *MappingMethods) Keys() (Object, error) {
	ret := C.PyMapping_Keys_(c(m))
	return obj2ObjErr(ret)
}

func (m *MappingMethods) Values() (Object, error) {
	ret := C.PyMapping_Values_(c(m))
	return obj2ObjErr(ret)
}

func (m *MappingMethods) Items() (Object, error) {
	ret := C.PyMapping_Items_(c(m))
	return obj2ObjErr(ret)
}

func (m *MappingMethods) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(c(m), cKey)
	return obj2ObjErr(ret)
}

func (m *MappingMethods) SetItemString(key string, v Object) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_SetItemString(c(m), cKey, c(v))
	return int2Err(ret)
}
