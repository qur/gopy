package py

// #cgo CFLAGS: -Werror
// #cgo pkg-config: python-3.13-embed libffi
//
// #include "utils.h"
import "C"

import (
	"log"
	"unsafe"
)

// Code generation commands:
//go:generate ./gen_exc.py exc.go python3.13
//go:generate ./gen_slots.py
//go:generate ./gen_types.py python3.13
//go:generate ./gen_arg.py

// nilValue is the value returned from String methods when the value is nil.
const nilValue = "<nil>"

// Initialize initialises the Python runtime.
//
// You probably want InitAndLockWithSignals though, as it doesn't require the
// caller to worry about goroutines or threads.
func Initialize() {
	InitializeEx(true)
}

// InitializeEx initialises the Python runtime.
//
// If initsigs is true then the Python runtime will install signal handlers.
//
// You probably want InitAndLock or InitAndLockWithSignals though, as they
// doesn't require the caller to worry about goroutines or threads.
func InitializeEx(initsigs bool) {
	if initsigs {
		C.Py_InitializeEx(1)
	} else {
		C.Py_InitializeEx(0)
	}

	if err := setupImporter(); err != nil {
		log.Printf("failed to setup importer: %s", err)
	}
}

// Finalize shuts down the Python runtime.
//
// You probably want to call the Lock.Finalize method though, as it will ensure
// that goroutines and threads are managed correctly.
func Finalize() {
	C.Py_Finalize()
}

// AddToPath appends the given directory to sys.path.
func AddToPath(dir string) error {
	p := C.CString("path")
	defer cfree(p)

	sysPath := C.PySys_GetObject(p)
	if sysPath == nil {
		return AttributeError.Err("path")
	}

	s := C.CString(dir)
	defer cfree(s)

	pDir := C.PyUnicode_FromString(s)
	if pDir == nil {
		return exception()
	}
	defer decref(pDir)

	return int2Err(C.PyList_Append(sysPath, pDir))
}

// PrependToPath prepends the given directory to sys.path.
func PrependToPath(dir string) error {
	p := C.CString("path")
	defer cfree(p)

	sysPath := C.PySys_GetObject(p)
	if sysPath == nil {
		return AttributeError.Err("path")
	}

	s := C.CString(dir)
	defer cfree(s)

	pDir := C.PyUnicode_FromString(s)
	if pDir == nil {
		return exception()
	}
	defer decref(pDir)

	return int2Err(C.PyList_Insert(sysPath, 0, pDir))
}

// Main is the main Python interpreter entrypoint.
//
// Once this function returns, the Python runtime is shutdown.
func Main(args []string) int {
	argv := make([]*C.char, len(args))

	for i, arg := range args {
		argv[i] = C.CString(arg)
		defer cfree(argv[i])
	}

	return int(C.Py_BytesMain(C.int(len(argv)), &argv[0]))
}

// Main is the main Python interpreter entrypoint used when using the newer
// Config based initialisation.
//
// Once this function returns, the Python runtime is shutdown.
func RunMain() int {
	return int(C.Py_RunMain())
}

// EnterRecursiveCall marks a point where a recursive Go-level call is about to
// be performed.  It returns true if the recursive call is permitted, otherwise
// a Python exception is set and false is returned.  where is a string that will
// be appended to the RuntimeError set if the recursion limit has been exceeded
// (e.g. " in instance check").  This function needs to be called if the
// recursive function may not invoke Python code (which automatically tracks
// recursion depth).
func EnterRecursiveCall(where string) bool {
	s := C.CString(where)
	defer cfree(s)

	return C.enterRecursive(s) == 0
}

// LeaveRecursiveCall must be called after a recursive call that was indicated
// by EnterRecursiveCall.
func LeaveRecursiveCall() {
	C.leaveRecursive()
}

func cfree[T any](value *T) {
	C.free(unsafe.Pointer(value))
}
