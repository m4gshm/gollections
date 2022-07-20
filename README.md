# Gollections

Golang generic containers and functions.

Supports Go version 1.18 or higher.

## Installation

```bash
go get -u github.com/m4gshm/gollections
```

## [Slice API](./slice/api.go)

The package provides various functions useful with slices.
Additionally contains subpackages as shortcuts of some functions from the base package.

Just looks some examples:

```go
package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/ordered"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/slice/sort"
	"github.com/m4gshm/gollections/sum"
)

func Test_SortInt(t *testing.T) {
	c := ordered.Sort([]int{-1, 0, 1, 2, 3})
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, c)
}

func Test_SortStructs(t *testing.T) {
	type User struct {
		name string
		age  int
	}
	var (
		users  = []User{{"Bob", 26}, {"Alice", 35}, {"Tom", 18}}
		byName = sort.ByOrdered(clone.Of(users), func(u User) string { return u.name })
		byAge  = sort.ByOrdered(clone.Of(users), func(u User) int { return u.age })
	)
	assert.Equal(t, []User{{"Alice", 35}, {"Bob", 26}, {"Tom", 18}}, byName)
	assert.Equal(t, []User{{"Tom", 18}, {"Bob", 26}, {"Alice", 35}}, byAge)
}

func Test_Reverse(t *testing.T) {
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, slice.Reverse([]int{3, 2, 1, 0, -1}))
}

var even = func(v int) bool { return v%2 == 0 }

func Test_Slice_Filter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	f := slice.Filter(s, even)
	e := []int{2, 4, 6}
	assert.Equal(t, e, f)
}

func Test_Slice_Group(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	g := slice.Group(s, even)
	e := map[bool][]int{false: {1, 3, 5}, true: {2, 4, 6}}
	assert.Equal(t, e, g)
}

func Test_Slice_ReduceSum(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	sum := slice.Reduce(s, sum.Of[int])
	e := 1 + 2 + 3 + 4 + 5 + 6
	assert.Equal(t, e, sum)
}

func Test_Slice_Flatt(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.Flatt(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_Range(t *testing.T) {
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, range_.Of(-1, 3))
	assert.Equal(t, []int{3, 2, 1, 0, -1}, range_.Of(3, -1))
	assert.Equal(t, []int{1}, range_.Of(1, 1))
}
```

## [Mutable containers](./mutable/api.go)

