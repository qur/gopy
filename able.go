package py

type Iterable interface {
	Object
	Next() (Object, error)
}

type Hashable interface {
	Object
	Hash() (Object, error)
}
