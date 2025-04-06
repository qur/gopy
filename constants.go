package py

// #include "utils.h"
import "C"

// ConstantType defines the type of constants that can be returned by
// [GetConstant]/[BorrowConstant].
type ConstantType uint

// The types of constant that can be returned by
// [GetConstant]/[BrorrowConstant].
const (
	ConstantNone ConstantType = iota
	ConstantFalse
	ConstantTrue
	ConstantEllipsis
	ConstantNotImplemented
	ConstantZero
	ConstantOne
	ConstantEmptyStr
	ConstantEmptyBytes
	ConstantEmptyTuple
)

// GetConstant returns a reference to the requested constant. An error returned
// if the requested constant c is not valid.
//
// Return value: New Reference.
func GetConstant(c ConstantType) (Object, error) {
	ret := C.Py_GetConstant(C.uint(c))
	return obj2ObjErr(ret)
}

// GetConstant returns a reference to the requested constant. An error returned
// if the requested constant c is not valid.
//
// Return value: Borrowed Reference.
func BorrowConstant(c ConstantType) (Object, error) {
	ret := C.Py_GetConstantBorrowed(C.uint(c))
	return obj2ObjErr(ret)
}
