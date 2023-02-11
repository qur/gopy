// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

// *Function represents a Python function.  In Python this is a function created
// using the "def" statement.
type Function struct {
	AbstractObject
	o C.PyFunctionObject
}

// FunctionType is the Type object that represents the Function type.
var FunctionType = (*Type)(unsafe.Pointer(&C.PyFunction_Type))

func newFunction(obj *C.PyObject) *Function {
	return (*Function)(unsafe.Pointer(obj))
}

func functionCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.functionCheck(c(obj)) != 0
}

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

// Code returns the code object associated with the function "fn".
//
// Return value: Borrowed Reference.
func (fn *Function) Code() (Object, error) {
	ret := C.PyFunction_GetCode(c(fn))
	return obj2ObjErr(ret)
}

// Globals returns the globals dictionary associated with the function "fn".
//
// Return value: Borrowed Reference.
func (fn *Function) Globals() (Object, error) {
	ret := C.PyFunction_GetGlobals(c(fn))
	return obj2ObjErr(ret)
}

// Module returns the __module__ attribute of the function "fn".
//
// Return value: Borrowed Reference.
func (fn *Function) Module() (Object, error) {
	ret := C.PyFunction_GetModule(c(fn))
	return obj2ObjErr(ret)
}

// Defaults returns the argument default values for the function "fn".  This may
// be nil or a Tuple of values.
//
// Return value: Borrowed Reference.
func (fn *Function) Defaults() (Object, error) {
	ret := C.PyFunction_GetDefaults(c(fn))
	return obj2ObjErr(ret)
}

// SetDefaults sets teh arguement default values for the function "fn".  "o"
// must be either a Tuple, or None.
func (fn *Function) SetDefaults(o Object) error {
	ret := C.PyFunction_SetDefaults(c(fn), c(o))
	return int2Err(ret)
}

// Closure returns the closure associated with function "fn".  This may be nil
// or a Tuple of Cell objects.
//
// Return value: Borrowed Reference.
func (fn *Function) Closure() (Object, error) {
	ret := C.PyFunction_GetClosure(c(fn))
	return obj2ObjErr(ret)
}

// SetClosure sets the closure associated with function "fn".  "o" must be
// either a Tuple of Cell objects, or None.
func (fn *Function) SetClosure(o Object) error {
	ret := C.PyFunction_SetClosure(c(fn), c(o))
	return int2Err(ret)
}
