// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
// static inline PyTypeObject *newType(void) {
//     return calloc(1, sizeof(PyTypeObject));
// }
// static inline int typeReady(PyTypeObject *o) {
//     int ret;
//     if (o->tp_new == NULL) {
//         o->tp_new = PyType_GenericNew;
//     }
//     if (o->tp_flags & Py_TPFLAGS_HAVE_GC) {
//         enableClassGc(o);
//     }
//     ret = PyType_Ready(o);
//     if (ret == 0) {
//         // We don't use tp_methods, and it is read when calling PyType_Ready
//         // - so we use it to hide a classContext struct.  The classContext
//         // starts with a NULL pointer just in case, so it looks like an
//         // empty methods list if Python does try to process it.
//         o->tp_methods = calloc(1, sizeof(ClassContext));
//     }
//     return ret;
// }
// static inline int setTypeAttr(PyTypeObject *tp, char *name, PyObject *o) {
//     return PyDict_SetItemString(tp->tp_dict, name, o);
// }
// static inline PyObject *typeAlloc(PyTypeObject *type, Py_ssize_t n) { return type->tp_alloc(type, n); }
// static inline int doVisit(PyObject *o, void *v, void *a) {
//     visitproc visit = v;
//     return visit(o, a);
// }
import "C"

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

const (
	TPFLAGS_HAVE_GC = uint32(C.Py_TPFLAGS_HAVE_GC)
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
// Type holds a Pointer to the Type instance for this class, this is filled in
// by calling Create().
//
// Pointer should be set to a pointer of the struct type that will represent an
// instance of the Python class.  This struct must contain an embedded
// py.BaseObject as its first member.  The easiest ways to set Pointer are
// either to use a struct literal (i.e. &MyClass{}), or to cast nil (i.e.
// (*MyClass)(nil)), if the struct is large then the latter method is more
// efficient (as an instance of the struct is not created).
//
// This struct may have the following special methods (the equivalent Python
// methods are also indicated):
//
//   PyInit(args *py.Tuple, kwds *py.Dict) os.Error              // __init__
//   PyCall(args *py.Tuple, kwds *py.Dict) (py.Object, os.Error) // __call__
//   PyRepr() string                                             // __repr__
//   PyStr() string                                              // __str__
//   PyCompare(obj py.Object) (int, os.Error)                    // __cmp__
//
// Properties are also supported, by implementing get and set methods:
//
//   PyGet_XXX() (py.Object, os.Error)
//   PySet_XXX(value py.Object) os.Error
//
// Methods on the Python class are implemented by methods with the Py_ prefix:
//
//   Py_XXX(args *py.Tuple, kwds *py.Dict) (py.Object, os.Error)
//
type Class struct {
	Name    string
	Flags   uint32
	Doc     string
	Type    *Type
	Pointer interface{}
	New     func(*Class, *Tuple, *Dict) (Object, os.Error)
}

var otyp = reflect.TypeOf(new(Object)).Elem()

//export goClassCallMethod
func goClassCallMethod(obj, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Unpack context and self pointer from obj
	t := (*C.PyObject)(obj)
	pyobj := unsafe.Pointer(C.PyTuple_GetItem(t, 0))
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 1), nil)

	// Get args and kwds ready to use, by turning them into pointers of the
	// appropriate type
	a := newTuple((*C.PyObject)(args))
	k := newDict((*C.PyObject)(kwds))

	// Now call the actual struct method by pulling the method out of the
	// reflect.Type object stored in the context
	f := (*func(p unsafe.Pointer, a *Tuple, k *Dict) (Object, os.Error))(unsafe.Pointer(&m))

	ret, err := (*f)(pyobj, a, k)
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
	a := newBaseObject((*C.PyObject)(arg)).actual()

	// Turn the function into something we can call
	f := (*func(p unsafe.Pointer, a Object) os.Error)(unsafe.Pointer(&m))

	err := (*f)(obj, a)
	if err != nil {
		raise(err)
		return -1
	}

	return 0
}

//export goClassGetProp
func goClassGetProp(obj, closure unsafe.Pointer) unsafe.Pointer {
	// Unpack set function from closure
	t := (*C.PyObject)(closure)
	m := C.PyCapsule_GetPointer(C.PyTuple_GetItem(t, 0), nil)

	// Turn the function into something we can call
	f := (*func(p unsafe.Pointer) (Object, os.Error))(unsafe.Pointer(&m))

	ret, err := (*f)(obj)
	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(ret))
}

//export goClassObjGet
func goClassObjGet(obj unsafe.Pointer, idx int) unsafe.Pointer {
	field := fields[idx]
	item := unsafe.Pointer(uintptr(obj) + field.Offset)

	var o Object

	if field.Type == otyp {
		o = *(*Object)(item)
	} else {
		o = unsafe.Unreflect(field.Type, item).(Object)
	}

	o.Incref()
	return unsafe.Pointer(c(o))
}

