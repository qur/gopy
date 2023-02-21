package py

import (
	"sync"
)

var (
	goModLock sync.Mutex
	goModule  *Module

	modDef = ModuleDef{
		Name: "go",
		Doc:  "go interface module",
	}
)

// InitGoModule creates and registers (and returns) the special built-in "go"
// module. This module exports go specific items to Python that can be accessed
// by calling "import go".  Since this function also initialises the Python
// representations that are exported by the "go" module it must be called before
// any of these types (e.g. Chan) can be used.
//
// This function may be called more than once to get at the *Module, the module
// will only be created and registered once.
//
// Note: The return value is a borrowed ref.
func InitGoModule() (*Module, error) {
	goModLock.Lock()
	defer goModLock.Unlock()

	if goModule != nil {
		return goModule, nil
	}

	goModule, err := CreateModule(&modDef)
	if err != nil {
		return nil, err
	}

	if err := chanClass.Create(); err != nil {
		goModule.Decref()
		goModule = nil
		return nil, err
	}

	if err := goModule.AddObject("Chan", &chanClass); err != nil {
		goModule.Decref()
		goModule = nil
		return nil, err
	}

	if err := goModule.Register(); err != nil {
		goModule.Decref()
		goModule = nil
		return nil, err
	}

	return goModule, nil
}
