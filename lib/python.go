// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #cgo CFLAGS: -Werror
// #cgo LDFLAGS: -lpython2.6 -lffi
// #include "utils.h"
import "C"

import (
	"unsafe"
)

func Initialize() {
	C.Py_Initialize()
}

func InitializeEx(initsigs bool) {
	if initsigs {
		C.Py_InitializeEx(1)
	} else {
		C.Py_InitializeEx(0)
	}
}

func AddToPath(dir string) {
	p := C.CString("path")
	defer C.free(unsafe.Pointer(p))

	sys_path := C.PySys_GetObject(p)
	if sys_path == nil {
		return
	}

	s := C.CString(dir)
	defer C.free(unsafe.Pointer(s))

	pDir := C.PyString_FromString(s)
	if pDir == nil {
		return
	}

	C.PyList_Append(sys_path, pDir)
}

func Main(args []string) int {
	argv := make([]*C.char, len(args))

	for i, arg := range args {
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	return int(C.Py_Main(C.int(len(argv)), &argv[0]))
}
