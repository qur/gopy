package py

// #include "utils.h"
import "C"

import (
	"sync"
	"unsafe"

	"golang.org/x/exp/constraints"
)

type Op int

const (
	LT = Op(C.Py_LT)
	LE = Op(C.Py_LE)
	EQ = Op(C.Py_EQ)
	NE = Op(C.Py_NE)
	GT = Op(C.Py_GT)
	GE = Op(C.Py_GE)
)

var NotImplemented = newObject(&C._Py_NotImplementedStruct)

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

// Object is the generic interface that represents a Python object.  All of the
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

func registerType(pyType *C.PyTypeObject, class *Class) {
	typeLock.Lock()
	defer typeLock.Unlock()

	types[pyType] = class
}

func getType(pyType *C.PyTypeObject) (*Class, bool) {
	typeLock.RLock()
	defer typeLock.RUnlock()

	class, ok := types[pyType]

	return class, ok
}

func newObject(obj *C.PyObject) Object {
	if obj == nil {
		return nil
	}

	o := unsafe.Pointer(obj)
	if o == unsafe.Pointer(None) {
		return None
	}

	cObj := getClassObject(o)
	if cObj != nil {
		return cObj
	}

	class, ok := getType((*C.PyTypeObject)(o))
	if ok {
		return class
	}

	if C.exceptionClassCheck(obj) != 0 {
		return newException(obj)
	}

	natObj := getNativeType(obj)
	if natObj != nil {
		return natObj
	}

	return newBaseObject(obj)
}
