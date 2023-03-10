package py

import (
	"fmt"
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

	rm := NewRefManager()
	defer rm.Decref()

	mod, err := CreateModule(&modDef)
	rm.Add(mod)
	if err != nil {
		return nil, fmt.Errorf("failed to create go module: %w", err)
	}

	if err := chanClass.Create(); err != nil {
		return nil, fmt.Errorf("failed to create chan class: %w", err)
	}

	if err := mod.AddObjectRef("Chan", &chanClass); err != nil {
		return nil, fmt.Errorf("failed to add Chan to go module (%p%T): %w", mod, mod, err)
	}

	chanClosedError, err := NewException("go.ChanClosedError", nil, nil)
	rm.Add(chanClosedError)
	if err != nil {
		return nil, err
	}

	if err := mod.AddObjectRef("ChanClosedError", chanClosedError); err != nil {
		return nil, fmt.Errorf("failed to add ChanCloserError to go module: %w", err)
	}

	if err := mod.Register(); err != nil {
		return nil, fmt.Errorf("failed to register go module: %w", err)
	}

	rm.Clear()
	ChanClosedError = chanClosedError
	goModule = mod

	return mod, nil
}
