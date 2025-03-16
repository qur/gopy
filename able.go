package py

// #include "utils.h"
import "C"

import "unsafe"

type Iterable interface {
	Object
	Iter() (Iterator, error)
}

// AsIterable returns obj as an Iterable. If obj is not an Iterable, then nil
// will be returned.
func AsIterable(obj Object) Iterable {
	if n, ok := obj.(Iterable); ok {
		return n
	}

	if obj.Type().o.tp_iter != nil {
		return (*IterableMethods)(unsafe.Pointer(obj.Base()))
	}

	// Sequences are always iterable, there is actually a dedicated iterator
	if s := AsSequenceMethods(obj); s != nil {
		return s
	}

	return nil
}

// Iter returns an iterator for the object represented by i.
//
// Return value: New Reference.
func (i *IterableMethods) Iter() (Iterator, error) {
	obj, err := obj2ObjErr(C.PyObject_GetIter(c(i)))
	if err != nil {
		return nil, err
	}

	it := AsIterator(obj)
	if it == nil {
		Decref(obj)
		return nil, TypeError.Err("__iter__ did not return iterator")
	}

	return it, nil
}

type Hashable interface {
	Object
	Hash() (Object, error)
}

type Callable interface {
	Object
	Call(args *Tuple, kwargs *Dict) (Object, error)
	CallGo(args []Object, kwargs map[string]Object) (Object, error)
}

type Comparable interface {
	Object
	RichCompare(obj Object, op Op) (Object, error)
}
