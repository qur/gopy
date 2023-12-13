package main

// #include <python3.12/Python.h>
import "C"

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

type typeSettings struct {
	Doc      string
	Check    bool
	New      bool
	Type     bool
	Immortal bool
}

func doExamine(v reflect.Value, prefix string, funcs map[string]bool) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		f := t.Field(i)
		funcs[f.Name] = !fv.IsZero()
		// special case for being set to _PyObject_NextNotImplemented, as
		// PyIter_Check doesn't consider that to be a "valid" value
		if f.Type.Kind() == reflect.Func && !fv.IsNil() {
			switch fv.UnsafePointer() {
			case C._PyObject_NextNotImplemented, C.PyObject_HashNotImplemented:
				funcs[f.Name] = false
			}
		}
		if !fv.IsZero() && strings.HasPrefix(f.Name, "tp_as_") {
			doExamine(fv.Elem(), f.Name+".", funcs)
		}
	}
}

func examine(value any) map[string]bool {
	v := reflect.ValueOf(value)
	funcs := map[string]bool{}
	doExamine(v, "", funcs)
	return funcs
}

func shortName(name string) string {
	n := strings.ToLower(name[:1])
	if n == "c" {
		// using c would conflict with the the c() function that converts Object
		// to *C.PyObject.
		return strings.ToLower(name[:2])
	}
	return n
}

func generateBasic(name string, settings typeSettings) {
	path := fmt.Sprintf("%s_gen.go", strings.ToLower(name))
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file %s: %s", path, err)
	}
	defer f.Close()
	if err := code.Execute(f, map[string]any{
		"type":     name,
		"ctype":    "PyObject",
		"ltype":    strings.ToLower(name[:1]) + name[1:],
		"name":     shortName(name),
		"funcs":    map[string]bool{},
		"settings": settings,
	}); err != nil {
		log.Fatalf("Failed to generate template: %s", err)
	}
}

func generate(name string, funcs map[string]bool) {
	path := fmt.Sprintf("%s_gen.go", strings.ToLower(name))
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file %s: %s", path, err)
	}
	defer f.Close()
	if err := code.Execute(f, map[string]any{
		"type":  name,
		"ctype": fmt.Sprintf("Py%sObject", name),
		"ltype": strings.ToLower(name[:1]) + name[1:],
		"name":  shortName(name),
		"funcs": funcs,
		"settings": typeSettings{
			Check:    true,
			New:      true,
			Type:     true,
			Immortal: false,
		},
	}); err != nil {
		log.Fatalf("Failed to generate template: %s", err)
	}
}

func main() {
	// Generate minimal code for basic types
	for name, settings := range basic {
		generateBasic(name, settings)
	}

	// Generate full types by examining the types
	for name, t := range types {
		generate(name, examine(t))
	}
}
