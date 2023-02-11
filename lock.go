// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"runtime"
)

type GilState struct {
	state C.PyGILState_STATE
	set   bool
}

func GilState_Ensure() *GilState {
	ret := C.PyGILState_Ensure()
	return &GilState{ret, true}
}

func (g *GilState) Release() {
	if g.set {
		C.PyGILState_Release(g.state)
		g.set = false
	}
}

func AcquireGil() {
	C.PyEval_AcquireLock()
}

func ReleaseGil() {
	C.PyEval_ReleaseLock()
}

// InitAndLock is a convience function.  It initializes Python, enables thread
// support, and returns a locked Lock instance.
func InitAndLock() *Lock {
	return initAndLock(false)
}

// InitAndLockWithSignals is similar to InitAndLock, except that it initializes
// the Python signal handling too.
func InitAndLockWithSignals() *Lock {
	return initAndLock(true)
}

func initAndLock(initsigs bool) *Lock {
	// Lock the current goroutine to the current OS thread, until we have
	// released the GIL (as CPython uses per-thread state)
	runtime.LockOSThread()

	// Initialize the default Python interpreter
	if initsigs {
		C.Py_InitializeEx(1)
	} else {
		C.Py_InitializeEx(0)
	}

	// Enable Python thread support, and then immediately release the GIL (and
	// thus "deativate" and per-thread state associated with the current thread
	C.PyEval_InitThreads()
	C.PyEval_SaveThread()

	// We can now unlock the current goroutine from the current OS thread, as
	// there is no active per-thread state
	runtime.UnlockOSThread()

	// Now that Python is setup, we can return a locked Lock, ready for the
	// calling code to use
	return NewLock()
}

// Lock is a high-level representation of the Python Global Interpreter Lock
// (GIL) and thread state.  When calling from Go into Python the GIL needs to be
// held and the appropriate thread state loaded.  Lock also calls
// runtime.LockOSThread() to make sure that the calling goroutine doesn't move
// thread whilst calling Python code which would invalidate the per-thread
// state.
//
// Basic usage is:
//
//   lock = py.NewLock()
//
//   // Call Python code ...
//
//   lock.Unlock()
//
// If it appropriate to let other Python threads run (e.g. during a long
// computation, or blocking operation), then there are two options.  Either
// unlock:
//
//   lock = py.NewLock()
//
//   // Call Python code ...
//
//   lock.Unlock()
//
//   // Slow or blocking Go operation
//
//   lock.Lock()
//
//   // Call Python code ...
//
//   lock.Unlock()
//
// or unblock threads, which will not call runtme.UnlockOSThread() but it less
// expensive, as we do not free and then recreate a thread state variable:
//
//   lock = py.NewLock()
//
//   // Call Python code ...
//
//   lock.UnblockThreads()
//
//   // Slow or blocking Go operation
//
//   lock.BlockThreads()
//
//   // Call Python code ...
//
//   lock.Unlock()
type Lock struct {
	gilState *GilState
	thState  *C.PyThreadState
}

// NewLock returns a new locked Lock
func NewLock() (lock *Lock) {
	lock = new(Lock)
	lock.Lock()
	return
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
	lock.gilState = GilState_Ensure()
	lock.inc()
}

func (lock *Lock) setCount(l int64) {
	dict := newDict(C.PyThreadState_GetDict())
	pl := NewLong(l)
	defer pl.Decref()
	err := dict.SetItemString("gopy.count", pl)
	if err != nil {
		panic(err)
	}
}

func (lock *Lock) getCount() int64 {
	dict := newDict(C.PyThreadState_GetDict())
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

// Unlock unlocks the lock.  When it returns no calls into Python may be made.
//
// If the lock is not locked when this function is called, then nothing happens,
// and the function returns immediately.  Also, it is not necessay to call
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

	releaseOsThread := lock.dec()

	lock.gilState.Release()
	if releaseOsThread {
		runtime.UnlockOSThread()
	}
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
func (lock *Lock) UnblockThreads() {
	if lock.gilState == nil {
		panic("UnblockThreads() called on Unlocked Lock")
	}

	if lock.thState == nil {
		lock.thState = C.PyEval_SaveThread()
	}
}

// BlockThreads() reclaims the GIL (and restores per-thread state), after is has
// been released by UnblockThreads().
//
// If this function is called without UnblockThreads() having been called, then
// nothing happens and the function returns immediately.
func (lock *Lock) BlockThreads() {
	if lock.gilState == nil {
		panic("BlockThreads() called on Unlocked Lock")
	}

	if lock.thState != nil {
		C.PyEval_RestoreThread(lock.thState)
		lock.thState = nil
	}
}
