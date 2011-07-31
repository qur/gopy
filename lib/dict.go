// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include <Python.h>
// static inline int dictCheck(PyObject *o) { return PyDict_Check(o); }
// static inline int dictCheckE(PyObject *o) { return PyDict_CheckExact(o); }
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type Dict struct {
	BaseObject
}

func dictCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.dictCheck(c(obj)) != 0
}

func newDict(obj *C.PyObject) *Dict {
	if obj == nil {
		return nil
	}
	return &Dict{BaseObject{obj}}
}

func Dict_New() (*Dict, os.Error) {
	ret := C.PyDict_New()
	if ret == nil {
		return nil, exception()
	}
	return newDict(ret), nil
}

func DictProxy_New(obj Object) (*Dict, os.Error) {
	ret := C.PyDictProxy_New(c(obj))
	if ret == nil {
		return nil, exception()
	}
	return newDict(ret), nil
}

func (d *Dict) CheckExact() bool {
	ret := C.dictCheckE(c(d))
	if int(ret) != 0 {
		return true
	}
	return false
}

func (d *Dict) Clear() {
	C.PyDict_Clear(c(d))
}

func (d *Dict) Contains(key Object) (bool, os.Error) {
	ret := C.PyDict_Contains(c(d), c(key))
	return int2BoolErr(ret)
}

func (d *Dict) Copy() (Object, os.Error) {
	ret := C.PyDict_Copy(c(d))
	return obj2ObjErr(ret)
}

func (d *Dict) SetItem(key, val Object) os.Error {
	ret := C.PyDict_SetItem(c(d), c(key), c(val))
	return int2Err(ret)
}

func (d *Dict) SetItemString(key string, val Object) os.Error {
	s := C.CString(key)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyDict_SetItemString(c(d), s, c(val))
	return int2Err(ret)
}

func (d *Dict) DelItem(key Object) os.Error {
	ret := C.PyDict_DelItem(c(d), c(key))
	return int2Err(ret)
}

func (d *Dict) DelItemString(key string) os.Error {
	s := C.CString(key)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyDict_DelItemString(c(d), s)
	return int2Err(ret)
}

func (d *Dict) GetItem(key Object) (Object, os.Error) {
	ret := C.PyDict_GetItem(c(d), c(key))
	return obj2ObjErr(ret)
}

func (d *Dict) GetItemString(key string) (Object, os.Error) {
	s := C.CString(key)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyDict_GetItemString(c(d), s)
	return obj2ObjErr(ret)
}

func (d *Dict) Items() (Object, os.Error) {
	ret := C.PyDict_Items(c(d))
	return obj2ObjErr(ret)
}

func (d *Dict) Keys() (Object, os.Error) {
	ret := C.PyDict_Keys(c(d))
	return obj2ObjErr(ret)
}

func (d *Dict) Values() (Object, os.Error) {
	ret := C.PyDict_Values(c(d))
	return obj2ObjErr(ret)
}

func (d *Dict) Size() int64 {
	ret := C.PyDict_Size(c(d))
	if ret < 0 {
		panic(exception())
	}
	return int64(ret)
}

// PyDict_Next

func (d *Dict) Merge(o Object, override bool) os.Error {
	over := 0
	if override {
		over = 1
	}
	ret := C.PyDict_Merge(c(d), c(o), C.int(over))
	return int2Err(ret)
}

func (d *Dict) Update(o Object) os.Error {
	ret := C.PyDict_Update(c(d), c(o))
	return int2Err(ret)
}

func (d *Dict) MergeFromSeq2(o Object, override bool) os.Error {
	over := 0
	if override {
		over = 1
	}
	ret := C.PyDict_MergeFromSeq2(c(d), c(o), C.int(over))
	return int2Err(ret)
}

// Note: the map holds borrowed references
func (d *Dict) Map() map[Object]Object {
	m := make(map[Object]Object, d.Size())
	var p C.Py_ssize_t
	var k *C.PyObject
	var v *C.PyObject
	for int(C.PyDict_Next(c(d), &p, &k, &v)) != 0 {
		key := newBaseObject(k).actual()
		value := newBaseObject(v).actual()
		m[key] = value
	}
	return m
}

// Note: the map holds borrowed references
func (d *Dict) MapString() (map[string]Object, os.Error) {
	m := make(map[string]Object, d.Size())
	var p C.Py_ssize_t
	var k *C.PyObject
	var v *C.PyObject
	for int(C.PyDict_Next(c(d), &p, &k, &v)) != 0 {
		key := newBaseObject(k).actual()
		value := newBaseObject(v).actual()
		s, ok := key.(*String)
		if !ok {
			return nil, fmt.Errorf("TypeError: %v is not a string", key)
		}
		m[s.String()] = value
	}
	return m, nil
}

func (d *Dict) String() string {
	return fmt.Sprintf("%v", d.Map())
}
