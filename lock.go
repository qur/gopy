package py

// #include "utils.h"
import "C"

import (
	"runtime"
)

// GILState holds the GIL state for the current thread.
type GILState struct {
	state C.PyGILState_STATE
	set   bool
}

// GILStateEnsure will ensure that the current thread holds the GIL.
//
// Use of Lock is recommend over direct use of this function, as it will
// correctly manage the interaction with goroutines.
//
// When the function returns, the current thread will hold the GIL and be able
// to call arbitrary Python code. Failure is a fatal error.
//
// Note: This function is just a convenience for creating an empty GILState and
// calling Ensure.
func GILStateEnsure() *GILState {
	g := &GILState{}
	g.Ensure()
	return g
}

// GILStateEnsure will ensure that the current thread holds the GIL.
//
// Use of Lock is recommend over direct use of this function, as it will
// correctly manage the interaction with goroutines.
//
// When the function returns, the current thread will hold the GIL and be able
// to call arbitrary Python code. Failure is a fatal error.
func (g *GILState) Ensure() {
	if !g.set {
		g.state = C.PyGILState_Ensure()
		g.set = true
	}
}

// Release will release the GIL and restore the previously saved state.
//
// Use of Lock is recommend over direct use of this function, as it will
// correctly manage the interaction with goroutines.
//
// No Python calls should be until the GIL is obtained again.
func (g *GILState) Release() {
	if g.set {
		C.PyGILState_Release(g.state)
		g.set = false
	}
}

func initAndLock(initsigs bool) *Lock {
	// Lock the current goroutine to the current OS thread, until we have
	// released the GIL (as CPython uses per-thread state)
	runtime.LockOSThread()

	// Initialize the default Python interpreter
	InitializeEx(initsigs)

	// Immediately release the GIL (and thus "deactivate" any per-thread state
	// associated with the current thread
	C.PyEval_SaveThread()

	// We can now unlock the current goroutine from the current OS thread, as
	// there is no active per-thread state
	runtime.UnlockOSThread()

	// Now that Python is setup, we can return a locked Lock, ready for the
	// calling code to use
	return NewLock()
}

// InitAndLock is a convenience function.  It initializes Python, enables thread
// support, and returns a locked Lock instance.
func InitAndLock() *Lock {
	return initAndLock(false)
}

// InitAndLockWithSignals is similar to InitAndLock, except that it initializes
// the Python signal handling too.
func InitAndLockWithSignals() *Lock {
	return initAndLock(true)
}

// Lock is a high-level representation of the Python Global Interpreter Lock
// (GIL) and thread state.  When calling from Go into Python the GIL needs to be
// held and the appropriate thread state loaded.  Lock also calls
// runtime.LockOSThread() to make sure that the calling goroutine doesn't move
// thread whilst calling Python code which would invalidate the per-thread
// state.
//
// A single instance of Lock is not intended to be shared, but instead
// represents a private view of the GIL. If the Lock is locked then you have the
// GIL and can execute Python code, if the Lock is unlocked (or unblocked) then
// you don't have the GIL and can't execute Python code.
//
// Basic usage is:
//
//	lock = py.NewLock()
//
//	// Call Python code ...
//
//	lock.Unlock()
//
// If it appropriate to let other Python threads run (e.g. during a long
// computation, or blocking operation), then there are two options.  Either
// unlock:
//
//	lock = py.NewLock()
//
//	// Call Python code ...
//
//	lock.Unlock()
//
//	// Slow or blocking Go operation
//
//	lock.Lock()
//
//	// Call Python code ...
//
//	lock.Unlock()
//
// or unblock threads, which will not call runtime.UnlockOSThread() but it less
// expensive, as we do not free and then recreate a thread state variable:
//
//	lock = py.NewLock()
//
//	// Call Python code ...
//
//	lock.UnblockThreads()
//
//	// Slow or blocking Go operation
//
//	lock.BlockThreads()
//
//	// Call Python code ...
//
//	lock.Unlock()
//
// Note: The Lock and Unlock methods will panic if Python is not initialised, so
// Lock and Unlock should not be called before Python is initialised, or after
// Python is finalized.
//
// The simplest way to embed and use Python is using the InitAndLock and
// Lock.Finalize wrappers:
//
//	lock := py.InitAndLock()
//	defer py.Finalize()
type Lock struct {
	gilState *GILState
	thState  *C.PyThreadState
}

// NewLock returns a new locked Lock
func NewLock() (lock *Lock) {
	lock = &Lock{}
	lock.Lock()
	return
}

