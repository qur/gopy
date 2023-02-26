package py

// #include "utils.h"
import "C"

import (
	"reflect"
)

//export goClassGetProp
func goClassGetProp(obj, closure *C.PyObject) *C.PyObject {
	// Unpack get function from closure
	idx := int(C.PyLong_AsLong(C.PyTuple_GetItem(closure, 0)))

	// Turn obj into the ClassObject instead of the proxy
	co := newObject(obj).(ClassObject)
	m := reflect.ValueOf(co).Method(idx).Interface()

	// Turn the function into something we can call
	f := m.(func() (Object, error))

	return ce(f())
}

//export goClassSetProp
func goClassSetProp(obj, arg, closure *C.PyObject) C.int {
	// Unpack set function from closure
	idx := int(C.PyLong_AsLong(C.PyTuple_GetItem(closure, 1)))

	// Turn obj into the ClassObject instead of the proxy
	co := newObject(obj).(ClassObject)
	m := reflect.ValueOf(co).Method(idx).Interface()

	// Turn arg into something usable
	a := newObject(arg)

	// Turn the function into something we can call
	f := m.(func(Object) error)

	return err2Int(f(a))
}
