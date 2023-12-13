// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Gen represents objects of the GenType (or PyGen_Type
// in the Python API) type.
//
// This type implements the Iterator protocol.
type Gen struct {
	o C.PyGenObject
}

var _ Object = (*Gen)(nil)

// GenType is the Type object that represents the Gen type.
var GenType = newType(&C.PyGen_Type)

func genCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.genCheck(c(obj)) != 0
}

// AsGen casts the given obj to a Gen (i.e. the underlying
// Python Object is the same, just the type is changed). If the value cannot be
// cast to a Gen, then nil is returned.
//
// Return value: Borrowed Reference.
func AsGen(obj Object) *Gen {
	if obj == nil {
		return nil
	}
	o := c(obj)
	if C.genCheck(o) == 0 {
		return nil
	}
	return (*Gen)(unsafe.Pointer(o))
}

func newGen(obj *C.PyObject) *Gen {
	return (*Gen)(unsafe.Pointer(obj))
}



func (g *Gen) c() *C.PyGenObject {
	return (*C.PyGenObject)(unsafe.Pointer(g))
}

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (g *Gen) Base() *BaseObject {
	return (*BaseObject)(unsafe.Pointer(g))
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (g *Gen) Type() *Type {
	return newType(c(g).ob_type)
}


// Decref decrements g's reference count, g may not be nil.
func (g *Gen) Decref() {
	obj := (*C.PyObject)(unsafe.Pointer(g))
	refcnt := (*int)(unsafe.Pointer(&obj.anon0[0]))
	*refcnt--
	if *refcnt == 0 {
		C._Py_Dealloc(obj)
	}
}

// Incref increments g's reference count, g may not be nil.
func (g *Gen) Incref() {
	refcnt := (*int)(unsafe.Pointer(&(*C.PyObject)(unsafe.Pointer(g)).anon0[0]))
	*refcnt++
}

// Repr returns a String representation of "g". This is equivalent to the
// Python "repr(g)".
//
// Return value: New Reference.
func (g *Gen) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(g))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// HasAttr returns true if "g" has the attribute "name".  This is equivalent
// to the Python "hasattr(g, name)".
func (g *Gen) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(g), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "g" with the name "name".  This is
// equivalent to the Python "g.name".
//
// Return value: New Reference.
func (g *Gen) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(g), c(name))
	return obj2ObjErr(ret)
}

// Iter returns an Iterator that will iterate over the members of g.
func (g *Gen) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c(g))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

// AsIteratorMethods returns a IteratorMethods instance that refers to the same
// underlying Python object as g.
//
// This method also means that Gen implements the Iterator interface.
func (g *Gen) AsIteratorMethods() *IteratorMethods {
	return (*IteratorMethods)(unsafe.Pointer(g.Base()))
}

// Next is the Iterator method, it returns the next item from the Object being
// iterated. When the end is reached then both the Object and the error will be
// nil.
func (g *Gen) Next() (Object, error) {
	ret := C.PyIter_Next(c(g))
	return obj2ObjErr(ret)
}



/*
set fields:
  am_send
  ob_base
  tp_as_async
  tp_basicsize
  tp_dealloc
  tp_finalize
  tp_flags
  tp_getattro
  tp_getset
  tp_itemsize
  tp_iter
  tp_iternext
  tp_members
  tp_methods
  tp_name
  tp_repr
  tp_traverse
  tp_weaklistoffset
*/

