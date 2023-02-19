package py

// #include "utils.h"
import "C"

func GetBuiltins() (Object, error) {
	ret := C.PyEval_GetBuiltins()
	return obj2ObjErr(ret)
}

func GetLocals() (Object, error) {
	ret := C.PyEval_GetLocals()
	return obj2ObjErr(ret)
}

func GetGlobals() (Object, error) {
	ret := C.PyEval_GetGlobals()
	return obj2ObjErr(ret)
}
