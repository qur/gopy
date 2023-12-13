// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Method represents objects of the MethodType (or PyMethod_Type
// in the Python API) type.
type Method struct {
	o C.PyMethodObject
}

var _ Object = (*Method)(nil)

// MethodType is the Type object that represents the Method type.
var MethodType = newType(&C.PyMethod_Type)

func methodCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.methodCheck(c(obj)) != 0
}

// AsMethod casts the given obj to a Method (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Method, then nil is returned.
//
// Return value: Borrowed Reference.
func AsMethod(obj Object) *Method {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.methodCheck(o) == 0 {
		return nil
	}
	return (*Method)(unsafe.Pointer(o))
}

func newMethod(obj *C.PyObject) *Method {
	return (*Method)(unsafe.Pointer(obj))
}



func (m *Method) c() *C.PyMethodObject {
	return (*C.PyMethodObject)(unsafe.Pointer(m))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (m *Method) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(m))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (m *Method) Type() *Type {
	return newType(c(m).ob_type)
}

// Decref decrements m's reference count, m may not be nil.
func (m *Method) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(m))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments m's reference count, m may not be nil.
func (m *Method) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(m)).anon0[0]))
	*refcnt++
}

// Repr returns a String representation of "m". This is equivalent to the
// Python "repr(m)".
//
// Return value: New Reference.
func (m *Method) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(m))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of m. The equivalent
// Python is "hash(m)".
func (m *Method) Hash() (int, error) {
	ret := C.PyObject_Hash(c(m))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// Call calls m with the given args and kwds. kwds may be nil, args may not
// (an empty Tuple must be used if no arguments are wanted). Returns the result
// of the call, or an Error on failure.  This is equivalent to
// "m(*args, **kwds)" in Python.
//
// Return value: New Reference.
func (m *Method) Call(args *Tuple, kwds *Dict) (Object, error) {
	ret := C.PyObject_Call(c(m), c(args), c(kwds))
	return obj2ObjErr(ret)
}

// CallGo calls m with the given args and kwds, either may be nil. Returns the
// result of the call, or an Error on failure.  This is equivalent to
// "m(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (m *Method) CallGo(args []any, kwds map[string]any) (Object, error) {
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
	ret := C.PyObject_Call(c(m), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// HasAttr returns true if "m" has the attribute "name".  This is equivalent
// to the Python "hasattr(m, name)".
func (m *Method) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(m), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "m" with the name "name".  This is
// equivalent to the Python "m.name".
//
// Return value: New Reference.
func (m *Method) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(m), c(name))
	return obj2ObjErr(ret)
}

// SetAttr sets the attribute of "m" with the name "name" to "value".  This is
// equivalent to the Python "m.name = value".
func (m *Method) SetAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c(m), c(name), c(value))
	return int2Err(ret)
}

// DelAttr deletes the attribute with the name "name" from "m".  This is
// equivalent to the Python "del m.name".
func (m *Method) DelAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c(m), c(name), nil)
	return int2Err(ret)
}

// RichCompare compares "m" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "m op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (m *Method) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(m), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "m" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "m op obj",
// where op is the corresponding Python operator for op.
func (m *Method) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(m), c(obj), C.int(op))
	return int2BoolErr(ret)
}



/*
set fields:
  ob_base
  tp_basicsize
  tp_call
  tp_dealloc
  tp_doc
  tp_flags
  tp_getattro
  tp_getset
  tp_hash
  tp_members
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_setattro
  tp_traverse
  tp_vectorcall_offset
  tp_weaklistoffset
*/

