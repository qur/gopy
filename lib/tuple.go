// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int tupleCheck(PyObject *o) { return PyTuple_Check(o); }
// static inline int tupleCheckE(PyObject *o) { return PyTuple_CheckExact(o); }
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type Tuple struct {
	BaseObject
}

func tupleCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.tupleCheck(c(obj)) != 0
}

func newTuple(obj *C.PyObject) *Tuple {
	return (*Tuple)(unsafe.Pointer(obj))
}

func buildTuple(format string, args ...interface{}) (*Tuple, os.Error) {
	if format == "" {
		return NewTuple(0)
	}
	bv, err := BuildValue(format, args...)
	if err != nil {
		return nil, err
	}
	t, ok := bv.(*Tuple)
	if ok {
		return t, nil
	}
	return Tuple_Pack(bv)
}

func NewTuple(size int64) (*Tuple, os.Error) {
	ret := C.PyTuple_New(C.Py_ssize_t(size))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

func Tuple_Pack(items ...Object) (*Tuple, os.Error) {
	// TODO: ...
	return nil, fmt.Errorf("Not Implemented")
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

func (t *Tuple) GetItem(pos int64) (Object, os.Error) {
	ret := C.PyTuple_GetItem(c(t), C.Py_ssize_t(pos))
	return obj2ObjErr(ret)
}

func (t *Tuple) GetSlice(low, high int64) (*Tuple, os.Error) {
	ret := C.PyTuple_GetSlice(c(t), C.Py_ssize_t(low), C.Py_ssize_t(high))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

func (t *Tuple) SetItem(pos int64, obj Object) os.Error {
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
