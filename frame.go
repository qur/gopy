// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

func GetFrame() (*Frame, error) {
	ret := C.PyEval_GetFrame()
	if ret == nil {
		return nil, exception()
	}
	return (*Frame)(unsafe.Pointer(ret)), nil
}

func (f *Frame) GetLineNumber() int {
	return int(C.PyFrame_GetLineNumber(&f.o))
}
