// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

type Unicode struct {
	AbstractObject
	SequenceProtocol
	o C.PyUnicodeObject
}

// UnicodeType is the Type object that represents the Unicode type.
var UnicodeType = (*Type)(unsafe.Pointer(&C.PyUnicode_Type))

func unicodeCheck(obj Object) bool {
	return C.unicodeCheck(c(obj)) != 0
}

func newUnicode(obj *C.PyObject) *Unicode {
	return (*Unicode)(unsafe.Pointer(obj))
}

func NewUnicode(s string) (*Unicode, error) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	ret := C.PyUnicode_FromString(cs)
	if ret == nil {
		return nil, exception()
	}
	return newUnicode(ret), nil
}

func (u *Unicode) String() string {
	if u == nil {
		return "<nil>"
	}
	s := C.PyUnicode_AsUTF8String(c(u))
	if s == nil {
		panic(exception())
	}
	defer C.decref(s)
	ret := C.PyString_AsString(s)
	if ret == nil {
		panic(exception())
	}
	return C.GoString(ret)
}

func (u *Unicode) Size() int64 {
	ret := C.PyUnicode_GetSize(c(u))
	return int64(ret)
}

func (u *Unicode) EncodeString(encoding, errors string) (Object, error) {
	var cEncoding, cErrors *C.char
	if encoding == "" {
		cEncoding = C.CString(encoding)
		defer C.free(unsafe.Pointer(cEncoding))
	}
	if errors != "" {
		cErrors = C.CString(errors)
		defer C.free(unsafe.Pointer(cErrors))
	}
	ret := C.PyUnicode_AsEncodedString(c(u), cEncoding, cErrors)
	return obj2ObjErr(ret)
}

func (u *Unicode) UTF8String() (Object, error) {
	ret := C.PyUnicode_AsUTF8String(c(u))
	return obj2ObjErr(ret)
}

func (u *Unicode) UTF16String() (Object, error) {
	ret := C.PyUnicode_AsUTF16String(c(u))
	return obj2ObjErr(ret)
}

func (u *Unicode) UTF32String() (Object, error) {
	ret := C.PyUnicode_AsUTF32String(c(u))
	return obj2ObjErr(ret)
}

func (u *Unicode) UnicodeEscapeString() (Object, error) {
	ret := C.PyUnicode_AsUnicodeEscapeString(c(u))
	return obj2ObjErr(ret)
}

func (u *Unicode) RawUnicodeEscapeString() (Object, error) {
	ret := C.PyUnicode_AsRawUnicodeEscapeString(c(u))
	return obj2ObjErr(ret)
}

func (u *Unicode) Latin1String() (Object, error) {
	ret := C.PyUnicode_AsLatin1String(c(u))
	return obj2ObjErr(ret)
}

func (u *Unicode) ASCIIString() (Object, error) {
	ret := C.PyUnicode_AsASCIIString(c(u))
	return obj2ObjErr(ret)
}

func (u *Unicode) CharmapString(mapping Object) (Object, error) {
	ret := C.PyUnicode_AsCharmapString(c(u), c(mapping))
	return obj2ObjErr(ret)
}

func (u *Unicode) Encode(encoding, errors string) (Object, error) {
	var cEncoding, cErrors *C.char
	if encoding == "" {
		cEncoding = C.CString(encoding)
		defer C.free(unsafe.Pointer(cEncoding))
	}
	if errors != "" {
		cErrors = C.CString(errors)
		defer C.free(unsafe.Pointer(cErrors))
	}
	ret := C.PyUnicode_AsEncodedObject(c(u), cEncoding, cErrors)
	return obj2ObjErr(ret)
}

func (left *Unicode) Concat(right Object) (Object, error) {
	ret := C.PyUnicode_Concat(c(left), c(right))
	return obj2ObjErr(ret)
}

func (u *Unicode) Split(sep Object, maxsplit int64) (Object, error) {
	ret := C.PyUnicode_Split(c(u), c(sep), C.Py_ssize_t(maxsplit))
	return obj2ObjErr(ret)
}

func (u *Unicode) Splitlines(keepend bool) (Object, error) {
	cKeepend := C.int(0)
	if keepend {
		cKeepend = 1
	}
	ret := C.PyUnicode_Splitlines(c(u), cKeepend)
	return obj2ObjErr(ret)
}

func (u *Unicode) Translate(table Object, errors string) (Object, error) {
	var cErrors *C.char
	if errors != "" {
		cErrors = C.CString(errors)
		defer C.free(unsafe.Pointer(cErrors))
	}
	ret := C.PyUnicode_Translate(c(u), c(table), cErrors)
	return obj2ObjErr(ret)
}

func (u *Unicode) Join(seq Object) (Object, error) {
	ret := C.PyUnicode_Join(c(u), c(seq))
	return obj2ObjErr(ret)
}

func (u *Unicode) Tailmatch(substr Object, start, end int64, direction int) (bool, error) {
	ret := C.PyUnicode_Tailmatch(c(u), c(substr), C.Py_ssize_t(start), C.Py_ssize_t(end), C.int(direction))
	return int2BoolErr(C.int(ret))
}

func (u *Unicode) Find(substr Object, start, end int64, direction int) (int64, bool, error) {
	ret := C.PyUnicode_Find(c(u), c(substr), C.Py_ssize_t(start), C.Py_ssize_t(end), C.int(direction))
	if ret >= 0 {
		return int64(ret), true, nil
	} else if ret == -1 {
		return 0, false, nil
	}
	return 0, false, exception()
}

func (u *Unicode) Count(substr Object, start, end int64) (int64, error) {
	ret := C.PyUnicode_Count(c(u), c(substr), C.Py_ssize_t(start), C.Py_ssize_t(end))
	return ssize_t2Int64Err(ret)
}

func (u *Unicode) Replace(substr, replstr Object, maxcount int64) (Object, error) {
	ret := C.PyUnicode_Replace(c(u), c(substr), c(replstr), C.Py_ssize_t(maxcount))
	return obj2ObjErr(ret)
}

func (u *Unicode) Compare(right Object) (int, error) {
	ret := C.PyUnicode_Compare(c(u), c(right))
	return int(ret), exception()
}

func (u *Unicode) RichCompare(right Object, op Op) (Object, error) {
	ret := C.PyUnicode_RichCompare(c(u), c(right), C.int(op))
	return obj2ObjErr(ret)
}

func (u *Unicode) Format(args *Tuple) (*Unicode, error) {
	ret := C.PyUnicode_Format(c(u), c(args))
	if ret == nil {
		return nil, exception()
	}
	return newUnicode(ret), nil
}

func (u *Unicode) Contains(element Object) (bool, error) {
	ret := C.PyUnicode_Contains(c(u), c(element))
	return int2BoolErr(ret)
}
