// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

type MethodObj struct {
	AbstractObject
	o C.PyMethodObject
}

// MethodType is the Type object that represents the MethodObj type.
var MethodType = (*Type)(unsafe.Pointer(&C.PyMethod_Type))

func newMethod(obj *C.PyObject) *MethodObj {
	return (*MethodObj)(unsafe.Pointer(obj))
}

func methodCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.methodCheck(c(obj)) != 0
}

func NewMethodObj(function, self, class Object) (*MethodObj, error) {
	ret := C.PyMethod_New(c(function), c(self), c(class))
	if ret == nil {
		return nil, exception()
	}
	return newMethod(ret), nil
}

func (m *MethodObj) Class() Object {
	ret := C.PyMethod_Class(c(m))
	return newObject(ret)
}

func (m *MethodObj) Self() Object {
	ret := C.PyMethod_Self(c(m))
	return newObject(ret)
}