func (lock *Lock) setCount(l int64) {
	dict := newDict(C.PyThreadState_GetDict())
	if dict == nil {
		// no current thread state, Python is not initialised on this thread.
		panic("Python not initialised on this thread")
	}
	pl := NewLong(l)
	defer pl.Decref()
	err := dict.SetItemString("gopy.count", pl)
	if err != nil {
		panic(err)
	}
}

func (lock *Lock) getCount() int64 {
	dict := newDict(C.PyThreadState_GetDict())
	if dict == nil {
		// no current thread state, Python is not initialised on this thread.
		panic("Python not initialised on this thread")
	}
	val, err := dict.GetItemString("gopy.count")
	if err != nil {
		return 0
	}
	count, ok := val.(*Long)
	if !ok {
		return 0
	}
	return count.Int64()
}

func (lock *Lock) inc() {
	lock.setCount(lock.getCount() + 1)
}

func (lock *Lock) dec() bool {
	count := lock.getCount()
	if count <= 0 {
		panic("Lock.dec() called with count <1!")
	}
	count--
	lock.setCount(count)
	return count == 0
}

// Lock locks the lock.  When it returns everything is setup for calling into
// Python.  No other Python threads will run until either Unlock() or
// UnblockThreads() are called.
//
// If the lock is already locked when this function is called, then nothing
// happens, and the function will return immediately.
func (lock *Lock) Lock() {
	if lock.gilState != nil {
		return
	}

	if lock.thState != nil {
		panic("Lock() called with threads unblocked!")
	}

	runtime.LockOSThread()
	lock.gilState = GILStateEnsure()
	lock.inc()
}

// Unlock unlocks the lock.  When it returns no calls into Python may be made.
//
// If the lock is not locked when this function is called, then nothing happens,
// and the function returns immediately.  Also, it is not necessary to call
// BlockThreads() before calling Unlock(), even if UnblockThreads() has been
// called.
func (lock *Lock) Unlock() {
	if lock.gilState == nil {
		return
	}

	if lock.thState != nil {
		C.PyEval_RestoreThread(lock.thState)
		lock.thState = nil
	}

	lock.dec()
	lock.gilState.Release()
	runtime.UnlockOSThread()

	lock.gilState = nil
}

// UnblockThreads() releases the GIL so that other Python threads may run.  It
// does not free the per-thread state created by Lock, nor does it call
// runtime.UnlockOSThread().  This function is intended to allow other Python
// threads to run whilst the calling code is either performing a slow/long
// running operation or is going to block.
//
// Nothing happens if this function is called more than once, all calls but the
// first will be ignored.
//
// If the lock is not locked, this function will panic.
func (lock *Lock) UnblockThreads() {
	if lock.gilState == nil {
		panic("UnblockThreads() called on Unlocked Lock")
	}

	if lock.thState == nil {
		lock.thState = C.PyEval_SaveThread()
	}
}

// BlockThreads() reclaims the GIL (and restores per-thread state), after it has
// been released by UnblockThreads().
//
// If this function is called without UnblockThreads() having been called, then
// nothing happens and the function returns immediately.
//
// If the lock is not locked, this function will panic.
func (lock *Lock) BlockThreads() {
	if lock.gilState == nil {
		panic("BlockThreads() called on Unlocked Lock")
	}

	if lock.thState != nil {
		C.PyEval_RestoreThread(lock.thState)
		lock.thState = nil
	}
}

// Finalize shuts down the Python runtime.
//
// Finalize also unlocks the Lock, also calling Unlock will cause Python to
// crash.
//
// If the Lock is not locked, or if this is not the outer-most lock of a nested
// set then this function will panic.
//
// However, it is not necessary to call BlockThreads() before calling
// Finalize(), even if UnblockThreads() has been called.
func (lock *Lock) Finalize() {
	if lock.gilState == nil {
		panic("Lock.Finalize called when Lock is not locked")
	}

	if lock.thState != nil {
		C.PyEval_RestoreThread(lock.thState)
		lock.thState = nil
	}

	locksActive := lock.dec()
	if !locksActive {
		panic("Lock.Finalize called when Locks are still active")
	}

	C.Py_Finalize()

	runtime.UnlockOSThread()
	lock.gilState = nil
}

// UnblockThreads will release the GIL to allow other threads to run. It it
// intended only for uses where using a Lock to ensure the GIL is held is not
// needed. The return value is a niladic function to call to regain the GIL.
//
// No check is made that the GIL is currently held.
//
// An example usage might be:
//
//	// Call Python code ...
//
//	blockThreads := py.UnblockThreads()
//
//	// Slow or blocking Go operation
//
//	blockThreads()
//
//	// Call Python code ...
func UnblockThreads() func() {
	thState := C.PyEval_SaveThread()
	return func() {
		C.PyEval_RestoreThread(thState)
	}
}
