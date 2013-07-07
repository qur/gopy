// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

//export goClassCall
func goClassCall(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, *Tuple, *Dict) (Object, error))(unsafe.Pointer(&ctxt.call))

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	ret, err := (*f)(obj, a, k)
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
	} else {
		(*AbstractObject)(obj).Free()
	}
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

//export goClassInit
func goClassInit(obj, args, kwds unsafe.Pointer) int {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, *Tuple, *Dict) error)(unsafe.Pointer(&ctxt.init))

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	err := (*f)(obj, a, k)
	if err != nil {
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

//export goClassRepr
func goClassRepr(obj unsafe.Pointer) unsafe.Pointer {
	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) string)(unsafe.Pointer(&ctxt.repr))

	s := C.CString((*f)(obj))
	defer C.free(unsafe.Pointer(s))

	return unsafe.Pointer(C.PyString_FromString(s))
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

	return unsafe.Pointer(C.PyString_FromString(s))
}
