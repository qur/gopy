package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

func getField(obj, arg unsafe.Pointer) (reflect.Value, error) {
	o := getClassObject(obj)
	if o == nil {
		return reflect.Value{}, fmt.Errorf("unknown object")
	}

	idx := int(C.PyLong_AsLong((*C.PyObject)(arg)))

	return reflect.ValueOf(o).Elem().Field(idx), nil
}

//export goClassNatGet
func goClassNatGet(obj, idx unsafe.Pointer) *C.PyObject {
	f, err := getField(obj, idx)
	if err != nil {
		raise(err)
		return nil
	}

	switch f.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return c(NewLong(f.Int()))
	}

	raise(NotImplementedError.ErrV(None))
	return nil
}

//export goClassNatSet
func goClassNatSet(obj, idx, obj2 unsafe.Pointer) int {
	f, err := getField(obj, idx)
	if err != nil {
		raise(err)
		return -1
	}

	// This is the new value we are being asked to set
	value := (*C.PyObject)(obj2)

	switch f.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := int64(C.PyLong_AsLong(value))
		if exceptionRaised() {
			return -1
		}
		f.SetInt(v)
		return 0
	}

	raise(NotImplementedError.ErrV(None))
	return -1
}

//export goClassObjGet
func goClassObjGet(obj, idx unsafe.Pointer) *C.PyObject {
	f, err := getField(obj, idx)
	if err != nil {
		raise(err)
		return nil
	}

	if f.IsNil() {
		None.Incref()
		return c(None)
	}

	o := f.Interface().(Object)
	o.Incref()
	return c(o)
}

//export goClassObjSet
func goClassObjSet(obj, idx, obj2 unsafe.Pointer) int {
	f, err := getField(obj, idx)
	if err != nil {
		raise(err)
		return -1
	}

	value := newObject((*C.PyObject)(obj2))
	v := reflect.ValueOf(value)

	// If the given value wasn't assignable to the field - raise an error
	if !v.Type().AssignableTo(f.Type()) {
		tn := f.Type().Elem().Name()
		raise(TypeError.Err("Cannot assign '%T' to '*%v'", value, tn))
		return -1
	}

	// If f is nil, then assignment is simple.
	if f.IsNil() {
		Incref(value)
		f.Set(v)
		return 0
	}

	// If f is not nil, then we have to be careful with refcounts, as decref
	// could invoke destructor code etc.
	tmp := f.Interface().(Object)
	Incref(value)
	f.Set(v)
	Decref(tmp)
	return 0
}
