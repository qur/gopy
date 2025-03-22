package py

// #include "utils.h"
import "C"

import (
	"fmt"
)

type StartToken int

const (
	EvalInput StartToken = iota
	FileInput
	SingleInput
)

func (s StartToken) c() (C.int, error) {
	switch s {
	case EvalInput:
		return C.Py_eval_input, nil
	case FileInput:
		return C.Py_file_input, nil
	case SingleInput:
		return C.Py_single_input, nil
	default:
		return 0, fmt.Errorf("invalid StartToken: %d", s)
	}
}

func RunString(code string, start StartToken, globals, locals Object) (Object, error) {
	codestr := C.CString(code)
	defer cfree(codestr)

	token, err := start.c()
	if err != nil {
		return nil, err
	}

	obj := C.PyRun_StringFlags(codestr, token, c(globals), c(locals), nil)
	if obj == nil {
		return nil, exception()
	}

	return newObject(obj), nil
}

func RunFile(filename string, start StartToken, globals, locals Object) (Object, error) {
	name := C.CString(filename)
	defer cfree(name)

	mode := C.CString("r")
	defer cfree(mode)

	token, err := start.c()
	if err != nil {
		return nil, err
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
