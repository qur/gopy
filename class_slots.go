// Code generated by gen_slots.py. DO NOT EDIT.

// This file is automatically generated.  To regenerate:
//   go generate ./...

package py

// #include "utils.h"
import "C"

import (
	"reflect"
	"unsafe"
)

// ===============================================================

type tp_repr interface {
	PyRepr() string
}
type tp_hash interface {
	PyHash() (uint32, error)
}
type tp_call interface {
	PyCall(*Tuple, *Dict) (Object, error)
}
type tp_str interface {
	PyStr() string
}
type tp_getattro interface {
	PyGetAttr(Object) (Object, error)
}
type tp_setattro interface {
	PySetAttr(Object, Object) error
}
type tp_iter interface {
	PyIter() (Object, error)
}
type tp_iternext interface {
	PyIterNext() (Object, error)
}
type tp_descr_get interface {
	PyDescrGet(Object, Object) (Object, error)
}
type tp_descr_set interface {
	PyDescrSet(Object, Object) error
}
type tp_init interface {
	PyInit(*Tuple, *Dict) error
}
type am_await interface {
	PyAwait() (Object, error)
}
type am_aiter interface {
	PyAsyncIter() (Object, error)
}
type am_anext interface {
	PyAsyncNext() (Object, error)
}
type nb_add interface {
	PyAdd(Object) (Object, error)
}
type nb_inplace_add interface {
	PyInplaceAdd(Object) (Object, error)
}
type nb_subtract interface {
	PySubtract(Object) (Object, error)
}
type nb_inplace_subtract interface {
	PyInplaceSubtract(Object) (Object, error)
}
type nb_multiply interface {
	PyMultiply(Object) (Object, error)
}
type nb_inplace_multiply interface {
	PyInplaceMultiply(Object) (Object, error)
}
type nb_remainder interface {
	PyRemainder(Object) (Object, error)
}
type nb_inplace_remainder interface {
	PyInplaceRemainder(Object) (Object, error)
}
type nb_divmod interface {
	PyDivmod(Object) (Object, error)
}
type nb_power interface {
	PyPower(Object, Object) (Object, error)
}
type nb_inplace_power interface {
	PyInplacePower(Object, Object) (Object, error)
}
type nb_negative interface {
	PyNegative() (Object, error)
}
type nb_positive interface {
	PyPositive() (Object, error)
}
type nb_absolute interface {
	PyAbsolute() (Object, error)
}
type nb_bool interface {
	PyBool() (bool, error)
}
type nb_invert interface {
	PyInvert() (Object, error)
}
type nb_lshift interface {
	PyLShift(Object) (Object, error)
}
type nb_inplace_lshift interface {
	PyInplaceLShift(Object) (Object, error)
}
type nb_rshift interface {
	PyRShift(Object) (Object, error)
}
type nb_inplace_rshift interface {
	PyInplaceRShift(Object) (Object, error)
}
type nb_and interface {
	PyAnd(Object) (Object, error)
}
type nb_inplace_and interface {
	PyInplaceAnd(Object) (Object, error)
}
type nb_xor interface {
	PyXor(Object) (Object, error)
}
type nb_inplace_xor interface {
	PyInplaceXor(Object) (Object, error)
}
type nb_or interface {
	PyOr(Object) (Object, error)
}
type nb_inplace_or interface {
	PyInplaceOr(Object) (Object, error)
}
type nb_int interface {
	PyInt() (Object, error)
}
type nb_float interface {
	PyFloat() (Object, error)
}
type nb_floor_divide interface {
	PyFloorDivide(Object) (Object, error)
}
type nb_inplace_floor_divide interface {
	PyInplaceFloorDivide(Object) (Object, error)
}
type nb_true_divide interface {
	PyTrueDivide(Object) (Object, error)
}
type nb_inplace_true_divide interface {
	PyInplaceTrueDivide(Object) (Object, error)
}
type nb_index interface {
	PyIndex() (Object, error)
}
type nb_matrix_multiply interface {
	PyMatrixMultiply(Object) (Object, error)
}
type nb_inplace_matrix_multiply interface {
	PyInplaceMatrixMultiply(Object) (Object, error)
}
type mp_length interface {
	PyMappingLen() int
}
type mp_subscript interface {
	PySubscript(Object) (Object, error)
}
type mp_ass_subscript interface {
	PyAssSubscript(Object, Object) error
}
type sq_length interface {
	PyLen() int
}
type sq_concat interface {
	PyConcat(Object) (Object, error)
}
type sq_repeat interface {
	PyRepeat(Object, int) (Object, error)
}
type sq_item interface {
	PyItem(Object, int) (Object, error)
}
type sq_ass_item interface {
	PyAssItem(Object, int, Object) error
}
type sq_contains interface {
	PyContains(Object) (bool, error)
}
type sq_inplace_concat interface {
	PyInplaceConcat(Object) (Object, error)
}
type sq_inplace_repeat interface {
	PyInplaceRepeat(Object, int) (Object, error)
}
type bf_getbuffer interface {
	PyGetBuffer(Object, int) error
}
type bf_releasebuffer interface {
	PyReleaseBuffer(Object)
}

