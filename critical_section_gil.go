//go:build !threaded

package py

// CriticalSection is a mechanism for locking objects when making changes to
// internal state. It is part of the threaded replacements for the GIL, and is a
// no-op in non-threaded builds.
type CriticalSection struct{}

// NewCriticalSection creates a new CriticalSection ready for use.
func NewCriticalSection() *CriticalSection {
	return &CriticalSection{}
}

// Begin starts a critical section for the given object.
func (cs *CriticalSection) Begin(Object) {}

// BeginMutex starts a critical section for the given Mutex.
func (cs *CriticalSection) BeginMutex(*Mutex) {}

// End ends the critical section.
func (cs *CriticalSection) End() {}

// CriticalSection2 is a mechanism for locking two objects when making changes
// to internal state. It is part of the threaded replacements for the GIL, and
// is a no-op in non-threaded builds.
type CriticalSection2 struct{}

// NewCriticalSection2 creates a new CriticalSection2 ready for use.
func NewCriticalSection2() *CriticalSection2 {
	return &CriticalSection2{}
}

// Begin starts a critical section for the given objects.
func (cs *CriticalSection2) Begin(_, _ Object) {}

// BeginMutex starts a critical section for the given Mutexs.
func (cs *CriticalSection2) BeginMutex2(_, _ *Mutex) {}

// End ends the critical section.
func (cs *CriticalSection2) End() {}
