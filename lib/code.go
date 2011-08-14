// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int codeCheck(PyObject *o) { return PyCode_Check(o); }
import "C"

import (
	"unsafe"
)

type Code struct {
	AbstractObject
	C.PyCodeObject
}

var CodeType = (*Type)(unsafe.Pointer(&C.PyCode_Type))

func codeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.codeCheck(c(obj)) != 0
}
