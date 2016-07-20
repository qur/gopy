// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// #include "gopy_types.h"
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
// static inline int exceptionCheck(PyObject *obj) {
//     return PyExceptionClass_Check(obj);
// }
import "C"

import (
	"reflect"
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

func stringify(obj Object) string {
	tpyS := C.PyObject_Str(c(obj))
	defer C.decref(tpyS)
	u := C.PyUnicode_AsUTF8String(tpyS)
	defer C.decref(u)
	return C.GoString(C.PyBytes_AsString(u))
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

func obj2Class(c *Class, obj *C.PyObject) (Object, bool) {
	vp := reflect.NewAt(reflect.TypeOf(c.Pointer), unsafe.Pointer(&obj))
	o, ok := vp.Elem().Interface().(Object)
	return o, ok
}

func newObject(obj *C.PyObject) Object {
	if obj == nil {
		return nil
	}

	o := unsafe.Pointer(obj)
	if o == unsafe.Pointer(None) {
		return None
	}

	pyType := (*C.PyTypeObject)(obj.ob_type)
	class, ok := getType(pyType)
	if ok {
		ret, ok := obj2Class(class, obj)
		if ok {
			return ret
		}
	}
	for pyType.tp_base != nil {
		pyType = (*C.PyTypeObject)(unsafe.Pointer(pyType.tp_base))
		class, ok := getType(pyType)
		if ok {
			ret, ok := obj2Class(class, obj)
			if ok {
				return ret
			}
		}
	}

	switch C.getBaseGoPyType(obj) {
	case C.GoPyList_Type:
		return (*List)(o)
	case C.GoPyTuple_Type:
		return (*Tuple)(o)
	case C.GoPyDict_Type:
		return (*Dict)(o)
	case C.GoPyUnicode_Type:
		return (*Unicode)(o)
	case C.GoPyBool_Type:
		return newBool(obj)
	case C.GoPyLong_Type:
		return (*Long)(o)
	case C.GoPyFloat_Type:
		return (*Float)(o)
	case C.GoPyModule_Type:
		return (*Module)(o)
	case C.GoPyType_Type:
		return (*Type)(o)
	case C.GoPyCode_Type:
		return (*Code)(o)
	case C.GoPyCFunction_Type:
		return (*CFunction)(o)
	case C.GoPyComplex_Type:
		return (*Complex)(o)
	case C.GoPyFrozenSet_Type:
		return (*FrozenSet)(o)
	case C.GoPySet_Type:
		return (*Set)(o)
	case C.GoPyFunction_Type:
		return (*Function)(o)
	}
	if C.exceptionCheck(obj) != 0 {
		return newException(obj)
	}
	return newBaseObject(obj)
}