Supports write operations (append, delete, replace).

  * [Vector](./mutable/vector/api.go) - the simplest based on built-in slice collection.
  * [Set](./mutable/set/api.go) - collection of unique items, prevents duplicates.
  * [Map](./mutable/map_/api.go) - built-in map wrapper that supports [container functions](#container-functions).
  * [OrderedSet](./mutable/oset/api.go) - collection of unique items, prevents duplicates, provides iteration in order of addition.
  * [OrderedMap](./mutable/omap/api.go) - same as the [Map](./mutable/map_/api.go), but supports iteration in the order in which elements are added.
  * [sync.Map](./mutable/sync/map.go) - generic wrapper of built-in embedded sync.Map.

## [Immutable containers](./immutable/api.go) 

The same interfaces as in the mutable package but for read-only purposes.

## Containers creating
### Immutable
```go
package examples

import (
	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/map_"
	"github.com/m4gshm/gollections/immutable/omap"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/immutable/vector"
)

func _() {
	var (
		_ immutable.Vector[int] = vector.Of(1, 2, 3)
		_ c.Vector[int]         = vector.New([]int{1, 2, 3})
	)
	var (
		_ immutable.Set[int] = set.Of(1, 2, 3)
		_ c.Set[int]         = set.New([]int{1, 2, 3})
	)
	var (
		_ ordered.Set[int] = oset.Of(1, 2, 3)
		_ c.Set[int]       = oset.New([]int{1, 2, 3})
	)
	var (
		_ immutable.Map[int, string] = map_.Of(K.V(1, "1"), K.V(2, "2"), K.V(3, "3"))
		_ c.Map[int, string]         = map_.New(map[int]string{1: "2", 2: "2", 3: "3"})
	)
	var (
		_ ordered.Map[int, string] = omap.Of(K.V(1, "1"), K.V(2, "2"), K.V(3, "3"))
		_ c.Map[int, string]       = omap.New(map[int]string{1: "2", 2: "2", 3: "3"})//source map order is unpredictable
	)
}
```
where [vector](./immutable/vector/api.go), [set](./immutable/set/api.go), [oset](./immutable/oset/api.go), [map_](./immutable/map_/api.go), [omap](./immutable/omap/api.go) are packages from [github.com/m4gshm/gollections/immutable](./immutable/) and [K.V](./K/v.go) is the method V from the package [K](./K/)
### Mutable
```go
package examples

import (
	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/mutable/map_"
	"github.com/m4gshm/gollections/mutable/omap"
	"github.com/m4gshm/gollections/mutable/ordered"
	"github.com/m4gshm/gollections/mutable/oset"
	"github.com/m4gshm/gollections/mutable/set"
	"github.com/m4gshm/gollections/mutable/vector"
)

func _() {
	capacity := 10

	var (
		_ *mutable.Vector[int] = vector.Of(1, 2, 3)
		_ c.Vector[int]        = vector.New[int](capacity)
		_ c.Vector[int]        = vector.Empty[int]()
	)
	var (
		_ *mutable.Set[int] = set.Of(1, 2, 3)
		_ c.Set[int]        = set.New[int](capacity)
		_ c.Set[int]        = set.Empty[int]()
	)
	var (
		_ *ordered.Set[int] = oset.Of(1, 2, 3)
		_ c.Set[int]        = oset.New[int](capacity)
		_ c.Set[int]        = oset.Empty[int]()
	)
	var (
		_ *mutable.Map[int, string] = map_.Of(K.V(1, "1"), K.V(2, "2"), K.V(3, "3"))
		_ c.Map[int, string]        = map_.New[int, string](capacity)
		_ c.Map[int, string]        = map_.Empty[int, string]()
	)
	var (
		_ *ordered.Map[int, string] = omap.Of(K.V(1, "1"), K.V(2, "2"), K.V(3, "3"))
		_ c.Map[int, string]        = omap.New[int, string](capacity)
		_ c.Map[int, string]        = omap.Empty[int, string]()
	)
}
```
where [vector](./mutable/vector/api.go), [set](./mutable/set/api.go), [oset](./mutable/oset/api.go), [map_](./mutable/map_/api.go), [omap](./mutable/omap/api.go) are packages from [github.com/m4gshm/gollections/mutable](./mutable/) and [K.V](./K/v.go) is the method V from the package [K](./K/)

## Pipe functions

There are three groups of operations:
 * Immediate - retrieves the result in place ([Sort](./immutable/vector.go#L112), [Reduce](./immutable/vector.go#L107) (of containers), [Track](./immutable/vector.go#L81), [TrackEach](./immutable/ordered/map.go#L136), [For](./immutable/vector.go#L89), [ForEach](./immutable/ordered/map.go#L144))
 * Intermediate - only defines a computation ([Wrap](./it/api.go#L17), [Map](./c/op/api.go#L11), [Flatt](./c/op/api.go#L21), [Filter](./c/op/api.go#L33), [Group](./c/op/api.go#L53)).
 * Final - applies intermediates and retrieves a result ([ForEach](./it/api.go#L75), [Slice](./it/api.go#L65), [Reduce](./it/api.go#L55) (of iterators))

Intermediates should wrap one by one to make a lazy computation chain that can be applied to the latest final operation.

```go
//Example 'filter', 'map', 'reduce' for an iterative container of 'items'

var items immutable.Vector[Item]

var (
    condition c.Predicate[Item]  = func(item Item) ...
    max       op.Binary[Attribute] = func(attribute1 Attribute, attribute2 Attribute) ...
) 

maxItemAttribute := it.Reduce(it.Map(c.Filer(items, condition), Item.GetAttribute), max)
```
Functions grouped into packages by applicable type ([container](./c/api.go), [map](./c/map_/api.go), [iterator](./it/api.go), [slice](slice/api.go))

## Packages
### [Common interfaces](./c/iface.go)

Iterator, Iterable, Container, Vector, Map, Set and so on.

### [Iterable container API](./c/op/api.go)
Declarative style API over 'Iterable' interface. Based on 'Iterator API' (see below).

### [Iterator API](./it/api.go)
Declarative style API over 'Iterator' interface. 

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
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/sum"
	"github.com/m4gshm/gollections/walk/group"
)

func Test_Slice_filtering(t *testing.T) {
	f := slice.Filter([]int{1, 2, 3, 4, 5, 6}, func(i int) bool { return i%2 == 0 })
	assert.Equal(t, []int{2, 4, 6}, f)
}

func Test_Set(t *testing.T) {
	var (
		s      *immutable.Set[int] = set.Of(1, 1, 2, 4, 3, 1)
		values []int               = s.Collect()
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
	oddSum := it.Reduce(it.Filter(it.Flatt(slc.Flatt(multiDimension, conv.To[[][]int]), conv.To[[]int]), odds), sum.Of[int])
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