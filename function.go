package py

// #include "utils.h"
import "C"

// NewFunction returns a new Function object that is associated with the given
// "code" and "globals".  "globals" must be a dictionary, and it should hold the
// global variables for the function.
//
// The function inherits its docstring, name and __module__ from the "code"
// object, the argument defaults and closure are initialised to nil.
//
// Return value: New Reference.
func NewFunction(code Object, globals Object) (*Function, error) {
	ret := C.PyFunction_New(c(code), c(globals))
	return newFunction(ret), nil
}

// Code returns the code object associated with the function "f".
//
// Return value: Borrowed Reference.
func (f *Function) Code() (Object, error) {
	ret := C.PyFunction_GetCode(c(f))
	return obj2ObjErr(ret)
}

// Globals returns the globals dictionary associated with the function "f".
//
// Return value: Borrowed Reference.
func (f *Function) Globals() (Object, error) {
	ret := C.PyFunction_GetGlobals(c(f))
	return obj2ObjErr(ret)
}

// Module returns the __module__ attribute of the function "f".
//
// Return value: Borrowed Reference.
func (f *Function) Module() (Object, error) {
	ret := C.PyFunction_GetModule(c(f))
	return obj2ObjErr(ret)
}

// Defaults returns the argument default values for the function "f". This may
// be nil or a Tuple of values.
//
// Return value: Borrowed Reference.
func (f *Function) Defaults() (Object, error) {
	ret := C.PyFunction_GetDefaults(c(f))
	return obj2ObjErr(ret)
}

// SetDefaults sets the argument default values for the function "f". "o"
// must be either a Tuple, or None.
func (f *Function) SetDefaults(o Object) error {
	ret := C.PyFunction_SetDefaults(c(f), c(o))
	return int2Err(ret)
}

// Closure returns the closure associated with function "f". This may be nil
// or a Tuple of Cell objects.
//
// Return value: Borrowed Reference.
func (f *Function) Closure() (Object, error) {
	ret := C.PyFunction_GetClosure(c(f))
	return obj2ObjErr(ret)
}

// SetClosure sets the closure associated with function "f". "o" must be
// either a Tuple of Cell objects, or None.
func (f *Function) SetClosure(o Object) error {
	ret := C.PyFunction_SetClosure(c(f), c(o))
	return int2Err(ret)
}
