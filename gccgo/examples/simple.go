package simple

import (
	"fmt"
	"gopy"
)

func example(args *py.Tuple) (py.Object, error) {
	fmt.Printf("simple.example: %v\n", args)
	py.None.Incref()
	return py.None, nil
}

func init() {
	methods := []py.Method{
		{"example", example, "example function"},
	}

	_, err := py.InitModule("simple", methods)
	if err != nil {
		panic(err)
	}
}