// ===============================================================
// ===============================================================

//export goClassSlot_tp_repr
func goClassSlot_tp_repr(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(tp_repr)

	s := C.CString(co.PyRepr())
	defer C.free(unsafe.Pointer(s))

	return unsafe.Pointer(C.PyUnicode_FromString(s))
}
//export goClassSlot_tp_hash
func goClassSlot_tp_hash(obj unsafe.Pointer) C.long {
	co := newObject((*C.PyObject)(obj)).(tp_hash)

	ret, err := co.PyHash()
	if err != nil {
		raise(err)
		return -1
	} else if C.long(ret) == -1 {
		return -2
	}

	return C.long(ret)
}
//export goClassSlot_tp_call
func goClassSlot_tp_call(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(tp_call)

	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))
	ret, err := co.PyCall(a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_tp_str
func goClassSlot_tp_str(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(tp_str)

	s := C.CString(co.PyStr())
	defer C.free(unsafe.Pointer(s))

	return unsafe.Pointer(C.PyUnicode_FromString(s))
}
//export goClassSlot_tp_getattro
func goClassSlot_tp_getattro(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(tp_getattro)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyGetAttr(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_tp_setattro
func goClassSlot_tp_setattro(obj, arg1, arg2 unsafe.Pointer) int {
	co := newObject((*C.PyObject)(obj)).(tp_setattro)

	o1 := newObject((*C.PyObject)(arg1))
	o2 := newObject((*C.PyObject)(arg2))
	if err := co.PySetAttr(o1, o2); err != nil {
		raise(err)
		return -1
	}

	return 0
}
//export goClassSlot_tp_iter
func goClassSlot_tp_iter(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(tp_iter)

	ret, err := co.PyIter()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_tp_iternext
func goClassSlot_tp_iternext(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(tp_iternext)

	ret, err := co.PyIterNext()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_tp_descr_get
func goClassSlot_tp_descr_get(obj, arg1, arg2 unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(tp_descr_get)

	o1 := newObject((*C.PyObject)(arg1))
	o2 := newObject((*C.PyObject)(arg2))
	ret, err := co.PyDescrGet(o1, o2)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_tp_descr_set
func goClassSlot_tp_descr_set(obj, arg1, arg2 unsafe.Pointer) int {
	co := newObject((*C.PyObject)(obj)).(tp_descr_set)

	o1 := newObject((*C.PyObject)(arg1))
	o2 := newObject((*C.PyObject)(arg2))
	if err := co.PyDescrSet(o1, o2); err != nil {
		raise(err)
		return -1
	}

	return 0
}
//export goClassSlot_tp_init
func goClassSlot_tp_init(obj, args, kwds unsafe.Pointer) int {
	co := newObject((*C.PyObject)(obj)).(tp_init)

	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))
	if err := co.PyInit(a, k); err != nil {
		raise(err)
		return -1
	}

	return 0
}
//export goClassSlot_am_await
func goClassSlot_am_await(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(am_await)

	ret, err := co.PyAwait()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_am_aiter
func goClassSlot_am_aiter(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(am_aiter)

	ret, err := co.PyAsyncIter()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_am_anext
func goClassSlot_am_anext(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(am_anext)

	ret, err := co.PyAsyncNext()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_add
func goClassSlot_nb_add(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_add)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyAdd(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_add
func goClassSlot_nb_inplace_add(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_add)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceAdd(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_subtract
func goClassSlot_nb_subtract(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_subtract)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PySubtract(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_subtract
func goClassSlot_nb_inplace_subtract(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_subtract)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceSubtract(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_multiply
func goClassSlot_nb_multiply(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_multiply)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyMultiply(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_multiply
func goClassSlot_nb_inplace_multiply(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_multiply)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceMultiply(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_remainder
func goClassSlot_nb_remainder(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_remainder)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyRemainder(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_remainder
func goClassSlot_nb_inplace_remainder(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_remainder)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceRemainder(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_divmod
func goClassSlot_nb_divmod(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_divmod)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyDivmod(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_power
func goClassSlot_nb_power(obj, arg1, arg2 unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_power)

	o1 := newObject((*C.PyObject)(arg1))
	o2 := newObject((*C.PyObject)(arg2))
	ret, err := co.PyPower(o1, o2)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_power
func goClassSlot_nb_inplace_power(obj, arg1, arg2 unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_power)

	o1 := newObject((*C.PyObject)(arg1))
	o2 := newObject((*C.PyObject)(arg2))
	ret, err := co.PyInplacePower(o1, o2)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_negative
func goClassSlot_nb_negative(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_negative)

	ret, err := co.PyNegative()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_positive
func goClassSlot_nb_positive(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_positive)

	ret, err := co.PyPositive()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_absolute
func goClassSlot_nb_absolute(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_absolute)

	ret, err := co.PyAbsolute()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_bool
func goClassSlot_nb_bool(obj unsafe.Pointer) int {
	co := newObject((*C.PyObject)(obj)).(nb_bool)

	ret, err := co.PyBool()
	if err != nil {
		raise(err)
		return -1
	}

	if ret {
		return 1
	}

	return 0
}
//export goClassSlot_nb_invert
func goClassSlot_nb_invert(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_invert)

	ret, err := co.PyInvert()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_lshift
func goClassSlot_nb_lshift(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_lshift)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyLShift(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_lshift
func goClassSlot_nb_inplace_lshift(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_lshift)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceLShift(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_rshift
func goClassSlot_nb_rshift(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_rshift)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyRShift(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_rshift
func goClassSlot_nb_inplace_rshift(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_rshift)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceRShift(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_and
func goClassSlot_nb_and(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_and)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyAnd(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_and
func goClassSlot_nb_inplace_and(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_and)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceAnd(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_xor
func goClassSlot_nb_xor(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_xor)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyXor(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_xor
func goClassSlot_nb_inplace_xor(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_xor)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceXor(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_or
func goClassSlot_nb_or(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_or)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyOr(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_or
func goClassSlot_nb_inplace_or(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_or)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceOr(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_int
func goClassSlot_nb_int(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_int)

	ret, err := co.PyInt()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_float
func goClassSlot_nb_float(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_float)

	ret, err := co.PyFloat()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_floor_divide
func goClassSlot_nb_floor_divide(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_floor_divide)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyFloorDivide(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_floor_divide
func goClassSlot_nb_inplace_floor_divide(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_floor_divide)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceFloorDivide(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_true_divide
func goClassSlot_nb_true_divide(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_true_divide)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyTrueDivide(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_true_divide
func goClassSlot_nb_inplace_true_divide(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_true_divide)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceTrueDivide(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_index
func goClassSlot_nb_index(obj unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_index)

	ret, err := co.PyIndex()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_matrix_multiply
func goClassSlot_nb_matrix_multiply(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_matrix_multiply)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyMatrixMultiply(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_nb_inplace_matrix_multiply
func goClassSlot_nb_inplace_matrix_multiply(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(nb_inplace_matrix_multiply)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceMatrixMultiply(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_mp_length
func goClassSlot_mp_length(obj unsafe.Pointer) C.Py_ssize_t {
	co := newObject((*C.PyObject)(obj)).(mp_length)

	return C.Py_ssize_t(co.PyMappingLen())
}
//export goClassSlot_mp_subscript
func goClassSlot_mp_subscript(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(mp_subscript)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PySubscript(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_mp_ass_subscript
func goClassSlot_mp_ass_subscript(obj, arg1, arg2 unsafe.Pointer) int {
	co := newObject((*C.PyObject)(obj)).(mp_ass_subscript)

	o1 := newObject((*C.PyObject)(arg1))
	o2 := newObject((*C.PyObject)(arg2))
	if err := co.PyAssSubscript(o1, o2); err != nil {
		raise(err)
		return -1
	}

	return 0
}
//export goClassSlot_sq_length
func goClassSlot_sq_length(obj unsafe.Pointer) C.Py_ssize_t {
	co := newObject((*C.PyObject)(obj)).(sq_length)

	return C.Py_ssize_t(co.PyLen())
}
//export goClassSlot_sq_concat
func goClassSlot_sq_concat(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(sq_concat)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyConcat(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_sq_repeat
func goClassSlot_sq_repeat(obj, arg1 unsafe.Pointer, arg2 C.Py_ssize_t) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(sq_repeat)

	o := newObject((*C.PyObject)(arg1))
	ret, err := co.PyRepeat(o, int(arg2))
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_sq_item
func goClassSlot_sq_item(obj, arg1 unsafe.Pointer, arg2 C.Py_ssize_t) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(sq_item)

	o := newObject((*C.PyObject)(arg1))
	ret, err := co.PyItem(o, int(arg2))
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_sq_ass_item
func goClassSlot_sq_ass_item(obj, arg1 unsafe.Pointer, arg2 C.Py_ssize_t, arg3 unsafe.Pointer) C.int {
	co := newObject((*C.PyObject)(obj)).(sq_ass_item)

	o1 := newObject((*C.PyObject)(arg1))
	o3 := newObject((*C.PyObject)(arg1))
	if err := co.PyAssItem(o1, int(arg2), o3); err != nil {
		raise(err)
		return -1
	}

	return 0
}
//export goClassSlot_sq_contains
func goClassSlot_sq_contains(obj, arg unsafe.Pointer) int {
	co := newObject((*C.PyObject)(obj)).(sq_contains)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyContains(o)
	if err != nil {
		raise(err)
		return -1
	}

	if ret {
		return 1
	}

	return 0
}
//export goClassSlot_sq_inplace_concat
func goClassSlot_sq_inplace_concat(obj, arg unsafe.Pointer) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(sq_inplace_concat)

	o := newObject((*C.PyObject)(arg))
	ret, err := co.PyInplaceConcat(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_sq_inplace_repeat
func goClassSlot_sq_inplace_repeat(obj, arg1 unsafe.Pointer, arg2 C.Py_ssize_t) unsafe.Pointer {
	co := newObject((*C.PyObject)(obj)).(sq_inplace_repeat)

	o := newObject((*C.PyObject)(arg1))
	ret, err := co.PyInplaceRepeat(o, int(arg2))
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}
//export goClassSlot_bf_getbuffer
func goClassSlot_bf_getbuffer(obj, arg1 unsafe.Pointer, arg2 C.int) int {
	co := newObject((*C.PyObject)(obj)).(bf_getbuffer)

	o := newObject((*C.PyObject)(arg1))
	if err := co.PyGetBuffer(o, int(arg2)); err != nil {
		raise(err)
		return -1
	}

	return 0
}
//export goClassSlot_bf_releasebuffer
func goClassSlot_bf_releasebuffer(obj, arg unsafe.Pointer) {
	co := newObject((*C.PyObject)(obj)).(bf_releasebuffer)

	o := newObject((*C.PyObject)(arg))
	co.PyReleaseBuffer(o)
}

// ===============================================================
// ===============================================================

var slotMap = map[C.uint64_t]reflect.Type{
	C.CLASS_HAS_TP_REPR: reflect.TypeOf((*tp_repr)(nil)).Elem(),
	C.CLASS_HAS_TP_HASH: reflect.TypeOf((*tp_hash)(nil)).Elem(),
	C.CLASS_HAS_TP_CALL: reflect.TypeOf((*tp_call)(nil)).Elem(),
	C.CLASS_HAS_TP_STR: reflect.TypeOf((*tp_str)(nil)).Elem(),
	C.CLASS_HAS_TP_GETATTRO: reflect.TypeOf((*tp_getattro)(nil)).Elem(),
	C.CLASS_HAS_TP_SETATTRO: reflect.TypeOf((*tp_setattro)(nil)).Elem(),
	C.CLASS_HAS_TP_ITER: reflect.TypeOf((*tp_iter)(nil)).Elem(),
	C.CLASS_HAS_TP_ITERNEXT: reflect.TypeOf((*tp_iternext)(nil)).Elem(),
	C.CLASS_HAS_TP_DESCR_GET: reflect.TypeOf((*tp_descr_get)(nil)).Elem(),
	C.CLASS_HAS_TP_DESCR_SET: reflect.TypeOf((*tp_descr_set)(nil)).Elem(),
	C.CLASS_HAS_TP_INIT: reflect.TypeOf((*tp_init)(nil)).Elem(),
	C.CLASS_HAS_AM_AWAIT: reflect.TypeOf((*am_await)(nil)).Elem(),
	C.CLASS_HAS_AM_AITER: reflect.TypeOf((*am_aiter)(nil)).Elem(),
	C.CLASS_HAS_AM_ANEXT: reflect.TypeOf((*am_anext)(nil)).Elem(),
	C.CLASS_HAS_NB_ADD: reflect.TypeOf((*nb_add)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_ADD: reflect.TypeOf((*nb_inplace_add)(nil)).Elem(),
	C.CLASS_HAS_NB_SUBTRACT: reflect.TypeOf((*nb_subtract)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_SUBTRACT: reflect.TypeOf((*nb_inplace_subtract)(nil)).Elem(),
	C.CLASS_HAS_NB_MULTIPLY: reflect.TypeOf((*nb_multiply)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_MULTIPLY: reflect.TypeOf((*nb_inplace_multiply)(nil)).Elem(),
	C.CLASS_HAS_NB_REMAINDER: reflect.TypeOf((*nb_remainder)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_REMAINDER: reflect.TypeOf((*nb_inplace_remainder)(nil)).Elem(),
	C.CLASS_HAS_NB_DIVMOD: reflect.TypeOf((*nb_divmod)(nil)).Elem(),
	C.CLASS_HAS_NB_POWER: reflect.TypeOf((*nb_power)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_POWER: reflect.TypeOf((*nb_inplace_power)(nil)).Elem(),
	C.CLASS_HAS_NB_NEGATIVE: reflect.TypeOf((*nb_negative)(nil)).Elem(),
	C.CLASS_HAS_NB_POSITIVE: reflect.TypeOf((*nb_positive)(nil)).Elem(),
	C.CLASS_HAS_NB_ABSOLUTE: reflect.TypeOf((*nb_absolute)(nil)).Elem(),
	C.CLASS_HAS_NB_BOOL: reflect.TypeOf((*nb_bool)(nil)).Elem(),
	C.CLASS_HAS_NB_INVERT: reflect.TypeOf((*nb_invert)(nil)).Elem(),
	C.CLASS_HAS_NB_LSHIFT: reflect.TypeOf((*nb_lshift)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_LSHIFT: reflect.TypeOf((*nb_inplace_lshift)(nil)).Elem(),
	C.CLASS_HAS_NB_RSHIFT: reflect.TypeOf((*nb_rshift)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_RSHIFT: reflect.TypeOf((*nb_inplace_rshift)(nil)).Elem(),
	C.CLASS_HAS_NB_AND: reflect.TypeOf((*nb_and)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_AND: reflect.TypeOf((*nb_inplace_and)(nil)).Elem(),
	C.CLASS_HAS_NB_XOR: reflect.TypeOf((*nb_xor)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_XOR: reflect.TypeOf((*nb_inplace_xor)(nil)).Elem(),
	C.CLASS_HAS_NB_OR: reflect.TypeOf((*nb_or)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_OR: reflect.TypeOf((*nb_inplace_or)(nil)).Elem(),
	C.CLASS_HAS_NB_INT: reflect.TypeOf((*nb_int)(nil)).Elem(),
	C.CLASS_HAS_NB_FLOAT: reflect.TypeOf((*nb_float)(nil)).Elem(),
	C.CLASS_HAS_NB_FLOOR_DIVIDE: reflect.TypeOf((*nb_floor_divide)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_FLOOR_DIVIDE: reflect.TypeOf((*nb_inplace_floor_divide)(nil)).Elem(),
	C.CLASS_HAS_NB_TRUE_DIVIDE: reflect.TypeOf((*nb_true_divide)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_TRUE_DIVIDE: reflect.TypeOf((*nb_inplace_true_divide)(nil)).Elem(),
	C.CLASS_HAS_NB_INDEX: reflect.TypeOf((*nb_index)(nil)).Elem(),
	C.CLASS_HAS_NB_MATRIX_MULTIPLY: reflect.TypeOf((*nb_matrix_multiply)(nil)).Elem(),
	C.CLASS_HAS_NB_INPLACE_MATRIX_MULTIPLY: reflect.TypeOf((*nb_inplace_matrix_multiply)(nil)).Elem(),
	C.CLASS_HAS_MP_LENGTH: reflect.TypeOf((*mp_length)(nil)).Elem(),
	C.CLASS_HAS_MP_SUBSCRIPT: reflect.TypeOf((*mp_subscript)(nil)).Elem(),
	C.CLASS_HAS_MP_ASS_SUBSCRIPT: reflect.TypeOf((*mp_ass_subscript)(nil)).Elem(),
	C.CLASS_HAS_SQ_LENGTH: reflect.TypeOf((*sq_length)(nil)).Elem(),
	C.CLASS_HAS_SQ_CONCAT: reflect.TypeOf((*sq_concat)(nil)).Elem(),
	C.CLASS_HAS_SQ_REPEAT: reflect.TypeOf((*sq_repeat)(nil)).Elem(),
	C.CLASS_HAS_SQ_ITEM: reflect.TypeOf((*sq_item)(nil)).Elem(),
	C.CLASS_HAS_SQ_ASS_ITEM: reflect.TypeOf((*sq_ass_item)(nil)).Elem(),
	C.CLASS_HAS_SQ_CONTAINS: reflect.TypeOf((*sq_contains)(nil)).Elem(),
	C.CLASS_HAS_SQ_INPLACE_CONCAT: reflect.TypeOf((*sq_inplace_concat)(nil)).Elem(),
	C.CLASS_HAS_SQ_INPLACE_REPEAT: reflect.TypeOf((*sq_inplace_repeat)(nil)).Elem(),
	C.CLASS_HAS_BF_GETBUFFER: reflect.TypeOf((*bf_getbuffer)(nil)).Elem(),
	C.CLASS_HAS_BF_RELEASEBUFFER: reflect.TypeOf((*bf_releasebuffer)(nil)).Elem(),
}

// ===============================================================
