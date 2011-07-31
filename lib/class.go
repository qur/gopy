// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include <Python.h>
// #include "utils.h"
// static inline PyTypeObject *newType(void) {
//     return calloc(1, sizeof(PyTypeObject));
// }
// static inline int typeReady(PyTypeObject *o) {
//     int ret;
//     o->tp_basicsize += sizeof(PyObject);
//     if (o->tp_new == NULL) {
//         o->tp_new = PyType_GenericNew;
//     }
//     ret = PyType_Ready(o);
//     if (ret == 0) {
//         o->tp_methods = calloc(1, sizeof(ClassContext));
//     }
//     return ret;
// }
// static inline int setTypeAttr(PyTypeObject *tp, char *name, PyObject *o) {
//     return PyDict_SetItemString(tp->tp_dict, name, o);
// }
// static inline uintptr_t headSize(void) { return sizeof(PyObject); }
// static inline PyObject *typeAlloc(PyTypeObject *type, Py_ssize_t n) { return type->tp_alloc(type, n); }
import "C"

import (
	"os"
	"reflect"
	"strings"
	"unsafe"
)

// A Class struct instance is used to define a Python class that has been
// implemented in Go.
//
// Name should be the name of the type in Python, including the package name,
// e.g. "test.MyClass"
//
// Flags and Doc are currently unused.
//
// Pointer should be set to a pointer of the struct type that will represent an
// instance of the Python class.  This struct must contain an embedded
// py.BaseObject as its first member.  The easiest ways to set Pointer are
// either to use a struct literal (i.e. &MyClass{}), or to cast nil (i.e.
// (*MyClass)(nil)), if the struct is large then the latter method is more
// efficient (as an instance of the struct is not created).
//
// This struct may have the following special methods (the equivalent Python
// methods are also indicated):
//
//   PyInit(args *py.Tuple, kwds *py.Dict) os.Error              // __init__
//   PyCall(args *py.Tuple, kwds *py.Dict) (py.Object, os.Error) // __call__
//   PyRepr() string                                             // __repr__
//   PyStr() string                                              // __str__
//   PyCompare(obj py.Object) (int, os.Error)                    // __cmp__
//
// Properties are also supported, by implementing get and set methods:
//
//   PyGet_XXX() (py.Object, os.Error)
//   PySet_XXX(value py.Object) os.Error
//
// Methods on the Python class are implemented by methods with the Py_ prefix:
//
//   Py_XXX(args *py.Tuple, kwds *py.Dict) (py.Object, os.Error)
//
type Class struct {
	Name    string
	Flags   int
	Doc     string
	Pointer interface{}
}

//export callClassMethod
func callClassMethod(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Unpack context and self pointer from obj
	t := (*C.PyObject)(obj)
	pyobj := unsafe.Pointer(C.PyTuple_GetItem(t, 0))
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 1), nil)

	// Turn the Python self into a Go struct self
	self := unsafe.Pointer(uintptr(pyobj) + uintptr(C.headSize()))

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := (*func(p unsafe.Pointer, a *Tuple, k *Dict) (Object, os.Error))(unsafe.Pointer(&m))

	ret, err := (*f)(self, a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export setClassProperty
func setClassProperty(obj, arg, closure unsafe.Pointer) int {
	// Turn the Python self into a Go struct self
	self := unsafe.Pointer(uintptr(obj) + uintptr(C.headSize()))

	// Unpack set function from closure
	t := (*C.PyObject)(closure)
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 1), nil)

	// Turn arg into something usable
	a := newBaseObject((*C.PyObject)(arg)).actual()

	// Turn the function into something we can call
	f := (*func(p unsafe.Pointer, a Object) os.Error)(unsafe.Pointer(&m))

	err := (*f)(self, a)
	if err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export getClassProperty
func getClassProperty(obj, closure unsafe.Pointer) unsafe.Pointer {
	// Turn the Python self into a Go struct self
	self := unsafe.Pointer(uintptr(obj) + uintptr(C.headSize()))

	// Unpack set function from closure
	t := (*C.PyObject)(closure)
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 0), nil)

	// Turn the function into something we can call
	f := (*func(p unsafe.Pointer) (Object, os.Error))(unsafe.Pointer(&m))

	ret, err := (*f)(self)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

func getClassContext(obj unsafe.Pointer) *C.ClassContext {
	o := (*C.PyObject)(obj)
	return (*C.ClassContext)(unsafe.Pointer(o.ob_type.tp_methods))
}

//export newGoClass
func newGoClass(typ, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Get the Python type object
	pyType := (*C.PyTypeObject)(typ)

	// Get the class context
	ctxt := (*C.ClassContext)(unsafe.Pointer(pyType.tp_methods))

	var obj *C.PyObject
	if ctxt.new != nil {
		// Get args and kwds ready to use, by turning them into pointers of the
		// appropriate type
		a := newTuple((*C.PyObject)(args))
		k := newDict((*C.PyObject)(kwds))

		// Turn the function into something we can call
		f := (*func(*Class, *Tuple, *Dict) (Object, os.Error))(unsafe.Pointer(&ctxt.new))

		// TODO: need to get hold of a pointer the the Class that defines this
		// type - and add an Alloc method to it.
		o, err := (*f)(nil, a, k)
		if err != nil {
			raise(err)
			return nil
		}
		obj = c(o)
	} else {
		// Create a new Python instance
		obj = C.typeAlloc(pyType, 0)
	}

	// Turn the Python self into a Go struct self
	self := (*BaseObject)(unsafe.Pointer(uintptr(unsafe.Pointer(obj)) + uintptr(C.headSize())))

	// We need to setup the internal pointer, so that the instance can act like
	// an Object
	self.o = obj

	return unsafe.Pointer(obj)
}

