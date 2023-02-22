package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"strings"
	"sync"
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

func (mod *Module) Register() error {
	name, err := mod.Name()
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
	addImport(name, mod)
	return nil
}

func (mod *Module) CheckExact() bool {
	return C.moduleCheckE(c(mod)) != 0
}

func (mod *Module) Dict() *Dict {
	ret := C.PyModule_GetDict(c(mod))
	return newDict(ret)
}

func (mod *Module) Name() (string, error) {
	ret := C.PyModule_GetName(c(mod))
	if ret == nil {
		return "", exception()
	}
	return C.GoString(ret), nil
}

func (mod *Module) isPackage() bool {
	ret, err := mod.GetAttrString("__path__")
	return ret != None && err == nil
}

// GetAttr is a convenience wrapper that is equivalent to
// mod.Base().GetAttr(name).
//
// Return value: New Reference.
func (mod *Module) GetAttr(name Object) (Object, error) {
	return mod.Base().GetAttr(name)
}

// GetAttrString is a convenience wrapper that is equivalent to
// mod.Base().GetAttrString(name).
//
// Return value: New Reference.
func (mod *Module) GetAttrString(name string) (Object, error) {
	return mod.Base().GetAttrString(name)
}

// TODO(jp3): PyModule_GetFilenameObject is the new API
// func (mod *Module) Filename() (string, error) {
// 	ret := C.PyModule_GetFilename(c(mod))
// 	if ret == nil {
// 		return "", exception()
// 	}
// 	return C.GoString(ret), nil
// }

func (mod *Module) AddObject(name string, obj Object) error {
	if obj == nil {
		return AssertionError.Err("ValueError: obj == nil!")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddObject(c(mod), cname, c(obj))
	return int2Err(ret)
}

func (mod *Module) AddObjectRef(name string, obj Object) error {
	if obj == nil {
		return AssertionError.Err("ValueError: obj == nil!")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddObjectRef(c(mod), cname, c(obj))
	return int2Err(ret)
}

func (mod *Module) AddIntConstant(name string, value int) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ret := C.PyModule_AddIntConstant(c(mod), cname, C.long(value))
	if ret < 0 {
		return exception()
	}

	return nil
}

func (mod *Module) AddStringConstant(name, value string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	ret := C.PyModule_AddStringConstant(c(mod), cname, cvalue)
	if ret < 0 {
		return exception()
	}

	return nil
}

func (mod *Module) AddType(t *Type) error {
	ret := C.PyModule_AddType(c(mod), t.c())
	return int2Err(ret)
}

func getParentName(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[:i]
	}
	return ""
}

var (
	importLock   sync.Mutex
	importMap    = map[string]*Module{}
	importOrigin *Unicode
)

func addImport(name string, mod *Module) {
	importLock.Lock()
	defer importLock.Unlock()

	mod.Incref()
	importMap[name] = mod
}

func getImport(name string) *Module {
	importLock.Lock()
	defer importLock.Unlock()

	return importMap[name]
}

func getParent(name string) *Module {
	if parent := getParentName(name); parent != "" {
		return getImport(parent)
	}
	return nil
}

func importerFindSpec(cls *Class, args *Tuple) (Object, error) {
	var name string
	var path, target Object
	ParseTuple(args, "sO|O", &name, &path, &target)

	// If this is a sub-package, we will only import it if we also own the
	// parent.
	if path != None && getParent(name) == nil {
		None.Incref()
		return None, nil
	}

	if getImport(name) == nil {
		// we don't have the requested module
		None.Incref()
		return None, nil
	}

	util, err := Import("importlib.util")
	if err != nil {
		return nil, err
	}

	sfl, err := util.GetAttrString("spec_from_loader")
	if err != nil {
		return nil, err
	}

	sflArgs, err := BuildValue("sO", name, cls)
	if err != nil {
		return nil, err
	}

	sflKwds, err := BuildValue("{sO}", "origin", importOrigin)
	if err != nil {
		return nil, err
	}

	return sfl.Base().Call(sflArgs.(*Tuple), sflKwds.(*Dict))
}

// importerExecModule starts the setup of a module.
//
// Shamelessly copied from _imp_create_builtin and create_builtin from import.c.
func importerCreateModule(args *Tuple) (Object, error) {
	var spec Object
	ParseTuple(args, "O", &spec)

	nameObj, err := spec.Base().GetAttrString("name")
	if err != nil {
		return nil, err
	}

	name, err := nameObj.(*Unicode).AsString()
	if err != nil {
		return nil, err
	}

	mod := getImport(name)
	if mod == nil {
		// we don't have the requested module
		None.Incref()
		return None, nil
	}

	return mod, nil
}

// importerExecModule completes the setup of a module.
//
// Shamelessly copied from exec_builtin_or_dynamic from import.c.
func importerExecModule(args *Tuple) (Object, error) {
	var mod Object
	ParseTuple(args, "O", &mod)

	if _, ok := mod.(*Module); !ok {
		// not actually a module, ignore it
		None.Incref()
		return None, nil
	}

	def := C.PyModule_GetDef(c(mod))
	if def == nil {
		None.Incref()
		return None, nil
	}

	state := C.PyModule_GetState(c(mod))
	if state != nil {
		// already initialised
		None.Incref()
		return None, nil
	}

	if C.PyModule_ExecDef(c(mod), def) < 0 {
		return nil, exception()
	}

	None.Incref()
	return None, nil
}

var importer = Class{
	Name: "GoImporter",
	Static: map[string]any{
		"create_module": importerCreateModule,
		"exec_module":   importerExecModule,
	},
	Class: map[string]any{
		"find_spec": importerFindSpec,
	},
}

func initModules() error {
	origin, err := NewUnicode("gopy")
	if err != nil {
		return err
	}
	importOrigin = origin

	sys, err := Import("sys")
	if err != nil {
		return err
	}

	metaPathObj, err := sys.GetAttrString("meta_path")
	if err != nil {
		return err
	}

	metaPath, ok := metaPathObj.(*List)
	if !ok {
		return fmt.Errorf("sys.meta_path should be list, got %T", metaPathObj)
	}

	if err := importer.Create(); err != nil {
		return err
	}

	return metaPath.Insert(0, &importer)
}

func InitExtension(f func() (*Module, error)) *BaseObject {
	ret, err := f()
	if err != nil {
		raise(err)
		return nil
	}

	return ret.Base()
}
