// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Function represents objects of the FunctionType (or PyFunction_Type
// in the Python API) type.
type Function struct {
	o C.PyFunctionObject
}

var _ Object = (*Function)(nil)

// FunctionType is the Type object that represents the Function type.
var FunctionType = newType(&C.PyFunction_Type)

func functionCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.functionCheck(c(obj)) != 0
}

// AsFunction casts the given obj to a Function (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Function, then nil is returned.
//
// Return value: Borrowed Reference.
func AsFunction(obj Object) *Function {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.functionCheck(o) == 0 {
		return nil
	}
	return (*Function)(unsafe.Pointer(o))
}

func newFunction(obj *C.PyObject) *Function {
	return (*Function)(unsafe.Pointer(obj))
}



func (f *Function) c() *C.PyFunctionObject {
	return (*C.PyFunctionObject)(unsafe.Pointer(f))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (f *Function) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(f))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (f *Function) Type() *Type {
	return newType(c(f).ob_type)
}

// Decref decrements f's reference count, f may not be nil.
func (f *Function) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(f))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments f's reference count, f may not be nil.
func (f *Function) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(f)).anon0[0]))
	*refcnt++
}

// Repr returns a String representation of "f". This is equivalent to the
// Python "repr(f)".
//
// Return value: New Reference.
func (f *Function) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(f))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Call calls f with the given args and kwds. kwds may be nil, args may not
// (an empty Tuple must be used if no arguments are wanted). Returns the result
// of the call, or an Error on failure.  This is equivalent to
// "f(*args, **kwds)" in Python.
//
// Return value: New Reference.
func (f *Function) Call(args *Tuple, kwds *Dict) (Object, error) {
	ret := C.PyObject_Call(c(f), c(args), c(kwds))
	return obj2ObjErr(ret)
}

// CallGo calls f with the given args and kwds, either may be nil. Returns the
// result of the call, or an Error on failure.  This is equivalent to
// "f(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (f *Function) CallGo(args []any, kwds map[string]any) (Object, error) {
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
	ret := C.PyObject_Call(c(f), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}



/*
set fields:
  ob_base
  tp_basicsize
  tp_call
  tp_clear
  tp_dealloc
  tp_descr_get
  tp_dictoffset
  tp_doc
  tp_flags
  tp_getset
  tp_members
  tp_name
  tp_new
  tp_repr
  tp_traverse
  tp_vectorcall_offset
  tp_weaklistoffset
*/

