package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Sequence is an interface that is implemented by types that implement the
// Python "Sequence Protocol".
type Sequence interface {
	Object
	AsSequence() *SequenceMethods
}

// MappingMethods is a concrete realisation of the full set of Sequence Protocol
// methods.  A type that implements the "Sequence Protocol" can be turned into a
// SequenceMethods instance using AsSequence.
//
// Note that the methods not already implemented on the type itself may return
// an error, as not all methods are implemented by all types that support the
// protocol.
type SequenceMethods struct {
	abstractObject
	o C.PyObject
}

// AsSequence returns a SequenceMethods instance that refers to the same
// underlying Python object as obj. If obj doesn't implement the "Sequence
// Protocol" (i.e. the Sequence interface), then nil is returned.
func AsSequence(obj Object) *SequenceMethods {
	if n, ok := obj.(Sequence); ok {
		return n.AsSequence()
	}
	if C.sequenceCheck(c(obj)) > 0 {
		return (*SequenceMethods)(unsafe.Pointer(obj.Base()))
	}
	return nil
}

func (s *SequenceMethods) Size() (int, error) {
	ret := C.PySequence_Size(c(s))
	return ssize_t2IntErr(ret)
}

func (s *SequenceMethods) Concat(obj Object) (Object, error) {
	ret := C.PySequence_Concat(c(s), c(obj))
	return obj2ObjErr(ret)
}

func (s *SequenceMethods) Repeat(count int) (Object, error) {
	ret := C.PySequence_Repeat(c(s), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (s *SequenceMethods) InPlaceConcat(obj Object) (Object, error) {
	ret := C.PySequence_InPlaceConcat(c(s), c(obj))
	return obj2ObjErr(ret)
}

func (s *SequenceMethods) InPlaceRepeat(count int) (Object, error) {
	ret := C.PySequence_InPlaceRepeat(c(s), C.Py_ssize_t(count))
	return obj2ObjErr(ret)
}

func (s *SequenceMethods) GetItem(i int) (Object, error) {
	ret := C.PySequence_GetItem(c(s), C.Py_ssize_t(i))
	return obj2ObjErr(ret)
}

func (s *SequenceMethods) GetSlice(i1, i2 int) (Object, error) {
	ret := C.PySequence_GetSlice(c(s), C.Py_ssize_t(i1), C.Py_ssize_t(i2))
	return obj2ObjErr(ret)
}

func (s *SequenceMethods) SetItem(i int, v Object) error {
	ret := C.PySequence_SetItem(c(s), C.Py_ssize_t(i), c(v))
	return int2Err(ret)
}

func (s *SequenceMethods) DelItem(i int) error {
	ret := C.PySequence_DelItem(c(s), C.Py_ssize_t(i))
	return int2Err(ret)
}

func (s *SequenceMethods) SetSlice(i1, i2 int, v Object) error {
	ret := C.PySequence_SetSlice(c(s), C.Py_ssize_t(i1), C.Py_ssize_t(i2), c(v))
	return int2Err(ret)
}

func (s *SequenceMethods) DelSlice(i1, i2 int) error {
	ret := C.PySequence_DelSlice(c(s), C.Py_ssize_t(i1), C.Py_ssize_t(i2))
	return int2Err(ret)
}

func (s *SequenceMethods) Count(value Object) (int, error) {
	ret := C.PySequence_Count(c(s), c(value))
	return ssize_t2IntErr(ret)
}

func (s *SequenceMethods) Contains(value Object) (bool, error) {
	ret := C.PySequence_Contains(c(s), c(value))
	return int2BoolErr(ret)
}

func (s *SequenceMethods) Index(value Object) (int, error) {
	ret := C.PySequence_Index(c(s), c(value))
	return ssize_t2IntErr(ret)
}

func (s *SequenceMethods) List() (*List, error) {
	ret := C.PySequence_List(c(s))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

func (s *SequenceMethods) Tuple() (*Tuple, error) {
	ret := C.PySequence_Tuple(c(s))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}
