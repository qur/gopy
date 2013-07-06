package pytesting

import (
	"github.com/qur/gopy/lib"
	"testing"
)

func TestFunction(t *testing.T) {
	py.Initialize()

	called := false
	f := func() (py.Object, error) {
		called = true
		return py.None, nil
	}
	if m, err := py.InitModule("test", []py.Method{{"test", f, ""}}); err != nil {
		t.Fatal(err)
	} else if t2, err := m.Dict().GetItemString("test"); err != nil {
		t.Fatal(err)
	} else {
		t2.Base().CallObject(nil)
	}
	if !called {
		t.Error("Function wasn't called")
	}
}

type ExampleClass struct {
	py.BaseObject
	called bool
}

func (e *ExampleClass) Py_Test() (py.Object, error) {
	panic("called")
}

func (e *ExampleClass) Py_Test2(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	if v, err := args.GetItem(0); err != nil {
		panic(err)
	} else if i, ok := v.(*py.Int); !ok {
		panic(v)
	} else if i.Int() != 10 {
		panic(i)
	}
	panic("called2")
}

func (e *ExampleClass) PyStr() string {
	panic("strcalled")
}

var exampleClass = py.Class{
	Name:    "test.test",
	Pointer: &ExampleClass{},
}

func TestMethod(t *testing.T) {
	py.Initialize()

	if main, err := py.NewDict(); err != nil {
		t.Fatal(err)
	} else if m, err := py.InitModule("test", nil); err != nil {
		t.Fatal(err)
	} else if c, err := exampleClass.Create(); err != nil {
		t.Fatal(err)
	} else if g, err := py.GetBuiltins(); err != nil {
		t.Fatal(err)
	} else if err := main.SetItemString("__builtins__", g); err != nil {
		t.Fatal(err)
	} else if err := m.AddObject("test", c); err != nil {
		t.Fatal(err)
	} else if _, err := py.RunString("import test; a = test.test()", py.SingleInput, main, nil); err != nil {
		t.Fatal(err)
	} else if a, err := main.GetItemString("a"); err != nil {
		t.Fatal(err)
	} else if a == py.None || a.Type().String() != "<type 'test.test'>" {
		t.Error(a.Type().String())
	}
}

func TestMethod2(t *testing.T) {
	py.Initialize()

	if main, err := py.NewDict(); err != nil {
		t.Fatal(err)
	} else if m, err := py.InitModule("test", nil); err != nil {
		t.Fatal(err)
	} else if c, err := exampleClass.Create(); err != nil {
		t.Fatal(err)
	} else if g, err := py.GetBuiltins(); err != nil {
		t.Fatal(err)
	} else if err := main.SetItemString("__builtins__", g); err != nil {
		t.Fatal(err)
	} else if err := m.AddObject("test", c); err != nil {
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
		for _, test := range tests {
			func() {
				defer func() {
					if i := recover(); i != test.pan {
						t.Error("Paniced for some other reason:", i)
					}
				}()
				a.Base().CallMethod(test.m, test.f, test.args...)
			}()
		}
	}
}
