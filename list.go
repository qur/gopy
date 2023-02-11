// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// *List represents a Python list.  In addition to satisfying the Object
// interface, List pointers also have a number of methods defined - representing
// the PyList_XXX functions from the Python C API.
type List struct {
	AbstractObject
	SequenceProtocol
	o C.PyListObject
}

// ListType is the Type object that represents the List type.
var ListType = (*Type)(unsafe.Pointer(&C.PyList_Type))

func listCheck(obj Object) bool {
	return C.listCheck(c(obj)) != 0
}

func newList(obj *C.PyObject) *List {
	return (*List)(unsafe.Pointer(obj))
}

// NewList creates a new Python List instance.  The created list has initial
// length "size".
//
// Note: If size > 0, then the objects in the returned list are initialised to
// nil.  Thus you cannot use Abstract API functions, or expose the object to
// Python code without first filling in all the created slots with
// list.SetItem().
//
// Return value: New Reference.
func NewList(size int64) (*List, error) {
	ret := C.PyList_New(C.Py_ssize_t(size))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

// CheckExact returns true if if l is an actual Python list, and not a sub type.
func (l *List) CheckExact() bool {
	ret := C.listCheckE(c(l))
	if int(ret) != 0 {
		return true
	}
	return false
}

// Size returns the number of elements in the list l.  This is equivalent to the
// Python "len(l)".
func (l *List) Size() int64 {
	ret := C.PyList_Size(c(l))
	if ret < 0 {
		panic(exception())
	}
	return int64(ret)
}

// GetItem returns the Object contained in list l at index idx.  If idx is out
// of bounds for l, then an IndexError will be returned.
//
// Return value: Borrowed Reference.
func (l *List) GetItem(idx int64) (Object, error) {
	ret := C.PyList_GetItem(c(l), C.Py_ssize_t(idx))
	return obj2ObjErr(ret)
}

// SetItem sets the Object at index idx in list l to Object obj.
//
// Note: This method "steals" a reference to obj, and discards a reference to
// the current value of idx in l (if there is one).
func (l *List) SetItem(idx int64, obj Object) error {
	ret := C.PyList_SetItem(c(l), C.Py_ssize_t(idx), c(obj))
	return int2Err(ret)
}

// Insert adds the Object obj to list l, by inserting it before the value
// currently stored at index idx (making obj the new value with index idx).
// This is equivalent to the Python "l.insert(idx, obj)".
func (l *List) Insert(idx int64, obj Object) error {
	ret := C.PyList_Insert(c(l), C.Py_ssize_t(idx), c(obj))
	return int2Err(ret)
}

// Append adds the Object obj to list l, by appending it to the end of the list.
// This is equivalent to the Python "l.append(obj)"
func (l *List) Append(obj Object) error {
	ret := C.PyList_Append(c(l), c(obj))
	return int2Err(ret)
}

func (l *List) GetSlice(low, high int64) (*List, error) {
	ret := C.PyList_GetSlice(c(l), C.Py_ssize_t(low), C.Py_ssize_t(high))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

func (l *List) SetSlice(low, high int64, items *List) error {
	ret := C.PyList_SetSlice(c(l), C.Py_ssize_t(low), C.Py_ssize_t(high), c(items))
	return int2Err(ret)
}

func (l *List) Sort() error {
	ret := C.PyList_Sort(c(l))
	return int2Err(ret)
}

func (l *List) Reverse() error {
	ret := C.PyList_Reverse(c(l))
	return int2Err(ret)
}

func (l *List) Tuple() *Tuple {
	ret := C.PyList_AsTuple(c(l))
	if ret == nil {
		panic(exception())
	}
	return newTuple(ret)
}

// Slice returns the list l as a Go Object slice.  The order of objects is
// copied from the Python list, but changes to the slice are not reflected in
// the Python list.
//
// Note: The returned slice contains borrowed references to the values.
func (l *List) Slice() []Object {
	size := l.Size()
	s := make([]Object, size)
	for i := int64(0); i < size; i++ {
		o, err := l.GetItem(i)
		if err != nil {
			panic(err)
		}
		s[i] = o
	}
	return s
}

// String returns a string representation of the list l.
func (l *List) String() string {
	if l == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", l.Slice())
}
