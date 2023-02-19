// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

type Gen struct {
	abstractObject
	o C.PyGenObject
}

// GenType is the Type object that represents the Gen type.
var GenType = (*Type)(unsafe.Pointer(&C.PyGen_Type))

func genCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.genCheck(c(obj)) != 0
}

func newGen(obj *C.PyObject) *Gen {
	return (*Gen)(unsafe.Pointer(obj))
}







/*
am_send = true
ob_base = true
tp_as_async = true
tp_basicsize = true
tp_dealloc = true
tp_finalize = true
tp_flags = true
tp_getattro = true
tp_getset = true
tp_itemsize = true
tp_iter = true
tp_iternext = true
tp_members = true
tp_methods = true
tp_name = true
tp_repr = true
tp_traverse = true
tp_weaklistoffset = true
*/

