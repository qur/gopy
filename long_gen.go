// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Long represents objects of the LongType (or PyLongType
// in the Python API) type.
//
// This type implements the Number protocol.
type Long struct {
	abstractObject
	o C.PyLongObject
}

// LongType is the Type object that represents the Long type.
var LongType = (*Type)(unsafe.Pointer(&C.PyLong_Type))

func longCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.longCheck(c(obj)) != 0
}

func newLong(obj *C.PyObject) *Long {
	return (*Long)(unsafe.Pointer(obj))
}

// Repr returns a String representation of "l". This is equivalent to the
// Python "repr(l)".
//
// Return value: New Reference.
func (l *Long) Repr() (Object, error) {
	ret := C.PyObject_Repr(c(l))
	return obj2ObjErr(ret)
}

// Hash computes and returns the hash value of l. The equivalent
// Python is "hash(l)".
func (l *Long) Hash() (int, error) {
	ret := C.PyObject_Hash(c(l))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// HasAttr returns true if "l" has the attribute "name".  This is equivalent
// to the Python "hasattr(l, name)".
func (l *Long) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(l), c(name))
	if ret == 1 {
		return true
	}
	return false
}

// GetAttr returns the attribute of "l" with the name "name".  This is
// equivalent to the Python "l.name".
//
// Return value: New Reference.
func (l *Long) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(l), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "l" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "l op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (l *Long) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(l), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "obj" with "obj2" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "obj op obj2",
// where op is the corresponding Python operator for op.
func (l *Long) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(l), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// AsNumber returns a NumberMethods instance that refers to the same underlying
// Python object as l.
//
// This method also means that Long implements the Number interface.
func (l *Long) AsNumber() *NumberMethods {
	return (*NumberMethods)(unsafe.Pointer(l.Base()))
}

// Add returns the result of adding l and obj. The equivalent Python is
// "l + obj".
//
// Return value: New Reference.
func (l *Long) Add(obj Object) (Object, error) {
	ret := C.PyNumber_Add(c(l), c(obj))
	return obj2ObjErr(ret)
}

// Subtract returns the result of subtracting obj from l. The equivalent Python
// is "l - obj".
//
// Return value: New Reference.
func (l *Long) Subtract(obj Object) (Object, error) {
	ret := C.PyNumber_Subtract(c(l), c(obj))
	return obj2ObjErr(ret)
}

// Multiply returns the result of multiplying l by obj. The equivalent Python
// is "l * obj".
//
// Return value: New Reference.
func (l *Long) Multiply(obj Object) (Object, error) {
	ret := C.PyNumber_Multiply(c(l), c(obj))
	return obj2ObjErr(ret)
}

// FloorDivide returns the floor of dividing l by obj. The equivalent Python is
// "l // obj".
//
// Return value: New Reference.
func (l *Long) FloorDivide(obj Object) (Object, error) {
	ret := C.PyNumber_FloorDivide(c(l), c(obj))
	return obj2ObjErr(ret)
}

// TrueDivide returns the approximate result of dividing l by obj. The result is
// approximate due to the limited representational accuracy of binary floating
// point numbers. The equivalent Python is "l / obj".
//
// Return value: New Reference.
func (l *Long) TrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_TrueDivide(c(l), c(obj))
	return obj2ObjErr(ret)
}

// Remainder returns the remainder of dividing l by obj. The equivalent Python
// is "l % obj".
//
// Return value: New Reference.
func (l *Long) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(c(l), c(obj))
	return obj2ObjErr(ret)
}

// Divmod returns the result of the Python "divmod(l, obj)".
//
// Return value: New Reference.
func (l *Long) Divmod(obj Object) (Object, error) {
	ret := C.PyNumber_Divmod(c(l), c(obj))
	return obj2ObjErr(ret)
}

