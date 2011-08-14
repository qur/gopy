// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int boolCheck(PyObject *o) { return PyBool_Check(o); }
// static inline void *pyTrue(void) { return Py_True; }
// static inline void *pyFalse(void) { return Py_False; }
import "C"

import (
	"fmt"
	"unsafe"
)

type Bool struct {
	AbstractObject
	o C.PyBoolObject
}

var BoolType = (*Type)(unsafe.Pointer(&C.PyBool_Type))

var True = (*Bool)(C.pyTrue())
var False = (*Bool)(C.pyFalse())

func boolCheck(obj Object) bool {
	return C.boolCheck(c(obj)) != 0
}

func newBool(obj *C.PyObject) *Bool {
	if obj == c(True) {
		return True
	}
	if obj == c(False) {
		return False
	}
	panic(fmt.Errorf("TypeError: not a bool"))
}

func (b *Bool) Bool() bool {
	if b == True {
		return true
	}
	if b == False {
		return false
	}
	panic(fmt.Errorf("TypeError: not a bool"))
}

func (b *Bool) String() string {
	if b == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", b.Bool())
}
