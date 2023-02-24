package py

// #include "utils.h"
import "C"

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

// Repr returns a String representation of "s". This is equivalent to the
// Python "repr(s)".
//
// Return value: New Reference.
func (f *FrozenSet) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(f))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Hash computes and returns the hash value of s. The equivalent
// Python is "hash(s)".
func (f *FrozenSet) Hash() (int, error) {
	ret := C.PyObject_Hash(c(f))
	if ret == -1 {
		return 0, exception()
	}
	return int(ret), nil
}

// HasAttr returns true if "s" has the attribute "name".  This is equivalent
// to the Python "hasattr(s, name)".
func (f *FrozenSet) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(f), c(name))
	return ret == 1
}

// GetAttr returns the attribute of "s" with the name "name".  This is
// equivalent to the Python "s.name".
//
// Return value: New Reference.
func (f *FrozenSet) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(f), c(name))
	return obj2ObjErr(ret)
}

// RichCompare compares "f" with "obj" using the specified operation (LE, GE
// etc.), and returns the result (True or False). The equivalent Python is "f
// op obj", where op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (f *FrozenSet) RichCompare(obj Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(f), c(obj), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompareBool compares "f" with "obj" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "f op obj", where
// op is the corresponding Python operator for op.
func (f *FrozenSet) RichCompareBool(obj Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(f), c(obj), C.int(op))
	return int2BoolErr(ret)
}

// Iter returns an Iterator that will iterate over the members of f.
func (f *FrozenSet) Iter() (Iterator, error) {
	ret := C.PyObject_GetIter(c(f))
	if ret == nil {
		return nil, exception()
	}
	return newIterator(ret)
}

func (f *FrozenSet) Size() int {
	ret := C.PyObject_Size(c(f))
	if ret < 0 {
		clearErr()
		return 0
	}
	return int(ret)
}

func (f *FrozenSet) Contains(obj Object) (bool, error) {
	ret := C.PySequence_Contains(c(f), c(obj))
	return int2BoolErr(ret)
}
