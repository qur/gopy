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

type Code struct {
	abstractObject
	o C.PyCodeObject
}

// CodeType is the Type object that represents the Code type.
var CodeType = (*Type)(unsafe.Pointer(&C.PyCode_Type))

func newCode(obj *C.PyObject) *Code {
	return (*Code)(unsafe.Pointer(obj))
}

func CompileFile(name string) (*Code, error) {
	// TODO(jp3): this should now be parsing the file directly, and then using
	// Py_CompileString(s, filename, C.Py_file_input)
	// fn := C.CString(name)
	// defer C.free(unsafe.Pointer(fn))
	// ret := C.compileFile(fn)
	// if ret == nil {
	// 	return nil, exception()
	// }
	// return newCode(ret), nil
	return nil, fmt.Errorf("func CompileFile not implemented")
}

func codeCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.codeCheck(c(obj)) != 0
}

func (code *Code) Eval(globals, locals Object) (Object, error) {
	ret := C.PyEval_EvalCode(c(code), c(globals), c(locals))
	return obj2ObjErr(ret)
}

func (code *Code) NumFree() int {
	pyCode := (*C.PyCodeObject)(unsafe.Pointer(code))
	n := C.PyCode_GetNumFree_(pyCode)
	return int(n)
}
