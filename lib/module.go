// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

type Module struct {
	AbstractObject
	o C.PyObject
}

// ModuleType is the Type object that represents the Module type.
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

func Import(name string) (*Module, error) {
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

func InitModule(name string, methods []Method) (*Module, error) {
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
		pyF, err := makeCFunction(method.Name, method.Func, method.Doc, n)
		if err != nil {
			return nil, err
		}

		if C.PyDict_SetItemString(d, C.CString(method.Name), c(pyF)) != 0 {
			return nil, exception()
		}
	}

	return newModule(m), nil
}

func ExecCodeModule(name string, code Object) (*Module, error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyImport_ExecCodeModule(s, c(code))
	if ret == nil {
		return nil, exception()
	}
	return newModule(ret), nil
}

func NewModule(name string) (*Module, error) {
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

func (mod *Module) Name() (string, error) {
	ret := C.PyModule_GetName(c(mod))
	if ret == nil {
		return "", exception()
	}
	return C.GoString(ret), nil
}

func (mod *Module) Filename() (string, error) {
	ret := C.PyModule_GetFilename(c(mod))
	if ret == nil {
		return "", exception()
	}
	return C.GoString(ret), nil
}

func (mod *Module) AddObject(name string, obj Object) error {
	if obj == nil {
		return AssertionError.Err("ValueError: obj == nil!")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddObject(c(mod), cname, c(obj))
	if ret < 0 {
		return exception()
	}

	return nil
}

func (mod *Module) AddIntConstant(name string, value int) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddIntConstant(c(mod), cname, C.long(value))
	if ret < 0 {
		return exception()
	}

	return nil
}

func (mod *Module) AddStringConstant(name, value string) error {
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
