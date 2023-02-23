package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"reflect"
	"strings"
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
// Once Create has been called, then Class is a valid Object that maps to the
// Python type instance for this class. Calling RawType will return the same
// Python instance as a Type.
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
//	PyInit(args *py.Tuple, kwds *py.Dict) error                 // __init__
//	PyCall(args *py.Tuple, kwds *py.Dict) (py.Object, error)    // __call__
//	PyRepr() string                                             // __repr__
//	PyStr() string                                              // __str__
//	PyRichCompare(obj py.Object, op py.Op) (py.Object, error)   // __cmp__
//
// If control over the deallocation process is desired, then the struct can
// implement:
//
//	PyDealloc()
//
// This method will then be called when the instance is being deallocated. If
// the method is implemented then it takes responsibility for clearing any
// contained Objects. The py.Clear and py.ClearClassObject functions can be used
// to assist.
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
//
// Class defines class methods for the Python class. In Go these will be
// functions that take a *Class as the first argument.
//
// New is an optional constructor for the instance type. The type returned
// should match that of Object. If New is not provided then a new instance of
// Object's type will be created.
//
// User is not used by the library code. It can be used to store state for the
// class methods, etc.
type Class struct {
	Name   string
	Flags  uint32
	Doc    string
	Object ClassObject
	Static map[string]any
	Class  map[string]any
	New    func(*Class, *Tuple, *Dict) (ClassObject, error)
	User   any
	base   *Type
}

func (cls *Class) Base() *BaseObject {
	return cls.base.Base()
}

func (cls *Class) Type() *Type {
	return cls.base.Type()
}

func (cls *Class) RawType() *Type {
	return cls.base
}

func (cls *Class) Decref() {
	cls.base.Decref()
}

func (cls *Class) Incref() {
	cls.base.Incref()
}

func (cls *Class) IsTrue() bool {
	return cls.base.IsTrue()
}

func (cls *Class) Not() bool {
	return cls.base.Not()
}

func (cls *Class) Free() {
	cls.base.Free()
}

func (cls *Class) newObject(args *Tuple, kwds *Dict) (ClassObject, error) {
	// simple case, a New method has been provided
	if cls.New != nil {
		return cls.New(cls, args, kwds)
	}

	// no New provided, so we need to create an instance of the correct type
	t := reflect.TypeOf(cls.Object).Elem()
	v := reflect.New(t)
	return v.Interface().(ClassObject), nil
}

var otyp = reflect.TypeOf((*Object)(nil)).Elem()

// Clear clear the given Object field correctly. This is equivalent to Py_CLEAR
// from the Python C API.
//
// To clear a field called foo in a struct called self:
//
//	py.Clear(&self.foo)
//
// This will set self.foo to nil, and decrement the reference count of foo.
func Clear[T Object](f *T) {
	tmp := *f
	f = nil
	tmp.Decref()
}

// ClearClassObject clears any contained Objects in the given instance. This
// function will correctly clear (i.e. decref and set to nil) any contained
// Objects in the supplied ClassObject.
func ClearClassObject(co ClassObject) {
	v := reflect.ValueOf(co).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.Type().Implements(otyp) || f.IsNil() {
			// only care about non-nil Object values
			continue
		}
		// Copy the behaviour of Py_CLEAR to avoid issues with loops when
		// calling decref
		tmp := f.Interface().(Object)
		f.SetZero()
		tmp.Decref()
	}
}

type prop struct {
	get Object
	set Object
}

