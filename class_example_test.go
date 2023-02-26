package py_test

import (
	"fmt"

	"gopython.xyz/py/v3"
)

type myClass struct {
	py.ClassBaseObject
}

func (m *myClass) Call(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	var s string
	if err := py.ParseTuple(args, "s", &s); err != nil {
		return nil, err
	}

	fmt.Printf("hello %s", s)

	return py.ReturnNone(), nil
}

var myClassDef = py.Class{
	Name:   "MyClass",
	Flags:  py.ClassBaseType,
	Doc:    "An example Class",
	Object: (*myClass)(nil),
}

func ExampleClass() {
	lock := py.InitAndLock()
	defer lock.Finalize()

	if err := myClassDef.Create(); err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}

	// Python: "m = MyClass()"
	m, err := myClassDef.CallGo(nil, nil)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}
	defer m.Decref()

	// Python: "m("world")"
	o, err := m.Base().CallGo(py.A{"world"}, nil)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}
	defer o.Decref()

	// Output:
	// hello world
}
