package py

// #include "utils.h"
import "C"

import (
	"errors"
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
	ClassHeapType = ClassFlags(C.Py_TPFLAGS_HEAPTYPE)
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
// BaseType is the base type of the class. Although it is an Object, this value
// must be either a *Type or *Class (or nil for no base type).
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
	Name     string
	Flags    ClassFlags
	Doc      string
	BaseType Object
	Object   ClassObject
	Static   map[string]any
	Class    map[string]any
	New      func(*Class, *Tuple, *Dict) (ClassObject, error)
	User     any
	base     *Type
}

var _ Object = (*Class)(nil)

// Base returns a BaseObject pointer that gives access to the generic methods on
// that type for this object.
func (cls *Class) Base() *BaseObject {
	return cls.base.Base()
}

// Type returns a pointer to the Type that represents the type of this object in
// Python. This will always be TypeType.
func (cls *Class) Type() *Type {
	return cls.base.Type()
}

// RawType returns the Type that represents this class in Python. The returned
// value is the same Python object, but as *Type instead of *Class.
//
// Return value: Borrowed Reference.
func (cls *Class) RawType() *Type {
	return cls.base
}

// Decref decrements cls's reference count, cls may not be nil. This should only
// be used when the ClassHeapType flag is set.
func (cls *Class) Decref() {
	cls.base.Decref()
}

// Incref increments cls's reference count, cls may not be nil. This should only
// be used when the ClassHeapType flag is set.
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
// The returned value will be a instance of the Class's Object type.
//
// Return value: New Reference.
func (cls *Class) Call(args *Tuple, kwds *Dict) (Object, error) {
	ret := C.PyObject_Call(c(cls), c(args), c(kwds))
	return obj2ObjErr(ret)
}

// CallGo calls cls with the given args and kwds, either may be nil. Returns the
// result of the call, or an Error on failure. This is equivalent to "cls(*args,
// **kwds)" in Python.
//
// The values are converted to Objects using NewValue. A TypeError will be
// returned if a value cannot be converted.
//
// The returned value will be a instance of the Class's Object type.
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

