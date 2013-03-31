// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

type StartToken int

const (
	EvalInput StartToken = iota
	FileInput
	SingleInput
)

func RunString(code string, start StartToken, globals, locals Object) (Object, error) {
	codestr := C.CString(code)
	defer C.free(unsafe.Pointer(codestr))

	var token C.int
	switch start {
	case EvalInput:
		token = C.Py_eval_input
	case FileInput:
		token = C.Py_file_input
	case SingleInput:
		token = C.Py_single_input
	}

	obj := C.PyRun_StringFlags(codestr, token, c(globals), c(locals), nil)
	if obj == nil {
		return nil, exception()
	}

	return newObject(obj), nil
}

func RunFile(filename string, start StartToken, globals, locals Object) (Object, error) {
	name := C.CString(filename)
	defer C.free(unsafe.Pointer(name))

	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))

	var token C.int
	switch start {
	case EvalInput:
		token = C.Py_eval_input
	case FileInput:
		token = C.Py_file_input
	case SingleInput:
		token = C.Py_single_input
	}

	file, err := C.fopen(name, mode)
	if file == nil {
		return nil, err
	}
	defer C.fclose(file)

	obj := C.PyRun_FileExFlags(file, name, token, c(globals), c(locals), 0, nil)
	if obj == nil {
		return nil, exception()
	}

	return newObject(obj), nil
}
