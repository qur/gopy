// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

// Chan is a Python object that wraps a Go channel (specifically a "chan
// Object").
type Chan struct {
	BaseObject
	c chan Object
}

func newChan(t *Type, args *Tuple, kw *Dict) (Object, error) {
	obj, err := t.Alloc(0)
	if err != nil {
		return nil, err
	}

	self, ok := obj.(*Chan)
	if !ok {
		defer obj.Decref()
		return nil, TypeError.Err("Alloc returned wrong type: %T", obj)
	}

	var buffer int

	err = ParseTuple(args, "i", &buffer)
	if err != nil {
		return nil, err
	}

	self.c = make(chan Object, buffer)

	return self, nil
}

// NewChan returns a new Chan object, with the channel created using "make(chan
// Object, buffer)".
//
// Note: This function should only be called if InitGoModule() has been called,
// as otherwise the Chan Python type has not been initialized.
//
// Return value: New Reference.
func NewChan(buffer int) (*Chan, error) {
	if goModule == nil {
		return nil, AssertionError.Err("go module has not been initialized!")
	}

	obj, err := chanClass.Alloc(0)
	if err != nil {
		return nil, err
	}

	self, ok := obj.(*Chan)
	if !ok {
		defer obj.Decref()
		return nil, TypeError.Err("Alloc returned wrong type: %T", obj)
	}

	self.c = make(chan Object, buffer)

	return self, nil
}

// PyDealloc is the deallocator for a Chan instance, it is used internally -
// c.Decref() should be used as normal.
func (c *Chan) PyDealloc() {
	c.c = nil
	c.Free()
}

// Py_put provides a c.put() method when this object is used in Python.
func (c *Chan) Py_put(args *Tuple, kw *Dict) (Object, error) {
	var obj Object
	kwlist := []string{"obj"}

	err := ParseTupleAndKeywords(args, kw, "O", kwlist, &obj)
	if err != nil {
		return nil, err
	}

	obj.Incref()
	c.c <- obj

	None.Incref()
	return None, nil
}

// Py_get provides a c.get() method when this object is used in Python.
func (c *Chan) Py_get(args *Tuple, kw *Dict) (Object, error) {
	err := ParseTupleAndKeywords(args, kw, "", []string{})
	if err != nil {
		return nil, err
	}

	obj := <-c.c

	return obj, nil
}

// Chan returns the channel inside the Chan object.  This is how Go code can get
// at the channel to communicate with Python code.
//
// When sending values, a reference count will be "stollen" - i.e. the reference
// count should be incremented if you want to continue to hold a reference after
// the send.
//
// When receiveing values, you get a new reference - i.e. the reference count
// should be decremented when you have finished with the value.
func (c *Chan) Chan() chan Object {
	return c.c
}

var chanClass = Class{
	Name:    "go.Chan",
	Pointer: (*Chan)(nil),
	New:     newChan,
}
