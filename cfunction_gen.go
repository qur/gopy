// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// CFunction represents objects of the CFunctionType (or PyCFunctionType
// in the Python API) type.
type CFunction struct {
	o C.PyCFunctionObject
}

var _ Object = (*CFunction)(nil)

// CFunctionType is the Type object that represents the CFunction type.
var CFunctionType = newType(&C.PyCFunction_Type)

func cFunctionCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.cFunctionCheck(c(obj)) != 0
}

func newCFunction(obj *C.PyObject) *CFunction {
	return (*CFunction)(unsafe.Pointer(obj))
}



func (cf *CFunction) c() *C.PyCFunctionObject {
	return (*C.PyCFunctionObject)(unsafe.Pointer(cf))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (cf *CFunction) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(cf))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (cf *CFunction) Type() *Type {
	return newType(c(cf).ob_type)
}

// Decref decrements cf's reference count, cf may not be nil.
func (cf *CFunction) Decref() {
	C.decref(c(cf))
}

// Incref increments cf's reference count, cf may not be nil.
func (cf *CFunction) Incref() {
	C.incref(c(cf))
}

func (cf *CFunction) raw() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(cf))
}

// Repr returns a String representation of "cf". This is equivalent to the
// Python "repr(cf)".
//
// Return value: New Reference.
func (cf *CFunction) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(cf))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of cf. The equivalent
// Python is "hash(cf)".
func (cf *CFunction) Hash() (int, error) {
	ret := C.PyObject_Hash(c(cf))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// Call calls cf with the given args and kwds. kwds may be nil, args may not
// (an empty Tuple must be used if no arguments are wanted). Returns the result
// of the call, or an Error on failure.  This is equivalent to
// "cf(*args, **kwds)" in Python.
//
// Return value: New Reference.
func (cf *CFunction) Call(args *Tuple, kwds *Dict) (Object, error) {
	ret := C.PyObject_Call(c(cf), c(args), c(kwds))
	return obj2ObjErr(ret)
}

// CallGo calls cf with the given args and kwds, either may be nil. Returns the
// result of the call, or an Error on failure.  This is equivalent to
// "cf(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (cf *CFunction) CallGo(args []any, kwds map[string]any) (Object, error) {
	obj1, err := NewTupleFromValues(args...)
	if err != nil {
		return nil, err
	}
	defer obj1.Decref()
	obj2, err := NewDictFromValuesString(kwds)
	if err != nil {
		return nil, err
	}
	defer obj2.Decref()
	ret := C.PyObject_Call(c(cf), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// HasAttr returns true if "cf" has the attribute "name".  This is equivalent
// to the Python "hasattr(cf, name)".
func (cf *CFunction) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(cf), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "cf" with the name "name".  This is
// equivalent to the Python "cf.name".
//
// Return value: New Reference.
func (cf *CFunction) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(cf), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "cf" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "cf op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (cf *CFunction) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(cf), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "cf" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "cf op obj",
// where op is the corresponding Python operator for op.
func (cf *CFunction) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(cf), c(obj), C.int(op))
	return int2BoolErr(ret)
}



/*
set fields:
  ob_base
  tp_basicsize
  tp_call
  tp_dealloc
  tp_flags
  tp_getattro
  tp_getset
  tp_hash
  tp_members
  tp_methods
  tp_name
  tp_repr
  tp_richcompare
  tp_traverse
  tp_vectorcall_offset
  tp_weaklistoffset
*/