// Power returns the result of the Python "pow(l, obj1, obj2)", where
// obj2 is optional.
//
// Return value: New Reference.
func (l *Long) Power(obj1, obj2 Object) (Object, error) {
	if obj2 == nil {
		None.Incref()
		obj2 = None
	}
	ret := C.PyNumber_Power(c(l), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// Negative returns the negation of l. The equivalent Python is "-l".
//
// Return value: New Reference.
func (l *Long) Negative() (Object, error) {
	ret := C.PyNumber_Negative(c(l))
	return obj2ObjErr(ret)
}

// Positive returns the positive of l. The equivalent Python is "+l".
//
// Return value: New Reference.
func (l *Long) Positive() (Object, error) {
	ret := C.PyNumber_Positive(c(l))
	return obj2ObjErr(ret)
}

// Absolute returns the absolute value of l. The equivalent Python is "abs(l)".
//
// Return value: New Reference.
func (l *Long) Absolute() (Object, error) {
	ret := C.PyNumber_Absolute(c(l))
	return obj2ObjErr(ret)
}

// Invert returns the bitwise negation of l. The equivalent Python is "-l".
//
// Return value: New Reference.
func (l *Long) Invert() (Object, error) {
	ret := C.PyNumber_Invert(c(l))
	return obj2ObjErr(ret)
}

// LShift returns the result of left shifting l by obj. The equivalent Python
// is "l << obj".
//
// Return value: New Reference.
func (l *Long) LShift(obj Object) (Object, error) {
	ret := C.PyNumber_Lshift(c(l), c(obj))
	return obj2ObjErr(ret)
}

// RShift returns the result of right shifting l by obj. The equivalent Python
// is "l << obj".
//
// Return value: New Reference.
func (l *Long) RShift(obj Object) (Object, error) {
	ret := C.PyNumber_Rshift(c(l), c(obj))
	return obj2ObjErr(ret)
}

// And returns the bitwise and of l and obj. The equivalent Python is
// "l & obj".
//
// Return value: New Reference.
func (l *Long) And(obj Object) (Object, error) {
	ret := C.PyNumber_And(c(l), c(obj))
	return obj2ObjErr(ret)
}

// Xor returns the bitwise xor of l and obj. The equivalent Python is
// "l ^ obj".
//
// Return value: New Reference.
func (l *Long) Xor(obj Object) (Object, error) {
	ret := C.PyNumber_Xor(c(l), c(obj))
	return obj2ObjErr(ret)
}

// Or returns the bitwise or of l and obj. The equivalent Python is
// "l | obj".
//
// Return value: New Reference.
func (l *Long) Or(obj Object) (Object, error) {
	ret := C.PyNumber_Or(c(l), c(obj))
	return obj2ObjErr(ret)
}

func (l *Long) Float() (*Float, error) {
	ret := C.PyNumber_Float(c(l))
	return newFloat(ret), exception()
}

func (l *Long) Index() (*Long, error) {
	ret := C.PyNumber_Index(c(l))
	return newLong(ret), exception()
}

func (l *Long) ToBase(base int) (*Long, error) {
	ret := C.PyNumber_ToBase(c(l), C.int(base))
	return newLong(ret), exception()
}

func (l *Long) AsInt(exc *ExceptionClass) (int, error) {
	ret := C.PyNumber_AsSsize_t(c(l), c(exc))
	return int(ret), exception()
}



/*
set fields:
  nb_absolute
  nb_add
  nb_and
  nb_bool
  nb_divmod
  nb_float
  nb_floor_divide
  nb_index
  nb_int
  nb_invert
  nb_lshift
  nb_multiply
  nb_negative
  nb_or
  nb_positive
  nb_power
  nb_remainder
  nb_rshift
  nb_subtract
  nb_true_divide
  nb_xor
  ob_base
  tp_as_number
  tp_basicsize
  tp_doc
  tp_flags
  tp_free
  tp_getattro
  tp_getset
  tp_hash
  tp_itemsize
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
*/

