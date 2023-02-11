// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

// SequenceProtocol is a 0-sized type that can be embedded in concrete types
// after the AbstractObject to provide access to the suite of methods that
// Python calls the "Sequence Protocol".
type SequenceProtocol struct{}

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
// implements the "Sequence Protocol" but doesn't embed SequenceProtocol can be
// turned into a Sequence by calling AsSequence.
type sequence struct {
	AbstractObject
	SequenceProtocol
	o C.PyObject
}

func csp(s *SequenceProtocol) *C.PyObject {
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

func (s *SequenceProtocol) Size() (int64, error) {
	ret := C.PySequence_Size(csp(s))
	return ssize_t2Int64Err(ret)
}

func (s *SequenceProtocol) Length() (int64, error) {
	ret := C.PySequence_Length(csp(s))
	return ssize_t2Int64Err(ret)
}

func (s *SequenceProtocol) Concat(obj Object) (Object, error) {
	ret := C.PySequence_Concat(csp(s), c(obj))
	return obj2ObjErr(ret)
}

func (s *SequenceProtocol) Repeat(count int64) (Object, error) {
	ret := C.PySequence_Repeat(csp(s), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (s *SequenceProtocol) InPlaceConcat(obj Object) (Object, error) {
	ret := C.PySequence_InPlaceConcat(csp(s), c(obj))
	return obj2ObjErr(ret)
}

func (s *SequenceProtocol) InPlaceRepeat(count int64) (Object, error) {
	ret := C.PySequence_InPlaceRepeat(csp(s), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (s *SequenceProtocol) GetItem(i int64) (Object, error) {
	ret := C.PySequence_GetItem(csp(s), C.Py_ssize_t(i))
	return obj2ObjErr(ret)
}

func (s *SequenceProtocol) GetSlice(i1, i2 int64) (Object, error) {
	ret := C.PySequence_GetSlice(csp(s), C.Py_ssize_t(i1), C.Py_ssize_t(i2))
	return obj2ObjErr(ret)
}

func (s *SequenceProtocol) SetItem(i int64, v Object) error {
	ret := C.PySequence_SetItem(csp(s), C.Py_ssize_t(i), c(v))
	return int2Err(ret)
}

func (s *SequenceProtocol) DelItem(i int64) error {
	ret := C.PySequence_DelItem(csp(s), C.Py_ssize_t(i))
	return int2Err(ret)
}

func (s *SequenceProtocol) SetSlice(i1, i2 int64, v Object) error {
	ret := C.PySequence_SetSlice(csp(s), C.Py_ssize_t(i1), C.Py_ssize_t(i2), c(v))
	return int2Err(ret)
}

func (s *SequenceProtocol) DelSlice(i1, i2 int64) error {
	ret := C.PySequence_DelSlice(csp(s), C.Py_ssize_t(i1), C.Py_ssize_t(i2))
	return int2Err(ret)
}

func (s *SequenceProtocol) Count(value Object) (int64, error) {
	ret := C.PySequence_Count(csp(s), c(value))
	return ssize_t2Int64Err(ret)
}

func (s *SequenceProtocol) Contains(value Object) (bool, error) {
	ret := C.PySequence_Contains(csp(s), c(value))
	return int2BoolErr(ret)
}

func (s *SequenceProtocol) Index(value Object) (int64, error) {
	ret := C.PySequence_Index(csp(s), c(value))
	return ssize_t2Int64Err(ret)
}

func (s *SequenceProtocol) List() (Object, error) {
	ret := C.PySequence_List(csp(s))
	return obj2ObjErr(ret)
}

func (s *SequenceProtocol) Tuple() (Object, error) {
	ret := C.PySequence_Tuple(csp(s))
	return obj2ObjErr(ret)
}
