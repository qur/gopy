// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #cgo CFLAGS: -Werror
// #cgo LDFLAGS: -lpython2.7
// #cgo pkg-config: libffi
//
// #include "utils.h"
import "C"

import (
	"os"
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

func Finalize() {
	C.Py_Finalize()
}

func Exit(status int) {
	C.Py_Finalize()
	os.Exit(status)
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

// EnterRecusiveCall marks a point where a recursive Go-level call is about to
// be performed.  It returns true if the recursive call is permitted, otherwise
// a Python exception is set and false is returned.  where is a string that will
// be appended to the RuntimeError set if the recursion limit has been exceeded
// (e.g. " in instance check").  This function needs to be called if the
// recursive function may not invoke Python code (which automatically tracks
// recursion depth).
func EnterRecursiveCall(where string) bool {
	s := C.CString(where)
	defer C.free(unsafe.Pointer(s))
	return C.enterRecursive(s) == 0
}

// LeaveRecursiveCall must be called after a recursive call that was indicated
// by EnterRecursiveCall.
func LeaveRecursiveCall() {
	C.leaveRecursive()
}
