package py

// #include "utils.h"
import "C"

import (
	"fmt"
)

func buildTuple(format string, args ...interface{}) (*Tuple, error) {
	if format == "" {
		return NewTuple(0)
	}
	bv, err := BuildValue(format, args...)
	if err != nil {
		return nil, err
	}
	t, ok := bv.(*Tuple)
	if ok {
		return t, nil
	}
	return PackTuple(bv)
}

// NewTuple returns a new *Tuple of the specified size.  However the entries are
// all set to NULL, so the tuple should not be shared, especially with Python
// code, until the entries have all been set.
//
// Return value: New Reference.
func NewTuple(size int) (*Tuple, error) {
	ret := C.PyTuple_New(C.Py_ssize_t(size))
	if ret == nil {
		return nil, exception()
	}
	return newTuple(ret), nil
}

// PackTuple returns a new *Tuple which contains the arguments. This tuple is
// ready to use. This function steals the references to the Objects.
//
// Return value: New Reference.
func PackTuple(items ...Object) (*Tuple, error) {
	t, err := NewTuple(len(items))
	if err != nil {
		return nil, err
	}

	for i, item := range items {
		if err := t.SetIndexSteal(i, item); err != nil {
			t.Decref()
			return nil, err
		}
	}

	return t, nil
}

// NewTupleFromValues creates a new Python Tuple from the supplied values. The
// values are converted to Objects using NewValue.
//
// Return value: New Reference.
func NewTupleFromValues(values ...any) (*Tuple, error) {
	t, err := NewTuple(len(values))
	if err != nil {
		return nil, err
	}
	for i, v := range values {
		o, err := NewValue(v)
		if err != nil {
			t.Decref()
			return nil, err
		}
		if err := t.SetIndexSteal(i, o); err != nil {
			t.Decref()
			return nil, err
		}
	}
	return t, nil
}

// CheckExact returns true if the object is actually a Tuple, not a subtype.
func (t *Tuple) CheckExact() bool {
	ret := C.tupleCheckE(c(t))
	return int(ret) != 0
}

// BorrowIndex returns the Object contained in tuple t at index idx. If idx is
// out of bounds for t, then an IndexError will be returned.
//
// Return value: Borrowed Reference.
func (t *Tuple) BorrowIndex(pos int) (Object, error) {
	ret := C.PyTuple_GetItem(c(t), C.Py_ssize_t(pos))
	return obj2ObjErr(ret)
}

// func (t *Tuple) GetSlice(low, high int) (*Tuple, error) {
// 	ret := C.PyTuple_GetSlice(c(t), C.Py_ssize_t(low), C.Py_ssize_t(high))
// 	if ret == nil {
// 		return nil, exception()
// 	}
// 	return newTuple(ret), nil
// }

// SetIndexSteal sets the Object at index idx in tuple t to Object obj.
//
// Note: This method "steals" a reference to obj, and discards a reference to
// the current value of idx in t (if there is one).
func (t *Tuple) SetIndexSteal(pos int, obj Object) error {
	ret := C.PyTuple_SetItem(c(t), C.Py_ssize_t(pos), c(obj))
	return int2Err(ret)
}

// _PyTuple_Resize

// PyTuple_ClearFreeList()

// Slice returns the Tuple contents as a slice of Object.
//
// Return value: Borrowed References.
func (t *Tuple) Slice() []Object {
	l := t.Size()
	s := make([]Object, l)
	for i := 0; i < l; i++ {
		o, err := t.BorrowIndex(i)
		if err != nil {
			panic(err)
		}
		s[i] = o
	}
	return s
}

// String will return a string representation of t.
func (t *Tuple) String() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", t.Slice())
}
