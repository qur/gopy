// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int intCheck(PyObject *o) { return PyInt_Check(o); }
import "C"

import (
	"fmt"
	"unsafe"
)

type Int struct {
	BaseObject
}

func intCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.intCheck(c(obj)) != 0
}

func newInt(obj *C.PyObject) *Int {
	return (*Int)(unsafe.Pointer(obj))
}

func Int_FromInt(i int) *Int {
	return newInt(C.PyInt_FromLong(C.long(i)))
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
