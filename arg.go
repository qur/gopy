// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"math"
	"unsafe"
)

func packValues(values []interface{}) ([]unsafe.Pointer, error) {
	cValues := make([]unsafe.Pointer, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case *string:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((*C.char)(nil)))))
		case *Object:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((*C.PyObject)(nil)))))
		case *int:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.int)(0)))))
		case *int8:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.int8_t)(0)))))
		case *int16:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.int16_t)(0)))))
		case *int32:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.int32_t)(0)))))
		case *int64:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.int64_t)(0)))))
		case *uint:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.uint)(0)))))
		case *uint8:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.uint8_t)(0)))))
		case *uint16:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.uint16_t)(0)))))
		case *uint32:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.uint32_t)(0)))))
		case *uint64:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.uint64_t)(0)))))
		case *float32:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.float)(0)))))
		case *float64:
			cValues[i] = unsafe.Pointer(C.calloc(1, C.size_t(unsafe.Sizeof((C.double)(0)))))
		default:
			return nil, TypeError.Err("Unsupported type: %T", v)
		}
	}
	return cValues, nil
}

func unpackValues(cValues []unsafe.Pointer, values []interface{}) error {
	defer func() {
		for _, v := range cValues {
			C.free(v)
		}
	}()
	for i, value := range values {
		switch v := value.(type) {
		case *string:
			*v = C.GoString(*(**C.char)(cValues[i]))
		case *Object:
			*v = newObject(*(**C.PyObject)(cValues[i]))
		case *int:
			*v = int(*(*C.int)(cValues[i]))
		case *int8:
			*v = int8(*(*C.int8_t)(cValues[i]))
		case *int16:
			*v = int16(*(*C.int16_t)(cValues[i]))
		case *int32:
			*v = int32(*(*C.int32_t)(cValues[i]))
		case *int64:
			*v = int64(*(*C.int64_t)(cValues[i]))
		case *uint:
			*v = uint(*(*C.uint)(cValues[i]))
		case *uint8:
			*v = uint8(*(*C.uint8_t)(cValues[i]))
		case *uint16:
			*v = uint16(*(*C.uint16_t)(cValues[i]))
		case *uint32:
			*v = uint32(*(*C.uint32_t)(cValues[i]))
		case *uint64:
			*v = uint64(*(*C.uint64_t)(cValues[i]))
		case *float32:
			*v = float32(*(*C.float)(cValues[i]))
		case *float64:
			*v = float64(*(*C.double)(cValues[i]))
		default:
			return TypeError.Err("Unsupported type: %T", v)
		}
	}
	return nil
}

func ParseTuple(args *Tuple, format string, values ...interface{}) error {
	if args == nil {
		return AssertionError.Err("ParseTuple: args was nil")
	}

	cv := (*unsafe.Pointer)(nil)
	cValues, err := packValues(values)
	if err != nil {
		return nil
	}

	if len(cValues) > 0 {
		cv = &cValues[0]
	}

	f := C.CString(format)
	defer C.free(unsafe.Pointer(f))

	ret := C.doParseTuple(c(args), f, cv, C.int(len(cValues)))
	if ret == 0 {
		return exception()
	}

	return unpackValues(cValues, values)
}

func ParseTupleAndKeywords(args *Tuple, kw *Dict, format string, kwlist []string, values ...interface{}) error {
	if args == nil {
		return AssertionError.Err("ParseTupleAndKeywords: args was nil")
	}

	cv := (*unsafe.Pointer)(nil)
	cValues, err := packValues(values)
	if err != nil {
		return nil
	}

	if len(cValues) > 0 {
		cv = &cValues[0]
	}

	f := C.CString(format)
	defer C.free(unsafe.Pointer(f))

	klist := make([]*C.char, len(kwlist)+1)

	for i, k := range kwlist {
		klist[i] = C.CString(k)
		defer C.free(unsafe.Pointer(klist[i]))
	}

	ret := C.doParseTupleKwds(c(args), c(kw), f, &klist[0], cv, C.int(len(cValues)))
	if ret == 0 {
		return exception()
	}

	return unpackValues(cValues, values)
}

func BuildValue(format string, values ...interface{}) (Object, error) {
	cValues := make([]C.ArgValue, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case string:
			s := C.CString(v)
			defer C.free(unsafe.Pointer(s))
			cValues[i]._type = &C.ffi_type_pointer
			cValues[i].value = unsafe.Pointer(&s)
		case Object:
			b := v.Base()
			cValues[i]._type = &C.ffi_type_pointer
			cValues[i].value = unsafe.Pointer(&b)
		case int:
			iv := C.int(v)
			if math.MaxInt == math.MaxInt64 {
				cValues[i]._type = &C.ffi_type_sint64
			} else {
				cValues[i]._type = &C.ffi_type_sint32
			}
			cValues[i].value = unsafe.Pointer(&iv)
		case int8:
			iv := C.int8_t(v)
			cValues[i]._type = &C.ffi_type_sint8
			cValues[i].value = unsafe.Pointer(&iv)
		case int16:
			iv := C.int16_t(v)
			cValues[i]._type = &C.ffi_type_sint16
			cValues[i].value = unsafe.Pointer(&iv)
		case int32:
			iv := C.int32_t(v)
			cValues[i]._type = &C.ffi_type_sint32
			cValues[i].value = unsafe.Pointer(&iv)
		case int64:
			iv := C.int64_t(v)
			cValues[i]._type = &C.ffi_type_sint64
			cValues[i].value = unsafe.Pointer(&iv)
		case uint:
			iv := C.uint(v)
			if math.MaxUint == math.MaxUint64 {
				cValues[i]._type = &C.ffi_type_uint64
			} else {
				cValues[i]._type = &C.ffi_type_uint32
			}
			cValues[i].value = unsafe.Pointer(&iv)
		case uint8:
			iv := C.uint8_t(v)
			cValues[i]._type = &C.ffi_type_uint8
			cValues[i].value = unsafe.Pointer(&iv)
		case uint16:
			iv := C.uint16_t(v)
			cValues[i]._type = &C.ffi_type_uint16
			cValues[i].value = unsafe.Pointer(&iv)
		case uint32:
			iv := C.uint32_t(v)
			cValues[i]._type = &C.ffi_type_uint32
			cValues[i].value = unsafe.Pointer(&iv)
		case uint64:
			iv := C.uint64_t(v)
			cValues[i]._type = &C.ffi_type_uint64
			cValues[i].value = unsafe.Pointer(&iv)
		case float32:
			fv := C.float(v)
			cValues[i]._type = &C.ffi_type_float
			cValues[i].value = unsafe.Pointer(&fv)
		case float64:
			fv := C.double(v)
			cValues[i]._type = &C.ffi_type_double
			cValues[i].value = unsafe.Pointer(&fv)
		default:
			return nil, TypeError.Err("Unsupported type: %T", v)
		}
	}
	f := C.CString(format)
	defer C.free(unsafe.Pointer(f))
	ret := C.doBuildValue(f, &cValues[0], C.int(len(cValues)))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret), nil
}
