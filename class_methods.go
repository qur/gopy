package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

func getMethod(obj unsafe.Pointer) (any, error) {
	t := (*C.PyObject)(obj)
	pyobj := unsafe.Pointer(C.PyTuple_GetItem(t, 0))
	idx := int(C.PyLong_AsLong(C.PyTuple_GetItem(t, 1)))

	o := getClassObject(pyobj)
	if o == nil {
		return nil, fmt.Errorf("unknown object")
	}

	return reflect.ValueOf(o).Method(idx).Interface(), nil
}

//export goClassCallMethod
func goClassCallMethod(obj, unused unsafe.Pointer) unsafe.Pointer {
	// Unpack context and self pointer from obj
	m, err := getMethod(obj)
	if err != nil {
		raise(err)
		return nil
	}

	f := m.(func() (Object, error))

	ret, err := f()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallMethodArgs
func goClassCallMethodArgs(obj, args unsafe.Pointer) unsafe.Pointer {
	// Unpack context and self pointer from obj
	m, err := getMethod(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args ready to use, by turning it into a pointer of the appropriate
	// type
	a := newTuple((*C.PyObject)(args))

	f := m.(func(*Tuple) (Object, error))

	ret, err := f(a)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallMethodKwds
func goClassCallMethodKwds(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Unpack context and self pointer from obj
	m, err := getMethod(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	f := m.(func(*Tuple, *Dict) (Object, error))

	ret, err := f(a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

func getStaticMethod(obj unsafe.Pointer) any {
	t := (*C.PyObject)(obj)
	c, _ := getType((*C.PyTypeObject)(unsafe.Pointer(C.PyTuple_GetItem(t, 0))))
	name := C.GoString(C.PyUnicode_AsUTF8(C.PyTuple_GetItem(t, 1)))
	return c.Static[name]
}

//export goClassCallStaticMethod
func goClassCallStaticMethod(obj, unused unsafe.Pointer) unsafe.Pointer {
	// Unpack function from obj
	m := getStaticMethod(obj)

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func() (Object, error))

	ret, err := f()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallStaticMethodArgs
func goClassCallStaticMethodArgs(obj, args unsafe.Pointer) unsafe.Pointer {
	// Unpack function from obj
	m := getStaticMethod(obj)

	// Get args ready to use, by turning it into a pointer of the appropriate
	// type
	a := newTuple((*C.PyObject)(args))

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(a *Tuple) (Object, error))

	ret, err := f(a)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallStaticMethodKwds
func goClassCallStaticMethodKwds(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Unpack function from obj
	m := getStaticMethod(obj)

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(a *Tuple, k *Dict) (Object, error))

	ret, err := f(a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

func getMethodAndClass(obj unsafe.Pointer) (any, *Class, error) {
	t := (*C.PyObject)(obj)
	pyobj := unsafe.Pointer(C.PyTuple_GetItem(t, 0))

	c, _ := getType((*C.PyTypeObject)(pyobj))
	if c == nil {
		return nil, nil, fmt.Errorf("unknown class")
	}

	name := C.GoString(C.PyUnicode_AsUTF8(C.PyTuple_GetItem(t, 1)))
	return c.Class[name], c, nil
}

//export goClassCallClassMethod
func goClassCallClassMethod(obj, unused unsafe.Pointer) unsafe.Pointer {
	// Unpack context and class pointer from obj
	m, o, err := getMethodAndClass(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(p *Class) (Object, error))

	ret, err := f(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallClassMethodArgs
func goClassCallClassMethodArgs(obj, args unsafe.Pointer) unsafe.Pointer {
	// Unpack context and class pointer from obj
	m, o, err := getMethodAndClass(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args ready to use, by turning it into a pointer of the appropriate
	// type
	a := newTuple((*C.PyObject)(args))

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(o *Class, a *Tuple) (Object, error))

	ret, err := f(o, a)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallClassMethodKwds
func goClassCallClassMethodKwds(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Unpack context and class pointer from obj
	m, o, err := getMethodAndClass(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := m.(func(o *Class, a *Tuple, k *Dict) (Object, error))

	ret, err := f(o, a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
