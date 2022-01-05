# Gollections

Golang generic container data structures and functions.

Need use Go version 1.18 beta 1 or newer.

## [Mutable containers](./mutable/api.go)

Supports write operations (append, delete, replace).

  * [Vector](./mutable/vector/api.go) - the simplest based on slices collection.
  * [OrderedSet](./mutable/set/api.go) - collection of unique items, prevents duplicates, provides iteration in order of addition.
  * [OrderedMap](./mutable/dict/api.go) - same as embedded map, but supports iteration in the order in which elements are added.

## [Immutable containers](./immutable/api.go) 

The same interfaces as in the mutable package but for read-only purposes.


## Functions

Consists of two groups of operations:
 * Intermediate - only defines a computation (Wrap, Map, Flatt, Filter).
 * Final - applies intermediates and retrieves the result (ForEach, Slice, Reudce)

Intermediates should wrap one by one to make a lazy computation chain that can be applied to the latest final operation.

```go
//Example 'filter', 'map', 'reduce' for an iterative container of 'items'


var items immutable.Vector[Item]

func condition(item Item) = ...
func max(attribute1 Attribute, attribute2 Attribute) = ... 

maxItemAttribute := it.Reduce(it.Map(c.Filer(items, condition), Item.GetAttribute), max)
```
Functions grouped into packages by applicable type ([iterable container](./c/api.go), [iterator](./it/api.go), [walker](walk/api.go), [slice](slice/api.go))

## Packages
### [Common interfaces](./typ/iface.go)

Iterator, Iterable, Container, Vector, Map, Set and so on.

### [Iterable container API](./c/api.go)
Declarative style API over 'Iterable' interface. Based on 'Iterator API' (see below).

### [Iterator API](./it/api.go)
Declarative style API over 'Iterator' interface. 

### [Slice API](./slice/api.go)
Same as 'Iterator API' but specifically for slices. Generally more performant than 'Iterator' but only as the first in a chain of intermediate operations.



## Examples
```go
package examples

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/immutable/set"
	"github.com/m4gshm/container/it"
	"github.com/m4gshm/container/op"
	slc "github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/walk"
	"github.com/stretchr/testify/assert"
)

func Test_OrderedSet(t *testing.T) {
	s := set.Of(1, 1, 2, 4, 3, 1)
	values := s.Elements()
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