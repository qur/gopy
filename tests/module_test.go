package pytesting

import (
	"fmt"
	"testing"

	"gopython.xyz/py/v3"
)

func TestFunction(t *testing.T) {
	lock := py.InitAndLock()
	defer lock.Finalize()

	called := false
	f := func() (py.Object, error) {
		called = true
		return py.None, nil
	}

	modDef := py.ModuleDef{
		Name: "test",
		Methods: []py.GoMethod{
			{"test", f, ""},
		},
	}

	if m, err := py.CreateModule(&modDef); err != nil {
		t.Fatal(err)
	} else if t2, err := m.GetAttrString("test"); err != nil {
		t.Fatal(err)
	} else {
		t2.Base().CallObject(nil)
	}
	if !called {
		t.Error("Function wasn't called")
	}
}

type ExampleClass struct {
	py.ClassBaseObject
	called bool
}

func (e *ExampleClass) Py_Test() (py.Object, error) {
	panic("called")
}

func (e *ExampleClass) Py_Test2(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	if v, err := args.GetIndex(0); err != nil {
		panic(err)
	} else if i, ok := v.(*py.Long); !ok {
		panic(v)
	} else if i.Int64() != 10 {
		panic(i)
	}
	panic("called2")
}

func (e *ExampleClass) Str() (py.Object, error) {
	panic("strcalled")
}

var exampleClass = py.Class{
	Name:   "test.test",
	Object: &ExampleClass{},
}

func TestMethod(t *testing.T) {
	lock := py.InitAndLock()
	defer lock.Finalize()

	modDef := py.ModuleDef{
		Name: "test",
	}

	if main, err := py.NewDict(); err != nil {
		t.Fatal(err)
	} else if m, err := py.CreateModule(&modDef); err != nil {
		t.Fatal(err)
	} else if err := m.Register(); err != nil {
		t.Fatal(err)
	} else if err := exampleClass.Create(); err != nil {
		t.Fatal(err)
	} else if g, err := py.GetBuiltins(); err != nil {
		t.Fatal(err)
	} else if err := main.SetItemString("__builtins__", g); err != nil {
		t.Fatal(err)
	} else if err := m.AddObjectRef("test", &exampleClass); err != nil {
		t.Fatal(err)
	} else if _, err := py.RunString("import test; a = test.test()", py.SingleInput, main, nil); err != nil {
		t.Fatal(err)
	} else if a, err := main.GetItemString("a"); err != nil {
		t.Fatal(err)
	} else if a == py.None || a.Type().String() != "<class 'test.test'>" {
		t.Error(a.Type().String())
	} else if _, ok := a.(*ExampleClass); !ok {
		t.Error(fmt.Sprintf("wanted *ExampleClass, got %T", a))
	}
}

func TestMethod2(t *testing.T) {
	lock := py.InitAndLock()
	defer lock.Finalize()

	modDef := py.ModuleDef{
		Name: "test",
	}

	if main, err := py.NewDict(); err != nil {
		t.Fatal(err)
	} else if m, err := py.CreateModule(&modDef); err != nil {
		t.Fatal(err)
	} else if err := m.Register(); err != nil {
		t.Fatal(err)
	} else if err := exampleClass.Create(); err != nil {
		t.Fatal(err)
	} else if g, err := py.GetBuiltins(); err != nil {
		t.Fatal(err)
	} else if err := main.SetItemString("__builtins__", g); err != nil {
		t.Fatal(err)
	} else if err := m.AddObjectRef("test", &exampleClass); err != nil {
		t.Fatal(err)
	} else if _, err := py.RunString("import test; a = test.test()", py.SingleInput, main, nil); err != nil {
		t.Fatal(err)
	} else if a, err := main.GetItemString("a"); err != nil {
		t.Fatal(err)
	} else {
		type Test struct {
			m    string
			pan  string
			f    string
			args []interface{}
		}
		tests := []Test{
			{"Test", "called", "", nil},
			{"Test2", "called2", "i", []interface{}{10}},
			{"__str__", "strcalled", "", nil},
		}
		// t.Run uses goroutines, so we need to allow other goroutines to grab
		// the GIL ...
		lock.UnblockThreads()
		defer lock.BlockThreads()
		for _, test := range tests {
			t.Run(test.m, func(t *testing.T) {
				// make sure that we have the GIL before doing anything else.
				lock := py.NewLock()
				defer lock.Unlock()
				defer func() {
					if i := recover(); i != test.pan {
						t.Error("Panicked for some other reason:", i)
					}
				}()
				a.Base().CallMethod(test.m, test.f, test.args...)
			})
		}
	}
}
