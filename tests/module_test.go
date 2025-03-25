package pytesting

import (
	"fmt"
	"testing"

	"gopython.xyz/py/v3"
	"gopython.xyz/py/v3/tests/assert"
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
			{Name: "test", Func: f, Doc: ""},
		},
	}

	if m, err := py.CreateModule(&modDef); err != nil {
		t.Fatal(err)
	} else if t2, err := m.GetAttrString("test"); err != nil {
		t.Fatal(err)
	} else if _, err := t2.Base().CallObject(nil); err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("Function wasn't called")
	}
}

type ExampleClass struct {
	py.ClassBaseObject
}

func (e *ExampleClass) Test() (py.Object, error) {
	panic("called")
}

func (e *ExampleClass) Test2(args *py.Tuple, _ *py.Dict) (py.Object, error) {
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
	Methods: map[string]string{
		"Test":     "Test",
		"test_two": "Test2",
	},
}

func TestMethod(t *testing.T) {
	lock := py.InitAndLock()
	defer lock.Finalize()

	m, err := setupTestModule()
	assert.NilErr(t, err)

	defer m.Decref()

	main, err := setupMain()
	assert.NilErr(t, err)

	defer main.Decref()

	_, err = py.RunString("import test; a = test.test()", py.SingleInput, main, nil)
	assert.NilErr(t, err)

	a, err := main.GetItemString("a")
	assert.NilErr(t, err)

	if a == py.None || a.Type().String() != "<class 'test.test'>" {
		t.Error(a.Type().String())
	}

	if _, ok := a.(*ExampleClass); !ok {
		t.Errorf("wanted *ExampleClass, got %T", a)
	}
}

func TestMethod2(t *testing.T) {
	lock := py.InitAndLock()
	defer lock.Finalize()

	m, err := setupTestModule()
	assert.NilErr(t, err)

	defer m.Decref()

	main, err := setupMain()
	assert.NilErr(t, err)

	defer main.Decref()

	_, err = py.RunString("import test; a = test.test()", py.SingleInput, main, nil)
	assert.NilErr(t, err)

	a, err := main.GetItemString("a")
	assert.NilErr(t, err)

	defer a.Decref()

	type Test struct {
		m    string
		pan  string
		f    string
		args []interface{}
	}

	tests := []Test{
		{"Test", "called", "", nil},
		{"test_two", "called2", "i", []interface{}{10}},
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
				if i := recover(); i != nil && i != test.pan {
					t.Error("Panicked for some other reason:", i)
				}
			}()

			_, err := a.Base().CallMethod(test.m, test.f, test.args...)
			assert.NilErr(t, err)
		})
	}
}

func setupTestModule() (*py.Module, error) {
	modDef := py.ModuleDef{
		Name: "test",
	}

	m, err := py.CreateModule(&modDef)
	if err != nil {
		return nil, fmt.Errorf("failed to create module: %w", err)
	}

	if err := m.Register(); err != nil {
		return nil, fmt.Errorf("failed to register module: %w", err)
	}

	if err := exampleClass.Create(); err != nil {
		return nil, fmt.Errorf("failed to create exampleClass: %w", err)
	}

	if err := m.AddObjectRef("test", &exampleClass); err != nil {
		return nil, fmt.Errorf("failed to add exampleClass to module: %w", err)
	}

	return m, nil
}
