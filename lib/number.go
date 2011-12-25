// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

// NumberProtocol is a 0-sized type that can be embedded in concrete types after
// the AbstractObject to provide access to the suite of methods that Python
// calls the "Number Protocol".
type NumberProtocol struct{}

// Number is an interface that defines the Python "Number Protocol".
type Number interface {
	Object
	Add(obj Object) (Object, error)
	Subtract(obj Object) (Object, error)
	Multiply(obj Object) (Object, error)
	Divide(obj Object) (Object, error)
	FloorDivide(obj Object) (Object, error)
	TrueDivide(obj Object) (Object, error)
	Remainder(obj Object) (Object, error)
	Divmod(obj Object) (Object, error)
	Power(obj1, obj2 Object) (Object, error)
	Negative() (Object, error)
	Positive() (Object, error)
	Absolute() (Object, error)
	Invert() (Object, error)
	Lshift(obj Object) (Object, error)
	Rshift(obj Object) (Object, error)
	And(obj Object) (Object, error)
	Xor(obj Object) (Object, error)
	Or(obj Object) (Object, error)
	InPlaceAdd(obj Object) (Object, error)
	InPlaceSubtract(obj Object) (Object, error)
	InPlaceMultiply(obj Object) (Object, error)
	InPlaceDivide(obj Object) (Object, error)
	InPlaceFloorDivide(obj Object) (Object, error)
	InPlaceTrueDivide(obj Object) (Object, error)
	InPlaceRemainder(obj Object) (Object, error)
	InPlacePower(obj1, obj2 Object) (Object, error)
	InPlaceLshift(obj Object) (Object, error)
	InPlaceRshift(obj Object) (Object, error)
	InPlaceAnd(obj Object) (Object, error)
	InPlaceXor(obj Object) (Object, error)
	InPlaceOr(obj Object) (Object, error)
}

// number is a concrete realisation of the Number Protocol.  A type that
// implements the "Number Protocol" but doesn't embed NumberProtocol can be
// turned into a Number by calling AsNumber.
type number struct {
	AbstractObject
	NumberProtocol
	o C.PyObject
}

func cnp(n *NumberProtocol) *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(n))
}

// AsNumber returns a struct pointer that satisfies the Number interface.  It
// will refer to the same underlying object as obj.  If obj doesn't implement
// the "Number Protocol", then nil is returned.
func AsNumber(obj Object) Number {
	if C.PyNumber_Check(c(obj)) != 1 {
		return nil
	}
	if n, ok := obj.(Number); ok {
		return n
	}
	return (*number)(unsafe.Pointer(obj.Base()))
}

