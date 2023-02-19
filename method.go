// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

func NewMethod(function, self Object) (*Method, error) {
	ret := C.PyMethod_New(c(function), c(self))
	if ret == nil {
		return nil, exception()
	}
	return newMethod(ret), nil
}

func (m *Method) Self() Object {
	ret := C.PyMethod_Self(c(m))
	return newObject(ret)
}
