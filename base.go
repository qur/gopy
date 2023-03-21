package py

// #include "utils.h"
import "C"

import "unsafe"

// BaseType is the Type object that represents the BaseObject type.
var BaseType = newType(&C.PyBaseObject_Type)

// Free deallocates the storage (in Python) for obj.  After calling this method,
// obj should no longer be used.
func (obj *BaseObject) Free() {
	free(obj)
}

func (obj *BaseObject) raw() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(obj))
}

// HasAttr returns true if "obj" has the attribute "name".  This is equivalent
// to the Python "hasattr(obj, name)".
func (obj *BaseObject) HasAttr(name Object) bool {
	ret := C.PyObject_HasAttr(c(obj), c(name))
	return ret == 1
}

// HasAttrString returns true if "obj" has the attribute "name".  This is
// equivalent to the Python "hasattr(obj, name)".
func (obj *BaseObject) HasAttrString(name string) bool {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyObject_HasAttrString(c(obj), s)
	return ret == 1
}

// GetAttr returns the attribute of "obj" with the name "name".  This is
// equivalent to the Python "obj.name".
//
// Return value: New Reference.
func (obj *BaseObject) GetAttr(name Object) (Object, error) {
	ret := C.PyObject_GetAttr(c(obj), c(name))
	return obj2ObjErr(ret)
}

// GetAttrString returns the attribute of "obj" with the name "name".  This is
// equivalent to the Python "obj.name".
//
// Return value: New Reference.
func (obj *BaseObject) GetAttrString(name string) (Object, error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyObject_GetAttrString(c(obj), s)
	return obj2ObjErr(ret)
}

// PyObject_GenericGetAttr : This is an internal helper function - we shouldn't
// need to expose it ...

// SetAttr sets the attribute of "obj" with the name "name" to "value".  This is
// equivalent to the Python "obj.name = value".
func (obj *BaseObject) SetAttr(name, value Object) error {
	ret := C.PyObject_SetAttr(c(obj), c(name), c(value))
	return int2Err(ret)
}

// SetAttrString sets the attribute of "obj" with the name "name" to "value".
// This is equivalent to the Python "obj.name = value".
func (obj *BaseObject) SetAttrString(name string, value Object) error {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyObject_SetAttrString(c(obj), s, c(value))
	return int2Err(ret)
}

// PyObject_GenericSetAttr : This is an internal helper function - we shouldn't
// need to expose it ...

// DelAttr deletes the attribute with the name "name" from "obj".  This is
// equivalent to the Python "del obj.name".
func (obj *BaseObject) DelAttr(name Object) error {
	ret := C.PyObject_SetAttr(c(obj), c(name), nil)
	return int2Err(ret)
}

// DelAttrString deletes the attribute with the name "name" from "obj".  This is
// equivalent to the Python "del obj.name".
func (obj *BaseObject) DelAttrString(name string) error {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyObject_SetAttrString(c(obj), s, nil)
	return int2Err(ret)
}

// RichCompare compares "obj" with "obj2" using the specified operation (LE, GE
// etc.), and returns the result.  The equivalent Python is "obj op obj2", where
// op is the corresponding Python operator for op.
//
// Return value: New Reference.
func (obj *BaseObject) RichCompare(obj2 Object, op Op) (Object, error) {
	ret := C.PyObject_RichCompare(c(obj), c(obj2), C.int(op))
	return obj2ObjErr(ret)
}

// RichCompare compares "obj" with "obj2" using the specified operation (LE, GE
// etc.), and returns true or false.  The equivalent Python is "obj op obj2",
// where op is the corresponding Python operator for op.
func (obj *BaseObject) RichCompareBool(obj2 Object, op Op) (bool, error) {
	ret := C.PyObject_RichCompareBool(c(obj), c(obj2), C.int(op))
	return int2BoolErr(ret)
}

