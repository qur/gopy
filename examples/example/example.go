package main

import (
	"fmt"
	"log"
	"os"

	"qur.me/py/v3"
)

func example(args *py.Tuple) (py.Object, error) {
	fmt.Printf("example.example\n")
	var s string
	var i int
	var o py.Object
	err := py.ParseTuple(args, "siO", &s, &i, &o)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return nil, err
	}
	l, ok := o.(*py.List)
	if !ok {
		return nil, fmt.Errorf("TypeError: wanted *py.List, not %T", o)
	}
	fmt.Printf("s: %s\n", s)
	fmt.Printf("i: %d\n", i)
	fmt.Printf("l: %v\n", l)
	for i, v := range l.Slice() {
		fmt.Printf("l[%d]: %v\n", i, v)
	}
	ret, err := py.BuildValue("siO", s, i, l)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type Example struct {
	py.ClassBaseObject
	wibble int64
	Foo    int       `py:"foo"`
	X      py.Object `py:"x"`
	Y      *py.Tuple `py:"y"`
}

func (e *Example) PyInit(args *py.Tuple, kwds *py.Dict) error {
	fmt.Printf("ExampleClass.PyInit: e=%p args=%v, kwds=%v\n", e, args, kwds)
	return nil
}

func (e *Example) PyRepr() string {
	return fmt.Sprintf("<example.ExampleClass self=%p wibble=%d>", e, e.wibble)
}

func (e *Example) PyStr() string {
	return fmt.Sprintf("example.ExampleClass(wibble=%d)", e.wibble)
}

func (e *Example) PyCall(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	fmt.Printf("ExampleClass.PyCall(self=%p, args=%v, kwds=%v)\n", e, args, kwds)
	py.None.Incref()
	return py.None, nil
}

func (e *Example) PyCompare(obj py.Object) (int, error) {
	o, ok := obj.(*Example)
	if !ok {
		return 0, fmt.Errorf("TypeError: not a example.ExampleClass instance")
	}
	if e.wibble < o.wibble {
		return -1, nil
	} else if e.wibble > o.wibble {
		return 1, nil
	}
	return 0, nil
}

func (e *Example) Py_bar(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	fmt.Printf("ExampleClass.Bar: %v %v %v\n", e, args, kwds)
	if kwds != nil {
		m, err := kwds.MapString()
		if err != nil {
			return nil, err
		}
		for k, v := range m {
			fmt.Printf("kwds[\"%v\"] = %v(%T)\n", k, v, v)
			switch value := v.(type) {
			case *py.Bool:
				if value.Bool() {
					fmt.Printf("%s it is then!\n", k)
				} else {
					fmt.Printf("No %s?  That's sad ...\n", k)
				}
			case *py.Long:
				fmt.Printf("%s is %d\n", k, value.Int64())
			}
		}
	}
	py.None.Incref()
	return py.None, nil
}

func (e *Example) PyGet_wibble() (py.Object, error) {
	i := py.NewLong(e.wibble)
	if i == nil {
		return nil, fmt.Errorf("TypeError: failed to convert wibble to Int")
	}
	return i, nil
}

func (e *Example) PySet_wibble(arg py.Object) error {
	i, ok := arg.(*py.Long)
	if !ok {
		return fmt.Errorf("TypeError: need a *py.Long, not %T", arg)
	}
	e.wibble = i.Int64()
	return nil
}

var exampleClass = py.Class{
	Name:   "example.ExampleClass",
	Object: (*Example)(nil),
}

var modDef = py.ModuleDef{
	Name:    "example",
	Package: true,
	Methods: []py.GoMethod{
		{"example", example, "example function"},
	},
}

var subDef = py.ModuleDef{
	Name: "example.sub",
}

func main() {
	// We don't want to call unlock, as Python will not be initialised after
	// Main returns, and that will cause Unlock to panic
	py.InitAndLockWithSignals()

	m, err := py.CreateModule(&modDef)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer m.Decref()

	if err := exampleClass.Create(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	if err = m.AddObjectRef("ExampleClass", &exampleClass); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	if err := m.Register(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	s, err := py.CreateModule(&subDef)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer s.Decref()

	if err := s.Register(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	if _, err := py.InitGoModule(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	py.Main(os.Args)
}
