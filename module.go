package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

func Import(name string) (*Module, error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))

	pyName := C.PyUnicode_FromString(s)
	defer C.decref(pyName)

	obj := C.PyImport_Import(pyName)
	if obj == nil {
		return nil, exception()
	}

	return newModule(obj), nil
}

// TODO(jp3): add support for other fields of PyModuleDef

type ModuleDef struct {
	Name    string
	Doc     string
	Package bool
	Methods []GoMethod
}

func CreateModule(md *ModuleDef) (*Module, error) {
	pyMD := C.newModuleDef()

	pyMD.m_name = C.CString(md.Name)
	if md.Doc != "" {
		pyMD.m_doc = C.CString(md.Doc)
	}

	m := C.PyModule_Create2(pyMD, C.PYTHON_API_VERSION)
	if m == nil {
		return nil, exception()
	}
	mod := newModule(m)

	if md.Package {
		// mark module as package by adding an empty list as __path__
		l, err := NewList(0)
		if err != nil {
			mod.Decref()
			return nil, err
		}
		defer l.Decref()

		if err := mod.AddObjectRef("__path__", l); err != nil {
			mod.Decref()
			return nil, err
		}
	}

	if len(md.Methods) == 0 {
		return mod, nil
	}

	n := C.PyUnicode_FromString(pyMD.m_name)
	if n == nil {
		mod.Decref()
		return nil, exception()
	}
	defer newObject(n).Decref()

	d := C.PyModule_GetDict(m)
	if d == nil {
		mod.Decref()
		return nil, exception()
	}

	for _, method := range md.Methods {
		pyF, err := makeCFunction(method.Name, method.Func, method.Doc, n)
		if err != nil {
			mod.Decref()
			return nil, err
		}
		defer pyF.Decref()

		if C.PyDict_SetItemString(d, C.CString(method.Name), c(pyF)) != 0 {
			mod.Decref()
			return nil, exception()
		}
	}

	return mod, nil
}

func ExecCodeModule(name string, code Object) (*Module, error) {
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	ret := C.PyImport_ExecCodeModule(s, c(code))
	if ret == nil {
		return nil, exception()
	}
	return newModule(ret), nil
}

func NewModule(name string) (*Module, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_New(cname)
	if ret == nil {
		return nil, exception()
	}

	return newModule(ret), nil
}

func (m *Module) Register() error {
	name, err := m.Name()
	if err != nil {
		return err
	}
	if parent := getParentName(name); parent != "" {
		if pMod := getImport(parent); pMod == nil {
			return fmt.Errorf("parent module '%s' isn't registered", parent)
		} else if !pMod.isPackage() {
			return fmt.Errorf("parent module '%s' is not a package", parent)
		}
	}
	addImport(name, m)
	return nil
}

func (m *Module) CheckExact() bool {
	return C.moduleCheckE(c(m)) != 0
}

func (m *Module) Dict() *Dict {
	ret := C.PyModule_GetDict(c(m))
	return newDict(ret)
}

func (m *Module) Name() (string, error) {
	ret := C.PyModule_GetName(c(m))
	if ret == nil {
		return "", exception()
	}
	return C.GoString(ret), nil
}

func (m *Module) isPackage() bool {
	ret, err := m.GetAttrString("__path__")
	return ret != None && err == nil
}

// GetAttr is a convenience wrapper that is equivalent to
// mod.Base().GetAttr(name).
//
// Return value: New Reference.
func (m *Module) GetAttr(name Object) (Object, error) {
	return m.Base().GetAttr(name)
}

// GetAttrString is a convenience wrapper that is equivalent to
// mod.Base().GetAttrString(name).
//
// Return value: New Reference.
func (m *Module) GetAttrString(name string) (Object, error) {
	return m.Base().GetAttrString(name)
}

// TODO(jp3): PyModule_GetFilenameObject is the new API
// func (mod *Module) Filename() (string, error) {
// 	ret := C.PyModule_GetFilename(c(mod))
// 	if ret == nil {
// 		return "", exception()
// 	}
// 	return C.GoString(ret), nil
// }

func (m *Module) AddObject(name string, obj Object) error {
	if obj == nil {
		return AssertionError.Err("ValueError: obj == nil!")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddObject(c(m), cname, c(obj))
	return int2Err(ret)
}

func (m *Module) AddObjectRef(name string, obj Object) error {
	if obj == nil {
		return AssertionError.Err("ValueError: obj == nil!")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddObjectRef(c(m), cname, c(obj))
	return int2Err(ret)
}

func (m *Module) AddIntConstant(name string, value int) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddIntConstant(c(m), cname, C.long(value))
	if ret < 0 {
		return exception()
	}

	return nil
}

func (m *Module) AddStringConstant(name, value string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	ret := C.PyModule_AddStringConstant(c(m), cname, cvalue)
	if ret < 0 {
		return exception()
	}

	return nil
}

func (m *Module) AddType(t *Type) error {
	ret := C.PyModule_AddType(c(m), t.c())
	return int2Err(ret)
}

func getParentName(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[:i]
	}
	return ""
}

func InitExtension(f func() (*Module, error)) *BaseObject {
	ret, err := f()
	if err != nil {
		raise(err)
		return nil
	}

	return ret.Base()
}
