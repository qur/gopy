// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int listCheck(PyObject *o) { return PyList_Check(o); }
// static inline int listCheckE(PyObject *o) { return PyList_CheckExact(o); }
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type List struct {
	BaseObject
}

func listCheck(obj Object) bool {
	return C.listCheck(c(obj)) != 0
}

func newList(obj *C.PyObject) *List {
	return (*List)(unsafe.Pointer(obj))
}

func List_New(size int64) (*List, os.Error) {
	ret := C.PyList_New(C.Py_ssize_t(size))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

func (l *List) CheckExact() bool {
	ret := C.listCheckE(c(l))
	if int(ret) != 0 {
		return true
	}
	return false
}

func (l *List) Size() int64 {
	ret := C.PyList_Size(c(l))
	if ret < 0 {
		panic(exception())
	}
	return int64(ret)
}

func (l *List) GetItem(idx int64) (Object, os.Error) {
	ret := C.PyList_GetItem(c(l), C.Py_ssize_t(idx))
	return obj2ObjErr(ret)
}

func (l *List) SetItem(idx int64, obj Object) os.Error {
	ret := C.PyList_SetItem(c(l), C.Py_ssize_t(idx), c(obj))
	return int2Err(ret)
}

func (l *List) Insert(idx int64, obj Object) os.Error {
	ret := C.PyList_Insert(c(l), C.Py_ssize_t(idx), c(obj))
	return int2Err(ret)
}

func (l *List) Append(obj Object) os.Error {
	ret := C.PyList_Append(c(l), c(obj))
	return int2Err(ret)
}

func (l *List) GetSlice(low, high int64) (*List, os.Error) {
	ret := C.PyList_GetSlice(c(l), C.Py_ssize_t(low), C.Py_ssize_t(high))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

func (l *List) SetSlice(low, high int64, items *List) os.Error {
	ret := C.PyList_SetSlice(c(l), C.Py_ssize_t(low), C.Py_ssize_t(high), c(items))
	return int2Err(ret)
}

func (l *List) Sort() os.Error {
	ret := C.PyList_Sort(c(l))
	return int2Err(ret)
}

func (l *List) Reverse() os.Error {
	ret := C.PyList_Reverse(c(l))
	return int2Err(ret)
}

func (l *List) Tuple() *Tuple {
	ret := C.PyList_AsTuple(c(l))
	if ret == nil {
		panic(exception())
	}
	return newTuple(ret)
}

func (l *List) Slice() []Object {
	size := l.Size()
	s := make([]Object, size)
	for i := int64(0); i < size; i++ {
		o, err := l.GetItem(i)
		if err != nil {
			panic(err)
		}
		s[i] = o
	}
	return s
}

func (l *List) String() string {
	return fmt.Sprintf("%v", l.Slice())
}
