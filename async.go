package py

// #include "utils.h"
import "C"
import "unsafe"

type SendResult int

const (
	SendReturn SendResult = C.PYGEN_RETURN
	SendError  SendResult = C.PYGEN_ERROR
	SendNext   SendResult = C.PYGEN_NEXT
)

type AsyncIterator interface {
	Object
	AsAsyncIterator() *AsyncIteratorMethods
}

type AsyncIteratorMethods struct {
	abstractObject
	o C.PyObject
}

func AsAsyncIterator(obj Object) *AsyncIteratorMethods {
	if n, ok := obj.(AsyncIterator); ok {
		return n.AsAsyncIterator()
	}
	if C.aIterCheck(c(obj)) > 0 {
		return (*AsyncIteratorMethods)(unsafe.Pointer(obj.Base()))
	}
	return nil
}