// Repr returns a String representation of "obj".  This is equivalent to the
// Python "repr(obj)".
//
// Return value: New Reference.
func (obj *BaseObject) Repr() (*Unicode, error) {
	ret := C.PyObject_Repr(c(obj))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Str returns a String representation of "obj".  This is equivalent to the
// Python "str(obj)".
//
// Return value: New Reference.
func (obj *BaseObject) Str() (*Unicode, error) {
	ret := C.PyObject_Str(c(obj))
	if ret == nil {
		return nil, exception()
	}
	return newObject(ret).(*Unicode), nil
}

// Bytes returns a Bytes representation of "obj".  This is equivalent to the
// Python "bytes(obj)".
//
// Return value: New Reference.
func (obj *BaseObject) Bytes() (Object, error) {
	ret := C.PyObject_Bytes(c(obj))
	return obj2ObjErr(ret)
}

// PyObject_Unicode : TODO

// IsInstance returns true if "obj" is an instance of "cls", false otherwise.
// If "cls" is a Type instead of a class, then true will be return if "obj" is
// of that type.  If "cls" is a Tuple then true will be returned if "obj" is an
// instance of any of the Objects in the tuple.  This is equivalent to the
// Python "isinstance(obj, cls)".
func (obj *BaseObject) IsInstance(cls Object) (bool, error) {
	ret := C.PyObject_IsInstance(c(obj), c(cls))
	return int2BoolErr(ret)
}

// IsSubclass returns true if "obj" is a Subclass of "cls", false otherwise.  If
// "cls" is a Tuple, then true is returned if "obj" is a Subclass of any member
// of "cls".  This is equivalent to the Python "issubclass(obj, cls)".
func (obj *BaseObject) IsSubclass(cls Object) (bool, error) {
	ret := C.PyObject_IsSubclass(c(obj), c(cls))
	return int2BoolErr(ret)
}

// Call calls obj with the given args and kwds. kwds may be nil, args may not
// (an empty Tuple must be used if no arguments are wanted).  Returns the result
// of the call, or an Error on failure.  This is equivalent to
// "obj(*args, **kwds)" in Python.
//
// Return value: New Reference.
func (obj *BaseObject) Call(args *Tuple, kwds *Dict) (Object, error) {
	ret := C.PyObject_Call(c(obj), c(args), c(kwds))
	return obj2ObjErr(ret)
}

// Convenience types for using things like CallGo
type (
	A []any
	K map[string]any
)

// CallGo calls obj with the given args and kwds, either may be nil. Returns the
// result of the call, or an Error on failure.  This is equivalent to
// "obj(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (obj *BaseObject) CallGo(args []any, kwds map[string]any) (Object, error) {
	obj1, err := NewTupleFromValues(args...)
	if err != nil {
		return nil, err
	}
	defer obj1.Decref()
	obj2, err := NewDictFromValuesString(kwds)
	if err != nil {
		return nil, err
	}
	defer obj2.Decref()
	ret := C.PyObject_Call(c(obj), c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// CallGoArgs calls obj with the given args. Returns the result of the call, or
// an Error on failure.  This is equivalent to "obj(*args)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (obj *BaseObject) CallGoArgs(args ...any) (Object, error) {
	return obj.CallGo(args, nil)
}

// CallObject calls obj with the given args.  args may be nil.  Returns the
// result of the call, or an Error on failure.  This is equivalent to
// "obj(*args)" in Python.
//
// Return value: New Reference.
func (obj *BaseObject) CallObject(args *Tuple) (Object, error) {
	ret := C.PyObject_CallObject(c(obj), c(args))
	return obj2ObjErr(ret)
}

func (obj *BaseObject) CallFunction(format string, args ...interface{}) (Object, error) {
	t, err := buildTuple(format, args...)
	if err != nil {
		return nil, err
	}
	defer t.Decref()
	return obj.CallObject(t)
}

func (obj *BaseObject) CallMethod(name string, format string, args ...interface{}) (Object, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	f := C.PyObject_GetAttrString(c(obj), cname)
	if f == nil {
		return nil, AttributeError.Err(name)
	}
	defer C.decref(f)

	if C.PyCallable_Check(f) == 0 {
		return nil, TypeError.Err("attribute of type '%s' is not callable", name)
	}

	t, err := buildTuple(format, args...)
	if err != nil {
		return nil, err
	}
	defer t.Decref()

	ret := C.PyObject_CallObject(f, c(t))
	return obj2ObjErr(ret)
}

func (obj *BaseObject) CallFunctionObjArgs(args ...Object) (Object, error) {
	t, err := PackTuple(args...)
	if err != nil {
		return nil, err
	}
	defer t.Decref()
	return obj.CallObject(t)
}

func (obj *BaseObject) CallMethodObjArgs(name string, args ...Object) (Object, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	f := C.PyObject_GetAttrString(c(obj), cname)
	if f == nil {
		return nil, AttributeError.Err(name)
	}
	defer C.decref(f)

	if C.PyCallable_Check(f) == 0 {
		return nil, TypeError.Err("attribute of type '%s' is not callable", name)
	}

	t, err := PackTuple(args...)
	if err != nil {
		return nil, err
	}
	defer t.Decref()

	ret := C.PyObject_CallObject(f, c(t))
	return obj2ObjErr(ret)
}

// CallMethodGo calls the method of obj called name with the given args and
// kwds, either may be nil. Returns the result of the call, or an Error on
// failure. This is equivalent to "obj.name(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (obj *BaseObject) CallMethodGo(name string, args []any, kwds map[string]any) (Object, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	f := C.PyObject_GetAttrString(c(obj), cname)
	if f == nil {
		return nil, AttributeError.Err(name)
	}
	defer C.decref(f)

	if C.PyCallable_Check(f) == 0 {
		return nil, TypeError.Err("attribute of type '%s' is not callable", name)
	}

	obj1, err := NewTupleFromValues(args...)
	if err != nil {
		return nil, err
	}
	defer obj1.Decref()

	obj2, err := NewDictFromValuesString(kwds)
	if err != nil {
		return nil, err
	}
	defer obj2.Decref()

	ret := C.PyObject_Call(f, c(obj1), c(obj2))
	return obj2ObjErr(ret)
}

