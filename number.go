package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Number is an interface that is implemented by types that implement the
// Python "Number Protocol".
type Number interface {
	Object
	AsNumber() *NumberMethods
}

type Index interface {
	Object
	Index() (*Long, error)
}

// NumberMethods is a concrete realisation of the full set of Number Protocol
// methods.  A type that implements the "Number Protocol" can be turned into a
// NumberMethods instance using AsNumber.
//
// Note that the methods not already implemented on the type itself may return
// an error, as not all methods are implemented by all types that support the
// protocol.
type NumberMethods struct {
	abstractObject
	o C.PyObject
}

// AsNumber returns a NumberMethods instance that refers to the same
// underlying Python object as obj. If obj doesn't implement the "Number
// Protocol" (i.e. the Number interface), then nil is returned.
func AsNumber(obj Object) *NumberMethods {
	if n, ok := obj.(Number); ok {
		return n.AsNumber()
	}
	if C.numberCheck(c(obj)) > 0 {
		return (*NumberMethods)(unsafe.Pointer(obj.Base()))
	}
	return nil
}

// Add returns the result of adding n and obj.  The equivalent Python is "n +
// obj".
//
// Return value: New Reference.
func (n *NumberMethods) Add(obj Object) (Object, error) {
	ret := C.PyNumber_Add(c(n), c(obj))
	return obj2ObjErr(ret)
}

// Subtract returns the result of subtracting obj from n.  The equivalent Python
// is "n - obj".
//
// Return value: New Reference.
func (n *NumberMethods) Subtract(obj Object) (Object, error) {
	ret := C.PyNumber_Subtract(c(n), c(obj))
	return obj2ObjErr(ret)
}

// Multiply returns the result of multiplying n by obj.  The equivalent Python
// is "n * obj".
//
// Return value: New Reference.
func (n *NumberMethods) Multiply(obj Object) (Object, error) {
	ret := C.PyNumber_Multiply(c(n), c(obj))
	return obj2ObjErr(ret)
}

// FloorDivide returns the floor of dividing n by obj. The equivalent Python is
// "n // obj".
//
// Return value: New Reference.
func (n *NumberMethods) FloorDivide(obj Object) (Object, error) {
	ret := C.PyNumber_FloorDivide(c(n), c(obj))
	return obj2ObjErr(ret)
}

// TrueDivide returns the approximate result of dividing n by obj. The result is
// approximate due to the limited representational accuracy of binary floating
// point numbers. The equivalent Python is "n / obj".
//
// Return value: New Reference.
func (n *NumberMethods) TrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_TrueDivide(c(n), c(obj))
	return obj2ObjErr(ret)
}

// Remainder returns the remainder of dividing n by obj.  The equivalent Python
// is "n % obj".
//
// Return value: New Reference.
func (n *NumberMethods) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(c(n), c(obj))
	return obj2ObjErr(ret)
}

// Divmod returns the result of the Python "divmod(n, obj)".
//
// Return value: New Reference.
func (n *NumberMethods) Divmod(obj Object) (Object, error) {
	ret := C.PyNumber_Divmod(c(n), c(obj))
	return obj2ObjErr(ret)
}

