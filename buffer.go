package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

type BufferProtocol interface {
	Object
	AsBufferMethods() *BufferMethods
}

// AsBuffer returns a BufferMethods instance that refers to the same
// underlying Python object as obj. If obj doesn't implement the "Buffer
// Protocol", then nil is returned.
//
// This method is more complete than the BufferProtocol interface, as it will
// also work with unknown or dynamic types that implement the "Buffer
// Protocol".
func AsBufferMethods(obj Object) *BufferMethods {
	if n, ok := obj.(BufferProtocol); ok {
		return n.AsBufferMethods()
	}
	if C.PyObject_CheckBuffer(c(obj)) > 0 {
		return (*BufferMethods)(unsafe.Pointer(obj.Base()))
	}
	return nil
}

func (b *BufferMethods) GetBuffer(flags BufferFlags) (*Buffer, error) {
	buf := newBuffer()
	ret := C.PyObject_GetBuffer(c(b), buf.c(), C.int(flags))
	if ret < 0 {
		return nil, exception()
	}
	return buf, nil
}

func GetBuffer(obj Object, flags BufferFlags) (*Buffer, error) {
	bm := AsBufferMethods(obj)
	if bm == nil {
		return nil, TypeError.Err("%s does implement Buffer Protocol", obj.Type())
	}
	return bm.GetBuffer(flags)
}

type BufferFlags int

// TODO(jp3): buffer flags

type BufferOrder byte

const (
	BufferOrderC       BufferOrder = 'C'
	BufferOrderFortran BufferOrder = 'F'
	BufferOrderAny     BufferOrder = 'A'
)

type Buffer struct {
	buf C.Py_buffer
}

func newBuffer() *Buffer {
	return &Buffer{}
}

func (b *Buffer) c() *C.Py_buffer {
	return &b.buf
}

func (b *Buffer) Release() {
	C.PyBuffer_Release(&b.buf)
}

func (b *Buffer) IsContiguous(order BufferOrder) bool {
	ret := C.PyBuffer_IsContiguous(b.c(), C.char(order))
	return ret == 1
}

func (b *Buffer) GetPointer(indicies ...int) (*byte, error) {
	if len(indicies) != int(b.buf.ndim) {
		return nil, ValueError.Err("wrong number of indicies: %d (wanted %d)", len(indicies), b.buf.ndim)
	}
	ind := make([]C.Py_ssize_t, len(indicies))
	for i, index := range indicies {
		ind[i] = C.Py_ssize_t(index)
	}
	ret := C.PyBuffer_GetPointer(b.c(), &ind[0])
	return (*byte)(ret), nil
}

func (b *Buffer) FromContinguous(buf []byte, order BufferOrder) error {
	ret := C.PyBuffer_FromContiguous(b.c(), unsafe.Pointer(&buf[0]), C.Py_ssize_t(len(buf)), C.char(order))
	return int2Err(ret)
}

func (b *Buffer) ToContinguous(buf []byte, order BufferOrder) error {
	ret := C.PyBuffer_ToContiguous(unsafe.Pointer(&buf[0]), b.c(), C.Py_ssize_t(len(buf)), C.char(order))
	return int2Err(ret)
}

// TODO(jp3): fill contiguous strides

func (b *Buffer) FillInfo(exporter Object, buf []byte, readonly bool, flags BufferFlags) error {
	ro := C.int(0)
	if readonly {
		ro = C.int(1)
	}
	ret := C.PyBuffer_FillInfo(b.c(), c(exporter), unsafe.Pointer(&buf[0]), C.Py_ssize_t(len(buf)), ro, C.int(flags))
	return int2Err(ret)
}
