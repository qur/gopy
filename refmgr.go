package py

// #include "utils.h"
import "C"

// RefManager is a helper for manager objects references. Objects can be added,
// and when Decref is called then Decref will be called on all the contained
// objects.
//
// Note that RefManager is not thread-safe.
type RefManager struct {
	objects map[Object]bool
}

// NewRefManager returns a new RefManager ready to use.
func NewRefManager() *RefManager {
	return &RefManager{
		objects: make(map[Object]bool),
	}
}

// Add adds the given object to RefManager stealing a reference.
//
// For convenience the passed in value is returned.
//
// Return value: Borrowed reference.
func (r *RefManager) Add(obj Object) Object {
	r.objects[obj] = true
	return obj
}

func (r *RefManager) add(o *C.PyObject) Object {
	obj := newObject(o)
	r.objects[obj] = true
	return obj
}

// AddE adds the given object to RefManager stealing a reference. The error
// value is ignored (i.e. if obj and err are both != nil then obj will still be
// added).
//
// For convenience the passed in values are returned.
//
// Return value: Borrowed reference.
func (r *RefManager) AddE(obj Object, err error) (Object, error) {
	return r.Add(obj), err
}

// AddRef adds the given object to RefManager. If the object is not already in
// the RefManager, then the reference count is also incremented.
//
// For convenience the passed in value is returned.
//
// Return value: Borrowed reference.
func (r *RefManager) AddRef(obj Object) Object {
	if r.objects[obj] {
		return obj
	}
	Incref(obj)
	r.objects[obj] = true
	return obj
}

func (r *RefManager) addRef(o *C.PyObject) Object {
	obj := newObject(o)
	if r.objects[obj] {
		return obj
	}
	Incref(obj)
	r.objects[obj] = true
	return obj
}

// AddRefE adds the given object to RefManager. If the object is not already in
// the RefManager, then the reference count is also incremented. The error
// value is ignored (i.e. if obj and err are both != nil then obj will still be
// added).
//
// For convenience the passed in values are returned.
//
// Return value: Borrowed reference.
func (r *RefManager) AddRefE(obj Object, err error) (Object, error) {
	return r.AddRef(obj), err
}

// AddAll adds the given objects to RefManager, stealing references.
func (r *RefManager) AllAll(objects ...Object) {
	for _, obj := range objects {
		r.Add(obj)
	}
}

// AddAllRef adds the given objects to RefManager. If an object is not already
// in the RefManager, then the reference count is also incremented.
func (r *RefManager) AllAllRef(objects ...Object) {
	for _, obj := range objects {
		r.AddRef(obj)
	}
}

// Decref will call decref on all of the Objects stored in the RefManager. The
// RefManager will also be emptied, so it is ok to keep using it after the call
// as only object added after the call will be tracked.
func (r *RefManager) Decref() {
	for obj := range r.objects {
		obj.Decref()
		delete(r.objects, obj)
	}
}

// Clear empties the RefManager without changing the references.
//
// This can be called between Adding item and calling Decref, to prevent Decref
// from taking any action.
func (r *RefManager) Clear() {
	for obj := range r.objects {
		delete(r.objects, obj)
	}
}

// Remove removes the given object from the manager. The reference held by the
// RefManager is returned.
//
// If obj was not in RefManager, then a new reference to obj will be returned.
//
// Return value: New reference.
func (r *RefManager) Remove(obj Object) Object {
	if r.objects[obj] {
		delete(r.objects, obj)
		return obj
	}
	Incref(obj)
	return obj
}
