// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

const (
	TPFLAGS_HAVE_GC  = uint32(C.Py_TPFLAGS_HAVE_GC)
	TPFLAGS_BASETYPE = uint32(C.Py_TPFLAGS_BASETYPE)
)

// A Class struct instance is used to define a Python class that has been
// implemented in Go.
//
// Name should be the name of the type in Python, including the package name,
// e.g. "test.MyClass"
//
// Flags is or'ed with Py_TPFLAGS_DEFAULT and passed through to the tp_flags
// member
//
// Doc is currently unused.
//
// Type holds a pointer to the Type instance for this class, this is filled in
// by calling Create().
//
// Object should be set to a pointer of the struct type that will represent an
// instance of the Python class.  This struct must contain an embedded
// py.ClassBaseObject.  The easiest ways to set Object are either to use a
// struct literal (i.e. &MyClass{}), or to cast nil (i.e. (*MyClass)(nil)), if
// the struct is large then the latter method is more efficient (as an instance
// of the struct is not created).
//
// This struct may have the following special methods (the equivalent Python
// methods are also indicated):
//
//	PyInit(args *py.Tuple, kwds *py.Dict) os.Error              // __init__
//	PyCall(args *py.Tuple, kwds *py.Dict) (py.Object, os.Error) // __call__
//	PyRepr() string                                             // __repr__
//	PyStr() string                                              // __str__
//	PyCompare(obj py.Object) (int, os.Error)                    // __cmp__
//
// Properties are also supported, by implementing get and set methods:
//
//	PyGet_XXX() (py.Object, os.Error)
//	PySet_XXX(value py.Object) os.Error
//
// Instance methods on the Python class are implemented by methods on the
// struct type with the Py_ prefix:
//
//	Py_ABC(args *py.Tuple) (py.Object, os.Error)
//	Py_XYX(args *py.Tuple, kwds *py.Dict) (py.Object, os.Error)
//
// NOTE: All of the methods referred to above should use a pointer receiver.
//
// Static defines static methods for the Python class.
type Class struct {
	Name   string
	Flags  uint32
	Doc    string
	Type   *Type
	Object ClassObject
	Static map[string]any
	New    func(*Class, *Tuple, *Dict) (ClassObject, error)
}

func (c *Class) newObject(args *Tuple, kwds *Dict) (ClassObject, error) {
	// simple case, a New method has been provided
	if c.New != nil {
		return c.New(c, args, kwds)
	}

	// no New provided, so we need to create an instance of the correct type
	t := reflect.TypeOf(c.Object).Elem()
	v := reflect.New(t)
	return v.Interface().(ClassObject), nil
}

var otyp = reflect.TypeOf(new(Object)).Elem()

func getMethodAndObject(obj unsafe.Pointer) (unsafe.Pointer, ClassObject, error) {
	t := (*C.PyObject)(obj)
	pyobj := unsafe.Pointer(C.PyTuple_GetItem(t, 0))
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 1), nil)

	o := getClassObject(pyobj)
	if o == nil {
		return nil, nil, fmt.Errorf("unknown object")
	}

	return m, o, nil
}

