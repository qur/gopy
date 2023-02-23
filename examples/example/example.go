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

func iterate(args *py.Tuple) (py.Object, error) {
	var o py.Object
	if err := py.ParseTuple(args, "O", &o); err != nil {
		return nil, err
	}
	fmt.Printf("item: %T (%s)\n", o, o.Type().String())
	i, err := py.GetIterator(o)
	if err != nil {
		return nil, err
	}
	defer i.Decref()
	fmt.Printf("iterator: %T (%s)\n", i, i.Type().String())
	items, err := py.Iterate(i)
	// iterate through items first, to make sure decref gets called
	for i, item := range items {
		fmt.Printf("%d: %s (%T)\n", i, item, item)
		item.Decref()
	}
	if err != nil {
		return nil, err
	}
	py.None.Incref()
	return py.None, nil
}

func msg(args *py.Tuple) (py.Object, error) {
	builtins, err := py.Import("builtins")
	if err != nil {
		return nil, err
	}
	defer builtins.Decref()

	print, err := builtins.GetAttrString("print")
	if err != nil {
		return nil, err
	}
	defer print.Decref()

	log.Printf("print: %T", print)

	sys, err := py.Import("sys")
	if err != nil {
		return nil, err
	}
	defer sys.Decref()

	stderr, err := sys.GetAttrString("stderr")
	if err != nil {
		return nil, err
	}
	defer stderr.Decref()

	return print.Base().CallGo(py.A{"this is a message"}, py.K{"file": stderr})
}

type Example struct {
	py.ClassBaseObject
	wibble int64
	Foo    int       `py:"foo,ro"`
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

func (e *Example) PyRichCompare(obj py.Object, op py.Op) (py.Object, error) {
	o, ok := obj.(*Example)
	if !ok {
		return nil, py.TypeError.Err("not a example.ExampleClass instance")
	}
	return py.RichCompareNative(e.wibble, o.wibble, op)
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

func (e *Example) Py_doot(arg py.Object) (py.Object, error) {
	log.Printf("arg: %T", arg)
	l, ok := arg.(*py.Long)
	if !ok {
		return nil, py.TypeError.Err("expected int, not %s", arg.Type())
	}
	e.Foo = int(l.Int64())
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
	Flags:  py.ClassBaseType,
	Object: (*Example)(nil),
}

var modDef = py.ModuleDef{
	Name:    "example",
	Package: true,
	Methods: []py.GoMethod{
		{"example", example, "example function"},
		{"iterate", iterate, "iterate any iterable"},
		{"msg", msg, "print messages using Python's print function"},
	},
}

var subDef = py.ModuleDef{
	Name: "example.sub",
}

func main() {
	// We don't want to call unlock, as Python will not be initialised after
	// Main returns, and that will cause Unlock to panic. For the same reason,
	// we don't defer any Decref calls here.
	py.InitAndLockWithSignals()

	m, err := py.CreateModule(&modDef)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

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

	if err := s.Register(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	if _, err := py.InitGoModule(); err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	py.Main(os.Args)
}
