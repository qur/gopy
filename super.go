package py

// #include "utils.h"
import "C"

import "unsafe"

// CallMethod calls the method of s called name with the given args and kwds.
// kwds may be nil, args may not (an empty Tuple must be used if no arguments
// are wanted). Returns the result of the call, or an Error on failure. This is
// equivalent to "s.name(*args, **kwds)" in Python.
//
// Return value: New Reference.
func (s *Super) CallMethod(name string, args *Tuple, kwds *Dict) (Object, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	f := C.PyObject_GetAttrString(s.raw(), cname)
	if f == nil {
		return nil, AttributeError.Err(name)
	}
	defer C.decref(f)

	if C.PyCallable_Check(f) == 0 {
		return nil, TypeError.Err("attribute of type '%s' is not callable", name)
	}

	ret := C.PyObject_Call(f, c(args), c(kwds))
	return obj2ObjErr(ret)
}

// CallMethodGo calls the method of s called name with the given args and
// kwds, either may be nil. Returns the result of the call, or an Error on
// failure. This is equivalent to "s.name(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (s *Super) CallMethodGo(name string, args []any, kwds map[string]any) (Object, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	f := C.PyObject_GetAttrString(s.raw(), cname)
	if f == nil {
		return nil, AttributeError.Err(name)
	}
	defer C.decref(f)

	if C.PyCallable_Check(f) == 0 {
		return nil, TypeError.Err("attribute of type '%s' is not callable", name)
	}

	obj1, err := NewTupleFromValues(args...)
	if err != nil {
		return nil, err
	}
	defer obj1.Decref()

	obj2, err := NewDictFromValuesString(kwds)
	if err != nil {
		return nil, err
	}
	defer obj2.Decref()

	ret := C.PyObject_Call(f, c(obj1), c(obj2))
	return obj2ObjErr(ret)
}