//export goClassCallMethod
func goClassCallMethod(obj, unused unsafe.Pointer) unsafe.Pointer {
	// Unpack context and self pointer from obj
	m, o, err := getMethodAndObject(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := (*func(p ClassObject) (Object, error))(unsafe.Pointer(&m))

	ret, err := (*f)(o)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallMethodArgs
func goClassCallMethodArgs(obj, args unsafe.Pointer) unsafe.Pointer {
	// Unpack context and self pointer from obj
	m, o, err := getMethodAndObject(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args ready to use, by turning it into a pointer of the appropriate
	// type
	a := newTuple((*C.PyObject)(args))

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := (*func(o ClassObject, a *Tuple) (Object, error))(unsafe.Pointer(&m))

	ret, err := (*f)(o, a)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallMethodKwds
func goClassCallMethodKwds(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Unpack context and self pointer from obj
	m, o, err := getMethodAndObject(obj)
	if err != nil {
		raise(err)
		return nil
	}

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := (*func(o ClassObject, a *Tuple, k *Dict) (Object, error))(unsafe.Pointer(&m))

	ret, err := (*f)(o, a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

func getStaticMethod(obj unsafe.Pointer) unsafe.Pointer {
	t := (*C.PyObject)(obj)
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 0), nil)
	return m
}

//export goClassCallStaticMethod
func goClassCallStaticMethod(obj, unused unsafe.Pointer) unsafe.Pointer {
	log.Printf("STATIC CALL NO ARGS: %p %p", obj, unused)

	// Unpack function from obj
	m := getStaticMethod(obj)

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := (*func() (Object, error))(unsafe.Pointer(&m))

	ret, err := (*f)()
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallStaticMethodArgs
func goClassCallStaticMethodArgs(obj, args unsafe.Pointer) unsafe.Pointer {
	log.Printf("STATIC CALL ARGS: %p %p", obj, args)

	// Unpack function from obj
	m := getStaticMethod(obj)

	// Get args ready to use, by turning it into a pointer of the appropriate
	// type
	a := newTuple((*C.PyObject)(args))

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := (*func(a *Tuple) (Object, error))(unsafe.Pointer(&m))

	ret, err := (*f)(a)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassCallStaticMethodKwds
func goClassCallStaticMethodKwds(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	log.Printf("STATIC CALL KWDS: %p %p", obj, args)

	// Unpack function from obj
	m := getStaticMethod(obj)

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	// Now call the actual static method by pulling the method out of the
	// reflect.Type object stored in the context
	f := (*func(a *Tuple, k *Dict) (Object, error))(unsafe.Pointer(&m))

	ret, err := (*f)(a, k)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassSetProp
func goClassSetProp(obj, arg, closure unsafe.Pointer) int {
	// Unpack set function from closure
	t := (*C.PyObject)(closure)
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 1), nil)

	// Turn arg into something usable
	a := newObject((*C.PyObject)(arg))

	// Turn the function into something we can call
	f := (*func(p unsafe.Pointer, a Object) error)(unsafe.Pointer(&m))

	err := (*f)(obj, a)
	if err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export goClassGetProp
func goClassGetProp(obj, closure unsafe.Pointer) unsafe.Pointer {
	// Unpack get function from closure
	t := (*C.PyObject)(closure)
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 0), nil)

	// Turn the function into something we can call
	f := (*func(p unsafe.Pointer) (Object, error))(unsafe.Pointer(&m))

	ret, err := (*f)(obj)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassObjGet
func goClassObjGet(obj unsafe.Pointer, idx int) unsafe.Pointer {
	field := getField(idx)
	item := unsafe.Pointer(uintptr(obj) + field.Offset)

	var o Object

	if field.Type == otyp {
		o = *(*Object)(item)
	} else {
		o = reflect.NewAt(field.Type, item).Elem().Interface().(Object)
	}

	o.Incref()
	return unsafe.Pointer(c(o))
}

//export goClassObjSet
func goClassObjSet(obj unsafe.Pointer, idx int, obj2 unsafe.Pointer) int {
	field := getField(idx)
	item := unsafe.Pointer(uintptr(obj) + field.Offset)

	// This is the new value we are being asked to set
	value := newObject((*C.PyObject)(obj2))

	// Special case for Object fields, we don't need reflect for these.  We have
	// to be careful with refcounts, as decref could invoke desctructor code
	// etc.
	if field.Type == otyp {
		o := (*Object)(item)
		tmp := *o
		Incref(value)
		*o = value
		Decref(tmp)
		return 0
	}

	vt := reflect.TypeOf(value)
	ov := reflect.NewAt(field.Type, unsafe.Pointer(item)).Elem()

	// If the value is assignable to the field, then we do it, with the same
	// refcount dance as above.
	if vt.AssignableTo(field.Type) {
		tmp := ov.Interface().(Object)
		Incref(value)
		ov.Set(reflect.ValueOf(value))
		Decref(tmp)
		return 0
	}

	// The given value wasn't assignable to the field - raise an error
	tn := ov.Type().Elem().Name()
	raise(TypeError.Err("Cannot assign '%T' to '*%v'", value, tn))
	return -1
}

//export goClassNatGet
func goClassNatGet(obj unsafe.Pointer, idx int) unsafe.Pointer {
	field := getField(idx)
	item := unsafe.Pointer(uintptr(obj) + field.Offset)

	switch field.Type.Kind() {
	case reflect.Int:
		i := (*int)(item)
		return unsafe.Pointer(C.PyLong_FromLong(C.long(*i)))
	}

	raise(NotImplementedError.ErrV(None))
	return nil
}

//export goClassNatSet
func goClassNatSet(obj unsafe.Pointer, idx int, obj2 unsafe.Pointer) int {
	field := getField(idx)
	item := unsafe.Pointer(uintptr(obj) + field.Offset)

	// This is the new value we are being asked to set
	value := newObject((*C.PyObject)(obj2))

	switch field.Type.Kind() {
	case reflect.Int:
		v := int(C.PyLong_AsLong(c(value)))
		if exceptionRaised() {
			return -1
		}
		i := (*int)(item)
		*i = v
		return 0
	}

	raise(NotImplementedError.ErrV(None))
	return -1
}

//export goClassTraverse
func goClassTraverse(obj, visit, arg unsafe.Pointer) int {
	// Get the Python type object
	pyType := (*C.PyTypeObject)((*C.PyObject)(obj).ob_type)

	class, ok := getType(pyType)
	if !ok {
		t := newType((*C.PyObject)(unsafe.Pointer(pyType)))
		raise(TypeError.Err("Not a recognised type: %s", t))
		return -1
	}

	st := reflect.TypeOf(class.Object).Elem()

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if !field.Type.AssignableTo(otyp) {
			continue
		}
		v := unsafe.Pointer(uintptr(obj) + field.Offset)
		var o Object
		if field.Type == otyp {
			o = *(*Object)(v)
		} else {
			o = *(**AbstractObject)(v)
		}
		ret := C.doVisit(c(o), visit, arg)
		if ret != 0 {
			return int(ret)
		}
	}

	return 0
}

//export goClassClear
func goClassClear(obj unsafe.Pointer) int {
	// Get the Python type object
	pyType := (*C.PyTypeObject)((*C.PyObject)(obj).ob_type)

	class, ok := getType(pyType)
	if !ok {
		t := newType((*C.PyObject)(unsafe.Pointer(pyType)))
		raise(TypeError.Err("Not a recognised type: %s", t))
		return -1
	}

	st := reflect.TypeOf(class.Object).Elem()

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if !field.Type.AssignableTo(otyp) {
			continue
		}
		v := unsafe.Pointer(uintptr(obj) + field.Offset)
		if field.Type == otyp {
			o := (*Object)(v)
			tmp := *o
			*o = nil
			Decref(tmp)
		} else {
			o := (**AbstractObject)(v)
			tmp := *o
			*o = nil
			Decref(tmp)
		}
	}

	return 0
}

var (
	ctxtLock sync.RWMutex
	contexts = map[uintptr]*C.ClassContext{}
)

func setClassContext(obj unsafe.Pointer, pyType *C.PyTypeObject) {
	ctxtLock.Lock()
	defer ctxtLock.Unlock()

	ctxt := (*C.ClassContext)(unsafe.Pointer(pyType.tp_methods))
	contexts[uintptr(obj)] = ctxt
}

func clearClassContext(obj unsafe.Pointer) {
	ctxtLock.Lock()
	defer ctxtLock.Unlock()

	delete(contexts, uintptr(obj))
}

func getClassContext(obj unsafe.Pointer) *C.ClassContext {
	ctxtLock.RLock()
	defer ctxtLock.RUnlock()

	ctxt := contexts[uintptr(obj)]
	if ctxt == nil {
		panic("Asked for context of unregistered object!")
	}
	return ctxt
}

//export goClassNew
func goClassNew(typ, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Get the Python type object
	pyType := (*C.PyTypeObject)(typ)

	class, _ := getType(pyType)

	for class == nil && pyType.tp_base != nil {
		pyType = (*C.PyTypeObject)(unsafe.Pointer(pyType.tp_base))
		class, _ = getType(pyType)
	}

	if class == nil {
		t := newType((*C.PyObject)(unsafe.Pointer(pyType)))
		raise(TypeError.Err("Not a recognised type: %s", t))
		return nil
	}

	// Get typ ready to use by turning into *Type
	t := newType((*C.PyObject)(typ))

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	// allocate the go object
	goObj, err := class.newObject(a, k)
	if err != nil {
		raise(err)
		return nil
	}

	// allocate the Python proxy object
	pyObj := unsafe.Pointer(C.typeAlloc(c(t), C.Py_ssize_t(0)))
	if pyObj == nil {
		return nil
	}

	// finalise the setup of the go object
	goObj.setBase((*BaseObject)(pyObj))
	registerClassObject(pyObj, goObj)

	// register class context against new object
	setClassContext(pyObj, pyType)

	return pyObj
}

type prop struct {
	get unsafe.Pointer
	set unsafe.Pointer
}

type method struct {
	f     unsafe.Pointer
	flags int
}

func methSigMatches(got reflect.Type, _want interface{}) error {
	// Note: Methods take the receiver as the first argument, which the want
	// signature doesn't include.

	if got.NumIn() == 0 {
		// The receiver is missing!
		return fmt.Errorf("Method without reciever!")
	}

	want := reflect.TypeOf(_want)

	if got.NumIn()-1 != want.NumIn() {
		return fmt.Errorf("Method should have %d arguments, not %d", want.NumIn(), got.NumIn()-1)
	}

	if got.NumOut() != want.NumOut() {
		return fmt.Errorf("Method should have %d return values, not %d", want.NumOut(), got.NumOut())
	}

	for i := 0; i < want.NumIn(); i++ {
		if got.In(i+1) != want.In(i) {
			return fmt.Errorf("Method argument %d should be %v, not %v", i+1, want.In(i), got.In(i+1))
		}
	}

	for i := 0; i < want.NumOut(); i++ {
		if got.Out(i) != want.Out(i) {
			return fmt.Errorf("Method return value %d should be %v, not %v", i+1, want.Out(i), got.Out(i))
		}
	}

	return nil
}

func getPythonCallFlags(f reflect.Type) (int, error) {
	log.Printf("GET FLAGS: %T", reflect.New(f).Interface())
	switch {
	case methSigMatches(f, pyUnaryFunc) == nil:
		return C.METH_NOARGS, nil
	case methSigMatches(f, pyBinaryCallFunc) == nil:
		return C.METH_VARARGS, nil
	case methSigMatches(f, pyTernaryCallFunc) == nil:
		return C.METH_VARARGS | C.METH_KEYWORDS, nil
	default:
		return 0, fmt.Errorf("invalid method signature")
	}
}

func funcSigMatches(got reflect.Type, _want interface{}) error {
	want := reflect.TypeOf(_want)

	if got.NumIn() != want.NumIn() {
		return fmt.Errorf("function should have %d arguments, not %d", want.NumIn(), got.NumIn())
	}

	if got.NumOut() != want.NumOut() {
		return fmt.Errorf("function should have %d return values, not %d", want.NumOut(), got.NumOut())
	}

	for i := 0; i < want.NumIn(); i++ {
		if got.In(i) != want.In(i) {
			return fmt.Errorf("function argument %d should be %v, not %v", i+1, want.In(i), got.In(i))
		}
	}

	for i := 0; i < want.NumOut(); i++ {
		if got.Out(i) != want.Out(i) {
			return fmt.Errorf("function return value %d should be %v, not %v", i+1, want.Out(i), got.Out(i))
		}
	}

	return nil
}

func getStaticCallFlags(f reflect.Type) (int, error) {
	log.Printf("GET STATIC FLAGS: %T", reflect.New(f).Interface())
	switch {
	case funcSigMatches(f, pyUnaryFunc) == nil:
		return C.METH_NOARGS, nil
	case funcSigMatches(f, pyBinaryCallFunc) == nil:
		return C.METH_VARARGS, nil
	case funcSigMatches(f, pyTernaryCallFunc) == nil:
		return C.METH_VARARGS | C.METH_KEYWORDS, nil
	default:
		return 0, fmt.Errorf("invalid function signature")
	}
}

// Alloc is a convenience function, so that Go code can create a new Object
// instance.
func (class *Class) Alloc(n int64) (obj Object, err error) {
	obj, err = class.Type.Alloc(n)

	// Since we are creating this object for Go code, this is probably the only
	// opportunity we will get to register this object instance.
	pyType := (*C.PyTypeObject)(unsafe.Pointer(c(class.Type)))
	setClassContext(unsafe.Pointer(c(obj)), pyType)

	return
}

func Clear(obj Object) error {
	ret := goClassClear(unsafe.Pointer(c(obj)))
	if ret < 0 {
		return exception()
	}
	return nil
}

var (
	fieldLock sync.RWMutex
	fields    []reflect.StructField
)

func registerField(field reflect.StructField) C.int {
	fieldLock.Lock()
	defer fieldLock.Unlock()

	fields = append(fields, field)
	return C.int(len(fields) - 1)
}

func getField(idx int) reflect.StructField {
	fieldLock.RLock()
	defer fieldLock.RUnlock()

	if idx >= len(fields) {
		panic("Request for unregistered field!")
	}

	return fields[idx]
}

var exportable = map[reflect.Kind]bool{
	reflect.Bool:    true,
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint32:  true,
	reflect.Uint64:  true,
	reflect.Uintptr: true,
	reflect.Float32: true,
	reflect.Float64: true,
}

type methodSignature struct {
	field string
	sig   interface{}
}

// Function signatures for methods that implement Python methods.  Note, the
// first argument is always the receiver, and is not included in these
// signatures (hence the names are one greater than the number of arguments
// taken).
var (
	pyInitFunc        = (func(*Tuple, *Dict) error)(nil)
	pyVoidFunc        = (func())(nil)
	pyReprFunc        = (func() string)(nil)
	pyLenFunc         = (func() int64)(nil)
	pyHashFunc        = (func() (uint32, error))(nil)
	pyInquiryFunc     = (func() (bool, error))(nil)
	pyUnaryFunc       = (func() (Object, error))(nil)
	pyBinaryFunc      = (func(Object) (Object, error))(nil)
	pyTernaryFunc     = (func(a, b Object) (Object, error))(nil)
	pyBinaryCallFunc  = (func(*Tuple) (Object, error))(nil)
	pyTernaryCallFunc = (func(*Tuple, *Dict) (Object, error))(nil)
	pyCompareFunc     = (func(Object) (int, error))(nil)
	pyRichCmpFunc     = (func(Object, Op) (Object, error))(nil)
	pyObjObjArgFunc   = (func(a, b Object) error)(nil)
	pySsizeArgFunc    = (func(int64) (Object, error))(nil)
	pySsizeObjArgFunc = (func(int64, Object) error)(nil)
	pyObjObjFunc      = (func(Object) (bool, error))(nil)
	pyGetAttrFunc     = (func(string) (Object, error))(nil)
	pyGetAttrObjFunc  = (func(Object) (Object, error))(nil)
	pySetAttrFunc     = (func(string, Object) error)(nil)
	pySetAttrObjFunc  = (func(Object, Object) error)(nil)
)

var methodMap = map[string]methodSignature{
	// Standard Methods
	"PyCall":        {"call", pyTernaryCallFunc},
	"PyCompare":     {"compare", pyCompareFunc},
	"PyDealloc":     {"dealloc", pyVoidFunc},
	"PyGetAttr":     {"getattr", pyGetAttrFunc},
	"PyGetAttrObj":  {"getattro", pyGetAttrObjFunc},
	"PyHash":        {"hash", pyHashFunc},
	"PyInit":        {"init", pyInitFunc},
	"PyIter":        {"iter", pyUnaryFunc},
	"PyIterNext":    {"iternext", pyUnaryFunc},
	"PyRepr":        {"repr", pyReprFunc},
	"PyRichCompare": {"richcmp", pyRichCmpFunc},
	"PySetAttr":     {"setattr", pySetAttrFunc},
	"PySetAttrObj":  {"setattro", pySetAttrObjFunc},
	"PyStr":         {"str", pyReprFunc},

	// Mapping Protocol
	"PyMapLen": {"mp_len", pyLenFunc},
	"PyMapGet": {"mp_get", pyBinaryFunc},
	"PyMapSet": {"mp_set", pyObjObjArgFunc},

	// Number Protocol
	"PyNumAdd":         {"nb_add", pyBinaryFunc},
	"PyNumSubtract":    {"nb_subtract", pyBinaryFunc},
	"PyNumMultiply":    {"nb_multiply", pyBinaryFunc},
	"PyNumDivide":      {"nb_divide", pyBinaryFunc},
	"PyNumRemainder":   {"nb_remainder", pyBinaryFunc},
	"PyNumDivmod":      {"nb_divmod", pyBinaryFunc},
	"PyNumPower":       {"nb_power", pyTernaryFunc},
	"PyNumNegative":    {"nb_negative", pyUnaryFunc},
	"PyNumPositive":    {"nb_positive", pyUnaryFunc},
	"PyNumAbsolute":    {"nb_absolute", pyUnaryFunc},
	"PyNumNonzero":     {"nb_nonzero", pyInquiryFunc},
	"PyNumInvert":      {"nb_invert", pyUnaryFunc},
	"PyNumLshift":      {"nb_lshift", pyBinaryFunc},
	"PyNumRshift":      {"nb_rshift", pyBinaryFunc},
	"PyNumAnd":         {"nb_and", pyBinaryFunc},
	"PyNumXor":         {"nb_xor", pyBinaryFunc},
	"PyNumOr":          {"nb_or", pyBinaryFunc},
	"PyNumInt":         {"nb_int", pyUnaryFunc},
	"PyNumLong":        {"nb_long", pyUnaryFunc},
	"PyNumFloat":       {"nb_float", pyUnaryFunc},
	"PyNumOct":         {"nb_oct", pyUnaryFunc},
	"PyNumHex":         {"nb_hex", pyUnaryFunc},
	"PyNumIpAdd":       {"nb_ip_add", pyBinaryFunc},
	"PyNumIpSubtract":  {"nb_ip_subtract", pyBinaryFunc},
	"PyNumIpMultiply":  {"nb_ip_multiply", pyBinaryFunc},
	"PyNumIpDivide":    {"nb_ip_divide", pyBinaryFunc},
	"PyNumIpRemainder": {"nb_ip_remainder", pyBinaryFunc},
	"PyNumIpPower":     {"nb_ip_power", pyTernaryFunc},
	"PyNumIpLshift":    {"nb_ip_lshift", pyBinaryFunc},
	"PyNumIpRshift":    {"nb_ip_rshift", pyBinaryFunc},
	"PyNumIpAnd":       {"nb_ip_and", pyBinaryFunc},
	"PyNumIpXor":       {"nb_ip_xor", pyBinaryFunc},
	"PyNumIpOr":        {"nb_ip_or", pyBinaryFunc},
	"PyNumFloorDiv":    {"nb_floordiv", pyBinaryFunc},
	"PyNumTrueDiv":     {"nb_truediv", pyBinaryFunc},
	"PyNumIpFloorDiv":  {"nb_ip_floordiv", pyBinaryFunc},
	"PyNumIpTrueDiv":   {"nb_ip_truediv", pyBinaryFunc},
	"PyNumIndex":       {"nb_index", pyUnaryFunc},

	// Sequence Protocol
	"PySeqLen":      {"sq_length", pyLenFunc},
	"PySeqConcat":   {"sq_concat", pyBinaryFunc},
	"PySeqRepeat":   {"sq_repeat", pySsizeArgFunc},
	"PySeqGet":      {"sq_get", pySsizeArgFunc},
	"PySeqSet":      {"sq_set", pySsizeObjArgFunc},
	"PySeqContains": {"sq_contains", pyObjObjFunc},
	"PySeqIpConcat": {"sq_ip_concat", pyBinaryFunc},
	"PySeqIpRepeat": {"sq_ip_repeat", pySsizeArgFunc},
}

func ctxtSet(ctxt *C.ClassContext, name string, fn unsafe.Pointer) {
	t := reflect.TypeOf(ctxt).Elem()
	f, ok := t.FieldByName(name)
	if !ok {
		panic("Tried to set a non-existant context field")
	}
	base := uintptr(unsafe.Pointer(ctxt))
	*(*unsafe.Pointer)(unsafe.Pointer(base + f.Offset)) = fn
	parts := strings.Split(name, "_")
	if len(parts) > 1 {
		hf, ok := t.FieldByName("has_" + parts[0])
		if ok {
			*(*int)(unsafe.Pointer(base + hf.Offset)) = 1
		}
	}
}

var (
	directFnCall = (*bool)(nil)
	indirections = make([]*unsafe.Pointer, 0, 100)
)

func methodAsPointer(m reflect.Method) unsafe.Pointer {
	// Go 1.0 uses direct function calls, Go 1.1 uses indirect function calls
	// (so that contextual data for closures can be held).  Figure out which
	// this is, and set directFnCall as appropriate.
	if directFnCall == nil {
		m := unsafe.Pointer(reflect.ValueOf(methodAsPointer).Pointer())
		f := (*func(reflect.Method) unsafe.Pointer)(unsafe.Pointer(&m))
		direct := fmt.Sprintf("%v", m) == fmt.Sprintf("%v", *f)
		directFnCall = &direct
	}

	fp := unsafe.Pointer(m.Func.Pointer())
	if *directFnCall {
		return fp
	}
	ifp := &fp
	indirections = append(indirections, ifp)
	return unsafe.Pointer(ifp)
}

func funcAsPointer(v reflect.Value) unsafe.Pointer {
	fp := unsafe.Pointer(v.Pointer())
	ifp := &fp
	indirections = append(indirections, ifp)
	return unsafe.Pointer(ifp)
}

var typeMap = map[string]*Type{
	"Bool":   BoolType,
	"Code":   CodeType,
	"Dict":   DictType,
	"Float":  FloatType,
	"List":   ListType,
	"Long":   LongType,
	"Module": ModuleType,
	// TODO(jp3): replace with bytes
	// "String": StringType,
	"Tuple": TupleType,
	"Type":  TypeType,
}

// Create creates and returns a pointer to a PyTypeObject that is the Python
// representation of the class that has been implemented in Go.
func (c *Class) Create() (*Type, error) {
	pyType := C.newType()
	pyType.tp_name = C.CString(c.Name)
	pyType.tp_flags = C.Py_TPFLAGS_DEFAULT | C.ulong(c.Flags)

	if c.Object == nil {
		c.Object = &ClassBaseObject{}
	}

	typ := reflect.TypeOf(c.Object)
	btyp := typ.Elem()

	if btyp.NumField() == 0 {
		return nil, fmt.Errorf("%s does not embed an Object", btyp.Name())
	}

	// Get a new context structure
	ctxt := C.newContext()

	methods := make(map[string]method)
	props := make(map[string]prop)

	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if !strings.HasPrefix(m.Name, "Py") {
			continue
		}
		t := m.Func.Type()
		f := methodAsPointer(m)
		fn := fmt.Sprintf("%s.%s", typ.Elem().Name(), m.Name)
		meth, ok := methodMap[m.Name]
		if ok {
			err := methSigMatches(t, meth.sig)
			if err != nil {
				C.free(unsafe.Pointer(pyType))
				return nil, fmt.Errorf("%s: %s", fn, err)
			}
			ctxtSet(ctxt, meth.field, f)
			continue
		}
		parts := strings.SplitN(m.Name, "_", 2)
		switch parts[0] {
		case "Py":
			flags, err := getPythonCallFlags(t)
			if err != nil {
				return nil, fmt.Errorf("%s: %s", fn, err)
			}
			methods[parts[1]] = method{f, flags}
		case "PySet":
			if err := methSigMatches(t, (func(Object) error)(nil)); err != nil {
				C.free(unsafe.Pointer(pyType))
				return nil, fmt.Errorf("%s: %s", fn, err)
			}
			p := props[parts[1]]
			p.set = f
			props[parts[1]] = p
		case "PyGet":
			if err := methSigMatches(t, (func() (Object, error))(nil)); err != nil {
				C.free(unsafe.Pointer(pyType))
				return nil, fmt.Errorf("%s: %s", fn, err)
			}
			p := props[parts[1]]
			p.get = f
			props[parts[1]] = p
		}
	}

	log.Printf("CLASS: setup %d static methods", len(c.Static))

	for name, fn := range c.Static {
		f := reflect.ValueOf(fn)
		t := f.Type()
		flags, err := getStaticCallFlags(t)
		if err != nil {
			return nil, fmt.Errorf("static %s: %s", name, err)
		}
		log.Printf("STATIC: %s %d", name, flags)
		methods[name] = method{funcAsPointer(f), flags | C.METH_STATIC}
	}

	pyType.tp_basicsize = C.Py_ssize_t(unsafe.Sizeof(C.PyObject{}))

	C.setClassContext(pyType, ctxt)

	if C.typeReady(pyType) < 0 {
		C.free(unsafe.Pointer(ctxt))
		C.free(unsafe.Pointer(pyType.tp_name))
		C.free(unsafe.Pointer(pyType))
		return nil, exception()
	}

	C.storeContext(pyType, ctxt)

	for name, method := range methods {
		s := C.CString(name)
		C.setTypeAttr(pyType, s, C.newMethod(s, method.f, C.int(method.flags)))
	}

	for name, prop := range props {
		s := C.CString(name)
		C.setTypeAttr(pyType, s, C.newProperty(pyType, s, prop.get, prop.set))
	}

	for i := 0; i < btyp.NumField(); i++ {
		field := btyp.Field(i)
		pyname := field.Tag.Get("Py")
		pydoc := field.Tag.Get("PyDoc")
		if pyname == "" && pydoc == "" {
			continue
		}
		if pyname == "" {
			pyname = field.Name
		}
		if field.Type.Implements(otyp) {
			// field is some type of object, so we can use the generic object
			// member get/set code.
			s := C.CString(pyname)
			defer C.free(unsafe.Pointer(s))
			d := C.CString(pydoc)
			C.setTypeAttr(pyType, s, C.newObjMember(registerField(field), d))
			continue
		}
		if !exportable[field.Type.Kind()] {
			C.free(unsafe.Pointer(ctxt))
			C.free(unsafe.Pointer(pyType.tp_name))
			C.free(unsafe.Pointer(pyType))
			return nil, fmt.Errorf("Cannot export %s.%s to Python: type '%s' unsupported", btyp.Name(), field.Name, field.Type.Name())
		}
		s := C.CString(pyname)
		defer C.free(unsafe.Pointer(s))
		d := C.CString(pydoc)
		C.setTypeAttr(pyType, s, C.newNatMember(registerField(field), d))
	}

	c.Type = newType((*C.PyObject)(unsafe.Pointer(pyType)))

	registerType(pyType, c)

	log.Printf("CREATE DONE")

	return c.Type, nil
}
