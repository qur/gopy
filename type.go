// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int typeCheck(PyObject *o) { return PyType_Check(o); }
// static inline int typeCheckE(PyObject *o) { return PyType_CheckExact(o); }
// static inline PyObject *typeAlloc(PyTypeObject *t, Py_ssize_t n) {
//    return t->tp_alloc(t, n);
// }
// static inline int typeInit(PyTypeObject *t, PyObject *o, PyObject *a, PyObject *k) {
//    return t->tp_init(o, a, k);
// }
import "C"

import "unsafe"

type Type struct {
	AbstractObject
	o *C.PyTypeObject
}

var typeObjMap = make(map[*C.PyObject]*Type)

// TypeType is the Type object that represents the Type type.
var TypeType = newType((*C.PyObject)(unsafe.Pointer(C.getBasePyType(C.GoPyType_Type))))

func typeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.typeCheck(c(obj)) != 0
}

func newType(obj *C.PyObject) *Type {
	if t, ok := typeObjMap[obj]; ok {
		return t
	}
	t := &Type{o: (*C.PyTypeObject)(unsafe.Pointer(obj))}
	typeObjMap[obj] = t
	return t
}

func (t *Type) Alloc(n int64) (Object, error) {
	ret := C.typeAlloc(t.o, C.Py_ssize_t(n))
	return obj2ObjErr(ret)
}

func (t *Type) Init(obj Object, args *Tuple, kw *Dict) error {
	ret := C.typeInit(t.o, c(obj), c(args), c(kw))
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
	return stringify(t)
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
	C.PyType_Modified(t.o)
}

// HasFeature returns true when "t" has the feature in question.
func (t *Type) HasFeature(feature uint32) bool {
	return (t.o.tp_flags & C.ulong(feature)) != 0
}

// IsGc returns true if the type "t" supports Cyclic Garbage Collection.
func (t *Type) IsGc() bool {
	return t.HasFeature(TPFLAGS_HAVE_GC)
}

// IsSubtype returns true if "t" is a subclass of "t2".
func (t *Type) IsSubtype(t2 *Type) bool {
	return C.PyType_IsSubtype(t.o, t2.o) == 1
}

// PyType_GenericAlloc : This is an internal function, which we should not need
// to expose.

// PyType_GenericNew : Another internal function we don't need to expose.

// PyType_Ready : This function is wrapped (along with a lot of other
// functionality) in the Create method of the Class stuct.
