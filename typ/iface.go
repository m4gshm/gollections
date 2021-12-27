package typ

import "constraints"

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	//checks ability on next element
	HasNext() bool
	//retrieves next element
	Get() T
}

type Walk[T any] interface {
	Begin() Iterator[T]
}

type Container[T any, S constraints.Integer] interface {
	Values() []T
	Len() S
}

type Resetable interface {
	Reset()
}
