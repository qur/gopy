package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// Mutex is a Python C API mutual exclusion lock. This is not a Python type, but
// is used as part of the API only.
type Mutex struct {
	m C.PyMutex
}

// NewMutex returns a new unlock Mutex.
func NewMutex() *Mutex {
	return (*Mutex)(unsafe.Pointer(C.newMutex()))
}

// Lock locks the Mutex m. If another thread has already locked it, the calling
// thread will block until the mutex is unlocked. While blocked, the thread will
// temporarily detach the thread state if one exists.
func (m *Mutex) Lock() {
	C.PyMutex_Lock(&m.m)
}

// Unlock unlock the Mutex m. The mutex must be locked — otherwise, the function
// will issue a fatal error.
func (m *Mutex) Unlock() {
	C.PyMutex_Unlock(&m.m)
}

// Returns true if the Mutex m is locked, false otherwise.
//
// Note: This function is intended for use in assertions and debugging only and
// should not be used to make concurrency control decisions, as the lock state
// may change immediately after the check.
func (m *Mutex) IsLocked() bool {
	return C.PyMutex_IsLocked(&m.m) != 0
}
