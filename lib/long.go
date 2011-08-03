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
	BaseObject
}

func longCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.longCheck(c(obj)) != 0
}

func newLong(obj *C.PyObject) *Long {
	return (*Long)(unsafe.Pointer(obj))
}

func Long_FromInt64(i int64) *Long {
	return newLong(C.PyLong_FromLongLong(C.longlong(i)))
}

func (l *Long) Int64() int64 {
	return int64(C.PyLong_AsLongLong(c(l)))
}

func (i *Long) String() string {
	return fmt.Sprintf("%v", i.Int64())
}
