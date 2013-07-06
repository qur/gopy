package rgo

import (
	"fmt"
	"gopy"
)

var c = make(chan int)

func s() {
	for i := 0; true; i++ {
		c <- i
	}
}

func one(args *py.Tuple) (py.Object, error) {
	i, err := args.GetItem(0)
	if err != nil {
		return nil, err
	}

	fn, err := args.GetItem(1)
	if err != nil {
		return nil, err
	}

	arg, err := args.GetItem(2)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s + go: %s %d\n", i, arg, <- c)

	bfn := fn.Base()
	ret, err := bfn.CallFunctionObjArgs(i, arg)
	if err != nil {
		return nil, err
	}
	ret.Decref()

	fmt.Printf("%s - go: %s\n", i, arg)

	py.None.Incref()
	return py.None, nil
}

func init() {
	methods := []py.Method{
		{"one", one, "a function"},
	}

	_, err := py.InitModule("rgo", methods)
	if err != nil {
		panic(err)
	}

	go s()
}
