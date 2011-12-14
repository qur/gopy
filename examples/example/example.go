// Copyright 2011 Julian Phillips.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/qur/gopy/lib"
	"os"
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
	py.None.Incref()
	ret, err := py.BuildValue("siO", s, i, l)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type ExampleClass struct {
	py.BaseObject
	wibble int
}

func (e *ExampleClass) PyInit(args *py.Tuple, kwds *py.Dict) error {
	fmt.Printf("ExampleClass.PyInit: args=%v, kwds=%v\n", args, kwds)
	return nil
}

func (e *ExampleClass) PyRepr() string {
	return fmt.Sprintf("<example.ExampleClass wibble=%d>", e.wibble)
}

func (e *ExampleClass) PyStr() string {
	return fmt.Sprintf("example.ExampleClass(wibble=%d)", e.wibble)
}

func (e *ExampleClass) PyCall(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	fmt.Printf("ExampleClass.PyCall(args=%v, kwds=%v)\n", args, kwds)
	py.None.Incref()
	return py.None, nil
}

func (e *ExampleClass) PyCompare(obj py.Object) (int, error) {
	o, ok := obj.(*ExampleClass)
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

func (e *ExampleClass) Py_bar(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
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
			case *py.Int:
				fmt.Printf("%s is %d\n", k, value.Int())
			}
		}
	}
	py.None.Incref()
	return py.None, nil
}

func (e *ExampleClass) PyGet_wibble() (py.Object, error) {
	i := py.NewInt(e.wibble)
	if i == nil {
		return nil, fmt.Errorf("TypeError: failed to convert wibble to Int")
	}
	return i, nil
}

func (e *ExampleClass) PySet_wibble(arg py.Object) error {
	i, ok := arg.(*py.Int)
	if !ok {
		return fmt.Errorf("TypeError: need a *py.Int, not %T", arg)
	}
	e.wibble = i.Int()
	return nil
}

var exampleClass = py.Class{
	Name:    "example.ExampleClass",
	Pointer: (*ExampleClass)(nil),
}

func main() {
	py.Initialize()

	methods := []py.Method{
		{"example", example, "example function"},
	}

	m, err := py.InitModule("example", methods)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	c, err := exampleClass.Create()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	err = m.AddObject("ExampleClass", c)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	py.Main(os.Args)
}
