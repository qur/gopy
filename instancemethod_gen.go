// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// InstanceMethod represents objects of the InstanceMethodType (or PyInstanceMethodType
// in the Python API) type.
type InstanceMethod struct {
	o C.PyInstanceMethodObject
}

var _ Object = (*InstanceMethod)(nil)

// InstanceMethodType is the Type object that represents the InstanceMethod type.
var InstanceMethodType = (*Type)(unsafe.Pointer(&C.PyInstanceMethod_Type))

func instanceMethodCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.instanceMethodCheck(c(obj)) != 0
}

func newInstanceMethod(obj *C.PyObject) *InstanceMethod {
	return (*InstanceMethod)(unsafe.Pointer(obj))
}

func (i *InstanceMethod) c() *C.PyInstanceMethodObject {
	return (*C.PyInstanceMethodObject)(unsafe.Pointer(i))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (i *InstanceMethod) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(i))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (i *InstanceMethod) Type() *Type {
	return newType((*C.PyObject)(unsafe.Pointer(c(i).ob_type)))
}

// Decref decrements i's reference count, i may not be nil.
func (i *InstanceMethod) Decref() {
	C.decref(c(i))
}

// Incref increments i's reference count, i may not be nil.
func (i *InstanceMethod) Incref() {
	C.incref(c(i))
}

// IsTrue returns true if the value of i is considered to be True. This is
// equivalent to "if i:" in Python.
func (i *InstanceMethod) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(i))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of i is considered to be False. This is
// equivalent to "if not i:" in Python.
func (i *InstanceMethod) Not() bool {
	ret := C.PyObject_Not(c(i))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Free deallocates the storage (in Python) for i. After calling this method,
// i should no longer be used.
func (i *InstanceMethod) Free() {
	free(i)
}

// Repr returns a String representation of "i". This is equivalent to the
// Python "repr(i)".
//
// Return value: New Reference.
func (i *InstanceMethod) Repr() (Object, error) {
	ret := C.PyObject_Repr(c(i))
	return obj2ObjErr(ret)
}

// Call calls i with the given args and kwds. kwds may be nil, args may not
// (an empty Tuple must be used if no arguments are wanted). Returns the result
// of the call, or an Error on failure.  This is equivalent to
// "i(*args, **kwds)" in Python.
//
// Return value: New Reference.
func (i *InstanceMethod) Call(args *Tuple, kwds *Dict) (Object, error) {
	ret := C.PyObject_Call(c(i), c(args), c(kwds))
	return obj2ObjErr(ret)
}

// CallGo calls i with the given args and kwds, either may be nil. Returns the
// result of the call, or an Error on failure.  This is equivalent to
// "i(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (i *InstanceMethod) CallGo(args []any, kwds map[string]any) (Object, error) {
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
	ret := C.PyObject_Call(c(i), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// HasAttr returns true if "i" has the attribute "name".  This is equivalent
// to the Python "hasattr(i, name)".
func (i *InstanceMethod) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(i), c(name))
	if ret == 1 {
		return true
	}
	return false
}

// GetAttr returns the attribute of "i" with the name "name".  This is
// equivalent to the Python "i.name".
//
// Return value: New Reference.
func (i *InstanceMethod) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(i), c(name))
	return obj2ObjErr(ret)
}

// SetAttr sets the attribute of "i" with the name "name" to "value".  This is
// equivalent to the Python "i.name = value".
func (i *InstanceMethod) SetAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c(i), c(name), c(value))
	return int2Err(ret)
}

// DelAttr deletes the attribute with the name "name" from "i".  This is
// equivalent to the Python "del i.name".
func (i *InstanceMethod) DelAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c(i), c(name), nil)
	return int2Err(ret)
}

// RichCompare compares "i" with "obj" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "i op obj", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (i *InstanceMethod) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(i), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "obj" with "obj2" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "obj op obj2",
// where op is the corresponding Python operator for op.
func (i *InstanceMethod) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(i), c(obj), C.int(op))
	return int2BoolErr(ret)
}



/*
set fields:
  ob_base
  tp_basicsize
  tp_call
  tp_dealloc
  tp_descr_get
  tp_doc
  tp_flags
  tp_getattro
  tp_getset
  tp_members
  tp_name
  tp_new
  tp_repr
  tp_richcompare
  tp_setattro
  tp_traverse
*/

