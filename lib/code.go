// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline int codeCheck(PyObject *o) { return PyCode_Check(o); }
import "C"

import (
	"os"
	"unsafe"
)

type Code struct {
	AbstractObject
	C.PyCodeObject
}

var CodeType = (*Type)(unsafe.Pointer(&C.PyCode_Type))

func newCode(obj *C.PyObject) *Code {
	return (*Code)(unsafe.Pointer(obj))
}

func CompileFile(name string) (*Code, os.Error) {
	fn := C.CString(name)
	defer C.free(unsafe.Pointer(fn))
	ret := C.compileFile(fn)
	if ret == nil {
		return nil, exception()
	}
	return newCode(ret), nil
}

func codeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.codeCheck(c(obj)) != 0
}

func (code *Code) Eval(globals, locals Object) (Object, os.Error) {
	pyCode := (*C.PyCodeObject)(unsafe.Pointer(code))
	ret := C.PyEval_EvalCode(pyCode, c(globals), c(locals))
	return obj2ObjErr(ret)
}
