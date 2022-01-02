# Gollections

Golang experiments with container data structures, slices, maps, and so on using generics.
Need use version 1.18 beta 1 or newer.

## Containers

### [Immutable](./immutable/api.go)

[Vector](./immutable/vector/api.go)

[OrderedMap](./immutable/map/api.go)

[OrderedSet](./immutable/set/api.go)


### [Mutable](./mutable/api.go)

[Vector](./mutable/vector/api.go)

[OrderedMap](./mutable/map/api.go)

[OrderedSet](./mutable/set/api.go)


## Packages
### [Interfaces](./typ/iface.go)

```go
/*
 *  Common interfaces
 */

type Container[T any, L constraints.Integer] interface {
	Walk[T]
	Finite[[]T, L]
}

type Vector[T any] interface {
	Container[T, int]
	RandomAccess[int, T]
}

type Set[T any] interface {
	Container[T, int]
	Checkable[T]
}

type Map[k comparable, v any, IT Iterator[*KV[k, v]]] interface {
	Track[v, k]
	Iterable[*KV[k, v], IT]
	Finite[map[k]v, int]
	Checkable[k]
	KeyAccess[k, v]
}

//Iterator base interface for containers, collections
type Iterator[T any] interface {
	//checks ability on next element or error
	HasNext() bool
	//retrieves next element
	Get() T
	//retrieves error
	Err() error
}

//Iterable iterator supplier
type Iterable[T any, It Iterator[T]] interface {
	Begin() It
}

//Walk touches all elements of the collection
type Walk[T any] interface {
	ForEach(func(element T))
}

//Track traverses container elements with position tracking (index, key, coordinates, etc.)
type Track[T any, P any] interface {
	ForEach(func(position P, element T))
}

//Checkable container with ability to check if an element is present
type Checkable[T any] interface {
	Contains(T) bool
}

//Finite not endless container that can be transformed to array or map of elements
type Finite[T any, L constraints.Integer] interface {
	Values() T
	Len() L
}

type Access[K any, V any] interface {
	Get(K) (V, bool)
}

type RandomAccess[K constraints.Integer, V any] interface {
	Access[K, V]
}

type KeyAccess[K comparable, V any] interface {
	Access[K, V]
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
package examples

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/immutable/set"
	it "github.com/m4gshm/container/iter"
	"github.com/m4gshm/container/op"
	slc "github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/walk"
	"github.com/stretchr/testify/assert"
)

func Test_OrderedSet(t *testing.T) {
	s := set.Of(1, 1, 2, 4, 3, 1)
	values := s.Values()
	fmt.Println(s) //[1, 2, 4, 3]

	assert.Equal(t, slc.Of(1, 2, 4, 3), values)
}

func Test_group_odd_even(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = walk.Group(set.Of(1, 1, 2, 4, 3, 1), even)
	)
	fmt.Println(groups) //map[false:[1 3] true:[2 4]]
	assert.Equal(t, map[bool][]int{false: {1, 3}, true: {2, 4}}, groups)
}

func Test_compute_odds_sum(t *testing.T) {
	var (
		odds           = func(v int) bool { return v%2 != 0 }
		multiDimension = [][][]int{{{1, 2, 3}, {4, 5, 6}}, {{7}, nil}, nil}
		expected       = 1 + 3 + 5 + 7
	)

	//declarative style
	oddSum := it.Reduce(it.Filter(it.Flatt(slc.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds), op.Sum[int])
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