package parallel

import (
	"gopy"
	"time"
)

func updateList(list *py.List) {
	lock := py.NewLock()
	defer lock.Unlock()

	for i := int64(1); true; i++ {
		lock.UnblockThreads()
		time.Sleep(time.Second)
		lock.BlockThreads()

		p := py.NewInt64(i)
		list.Append(p)
		p.Decref()
	}
}

func example(args *py.Tuple) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var o py.Object

	err := py.ParseTuple(args, "O", &o)
	if err != nil {
		return nil, err
	}

	l, ok := o.(*py.List)
	if !ok {
		return nil, py.TypeError.Err("Expected *py.List, not %T", o)
	}

	go updateList(l)

	py.None.Incref()
	return py.None, nil
}

func init() {
	methods := []py.Method{
		{"example", example, "example function"},
	}

	_, err := py.InitModule("parallel", methods)
	if err != nil {
		panic(err)
	}
}
