// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include <Python.h>
// #include "utils.h"
import "C"

import (
//	"os"
//	"unsafe"
)

func Eval_GetBuiltins() Object {
	ret := C.PyEval_GetBuiltins()
	if ret == nil {
		return nil
	}
	return newBaseObject(ret).actual()
}

func Eval_GetGlobals() Object {
	ret := C.PyEval_GetGlobals()
	if ret == nil {
		return nil
	}
	return newBaseObject(ret).actual()
}
