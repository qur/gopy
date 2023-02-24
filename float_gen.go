// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Float represents objects of the FloatType (or PyFloatType
// in the Python API) type.
//
// This type implements the Number protocol.
type Float struct {
	o C.PyFloatObject
}

var _ Object = (*Float)(nil)

// FloatType is the Type object that represents the Float type.
var FloatType = (*Type)(unsafe.Pointer(&C.PyFloat_Type))

func floatCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.floatCheck(c(obj)) != 0
}

func newFloat(obj *C.PyObject) *Float {
	return (*Float)(unsafe.Pointer(obj))
}

func (f *Float) c() *C.PyFloatObject {
	return (*C.PyFloatObject)(unsafe.Pointer(f))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (f *Float) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(f))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (f *Float) Type() *Type {
	return newType((*C.PyObject)(unsafe.Pointer(c(f).ob_type)))
}

// Decref decrements f's reference count, f may not be nil.
func (f *Float) Decref() {
	C.decref(c(f))
}

// Incref increments f's reference count, f may not be nil.
func (f *Float) Incref() {
	C.incref(c(f))
}

// IsTrue returns true if the value of f is considered to be True. This is
// equivalent to "if f:" in Python.
func (f *Float) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(f))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of f is considered to be False. This is
// equivalent to "if not f:" in Python.
func (f *Float) Not() bool {
	ret := C.PyObject_Not(c(f))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Free deallocates the storage (in Python) for f. After calling this method,
// f should no longer be used.
func (f *Float) Free() {
	free(f)
}

// Repr returns a String representation of "f". This is equivalent to the
// Python "repr(f)".
//
// Return value: New Reference.
func (f *Float) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(f))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of f. The equivalent
// Python is "hash(f)".
func (f *Float) Hash() (int, error) {
	ret := C.PyObject_Hash(c(f))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// HasAttr returns true if "f" has the attribute "name".  This is equivalent
// to the Python "hasattr(f, name)".
func (f *Float) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(f), c(name))
	if ret == 1 {
		return true
	}
	return false
}

// GetAttr returns the attribute of "f" with the name "name".  This is
// equivalent to the Python "f.name".
//
// Return value: New Reference.
func (f *Float) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(f), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "f" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "f op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (f *Float) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(f), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "obj" with "obj2" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "obj op obj2",
// where op is the corresponding Python operator for op.
func (f *Float) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(f), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// AsNumberMethods returns a NumberMethods instance that refers to the same
// underlying Python object as f.
//
// This method also means that Float implements the NumberProtocol
// interface.
func (f *Float) AsNumberMethods() *NumberMethods {
	return (*NumberMethods)(unsafe.Pointer(f.Base()))
}

// Add returns the result of adding f and obj. The equivalent Python is
// "f + obj".
//
// Return value: New Reference.
func (f *Float) Add(obj Object) (Object, error) {
	ret := C.PyNumber_Add(c(f), c(obj))
	return obj2ObjErr(ret)
}

// Subtract returns the result of subtracting obj from f. The equivalent Python
// is "f - obj".
//
// Return value: New Reference.
func (f *Float) Subtract(obj Object) (Object, error) {
	ret := C.PyNumber_Subtract(c(f), c(obj))
	return obj2ObjErr(ret)
}

// Multiply returns the result of multiplying f by obj. The equivalent Python
// is "f * obj".
//
// Return value: New Reference.
func (f *Float) Multiply(obj Object) (Object, error) {
	ret := C.PyNumber_Multiply(c(f), c(obj))
	return obj2ObjErr(ret)
}

// FloorDivide returns the floor of dividing f by obj. The equivalent Python is
// "f // obj".
//
// Return value: New Reference.
func (f *Float) FloorDivide(obj Object) (Object, error) {
	ret := C.PyNumber_FloorDivide(c(f), c(obj))
	return obj2ObjErr(ret)
}

// TrueDivide returns the approximate result of dividing f by obj. The result is
// approximate due to the limited representational accuracy of binary floating
// point numbers. The equivalent Python is "f / obj".
//
// Return value: New Reference.
func (f *Float) TrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_TrueDivide(c(f), c(obj))
	return obj2ObjErr(ret)
}

// Remainder returns the remainder of dividing f by obj. The equivalent Python
// is "f % obj".
//
// Return value: New Reference.
func (f *Float) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(c(f), c(obj))
	return obj2ObjErr(ret)
}

// Divmod returns the result of the Python "divmod(f, obj)".
//
// Return value: New Reference.
func (f *Float) Divmod(obj Object) (Object, error) {
	ret := C.PyNumber_Divmod(c(f), c(obj))
	return obj2ObjErr(ret)
}

// Power returns the result of the Python "pow(f, obj1, obj2)", where
// obj2 is optional.
//
// Return value: New Reference.
func (f *Float) Power(obj1, obj2 Object) (Object, error) {
	if obj2 == nil {
		None.Incref()
		obj2 = None
	}
	ret := C.PyNumber_Power(c(f), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// Negative returns the negation of f. The equivalent Python is "-f".
//
// Return value: New Reference.
func (f *Float) Negative() (Object, error) {
	ret := C.PyNumber_Negative(c(f))
	return obj2ObjErr(ret)
}

// Positive returns the positive of f. The equivalent Python is "+f".
//
// Return value: New Reference.
func (f *Float) Positive() (Object, error) {
	ret := C.PyNumber_Positive(c(f))
	return obj2ObjErr(ret)
}

// Absolute returns the absolute value of f. The equivalent Python is "abs(f)".
//
// Return value: New Reference.
func (f *Float) Absolute() (Object, error) {
	ret := C.PyNumber_Absolute(c(f))
	return obj2ObjErr(ret)
}

func (f *Float) Long() (*Long, error) {
	ret := C.PyNumber_Long(c(f))
	return newLong(ret), exception()
}



/*
set fields:
  nb_absolute
  nb_add
  nb_bool
  nb_divmod
  nb_float
  nb_floor_divide
  nb_int
  nb_multiply
  nb_negative
  nb_positive
  nb_power
  nb_remainder
  nb_subtract
  nb_true_divide
  ob_base
  tp_as_number
  tp_basicsize
  tp_dealloc
  tp_doc
  tp_flags
  tp_getattro
  tp_getset
  tp_hash
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_vectorcall
*/

