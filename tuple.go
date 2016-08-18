// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int tupleCheck(PyObject *o) { return PyTuple_Check(o); }
// static inline int tupleCheckE(PyObject *o) { return PyTuple_CheckExact(o); }
// static inline size_t tupleItemSize(void) { return sizeof(PyObject *); }
import "C"

import (
	"fmt"
	"unsafe"
)

type Tuple struct {
	AbstractObject
	o *C.PyTupleObject
}

var tupleObjMap = make(map[*C.PyObject]*Tuple)

// TupleType is the Type object that represents the Tuple type.
var TupleType = newType((*C.PyObject)(unsafe.Pointer(C.getBasePyType(C.GoPyTuple_Type))))

func tupleCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.tupleCheck(c(obj)) != 0
}

func newTuple(obj *C.PyObject) *Tuple {
	if t, ok := tupleObjMap[obj]; ok {
		return t
	}
	t := &Tuple{o: (*C.PyTupleObject)(unsafe.Pointer(obj))}
	tupleObjMap[obj] = t
	return t
}

// NewTuple returns a new *Tuple of the specified size.  However the entries are
// all set to NULL, so the tuple should not be shared, especially with Python
// code, until the entries have all been set.
//
// Return value: New Reference.
func NewTuple(size int64) (*Tuple, error) {
	ret := C.PyTuple_New(C.Py_ssize_t(size))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

// PackTuple returns a new *Tuple which contains the arguments.  This tuple is
// ready to use.
//
// Return value: New Reference.
func PackTuple(items ...Object) (*Tuple, error) {
	ret := C.PyTuple_New(C.Py_ssize_t(len(items)))
	if ret == nil {
		return nil, exception()
	}

	// Since the ob_item array has a size of 1, Go won't let us index more than
	// a single entry, and if we try and use our own local type definition with
	// a flexible array member then cgo converts it to [0]byte which is even
	// less useful.  So, we resort to pointer manipulation - which is
	// unfortunate, as it's messy in Go.

	// base is a pointer to the first item in the array of PyObject pointers.
	// step is the size of a PyObject * (i.e. the number of bytes we need to add
	// to get to the next item).
	base := unsafe.Pointer(&(*C.PyTupleObject)(unsafe.Pointer(ret)).ob_item[0])
	step := uintptr(C.tupleItemSize())

	for _, item := range items {
		item.Incref()
		*(**C.PyObject)(base) = c(item)

		// Move base to point to the next item, by incrementing by step bytes
		base = unsafe.Pointer(uintptr(base) + step)
	}

	return newTuple(ret), nil
}

func (t *Tuple) CheckExact() bool {
	ret := C.tupleCheckE(c(t))
	if int(ret) != 0 {
		return true
	}
	return false
}

func (t *Tuple) Size() int64 {
	ret := C.PyTuple_Size(c(t))
	if ret < 0 {
		panic(exception())
	}
	return int64(ret)
}

func (t *Tuple) GetItem(pos int64) (Object, error) {
	ret := C.PyTuple_GetItem(c(t), C.Py_ssize_t(pos))
	return obj2ObjErr(ret)
}

func (t *Tuple) GetSlice(low, high int64) (*Tuple, error) {
	ret := C.PyTuple_GetSlice(c(t), C.Py_ssize_t(low), C.Py_ssize_t(high))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

func (t *Tuple) SetItem(pos int64, obj Object) error {
	ret := C.PyTuple_SetItem(c(t), C.Py_ssize_t(pos), c(obj))
	return int2Err(ret)
}

// _PyTuple_Resize

// PyTuple_ClearFreeList()

func (t *Tuple) Slice() []Object {
	l := t.Size()
	s := make([]Object, l)
	for i := int64(0); i < l; i++ {
		o, err := t.GetItem(i)
		if err != nil {
			panic(err)
		}
		s[i] = o
	}
	return s
}

func (t *Tuple) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", t.Slice())
}
