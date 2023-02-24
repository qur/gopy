package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

func newException(obj *C.PyObject) *ExceptionClass {
	return (*ExceptionClass)(unsafe.Pointer(obj))
}

// ErrV returns a new Error of the specified kind (e), and with the given value.
func (e *ExceptionClass) ErrV(obj Object) *Error {
	return NewErrorV(e, obj)
}

// Err returns a new Error of the specified kind (e), and with the value being a
// new String containing the string created with the given format and args.
func (e *ExceptionClass) Err(format string, args ...interface{}) *Error {
	return NewError(e, format, args...)
}

// NewException creates and returns a new exception class. The name argument
// must be the name of the new exception, of the form "module.classname". The
// base and dict arguments are normally nil. This creates a class object derived
// from Exception.
//
// The __module__ attribute of the new class is set to the first part (up to the
// last dot) of the name argument, and the class name is set to the last part
// (after the last dot). The base argument can be used to specify alternate base
// classes; it can either be only one class or a tuple of classes. The dict
// argument can be used to specify a dictionary of class variables and methods.
//
// Return value: New Reference.
func NewException(name string, base, dict Object) (*ExceptionClass, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ret := C.PyErr_NewException(cName, c(base), c(dict))
	if ret == nil {
		return nil, exception()
	}
	return newException(ret), nil
}

// NewExceptionWithDoc creates and returns a new exception class. This is the
// same as NewException, except that docstring can also be supplied. If
// non-empty it will be used as the docstring for the exception class.
//
// Return value: New Reference.
func NewExceptionWithDoc(name, doc string, base, dict Object) (*ExceptionClass, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var cDoc *C.char
	if doc != "" {
		cDoc = C.CString(doc)
		defer C.free(unsafe.Pointer(cDoc))
	}
	ret := C.PyErr_NewExceptionWithDoc(cName, cDoc, c(base), c(dict))
	if ret == nil {
		return nil, exception()
	}
	return newException(ret), nil
}
