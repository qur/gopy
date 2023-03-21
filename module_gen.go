// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Bool is the representation of the Python module type.
type Module struct {
	o C.PyObject
}

var _ Object = (*Module)(nil)

// ModuleType is the Type object that represents the Module type.
var ModuleType = newType(&C.PyModule_Type)

func moduleCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.moduleCheck(c(obj)) != 0
}

// AsModule casts the given obj to a Module (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Module, then nil is returned.
//
// Return value: Borrowed Reference.
func AsModule(obj Object) *Module {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.moduleCheck(o) == 0 {
		return nil
	}
	return (*Module)(unsafe.Pointer(o))
}

func newModule(obj *C.PyObject) *Module {
	return (*Module)(unsafe.Pointer(obj))
}



func (m *Module) c() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(m))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (m *Module) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(m))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (m *Module) Type() *Type {
	return newType(c(m).ob_type)
}

// Decref decrements m's reference count, m may not be nil.
func (m *Module) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(m))
	obj.ob_refcnt--
	if obj.ob_refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments m's reference count, m may not be nil.
func (m *Module) Incref() {
	obj := (*C.PyObject)(unsafe.Pointer(m))
	obj.ob_refcnt++
}



/*
set fields:
*/

