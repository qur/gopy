// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include <Python.h>
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
// static inline void xdecref(PyObject *obj) { Py_XDECREF(obj); }
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type Error struct {
	kind *C.PyObject
	msg  string
}

func (e *Error) String() string {
	return fmt.Sprintf("<TYPE>: %s", e.msg)
}

func exception() (e os.Error) {
	e = nil
	var t *C.PyObject
	var v *C.PyObject
	var tb *C.PyObject
	C.PyErr_Fetch(&t, &v, &tb)
	if v != nil {
		s := C.PyObject_Str(v)
		cs := C.PyString_AsString(s)
		e = fmt.Errorf("%s", C.GoString(cs))
		C.decref(s)
	}
	C.xdecref(t)
	C.xdecref(v)
	C.xdecref(tb)
	return
}

func raise(err os.Error) {
	var msg string
	var exc = C.PyExc_Exception

	e, ok := err.(*Error)
	if ok {
		exc = e.kind
		msg = e.msg
	} else {
		msg = err.String()
	}

	s := C.CString(msg)
	defer C.free(unsafe.Pointer(s))

	C.PyErr_SetString(exc, s)
}

func Err_Format(f string, args ...interface{}) {
	s := C.CString(fmt.Sprintf(f, args))
	defer C.free(unsafe.Pointer(s))
	C.PyErr_SetString(C.PyExc_Exception, s)
}

func KeyError(format string, args ...interface{}) os.Error {
	msg := fmt.Sprintf(format, args...)
	return &Error{C.PyExc_KeyError, msg}
}

func int2Err(i C.int) os.Error {
	if i < 0 {
		return exception()
	}
	return nil
}

func int2BoolErr(i C.int) (bool, os.Error) {
	if i < 0 {
		return false, exception()
	}
	return i > 0, nil
}

func ssize_t2Int64Err(s C.Py_ssize_t) (int64, os.Error) {
	if s < 0 {
		return 0, exception()
	}
	return int64(s), nil
}

func obj2ObjErr(obj *C.PyObject) (Object, os.Error) {
	if obj == nil {
		return nil, exception()
	}
	return newBaseObject(obj).actual(), nil
}
