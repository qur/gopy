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
	idx := int(C.PyLong_AsLong(C.PyTuple_GetItem(t, 1)))

	// Turn obj into the ClassObject instead of the proxy
	co := newObject((*C.PyObject)(obj)).(ClassObject)
	m := reflect.ValueOf(co).Method(idx).Interface()

	// Turn arg into something usable
	a := newObject((*C.PyObject)(arg))

	// Turn the function into something we can call
	f := m.(func(Object) error)

	if err := f(a); err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export goClassGetProp
func goClassGetProp(obj, closure unsafe.Pointer) unsafe.Pointer {
	// Unpack get function from closure
	t := (*C.PyObject)(closure)
	idx := int(C.PyLong_AsLong(C.PyTuple_GetItem(t, 0)))

	// Turn obj into the ClassObject instead of the proxy
	co := newObject((*C.PyObject)(obj)).(ClassObject)
	m := reflect.ValueOf(co).Method(idx).Interface()

	// Turn the function into something we can call
	f := m.(func() (Object, error))

	ret, err := f()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