//export goClassObjSet
func goClassObjSet(obj unsafe.Pointer, idx int, obj2 unsafe.Pointer) int {
	field := fields[idx]
	item := unsafe.Pointer(uintptr(obj) + field.Offset)

	// This is the new value we are being asked to set
	value := newBaseObject((*C.PyObject)(obj2)).actual()

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
	o := unsafe.Unreflect(reflect.PtrTo(field.Type), unsafe.Pointer(&item))
	ov := reflect.ValueOf(o).Elem()

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
	raise(TypeError("Cannot assign '%T' to '*%v'", value, tn))
	return -1
}

//export goClassNatGet
func goClassNatGet(obj unsafe.Pointer, idx int) unsafe.Pointer {
	field := fields[idx]
	item := unsafe.Pointer(uintptr(obj) + field.Offset)

	switch field.Type.Kind() {
	case reflect.Int:
		i := (*int)(item)
		return unsafe.Pointer(C.PyInt_FromLong(C.long(*i)))
	}

	raise(fmt.Errorf("Not Implemented"))
	return nil
}

//export goClassNatSet
func goClassNatSet(obj unsafe.Pointer, idx int, obj2 unsafe.Pointer) int {
	field := fields[idx]
	item := unsafe.Pointer(uintptr(obj) + field.Offset)

	// This is the new value we are being asked to set
	value := newBaseObject((*C.PyObject)(obj2)).actual()

	switch field.Type.Kind() {
	case reflect.Int:
		v := int(C.PyInt_AsLong(c(value)))
		if exceptionRaised() {
			return -1
		}
		i := (*int)(item)
		*i = v
		return 0
	}

	raise(fmt.Errorf("Not Implemented"))
	return -1
}

//export goClassTraverse
func goClassTraverse(obj, visit, arg unsafe.Pointer) int {
	// Get the Python type object
	pyType := (*C.PyTypeObject)((*C.PyObject)(obj).ob_type)

	class, ok := types[pyType]
	if !ok {
		raise(fmt.Errorf("TypeError: Not a recognised type"))
		return -1
	}

	typ := unsafe.Typeof(class.Pointer)
	st := reflect.TypeOf(unsafe.Unreflect(typ, unsafe.Pointer(&obj))).Elem()

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
			o = *(**BaseObject)(v)
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

	class, ok := types[pyType]
	if !ok {
		raise(fmt.Errorf("TypeError: Not a recognised type"))
		return -1
	}

	typ := unsafe.Typeof(class.Pointer)
	st := reflect.TypeOf(unsafe.Unreflect(typ, unsafe.Pointer(&obj))).Elem()

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
			o := (**BaseObject)(v)
			tmp := *o
			*o = nil
			Decref(tmp)
		}
	}

	return 0
}

func getClassContext(obj unsafe.Pointer) *C.ClassContext {
	o := (*C.PyObject)(obj)
	return (*C.ClassContext)(unsafe.Pointer(o.ob_type.tp_methods))
}

//export goClassNew
func goClassNew(typ, args, kwds unsafe.Pointer) unsafe.Pointer {
	// Get the Python type object
	pyType := (*C.PyTypeObject)(typ)

	class, ok := types[pyType]
	if !ok {
		raise(fmt.Errorf("TypeError: Not a recognised type"))
		return nil
	}

	var obj Object
	var err os.Error

	if class.New != nil {
		// Get args and kwds ready to use, by turning them into pointers of the
		// appropriate type
		a := newTuple((*C.PyObject)(args))
		k := newDict((*C.PyObject)(kwds))

		obj, err = class.New(class, a, k)
	} else {
		// Create a new Python instance
		obj, err = class.Alloc(0)
	}

	if err != nil {
		raise(err)
		return nil
	}

	return unsafe.Pointer(c(obj))
}

type prop struct {
	get unsafe.Pointer
	set unsafe.Pointer
}

