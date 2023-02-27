package py

// #include "utils.h"
import "C"

import (
	"os"
	"unsafe"
)

// CompileFile compiles the Python file at the given path.
//
// The returned Code can be turned into a Module using ExecCodeModule.
//
// Return value: New Reference.
func CompileFile(name string) (*Code, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	// we need data to be NUL terminated, and using append might be able to use
	// some spare space whereas using C.CString(string(data)) will definitely
	// copy at least once, maybe even twice.
	data = append(data, 0)

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.Py_CompileString((*C.char)(unsafe.Pointer(&data[0])), cName, C.Py_file_input)
	if ret == nil {
		return nil, exception()
	}

	return newCode(ret), nil
}

// Eval is a simplified function for evaluating a Code object. The code is
// executed in the environment of the given globals and locals.
//
// Return value: New Reference.
func (co *Code) Eval(globals, locals Object) (Object, error) {
	ret := C.PyEval_EvalCode(c(co), c(globals), c(locals))
	return obj2ObjErr(ret)
}
