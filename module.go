// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"log"
	"sync"
	"unsafe"
)

type Module struct {
	AbstractObject
	o C.PyObject
}

// ModuleType is the Type object that represents the Module type.
var ModuleType = (*Type)(unsafe.Pointer(&C.PyModule_Type))

func moduleCheck(obj Object) bool {
	if obj == nil {
		return false
	}
	return C.moduleCheck(c(obj)) != 0
}

func newModule(obj *C.PyObject) *Module {
	return (*Module)(unsafe.Pointer(obj))
}

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

	if len(md.Methods) == 0 {
		return newModule(m), nil
	}

	n := C.PyUnicode_FromString(pyMD.m_name)
	if n == nil {
		return nil, exception()
	}

	d := C.PyModule_GetDict(m)
	if d == nil {
		return nil, exception()
	}

	for _, method := range md.Methods {
		pyF, err := makeCFunction(method.Name, method.Func, method.Doc, n)
		if err != nil {
			return nil, err
		}

		if C.PyDict_SetItemString(d, C.CString(method.Name), c(pyF)) != 0 {
			return nil, exception()
		}
	}

	return newModule(m), nil
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
	if ret < 0 {
		return exception()
	}

	return nil
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

var (
	importLock   sync.Mutex
	importMap    = map[string]*Module{}
	importOrigin *Unicode
)

func addImport(name string, mod *Module) {
	importLock.Lock()
	defer importLock.Unlock()

	importMap[name] = mod
}

func getImport(name string) *Module {
	importLock.Lock()
	defer importLock.Unlock()

	return importMap[name]
}

func importerFindSpec(args *Tuple) (Object, error) {
	var name string
	var path, target Object
	ParseTuple(args, "sO|O", &name, &path, &target)
	log.Printf("importerFindSpec: %p %s %T %T", args, name, path, target)

	if path != None {
		// we don't support lookups with path - basically this is a variant of
		// BuiltinImporter
		None.Incref()
		return None, nil
	}

	mod := getImport(name)
	if mod == nil {
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

	log.Printf("SFL: %p %T", sfl, sfl)

	return sfl.Base().CallFunction("sOO", name, mod, importOrigin)
}

// importerExecModule starts the setup of a module.
//
// Shamelessly copied from _imp_create_builtin and create_builtin from import.c.
func importerCreateModule(args *Tuple) (Object, error) {
	var spec Object
	ParseTuple(args, "O", &spec)
	log.Printf("importerCreateModule: %p %T", args, spec)

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
	log.Printf("importerExecModule: %p %T", args, mod)

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
		"find_spec":     importerFindSpec,
		"create_module": importerCreateModule,
		"exec_module":   importerExecModule,
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

	log.Printf("sys.meta_path: %d items", metaPath.Size())

	for i, obj := range metaPath.Slice() {
		r, err := obj.Base().Repr()
		if err != nil {
			return err
		}
		log.Printf("  %d: %T %v", i, r, r)
	}

	if err := importer.Create(); err != nil {
		return err
	}

	return metaPath.Insert(0, &importer)
}
