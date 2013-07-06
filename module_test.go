package pytesting

import (
	"github.com/qur/gopy/lib"
	"testing"
)

func TestFunction(t *testing.T) {
	py.Initialize()
	defer py.Finalize()

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

func (e *ExampleClass) PyStr() string {
	panic("called")
	return "my test"
}

var exampleClass = py.Class{
	Name:    "test.test",
	Pointer: (*ExampleClass)(nil),
}

func TestMethod(t *testing.T) {
	py.Initialize()
	defer py.Finalize()

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
	defer py.Finalize()

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
	} else if _, err := py.RunString("import test; a = test.test().__str__", py.SingleInput, main, nil); err != nil {
		t.Fatal(err)
	} else if a, err := main.GetItemString("a"); err != nil {
		t.Fatal(err)
	} else {
		_ = a
		// Hangs forever...
		//		t.Log(a.Base().CallObject(nil))
	}
}
