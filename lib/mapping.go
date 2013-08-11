// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

// MappingProtocol is a 0-sized type that can be embedded in concrete types
// after the AbstractObject to provide access to the suite of methods that
// Python calls the "Mapping Protocol".
type MappingProtocol struct{}

// Mapping is an interface that defines the Python "Mapping Protocol".
type Mapping interface {
	Object
	Size() (int64, error)
	Length() (int64, error)
	DelItemString(key string) error
	DelItem(key Object) error
	HasKeyString(key string) bool
	HasKey(key Object) bool
	Keys() (Object, error)
	Values() (Object, error)
	Items() (Object, error)
	GetItemString(key string) (Object, error)
	SetItemString(key string, v Object) error
}

// mapping is a concrete realisation of the Mapping Protocol.  A type that
// implements the "Mapping Protocol" but doesn't embed MappingProtocol can be
// turned into a Mapping by calling AsMapping.
type mapping struct {
	AbstractObject
	MappingProtocol
	o C.PyObject
}

func cmp(m *MappingProtocol) *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(m))
}

// AsMapping returns a struct pointer that satisfies the Mapping interface.
// It will refer to the same underlying object as obj.  If obj doesn't implement
// the "Mapping Protocol", then nil is returned.
func AsMapping(obj Object) Mapping {
	if C.PyMapping_Check(c(obj)) != 1 {
		return nil
	}
	if n, ok := obj.(Mapping); ok {
		return n
	}
	return (*mapping)(unsafe.Pointer(obj.Base()))
}

func (m *MappingProtocol) Size() (int64, error) {
	ret := C.PyMapping_Size(cmp(m))
	return ssize_t2Int64Err(ret)
}

func (m *MappingProtocol) Length() (int64, error) {
	ret := C.PyMapping_Length(cmp(m))
	return ssize_t2Int64Err(ret)
}

func (m *MappingProtocol) DelItemString(key string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyObject_DelItemString(cmp(m), cKey)
	return int2Err(ret)
}

func (m *MappingProtocol) DelItem(key Object) error {
	ret := C.PyObject_DelItem(cmp(m), c(key))
	return int2Err(ret)
}

func (m *MappingProtocol) HasKeyString(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_HasKeyString(cmp(m), cKey)
	return ret > 0
}

func (m *MappingProtocol) HasKey(key Object) bool {
	ret := C.PyMapping_HasKey(cmp(m), c(key))
	return ret > 0
}

func (m *MappingProtocol) Keys() (Object, error) {
	ret := C.PyMapping_Keys_(cmp(m))
	return obj2ObjErr(ret)
}

func (m *MappingProtocol) Values() (Object, error) {
	ret := C.PyMapping_Values_(cmp(m))
	return obj2ObjErr(ret)
}

func (m *MappingProtocol) Items() (Object, error) {
	ret := C.PyMapping_Items_(cmp(m))
	return obj2ObjErr(ret)
}

func (m *MappingProtocol) GetItemString(key string) (Object, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_GetItemString(cmp(m), cKey)
	return obj2ObjErr(ret)
}

func (m *MappingProtocol) SetItemString(key string, v Object) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ret := C.PyMapping_SetItemString(cmp(m), cKey, c(v))
	return int2Err(ret)
}