// Power returns the result of the Python "pow(n, obj1, obj2)".
//
// Return value: New Reference.
func (n *NumberMethods) Power(obj1, obj2 Object) (Object, error) {
	ret := C.PyNumber_Power(c(n), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// Negative returns the negation of n.  The equivalent Python is "-n".
//
// Return value: New Reference.
func (n *NumberMethods) Negative() (Object, error) {
	ret := C.PyNumber_Negative(c(n))
	return obj2ObjErr(ret)
}

// Positive returns the positive of n.  The equivalent Python is "+n".
//
// Return value: New Reference.
func (n *NumberMethods) Positive() (Object, error) {
	ret := C.PyNumber_Positive(c(n))
	return obj2ObjErr(ret)
}

// Absolute returns the absolute value of n.  The equivalent Python is "abs(n)".
//
// Return value: New Reference.
func (n *NumberMethods) Absolute() (Object, error) {
	ret := C.PyNumber_Absolute(c(n))
	return obj2ObjErr(ret)
}

// Invert returns the bitwise negation of n.  The equivalent Python is "-n".
//
// Return value: New Reference.
func (n *NumberMethods) Invert() (Object, error) {
	ret := C.PyNumber_Invert(c(n))
	return obj2ObjErr(ret)
}

// LShift returns the result of left shifting n by obj.  The equivalent Python
// is "n << obj".
//
// Return value: New Reference.
func (n *NumberMethods) LShift(obj Object) (Object, error) {
	ret := C.PyNumber_Lshift(c(n), c(obj))
	return obj2ObjErr(ret)
}

// RShift returns the result of right shifting n by obj.  The equivalent Python
// is "n << obj".
//
// Return value: New Reference.
func (n *NumberMethods) RShift(obj Object) (Object, error) {
	ret := C.PyNumber_Rshift(c(n), c(obj))
	return obj2ObjErr(ret)
}

// And returns the bitwise and of n and obj.  The equivalent Python is "n &
// obj".
//
// Return value: New Reference.
func (n *NumberMethods) And(obj Object) (Object, error) {
	ret := C.PyNumber_And(c(n), c(obj))
	return obj2ObjErr(ret)
}

// Xor returns the bitwise xor of n and obj.  The equivalent Python is "n ^
// obj".
//
// Return value: New Reference.
func (n *NumberMethods) Xor(obj Object) (Object, error) {
	ret := C.PyNumber_Xor(c(n), c(obj))
	return obj2ObjErr(ret)
}

// Or returns the bitwise or of n and obj.  The equivalent Python is "n | obj".
//
// Return value: New Reference.
func (n *NumberMethods) Or(obj Object) (Object, error) {
	ret := C.PyNumber_Or(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceAdd returns the result of adding n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n += obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceAdd(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceAdd(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceSubtract returns the result of subtracting obj from n.  This is done
// in place if supported by n.  The equivalent Python is "n -= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceSubtract(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceSubtract(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceMultiply returns the result of multiplying n by obj.  This is done in
// place if supported by n.  The equivalent Python is "n *= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceMultiply(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceMultiply(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceFloorDivide returns the floor of dividing n by obj. This is done in
// place if supported by n. The equivalent Python is "n //= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceFloorDivide(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceFloorDivide(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceTrueDivide returns the approximate result of dividing n by obj. This
// is done in place if supported by n. The result is approximate due to the
// limited representational accuracy of binary floating point numbers. The
// equivalent Python is "n /= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceTrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceTrueDivide(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceRemainder returns the remainder of n divided by obj.  This is done in
// place if supported by n.  The equivalent Python is "n %= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceRemainder(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceRemainder(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlacePower returns the result of the Python "pow(n, obj1, obj2)".  This is
// done in place if supported by n.  If obj2 is None, then the Python "n **=
// obj" is also equivalent, if obj2 is not None, there is no equivalent in
// Python.
//
// Return value: New Reference.
func (n *NumberMethods) InPlacePower(obj1, obj2 Object) (Object, error) {
	ret := C.PyNumber_InPlacePower(c(n), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// InPlaceLShift returns the result of left shifting n by obj.  This is done in
// place if supported by n.  The equivalent Python is "n <<= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceLShift(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceLshift(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceRShift returns the result of right shifting n by obj.  This is done in
// place if supported by n.  The equivalent Python is "n >>= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceRShift(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceRshift(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceAnd returns the bitwise and of n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n &= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceAnd(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceAnd(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceXor returns the bitwise xor of n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n ^= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceXor(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceXor(c(n), c(obj))
	return obj2ObjErr(ret)
}

// InPlaceOr returns the bitwise or of n and obj.  This is done in place if
// supported by n.  The equivalent Python is "n |= obj".
//
// Return value: New Reference.
func (n *NumberMethods) InPlaceOr(obj Object) (Object, error) {
	ret := C.PyNumber_InPlaceOr(c(n), c(obj))
	return obj2ObjErr(ret)
}

// PyNumber_Long: TODO
func (n *NumberMethods) Long() (*Long, error) {
	ret := C.PyNumber_Long(c(n))
	return newLong(ret), exception()
}

// PyNumber_Float: TODO
func (n *NumberMethods) Float() (*Float, error) {
	ret := C.PyNumber_Float(c(n))
	return newFloat(ret), exception()
}

// PyNumber_Index: TODO
func (n *NumberMethods) Index() (*Long, error) {
	ret := C.PyNumber_Index(c(n))
	return newLong(ret), exception()
}

// PyNumber_ToBase: TODO
func (n *NumberMethods) ToBase(base int) (*Long, error) {
	ret := C.PyNumber_ToBase(c(n), C.int(base))
	return newLong(ret), exception()
}

// PyNumber_AsSsize_t: TODO
func (n *NumberMethods) AsInt(exc *ExceptionClass) (int, error) {
	ret := C.PyNumber_AsSsize_t(c(n), c(exc))
	return int(ret), exception()
}
