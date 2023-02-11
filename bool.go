// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// *Bool is the representation of the Python bool type.  There are only two
// possible values for a Bool, True and False.  Every True value refers to the
// same instance, and every False value refers to the same value.
type Bool struct {
	AbstractObject
	o C.PyBoolObject
}

// BoolType is the Type object that represents the Bool type.
var BoolType = (*Type)(unsafe.Pointer(&C.PyBool_Type))

// True is the true value of the Bool type.  It is a singleton value, all true
// values refer to the same instance.
var True = (*Bool)(C.pyTrue())

// False is the false value of the Bool type.  It is a singleton value, all
// false values refer to the same instance.
var False = (*Bool)(C.pyFalse())

func boolCheck(obj Object) bool {
	return C.boolCheck(c(obj)) != 0
}

func newBool(obj *C.PyObject) *Bool {
	if obj == c(True) {
		return True
	}
	if obj == c(False) {
		return False
	}
	panic(TypeError.Err("not a bool"))
}

// Bool returns the value of "b" as a bool.  true for True, false for False.  If
// "b" is neither True nor False then this function will panic.
func (b *Bool) Bool() bool {
	if b == True {
		return true
	}
	if b == False {
		return false
	}
	panic(TypeError.Err("not a bool"))
}

// String returns a printable representation of the Bool "b".
func (b *Bool) String() string {
	if b == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", b.Bool())
}
