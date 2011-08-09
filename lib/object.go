// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline void incref(PyObject *obj) { Py_INCREF(obj); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
// static inline void typeFree(PyTypeObject *type, PyObject *o) { type->tp_free(o); }
import "C"

import (
	"fmt"
	"os"
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

type Object interface {
	Base() *BaseObject
	Type() *Type
	Decref()
	Incref()
	IsTrue() bool
	Not() bool
	Free()
}

var None = (*NoneObject)(unsafe.Pointer(&C._Py_NoneStruct))

type BaseObject struct {
	C.PyObject
}

type NoneObject struct {
	BaseObject
}

func (n *NoneObject) String() string {
	return "None"
}

func newBaseObject(obj *C.PyObject) *BaseObject {
	return (*BaseObject)(unsafe.Pointer(obj))
}

func c(obj Object) *C.PyObject {
	if obj == nil {
		return nil
	}
	return (*C.PyObject)(unsafe.Pointer(obj.Base()))
}

func (obj *BaseObject) Base() *BaseObject {
	return obj
}

func (obj *BaseObject) Type() *Type {
	o := c(obj).ob_type
	return newType((*C.PyObject)(unsafe.Pointer(o)))
}

func Decref(obj Object) {
	if obj != nil {
		C.decref(c(obj))
	}
}

func (obj *BaseObject) Decref() {
	C.decref(c(obj))
}

func Incref(obj Object) {
	if obj != nil {
		C.incref(c(obj))
	}
}

func (obj *BaseObject) Incref() {
	C.incref(c(obj))
}

func (obj *BaseObject) Free() {
	o := c(obj)
	pyType := (*C.PyTypeObject)(unsafe.Pointer(o.ob_type))
	C.typeFree(pyType, o)

	// Make sure this instance isn't registered anymore
	contexts[uintptr(unsafe.Pointer(o))] = nil, false
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
	t, err := Tuple_Pack(args...)
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

	t, err := Tuple_Pack(args...)
	if err != nil {
		return nil, err
	}

	ret := C.PyObject_CallObject(f, c(t))
	return obj2ObjErr(ret)
}

func (obj *BaseObject) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(obj))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

func (obj *BaseObject) Not() bool {
	ret := C.PyObject_Not(c(obj))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

func (obj *BaseObject) Dir() (Object, os.Error) {
	ret := C.PyObject_Dir(c(obj))
	return obj2ObjErr(ret)
}

var types = make(map[*C.PyTypeObject]*Class)

func registerType(pyType *C.PyTypeObject, class *Class) {
	types[pyType] = class
}

func (obj *BaseObject) actual() Object {
	if obj == nil {
		return nil
	}
	o := unsafe.Pointer(obj)
	if o == unsafe.Pointer(None) {
		return None
	}
	pyType := (*C.PyTypeObject)(c(obj).ob_type)
	class, ok := types[pyType]
	if ok {
		t := unsafe.Typeof(class.Pointer)
		ret, ok := unsafe.Unreflect(t, unsafe.Pointer(&obj)).(Object)
		if ok {
			return ret
		}
	}
	switch C.getBasePyType(c(obj)) {
	case &C.PyList_Type:
		return (*List)(o)
	case &C.PyTuple_Type:
		return (*Tuple)(o)
	case &C.PyDict_Type:
		return (*Dict)(o)
	case &C.PyString_Type:
		return (*String)(o)
	case &C.PyBool_Type:
		return newBool(c(obj))
	case &C.PyInt_Type:
		return (*Int)(o)
	case &C.PyLong_Type:
		return (*Long)(o)
	case &C.PyFloat_Type:
		return (*Float)(o)
	case &C.PyModule_Type:
		return (*Module)(o)
	case &C.PyType_Type:
		return (*Type)(o)
	case &C.PyCode_Type:
		return (*Code)(o)
	case &C.PyCFunction_Type:
		return (*CFunction)(o)
	}
	for pyType.tp_base != nil {
		pyType = (*C.PyTypeObject)(unsafe.Pointer(pyType.tp_base))
		class, ok := types[pyType]
		if ok {
			t := unsafe.Typeof(class.Pointer)
			ret, ok := unsafe.Unreflect(t, unsafe.Pointer(&obj)).(Object)
			if ok {
				return ret
			}
		}
	}
	return obj
}
