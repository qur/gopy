package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"unsafe"
)

type Bytes struct {
	AbstractObject
	o C.PyBytesObject
}

// BytesType is the Type object that represents the Bytes type.
var BytesType = (*Type)(unsafe.Pointer(&C.PyBytes_Type))

func bytesCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.bytesCheck(c(obj)) != 0
}

func newBytes(obj *C.PyObject) *Bytes {
	return (*Bytes)(unsafe.Pointer(obj))
}

func NewBytes(b []byte) *Bytes {
	d := (*C.char)(unsafe.Pointer(&b[0]))
	return newBytes(C.PyBytes_FromStringAndSize(d, C.Py_ssize_t(len(b))))
}

// Bytes returns the contents of the Bytes object
//
// NOTE: This function returns a reference to the internal buffer of the Bytes
// object, and MUST NOT be changed.
func (b *Bytes) Bytes() []byte {
	// we use PyBytes_AsString because PyBytes_AsStringAndSize throws a hissy
	// fit about embedded NULs, and we don't care, that's fine in []byte.
	data := C.PyBytes_AsString(c(b))
	if data == nil {
		// this should only happen if b fails PyBytes_Check, which shouldn't be
		// possible.
		return nil
	}
	length := C.PyBytes_Size(c(b))
	if length < 0 {
		// Again, this should only happen if b fails PyBytes_Check, which
		// shouldn't be possible.
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(data)), int(length))
}

func (b *Bytes) Size() int {
	ret := C.PyBytes_Size(c(b))
	if ret < 0 {
		// this should only happen if b fails PyBytes_Check, which shouldn't be
		// possible.
		return 0
	}
	return int(ret)
}

func (b *Bytes) String() string {
	if b == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", b.Bytes())
}
