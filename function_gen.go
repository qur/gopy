// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Function represents objects of the FunctionType (or PyFunctionType
// in the Python API) type.
type Function struct {
	abstractObject
	o C.PyFunctionObject
}

// FunctionType is the Type object that represents the Function type.
var FunctionType = (*Type)(unsafe.Pointer(&C.PyFunction_Type))

func functionCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.functionCheck(c(obj)) != 0
}

func newFunction(obj *C.PyObject) *Function {
	return (*Function)(unsafe.Pointer(obj))
}

// Repr returns a String representation of "f". This is equivalent to the
// Python "repr(f)".
//
// Return value: New Reference.
func (f *Function) Repr() (Object, error) {
	ret := C.PyObject_Repr(c(f))
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

