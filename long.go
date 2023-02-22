package py

// #include "utils.h"
import "C"

import (
	"fmt"
)

func NewLong[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32](i T) *Long {
	return newLong(C.PyLong_FromLongLong(C.longlong(int64(i))))
}

func (l *Long) Int64() int64 {
	return int64(C.PyLong_AsLongLong(c(l)))
}

func (l *Long) String() string {
	if l == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", l.Int64())
}
