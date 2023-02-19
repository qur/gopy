// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

type Code struct {
	abstractObject
	o C.PyCodeObject
}

// CodeType is the Type object that represents the Code type.
var CodeType = (*Type)(unsafe.Pointer(&C.PyCode_Type))

func codeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.codeCheck(c(obj)) != 0
}

func newCode(obj *C.PyObject) *Code {
	return (*Code)(unsafe.Pointer(obj))
}







/*
ob_base = true
tp_basicsize = true
tp_dealloc = true
tp_doc = true
tp_getattro = true
tp_getset = true
tp_hash = true
tp_itemsize = true
tp_members = true
tp_methods = true
tp_name = true
tp_new = true
tp_repr = true
tp_richcompare = true
tp_weaklistoffset = true
*/

