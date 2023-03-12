package py

// #include "utils.h"
import "C"
import (
	"log"
	"reflect"
	"unsafe"
)

//export goClassTraverse
func goClassTraverse(obj *C.PyObject, visit C.visitproc, arg unsafe.Pointer) C.int {
	co := getClassObject(unsafe.Pointer(obj))
	if co == nil {
		raise(TypeError.Err("not a recognised type: %s", newObject((*C.PyObject)(obj)).Type()))
		return -1
	}

	v := reflect.ValueOf(co).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.Type().Implements(otyp) || f.Type() == cboType || f.IsNil() {
			// only care about non-nil Object values that aren't ClassBaseObject
			continue
		}
		if ret := C.doVisit((*C.PyObject)(f.UnsafePointer()), visit, arg); ret != 0 {
			return ret
		}
	}

	base := co.Type().o.tp_base
	if base != nil {
		return C.typeTraverse(base, obj, visit, arg)
	}

	return 0
}

//export goClassClear
func goClassClear(obj *C.PyObject) C.int {
	co := getClassObject(unsafe.Pointer(obj))
	if co == nil {
		raise(TypeError.Err("not a recognised type: %s", newObject((*C.PyObject)(obj)).Type()))
		return -1
	}

	ClearClassObject(co)

	base := co.Type().o.tp_base
	if base != nil {
		return C.typeClear(base, obj)
	}

	return 0
}

type tpDealloc interface {
	PyDealloc()
}

//export goClassDealloc
func goClassDealloc(obj unsafe.Pointer) {
	// Turn obj into the ClassObject instead of the proxy.
	co := getClassObject(obj)
	if co == nil {
		// not a recognised type, but dealloc can't return an error ...
		return
	}

	class := getClass(co.Type().c())
	if class != nil && (class.Flags&ClassHaveGC != 0) {
		C.PyObject_GC_UnTrack(unsafe.Pointer(c(co)))
	}

	// If co implements tpDealloc then call the PyDealloc method, else use
	// ClearClassObject
	if dealloc, ok := co.(tpDealloc); ok {
		dealloc.PyDealloc()
	} else {
		ClearClassObject(co)
	}

	// we always want Python to _actually_ free the object, any registered hook
	// should just be tidying things up on the Go side.
	free(co)
}

//export goClassNew
func goClassNew(typ, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Get the Python type object
	pyType := (*C.PyTypeObject)(typ)
	base := pyType.tp_base

	class := getClass(pyType)

	// must be parent type that is a class
	for class == nil && pyType.tp_base != nil {
		pyType = (*C.PyTypeObject)(unsafe.Pointer(pyType.tp_base))
		class = getClass(pyType)
	}

	if class == nil {
		t := newType(pyType)
		raise(TypeError.Err("Not a recognised type: %s", t))
		return nil
	}

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	// allocate the go object
	goObj, err := class.newObject(a, k)
	if err != nil {
		raise(err)
		return nil
	}

	// allocate the Python proxy object
	pyObj := unsafe.Pointer(C.typeNew(base, (*C.PyTypeObject)(typ), (*C.PyObject)(args), (*C.PyObject)(kwds)))
	if pyObj == nil {
		return nil
	}

	// finalise the setup of the go object
	goObj.setBase((*BaseObject)(pyObj))
	v := reflect.ValueOf(goObj).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Type() {
		case cipType, cnpType, cspType, cmpType:
			cp := field.Addr().Interface().(classProtocol)
			cp.setCBO(goObj.getCBO())
		}
	}
	registerClassObject(pyObj, goObj)

	return pyObj
}
