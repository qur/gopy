// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

// sequenceProtocol is a 0-sized type that can be embedded in concrete types
// after the AbstractObject to provide access to the suite of methods that
// Python calls the "Sequence Protocol".
type sequenceProtocol struct{}

// Sequence is an interface that defines the Python "Sequence Protocol".
type Sequence interface {
	Object
	Size() (int64, error)
	Length() (int64, error)
	Concat(obj Object) (Object, error)
	Repeat(count int64) (Object, error)
	InPlaceConcat(obj Object) (Object, error)
	InPlaceRepeat(count int64) (Object, error)
	GetItem(i int64) (Object, error)
	GetSlice(i1, i2 int64) (Object, error)
	SetItem(i int64, v Object) error
	DelItem(i int64) error
	SetSlice(i1, i2 int64, v Object) error
	DelSlice(i1, i2 int64) error
	Count(value Object) (int64, error)
	Contains(value Object) (bool, error)
	Index(value Object) (int64, error)
	List() (Object, error)
	Tuple() (Object, error)
}

// sequence is a concrete realisation of the Sequence Protocol.  A type that
// implements the "Sequence Protocol" but doesn't implement Sequence can be
// turned into a Sequence by calling AsSequence.
type sequence struct {
	abstractObject
	sequenceProtocol
	o C.PyObject
}

func csp(s *sequenceProtocol) *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(s))
}

// AsSequence returns a struct pointer that satisfies the Sequence interface.
// It will refer to the same underlying object as obj.  If obj doesn't implement
// the "Sequence Protocol", then nil is returned.
func AsSequence(obj Object) Sequence {
	if C.PySequence_Check(c(obj)) != 1 {
		return nil
	}
	if n, ok := obj.(Sequence); ok {
		return n
	}
	return (*sequence)(unsafe.Pointer(obj.Base()))
}

func (s *sequenceProtocol) Size() (int64, error) {
	ret := C.PySequence_Size(csp(s))
	return ssize_t2Int64Err(ret)
}

func (s *sequenceProtocol) Length() (int64, error) {
	ret := C.PySequence_Length(csp(s))
	return ssize_t2Int64Err(ret)
}

func (s *sequenceProtocol) Concat(obj Object) (Object, error) {
	ret := C.PySequence_Concat(csp(s), c(obj))
	return obj2ObjErr(ret)
}

func (s *sequenceProtocol) Repeat(count int64) (Object, error) {
	ret := C.PySequence_Repeat(csp(s), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (s *sequenceProtocol) InPlaceConcat(obj Object) (Object, error) {
	ret := C.PySequence_InPlaceConcat(csp(s), c(obj))
	return obj2ObjErr(ret)
}

func (s *sequenceProtocol) InPlaceRepeat(count int64) (Object, error) {
	ret := C.PySequence_InPlaceRepeat(csp(s), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (s *sequenceProtocol) GetItem(i int64) (Object, error) {
	ret := C.PySequence_GetItem(csp(s), C.Py_ssize_t(i))
	return obj2ObjErr(ret)
}

func (s *sequenceProtocol) GetSlice(i1, i2 int64) (Object, error) {
	ret := C.PySequence_GetSlice(csp(s), C.Py_ssize_t(i1), C.Py_ssize_t(i2))
	return obj2ObjErr(ret)
}

func (s *sequenceProtocol) SetItem(i int64, v Object) error {
	ret := C.PySequence_SetItem(csp(s), C.Py_ssize_t(i), c(v))
	return int2Err(ret)
}

func (s *sequenceProtocol) DelItem(i int64) error {
	ret := C.PySequence_DelItem(csp(s), C.Py_ssize_t(i))
	return int2Err(ret)
}

func (s *sequenceProtocol) SetSlice(i1, i2 int64, v Object) error {
	ret := C.PySequence_SetSlice(csp(s), C.Py_ssize_t(i1), C.Py_ssize_t(i2), c(v))
	return int2Err(ret)
}

func (s *sequenceProtocol) DelSlice(i1, i2 int64) error {
	ret := C.PySequence_DelSlice(csp(s), C.Py_ssize_t(i1), C.Py_ssize_t(i2))
	return int2Err(ret)
}

func (s *sequenceProtocol) Count(value Object) (int64, error) {
	ret := C.PySequence_Count(csp(s), c(value))
	return ssize_t2Int64Err(ret)
}

func (s *sequenceProtocol) Contains(value Object) (bool, error) {
	ret := C.PySequence_Contains(csp(s), c(value))
	return int2BoolErr(ret)
}

func (s *sequenceProtocol) Index(value Object) (int64, error) {
	ret := C.PySequence_Index(csp(s), c(value))
	return ssize_t2Int64Err(ret)
}

func (s *sequenceProtocol) List() (Object, error) {
	ret := C.PySequence_List(csp(s))
	return obj2ObjErr(ret)
}

func (s *sequenceProtocol) Tuple() (Object, error) {
	ret := C.PySequence_Tuple(csp(s))
	return obj2ObjErr(ret)
}
