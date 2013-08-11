// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

type String struct {
	AbstractObject
	SequenceProtocol
	o C.PyStringObject
}

// StringType is the Type object that represents the String type.
var StringType = (*Type)(unsafe.Pointer(&C.PyString_Type))

func stringCheck(obj Object) bool {
	return C.stringCheck(c(obj)) != 0
}

func newString(obj *C.PyObject) *String {
	return (*String)(unsafe.Pointer(obj))
}

func NewString(s string) (*String, error) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	ret := C.PyString_FromString(cs)
	if ret == nil {
		return nil, exception()
	}
	return newString(ret), nil
}

func (s *String) String() string {
	if s == nil {
		return "<nil>"
	}
	ret := C.PyString_AsString(c(s))
	if ret == nil {
		panic(exception())
	}
	return C.GoString(ret)
}

func (s *String) Format(args *Tuple) (*String, error) {
	ret := C.PyString_Format(c(s), c(args))
	if ret == nil {
		return nil, exception()
	}
	return newString(ret), nil
}

func (s *String) Size() int64 {
	ret := C.PyString_Size(c(s))
	return int64(ret)
}

func (s *String) Decode(encoding, errors string) (Object, error) {
	var cEncoding, cErrors *C.char
	if encoding == "" {
		cEncoding = C.CString(encoding)
		defer C.free(unsafe.Pointer(cEncoding))
	}
	if errors != "" {
		cErrors = C.CString(errors)
		defer C.free(unsafe.Pointer(cErrors))
	}
	ret := C.PyString_AsDecodedObject(c(s), cEncoding, cErrors)
	return obj2ObjErr(ret)
}

func (s *String) Encode(encoding, errors string) (Object, error) {
	var cEncoding, cErrors *C.char
	if encoding == "" {
		cEncoding = C.CString(encoding)
		defer C.free(unsafe.Pointer(cEncoding))
	}
	if errors != "" {
		cErrors = C.CString(errors)
		defer C.free(unsafe.Pointer(cErrors))
	}
	ret := C.PyString_AsEncodedObject(c(s), cEncoding, cErrors)
	return obj2ObjErr(ret)
}
