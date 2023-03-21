// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Type represents objects of the TypeType (or PyType_Type
// in the Python API) type.
type Type struct {
	o C.PyTypeObject
}

var _ Object = (*Type)(nil)

// TypeType is the Type object that represents the Type type.
var TypeType = newType(&C.PyType_Type)

func typeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.typeCheck(c(obj)) != 0
}

// AsType casts the given obj to a Type (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Type, then nil is returned.
//
// Return value: Borrowed Reference.
func AsType(obj Object) *Type {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.typeCheck(o) == 0 {
		return nil
	}
	return (*Type)(unsafe.Pointer(o))
}

func newType(obj *C.PyTypeObject) *Type {
	return (*Type)(unsafe.Pointer(obj))
}

func (t *Type) c() *C.PyTypeObject {
	return (*C.PyTypeObject)(unsafe.Pointer(t))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (t *Type) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(t))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (t *Type) Type() *Type {
	return newType(c(t).ob_type)
}

// Decref decrements t's reference count, t may not be nil.
func (t *Type) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(t))
	obj.ob_refcnt--
	if obj.ob_refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments t's reference count, t may not be nil.
func (t *Type) Incref() {
	obj := (*C.PyObject)(unsafe.Pointer(t))
	obj.ob_refcnt++
}

// Repr returns a String representation of "t". This is equivalent to the
// Python "repr(t)".
//
// Return value: New Reference.
func (t *Type) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(t))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Call calls t with the given args and kwds. kwds may be nil, args may not
// (an empty Tuple must be used if no arguments are wanted). Returns the result
// of the call, or an Error on failure.  This is equivalent to
// "t(*args, **kwds)" in Python.
//
// Return value: New Reference.
func (t *Type) Call(args *Tuple, kwds *Dict) (Object, error) {
	ret := C.PyObject_Call(c(t), c(args), c(kwds))
	return obj2ObjErr(ret)
}

// CallGo calls t with the given args and kwds, either may be nil. Returns the
// result of the call, or an Error on failure.  This is equivalent to
// "t(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (t *Type) CallGo(args []any, kwds map[string]any) (Object, error) {
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
	ret := C.PyObject_Call(c(t), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// HasAttr returns true if "t" has the attribute "name".  This is equivalent
// to the Python "hasattr(t, name)".
func (t *Type) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(t), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "t" with the name "name".  This is
// equivalent to the Python "t.name".
//
// Return value: New Reference.
func (t *Type) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(t), c(name))
	return obj2ObjErr(ret)
}

// SetAttr sets the attribute of "t" with the name "name" to "value".  This is
// equivalent to the Python "t.name = value".
func (t *Type) SetAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c(t), c(name), c(value))
	return int2Err(ret)
}

// DelAttr deletes the attribute with the name "name" from "t".  This is
// equivalent to the Python "del t.name".
func (t *Type) DelAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c(t), c(name), nil)
	return int2Err(ret)
}

// Or returns the bitwise or of t and obj. The equivalent Python is
// "t | obj".
//
// Return value: New Reference.
func (t *Type) Or(obj Object) (Object, error) {
	ret := C.PyNumber_Or(c(t), c(obj))
	return obj2ObjErr(ret)
}



/*
set fields:
  nb_or
  ob_base
  tp_as_number
  tp_basicsize
  tp_call
  tp_clear
  tp_dealloc
  tp_dictoffset
  tp_doc
  tp_flags
  tp_free
  tp_getattro
  tp_getset
  tp_init
  tp_is_gc
  tp_itemsize
  tp_members
  tp_methods
  tp_name
  tp_new
  tp_repr
  tp_setattro
  tp_traverse
  tp_vectorcall
  tp_vectorcall_offset
  tp_weaklistoffset
*/

