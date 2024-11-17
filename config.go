package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

type ConfigFlag int

const (
	defaultFlag ConfigFlag = iota
	Enabled
	Disabled
)

func (f ConfigFlag) set(v *C.int) {
	switch f {
	case defaultFlag:
	case Enabled:
		*v = 1
	case Disabled:
		*v = 0
	default:
		panic(fmt.Sprintf("invalid ConfigFlag value: %d", f))
	}
}

type PreConfig struct {
	Args     []string
	Isolated ConfigFlag
}

func PythonPreConfig(args ...string) *PreConfig {
	return &PreConfig{
		Args:     args,
		Isolated: Disabled,
	}
}

func IsolatedPreConfig() *PreConfig {
	return &PreConfig{
		Isolated: Enabled,
	}
}

func (c *PreConfig) PreInitialize() error {
	cfg := C.PyPreConfig{}
	if c.Isolated == Enabled {
		C.PyPreConfig_InitIsolatedConfig(&cfg)
	} else {
		C.PyPreConfig_InitPythonConfig(&cfg)
	}

	c.Isolated.set(&cfg.isolated)

	if c.Args == nil {
		return status2Err(C.Py_PreInitialize(&cfg))
	}

	argv := make([]*C.char, len(c.Args))

	for i, arg := range c.Args {
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	return status2Err(C.Py_PreInitializeFromBytesArgs(&cfg, C.Py_ssize_t(len(argv)), &argv[0]))
}

func setArgs(args []string, cfg *C.PyConfig) error {
	argv := make([]*C.char, len(args))

	for i, arg := range args {
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	return status2Err(C.PyConfig_SetBytesArgv(cfg, C.Py_ssize_t(len(argv)), &argv[0]))
}

type Config struct {
	Args     []string
	SafePath ConfigFlag
	// base_exec_prefix
	// base_executable
	// base_prefix

	Isolated ConfigFlag

	InstallSignalHandlers ConfigFlag
}

func PythonConfig(args ...string) *Config {
	return &Config{
		Args:     args,
		Isolated: Disabled,
	}
}

func IsolatedConfig() *Config {
	return &Config{
		Isolated: Enabled,
	}
}

func (c *Config) Initialize() error {
	cfg := C.PyConfig{}
	if c.Isolated == Enabled {
		C.PyConfig_InitIsolatedConfig(&cfg)
	} else {
		C.PyConfig_InitPythonConfig(&cfg)
	}
	defer C.PyConfig_Clear(&cfg)

	if c.Args != nil {
		if err := setArgs(c.Args, &cfg); err != nil {
			return err
		}
	}

	c.SafePath.set(&cfg.safe_path)
	c.Isolated.set(&cfg.isolated)
	c.InstallSignalHandlers.set(&cfg.install_signal_handlers)

	return status2Err(C.Py_InitializeFromConfig(&cfg))
}

func (c *Config) InitAndLock() (*Lock, error) {
	// Lock the current goroutine to the current OS thread, until we have
	// released the GIL (as CPython uses per-thread state)
	runtime.LockOSThread()

	// Initialize the default Python interpreter
	if err := c.Initialize(); err != nil {
		runtime.UnlockOSThread()
		return nil, err
	}

	// Immediately release the GIL (and thus "deactivate" any per-thread state
	// associated with the current thread
	C.PyEval_SaveThread()

	// We can now unlock the current goroutine from the current OS thread, as
	// there is no active per-thread state
	runtime.UnlockOSThread()

	// Now that Python is setup, we can return a locked Lock, ready for the
	// calling code to use
	return NewLock(), nil
}
