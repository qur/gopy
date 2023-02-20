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
	abstractObject
	o C.PyFloatObject
}

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

func (f *Float) AsNumber() *NumberMethods {
	return (*NumberMethods)(unsafe.Pointer(f.Base()))
}

func (f *Float) Add(obj Object) (Object, error) {
	ret := C.PyNumber_Add(c(f), c(obj))
	return obj2ObjErr(ret)
}

func (f *Float) Subtract(obj Object) (Object, error) {
	ret := C.PyNumber_Subtract(c(f), c(obj))
	return obj2ObjErr(ret)
}

func (f *Float) Multiply(obj Object) (Object, error) {
	ret := C.PyNumber_Multiply(c(f), c(obj))
	return obj2ObjErr(ret)
}

func (f *Float) FloorDivide(obj Object) (Object, error) {
	ret := C.PyNumber_FloorDivide(c(f), c(obj))
	return obj2ObjErr(ret)
}

func (f *Float) TrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_TrueDivide(c(f), c(obj))
	return obj2ObjErr(ret)
}

func (f *Float) Remainder(obj Object) (Object, error) {
	ret := C.PyNumber_Remainder(c(f), c(obj))
	return obj2ObjErr(ret)
}

func (f *Float) Divmod(obj Object) (Object, error) {
	ret := C.PyNumber_Divmod(c(f), c(obj))
	return obj2ObjErr(ret)
}

func (f *Float) Power(obj1, obj2 Object) (Object, error) {
	ret := C.PyNumber_Power(c(f), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

func (f *Float) Negative() (Object, error) {
	ret := C.PyNumber_Negative(c(f))
	return obj2ObjErr(ret)
}

func (f *Float) Positive() (Object, error) {
	ret := C.PyNumber_Positive(c(f))
	return obj2ObjErr(ret)
}

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

