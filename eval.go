// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

func GetBuiltins() (Object, error) {
	ret := C.PyEval_GetBuiltins()
	return obj2ObjErr(ret)
}

func GetLocals() (Object, error) {
	ret := C.PyEval_GetLocals()
	return obj2ObjErr(ret)
}

func GetGlobals() (Object, error) {
	ret := C.PyEval_GetGlobals()
	return obj2ObjErr(ret)
}

func InitThreads() {
	C.PyEval_InitThreads()
}
