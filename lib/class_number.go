// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

func cnCallUnary(fn, obj unsafe.Pointer) unsafe.Pointer {
	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) (Object, error))(unsafe.Pointer(&fn))

	ret, err := (*f)(obj)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

func cnCallBinary(fn, obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object) (Object, error))(unsafe.Pointer(&fn))

	// Get obj2 ready to use
	arg := newObject((*C.PyObject)(obj2))

	ret, err := (*f)(obj1, arg)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

func cnCallTernary(fn, obj1, obj2, obj3 unsafe.Pointer) unsafe.Pointer {
	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object, Object) (Object, error))(unsafe.Pointer(&fn))

	// Get obj2 and obj3 ready to use
	arg1 := newObject((*C.PyObject)(obj2))
	arg2 := newObject((*C.PyObject)(obj3))

	ret, err := (*f)(obj1, arg1, arg2)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassNumAdd
func goClassNumAdd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_add, obj1, obj2)
}

//export goClassNumSubtract
func goClassNumSubtract(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_subtract, obj1, obj2)
}

//export goClassNumMultiply
func goClassNumMultiply(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_multiply, obj1, obj2)
}

//export goClassNumDivide
func goClassNumDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_divide, obj1, obj2)
}

//export goClassNumRemainder
func goClassNumRemainder(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_remainder, obj1, obj2)
}

//export goClassNumDivmod
func goClassNumDivmod(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_divmod, obj1, obj2)
}

//export goClassNumPower
func goClassNumPower(obj1, obj2, obj3 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallTernary(ctxt.nb_power, obj1, obj2, obj3)
}

//export goClassNumNegative
func goClassNumNegative(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_negative, obj)
}

//export goClassNumPositive
func goClassNumPositive(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_positive, obj)
}

//export goClassNumAbsolute
func goClassNumAbsolute(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_absolute, obj)
}

//export goClassNumNonzero
func goClassNumNonzero(obj unsafe.Pointer) int {
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) (bool, error))(unsafe.Pointer(&ctxt.nb_nonzero))

	ret, err := (*f)(obj)
	if err != nil {
		raise(err)
		return -1
	}

	if ret {
		return 1
	}

	return 0
}

//export goClassNumInvert
func goClassNumInvert(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_invert, obj)
}

//export goClassNumLshift
func goClassNumLshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_lshift, obj1, obj2)
}

//export goClassNumRshift
func goClassNumRshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_rshift, obj1, obj2)
}

//export goClassNumAnd
func goClassNumAnd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_and, obj1, obj2)
}

//export goClassNumXor
func goClassNumXor(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_xor, obj1, obj2)
}

//export goClassNumOr
func goClassNumOr(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_or, obj1, obj2)
}

//nb_coerce

//export goClassNumInt
func goClassNumInt(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_int, obj)
}

//export goClassNumLong
func goClassNumLong(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_long, obj)
}

//export goClassNumFloat
func goClassNumFloat(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_float, obj)
}

//export goClassNumOct
func goClassNumOct(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_oct, obj)
}

//export goClassNumHex
func goClassNumHex(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_hex, obj)
}

//export goClassNumInplaceAdd
func goClassNumInplaceAdd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_add, obj1, obj2)
}

//export goClassNumInplaceSubtract
func goClassNumInplaceSubtract(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_subtract, obj1, obj2)
}

//export goClassNumInplaceMultiply
func goClassNumInplaceMultiply(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_multiply, obj1, obj2)
}

//export goClassNumInplaceDivide
func goClassNumInplaceDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_divide, obj1, obj2)
}

//export goClassNumInplaceRemainder
func goClassNumInplaceRemainder(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_remainder, obj1, obj2)
}

//export goClassNumInplacePower
func goClassNumInplacePower(obj1, obj2, obj3 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallTernary(ctxt.nb_ip_power, obj1, obj2, obj3)
}

//export goClassNumInplaceLshift
func goClassNumInplaceLshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_lshift, obj1, obj2)
}

//export goClassNumInplaceRshift
func goClassNumInplaceRshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_rshift, obj1, obj2)
}

//export goClassNumInplaceAnd
func goClassNumInplaceAnd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_and, obj1, obj2)
}

//export goClassNumInplaceXor
func goClassNumInplaceXor(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_xor, obj1, obj2)
}

//export goClassNumInplaceOr
func goClassNumInplaceOr(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_or, obj1, obj2)
}

//export goClassNumFloorDivide
func goClassNumFloorDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_floordiv, obj1, obj2)
}

//export goClassNumTrueDivide
func goClassNumTrueDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_truediv, obj1, obj2)
}

//export goClassNumInplaceFloorDivide
func goClassNumInplaceFloorDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_floordiv, obj1, obj2)
}

//export goClassNumInplaceTrueDivide
func goClassNumInplaceTrueDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_truediv, obj1, obj2)
}

//export goClassNumIndex
func goClassNumIndex(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_index, obj)
}
