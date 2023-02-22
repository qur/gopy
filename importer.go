package py

// #include "utils.h"
import "C"

import (
	"fmt"
	"sync"
)

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

	return sfl.Base().CallGo(A{name, cls}, K{"origin": importOrigin})
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

func setupImporter() error {
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
