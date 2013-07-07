// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import "sync"

var (
	goModLock sync.Mutex
	goModule  *Module
)

// InitGoModules initializes (and returns) the special built-in "go" module.
// This module exports go specific items to Python that can be accessed by
// calling "import go".  Since this function also initialises the Python
// representations that are exported by the "go" module it must be called before
// any of these types (e.g. Chan) can be used.
//
// This function may be called more than once to get at the *Module, the module
// will only be created and initialized once.
func InitGoModule() (*Module, error) {
	goModLock.Lock()
	defer goModLock.Unlock()

	if goModule != nil {
		return goModule, nil
	}

	goModule, err := InitModule("go", []Method{})
	if err != nil {
		return nil, err
	}

	c, err := chanClass.Create()
	if err != nil {
		goModule.Decref()
		goModule = nil
		return nil, err
	}

	err = goModule.AddObject("Chan", c)
	if err != nil {
		goModule.Decref()
		goModule = nil
		c.Decref()
		return nil, err
	}

	return goModule, nil
}
