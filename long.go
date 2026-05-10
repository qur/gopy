package py

// #include "utils.h"
import "C"

import (
	"strconv"
)

// NewLongChecked return a new Long with the given value, or an error if an
// exception was raised.
//
// Return value: New Reference.
func NewLongChecked[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64](i T) (*Long, error) {
	if i >= 0 {
		return newLong(C.PyLong_FromUInt64(C.uint64_t(uint64(i)))), exception()
	}

	return newLong(C.PyLong_FromInt64(C.int64_t(int64(i)))), exception()
}

// NewLong returns a new Long with the given value, it will panic if an
// exception is raised.
//
// Return value: New Reference.
func NewLong[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64](i T) *Long {
	l, err := NewLongChecked(i)
	if err != nil {
		panic(err)
	}

	return l
}

// AsInt64 returns the value of l as an int64. If the value does not fit in an
// int64, then OverflowError will be returned.
func (l *Long) AsInt64() (int64, error) {
	var v C.int64_t

	if err := int2Err(C.PyLong_AsInt64(c(l), &v)); err != nil {
		return 0, err
	}

	return int64(v), nil
}

// Int64 returns the value of l as an int64, it will panic if AsInt64 would
// return an error.
func (l *Long) Int64() int64 {
	i, err := l.AsInt64()
	if err != nil {
		panic(err)
	}

	return i
}

// AsUint64 returns the value of l as a uint64. If the value does not fit in an
// uint64, then OverflowError will be returned, if the value is negative, then
// ValueError will be returned.
func (l *Long) AsUint64() (uint64, error) {
	var v C.uint64_t

	if err := int2Err(C.PyLong_AsUInt64(c(l), &v)); err != nil {
		return 0, err
	}

	return uint64(v), nil
}

// Uint64 returns the value of l as a uint64, it will panic if AsUint64 would
// return an error.
func (l *Long) Uint64() uint64 {
	i, err := l.AsUint64()
	if err != nil {
		panic(err)
	}

	return i
}

// String returns a string representation of the value l. It will panic if l
// can't be converted to either an in64 or a uint64.
func (l *Long) String() string {
	if l == nil {
		return nilValue
	}

	i64, err := l.AsInt64()
	if err == nil {
		return strconv.FormatInt(i64, 10)
	}

	u64, err := l.AsUint64()
	if err == nil {
		return strconv.FormatUint(u64, 10)
	}

	panic(NewErrorV(ValueError, l))
}
