package py

// #include "utils.h"
import "C"

// None is the Python equivalent to nil.
var None = newNoneObject(&C._Py_NoneStruct)

// String returns a string representation of the NoneObject type. This function
// always returns "None".
func (*NoneObject) String() string {
	return "None"
}

// ReturnNone returns None (after incrementing the reference count).
// This is a convenience function for returning None from a function.
//
// Return value: New Reference.
func ReturnNone() Object {
	None.Incref()
	return None
}
