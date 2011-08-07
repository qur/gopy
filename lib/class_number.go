// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

//export goClassNumAdd
func goClassNumAdd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumAdd"))
	return nil
}

//export goClassNumSubtract
func goClassNumSubtract(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumSubtract"))
	return nil
}

//export goClassNumMultiply
func goClassNumMultiply(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumMultiply"))
	return nil
}

//export goClassNumDivide
func goClassNumDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumDivide"))
	return nil
}

//export goClassNumRemainder
func goClassNumRemainder(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumRemainder"))
	return nil
}

//export goClassNumDivmod
func goClassNumDivmod(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumDivmod"))
	return nil
}

//export goClassNumPower
func goClassNumPower(obj1, obj2, obj3 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumPower"))
	return nil
}

//export goClassNumNegative
func goClassNumNegative(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumNegative"))
	return nil
}

//export goClassNumPositive
func goClassNumPositive(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumPositive"))
	return nil
}

//export goClassNumAbsolute
func goClassNumAbsolute(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumAbsolute"))
	return nil
}

//export goClassNumNonzero
func goClassNumNonzero(obj unsafe.Pointer) int {
	raise(NotImplemented("goClassNumNonzero"))
	return -1
}

//export goClassNumInvert
func goClassNumInvert(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInvert"))
	return nil
}

//export goClassNumLshift
func goClassNumLshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumLshift"))
	return nil
}

//export goClassNumRshift
func goClassNumRshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumRshift"))
	return nil
}

//export goClassNumAnd
func goClassNumAnd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumAnd"))
	return nil
}

//export goClassNumXor
func goClassNumXor(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumXor"))
	return nil
}

//export goClassNumOr
func goClassNumOr(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumOr"))
	return nil
}

//nb_coerce

//export goClassNumInt
func goClassNumInt(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInt"))
	return nil
}

//export goClassNumLong
func goClassNumLong(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumLong"))
	return nil
}

//export goClassNumFloat
func goClassNumFloat(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumFloat"))
	return nil
}

//export goClassNumOct
func goClassNumOct(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumOct"))
	return nil
}

//export goClassNumHex
func goClassNumHex(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumHex"))
	return nil
}

//export goClassNumInplaceAdd
func goClassNumInplaceAdd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceAdd"))
	return nil
}

//export goClassNumInplaceSubtract
func goClassNumInplaceSubtract(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceSubtract"))
	return nil
}

//export goClassNumInplaceMultiply
func goClassNumInplaceMultiply(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceMultiply"))
	return nil
}

//export goClassNumInplaceDivide
func goClassNumInplaceDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceDivide"))
	return nil
}

//export goClassNumInplaceRemainder
func goClassNumInplaceRemainder(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceRemainder"))
	return nil
}

//export goClassNumInplacePower
func goClassNumInplacePower(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplacePower"))
	return nil
}

//export goClassNumInplaceLshift
func goClassNumInplaceLshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceLshift"))
	return nil
}

//export goClassNumInplaceRshift
func goClassNumInplaceRshift(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceRshift"))
	return nil
}

//export goClassNumInplaceAnd
func goClassNumInplaceAnd(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceAnd"))
	return nil
}

//export goClassNumInplaceXor
func goClassNumInplaceXor(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceXor"))
	return nil
}

//export goClassNumInplaceOr
func goClassNumInplaceOr(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceOr"))
	return nil
}

//export goClassNumFloorDivide
func goClassNumFloorDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNum"))
	return nil
}

//export goClassNumTrueDivide
func goClassNumTrueDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNum"))
	return nil
}

//export goClassNumInplaceFloorDivide
func goClassNumInplaceFloorDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceFloorDivide"))
	return nil
}

//export goClassNumInplaceTrueDivide
func goClassNumInplaceTrueDivide(obj1, obj2 unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumInplaceTrueDivide"))
	return nil
}

//export goClassNumIndex
func goClassNumIndex(obj unsafe.Pointer) unsafe.Pointer {
	raise(NotImplemented("goClassNumIndex"))
	return nil
}
