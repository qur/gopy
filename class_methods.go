package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

//////////////////////////////////////////////////////////////////////////////
//
// METHODS
//
//////////////////////////////////////////////////////////////////////////////

func getMethod(obj *C.PyObject) (any, error) {
	pyobj := unsafe.Pointer(C.PyTuple_GetItem(obj, 0))
	idx := int(C.PyLong_AsLong(C.PyTuple_GetItem(obj, 1)))

	o := getClassObject(pyobj)
	if o == nil {
		return nil, fmt.Errorf("unknown object")
	}

	return reflect.ValueOf(o).Method(idx).Interface(), nil
}

//export goClassCallMethod
func goClassCallMethod(obj, unused *C.PyObject) *C.PyObject {
	// Unpack context and self pointer from obj
	m, err := getMethod(obj)
	if err != nil {
		raise(err)
		return nil
	}

	f := m.(func() (Object, error))

	return ce(f())
}

//export goClassCallMethodSingle
func goClassCallMethodSingle(obj, arg *C.PyObject) *C.PyObject {
	// Unpack context and self pointer from obj
	m, err := getMethod(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get arg ready to use, by turning it into a pointer of the appropriate
	// type
	o := newObject(arg)

	f := m.(func(Object) (Object, error))

	return ce(f(o))
}

//export goClassCallMethodArgs
func goClassCallMethodArgs(obj, args *C.PyObject) *C.PyObject {
	// Unpack context and self pointer from obj
	m, err := getMethod(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args ready to use, by turning it into a pointer of the appropriate
	// type
	a := newTuple(args)

	f := m.(func(*Tuple) (Object, error))

	return ce(f(a))
}

//export goClassCallMethodKwds
func goClassCallMethodKwds(obj, args, kwds *C.PyObject) *C.PyObject {
	// Unpack context and self pointer from obj
	m, err := getMethod(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple(args)
	k := newDict(kwds)

	f := m.(func(*Tuple, *Dict) (Object, error))

	return ce(f(a, k))
}

//////////////////////////////////////////////////////////////////////////////
//
// STATIC
//
//////////////////////////////////////////////////////////////////////////////

func getStaticMethod(obj *C.PyObject) any {
	c := getClass((*C.PyTypeObject)(unsafe.Pointer(C.PyTuple_GetItem(obj, 0))))
	name := C.GoString(C.PyUnicode_AsUTF8(C.PyTuple_GetItem(obj, 1)))
	return c.Static[name]
}

//export goClassCallStaticMethod
func goClassCallStaticMethod(obj, unused *C.PyObject) *C.PyObject {
	// Unpack function from obj
	m := getStaticMethod(obj)

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func() (Object, error))

	return ce(f())
}

//export goClassCallStaticMethodSingle
func goClassCallStaticMethodSingle(obj, arg *C.PyObject) *C.PyObject {
	// Unpack function from obj
	m := getStaticMethod(obj)

	// Get arg ready to use, by turning it into a pointer of the appropriate
	// type
	o := newObject(arg)

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(Object) (Object, error))

	return ce(f(o))
}

//export goClassCallStaticMethodArgs
func goClassCallStaticMethodArgs(obj, args *C.PyObject) *C.PyObject {
	// Unpack function from obj
	m := getStaticMethod(obj)

	// Get args ready to use, by turning it into a pointer of the appropriate
	// type
	a := newTuple(args)

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(a *Tuple) (Object, error))

	return ce(f(a))
}

//export goClassCallStaticMethodKwds
func goClassCallStaticMethodKwds(obj, args, kwds *C.PyObject) *C.PyObject {
	// Unpack function from obj
	m := getStaticMethod(obj)

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple(args)
	k := newDict(kwds)

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(a *Tuple, k *Dict) (Object, error))

	return ce(f(a, k))
}

//////////////////////////////////////////////////////////////////////////////
//
// CLASS
//
//////////////////////////////////////////////////////////////////////////////

func getMethodAndClass(obj *C.PyObject) (any, *Class, error) {
	t := (*C.PyObject)(obj)
	pyobj := unsafe.Pointer(C.PyTuple_GetItem(t, 0))

	c := getClass((*C.PyTypeObject)(pyobj))
	if c == nil {
		return nil, nil, fmt.Errorf("unknown class")
	}

	name := C.GoString(C.PyUnicode_AsUTF8(C.PyTuple_GetItem(t, 1)))
	return c.Class[name], c, nil
}

//export goClassCallClassMethod
func goClassCallClassMethod(obj, unused *C.PyObject) *C.PyObject {
	// Unpack context and class pointer from obj
	m, o, err := getMethodAndClass(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(p *Class) (Object, error))

	return ce(f(o))
}

//export goClassCallClassMethodSingle
func goClassCallClassMethodSingle(obj, arg *C.PyObject) *C.PyObject {
	// Unpack context and class pointer from obj
	m, o, err := getMethodAndClass(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get arg ready to use, by turning it into a pointer of the appropriate
	// type
	a := newObject(arg)

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(*Class, Object) (Object, error))

	return ce(f(o, a))
}

//export goClassCallClassMethodArgs
func goClassCallClassMethodArgs(obj, args *C.PyObject) *C.PyObject {
	// Unpack context and class pointer from obj
	m, o, err := getMethodAndClass(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args ready to use, by turning it into a pointer of the appropriate
	// type
	a := newTuple(args)

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(o *Class, a *Tuple) (Object, error))

	return ce(f(o, a))
}

//export goClassCallClassMethodKwds
func goClassCallClassMethodKwds(obj, args, kwds *C.PyObject) *C.PyObject {
	// Unpack context and class pointer from obj
	m, o, err := getMethodAndClass(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple(args)
	k := newDict(kwds)

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(o *Class, a *Tuple, k *Dict) (Object, error))

	return ce(f(o, a, k))
}
