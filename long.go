// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline long longCheck(PyObject *o) { return PyLong_Check(o); }
import "C"

import (
	"fmt"
	"unsafe"
)

type Long struct {
	AbstractObject
	NumberProtocol
	o *C.PyLongObject
}

var longObjMap = make(map[*C.PyObject]*Long)

// LongType is the Type object that represents the Long type.
var LongType = newType((*C.PyObject)(unsafe.Pointer(C.getBasePyType(C.GoPyLong_Type))))

func longCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.longCheck(c(obj)) != 0
}

func newLong(obj *C.PyObject) *Long {
	if l, ok := longObjMap[obj]; ok {
		return l
	}
	l := &Long{o: (*C.PyLongObject)(unsafe.Pointer(obj))}
	longObjMap[obj] = l
	return l
}

func NewLong(i int64) *Long {
	return newLong(C.PyLong_FromLongLong(C.PY_LONG_LONG(i)))
}

func (l *Long) Int64() int64 {
	// TODO: AsLongLong doesn't work for me on windows...
	return int64(C.PyLong_AsLong(c(l)))
}

func (l *Long) String() string {
	if l == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", l.Int64())
}
