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

type CFunction struct {
	AbstractObject
	o C.PyCFunctionObject
}

func cfunctionCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.cFunctionCheck(c(obj)) != 0
}

func newCFunction(obj *C.PyObject) *CFunction {
	return (*CFunction)(unsafe.Pointer(obj))
}

func NewCFunction(name string, fn interface{}, doc string) (*CFunction, error) {
	return makeCFunction(name, fn, doc, nil)
}

func makeCFunction(name string, fn interface{}, doc string, mod_name *C.PyObject) (*CFunction, error) {
	ml := C.newMethodDef()

	switch fn.(type) {

	case func() (Object, error):
		C.set_call_noargs(&ml.ml_meth)
		ml.ml_flags = C.METH_NOARGS

	case func(a *Tuple) (Object, error):
		C.set_call_args(&ml.ml_meth)
		ml.ml_flags = C.METH_VARARGS

	case func(a *Tuple, k *Dict) (Object, error):
		C.set_call_keywords(&ml.ml_meth)
		ml.ml_flags = C.METH_VARARGS | C.METH_KEYWORDS

	default:
		C.free(unsafe.Pointer(ml))
		return nil, TypeError.Err("CFunction_New: unknown func type for %s", name)

	}

	ret := C.PyCFunction_NewEx(ml, saveFunc(fn), mod_name)
	if ret == nil {
		C.free(unsafe.Pointer(ml))
		return nil, exception()
	}

	ml.ml_name = C.CString(name)
	ml.ml_doc = C.CString(doc)

	return newCFunction(ret), nil
}

// PyCFunction_GetFunction

func (f *CFunction) Self() (Object, error) {
	ret := C.PyCFunction_GetSelf(c(f))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret), nil
}

func (f *CFunction) Flags() (int, error) {
	ret := C.PyCFunction_GetFlags(c(f))
	return int(ret), exception()
}

func (f *CFunction) Call(args *Tuple, kw *Dict) (Object, error) {
	ret := C.PyCFunction_Call(c(f), c(args), c(kw))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret), nil
}

type Method struct {
	Name string
	Func interface{}
	Doc  string
}

var (
	funcLock sync.RWMutex
	funcs    []interface{}
)

func saveFunc(f interface{}) *C.PyObject {
	funcLock.Lock()
	defer funcLock.Unlock()

	funcs = append(funcs, f)
	return C.PyInt_FromLong(C.long(len(funcs) - 1))
}

func getFunc(self unsafe.Pointer) interface{} {
	funcLock.RLock()
	defer funcLock.RUnlock()

	idx := int(C.PyInt_AsLong((*C.PyObject)(self)))

	if idx >= len(funcs) {
		return nil
	}

	return funcs[idx]
}

//export callWithoutArgs
func callWithoutArgs(self, args unsafe.Pointer) unsafe.Pointer {
	f, ok := getFunc(self).(func() (Object, error))
	if !ok {
		raise(AssertionError.Err("callWithoutArgs: wrong function type!!!"))
		return nil
	}
	ret, err := f()
	if err != nil {
		raise(err)
		return nil
	}
	return unsafe.Pointer(c(ret))
}

//export callWithArgs
func callWithArgs(self, args unsafe.Pointer) unsafe.Pointer {
	f, ok := getFunc(self).(func(a *Tuple) (Object, error))
	if !ok {
		raise(AssertionError.Err("callWithArgs: wrong function type!!!"))
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
	f, ok := getFunc(self).(func(a *Tuple, k *Dict) (Object, error))
	if !ok {
		raise(AssertionError.Err("callWithKeywords: wrong function type!!!"))
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
