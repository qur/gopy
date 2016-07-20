// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

func csCallSsizeArgFunc(fn, obj unsafe.Pointer, s C.Py_ssize_t) unsafe.Pointer {
	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, int64) (Object, error))(unsafe.Pointer(&fn))

	ret, err := (*f)(obj, int64(s))
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassSeqLength
func goClassSeqLength(obj unsafe.Pointer) C.Py_ssize_t {
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) int64)(unsafe.Pointer(&ctxt.sq_length))

	return C.Py_ssize_t((*f)(obj))
}

//export goClassSeqConcat
func goClassSeqConcat(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.sq_concat, obj1, obj2)
}

//export goClassSeqRepeat
func goClassSeqRepeat(obj unsafe.Pointer, count C.Py_ssize_t) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return csCallSsizeArgFunc(ctxt.sq_repeat, obj, count)
}

//export goClassSeqGetItem
func goClassSeqGetItem(obj unsafe.Pointer, idx C.Py_ssize_t) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return csCallSsizeArgFunc(ctxt.sq_get, obj, idx)
}

//export goClassSeqSetItem
func goClassSeqSetItem(obj unsafe.Pointer, i C.Py_ssize_t, v unsafe.Pointer) int {
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, int64, Object) error)(unsafe.Pointer(&ctxt.sq_set))

	// Get v ready to use
	a := newObject((*C.PyObject)(v))

	err := (*f)(obj, int64(i), a)
	if err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export goClassSeqContains
func goClassSeqContains(obj1, obj2 unsafe.Pointer) int {
	ctxt := getClassContext(obj1)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object) (bool, error))(unsafe.Pointer(&ctxt.sq_contains))

	// Get obj2 ready to use
	a := newObject((*C.PyObject)(obj2))

	ret, err := (*f)(obj1, a)
	if err != nil {
		raise(err)
		return -1
	}

	if ret {
		return 1
	}

	return 0
}

//export goClassSeqIpConcat
func goClassSeqIpConcat(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.sq_ip_concat, obj1, obj2)
}

//export goClassSeqIpRepeat
func goClassSeqIpRepeat(obj unsafe.Pointer, count C.Py_ssize_t) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return csCallSsizeArgFunc(ctxt.sq_ip_repeat, obj, count)
}