// Super returns a Super object for the Class. If obj is not nil, then it will
// be a bound Super to that object.
//
// This is equivalent to `super(cls, obj)` in Python, or `super(cls)` if obj is
// nil.
//
// Return value: New Reference.
func (cls *Class) Super(obj Object) (*Super, error) {
	args := make([]Object, 0, 2)

	cls.Incref()

	args = append(args, cls)

	if obj != nil {
		obj.Incref()
		args = append(args, obj)
	}

	t, err := PackTuple(args...)
	if err != nil {
		return nil, err
	}

	o, err := SuperType.Call(t, nil)
	if err != nil {
		return nil, err
	}

	return newSuper(c(o)), nil
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

// Instance returns o as an instance of the Class, i.e. as the same type as
// Class.Object. If o is not an instance of Class, an instance of a supertype of
// Class, or an instance of a subtype of Class which was previously down-casted
// from Class (or a superclass), then nil will be returned.
//
// Return value: Borrowed Reference.
func (cls *Class) Instance(o Object) ClassObject {
	return getClassObjectByType(c(o), cls.base.c())
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
	for i := range v.NumField() {
		f := v.Field(i)
		if !f.Type().Implements(otyp) || f.Type() == cboType || f.IsNil() || !v.Type().Field(i).IsExported() {
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
	if got.NumIn() == 0 {
		// The receiver is missing!
		return errors.New("method without receiver")
	}

	want := reflect.TypeOf(_want)

	// Note: Methods take the receiver as the first argument, which the want
	// signature doesn't include.

	if got.NumIn()-1 != want.NumIn() {
		return fmt.Errorf("method should have %d arguments, not %d", want.NumIn(), got.NumIn()-1)
	}

	if got.NumOut() != want.NumOut() {
		return fmt.Errorf("method should have %d return values, not %d", want.NumOut(), got.NumOut())
	}

	for i := range want.NumIn() {
		if got.In(i+1) != want.In(i) {
			return fmt.Errorf("method argument %d should be %v, not %v", i+1, want.In(i), got.In(i+1))
		}
	}

	for i := range want.NumOut() {
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
		return 0, errors.New("invalid method signature")
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

	for i := range want.NumIn() {
		if got.In(i) != want.In(i) {
			return fmt.Errorf("function argument %d should be %v, not %v", i+1, want.In(i), got.In(i))
		}
	}

	for i := range want.NumOut() {
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
		return 0, errors.New("invalid function signature")
	}
}

func getCallFlags(f reflect.Type, kind C.int) (C.int, error) {
	switch kind {
	case C.METH_CLASS:
		return getPythonCallFlags(f)
	case C.METH_STATIC:
		return getStaticCallFlags(f)
	default:
		return 0, fmt.Errorf("unsupported kind for getCallFlags: %d", kind)
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

const subclassFlags = (C.Py_TPFLAGS_LONG_SUBCLASS |
	C.Py_TPFLAGS_LIST_SUBCLASS |
	C.Py_TPFLAGS_TUPLE_SUBCLASS |
	C.Py_TPFLAGS_BYTES_SUBCLASS |
	C.Py_TPFLAGS_UNICODE_SUBCLASS |
	C.Py_TPFLAGS_DICT_SUBCLASS |
	C.Py_TPFLAGS_TYPE_SUBCLASS |
	C.Py_TPFLAGS_BASE_EXC_SUBCLASS)

func fastSubclassFlags(t *Type) C.ulong {
	switch t {
	case LongType:
		return C.Py_TPFLAGS_LONG_SUBCLASS
	case ListType:
		return C.Py_TPFLAGS_LIST_SUBCLASS
	case TupleType:
		return C.Py_TPFLAGS_TUPLE_SUBCLASS
	case BytesType:
		return C.Py_TPFLAGS_BYTES_SUBCLASS
	case UnicodeType:
		return C.Py_TPFLAGS_UNICODE_SUBCLASS
	case DictType:
		return C.Py_TPFLAGS_DICT_SUBCLASS
	case TypeType:
		return C.Py_TPFLAGS_TYPE_SUBCLASS
	}

	if t.o.tp_flags&subclassFlags != 0 {
		return t.o.tp_flags & subclassFlags
	}

	return 0
}

// Create completes the initialisation of the Class by creating the Python type.
// The created type is then stored in the Class and accessible via the RawType
// method. A Class is not a valid Python object until Create has been
// successfully called.
func (cls *Class) Create() (err error) {
	name := C.CString(cls.Name)
	defer cfree(name)

	pyHeapType := C.newType(C.ulong(cls.Flags))
	pyHeapType._ht_tpname = C.copyName(name)
	pyType := &pyHeapType.ht_type
	pyType.tp_name = pyHeapType._ht_tpname
	pyType.tp_flags = C.Py_TPFLAGS_DEFAULT | C.ulong(cls.Flags)
	pyType.tp_basicsize = C.Py_ssize_t(unsafe.Sizeof(C.PyObject{}))
	pyType.tp_itemsize = 0

	defer func() {
		if err != nil {
			if cls.Flags&ClassHeapType != 0 {
				xdecref(pyHeapType.ht_name)
				C.PyObject_Free(unsafe.Pointer(pyHeapType))
			} else {
				cfree(pyHeapType)
			}
		}
	}()

	if cls.Flags&ClassHeapType != 0 {
		// for heap allocated types we need to setup the ht_name and ht_qualname
		// members.
		pyHeapType.ht_name = C.typeName(pyHeapType._ht_tpname)
		if pyHeapType.ht_name == nil {
			return exception()
		}

		incref(pyHeapType.ht_name)
		pyHeapType.ht_qualname = pyHeapType.ht_name
	}

	// start by validating BaseType
	if err := cls.validateBaseType(pyType); err != nil {
		return err
	}

	if cls.Object == nil {
		cls.Object = &ClassBaseObject{}
	}

	typ := reflect.TypeOf(cls.Object)

	setupSlots(pyHeapType, typ)

	if err := cls.setupMethodsAndProperties(pyType, typ); err != nil {
		return err
	}

	if err := cls.setupFields(pyType, typ.Elem()); err != nil {
		return err
	}

	cls.base = newType(pyType)
	registerClass(pyType, cls)

	return nil
}

func (cls *Class) validateBaseType(pyType *C.PyTypeObject) error {
	switch b := cls.BaseType.(type) {
	case nil:
		// ok, no base type
	case *Type:
		// *Type is good, but shouldn't be nil
		if b == nil {
			return errors.New("BaseType set, but nil")
		}

		pyType.tp_base = b.c()
		pyType.tp_basicsize = b.o.tp_basicsize
		pyType.tp_itemsize = b.o.tp_itemsize
		pyType.tp_flags |= fastSubclassFlags(b)
	case *Class:
		// *Class is good, but should be initialised and not nil
		if b == nil {
			return errors.New("BaseType set, but nil")
		}

		raw := b.RawType()
		if raw == nil {
			return errors.New("can't use uninitialised *Class as BaseType")
		}

		pyType.tp_base = raw.c()
		pyType.tp_basicsize = raw.o.tp_basicsize
		pyType.tp_itemsize = raw.o.tp_itemsize
		pyType.tp_flags |= fastSubclassFlags(raw)
	default:
		return fmt.Errorf("%T is not a supported type for BaseType", b)
	}

	return nil
}

func setupSlots(pyHeapType *C.PyHeapTypeObject, typ reflect.Type) {
	slotFlags := C.uint64_t(0)

	for flag, slot := range slotMap {
		if typ.Implements(slot) {
			slotFlags |= flag
		}
	}

	C.setSlots(pyHeapType, slotFlags)
}

func (cls *Class) setupMethodsAndProperties(pyType *C.PyTypeObject, typ reflect.Type) error {
	methods := make(map[string]method)
	props := make(map[string]prop)

	if err := addMethods(methods, cls.Static, C.METH_STATIC); err != nil {
		return err
	}

	if err := addMethods(methods, cls.Class, C.METH_CLASS); err != nil {
		return err
	}

	if err := extractMethodsAndProperties(methods, props, typ); err != nil {
		return err
	}

	if C.typeReady(pyType) < 0 {
		return exception()
	}

	for name, method := range methods {
		s := C.CString(name)

		C.setTypeAttr(pyType, s, C.newMethod(pyType, s, c(method.f), method.flags))
	}

	for name, prop := range props {
		s := C.CString(name)

		C.setTypeAttr(pyType, s, C.newProperty(pyType, s, c(prop.get), c(prop.set)))
	}

	return nil
}

func addMethods(methods map[string]method, functions map[string]any, kind C.int) error {
	for name, fn := range functions {
		f := reflect.ValueOf(fn)
		t := f.Type()

		flags, err := getCallFlags(t, kind)
		if err != nil {
			return fmt.Errorf("static %s: %w", name, err)
		}

		key, err := NewUnicode(name)
		if err != nil {
			return fmt.Errorf("static %s: %w", name, err)
		}

		methods[name] = method{key, flags | kind}
	}

	return nil
}

func extractMethodsAndProperties(methods map[string]method, props map[string]prop, typ reflect.Type) error {
	for i := range typ.NumMethod() {
		m := typ.Method(i)
		if !strings.HasPrefix(m.Name, "Py") {
			continue
		}

		t := m.Func.Type()
		fn := fmt.Sprintf("%s.%s", typ.Elem().Name(), m.Name)

		switch parts := strings.SplitN(m.Name, "_", 2); parts[0] {
		case "Py":
			flags, err := getPythonCallFlags(t)
			if err != nil {
				return fmt.Errorf("%s: %w", fn, err)
			}

			methods[parts[1]] = method{NewLong(int64(i)), flags}
		case "PySet":
			if err := methSigMatches(t, (func(Object) error)(nil)); err != nil {
				return fmt.Errorf("%s: %w", fn, err)
			}

			p := props[parts[1]]
			p.set = NewLong(int64(i))
			props[parts[1]] = p
		case "PyGet":
			if err := methSigMatches(t, (func() (Object, error))(nil)); err != nil {
				return fmt.Errorf("%s: %w", fn, err)
			}

			p := props[parts[1]]
			p.get = NewLong(int64(i))
			props[parts[1]] = p
		}
	}

	return nil
}

func (cls *Class) setupFields(pyType *C.PyTypeObject, btyp reflect.Type) error {
	for i := range btyp.NumField() {
		field := btyp.Field(i)

		ignore, err := cls.ignoreEmbedded(field)
		if err != nil {
			return err
		}

		if ignore || !field.IsExported() {
			// We have some helper types that get embedded in the ClassObject
			// implementation. Don't export these to Python. We also ignore
			// anything that isn't exported.
			continue
		}

		def, err := newFieldDef(field)
		if err != nil {
			return err
		}

		if def.name == "-" {
			// tag explicitly set to ignore field
			continue
		}

		if err := def.exportToPython(i, pyType, btyp); err != nil {
			return err
		}
	}

	return nil
}

func (cls *Class) ignoreEmbedded(field reflect.StructField) (bool, error) {
	switch field.Type {
	case cipType:
		if _, ok := cls.Object.(tp_iternext); !ok {
			return false, fmt.Errorf("%T claimed to implement IteratorProtocol "+
				"by embedding ClassIteratorProtocol, but doesn't have required methods", cls.Object)
		}

		return true, nil

	case cspType:
		if _, ok := cls.Object.(sq_item); !ok {
			return false, fmt.Errorf("%T claimed to implement SequenceProtocol "+
				"by embedding ClassSequenceProtocol, but doesn't have required methods", cls.Object)
		}

		return true, nil

	case cmpType:
		if _, ok := cls.Object.(mp_subscript); !ok {
			return false, fmt.Errorf("%T claimed to implement MappingProtocol "+
				"by embedding ClassMappingProtocol, but doesn't have required methods", cls.Object)
		}

		return true, nil

	case cboType, cnpType:
		return true, nil
	}

	return false, nil
}

type fieldDef struct {
	name  string
	doc   string
	field reflect.StructField
	ro    C.int
}

func newFieldDef(field reflect.StructField) (fieldDef, error) {
	name := field.Tag.Get("py")
	if name == "" {
		name = field.Name
	}

	f := fieldDef{
		name:  name,
		doc:   field.Tag.Get("pyDoc"),
		field: field,
		ro:    0,
	}

	if err := f.parseFlags(); err != nil {
		return fieldDef{}, err
	}

	return f, nil
}

func (f *fieldDef) parseFlags() error {
	parts := strings.Split(f.name, ",")
	if len(parts) > 0 {
		f.name = parts[0]

		for _, opt := range parts[1:] {
			switch opt {
			case "ro":
				f.ro = 1
			default:
				return fmt.Errorf("unknown tag option: %s", opt)
			}
		}
	}

	return nil
}

func (f *fieldDef) exportToPython(i int, pyType *C.PyTypeObject, btyp reflect.Type) error {
	if f.field.Type.Implements(otyp) {
		// field is some type of object, so we can use the generic object
		// member get/set code.
		s := C.CString(f.name)
		defer cfree(s)

		C.setTypeAttr(pyType, s, C.newObjMember(pyType, s, c(NewLong(int64(i))), C.CString(f.doc), f.ro))

		return nil
	}

	if exportable[f.field.Type.Kind()] {
		// field is a simple exportable native type, we can use the native
		// member get/set code.
		s := C.CString(f.name)
		defer cfree(s)

		C.setTypeAttr(pyType, s, C.newNatMember(pyType, s, c(NewLong(int64(i))), C.CString(f.doc), f.ro))

		return nil
	}

	return fmt.Errorf("cannot export %s.%s to Python: type '%s' unsupported",
		btyp.Name(), f.field.Name, f.field.Type.Name())
}
