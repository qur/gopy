package pytesting

import (
	"fmt"
	"testing"

	"gopython.xyz/py/v14"
	"gopython.xyz/py/v14/tests/assert"
)

func TestRunString(t *testing.T) {
	lock := py.InitAndLock()
	defer lock.Finalize()

	main, err := setupMain()
	assert.NilErr(t, err)

	defer main.Decref()

	_, err = py.RunString("a = 'hello world!'", py.FileInput, main, nil)
	assert.NilErr(t, err)

	a, err := main.GetItemString("a")
	assert.NilErr(t, err)

	defer a.Decref()

	if b, ok := a.(*py.Unicode); !ok || b.String() != "hello world!" {
		t.Error(b, err)
	}
}

func setupMain() (*py.Dict, error) {
	main, err := py.NewDict()
	if err != nil {
		return nil, fmt.Errorf("failed to create main dict: %w", err)
	}

	g, err := py.GetBuiltins()
	if err != nil {
		return nil, fmt.Errorf("failed to get builtins: %w", err)
	}

	if err := main.SetItemString("__builtins__", g); err != nil {
		return nil, fmt.Errorf("failed to set __builtins__: %w", err)
	}

	return main, nil
}