// Add returns the result of adding n and obj.  The equivalent Python is "n +
// obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Add(obj Object) (Object, error) {
	ret := C.PyNumber_Add(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// Subtract returns the result of subtracting obj from n.  The equivalent Python
// is "n - obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Subtract(obj Object) (Object, error) {
	ret := C.PyNumber_Subtract(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// Multiply returns the result of multiplying n by obj.  The equivalent Python
// is "n * obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Multiply(obj Object) (Object, error) {
	ret := C.PyNumber_Multiply(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// Divide returns the result of dividing n by obj.  The equivalent Python is "n
// / obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Divide(obj Object) (Object, error) {
	ret := C.PyNumber_Divide(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// FloorDivide returns the floor of dividing n obj obj.
//
// Return value: New Reference.
func (n *NumberProtocol) FloorDivide(obj Object) (Object, error) {
	ret := C.PyNumber_FloorDivide(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// TrueDivide returns the ... TODO
//
// Return value: New Reference.
func (n *NumberProtocol) TrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_TrueDivide(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// Remainder returns the remainder of dividing n by obj.  The equivalent Python
// is "n % obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// Divmod returns the result of the Python "divmod(n, obj)".
//
// Return value: New Reference.
func (n *NumberProtocol) Divmod(obj Object) (Object, error) {
	ret := C.PyNumber_Divmod(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// Power returns the result of the Python "pow(n, obj1, obj2)".
//
// Return value: New Reference.
func (n *NumberProtocol) Power(obj1, obj2 Object) (Object, error) {
	ret := C.PyNumber_Power(cnp(n), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// Negative returns the negation of n.  The equivalent Python is "-n".
//
// Return value: New Reference.
func (n *NumberProtocol) Negative() (Object, error) {
	ret := C.PyNumber_Negative(cnp(n))
	return obj2ObjErr(ret)
}

// Positive returns the positive of n.  The equivalent Python is "+n".
//
// Return value: New Reference.
func (n *NumberProtocol) Positive() (Object, error) {
	ret := C.PyNumber_Positive(cnp(n))
	return obj2ObjErr(ret)
}

// Absolute returns the absolute value of n.  The equivalent Python is "abs(n)".
//
// Return value: New Reference.
func (n *NumberProtocol) Absolute() (Object, error) {
	ret := C.PyNumber_Absolute(cnp(n))
	return obj2ObjErr(ret)
}

// Invert returns the bitwise negation of n.  The equivalent Python is "-n".
//
// Return value: New Reference.
func (n *NumberProtocol) Invert() (Object, error) {
	ret := C.PyNumber_Invert(cnp(n))
	return obj2ObjErr(ret)
}

// Lshift returns the result of left shifting n by obj.  The equivalent Python
// is "n << obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Lshift(obj Object) (Object, error) {
	ret := C.PyNumber_Lshift(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// Rshift returns the result of right shifting n by obj.  The equivalent Python
// is "n << obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Rshift(obj Object) (Object, error) {
	ret := C.PyNumber_Rshift(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// And returns the bitwise and of n and obj.  The equivalent Python is "n &
// obj".
//
// Return value: New Reference.
func (n *NumberProtocol) And(obj Object) (Object, error) {
	ret := C.PyNumber_And(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// Xor returns the bitwise xor of n and obj.  The equivalent Python is "n ^
// obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Xor(obj Object) (Object, error) {
	ret := C.PyNumber_Xor(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// Or returns the bitwise or of n and obj.  The equivalent Python is "n | obj".
//
// Return value: New Reference.
func (n *NumberProtocol) Or(obj Object) (Object, error) {
	ret := C.PyNumber_Or(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceAdd returns the result of adding n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n += obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceAdd(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceAdd(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceSubtract returns the result of subtracting obj from n.  This is done
// in place if supported by n.  The equivalent Python is "n -= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceSubtract(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceSubtract(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceMultiply returns the result of multiplying n by obj.  This is done in
// place if supported by n.  The equivalent Python is "n *= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceMultiply(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceMultiply(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceDivide returns the result of dividing n by obj.  This is done in place
// if supported by n.  The equivalent Python is "n /= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceDivide(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceDivide(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// TODO returns the ...
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceFloorDivide(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceFloorDivide(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// TODO returns the ...
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceTrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceTrueDivide(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceRemainder returns the remainder of n divided by obj.  This is done in
// place if supported by n.  The equivalent Python is "n %= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceRemainder(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceRemainder(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlacePower returns the result of the Python "pow(n, obj1, obj2)".  This is
// done in place if supported by n.  If obj2 is None, then the Python "n **=
// obj" is also equivalent, if obj2 is not None, there is no equivalent in
// Python.
//
// Return value: New Reference.
func (n *NumberProtocol) InPlacePower(obj1, obj2 Object) (Object, error) {
	ret := C.PyNumber_InPlacePower(cnp(n), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// InPlaceLshift returns the result of left shifting n by obj.  This is done in
// place if supported by n.  The equivalent Python is "n <<= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceLshift(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceLshift(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceRshift returns the result of right shifting n by obj.  This is done in
// place if supported by n.  The equivalent Python is "n >>= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceRshift(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceRshift(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceAnd returns the bitwise and of n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n &= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceAnd(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceAnd(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceXor returns the bitwise xor of n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n ^= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceXor(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceXor(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceOr returns the bitwise or of n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n |= obj".
//
// Return value: New Reference.
func (n *NumberProtocol) InPlaceOr(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceOr(cnp(n), c(obj))
	return obj2ObjErr(ret)
}

// PyNumber_Coerce: TODO

// PyNumber_CoerceEx: TODO

// PyNumber_Int: TODO

// PyNumber_Long: TODO

// PyNumber_Float: TODO

// PyNumber_Index: TODO

// PyNumber_ToBase: TODO

// PyNumber_AsSsize_t: TODO
