package py

// #include "utils.h"
import "C"

// GetBuiltins returns a dictionary of the current Python builtins.
//
// Return value: Borrowed Reference.
func GetBuiltins() (Object, error) {
	ret := C.PyEval_GetBuiltins()
	return obj2ObjErr(ret)
}

// GetLocals returns a dictionary of the current Python local variables.
//
// Return value: Borrowed Reference.
func GetLocals() (Object, error) {
	ret := C.PyEval_GetLocals()
	return obj2ObjErr(ret)
}

// GetGlobals returns a dictionary of the current Python global variables.
//
// Return value: Borrowed Reference.
func GetGlobals() (Object, error) {
	ret := C.PyEval_GetGlobals()
	return obj2ObjErr(ret)
}
