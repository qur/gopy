// Code generated by gen_extension. DO NOT EDIT.

package main

// #include <python3.11/Python.h>
import "C"

import (
	"unsafe"

	"qur.me/py/v3"
)

//export PyInit_extension
func PyInit_extension() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(py.InitExtension(initExtension)))
}