//export initGoClass
func initGoClass(obj, args, kwds unsafe.Pointer) int {
	// Turn the Python self into a Go struct self
	self := unsafe.Pointer(uintptr(obj) + uintptr(C.headSize()))

	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, *Tuple, *Dict) os.Error)(unsafe.Pointer(&ctxt.init))

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	err := (*f)(self, a, k)
	if err != nil {
		// Turn err into exception
		return -1
	}

	return 0
}

//export reprGoClass
func reprGoClass(obj unsafe.Pointer) unsafe.Pointer {
	// Turn the Python self into a Go struct self
	self := unsafe.Pointer(uintptr(obj) + uintptr(C.headSize()))

	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) string)(unsafe.Pointer(&ctxt.repr))

	s := C.CString((*f)(self))
	defer C.free(unsafe.Pointer(s))

	return unsafe.Pointer(C.PyString_FromString(s))
}

//export strGoClass
func strGoClass(obj unsafe.Pointer) unsafe.Pointer {
	// Turn the Python self into a Go struct self
	self := unsafe.Pointer(uintptr(obj) + uintptr(C.headSize()))

	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer) string)(unsafe.Pointer(&ctxt.str))

	s := C.CString((*f)(self))
	defer C.free(unsafe.Pointer(s))

	return unsafe.Pointer(C.PyString_FromString(s))
}

//export callGoClass
func callGoClass(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Turn the Python self into a Go struct self
	self := unsafe.Pointer(uintptr(obj) + uintptr(C.headSize()))

	// Get the class context
	ctxt := getClassContext(obj)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, *Tuple, *Dict) (Object, os.Error))(unsafe.Pointer(&ctxt.call))

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	ret, err := (*f)(self, a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export compareGoClass
func compareGoClass(obj1, obj2 unsafe.Pointer) int {
	// Turn the Python self into a Go struct self
	self := unsafe.Pointer(uintptr(obj1) + uintptr(C.headSize()))

	// Get the class context
	ctxt := getClassContext(obj1)

	// Turn the function into something we can call
	f := (*func(unsafe.Pointer, Object) (int, os.Error))(unsafe.Pointer(&ctxt.compare))

	o := newBaseObject((*C.PyObject)(obj2)).actual()

	ret, err := (*f)(self, o)
	if err != nil {
		raise(err)
		return -1
	}

	return ret
}

type prop struct {
	get unsafe.Pointer
	set unsafe.Pointer
}

// Create creates and returns a pointer to a PyTypeObject that is the Python
// representation of the class that has been implemented in Go.
func (c Class) Create() (Object, os.Error) {
	typ := reflect.TypeOf(c.Pointer)

	pyType := C.newType()
	pyType.tp_name = C.CString(c.Name)
	pyType.tp_basicsize = C.Py_ssize_t(typ.Elem().Size())
	pyType.tp_flags = C.Py_TPFLAGS_DEFAULT

	if C.typeReady(pyType) < 0 {
		C.free(unsafe.Pointer(pyType.tp_name))
		C.free(unsafe.Pointer(pyType))
		return nil, exception()
	}

	registerType(pyType, unsafe.Typeof(c.Pointer))

	ctxt := (*C.ClassContext)(unsafe.Pointer(pyType.tp_methods))

	// We don't use tp_methods, and it is read when calling PyType_Ready - so we
	// use it to hide a classContext struct.  The classContext starts with a
	// NULL pointer just in case, so it looks like an empty methods list if
	// Python does try to process it.
	pyType.tp_methods = (*C.struct_PyMethodDef)(unsafe.Pointer(ctxt))

	props := make(map[string]prop)

	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		f := unsafe.Pointer(m.Func.Pointer())
		if !strings.HasPrefix(m.Name, "Py") {
			continue
		}
		parts := strings.SplitN(m.Name, "_", 2)
		switch parts[0] {
		case "PyInit":
			ctxt.init = f
		case "PyRepr":
			ctxt.repr = f
		case "PyStr":
			ctxt.str = f
		case "PyCall":
			ctxt.call = f
		case "PyCompare":
			ctxt.compare = f
		case "Py":
			s := C.CString(parts[1])
			C.setTypeAttr(pyType, s, C.newMethod(s, f))
		case "PySet":
			p := props[parts[1]]
			p.set = f
			props[parts[1]] = p
		case "PyGet":
			p := props[parts[1]]
			p.get = f
			props[parts[1]] = p
		}
	}

	C.setClassContext(pyType, ctxt)

	for name, prop := range props {
		s := C.CString(name)
		C.setTypeAttr(pyType, s, C.newProperty(pyType, s, prop.get, prop.set))
	}

	return newBaseObject((*C.PyObject)(unsafe.Pointer(pyType))), nil
}
