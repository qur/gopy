// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline void incref(PyObject *obj) { Py_INCREF(obj); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
// static inline void xdecref(PyObject *obj) { Py_XDECREF(obj); }
import "C"

import (
	"fmt"
	"os"
)

// Error represents a Python exception as a Go struct that implements the
// os.Error interface.  It allows Go code to handle Python exceptions in an
// idiomatic Go fashion.
type Error struct {
	kind *C.PyObject
	tb   *C.PyObject
	val  *C.PyObject
}

// String() returns a string representation of the Python exception represented
// by the Error e.  This is the same as the final line of the Python output from
// an uncaught exception.
func (e *Error) String() string {
	ts := ""
	en := C.excName(e.kind)
	if en.c == nil {
		tpyS := C.PyObject_Str(e.kind)
		defer C.decref(tpyS)
		ts = C.GoString(C.PyString_AsString(tpyS))
	} else {
		if en.m != nil {
			ts = C.GoString(en.m) + "."
		}
		ts += C.GoString(en.c)
	}

	pyS := C.PyObject_Str(e.val)
	defer C.decref(pyS)
	s := C.GoString(C.PyString_AsString(pyS))

	return fmt.Sprintf("%s: %s", ts, s)
}

// NewError returns a new Error of the specified kind, and with the given value.
func NewError(kind Object, value Object) *Error {
	Incref(kind)
	Incref(value)
	return &Error{c(kind), nil, c(value)}
}

// NewErrorString returns a new Error of the specified kind, and with the value
// being a new String containing msg.
func NewErrorString(kind Object, msg string) *Error {
	Incref(kind)
	val, _ := NewString(msg)
	return &Error{c(kind), nil, c(val)}
}

// NewErrorFormat returns a new Error of the specified kind, and with the value
// being a new String containing the string created the given format and args.
func NewErrorFormat(kind Object, format string, args ...interface{}) *Error {
	msg := fmt.Sprintf(format, args...)
	Incref(kind)
	val, _ := NewString(msg)
	return &Error{c(kind), nil, c(val)}
}

func exceptionRaised() bool {
	return C.PyErr_Occurred() != nil
}

func exception() os.Error {
	if C.PyErr_Occurred() == nil {
		return nil
	}

	var t, v, tb *C.PyObject

	defer C.xdecref(v)

	C.PyErr_Fetch(&t, &v, &tb)

	return &Error{t, tb, v}
}

func raise(err os.Error) {
	var val *C.PyObject
	var exc = C.PyExc_Exception

	e, ok := err.(*Error)
	if ok {
		exc = e.kind
		val = e.val
	} else {
		v, _ := NewString(err.String())
		val = c(v)
	}

	C.PyErr_SetObject(exc, val)
}

func TypeError(format string, args ...interface{}) os.Error {
	msg := fmt.Sprintf(format, args...)
	val, _ := NewString(msg)
	C.incref(C.PyExc_TypeError)
	return &Error{C.PyExc_TypeError, nil, c(val)}
}

func KeyError(format string, args ...interface{}) os.Error {
	msg := fmt.Sprintf(format, args...)
	val, _ := NewString(msg)
	C.incref(C.PyExc_KeyError)
	return &Error{C.PyExc_KeyError, nil, c(val)}
}

func AttributeError(format string, args ...interface{}) os.Error {
	msg := fmt.Sprintf(format, args...)
	val, _ := NewString(msg)
	C.incref(C.PyExc_AttributeError)
	return &Error{C.PyExc_AttributeError, nil, c(val)}
}

func NotImplemented(format string, args ...interface{}) os.Error {
	msg := fmt.Sprintf(format, args...)
	val, _ := NewString(msg)
	C.incref(C.PyExc_NotImplementedError)
	return &Error{C.PyExc_NotImplementedError, nil, c(val)}
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
	return newObject(obj), nil
}
