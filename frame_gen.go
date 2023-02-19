// Code generated by gen_types. DO NOT EDIT.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

type Frame struct {
	abstractObject
	o C.PyFrameObject
}

// FrameType is the Type object that represents the Frame type.
var FrameType = (*Type)(unsafe.Pointer(&C.PyFrame_Type))

func frameCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.frameCheck(c(obj)) != 0
}

func newFrame(obj *C.PyObject) *Frame {
	return (*Frame)(unsafe.Pointer(obj))
}







/*
ob_base = true
tp_basicsize = true
tp_clear = true
tp_dealloc = true
tp_flags = true
tp_getattro = true
tp_getset = true
tp_itemsize = true
tp_members = true
tp_methods = true
tp_name = true
tp_repr = true
tp_setattro = true
tp_traverse = true
*/

