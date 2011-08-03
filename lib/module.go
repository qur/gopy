// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int moduleCheck(PyObject *o) { return PyModule_Check(o); }
// static inline void decref(PyObject *obj) { Py_DECREF(obj); }
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type Module struct {
	BaseObject
}

func moduleCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.moduleCheck(c(obj)) != 0
}

func newModule(obj *C.PyObject) *Module {
	return (*Module)(unsafe.Pointer(obj))
}

type Method struct {
	Name string
	Func interface{}
	Doc  string
}

var funcs []interface{}

func saveFunc(f interface{}) *C.PyObject {
	funcs = append(funcs, f)
	return C.PyInt_FromLong(C.long(len(funcs) - 1))
}

//export callWithArgs
func callWithArgs(self, args unsafe.Pointer) unsafe.Pointer {
	_idx := (*C.PyObject)(self)
	idx := C.PyInt_AsLong(_idx)
	f, ok := funcs[idx].(func(a *Tuple) (Object, os.Error))
	if !ok {
		fmt.Printf("invalid index: %d\n", idx)
		return nil
	}
	a := newTuple((*C.PyObject)(args))
	ret, err := f(a)
	if err != nil {
		raise(err)
		return nil
	}
	return unsafe.Pointer(c(ret))
}

//export callWithKeywords
func callWithKeywords(self, args, kw unsafe.Pointer) unsafe.Pointer {
	_idx := (*C.PyObject)(self)
	idx := C.PyInt_AsLong(_idx)
	f, ok := funcs[idx].(func(a *Tuple, k *Dict) (Object, os.Error))
	if !ok {
		fmt.Printf("invalid index: %d\n", idx)
		return nil
	}
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kw))
	ret, err := f(a, k)
	if err != nil {
		raise(err)
		return nil
	}
	return unsafe.Pointer(c(ret))
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
