// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

type tpCall interface {
	PyCall(*Tuple, *Dict) (Object, error)
}

//export goClassCall
func goClassCall(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Turn obj into the ClassObject instead of the proxy, and it should
	// implement tpCall.
	co := newObject((*C.PyObject)(obj)).(tpCall)

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	ret, err := co.PyCall(a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCompare
func goClassCompare(obj1, obj2 unsafe.Pointer) int {
	// Get the class context
	ctxt := getClassContext(obj1)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object) (int, error))(unsafe.Pointer(&ctxt.compare))

	o := newObject((*C.PyObject)(obj2))

	ret, err := (*f)(obj1, o)
	if err != nil {
		raise(err)
		return -1
	}

	return ret
}

//export goClassGetAttr
func goClassGetAttr(obj unsafe.Pointer, name *C.char) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, string) (Object, error))(unsafe.Pointer(&ctxt.getattr))

	s := C.GoString(name)

	ret, err := (*f)(obj, s)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassGetAttrObj
func goClassGetAttrObj(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj1)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object) (Object, error))(unsafe.Pointer(&ctxt.getattro))

	o := newObject((*C.PyObject)(obj2))

	ret, err := (*f)(obj1, o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassDealloc
func goClassDealloc(obj unsafe.Pointer) {
	// Get the class context
	ctxt := getClassContext(obj)

	if ctxt.dealloc != nil {
		// Turn the function into something we can call
		f := (*func(unsafe.Pointer))(unsafe.Pointer(&ctxt.dealloc))

		(*f)(obj)
	}

	// we always want Python to _actually_ free the object, any registered hook
	// should just be tidying things up on the Go side.
	(*AbstractObject)(obj).Free()
}

//export goClassHash
func goClassHash(obj unsafe.Pointer) C.long {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) (uint32, error))(unsafe.Pointer(&ctxt.hash))

	ret, err := (*f)(obj)
	if err != nil {
		raise(err)
		return -1
	} else if C.long(ret) == -1 {
		return -2
	}

	return C.long(ret)
}

type tpInit interface {
	PyInit(*Tuple, *Dict) error
}

//export goClassInit
func goClassInit(obj, args, kwds unsafe.Pointer) int {
	// Turn obj into the ClassObject instead of the proxy, and it should
	// implement tpInit.
	co := newObject((*C.PyObject)(obj)).(tpInit)

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	if err := co.PyInit(a, k); err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export goClassIter
func goClassIter(obj unsafe.Pointer) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) (Object, error))(unsafe.Pointer(&ctxt.iter))

	ret, err := (*f)(obj)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassIterNext
func goClassIterNext(obj unsafe.Pointer) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) (Object, error))(unsafe.Pointer(&ctxt.iternext))

	ret, err := (*f)(obj)
	if err != nil {
		raise(err)
		return nil
	} else if ret == nil {
		return nil
	}

	return unsafe.Pointer(c(ret))
}

type tpRepr interface {
	PyRepr() string
}

//export goClassRepr
func goClassRepr(obj unsafe.Pointer) unsafe.Pointer {
	// Turn obj into the ClassObject instead of the proxy, and it should
	// implement tpRepr.
	co := newObject((*C.PyObject)(obj)).(tpRepr)

	s := C.CString(co.PyRepr())
	defer C.free(unsafe.Pointer(s))

	return unsafe.Pointer(C.PyUnicode_FromString(s))
}

//export goClassRichCmp
func goClassRichCmp(obj1, obj2 unsafe.Pointer, op int) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj1)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object, Op) (Object, error))(unsafe.Pointer(&ctxt.richcmp))

	// Get obj2 ready for use
	arg := newObject((*C.PyObject)(obj2))

	ret, err := (*f)(obj1, arg, Op(op))
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassSetAttr
func goClassSetAttr(obj unsafe.Pointer, name *C.char, obj2 unsafe.Pointer) int {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, string, Object) error)(unsafe.Pointer(&ctxt.setattr))

	s := C.GoString(name)
	o := newObject((*C.PyObject)(obj2))

	err := (*f)(obj, s, o)
	if err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export goClassSetAttrObj
func goClassSetAttrObj(obj1, obj2, obj3 unsafe.Pointer) int {
	// Get the class context
	ctxt := getClassContext(obj1)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object, Object) error)(unsafe.Pointer(&ctxt.setattro))

	o := newObject((*C.PyObject)(obj2))
	o2 := newObject((*C.PyObject)(obj3))

	err := (*f)(obj1, o, o2)
	if err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export goClassStr
func goClassStr(obj unsafe.Pointer) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) string)(unsafe.Pointer(&ctxt.str))

	s := C.CString((*f)(obj))
	defer C.free(unsafe.Pointer(s))

	return unsafe.Pointer(C.PyUnicode_FromString(s))
}
