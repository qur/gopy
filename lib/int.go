// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type Int struct {
	AbstractObject
	NumberProtocol
	o C.PyIntObject
}

// IntType is the Type object that represents the Int type.
var IntType = (*Type)(unsafe.Pointer(&C.PyInt_Type))

func intCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.intCheck(c(obj)) != 0
}

func newInt(obj *C.PyObject) *Int {
	return (*Int)(unsafe.Pointer(obj))
}

func NewInt(i int) *Int {
	return newInt(C.PyInt_FromLong(C.long(i)))
}

func NewInt64(i int64) *Int {
	return newInt(C.PyInt_FromSsize_t(C.Py_ssize_t(i)))
}

func (i *Int) Int() int {
	return int(C.PyInt_AsLong(c(i)))
}

func (i *Int) String() string {
	if i == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", i.Int())
}
