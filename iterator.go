package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

type Iterator interface {
	Object
	Next() (Object, error)
}

type IteratorProtocol interface {
	Object
	AsIterator() *IteratorMethods
}

type IteratorMethods struct {
	abstractObject
	o C.PyObject
}

func newIterator(o *C.PyObject) (Iterator, error) {
	obj := newObject(o)
	if n, ok := obj.(Iterator); ok {
		return n, nil
	}
	if C.iterCheck(c(obj)) > 0 {
		return (*IteratorMethods)(unsafe.Pointer(obj.Base())), nil
	}
	return nil, TypeError.Err("object of type '%T' is not an iterator", obj)
}

func GetIterator(obj Object) (Iterator, error) {
	if i := AsIterator(obj); i != nil {
		return i, nil
	}
	if n, ok := obj.(Iterable); ok {
		return n.Iter()
	}
	return nil, TypeError.Err("cannot get iterator from %T", obj)
}

func AsIterator(obj Object) Iterator {
	if n, ok := obj.(Iterator); ok {
		return n
	}
	if n, ok := obj.(IteratorProtocol); ok {
		return n.AsIterator()
	}
	if C.iterCheck(c(obj)) > 0 {
		return (*IteratorMethods)(unsafe.Pointer(obj.Base()))
	}
	return nil
}

func (i *IteratorMethods) Next() (Object, error) {
	ret := C.PyIter_Next(c(i))
	return obj2ObjErr(ret)
}

// Iterate will iterate through the given iterator and return the values. If an
// error is encountered at any point, then iteration stops, the error is
// returned with all values collected so far.
func Iterate(i Iterator) ([]Object, error) {
	values := []Object{}
	for {
		item, err := i.Next()
		if item != nil {
			values = append(values, item)
		}
		if err != nil || item == nil {
			return values, err
		}
	}
}
