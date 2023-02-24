package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// NewByteArray returns a new ByteArray instance that contains a copy of the
// supplied []byte.
//
// Return value: New Reference.
func NewByteArray(b []byte) *ByteArray {
	d := (*C.char)(unsafe.Pointer(&b[0]))
	return newByteArray(C.PyByteArray_FromStringAndSize(d, C.Py_ssize_t(len(b))))
}

// NewByteArrayFromObject returns a new ByteArray instance created from the
// supplied Object.
//
// Return value: New Reference.
func NewByteArrayFromObject(o Object) (*ByteArray, error) {
	ret := C.PyByteArray_FromObject(c(o))
	if ret == nil {
		return nil, exception()
	}
	return newByteArray(ret), nil
}

// Bytes returns the contents of the ByteArray object
//
// NOTE: This function returns a reference to the internal buffer of the
// ByteArray object, so append or other operations that change the length or
// capacity MUST NOT be used.
func (b *ByteArray) Bytes() []byte {
	data := C.PyByteArray_AsString(c(b))
	if data == nil {
		// this should only happen if b fails PyByteArray_Check, which shouldn't
		// be possible.
		return nil
	}
	length := C.PyByteArray_Size(c(b))
	if length < 0 {
		// Again, this should only happen if b fails PyByteArray_Check, which
		// shouldn't be possible.
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(data)), int(length))
}

// Resize changes the size of b's internal buffer to the given size.
func (b *ByteArray) Resize(s int) error {
	ret := C.PyByteArray_Resize(c(b), C.Py_ssize_t(s))
	return int2Err(ret)
}

// String returns a string representation of the ByteArray
func (b *ByteArray) String() string {
	if b == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", b.Bytes())
}
