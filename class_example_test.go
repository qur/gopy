package py_test

import (
	"fmt"

	"gopython.xyz/py/v3"
)

// myClass is the type that represents the custom class in Go.
type myClass struct {
	py.ClassBaseObject
}

// Call is the method called when the type is called in Python.
//
// This is one of many optional functions that can be implemented on the custom
// class type to respond to actions in Python.
func (m *myClass) Call(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	var s string
	if err := py.ParseTuple(args, "s", &s); err != nil {
		return nil, err
	}

	fmt.Printf("hello %s", s)

	return py.ReturnNone(), nil
}

// myClassType is a py.Class value that is the Type of myClass.
var myClassType = py.Class{
	Name:   "MyClass",
	Flags:  py.ClassBaseType,
	Doc:    "An example Class",
	Object: (*myClass)(nil),
}

func ExampleClass() {
	lock := py.InitAndLock()
	defer lock.Finalize()

	// we have to Create the Class before we can use it
	if err := myClassType.Create(); err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}

	// Python: "m = MyClass()"
	m, err := myClassType.CallGo(nil, nil)
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
