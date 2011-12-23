// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline void incref(PyObject *obj) { Py_INCREF(obj); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
// static inline void xdecref(PyObject *obj) { Py_XDECREF(obj); }
import "C"

import "fmt"

// Error represents a Python exception as a Go struct that implements the
// error interface.  It allows Go code to handle Python exceptions in an
// idiomatic Go fashion.
type Error struct {
	Kind  Object
	Value Object
	tb    *C.PyObject
}

// Error() returns a string representation of the Python exception represented
// by the Error e.  This is the same as the final line of the Python output from
// an uncaught exception.
func (e *Error) Error() string {
	ts := ""
	en := C.excName(c(e.Kind))
	if en.c == nil {
		tpyS := C.PyObject_Str(c(e.Kind))
		defer C.decref(tpyS)
		ts = C.GoString(C.PyString_AsString(tpyS))
	} else {
		if en.m != nil {
			ts = C.GoString(en.m) + "."
		}
		ts += C.GoString(en.c)
	}

	pyS := C.PyObject_Str(c(e.Value))
	defer C.decref(pyS)
	s := C.GoString(C.PyString_AsString(pyS))

	return fmt.Sprintf("%s: %s", ts, s)
}

// Matches returns true if e.Kind matches the exception in exc.  If exc is a
// Class, then true is returned if e.Kind is an instance.  If exc is a Tuple,
// then all elements (and recursively for sub elements) are searched for a
// match.
func (e *Error) Matches(exc Object) bool {
	return C.PyErr_GivenExceptionMatches(c(e.Kind), c(exc)) != 0
}

// Normalize adjusts e.Kind/e.Value in the case that the values aren't
// normalized to start with.  It's possible that an Error returned from Python
// might have e.Kind be a Class, with e.Value not being an instance of that
// class, Normalize will fix this.  The separate normalization is implemented in
// Python to improve performance.
func (e *Error) Normalize() {
	exc := c(e.Kind)
	val := c(e.Value)
	tb := e.tb
	C.PyErr_NormalizeException(&exc, &val, &tb)
	if exc != c(e.Kind) {
		e.Kind = newObject(exc)
	}
	if val != c(e.Value) {
		e.Value = newObject(val)
	}
	e.tb = tb
}

// NewError returns a new Error of the specified kind, and with the given value.
func NewError(kind Object, value Object) *Error {
	Incref(kind)
	Incref(value)
	return &Error{kind, value, nil}
}

// NewErrorString returns a new Error of the specified kind, and with the value
// being a new String containing msg.
func NewErrorString(kind Object, msg string) *Error {
	Incref(kind)
	val, _ := NewString(msg)
	return &Error{kind, val, nil}
}

// NewErrorFormat returns a new Error of the specified kind, and with the value
// being a new String containing the string created the given format and args.
func NewErrorFormat(kind Object, format string, args ...interface{}) *Error {
	msg := fmt.Sprintf(format, args...)
	Incref(kind)
	val, _ := NewString(msg)
	return &Error{kind, val, nil}
}

func exceptionRaised() bool {
	return C.PyErr_Occurred() != nil
}

func exception() error {
	if C.PyErr_Occurred() == nil {
		return nil
	}

	var t, v, tb *C.PyObject

	defer C.xdecref(v)

	C.PyErr_Fetch(&t, &v, &tb)

	return &Error{newObject(t), newObject(v), tb}
}

func raise(err error) {
	var val *C.PyObject
	var exc = C.PyExc_Exception

	e, ok := err.(*Error)
	if ok {
		exc = c(e.Kind)
		val = c(e.Value)
	} else {
		v, _ := NewString(err.Error())
		val = c(v)
	}

	C.PyErr_SetObject(exc, val)
}

func TypeError(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	val, _ := NewString(msg)
	Exc.TypeError.Incref()
	return &Error{Exc.TypeError, val, nil}
}

func KeyError(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	val, _ := NewString(msg)
	Exc.KeyError.Incref()
	return &Error{Exc.KeyError, val, nil}
}

func AttributeError(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	val, _ := NewString(msg)
	Exc.AttributeError.Incref()
	return &Error{Exc.AttributeError, val, nil}
}

func NotImplemented(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	val, _ := NewString(msg)
	Exc.NotImplementedError.Incref()
	return &Error{Exc.NotImplementedError, val, nil}
}

func int2Err(i C.int) error {
	if i < 0 {
		return exception()
	}
	return nil
}

func int2BoolErr(i C.int) (bool, error) {
	if i < 0 {
		return false, exception()
	}
	return i > 0, nil
}

func ssize_t2Int64Err(s C.Py_ssize_t) (int64, error) {
	if s < 0 {
		return 0, exception()
	}
	return int64(s), nil
}

func obj2ObjErr(obj *C.PyObject) (Object, error) {
	if obj == nil {
		return nil, exception()
	}
	return newObject(obj), nil
}
