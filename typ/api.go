package typ


//Iterator base interface for containers, collections
type Iterator[T any] interface {
	//checks ability on next element
	HasNext() bool
	//retrieves next element
	Get() T
}

type Resetable interface {
	Reset()
}


