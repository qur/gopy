// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

// *Set represents a Python set.  In addition to satisfying the Object
// interface, Set pointers also have a number of methods defined - representing
// the PySet_XXX functions from the Python C API.
type Set struct {
	AbstractObject
	o C.PySetObject
}

type FrozenSet struct {
	Set
}

// SetType is the Type object that represents the Set type.
var SetType = (*Type)(unsafe.Pointer(&C.PySet_Type))

// FrozenSetType is the Type object that represents the FrozenSet type.
var FrozenSetType = (*Type)(unsafe.Pointer(&C.PyFrozenSet_Type))

func setCheck(obj Object) bool {
	return C.setCheck(c(obj)) != 0
}

func frozenSetCheck(obj Object) bool {
	return C.frozenSetCheck(c(obj)) != 0
}

func newSet(obj *C.PyObject) *Set {
	return (*Set)(unsafe.Pointer(obj))
}

func newFrozenSet(obj *C.PyObject) *FrozenSet {
	return (*FrozenSet)(unsafe.Pointer(obj))
}

// NewSet create a new Python set instance.  The set contains the values from
// the passed Object if it is iterable (a TypeError is returned if "o" is not
// iterable).  An empty set is returned if o is nil.
//
// Return value: New Reference.
func NewSet(o Object) (*Set, error) {
	ret := C.PySet_New(c(o))
	if ret == nil {
		return nil, exception()
	}
	return newSet(ret), nil
}

// NewFrozenSet create a new Python frozenset instance.  The set contains the
// values from the passed Object if it is iterable (a TypeError is returned if
// "o" is not iterable).  An empty set is returned if o is nil.
//
// Return value: New Reference.
func NewFrozenSet(o Object) (*FrozenSet, error) {
	ret := C.PyFrozenSet_New(c(o))
	if ret == nil {
		return nil, exception()
	}
	return newFrozenSet(ret), nil
}

func (f *FrozenSet) CheckExact() bool {
	ret := C.frozenSetCheckE(c(f))
	return ret != 0
}

// Size returns the number of elements in the set "s".  This is equivalent to
// the Python "len(s)".
func (s *Set) Size() int64 {
	ret := C.PySet_Size(c(s))
	if ret < 0 {
		panic(exception())
	}
	return int64(ret)
}

// Contains returns true if the set "s" contains the Object "key".  "key" must
// be hashable, otherwise a TypeError is returned.  This is equivalent to the
// Python "key in s".
func (s *Set) Contains(key Object) (bool, error) {
	ret := C.PySet_Contains(c(s), c(key))
	return int2BoolErr(ret)
}

// Add adds "key" to the set "s".  "key" must be hashable, otherwise a TypeError
// is returned.
func (s *Set) Add(key Object) error {
	ret := C.PySet_Add(c(s), c(key))
	return int2Err(ret)
}

// Discard removes the specified "key" from the set "s".  It returns true if the
// Object was found and removed, false if it was not found.  "key" must be
// hashable, otherwise a TypeError is returned.
func (s *Set) Discard(key Object) (bool, error) {
	ret := C.PySet_Discard(c(s), c(key))
	return int2BoolErr(ret)
}

// Pop returns a new refereence to an arbitrary Object from the set, and removes
// it.  If the set is empty a KeyError is returned.
//
// Return value: New Reference.
func (s *Set) Pop() (Object, error) {
	ret := C.PySet_Pop(c(s))
	return obj2ObjErr(ret)
}

// Clear empties the set "s".
func (s *Set) Clear() error {
	ret := C.PySet_Clear(c(s))
	return int2Err(ret)
}
