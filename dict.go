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

// *Dict represents a Python dictionary.  In addition to satisfying the Object
// interface, Dict pointers also have a number of methods defined - representing
// the PyDict_XXX functions from the Python C API.
type Dict struct {
	AbstractObject
	MappingProtocol
	o C.PyDictObject
}

// DictType is the Type object that represents the Dict type.
var DictType = (*Type)(unsafe.Pointer(&C.PyDict_Type))

func dictCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.dictCheck(c(obj)) != 0
}

func newDict(obj *C.PyObject) *Dict {
	return (*Dict)(unsafe.Pointer(obj))
}

// NewDict creates a new empty dictionary.
//
// Return value: New Reference.
func NewDict() (*Dict, error) {
	ret := C.PyDict_New()
	if ret == nil {
		return nil, exception()
	}
	return newDict(ret), nil
}

func NewDictProxy(obj Object) (*Dict, error) {
	ret := C.PyDictProxy_New(c(obj))
	if ret == nil {
		return nil, exception()
	}
	return newDict(ret), nil
}

// CheckExact returns true if d is an actual dictionary object, and not an
// instance of a sub type.
func (d *Dict) CheckExact() bool {
	ret := C.dictCheckE(c(d))
	if int(ret) != 0 {
		return true
	}
	return false
}

// Clear empties the dictionary d of all key-value pairs.
func (d *Dict) Clear() {
	C.PyDict_Clear(c(d))
}

// Contains Returns true if the dictionary contains the given key.  This is
// equivalent to the Python expression "key in d".
func (d *Dict) Contains(key Object) (bool, error) {
	ret := C.PyDict_Contains(c(d), c(key))
	return int2BoolErr(ret)
}

// Copy returns a new dictionary that contains the same key-values pairs as d.
//
// Return value: New Reference.
func (d *Dict) Copy() (Object, error) {
	ret := C.PyDict_Copy(c(d))
	return obj2ObjErr(ret)
}

// SetItem inserts "val" into dictionary d with the key "key".  If "key" is not
// hashable, then a TypeError will be returned.
func (d *Dict) SetItem(key, val Object) error {
	ret := C.PyDict_SetItem(c(d), c(key), c(val))
	return int2Err(ret)
}

// SetItemString inserts "val" into dictionary d with the key "key" (or rather,
// with a *String with the value of "key" will be used as the key).  If "key" is
// not hashable, then a TypeError will be returned.
func (d *Dict) SetItemString(key string, val Object) error {
	s := C.CString(key)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyDict_SetItemString(c(d), s, c(val))
	return int2Err(ret)
}

// DelItem removes the entry with the key of "key" from the dictionary d.  If
// "key" is not hashable, a TypeError is returned.
func (d *Dict) DelItem(key Object) error {
	ret := C.PyDict_DelItem(c(d), c(key))
	return int2Err(ret)
}

// DelItem removes the entry with the key of "key" (or rather, with a *String
// with the value of "key" as the key) from the dictionary d.
func (d *Dict) DelItemString(key string) error {
	s := C.CString(key)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyDict_DelItemString(c(d), s)
	return int2Err(ret)
}

// GetItem returns the Object from dictionary d which has the key "key".  If
// there is no such Object, then nil is returned (without an error).
//
// Return value: Borrowed Reference.
func (d *Dict) GetItem(key Object) (Object, error) {
	ret := C.PyDict_GetItem(c(d), c(key))
	return obj2ObjErr(ret)
}

// GetItemString returns the Object from dictionary d which has the key "key"
// (or rather, which has a *String with the value of "key" as the key).  If
// there is no such Object, then nil is returned (without an error).
//
// Return value: Borrowed Reference.
func (d *Dict) GetItemString(key string) (Object, error) {
	s := C.CString(key)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyDict_GetItemString(c(d), s)
	return obj2ObjErr(ret)
}

// Items returns a *List containing all the items from the dictionary d, as with
// the Python "d.items()".
//
// Return value: New Reference.
func (d *Dict) Items() (*List, error) {
	ret := C.PyDict_Items(c(d))
	return newList(ret), exception()
}

// Keys returns a *List containing all the keys from the dictionary d, as with
// the Python "d.keys()".
//
// Return value: New Reference.
func (d *Dict) Keys() (*List, error) {
	ret := C.PyDict_Keys(c(d))
	return newList(ret), exception()
}

// Values returns a *List containing all the values from the dictionary d, as
// with the Python "d.values()".
//
// Return value: New Reference.
func (d *Dict) Values() (*List, error) {
	ret := C.PyDict_Values(c(d))
	return newList(ret), exception()
}

// Size returns the number of items in the dictionary d.  This is equivalent to
// the Python "len(d)".
func (d *Dict) Size() int64 {
	ret := C.PyDict_Size(c(d))
	if ret < 0 {
		panic(exception())
	}
	return int64(ret)
}

// PyDict_Next

// Merge merges key values pairs from Object o (which may be a dictionary, or an
// object that supports "o.keys()" and "o[key]") into the dictionary d.  If
// override is true then a matching key in d will have it's value replaced by
// the one in o, else the value in d will be left.
func (d *Dict) Merge(o Object, override bool) error {
	over := 0
	if override {
		over = 1
	}
	ret := C.PyDict_Merge(c(d), c(o), C.int(over))
	return int2Err(ret)
}

// Update replaces key values pairs in d with those from o.  It is equivalent to
// d.Merge(o, true) in Go, or "d.update(o)" in Python.
func (d *Dict) Update(o Object) error {
	ret := C.PyDict_Update(c(d), c(o))
	return int2Err(ret)
}

// MergeFromSeq2 merges key values pairs from the Object o (which must be an
// iterable object, where each item is an iterable of length 2 - the key value
// pairs).  If override is true then the last key value pair with the same key
// wins, otherwise the first instance does (where an instance already in d
// counts before any in o).
func (d *Dict) MergeFromSeq2(o Object, override bool) error {
	over := 0
	if override {
		over = 1
	}
	ret := C.PyDict_MergeFromSeq2(c(d), c(o), C.int(over))
	return int2Err(ret)
}

// Map returns a Go map that contains the values from the Python dictionary,
// indexed by the keys.  The keys and values are the same as in the Python
// dictionary, but changes to the Go map are not propogated back to the Python
// dictionary.
//
// Note: the map holds borrowed references
func (d *Dict) Map() map[Object]Object {
	m := make(map[Object]Object, d.Size())
	var p C.Py_ssize_t
	var k *C.PyObject
	var v *C.PyObject
	for int(C.PyDict_Next(c(d), &p, &k, &v)) != 0 {
		key := newObject(k)
		value := newObject(v)
		m[key] = value
	}
	return m
}

// MapString is similar to Map, except that the keys are first converted to
// strings.  If the keys are not all Python strings, then an error is returned.
//
// Note: the map holds borrowed references
func (d *Dict) MapString() (map[string]Object, error) {
	m := make(map[string]Object, d.Size())
	var p C.Py_ssize_t
	var k *C.PyObject
	var v *C.PyObject
	for int(C.PyDict_Next(c(d), &p, &k, &v)) != 0 {
		key := newObject(k)
		value := newObject(v)
		s, ok := key.(*String)
		if !ok {
			return nil, TypeError.Err("%v is not a string", key)
		}
		m[s.String()] = value
	}
	return m, nil
}

// String returns a string representation of the contents of the dictionary d.
func (d *Dict) String() string {
	if d == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", d.Map())
}
