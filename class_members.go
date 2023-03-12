package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"reflect"
)

func getField(obj, arg *C.PyObject) (reflect.Value, reflect.StructField, error) {
	o := getClassObject(obj)
	if o == nil {
		return reflect.Value{}, reflect.StructField{}, fmt.Errorf("unknown object")
	}

	idx := int(C.PyLong_AsLong(arg))

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
func goClassNatSet(obj, obj2, idx *C.PyObject) int {
	f, t, err := getField(obj, idx)
	if err != nil {
		raise(err)
		return -1
	}

	// This is the new value we are being asked to set
	o := newObject(obj2)

	switch f.Type().Kind() {
	case reflect.Bool:
		b, ok := o.(*Bool)
		if !ok {
			raise(TypeError.Err("field %s is bool, got %s", t.Name, o.Type()))
			return -1
		}
		f.SetBool(b.Bool())
		return 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := int64(C.PyLong_AsLong(obj2))
		if exceptionRaised() {
			return -1
		}
		f.SetInt(v)
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		v := int64(C.PyLong_AsLong(obj2))
		if exceptionRaised() {
			return -1
		}
		if v < 0 {
			raise(ValueError.Err("field %s is unsigned, can't set to negative value %d", t.Name, v))
			return -1
		}
		f.SetUint(uint64(v))
		return 0
	case reflect.Float32, reflect.Float64:
		v := float64(C.PyFloat_AsDouble(obj2))
		if exceptionRaised() {
			return -1
		}
		f.SetFloat(v)
		return 0
	case reflect.String:
		v := C.PyUnicode_AsUTF8(obj2)
		if exceptionRaised() {
			return -1
		}
		f.SetString(C.GoString(v))
		return 0
	case reflect.Complex64, reflect.Complex128:
		v, ok := o.(*Complex)
		if !ok {
			raise(TypeError.Err("field %s is complex, got %s", t.Name, o.Type()))
			return -1
		}
		f.SetComplex(v.Complex128())
		return 0
	}

	raise(NotImplementedError.ErrV(None))
	return -1
}

//export goClassObjGet
func goClassObjGet(obj, idx *C.PyObject) *C.PyObject {
	f, _, err := getField(obj, idx)
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
