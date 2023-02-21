package py

type Iterable interface {
	Object
	Iter() (Iterator, error)
}

type Hashable interface {
	Object
	Hash() (Object, error)
}

type Callable interface {
	Object
	Call(*Tuple, *Dict) (Object, error)
	CallGo([]Object, map[string]Object) (Object, error)
}

type Comparable interface {
	Object
	RichCompare(Object, Op) (Object, error)
}
