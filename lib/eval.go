// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

func GetBuiltins() (Object, error) {
	ret := C.PyEval_GetBuiltins()
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret), nil
}

func GetGlobals() (Object, error) {
	ret := C.PyEval_GetGlobals()
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret), nil
}

func InitThreads() {
	C.PyEval_InitThreads()
}
