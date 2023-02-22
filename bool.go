package py

// #include "utils.h"
import "C"

import (
	"fmt"
)

// True is the true value of the Bool type.  It is a singleton value, all true
// values refer to the same instance.
var True = (*Bool)(C.pyTrue())

// False is the false value of the Bool type.  It is a singleton value, all
// false values refer to the same instance.
var False = (*Bool)(C.pyFalse())

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
	switch b {
	case True:
		return true
	case False:
		return false
	default:
		panic(TypeError.Err("not a bool"))
	}
}

// String returns a printable representation of the Bool "b".
func (b *Bool) String() string {
	if b == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", b.Bool())
}
