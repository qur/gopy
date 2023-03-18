package main

import (
	"log"
	"os"
	"path/filepath"

	"gopython.xyz/py/v3"
)

type Wibble struct {
	py.ClassBaseObject
}

func (w *Wibble) Py_do(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	log.Printf("GO DO: %s %s", args, kwds)
	super, err := wibbleType.Super(w)
	if err != nil {
		return nil, err
	}
	return super.CallMethod("do", args, kwds)
}

var wibbleType = py.Class{
	Name:   "Wibble",
	Flags:  py.ClassBaseType | py.ClassHeapType,
	Object: (*Wibble)(nil),
}

func main() {
	lock := py.InitAndLock()
	defer lock.Finalize()

	base, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("failed to get base directory: %s", err)
	}
	code := filepath.Join(base, "code")

	log.Printf("code: %s", code)

	if err := py.AddToPath(code); err != nil {
		log.Fatalf("failed to prepend to path: %s", err)
	}

	m, err := py.Import("loadertest")
	log.Printf("IMPORT: %p %v", m, err)

	if err != nil {
		log.Fatalf("failed to import module: %s", err)
	}

	for k := range m.Dict().Map() {
		log.Printf("%s", k)
	}

	w, err := m.GetAttrString("Wibble")
	if err != nil {
		log.Fatalf("failed to get loadertest.Wibble: %s", err)
	}

	wibbleType.BaseType = w

	if err := wibbleType.Create(); err != nil {
		log.Fatalf("failed to create subclass of loadertest.Wibble: %s", err)
	}

	foobar, err := m.GetAttrString("foobar")
	if err != nil {
		log.Fatalf("failed to get loadertest.foobar: %s", err)
	}

	o, err := foobar.Base().CallGo(py.A{&wibbleType}, nil)
	if err != nil {
		log.Fatalf("failed to call loadertest.foobar: %s", err)
	}
	py.Decref(o)
}
