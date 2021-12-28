# Gollections

Golang experiments with container data structures, slices, maps, and so on using generics.
Need use version 1.18 beta 1 or newer.

## Containers

### [Immutable](./immutable/api.go)


[OrderedMap](./immutable/map.go)

[OrderedSet](./immutable/set.go)


### [Mutable](./mutable/api.go)

[OrderedMap](./mutable/map.go)

[OrderedSet](./mutable/set.go)


## Packages
### [Interfaces](./typ/iface.go)
```go
//Common interfaces


type Iterator[T any] interface {
	//checks ability on next element
	HasNext() bool
	//retrieves next element
	Get() T
}

type Iterable[T any, It Iterator[T]] interface {
	Begin() It
}

type Walker[P any, T any] func(position P, value T)
type Walk[T any, P any] interface {
	ForEach(Walker[P, T])
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
```

### [Iterator](./iter/)
API over 'Iterator' to make map, filter, flat, reduce operations in declarative style. 

Consists of two groups of operations:
 * Intermediate - only defines computation (Wrap,Map, Flatt, Filter).
 * Final - applies intermediate links (ToSlice, Reudce)
  

### [Slice](./slice/)
Same as [iter](./iter/) but specifically for slices. Generally more performant than [iter](./iter/) but only as the first in a chain of intermediate operations.



## Examples
```go
import (
	"testing"

	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/iter"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/stretchr/testify/assert"
)


func Test_compute_odds_sum(t *testing.T) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
	)

	expected := 1 + 3 + 5 + 7

    //declarative style
	oddSum := iter.Reduce(iter.Filter(iter.Flatt(slice.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds), op.Sum[int])
	assert.Equal(t, expected, oddSum)

	//plain old style
	oddSum = 0
	for _, i := range multiDimension {
		for _, ii := range i {
			for _, iii := range ii {
				if odds(iii) {
					oddSum += iii
				}
			}
		}
	}

	assert.Equal(t, expected, oddSum)
}
```