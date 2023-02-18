// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"sync"
	"unsafe"
)

type Op int

const (
	LT = Op(C.Py_LT)
	LE = Op(C.Py_LE)
	EQ = Op(C.Py_EQ)
	NE = Op(C.Py_NE)
	GT = Op(C.Py_GT)
	GE = Op(C.Py_GE)
)

// Object is the generic interface that represents a Python object.  All of the
// concrete types satisfy the Object interface.
type Object interface {
	Base() *BaseObject
	Type() *Type
	Decref()
	Incref()
	IsTrue() bool
	Not() bool
	Free()
}

// None is the Python equivalent to nil.
var None = (*NoneObject)(unsafe.Pointer(&C._Py_NoneStruct))

// NoneObject is the type of the None value.  The only value of this type is
// None.
type NoneObject struct {
	AbstractObject
}

func (n *NoneObject) String() string {
	return "None"
}

func c(obj Object) *C.PyObject {
	if obj == nil {
		return nil
	}
	return (*C.PyObject)(unsafe.Pointer(obj.Base()))
}

var (
	typeLock sync.RWMutex
	types    = make(map[*C.PyTypeObject]*Class)
)

func registerType(pyType *C.PyTypeObject, class *Class) {
	typeLock.Lock()
	defer typeLock.Unlock()

	types[pyType] = class
}

func getType(pyType *C.PyTypeObject) (*Class, bool) {
	typeLock.RLock()
	defer typeLock.RUnlock()

	class, ok := types[pyType]

	return class, ok
}

func newObject(obj *C.PyObject) Object {
	if obj == nil {
		return nil
	}

	o := unsafe.Pointer(obj)
	if o == unsafe.Pointer(None) {
		return None
	}

	cObj := getClassObject(o)
	if cObj != nil {
		return cObj
	}

	switch C.getBasePyType(obj) {
	case &C.PyList_Type:
		return (*List)(o)
	case &C.PyTuple_Type:
		return (*Tuple)(o)
	case &C.PyDict_Type:
		return (*Dict)(o)
	case &C.PyBool_Type:
		return newBool(obj)
	case &C.PyBytes_Type:
		return (*Bytes)(o)
	case &C.PyLong_Type:
		return (*Long)(o)
	case &C.PyFloat_Type:
		return (*Float)(o)
	case &C.PyModule_Type:
		return (*Module)(o)
	case &C.PyType_Type:
		class, ok := getType((*C.PyTypeObject)(o))
		if ok {
			return class
		}
		return (*Type)(o)
	case &C.PyCode_Type:
		return (*Code)(o)
	case &C.PyCFunction_Type:
		return (*CFunction)(o)
	case &C.PyComplex_Type:
		return (*Complex)(o)
	case &C.PyFrozenSet_Type:
		return (*FrozenSet)(o)
	case &C.PySet_Type:
		return (*Set)(o)
	case &C.PyFunction_Type:
		return (*Function)(o)
	case &C.PyFrame_Type:
		return (*Frame)(o)
	case &C.PyMethod_Type:
		return (*Method)(o)
	case &C.PyUnicode_Type:
		return (*Unicode)(o)
	}

	if C.exceptionCheck(obj) != 0 {
		return newException(obj)
	}

	return newBaseObject(obj)
}
