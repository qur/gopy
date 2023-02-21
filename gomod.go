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

// ChanClosedError is an exception raised by the Chan type when the channel has been
// closed. This variable is initialised as part of initialising the Go module,
// so it should not be used if the Go module is not initialised.
var ChanClosedError *ExceptionClass

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

	if err := goModule.AddObjectRef("Chan", &chanClass); err != nil {
		goModule.Decref()
		goModule = nil
		return nil, err
	}

	chanClosedError, err := NewException("go.ChanClosedError", nil, nil)
	if err != nil {
		goModule.Decref()
		goModule = nil
		return nil, err
	}

	if err := goModule.AddObjectRef("ChanClosedError", chanClosedError); err != nil {
		goModule.Decref()
		chanClosedError.Decref()
		goModule = nil
		return nil, err
	}

	if err := goModule.Register(); err != nil {
		goModule.Decref()
		chanClosedError.Decref()
		goModule = nil
		return nil, err
	}

	ChanClosedError = chanClosedError

	return goModule, nil
}
