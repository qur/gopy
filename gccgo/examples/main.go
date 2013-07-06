package Main

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
	fmt.Printf("Main.init()\n")
}

func Main() {
	fmt.Printf("Main.Main()\n")
	lock := py.InitAndLock()
	defer lock.Unlock()

	fmt.Println("Locked.")

	m, err := py.Import("simple")
	if err != nil {
		panic(err)
	}

	fmt.Println("Imported")

	b := m.Base()

	r, err := b.CallMethod("example", "s", "hello there")
	if err != nil {
		panic(err)
	}

	fmt.Println(r)
}
