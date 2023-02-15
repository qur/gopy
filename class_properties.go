package py

// #include "utils.h"
import "C"

import (
	"reflect"
	"unsafe"
)

//export goClassSetProp
func goClassSetProp(obj, arg, closure unsafe.Pointer) int {
	// Unpack set function from closure
	t := (*C.PyObject)(closure)
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 1), nil)

	// Turn obj into the ClassObject instead of the proxy
	co := newObject((*C.PyObject)(obj)).(ClassObject)
	v := reflect.ValueOf(co)

	// Turn arg into something usable
	a := newObject((*C.PyObject)(arg))

	// Turn the function into something we can call
	f := (*func(p unsafe.Pointer, a Object) error)(unsafe.Pointer(&m))

	err := (*f)(v.UnsafePointer(), a)
	if err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export goClassGetProp
func goClassGetProp(obj, closure unsafe.Pointer) unsafe.Pointer {
	// Unpack get function from closure
	t := (*C.PyObject)(closure)
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 0), nil)

	// Turn obj into the ClassObject instead of the proxy
	co := newObject((*C.PyObject)(obj)).(ClassObject)
	v := reflect.ValueOf(co)

	// Turn the function into something we can call
	f := (*func(p unsafe.Pointer) (Object, error))(unsafe.Pointer(&m))

	ret, err := (*f)(v.UnsafePointer())
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
