// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type BaseObject struct {
	AbstractObject
	C.PyObject
}

func newBaseObject(obj *C.PyObject) *BaseObject {
	return (*BaseObject)(unsafe.Pointer(obj))
}

func (obj *BaseObject) Call(args *Tuple, kwds *Dict) (Object, os.Error) {
	ret := C.PyObject_Call(c(obj), c(args), c(kwds))
	return obj2ObjErr(ret)
}

func (obj *BaseObject) CallObject(args *Tuple) (Object, os.Error) {
	var a *C.PyObject = nil
	if args != nil {
		a = c(args)
	}
	ret := C.PyObject_CallObject(c(obj), a)
	return obj2ObjErr(ret)
}

func (obj *BaseObject) CallFunction(format string, args ...interface{}) (Object, os.Error) {
	t, err := buildTuple(format, args...)
	if err != nil {
		return nil, err
	}
	return obj.CallObject(t)
}

func (obj *BaseObject) CallMethod(name string, format string, args ...interface{}) (Object, os.Error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	f := C.PyObject_GetAttrString(c(obj), cname)
	if f == nil {
		return nil, fmt.Errorf("AttributeError: %s", name)
	}
	defer C.decref(f)

	if C.PyCallable_Check(f) == 0 {
		return nil, fmt.Errorf("TypeError: attribute of type '%s' is not callable", name)
	}

	t, err := buildTuple(format, args...)
	if err != nil {
		return nil, err
	}

	ret := C.PyObject_CallObject(f, c(t))
	return obj2ObjErr(ret)
}

func (obj *BaseObject) CallFunctionObjArgs(args ...Object) (Object, os.Error) {
	t, err := PackTuple(args...)
	if err != nil {
		return nil, err
	}
	return obj.CallObject(t)
}

func (obj *BaseObject) CallMethodObjArgs(name string, args ...Object) (Object, os.Error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	f := C.PyObject_GetAttrString(c(obj), cname)
	if f == nil {
		return nil, fmt.Errorf("AttributeError: %s", name)
	}
	defer C.decref(f)

	if C.PyCallable_Check(f) == 0 {
		return nil, fmt.Errorf("TypeError: attribute of type '%s' is not callable", name)
	}

	t, err := PackTuple(args...)
	if err != nil {
		return nil, err
	}

	ret := C.PyObject_CallObject(f, c(t))
	return obj2ObjErr(ret)
}

func (obj *BaseObject) Dir() (Object, os.Error) {
	ret := C.PyObject_Dir(c(obj))
	return obj2ObjErr(ret)
}