// CallMethodGoArgs calls the method of obj called name with the given args.
// Returns the result of the call, or an Error on failure. This is equivalent
// to "obj.name(*args)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (obj *BaseObject) CallMethodGoArgs(name string, args ...any) (Object, error) {
	return obj.CallMethodGo(name, args, nil)
}

// PyObject_Hash : TODO

// PyObject_HashNotImplement : This is an internal function, that we probably
// don't need to export.

// IsTrue returns true if the value of obj is considered to be True.  This is
// equivalent to "if obj:" in Python.
func (obj *BaseObject) IsTrue() bool {
	ret := C.PyObject_IsTrue(c(obj))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// Not returns true if the value of obj is considered to be False.  This is
// equivalent to "if not obj:" in Python.
func (obj *BaseObject) Not() bool {
	ret := C.PyObject_Not(c(obj))
	if ret < 0 {
		panic(exception())
	}
	return ret != 0
}

// PyObject_TypeCheck : TODO

// Length returns the length of the Object.  This is equivalent to the Python
// "len(obj)".
func (obj *BaseObject) Length() (int64, error) {
	ret := C.PyObject_Length(c(obj))
	return int64(ret), exception()
}

// Size returns the length of the Object.  This is equivalent to the Python
// "len(obj)".
func (obj *BaseObject) Size() (int64, error) {
	ret := C.PyObject_Size(c(obj))
	return int64(ret), exception()
}

// GetItem returns the element of "obj" corresponding to "key".  This is
// equivalent to the Python "obj[key]".
//
// Return value: New Reference.
func (obj *BaseObject) GetItem(key Object) (Object, error) {
	ret := C.PyObject_GetItem(c(obj), c(key))
	return obj2ObjErr(ret)
}

// SetItem sets the element of "obj" corresponding to "key" to "value".  This is
// equivalent to the Python "obj[key] = value".
func (obj *BaseObject) SetItem(key, value Object) error {
	ret := C.PyObject_SetItem(c(obj), c(key), c(value))
	return int2Err(ret)
}

// DelItem deletes the element from "obj" that corresponds to "key".  This is
// equivalent to the Python "del obj[key]".
func (obj *BaseObject) DelItem(key Object) error {
	ret := C.PyObject_DelItem(c(obj), c(key))
	return int2Err(ret)
}

// PyObject_AsFileDescriptor : TODO

func (obj *BaseObject) Dir() (Object, error) {
	ret := C.PyObject_Dir(c(obj))
	return obj2ObjErr(ret)
}

// PyObject_GetIter : TODO

// CopyData copies the data from src to obj. Unless both obj and src implement
// the Buffer Protocol this method will return a TypeError. If both obj and src
// implement the Buffer Protocol, then the copy may still fail with an error.
func (obj *BaseObject) CopyData(src Object) error {
	ret := C.PyObject_CopyData(c(obj), c(src))
	return int2Err(ret)
}
