package py

// #include "utils.h"
import "C"

import (
	"fmt"
)

// NewList creates a new Python List instance.  The created list has initial
// length "size".
//
// Note: If size > 0, then the objects in the returned list are initialised to
// nil. Thus you cannot use Abstract API functions, or expose the object to
// Python code without first filling in all the created slots with
// list.SetItem() or list.SetItemSteal().
//
// Return value: New Reference.
func NewList(size int) (*List, error) {
	ret := C.PyList_New(C.Py_ssize_t(size))
	if ret == nil {
		return nil, exception()
	}
	return newList(ret), nil
}

// NewListFromObjects create a new Python List from the supplied Objects. The
// references are stolen.
//
// Return value: New Reference.
func NewListFromObjects(objects ...Object) (*List, error) {
	l, err := NewList(len(objects))
	if err != nil {
		return nil, err
	}
	for i, o := range objects {
		if err := l.SetIndexSteal(i, o); err != nil {
			l.Decref()
			return nil, err
		}
	}
	return l, nil
}

// NewListFromValues creates a new Python List from the supplied values. The
// values are converted to Objects using NewValue.
//
// Return value: New Reference.
func NewListFromValues(values ...any) (*List, error) {
	rm := NewRefManager()
	defer rm.Decref()

	l, err := NewList(len(values))
	if err != nil {
		return nil, err
	}
	rm.Add(l)

	for i, v := range values {
		o, err := NewValue(v)
		if err != nil {
			return nil, err
		}
		if err := l.SetIndexSteal(i, o); err != nil {
			rm.Add(o)
			return nil, err
		}
	}

	rm.Clear()
	return l, nil
}

// NewListFromSlice creates a new Python List from the supplied slice. The
// values are converted to Objects using NewValue.
//
// Return value: New Reference.
func NewListFromSlice[T any](values []T) (*List, error) {
	rm := NewRefManager()
	defer rm.Decref()

	l, err := NewList(len(values))
	if err != nil {
		return nil, err
	}
	rm.Add(l)

	for i, v := range values {
		o, err := NewValue(v)
		if err != nil {
			return nil, err
		}
		if err := l.SetIndexSteal(i, o); err != nil {
			rm.Add(o)
			return nil, err
		}
	}

	rm.Clear()
	return l, nil
}

// CheckExact returns true if if l is an actual Python list, and not a sub type.
func (l *List) CheckExact() bool {
	ret := C.listCheckE(c(l))
	return ret == 1
}

// BorrowIndex returns the Object contained in list l at index idx. If idx is
// out of bounds for l, then an IndexError will be returned.
//
// Return value: Borrowed Reference.
func (l *List) BorrowIndex(idx int) (Object, error) {
	ret := C.PyList_GetItem(c(l), C.Py_ssize_t(idx))
	return obj2ObjErr(ret)
}

// SetIndexSteal sets the Object at index idx in list l to Object obj.
//
// Note: This method "steals" a reference to obj, and discards a reference to
// the current value of idx in l (if there is one).
func (l *List) SetIndexSteal(idx int, obj Object) error {
	ret := C.PyList_SetItem(c(l), C.Py_ssize_t(idx), c(obj))
	return int2Err(ret)
}

// Insert adds the Object obj to list l, by inserting it before the value
// currently stored at index idx (making obj the new value with index idx).
// This is equivalent to the Python "l.insert(idx, obj)".
func (l *List) Insert(idx int, obj Object) error {
	ret := C.PyList_Insert(c(l), C.Py_ssize_t(idx), c(obj))
	return int2Err(ret)
}

// Append adds the Object obj to list l, by appending it to the end of the list.
// This is equivalent to the Python "l.append(obj)"
func (l *List) Append(obj Object) error {
	ret := C.PyList_Append(c(l), c(obj))
	return int2Err(ret)
}

// Sort will sort the values in the List in place. The equivalent Python is
// "l.sort()".
func (l *List) Sort() error {
	ret := C.PyList_Sort(c(l))
	return int2Err(ret)
}

// Reverse will reverse the values in the List in place. The equivalent Python
// is "l.reverse()".
func (l *List) Reverse() error {
	ret := C.PyList_Reverse(c(l))
	return int2Err(ret)
}

// Tuple returns a new Tuple object containing the same contents as l. The
// equivalent Python is "tuple(list)".
//
// Return value: New Reference.
// func (l *List) Tuple() *Tuple {
// 	ret := C.PyList_AsTuple(c(l))
// 	if ret == nil {
// 		panic(exception())
// 	}
// 	return newTuple(ret)
// }

// Slice returns the list l as a Go Object slice.  The order of objects is
// copied from the Python list, but changes to the slice are not reflected in
// the Python list.
//
// Note: The returned slice contains borrowed references to the values.
func (l *List) Slice() []Object {
	size := l.Size()
	s := make([]Object, size)
	for i := 0; i < size; i++ {
		o, err := l.GetIndex(i)
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
