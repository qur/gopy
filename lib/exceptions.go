// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

var Exc = _get_exceptions()

type BaseException struct {
	AbstractObject
	o *C.PyBaseExceptionObject
}

func newException(obj *C.PyObject) *BaseException {
	return (*BaseException)(unsafe.Pointer(obj))
}
