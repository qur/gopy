package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

func packValues(values []interface{}) ([]unsafe.Pointer, error) {
	cValues := make([]unsafe.Pointer, len(values))

	for i, value := range values {
		p := calloc(value)
		if p == nil {
			return nil, TypeError.Err("Unsupported type: %T", value)
		}

		cValues[i] = p
	}

	return cValues, nil
}

func ParseTuple(args *Tuple, format string, values ...interface{}) error {
	if args == nil {
		return AssertionError.Err("ParseTuple: args was nil")
	}

	cValues, err := packValues(values)
	if err != nil {
		return nil
	}

	cv := (*unsafe.Pointer)(nil)
	if len(cValues) > 0 {
		cv = &cValues[0]
	}

	f := C.CString(format)
	defer cfree(f)

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

	cValues, err := packValues(values)
	if err != nil {
		return nil
	}

	cv := (*unsafe.Pointer)(nil)
	if len(cValues) > 0 {
		cv = &cValues[0]
	}

	f := C.CString(format)
	defer cfree(f)

	klist := make([]*C.char, len(kwlist)+1)

	for i, k := range kwlist {
		klist[i] = C.CString(k)
		defer cfree(klist[i])
	}

	ret := C.doParseTupleKwds(c(args), c(kw), f, &klist[0], cv, C.int(len(cValues)))
	if ret == 0 {
		return exception()
	}

	return unpackValues(cValues, values)
}

func BuildValue(format string, values ...interface{}) (Object, error) {
	cValues, cleanup, err := buildCValues(values)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	f := C.CString(format)
	defer cfree(f)

	ret := C.doBuildValue(f, &cValues[0], C.int(len(cValues)))
	if ret == nil {
		return nil, exception()
	}

	return newObject(ret), nil
}
