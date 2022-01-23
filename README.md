# Gollections

Golang generic container data structures and functions.

Need use Go version 1.18 beta 1 or newer.

## [Mutable containers](./mutable/api.go)

Supports write operations (append, delete, replace).

  * [Vector](./mutable/vector/api.go) - the simplest based on built-in slice collection.
  * [Set](./mutable/set/api.go) - collection of unique items, prevents duplicates.
  * [Map](./mutable/map_/api.go) - built-in map wrapper that supports [container functions](#container-functions).
  * [OrderedSet](./mutable/oset/api.go) - collection of unique items, prevents duplicates, provides iteration in order of addition.
  * [OrderedMap](./mutable/omap/api.go) - same as the [Map](./mutable/map_/api.go), but supports iteration in the order in which elements are added.

## [Immutable containers](./immutable/api.go) 

The same interfaces as in the mutable package but for read-only purposes.

## Container functions

Consists of two groups of operations:
 * Intermediate - only defines a computation (Wrap, Map, Flatt, Filter, Group).
 * Final - applies intermediates and retrieves a result (ForEach, Slice, Reduce)

Intermediates should wrap one by one to make a lazy computation chain that can be applied to the latest final operation.

```go
//Example 'filter', 'map', 'reduce' for an iterative container of 'items'

var items immutable.Vector[Item]

var (
    condition typ.Predicate[Item]  = func(item Item) ...
    max       op.Binary[Attribute] = func(attribute1 Attribute, attribute2 Attribute) ...
) 

maxItemAttribute := it.Reduce(it.Map(c.Filer(items, condition), Item.GetAttribute), max)
```
Functions grouped into packages by applicable type ([container](./c/api.go), [iterator](./it/api.go), [slice](slice/api.go))

## Builder and util functions

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

	"github.com/stretchr/testify/assert"

	cgroup "github.com/m4gshm/gollections/c/group"
	"github.com/m4gshm/gollections/conv"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/it"
	slc "github.com/m4gshm/gollections/it/slice"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/walk/group"
)

func Test_Set(t *testing.T) {
	var (
		s      immutable.Set[int] = set.Of(1, 1, 2, 4, 3, 1)
		values []int              = s.Collect()
	)

	assert.Equal(t, 4, s.Len())
	assert.Equal(t, 4, len(values))

	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))
	assert.True(t, s.Contains(4))
	assert.False(t, s.Contains(5))
}

func Test_OrderedSet(t *testing.T) {
	s := oset.Of(1, 1, 2, 4, 3, 1)
	values := s.Collect()
	fmt.Println(s) //[1, 2, 4, 3]

	assert.Equal(t, slice.Of(1, 2, 4, 3), values)
}

func Test_group_orderset_odd_even(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = group.Of(oset.Of(1, 1, 2, 4, 3, 1), even)
	)
	fmt.Println(groups) //map[false:[1 3] true:[2 4]]
	assert.Equal(t, map[bool][]int{false: {1, 3}, true: {2, 4}}, groups)
}

func Test_group_orderset_with_filtering_by_stirng_len(t *testing.T) {
	var groups = cgroup.Of(oset.Of(
		"seventh", "seventh", //duplicated
		"first", "second", "third", "fourth",
		"fifth", "sixth", "eighth",
		"ninth", "tenth", "one", "two", "three", "1",
		"second", //duplicate
	), func(v string) int { return len(v) },
	).FilterKey(
		func(k int) bool { return k > 3 },
	).MapValue(
		func(v string) string { return v + "_" },
	).Collect()

	fmt.Println(groups) //map[int][]string{5:[]string{"first_", "third_", "fifth_", "sixth_", "ninth_", "tenth_", "three_"}, 6:[]string{"second_", "fourth_", "eighth_"}, 7:[]string{"seventh_"}}

	assert.Equal(t, map[int][]string{
		5: {"first_", "third_", "fifth_", "sixth_", "ninth_", "tenth_", "three_"},
		6: {"second_", "fourth_", "eighth_"},
		7: {"seventh_"},
	}, groups)
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