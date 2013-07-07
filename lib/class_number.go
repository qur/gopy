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

//export GoClassNumAdd
func GoClassNumAdd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_add, obj1, obj2)
}

//export GoClassNumSubtract
func GoClassNumSubtract(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_subtract, obj1, obj2)
}

//export GoClassNumMultiply
func GoClassNumMultiply(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_multiply, obj1, obj2)
}

//export GoClassNumDivide
func GoClassNumDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_divide, obj1, obj2)
}

//export GoClassNumRemainder
func GoClassNumRemainder(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_remainder, obj1, obj2)
}

//export GoClassNumDivmod
func GoClassNumDivmod(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_divmod, obj1, obj2)
}

//export GoClassNumPower
func GoClassNumPower(obj1, obj2, obj3 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallTernary(ctxt.nb_power, obj1, obj2, obj3)
}

//export GoClassNumNegative
func GoClassNumNegative(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_negative, obj)
}

//export GoClassNumPositive
func GoClassNumPositive(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_positive, obj)
}

//export GoClassNumAbsolute
func GoClassNumAbsolute(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_absolute, obj)
}

//export GoClassNumNonzero
func GoClassNumNonzero(obj unsafe.Pointer) int {
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

//export GoClassNumInvert
func GoClassNumInvert(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_invert, obj)
}

//export GoClassNumLshift
func GoClassNumLshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_lshift, obj1, obj2)
}

//export GoClassNumRshift
func GoClassNumRshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_rshift, obj1, obj2)
}

//export GoClassNumAnd
func GoClassNumAnd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_and, obj1, obj2)
}

//export GoClassNumXor
func GoClassNumXor(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_xor, obj1, obj2)
}

//export GoClassNumOr
func GoClassNumOr(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_or, obj1, obj2)
}

//nb_coerce

//export GoClassNumInt
func GoClassNumInt(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_int, obj)
}

//export GoClassNumLong
func GoClassNumLong(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_long, obj)
}

//export GoClassNumFloat
func GoClassNumFloat(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_float, obj)
}

//export GoClassNumOct
func GoClassNumOct(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_oct, obj)
}

//export GoClassNumHex
func GoClassNumHex(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_hex, obj)
}

//export GoClassNumInplaceAdd
func GoClassNumInplaceAdd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_add, obj1, obj2)
}

//export GoClassNumInplaceSubtract
func GoClassNumInplaceSubtract(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_subtract, obj1, obj2)
}

//export GoClassNumInplaceMultiply
func GoClassNumInplaceMultiply(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_multiply, obj1, obj2)
}

//export GoClassNumInplaceDivide
func GoClassNumInplaceDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_divide, obj1, obj2)
}

//export GoClassNumInplaceRemainder
func GoClassNumInplaceRemainder(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_remainder, obj1, obj2)
}

//export GoClassNumInplacePower
func GoClassNumInplacePower(obj1, obj2, obj3 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallTernary(ctxt.nb_ip_power, obj1, obj2, obj3)
}

//export GoClassNumInplaceLshift
func GoClassNumInplaceLshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_lshift, obj1, obj2)
}

//export GoClassNumInplaceRshift
func GoClassNumInplaceRshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_rshift, obj1, obj2)
}

//export GoClassNumInplaceAnd
func GoClassNumInplaceAnd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_and, obj1, obj2)
}

//export GoClassNumInplaceXor
func GoClassNumInplaceXor(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_xor, obj1, obj2)
}

//export GoClassNumInplaceOr
func GoClassNumInplaceOr(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_or, obj1, obj2)
}

//export GoClassNumFloorDivide
func GoClassNumFloorDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_floordiv, obj1, obj2)
}

//export GoClassNumTrueDivide
func GoClassNumTrueDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_truediv, obj1, obj2)
}

//export GoClassNumInplaceFloorDivide
func GoClassNumInplaceFloorDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_floordiv, obj1, obj2)
}

//export GoClassNumInplaceTrueDivide
func GoClassNumInplaceTrueDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj1)
	return cnCallBinary(ctxt.nb_ip_truediv, obj1, obj2)
}

//export GoClassNumIndex
func GoClassNumIndex(obj unsafe.Pointer) unsafe.Pointer {
	ctxt := getClassContext(obj)
	return cnCallUnary(ctxt.nb_index, obj)
}
