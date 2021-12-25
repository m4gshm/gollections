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

type Array interface {
	~[]any | ~[]uintptr |
		~[]int | ~[]int8 | []int16 | []int32 | []int64 |
		~[]uint | ~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 |
		~[]float32 | ~[]float64 |
		~[]complex64 | ~[]complex128 |
		~[]string | ~string
}
