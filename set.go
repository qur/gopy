package py

// #include "utils.h"
import "C"

import "unsafe"

type FrozenSet struct {
	Set
}

// FrozenSetType is the Type object that represents the FrozenSet type.
var FrozenSetType = (*Type)(unsafe.Pointer(&C.PyFrozenSet_Type))

func frozenSetCheck(obj Object) bool {
	return C.frozenSetCheck(c(obj)) != 0
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

// Contains returns true if the set "s" contains the Object "key".  "key" must
// be hashable, otherwise a TypeError is returned.  This is equivalent to the
// Python "key in s".
// func (s *Set) Contains(key Object) (bool, error) {
// 	ret := C.PySet_Contains(c(s), c(key))
// 	return int2BoolErr(ret)
// }

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