func methSigMatches(got reflect.Type, _want interface{}) os.Error {
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

func (class *Class) Alloc(n int64) (Object, os.Error) {
	pyType := (*C.PyTypeObject)(unsafe.Pointer(c(class.Type)))

	obj := C.typeAlloc(pyType, 0)
	if obj == nil {
		return nil, exception()
	}

	return newBaseObject(obj).actual(), nil
}

func Clear(obj Object) os.Error {
	ret := goClassClear(unsafe.Pointer(c(obj)))
	if ret < 0 {
		return exception()
	}
	return nil
}

var fields []reflect.StructField

func registerField(field reflect.StructField) C.int {
	fields = append(fields, field)
	return C.int(len(fields) - 1)
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

type goMethod struct {
	field string
	sig   interface{}
}

// Function signatures for methods that implement Python methods.  Note, the
// first argument is always the receiver, and is not included in these
// signatures (hence the names are one greater than the number of arguments
// taken).
var (
	pyInitFunc        = func(*Tuple, *Dict) os.Error(nil)
	pyVoidFunc        = (func())(nil)
	pyReprFunc        = func() string(nil)
	pyLenFunc         = func() int64(nil)
	pyInquiryFunc     = func() (bool, os.Error)(nil)
	pyUnaryFunc       = func() (Object, os.Error)(nil)
	pyBinaryFunc      = func(Object) (Object, os.Error)(nil)
	pyTernaryFunc     = func(a, b Object) (Object, os.Error)(nil)
	pyBinaryCallFunc  = func(*Tuple) (Object, os.Error)(nil)
	pyTernaryCallFunc = func(*Tuple, *Dict) (Object, os.Error)(nil)
	pyCompareFunc     = func(Object) (int, os.Error)(nil)
	pyObjObjArgFunc   = func(a, b Object) os.Error(nil)
	pySsizeArgFunc    = func(int64) (Object, os.Error)(nil)
	pySsizeObjArgFunc = func(int64, Object) os.Error(nil)
	pyObjObjFunc      = func(Object) (bool, os.Error)(nil)
)

var methodMap = map[string]goMethod{
	// Standard Methods
	"PyCall":    {"call", pyTernaryCallFunc},
	"PyCompare": {"compare", pyCompareFunc},
	"PyDealloc": {"dealloc", pyVoidFunc},
	//"PyGetAttr":    {"getattr", },
	//"PyGetAttrObj": {"getattro", },
	//"PyHash":       {"hash", },
	"PyInit": {"init", pyInitFunc},
	//"PyIter":       {"iter", },
	//"PyIterNext":   {"iternext", },
	//"PyPrint":      {"print", },
	"PyRepr": {"repr", pyReprFunc},
	//"PyRichCmp":    {"richcmp", },
	//"PySetAttr":    {"setattr", },
	//"PySetAttrObj": {"setattro", },
	"PyStr": {"str", pyReprFunc},

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

// Create creates and returns a pointer to a PyTypeObject that is the Python
// representation of the class that has been implemented in Go.
func (c *Class) Create() (*Type, os.Error) {
	typ := reflect.TypeOf(c.Pointer)

	pyType := C.newType()
	pyType.tp_name = C.CString(c.Name)
	pyType.tp_basicsize = C.Py_ssize_t(typ.Elem().Size())
	pyType.tp_flags = C.Py_TPFLAGS_DEFAULT | C.Py_TPFLAGS_CHECKTYPES | C.long(c.Flags)

	if C.typeReady(pyType) < 0 {
		C.free(unsafe.Pointer(pyType.tp_name))
		C.free(unsafe.Pointer(pyType))
		return nil, exception()
	}

	registerType(pyType, c)

	// Get the context
	ctxt := (*C.ClassContext)(unsafe.Pointer(pyType.tp_methods))

	btyp := typ.Elem()
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
			return nil, fmt.Errorf("Cannot export %s.%s to Python: type '%s' unsupported", btyp.Name(), field.Name, field.Type.Name())
		}
		s := C.CString(pyname)
		defer C.free(unsafe.Pointer(s))
		d := C.CString(pydoc)
		C.setTypeAttr(pyType, s, C.newNatMember(registerField(field), d))
	}

	props := make(map[string]prop)

	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if !strings.HasPrefix(m.Name, "Py") {
			continue
		}
		t := m.Func.Type()
		f := unsafe.Pointer(m.Func.Pointer())
		fn := fmt.Sprintf("%s.%s", typ.Elem().Name(), m.Name)
		meth, ok := methodMap[m.Name]
		if ok {
			err := methSigMatches(t, meth.sig)
			if err != nil {
				return nil, fmt.Errorf("%s: %s", fn, err)
			}
			ctxtSet(ctxt, meth.field, f)
			continue
		}
		parts := strings.SplitN(m.Name, "_", 2)
		switch parts[0] {
		case "Py":
			err := methSigMatches(t, func(a *Tuple, k *Dict) (Object, os.Error)(nil))
			if err != nil {
				return nil, fmt.Errorf("%s: %s", fn, err)
			}
			s := C.CString(parts[1])
			C.setTypeAttr(pyType, s, C.newMethod(s, f))
		case "PySet":
			err := methSigMatches(t, func(Object) os.Error(nil))
			if err != nil {
				return nil, fmt.Errorf("%s: %s", fn, err)
			}
			p := props[parts[1]]
			p.set = f
			props[parts[1]] = p
		case "PyGet":
			err := methSigMatches(t, func() (Object, os.Error)(nil))
			if err != nil {
				return nil, fmt.Errorf("%s: %s", fn, err)
			}
			p := props[parts[1]]
			p.get = f
			props[parts[1]] = p
		}
	}

	C.setClassContext(pyType, ctxt)

	for name, prop := range props {
		s := C.CString(name)
		C.setTypeAttr(pyType, s, C.newProperty(pyType, s, prop.get, prop.set))
	}

	c.Type = newType((*C.PyObject)(unsafe.Pointer(pyType)))

	return c.Type, nil
}
