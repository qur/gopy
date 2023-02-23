package py

// #include "utils.h"
import "C"

// Chan is a Python object that wraps a Go channel (specifically a "chan
// Object").
type Chan struct {
	ClassBaseObject
	c chan Object
}

// newChan is the New function for the Chan Class that is called as part of the
// object initialisation.
func newChan(c *Class, args *Tuple, kw *Dict) (ClassObject, error) {
	var buffer int
	if err := ParseTuple(args, "i", &buffer); err != nil {
		return nil, err
	}

	ch := &Chan{
		c: make(chan Object, buffer),
	}

	return ch, nil
}

// NewChan returns a new Chan object, with the channel created using "make(chan
// Object, buffer)".
//
// Note: This function should only be called if InitGoModule() has been called,
// as otherwise the Chan Python type has not been initialized. It will always
// return an error if InitGoModule() has not been called.
//
// Return value: New Reference.
func NewChan(buffer int) (*Chan, error) {
	goModLock.Lock()
	defer goModLock.Unlock()

	if goModule == nil {
		return nil, AssertionError.Err("go module has not been initialized!")
	}

	obj, err := chanClass.Call(nil, nil)
	if err != nil {
		return nil, err
	}

	self, ok := obj.(*Chan)
	if !ok {
		defer obj.Decref()
		return nil, TypeError.Err("Alloc returned wrong type: %T", obj)
	}

	return self, nil
}

// Py_put provides a c.put() method when this object is used in Python.
func (c *Chan) Py_put(args *Tuple, kw *Dict) (ret Object, err error) {
	var obj Object
	kwlist := []string{"obj"}

	if err := ParseTupleAndKeywords(args, kw, "O", kwlist, &obj); err != nil {
		return nil, err
	}

	obj.Incref()
	None.Incref()
	ret = None

	defer func() {
		if p := recover(); p != nil {
			obj.Decref()
			ret.Decref()
			ret = nil
			err = ChanClosedError.Err("Chan closed")
		}
	}()

	c.c <- obj

	return
}

// Py_get provides a c.get() method when this object is used in Python.
func (c *Chan) Py_get(args *Tuple, kw *Dict) (Object, error) {
	if err := ParseTupleAndKeywords(args, kw, "", []string{}); err != nil {
		return nil, err
	}

	obj := <-c.c

	if obj == nil {
		return nil, ChanClosedError.Err("Chan closed")
	}

	return obj, nil
}

// Py_close provides a c.close() method when this object is used in Python.
func (c *Chan) Py_close(args *Tuple, kw *Dict) (ret Object, err error) {
	if err := ParseTupleAndKeywords(args, kw, "", []string{}); err != nil {
		return nil, err
	}

	None.Incref()
	ret = None

	defer func() {
		if p := recover(); p != nil {
			ret.Decref()
			ret = nil
			err = ChanClosedError.Err("Chan already closed")
		}
	}()

	close(c.c)

	return
}

// PyIter is called to get an iterator for the item. This is used when running
// "for ... in c" in Python.
func (c *Chan) PyIter() (Object, error) {
	c.Incref()
	return c, nil
}

// PyIterNext is the iterator API. This function will get called to return the
// next value.
func (c *Chan) PyIterNext() (Object, error) {
	obj := <-c.c

	return obj, nil
}

// Chan returns the channel inside the Chan object.  This is how Go code can get
// at the channel to communicate with Python code.
//
// When sending values, a reference count will be "stolen" - i.e. the reference
// count should be incremented if you want to continue to hold a reference after
// the send.
//
// When receiving values, you get a new reference - i.e. the reference count
// should be decremented when you have finished with the value.
func (c *Chan) Chan() chan Object {
	return c.c
}

var chanClass = Class{
	Name:   "go.Chan",
	Object: (*Chan)(nil),
	New:    newChan,
}
