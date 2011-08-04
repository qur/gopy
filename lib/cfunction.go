// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int cfunctionCheck(PyObject *o) { return PyCFunction_Check(o); }
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type CFunction struct {
	BaseObject
}

func cfunctionCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.cfunctionCheck(c(obj)) != 0
}

func newCFunction(obj *C.PyObject) *CFunction {
	return (*CFunction)(unsafe.Pointer(obj))
}

func CFunction_New(name string, fn interface{}, doc string) (*CFunction, os.Error) {
	ml := C.newMethodDef()

	switch f := fn.(type) {

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
		C.free(unsafe.Pointer(ml))
		return nil, fmt.Errorf("CFunction_New: unknown func type for %s", name)

	}

	ret := C.PyCFunction_NewEx(ml, saveFunc(fn), nil)
	if ret == nil {
		C.free(unsafe.Pointer(ml))
		return nil, exception()
	}

	ml.ml_name = C.CString(name)
	ml.ml_doc = C.CString(doc)

	return newCFunction(ret), nil
}

// PyCFunction_GetFunction

func (f *CFunction) Self() (Object, os.Error) {
	ret := C.PyCFunction_GetSelf(c(f))
	if ret == nil {
		return nil, exception()
	}
	return newBaseObject(ret).actual(), nil
}

func (f *CFunction) Flags() (int, os.Error) {
	ret := C.PyCFunction_GetFlags(c(f))
	return int(ret), exception()
}

func (f *CFunction) Call(args *Tuple, kw *Dict) (Object, os.Error) {
	ret := C.PyCFunction_Call(c(f), c(args), c(kw))
	if ret == nil {
		return nil, exception()
	}
	return newBaseObject(ret).actual(), nil
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

//export callWithoutArgs
func callWithoutArgs(self, args unsafe.Pointer) unsafe.Pointer {
	_idx := (*C.PyObject)(self)
	idx := C.PyInt_AsLong(_idx)
	f, ok := funcs[idx].(func() (Object, os.Error))
	if !ok {
		fmt.Printf("invalid index: %d\n", idx)
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

