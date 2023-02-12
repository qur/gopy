// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

type Method struct {
	AbstractObject
	o C.PyMethodObject
}

// MethodType is the Type object that represents the MethodObj type.
var MethodType = (*Type)(unsafe.Pointer(&C.PyMethod_Type))

func newMethod(obj *C.PyObject) *Method {
	return (*Method)(unsafe.Pointer(obj))
}

func methodCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.methodCheck(c(obj)) != 0
}

func NewMethod(function, self Object) (*Method, error) {
	ret := C.PyMethod_New(c(function), c(self))
	if ret == nil {
		return nil, exception()
	}
	return newMethod(ret), nil
}

func (m *Method) Self() Object {
	ret := C.PyMethod_Self(c(m))
	return newObject(ret)
}