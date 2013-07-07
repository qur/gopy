// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

type Type struct {
	AbstractObject
	o C.PyTypeObject
}

// TypeType is the Type object that represents the Type type.
var TypeType = (*Type)(unsafe.Pointer(&C.PyType_Type))

func typeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.typeCheck(c(obj)) != 0
}

func newType(obj *C.PyObject) *Type {
	return (*Type)(unsafe.Pointer(obj))
}

func (t *Type) Alloc(n int64) (Object, error) {
	ret := C.typeAlloc(c(t), C.Py_ssize_t(n))
	return obj2ObjErr(ret)
}

func (t *Type) Init(obj Object, args *Tuple, kw *Dict) error {
	ret := C.typeInit(c(t), c(obj), c(args), c(kw))
	if ret < 0 {
		return exception()
	}
	return nil
}

func (t *Type) String() string {
	pyS := C.PyObject_Str(c(t))
	if pyS == nil {
		return "<unknown type>"
	}
	return C.GoString(C.PyString_AsString(pyS))
}

// CheckExact returns true when "t" is an actual Type object, and not some form
// of subclass.
func (t *Type) CheckExact() bool {
	return C.typeCheckE(c(t)) == 1
}

// PyType_ClearCache : TODO - ???

// Modified should be called after the attributes or base class of a Type have
// been changed.
func (t *Type) Modified() {
	C.PyType_Modified(&t.o)
}

// HasFeature returns true when "t" has the feature in question.
func (t *Type) HasFeature(feature uint32) bool {
	return (t.o.tp_flags & C.long(feature)) != 0
}

// IsGc returns true if the type "t" supports Cyclic Garbage Collection.
func (t *Type) IsGc() bool {
	return t.HasFeature(TPFLAGS_HAVE_GC)
}

// IsSubtype returns true if "t" is a subclass of "t2".
func (t *Type) IsSubtype(t2 *Type) bool {
	return C.PyType_IsSubtype(&t.o, &t2.o) == 1
}

// PyType_GenericAlloc : This is an internal function, which we should not need
// to expose.

// PyType_GenericNew : Another internal function we don't need to expose.

// PyType_Ready : This function is wrapped (along with a lot of other
// functionality) in the Create method of the Class stuct.
