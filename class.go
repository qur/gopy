package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

var otyp = reflect.TypeOf((*Object)(nil)).Elem()

// Class flags. A subset of Python type flags that are supported for types
// implemented in Go. These flags enable optional features, such as supporting
// sub-types.
type ClassFlags uint32

// Flags to enable optional features of classes (types).
const (
	ClassHaveGC   = ClassFlags(C.Py_TPFLAGS_HAVE_GC)
	ClassBaseType = ClassFlags(C.Py_TPFLAGS_BASETYPE)
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
//	Init(args *py.Tuple, kwds *py.Dict) error                 // __init__
//	Call(args *py.Tuple, kwds *py.Dict) (py.Object, error)    // __call__
//	Repr() (Object, error)                                    // __repr__
//	Str() (Object, error)                                     // __str__
//	RichCompare(obj py.Object, op py.Op) (py.Object, error)   // __cmp__
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
//
// To create a new instance of the Class in Go, then use the Callable methods
// (i.e. Call or CallGo), which map to the Python expression "cls(...)".
type Class struct {
	Name   string
	Flags  ClassFlags
	Doc    string
	Object ClassObject
	Static map[string]any
	Class  map[string]any
	New    func(*Class, *Tuple, *Dict) (ClassObject, error)
	User   any
	base   *Type
}

var _ Object = (*Class)(nil)

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (cls *Class) Base() *BaseObject {
	return cls.base.Base()
}

// Type returns a pointer to the Type that represents the type of this object in
// Python.
func (cls *Class) Type() *Type {
	return cls.base.Type()
}

// RawType returns the Typo that represents this class in Python.
//
// Return value: Borrowed Reference.
func (cls *Class) RawType() *Type {
	return cls.base
}

// Decref decrements cls's reference count, cls may not be nil.
func (cls *Class) Decref() {
	cls.base.Decref()
}

// Incref increments cks's reference count, cls may not be nil.
func (cls *Class) Incref() {
	cls.base.Incref()
}

func (cls *Class) raw() *C.PyObject {
	return (*C.PyObject)(unsafe.Pointer(cls.base))
}

// Call calls cls with the given args and kwds. kwds may be nil, args may not
// (an empty Tuple must be used if no arguments are wanted). Returns the result
// of the call, or an Error on failure. This is equivalent to
// "cls(*args, **kwds)" in Python.
//
// Return value: New Reference.
func (cls *Class) Call(args *Tuple, kwds *Dict) (Object, error) {
	ret := C.PyObject_Call(c(cls), c(args), c(kwds))
	return obj2ObjErr(ret)
}

// CallGo calls cf with the given args and kwds, either may be nil. Returns the
// result of the call, or an Error on failure. This is equivalent to
// "cf(*args, **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// Return value: New Reference.
func (cls *Class) CallGo(args []any, kwds map[string]any) (Object, error) {
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
	ret := C.PyObject_Call(c(cls), c(obj1), c(obj2))
	return obj2ObjErr(ret)
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
	reflect.ValueOf(f).Elem().SetZero()
	tmp.Decref()
}

// ClearClassObject clears any contained Objects in the given instance. This
// function will correctly clear (i.e. decref and set to nil) any exported
// contained Objects in the supplied ClassObject. If the ClassObject contains
// private Objects then they will need to Cleared explicitly using Clear from a
// Dealloc method.
func ClearClassObject(co ClassObject) {
	v := reflect.ValueOf(co).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.Type().Implements(otyp) || f.IsNil() || !v.Type().Field(i).IsExported() {
			// only care about exported non-nil Object values
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
		return fmt.Errorf("method without receiver")
	}

	want := reflect.TypeOf(_want)

	if got.NumIn()-1 != want.NumIn() {
		return fmt.Errorf("method should have %d arguments, not %d", want.NumIn(), got.NumIn()-1)
	}

	if got.NumOut() != want.NumOut() {
		return fmt.Errorf("method should have %d return values, not %d", want.NumOut(), got.NumOut())
	}

	for i := 0; i < want.NumIn(); i++ {
		if got.In(i+1) != want.In(i) {
			return fmt.Errorf("method argument %d should be %v, not %v", i+1, want.In(i), got.In(i+1))
		}
	}

	for i := 0; i < want.NumOut(); i++ {
		if got.Out(i) != want.Out(i) {
			return fmt.Errorf("method return value %d should be %v, not %v", i+1, want.Out(i), got.Out(i))
		}
	}

	return nil
}

func getPythonCallFlags(f reflect.Type) (C.int, error) {
	switch {
	case methSigMatches(f, pyNoArgsFunc) == nil:
		return C.METH_NOARGS, nil
	case methSigMatches(f, pySingleFunc) == nil:
		return C.METH_O, nil
	case methSigMatches(f, pyVarargsFunc) == nil:
		return C.METH_VARARGS, nil
	case methSigMatches(f, pyKeywordsFunc) == nil:
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
	case funcSigMatches(f, pyNoArgsFunc) == nil:
		return C.METH_NOARGS, nil
	case funcSigMatches(f, pySingleFunc) == nil:
		return C.METH_O, nil
	case funcSigMatches(f, pyVarargsFunc) == nil:
		return C.METH_VARARGS, nil
	case funcSigMatches(f, pyKeywordsFunc) == nil:
		return C.METH_VARARGS | C.METH_KEYWORDS, nil
	default:
		return 0, fmt.Errorf("invalid function signature")
	}
}

var exportable = map[reflect.Kind]bool{
	reflect.Bool:       true,
	reflect.Int:        true,
	reflect.Int8:       true,
	reflect.Int16:      true,
	reflect.Int32:      true,
	reflect.Int64:      true,
	reflect.Uint:       true,
	reflect.Uint8:      true,
	reflect.Uint16:     true,
	reflect.Uint32:     true,
	reflect.Float32:    true,
	reflect.Float64:    true,
	reflect.String:     true,
	reflect.Complex64:  true,
	reflect.Complex128: true,
}

type methodSignature struct {
	field string
	sig   interface{}
}

// Function signatures for methods that implement Python methods. Note, the
// first argument is always the receiver, and is not included in these
// signatures.
var (
	pyNoArgsFunc   = (func() (Object, error))(nil)
	pySingleFunc   = (func(Object) (Object, error))(nil)
	pyVarargsFunc  = (func(*Tuple) (Object, error))(nil)
	pyKeywordsFunc = (func(*Tuple, *Dict) (Object, error))(nil)
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
		pyEmbed := false
		switch field.Type {
		case cipType:
			if _, ok := cls.Object.(tp_iternext); !ok {
				C.free(unsafe.Pointer(ctxt))
				C.free(unsafe.Pointer(pyType.tp_name))
				C.free(unsafe.Pointer(pyType))
				return fmt.Errorf("%T claimed to implement IteratorProtocol by embedding ClassIteratorProtocol, but doesn't have required methods", cls.Object)
			}
			pyEmbed = true
		case cspType:
			if _, ok := cls.Object.(sq_item); !ok {
				C.free(unsafe.Pointer(ctxt))
				C.free(unsafe.Pointer(pyType.tp_name))
				C.free(unsafe.Pointer(pyType))
				return fmt.Errorf("%T claimed to implement SequenceProtocol by embedding ClassSequenceProtocol, but doesn't have required methods", cls.Object)
			}
			pyEmbed = true
		case cmpType:
			if _, ok := cls.Object.(mp_subscript); !ok {
				C.free(unsafe.Pointer(ctxt))
				C.free(unsafe.Pointer(pyType.tp_name))
				C.free(unsafe.Pointer(pyType))
				return fmt.Errorf("%T claimed to implement MappingProtocol by embedding ClassMappingProtocol, but doesn't have required methods", cls.Object)
			}
			pyEmbed = true
		case cboType, cnpType:
			pyEmbed = true
		}
		if pyEmbed || !field.IsExported() {
			// We have some helper types that get embedded in the ClassObject
			// implementation. Don't export these to Python. We also ignore
			// anything that isn't exported.
			continue
		}
		pyname := field.Tag.Get("py")
		if pyname == "-" {
			// tag explicitly set to ignore field
			continue
		}
		pydoc := field.Tag.Get("pyDoc")
		ro := C.int(0)
		parts := strings.Split(pyname, ",")
		if len(parts) > 0 {
			pyname = parts[0]
			for _, opt := range parts[1:] {
				switch opt {
				case "ro":
					ro = C.int(1)
				default:
					C.free(unsafe.Pointer(ctxt))
					C.free(unsafe.Pointer(pyType.tp_name))
					C.free(unsafe.Pointer(pyType))
					return fmt.Errorf("unknown tag option: %s", opt)
				}
			}
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
			C.setTypeAttr(pyType, s, C.newObjMember(pyType, s, c(NewLong(int64(i))), d, ro))
			continue
		}
		if exportable[field.Type.Kind()] {
			// field is a simple exportable native type, we can use the native
			// member get/set code.
			s := C.CString(pyname)
			defer C.free(unsafe.Pointer(s))
			d := C.CString(pydoc)
			C.setTypeAttr(pyType, s, C.newNatMember(pyType, s, c(NewLong(int64(i))), d, ro))
			continue
		}
		C.free(unsafe.Pointer(ctxt))
		C.free(unsafe.Pointer(pyType.tp_name))
		C.free(unsafe.Pointer(pyType))
		return fmt.Errorf("cannot export %s.%s to Python: type '%s' unsupported", btyp.Name(), field.Name, field.Type.Name())
	}

	cls.base = newType(pyType)

	registerClass(pyType, cls)

	return nil
}
