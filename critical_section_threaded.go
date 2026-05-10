//go:build threaded

package py

// #include "utils.h"
import "C"

import (
	"unsafe"
)

// CriticalSection is a mechanism for locking objects when making changes to
// internal state. It is part of the threaded replacements for the GIL, and is a
// no-op in non-threaded builds.
//
// Critical sections avoid deadlocks by implicitly suspending active critical
// sections, hence, they do not provide exclusive access such as provided by
// traditional locks like PyMutex. When a critical section is started, the
// per-object lock for the object is acquired. If the code executed inside the
// critical section calls C-API functions then it can suspend the critical
// section thereby releasing the per-object lock, so other threads can acquire
// the per-object lock for the same object.
//
// Variants that accept PyMutex pointers rather than Python objects are also
// available. Use these variants to start a critical section in a situation
// where there is no PyObject – for example, when working with a C type that
// does not extend or wrap PyObject but still needs to call into the C API in a
// manner that might lead to deadlocks.
type CriticalSection struct {
	cs C.PyCriticalSection
}

// NewCriticalSection creates a new CriticalSection ready for use.
func NewCriticalSection() *CriticalSection {
	return (*CriticalSection)(unsafe.Pointer(C.newCriticalSection()))
}

// Begin starts a critical section for the given object.
func (cs *CriticalSection) Begin(o Object) {
	C.PyCriticalSection_Begin(&cs.cs, c(o))
}

// BeginMutex starts a critical section for the given Mutex.
func (cs *CriticalSection) BeginMutex(m *Mutex) {
	C.PyCriticalSection_BeginMutex(&cs.cs, &m.m)
}

// End ends the critical section.
func (cs *CriticalSection) End() {
	C.PyCriticanSection_End(&cs.cs)
}

// CriticalSection2 is a mechanism for locking two objects when making changes
// to internal state. It is part of the threaded replacements for the GIL, and
// is a no-op in non-threaded builds.
//
// Critical sections avoid deadlocks by implicitly suspending active critical
// sections, hence, they do not provide exclusive access such as provided by
// traditional locks like PyMutex. When a critical section is started, the
// per-object lock for the object is acquired. If the code executed inside the
// critical section calls C-API functions then it can suspend the critical
// section thereby releasing the per-object lock, so other threads can acquire
// the per-object lock for the same object.
//
// Variants that accept PyMutex pointers rather than Python objects are also
// available. Use these variants to start a critical section in a situation
// where there is no PyObject – for example, when working with a C type that
// does not extend or wrap PyObject but still needs to call into the C API in a
// manner that might lead to deadlocks.
type CriticalSection2 struct {
	cs C.PyCriticalSection2
}

// NewCriticalSection2 creates a new CriticalSection2 ready for use.
func NewCriticalSection2() *CriticalSection2 {
	return (*CriticalSection2)(unsafe.Pointer(C.newCriticalSection2()))
}

// Begin starts a critical section for the given objects.
func (cs *CriticalSection2) Begin(o1, o2 Object) {
	C.PyCriticalSection2_Begin(&cs.cs, c(o1), c(o2))
}

// BeginMutex starts a critical section for the given Mutexs.
func (cs *CriticalSection2) BeginMutex2(m1, m2 *Mutex) {
	C.PyCriticalSection2_BeginMutex(&cs.cs, &m1.m, &m2.m)
}

// End ends the critical section.
func (cs *CriticalSection2) End() {
	C.PyCriticanSection2_End(&cs.cs)
}
