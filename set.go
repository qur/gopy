package py

// #include "utils.h"
import "C"

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
