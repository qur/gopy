package py

// #include "utils.h"
import "C"

import (
	"os"
	"unsafe"
)

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

func (code *Code) Eval(globals, locals Object) (Object, error) {
	ret := C.PyEval_EvalCode(c(code), c(globals), c(locals))
	return obj2ObjErr(ret)
}

func (code *Code) NumFree() int {
	pyCode := (*C.PyCodeObject)(unsafe.Pointer(code))
	n := C.PyCode_GetNumFree_(pyCode)
	return int(n)
}
