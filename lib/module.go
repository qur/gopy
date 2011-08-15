// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int moduleCheck(PyObject *o) { return PyModule_Check(o); }
// static inline int moduleCheckE(PyObject *o) { return PyModule_CheckExact(o); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type Module struct {
	AbstractObject
	o C.PyObject
}

var ModuleType = (*Type)(unsafe.Pointer(&C.PyModule_Type))

func moduleCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.moduleCheck(c(obj)) != 0
}

func newModule(obj *C.PyObject) *Module {
	return (*Module)(unsafe.Pointer(obj))
}

func Import(name string) (*Module, os.Error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))

	pyName := C.PyString_FromString(s)
	defer C.decref(pyName)

	obj := C.PyImport_Import(pyName)
	if obj == nil {
		return nil, exception()
	}

	return newModule(obj), nil
}

func InitModule(name string, methods []Method) (*Module, os.Error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	m := C.Py_InitModule4(cName, nil, nil, nil, C.PYTHON_API_VERSION)
	if m == nil {
		return nil, exception()
	}

	if len(methods) == 0 {
		return newModule(m), nil
	}

	n := C.PyString_FromString(cName)
	if n == nil {
		return nil, exception()
	}

	d := C.PyModule_GetDict(m)
	if d == nil {
		return nil, exception()
	}

	for _, method := range methods {
		ml := C.newMethodDef()

		ml.ml_name = C.CString(method.Name)
		ml.ml_doc = C.CString(method.Doc)

		switch f := method.Func.(type) {

		case func() (Object, os.Error):
			C.set_call_noargs(&ml.ml_meth)
			ml.ml_flags = C.METH_NOARGS

		case func(a *Tuple) (Object, os.Error):
			C.set_call_args(&ml.ml_meth)
			ml.ml_flags = C.METH_VARARGS

		case func(a *Tuple, k *Dict) (Object, os.Error):
			C.set_call_keywords(&ml.ml_meth)
			ml.ml_flags = C.METH_VARARGS | C.METH_KEYWORDS

		default:
			return nil, fmt.Errorf("InitModule: unknown func type for %s", method.Name)

		}

		v := C.PyCFunction_NewEx(ml, saveFunc(method.Func), n)
		if v == nil {
			return nil, exception()
		}

		if C.PyDict_SetItemString(d, ml.ml_name, v) != 0 {
			return nil, exception()
		}
	}

	return newModule(m), nil
}

func ExecCodeModule(name string, code Object) (*Module, os.Error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyImport_ExecCodeModule(s, c(code))
	if ret == nil {
		return nil, exception()
	}
	return newModule(ret), nil
}

func NewModule(name string) (*Module, os.Error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_New(cname)
	if ret == nil {
		return nil, exception()
	}

	return newModule(ret), nil
}

func (mod *Module) CheckExact() bool {
	return C.moduleCheckE(c(mod)) != 0
}

func (mod *Module) Dict() *Dict {
	ret := C.PyModule_GetDict(c(mod))
	return newDict(ret)
}

func (mod *Module) Name() (string, os.Error) {
	ret := C.PyModule_GetName(c(mod))
	if ret == nil {
		return "", exception()
	}
	return C.GoString(ret), nil
}

func (mod *Module) Filename() (string, os.Error) {
	ret := C.PyModule_GetFilename(c(mod))
	if ret == nil {
		return "", exception()
	}
	return C.GoString(ret), nil
}

func (mod *Module) AddObject(name string, obj Object) os.Error {
	if obj == nil {
		return fmt.Errorf("ValueError: obj == nil!")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddObject(c(mod), cname, c(obj))
	if ret < 0 {
		return exception()
	}

	return nil
}

func (mod *Module) AddIntConstant(name string, value int) os.Error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddIntConstant(c(mod), cname, C.long(value))
	if ret < 0 {
		return exception()
	}

	return nil
}

func (mod *Module) AddStringConstant(name, value string) os.Error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	ret := C.PyModule_AddStringConstant(c(mod), cname, cvalue)
	if ret < 0 {
		return exception()
	}

	return nil
}
