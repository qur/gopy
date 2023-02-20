package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

type Iterator interface {
	Object
	AsIterator() *IteratorMethods
}

type IteratorMethods struct {
	abstractObject
	o C.PyObject
}

func AsIterator(obj Object) *IteratorMethods {
	if n, ok := obj.(Iterator); ok {
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

type AsyncIterator interface {
	Object
	AsAsyncIterator() *AsyncIteratorMethods
}

type AsyncIteratorMethods struct {
	abstractObject
	o C.PyObject
}

func AsAsyncIterator(obj Object) *AsyncIteratorMethods {
	if n, ok := obj.(AsyncIterator); ok {
		return n.AsAsyncIterator()
	}
	if C.aIterCheck(c(obj)) > 0 {
		return (*AsyncIteratorMethods)(unsafe.Pointer(obj.Base()))
	}
	return nil
}
