package py

// #include "utils.h"
import "C"

import (
	"sync"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Object is the generic interface that represents a Python object. All of the
// concrete types satisfy the Object interface.
type Object interface {
	Base() *BaseObject
	Type() *Type
	Decref()
	Incref()
	IsTrue() bool
	Not() bool
	Free()
}

// Op is the type for the Python comparison operators. Used by the RichCompare
// functions.
type Op int

// Python comparison operators
const (
	LT = Op(C.Py_LT)
	LE = Op(C.Py_LE)
	EQ = Op(C.Py_EQ)
	NE = Op(C.Py_NE)
	GT = Op(C.Py_GT)
	GE = Op(C.Py_GE)
)

// NotImplemented is a special Object value than can be returned in some
// circumstances to indicate that a type method is not implemented.
var NotImplemented = newObject(&C._Py_NotImplementedStruct)

// RichCompareNativeBool is a helper function for implementing RichCompareBool.
// Given two comparable Go values it will compare them with the requested Op and
// return true or false. If the op is unknown, then a ValueError will be
// returned.
func RichCompareNativeBool[T constraints.Ordered](a, b T, op Op) (bool, error) {
	switch op {
	case LT:
		return a < b, nil
	case LE:
		return a <= b, nil
	case EQ:
		return a == b, nil
	case NE:
		return a != b, nil
	case GT:
		return a > b, nil
	case GE:
		return a >= b, nil
	default:
		return false, ValueError.Err("unknown compare op: %d", op)
	}
}

// RichCompareNative is a helper function for implementing RichCompare. Given
// two comparable Go values it will compare them with the requested Op and
// return True or False. If the op is unknown, then a ValueError will be
// returned.
//
// Return value: New Reference.
func RichCompareNative[T constraints.Ordered](a, b T, op Op) (Object, error) {
	ret, err := RichCompareNativeBool(a, b, op)
	if err != nil {
		return NotImplemented, nil
	}
	if ret {
		True.Incref()
		return True, nil
	}
	False.Incref()
	return False, nil
}

// NewValue will try to return an appropriate Python Object for the supplied
// value. If the type can't be converted, then a TypeError will be returned.
//
// If an Object is supplied, then a new reference to that Object will be
// returned.
//
// Return value: New Reference.
func NewValue(value any) (Object, error) {
	switch v := value.(type) {
	case Object:
		v.Incref()
		return v, nil
	case bool:
		return NewBool(v), nil
	case int:
		return NewLong(v), nil
	case int8:
		return NewLong(v), nil
	case int16:
		return NewLong(v), nil
	case int32:
		return NewLong(v), nil
	case int64:
		return NewLong(v), nil
	case uint8:
		return NewLong(v), nil
	case uint16:
		return NewLong(v), nil
	case uint32:
		return NewLong(v), nil
	case float32:
		return NewFloat(v)
	case float64:
		return NewFloat(v)
	case complex64:
		return NewComplex(v)
	case complex128:
		return NewComplex(v)
	case string:
		return NewUnicode(v)
	case []byte:
		return NewBytes(v), nil
	case []Object:
		return NewListFromObjects(v...)
	case map[Object]Object:
		return NewDictFromMap(v)
	case map[string]Object:
		return NewDictFromMapString(v)
	default:
		return nil, TypeError.Err("unsupported type %T", v)
	}
}

// None is the Python equivalent to nil.
var None = (*NoneObject)(unsafe.Pointer(&C._Py_NoneStruct))

// NoneObject is the type of the None value.  The only value of this type is
// None.
type NoneObject struct {
	BaseObject
}

func (n *NoneObject) String() string {
	return "None"
}

// Decref decrements obj's reference count, obj may be nil.
func Decref(obj Object) {
	if obj != nil {
		C.decref(c(obj))
	}
}

// Incref increments obj's reference count, obj may be nil.
func Incref(obj Object) {
	if obj != nil {
		C.incref(c(obj))
	}
}

func c(obj Object) *C.PyObject {
	if obj == nil {
		return nil
	}
	return (*C.PyObject)(unsafe.Pointer(obj.Base()))
}

func ce(obj Object, err error) *C.PyObject {
	if err != nil {
		raise(err)
		return nil
	}
	return c(obj)
}

// free deallocates the storage (in Python) for obj.  After calling this method,
// obj should no longer be used.
func free(obj Object) {
	o := c(obj)

	// This can happen if a PyDealloc method on a ClassObject calls Free
	if o == nil {
		return
	}

	// Make sure this instance isn't registered any more
	clearClassObject(unsafe.Pointer(o))

	// Call Python free function
	pyType := (*C.PyTypeObject)(unsafe.Pointer(o.ob_type))
	C.typeFree(pyType, o)
}

var (
	typeLock sync.RWMutex
	types    = make(map[*C.PyTypeObject]*Class)
)

func registerClass(pyType *C.PyTypeObject, class *Class) {
	typeLock.Lock()
	defer typeLock.Unlock()

	types[pyType] = class
}

func getClass(pyType *C.PyTypeObject) *Class {
	typeLock.RLock()
	defer typeLock.RUnlock()

	return types[pyType]
}

func newObject(obj *C.PyObject) Object {
	if obj == nil {
		return nil
	}

	o := unsafe.Pointer(obj)
	if o == unsafe.Pointer(None) {
		return None
	}

	if cObj := getClassObject(o); cObj != nil {
		return cObj
	}

	if class := getClass((*C.PyTypeObject)(o)); class != nil {
		return class
	}

	if C.exceptionClassCheck(obj) != 0 {
		return newException(obj)
	}

	if natObj := getNativeType(obj); natObj != nil {
		return natObj
	}

	return newBaseObject(obj)
}
