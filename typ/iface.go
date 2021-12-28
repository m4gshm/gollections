package typ

import "constraints"

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	//checks ability on next element
	HasNext() bool
	//retrieves next element
	Get() T
}

type Walk[T any, It Iterator[T]] interface {
	Begin() It
}

type Container[T any] interface {
	Values() T
}

type Measureable[L constraints.Integer] interface {
	Len() L
}

type Checkable[T any] interface {
	Contains(T) bool
}

type Resetable interface {
	Reset()
}

type Appendable[T any] interface {
	Add(T) bool
}

type Deletable[T any] interface {
	Delete(T) bool
}

type Access[K any, V any] interface {
	Get(K) V
}

type RandomAccess[K constraints.Integer, V any] interface {
	Access[K, V]
}

type KeyAccess[K comparable, V any] interface {
	Access[K, V]
}
