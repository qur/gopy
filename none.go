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

// ReturnNone returns None.
// This is a deprecated convenience function for returning None from a function.
// Previously the reference count for None had to be incremented, but now this
// function adds no value, and should not be used.
//
// Return value: New Reference.
func ReturnNone() Object {
	return None
}
