package py

// #include "utils.h"
import "C"

import (
	"errors"
	"reflect"
	"unsafe"
)

func getField(obj, arg *C.PyObject) (reflect.Value, reflect.StructField, error) {
	pyType := unsafe.Pointer(C.PyTuple_GetItem(arg, 0))

	o := getClassObjectByType(obj, (*C.PyTypeObject)(pyType))
	if o == nil {
		return reflect.Value{}, reflect.StructField{}, errors.New("unknown object")
	}

	idx := int(C.PyLong_AsLong(C.PyTuple_GetItem(arg, 1)))

	return reflect.ValueOf(o).Elem().Field(idx), reflect.TypeOf(o).Elem().Field(idx), nil
}

//export goClassNatGet
func goClassNatGet(obj, idx *C.PyObject) *C.PyObject {
	f, _, err := getField(obj, idx)
	if err != nil {
		raise(err)
		return nil
	}

	switch f.Type().Kind() {
	case reflect.Bool:
		return c(NewBool(f.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return c(NewLong(f.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return c(NewLong(int64(f.Uint())))
	case reflect.Float32, reflect.Float64:
		return ce(NewFloat(f.Float()))
	case reflect.String:
		return ce(NewUnicode(f.String()))
	case reflect.Complex64, reflect.Complex128:
		return ce(NewComplex(f.Complex()))
	}

	raise(NotImplementedError.ErrV(None))

	return nil
}

//export goClassNatSet
func goClassNatSet(obj, obj2, idx *C.PyObject) (ret int) {
	// obj is the destination object
	f, t, err := getField(obj, idx)
	if err != nil {
		raise(err)
		return -1
	}

	// it's very unlikely, but calling the `.SetXXX()` methods of
	// [reflect.Value] can panic, so convert any such panic into a Python
	// exception instead.
	defer func() {
		if v := recover(); v != nil {
			raise(AttributeError.Err("failed to set %s: %v", t.Name, v))

			ret = -1
		}
	}()

	// obj2 is the new value we are being asked to set into obj

	switch f.Type().Kind() {
	case reflect.Bool:
		return setBoolField(f, obj2, t.Name)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setIntField(f, obj2)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUintField(f, obj2)
	case reflect.Float32, reflect.Float64:
		return setFloatField(f, obj2)
	case reflect.String:
		return setStringField(f, obj2)
	case reflect.Complex64, reflect.Complex128:
		return setComplexField(f, obj2, t.Name)
	}

	raise(NotImplementedError.ErrV(None))

	return -1
}

func setBoolField(f reflect.Value, obj *C.PyObject, name string) int {
	o := newObject(obj)

	b, ok := o.(*Bool)
	if !ok {
		raise(TypeError.Err("field %s is bool, got %s", name, o.Type()))
		return -1
	}

	f.SetBool(b.Bool())

	return 0
}

func setIntField(f reflect.Value, obj *C.PyObject) int {
	v := int64(C.PyLong_AsLongLong(obj))

	if exceptionRaised() {
		return -1
	}

	f.SetInt(v)

	return 0
}

func setUintField(f reflect.Value, obj *C.PyObject) int {
	v := uint64(C.PyLong_AsUnsignedLongLong(obj))

	if exceptionRaised() {
		return -1
	}

	f.SetUint(v)

	return 0
}

func setFloatField(f reflect.Value, obj *C.PyObject) int {
	v := float64(C.PyFloat_AsDouble(obj))

	if exceptionRaised() {
		return -1
	}

	f.SetFloat(v)

	return 0
}

func setStringField(f reflect.Value, obj *C.PyObject) int {
	v := C.PyUnicode_AsUTF8(obj)

	if exceptionRaised() {
		return -1
	}

	f.SetString(C.GoString(v))

	return 0
}

func setComplexField(f reflect.Value, obj *C.PyObject, name string) int {
	o := newObject(obj)

	v, ok := o.(*Complex)
	if !ok {
		raise(TypeError.Err("field %s is complex, got %s", name, o.Type()))
		return -1
	}

	f.SetComplex(v.Complex128())

	return 0
}

//export goClassObjGet
func goClassObjGet(obj, idx *C.PyObject) *C.PyObject {
	f, _, err := getField(obj, idx)
	if err != nil {
		raise(err)
		return nil
	}

	if f.IsNil() {
		return c(None)
	}

	o := f.Interface().(Object)
	o.Incref()

	return c(o)
}

//export goClassObjSet
func goClassObjSet(obj, obj2, idx *C.PyObject) int {
	f, _, err := getField(obj, idx)
	if err != nil {
		raise(err)
		return -1
	}

	value := newObject(obj2)
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
