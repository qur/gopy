package py

// #include "utils.h"
import "C"
import (
	"reflect"
	"unsafe"
)

//export goClassTraverse
func goClassTraverse(obj *C.PyObject, visit C.visitproc, arg unsafe.Pointer) C.int {
	co := getClassObject(obj)
	if co == nil {
		raise(TypeError.Err("not a recognised type: %s", newObject((*C.PyObject)(obj)).Type()))
		return -1
	}

	return classTraverse(co, visit, arg)
}

func classTraverse(co ClassObject, visit C.visitproc, arg unsafe.Pointer) C.int {
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

	class := getClass(co.Type().c())
	switch b := class.BaseType.(type) {
	case *Class:
		co = getClassObjectByType(c(co), b.Type().c())
		if co != nil {
			return classTraverse(co, visit, arg)
		}
	case *Type:
		base := co.Type().o.tp_base
		if base != nil {
			return C.typeTraverse(base, c(co), visit, arg)
		}
	}

	return 0
}

//export goClassClear
func goClassClear(obj *C.PyObject) C.int {
	co := getClassObject(obj)
	if co == nil {
		raise(TypeError.Err("not a recognised type: %s", newObject((*C.PyObject)(obj)).Type()))
		return -1
	}

	return classClear(co)
}

func classClear(co ClassObject) C.int {
	ClearClassObject(co)

	class := getClass(co.Type().c())
	switch b := class.BaseType.(type) {
	case *Class:
		co = getClassObjectByType(c(co), b.Type().c())
		if co != nil {
			return classClear(co)
		}
	case *Type:
		base := co.Type().o.tp_base
		if base != nil {
			return C.typeClear(base, c(co))
		}
	}

	return 0
}

type tpDealloc interface {
	PyDealloc()
}

//export goClassDealloc
func goClassDealloc(obj *C.PyObject) {
	// Turn obj into the ClassObject instead of the proxy.
	co := getClassObject(obj)
	if co == nil {
		// not a recognised type, but dealloc can't return an error ...
		return
	}

	if classDealloc(co) {
		// classDealloc has already freed the Python memory, but we need to
		// unregister the Go instances
		clearClassObject(obj)
		return
	}

	// we always want Python to _actually_ free the object, any registered hook
	// should just be tidying things up on the Go side.
	free(co)
}

func classDealloc(co ClassObject) bool {
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

	switch b := class.BaseType.(type) {
	case *Class:
		co = getClassObjectByType(c(co), b.Type().c())
		if co != nil {
			return classDealloc(co)
		}
	case *Type:
		base := co.Type().o.tp_base
		if base != nil {
			C.typeDealloc(base, c(co))
			return true
		}
	}

	return false
}

//export goClassNew
func goClassNew(typ *C.PyTypeObject, args, kwds *C.PyObject) *C.PyObject {
	// Get the Python type object
	pyType := typ

	class := getClass(pyType)

	// might be parent type that is a class
	for class == nil && pyType.tp_base != nil {
		pyType = (*C.PyTypeObject)(unsafe.Pointer(pyType.tp_base))
		class = getClass(pyType)
	}

	if class == nil {
		t := newType(pyType)
		raise(TypeError.Err("Not a recognised type: %s", t))
		return nil
	}

	return class.new(typ, args, kwds)
}

func (cls *Class) new(typ *C.PyTypeObject, args, kwds *C.PyObject) *C.PyObject {
	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple(args)
	k := newDict(kwds)

	// allocate the go object
	goObj, err := cls.newObject(a, k)
	if err != nil {
		raise(err)
		return nil
	}

	// allocate the Python proxy object
	var pyObj *C.PyObject
	switch b := cls.BaseType.(type) {
	case nil:
		pyObj = C.typeAlloc(typ, 0)
	case *Class:
		pyObj = b.new(typ, args, kwds)
	default:
		pyObj = C.typeNew(cls.base.o.tp_base, typ, args, kwds)
	}
	if pyObj == nil {
		return nil
	}

	// finalise the setup of the go object
	goObj.setBase(newBaseObject(pyObj))
	v := reflect.ValueOf(goObj).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Type() {
		case cipType, cnpType, cspType, cmpType:
			cp := field.Addr().Interface().(classProtocol)
			cp.setCBO(goObj.getCBO())
		}
	}

	registerClassObject(pyObj, cls.base.c(), goObj)
	if typ != cls.base.c() {
		registerClassObject(pyObj, typ, goObj)
	}

	return pyObj
}
