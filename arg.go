// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "unsafe"

func packValues(values []interface{}) ([]unsafe.Pointer, error) {
	cValues := make([]unsafe.Pointer, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case *string:
			cValues[i] = unsafe.Pointer(new(*C.char))
		case *Object:
			cValues[i] = unsafe.Pointer(new(*C.PyObject))
		case *int:
			cValues[i] = unsafe.Pointer(new(C.int))
		case *int8:
			cValues[i] = unsafe.Pointer(new(C.int8_t))
		case *int16:
			cValues[i] = unsafe.Pointer(new(C.int16_t))
		case *int32:
			cValues[i] = unsafe.Pointer(new(C.int32_t))
		case *int64:
			cValues[i] = unsafe.Pointer(new(C.int64_t))
		case *uint:
			cValues[i] = unsafe.Pointer(new(C.uint))
		case *uint8:
			cValues[i] = unsafe.Pointer(new(C.uint8_t))
		case *uint16:
			cValues[i] = unsafe.Pointer(new(C.uint16_t))
		case *uint32:
			cValues[i] = unsafe.Pointer(new(C.uint32_t))
		case *uint64:
			cValues[i] = unsafe.Pointer(new(C.uint64_t))
		case *float32:
			cValues[i] = unsafe.Pointer(new(C.float))
		case *float64:
			cValues[i] = unsafe.Pointer(new(C.double))
		default:
			return nil, TypeError.Err("Unsupported type: %T", v)
		}
	}
	return cValues, nil
}

func unpackValues(cValues []unsafe.Pointer, values []interface{}) error {
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
