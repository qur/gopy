// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

type ExceptionClass struct {
	AbstractObject
	o C.PyBaseExceptionObject
}

func newException(obj *C.PyObject) *ExceptionClass {
	return (*ExceptionClass)(unsafe.Pointer(obj))
}

// ErrV returns a new Error of the specified kind, and with the given value.
func (kind *ExceptionClass) ErrV(obj Object) *Error {
	return NewErrorV(kind, obj)
}

// Err returns a new Error of the specified kind, and with the value being a
// new String containing the string created the given format and args.
func (kind *ExceptionClass) Err(format string, args ...interface{}) *Error {
	return NewError(kind, format, args...)
}