type method struct {
	f     Object
	flags C.int
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

func getPythonCallFlags(f reflect.Type) (C.int, error) {
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

func getStaticCallFlags(f reflect.Type) (C.int, error) {
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
	return nil, fmt.Errorf("TODO(jp3): how should Class.Alloc be working now?")
	// obj, err = class.base.Alloc(n)

	// // Since we are creating this object for Go code, this is probably the only
	// // opportunity we will get to register this object instance.
	// pyType := (*C.PyTypeObject)(unsafe.Pointer(c(class.base)))
	// setClassContext(unsafe.Pointer(c(obj)), pyType)

	// return
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

// Create completes the initialisation of the Class by creating the Python type.
// The created type is then stored in the Class and accessible via the Type
// method. A Class is not a valid Python object until Create has been
// successfully called.
func (cls *Class) Create() error {
	pyType := C.newType()
	pyType.tp_name = C.CString(cls.Name)
	pyType.tp_flags = C.Py_TPFLAGS_DEFAULT | C.ulong(cls.Flags)

	if cls.Object == nil {
		cls.Object = &ClassBaseObject{}
	}

	typ := reflect.TypeOf(cls.Object)
	btyp := typ.Elem()

	if btyp.NumField() == 0 {
		return fmt.Errorf("%s does not embed an Object", btyp.Name())
	}

	methods := make(map[string]method)
	props := make(map[string]prop)

	slotFlags := C.uint64_t(0)

	for flag, slot := range slotMap {
		if typ.Implements(slot) {
			slotFlags |= flag
		}
	}

	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if !strings.HasPrefix(m.Name, "Py") {
			continue
		}
		t := m.Func.Type()
		fn := fmt.Sprintf("%s.%s", typ.Elem().Name(), m.Name)
		parts := strings.SplitN(m.Name, "_", 2)
		switch parts[0] {
		case "Py":
			flags, err := getPythonCallFlags(t)
			if err != nil {
				return fmt.Errorf("%s: %s", fn, err)
			}
			methods[parts[1]] = method{NewLong(int64(i)), flags}
		case "PySet":
			if err := methSigMatches(t, (func(Object) error)(nil)); err != nil {
				C.free(unsafe.Pointer(pyType))
				return fmt.Errorf("%s: %s", fn, err)
			}
			p := props[parts[1]]
			p.set = NewLong(int64(i))
			props[parts[1]] = p
		case "PyGet":
			if err := methSigMatches(t, (func() (Object, error))(nil)); err != nil {
				C.free(unsafe.Pointer(pyType))
				return fmt.Errorf("%s: %s", fn, err)
			}
			p := props[parts[1]]
			p.get = NewLong(int64(i))
			props[parts[1]] = p
		}
	}

	for name, fn := range cls.Static {
		f := reflect.ValueOf(fn)
		t := f.Type()
		flags, err := getStaticCallFlags(t)
		if err != nil {
			return fmt.Errorf("static %s: %s", name, err)
		}
		key, err := NewUnicode(name)
		if err != nil {
			return fmt.Errorf("static %s: %s", name, err)
		}
		methods[name] = method{key, flags | C.METH_STATIC}
	}

	for name, fn := range cls.Class {
		f := reflect.ValueOf(fn)
		t := f.Type()
		flags, err := getPythonCallFlags(t)
		if err != nil {
			return fmt.Errorf("class %s: %s", name, err)
		}
		key, err := NewUnicode(name)
		if err != nil {
			return fmt.Errorf("class %s: %s", name, err)
		}
		methods[name] = method{key, flags | C.METH_CLASS}
	}

	pyType.tp_basicsize = C.Py_ssize_t(unsafe.Sizeof(C.PyObject{}))

	ctxt := C.setSlots(pyType, slotFlags)

	if C.typeReady(pyType) < 0 {
		C.free(unsafe.Pointer(ctxt))
		C.free(unsafe.Pointer(pyType.tp_name))
		C.free(unsafe.Pointer(pyType))
		return exception()
	}

	C.storeContext(pyType, ctxt)

	for name, method := range methods {
		s := C.CString(name)
		C.setTypeAttr(pyType, s, C.newMethod(s, c(method.f), C.int(method.flags)))
	}

	for name, prop := range props {
		s := C.CString(name)
		C.setTypeAttr(pyType, s, C.newProperty(pyType, s, c(prop.get), c(prop.set)))
	}

	for i := 0; i < btyp.NumField(); i++ {
		field := btyp.Field(i)
		pyname := field.Tag.Get("py")
		pydoc := field.Tag.Get("pyDoc")
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
			C.setTypeAttr(pyType, s, C.newObjMember(pyType, s, c(NewLong(int64(i))), d))
			continue
		}
		if !exportable[field.Type.Kind()] {
			C.free(unsafe.Pointer(ctxt))
			C.free(unsafe.Pointer(pyType.tp_name))
			C.free(unsafe.Pointer(pyType))
			return fmt.Errorf("cannot export %s.%s to Python: type '%s' unsupported", btyp.Name(), field.Name, field.Type.Name())
		}
		s := C.CString(pyname)
		defer C.free(unsafe.Pointer(s))
		d := C.CString(pydoc)
		C.setTypeAttr(pyType, s, C.newNatMember(pyType, s, c(NewLong(int64(i))), d))
	}

	cls.base = newType((*C.PyObject)(unsafe.Pointer(pyType)))

	registerClass(pyType, cls)

	return nil
}
