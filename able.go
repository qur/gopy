package py

type Iterable interface {
	Object
	Next() (Object, error)
}

type Hashable interface {
	Object
	Hash() (Object, error)
}

type Callable interface {
	Object
	Call(*Tuple, *Dict) (Object, error)
	CallGo(args []Object, kwds map[string]Object) (Object, error)
}
