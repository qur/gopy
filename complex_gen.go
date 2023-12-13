// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Complex represents objects of the ComplexType (or PyComplex_Type
// in the Python API) type.
//
// This type implements the Number protocol.
type Complex struct {
	o C.PyComplexObject
}

var _ Object = (*Complex)(nil)

// ComplexType is the Type object that represents the Complex type.
var ComplexType = newType(&C.PyComplex_Type)

func complexCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.complexCheck(c(obj)) != 0
}

// AsComplex casts the given obj to a Complex (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Complex, then nil is returned.
//
// Return value: Borrowed Reference.
func AsComplex(obj Object) *Complex {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.complexCheck(o) == 0 {
		return nil
	}
	return (*Complex)(unsafe.Pointer(o))
}

func newComplex(obj *C.PyObject) *Complex {
	return (*Complex)(unsafe.Pointer(obj))
}



func (co *Complex) c() *C.PyComplexObject {
	return (*C.PyComplexObject)(unsafe.Pointer(co))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (co *Complex) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(co))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (co *Complex) Type() *Type {
	return newType(c(co).ob_type)
}


// Decref decrements co's reference count, co may not be nil.
func (co *Complex) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(co))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	if *refcnt == C._Py_IMMORTAL_REFCNT {
		return
	}
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments co's reference count, co may not be nil.
func (co *Complex) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(co)).anon0[0]))
	if *refcnt == C._Py_IMMORTAL_REFCNT {
		return
	}
	*refcnt++
}

// Repr returns a String representation of "co". This is equivalent to the
// Python "repr(co)".
//
// Return value: New Reference.
func (co *Complex) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(co))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of co. The equivalent
// Python is "hash(co)".
func (co *Complex) Hash() (int, error) {
	ret := C.PyObject_Hash(c(co))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// HasAttr returns true if "co" has the attribute "name".  This is equivalent
// to the Python "hasattr(co, name)".
func (co *Complex) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(co), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "co" with the name "name".  This is
// equivalent to the Python "co.name".
//
// Return value: New Reference.
func (co *Complex) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(co), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "co" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "co op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (co *Complex) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(co), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "co" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "co op obj",
// where op is the corresponding Python operator for op.
func (co *Complex) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(co), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// AsNumberMethods returns a NumberMethods instance that refers to the same
// underlying Python object as co.
//
// This method also means that Complex implements the NumberProtocol
// interface.
func (co *Complex) AsNumberMethods() *NumberMethods {
	return (*NumberMethods)(unsafe.Pointer(co.Base()))
}

// Add returns the result of adding co and obj. The equivalent Python is
// "co + obj".
//
// Return value: New Reference.
func (co *Complex) Add(obj Object) (Object, error) {
	ret := C.PyNumber_Add(c(co), c(obj))
	return obj2ObjErr(ret)
}

// Subtract returns the result of subtracting obj from co. The equivalent Python
// is "co - obj".
//
// Return value: New Reference.
func (co *Complex) Subtract(obj Object) (Object, error) {
	ret := C.PyNumber_Subtract(c(co), c(obj))
	return obj2ObjErr(ret)
}

// Multiply returns the result of multiplying co by obj. The equivalent Python
// is "co * obj".
//
// Return value: New Reference.
func (co *Complex) Multiply(obj Object) (Object, error) {
	ret := C.PyNumber_Multiply(c(co), c(obj))
	return obj2ObjErr(ret)
}

// TrueDivide returns the approximate result of dividing co by obj. The result is
// approximate due to the limited representational accuracy of binary floating
// point numbers. The equivalent Python is "co / obj".
//
// Return value: New Reference.
func (co *Complex) TrueDivide(obj Object) (Object, error) {
	ret := C.PyNumber_TrueDivide(c(co), c(obj))
	return obj2ObjErr(ret)
}

// Power returns the result of the Python "pow(co, obj1, obj2)", where
// obj2 is optional.
//
// Return value: New Reference.
func (co *Complex) Power(obj1, obj2 Object) (Object, error) {
	if obj2 == nil {
		None.Incref()
		obj2 = None
	}
	ret := C.PyNumber_Power(c(co), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// Negative returns the negation of co. The equivalent Python is "-co".
//
// Return value: New Reference.
func (co *Complex) Negative() (Object, error) {
	ret := C.PyNumber_Negative(c(co))
	return obj2ObjErr(ret)
}

// Positive returns the positive of co. The equivalent Python is "+co".
//
// Return value: New Reference.
func (co *Complex) Positive() (Object, error) {
	ret := C.PyNumber_Positive(c(co))
	return obj2ObjErr(ret)
}

// Absolute returns the absolute value of co. The equivalent Python is "abs(co)".
//
// Return value: New Reference.
func (co *Complex) Absolute() (Object, error) {
	ret := C.PyNumber_Absolute(c(co))
	return obj2ObjErr(ret)
}



/*
set fields:
  nb_absolute
  nb_add
  nb_bool
  nb_multiply
  nb_negative
  nb_positive
  nb_power
  nb_subtract
  nb_true_divide
  ob_base
  tp_alloc
  tp_as_number
  tp_basicsize
  tp_doc
  tp_flags
  tp_free
  tp_getattro
  tp_hash
  tp_members
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
*/

