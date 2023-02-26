package py

// #include "utils.h"
import "C"

import (
	"sync"
	"unsafe"
)

func NewCFunction(name string, fn interface{}, doc string) (*CFunction, error) {
	return makeCFunction(name, fn, doc, nil)
}

func makeCFunction(name string, fn interface{}, doc string, mod_name *C.PyObject) (*CFunction, error) {
	ml := C.newMethodDef()

	switch fn.(type) {

	case func() (Object, error):
		C.set_call_noargs(&ml.ml_meth)
		ml.ml_flags = C.METH_NOARGS

	case func(Object) (Object, error):
		C.set_call_single(&ml.ml_meth)
		ml.ml_flags = C.METH_O

	case func(*Tuple) (Object, error):
		C.set_call_args(&ml.ml_meth)
		ml.ml_flags = C.METH_VARARGS

	case func(*Tuple, *Dict) (Object, error):
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

// GoMethod describes a Go function to be used in Python. Name is the name of
// the function in Python, Doc is the docstring to be used, and Func is the
// function to be called.
//
// Multiple signatures for the function are supported:
//
//	func() (py.Object, error)
//	func(py.Object) (py.Object, error)
//	func(*py.Tuple) (py.Object, error)
//	func(*py.Tuple, *py.Dict) (py.Object, error)
//
// These map to the following behaviours:
//
//   - A function with no arguments
//   - A function with a single argument, passed directly
//   - A function with positional arguments only, passed as a Tuple
//   - A function with positional and keyword arguments, passed as a Tuple and
//     Dict
type GoMethod struct {
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
	return C.PyLong_FromLong(C.long(len(funcs) - 1))
}

func getFunc(self *C.PyObject) interface{} {
	funcLock.RLock()
	defer funcLock.RUnlock()

	idx := int(C.PyLong_AsLong(self))

	if idx >= len(funcs) {
		return nil
	}

	return funcs[idx]
}

//export callWithoutArgs
func callWithoutArgs(self, args *C.PyObject) *C.PyObject {
	f, ok := getFunc(self).(func() (Object, error))
	if !ok {
		raise(AssertionError.Err("callWithoutArgs: wrong function type!!!"))
		return nil
	}
	return ce(f())
}

//export callWithSingle
func callWithSingle(self, arg *C.PyObject) *C.PyObject {
	f, ok := getFunc(self).(func(a Object) (Object, error))
	if !ok {
		raise(AssertionError.Err("callWithArgs: wrong function type!!!"))
		return nil
	}
	return ce(f(newObject(arg)))
}

//export callWithArgs
func callWithArgs(self, args *C.PyObject) *C.PyObject {
	f, ok := getFunc(self).(func(a *Tuple) (Object, error))
	if !ok {
		raise(AssertionError.Err("callWithArgs: wrong function type!!!"))
		return nil
	}
	return ce(f(newTuple(args)))
}

//export callWithKeywords
func callWithKeywords(self, args, kw *C.PyObject) *C.PyObject {
	f, ok := getFunc(self).(func(a *Tuple, k *Dict) (Object, error))
	if !ok {
		raise(AssertionError.Err("callWithKeywords: wrong function type!!!"))
		return nil
	}
	return ce(f(newTuple(args), newDict(kw)))
}
