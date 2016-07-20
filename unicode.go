// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int unicodeCheck(PyObject *o) { return PyUnicode_Check(o); }
import "C"

import "unsafe"

type Unicode struct {
	AbstractObject
	o C.PyUnicodeObject
}

// UnicodeType is the Type object that represents the Unicode type.
var UnicodeType = (*Type)(unsafe.Pointer(C.getBasePyType(C.GoPyUnicode_Type)))

func unicodeCheck(obj Object) bool {
	return C.unicodeCheck(c(obj)) != 0
}

func newUnicode(obj *C.PyObject) *Unicode {
	return (*Unicode)(unsafe.Pointer(obj))
}

func NewUnicode(s string) (*Unicode, error) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	ret := C.PyUnicode_FromStringAndSize(cs, C.Py_ssize_t(len(s)))
	if ret == nil {
		return nil, exception()
	}
	return newUnicode(ret), nil
}

func (s *Unicode) String() string {
	if s == nil {
		return "<nil>"
	}
	return stringify(s)
}

func (s *Unicode) Format(args *Tuple) (*Unicode, error) {
	ret := C.PyUnicode_Format(c(s), c(args))
	if ret == nil {
		return nil, exception()
	}
	return newUnicode(ret), nil
}
