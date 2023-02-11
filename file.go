// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

type File struct {
	AbstractObject
	o C.PyFileObject
}

// FileType is the Type object that represents the File type.
var FileType = (*Type)(unsafe.Pointer(&C.PyFile_Type))

func newFile(obj *C.PyObject) *File {
	return (*File)(unsafe.Pointer(obj))
}

func fileCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.fileCheck(c(obj)) != 0
}

func NewFile(filename, mode string) (*File, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	cMode := C.CString(mode)
	defer C.free(unsafe.Pointer(cMode))
	ret := C.PyFile_FromString(cFilename, cMode)
	if ret == nil {
		return nil, exception()
	}
	return newFile(ret), nil
}

func (f *File) IncUseCount() {
	pyFile := (*C.PyFileObject)(unsafe.Pointer(f))
	C.PyFile_IncUseCount(pyFile)
}

func (f *File) DecUseCount() {
	pyFile := (*C.PyFileObject)(unsafe.Pointer(f))
	C.PyFile_DecUseCount(pyFile)
}

func (f *File) GetLine(n int) (Object, error) {
	ret := C.PyFile_GetLine(c(f), C.int(n))
	return obj2ObjErr(ret)
}

func (f *File) Name() Object {
	ret := C.PyFile_Name(c(f))
	return newObject(ret)
}

func (f *File) SetBufSize(n int) {
	C.PyFile_SetBufSize(c(f), C.int(n))
}

func (f *File) SetEncoding(enc string) bool {
	cEnc := C.CString(enc)
	defer C.free(unsafe.Pointer(cEnc))
	ret := C.PyFile_SetEncoding(c(f), cEnc)
	if ret != 0 {
		return false
	}
	return true
}

func (f *File) SetEncodingAndErrors(enc, errors string) bool {
	cEnc := C.CString(enc)
	defer C.free(unsafe.Pointer(cEnc))
	cErrors := C.CString(enc)
	defer C.free(unsafe.Pointer(cErrors))
	ret := C.PyFile_SetEncodingAndErrors(c(f), cEnc, cErrors)
	if ret != 0 {
		return false
	}
	return true
}

func (f *File) SoftSpace(newflag int) int {
	ret := C.PyFile_SoftSpace(c(f), C.int(newflag))
	return int(ret)
}

func (f *File) WriteObject(obj Object, flags int) error {
	ret := C.PyFile_WriteObject(c(obj), c(f), C.int(flags))
	return int2Err(ret)
}

func (f *File) WriteString(s string) error {
	cS := C.CString(s)
	defer C.free(unsafe.Pointer(cS))
	ret := C.PyFile_WriteString(cS, c(f))
	return int2Err(ret)
}
