# Gollections

Golang generic containers and functions.

Supports Go version 1.18 or higher.

## Installation

```bash
go get -u github.com/m4gshm/gollections
```

or 

```bash
go get -u github.com/m4gshm/gollections@HEAD
```

## [Slice API](./slice/api.go)

The package provides various functions useful with slices.
Additionally contains subpackages as shortcuts of some functions from the base package.

Just looks some examples:

```go
package examples

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/m4gshm/gollections/first"
	"github.com/m4gshm/gollections/last"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/clone/sort"
	"github.com/m4gshm/gollections/slice/group"
	"github.com/m4gshm/gollections/slice/range_"
	"github.com/m4gshm/gollections/slice/reverse"
	"github.com/m4gshm/gollections/sum"
)

func Test_SortInt(t *testing.T) {
	c := sort.Of([]int{1, 3, -1, 2, 0})
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, c)
}

func Test_SortStructs(t *testing.T) {
	type User struct {
		name string
		age  int
	}
	var (
		users  = []User{{"Bob", 26}, {"Alice", 35}, {"Tom", 18}}
		byName = sort.By(users, func(u User) string { return u.name })
		byAge  = sort.By(users, func(u User) int { return u.age })
	)
	assert.Equal(t, []User{{"Alice", 35}, {"Bob", 26}, {"Tom", 18}}, byName)
	assert.Equal(t, []User{{"Tom", 18}, {"Bob", 26}, {"Alice", 35}}, byAge)
}

func Test_SortStructsByLess(t *testing.T) {
	type User struct {
		name string
		age  int
	}
	var (
		users        = []User{{"Bob", 26}, {"Alice", 35}, {"Tom", 18}}
		byName       = sort.ByLess(users, func(u1, u2 User) bool { return u1.name < u2.name })
		byAgeReverse = sort.ByLess(users, func(u1, u2 User) bool { return u1.age > u2.age })
	)
	assert.Equal(t, []User{{"Alice", 35}, {"Bob", 26}, {"Tom", 18}}, byName)
	assert.Equal(t, []User{{"Alice", 35}, {"Bob", 26}, {"Tom", 18}}, byAgeReverse)
}

func Test_Reverse(t *testing.T) {
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, slice.Reverse([]int{3, 2, 1, 0, -1}))
	assert.Equal(t, []int{-1, 0, 1, 2, 3}, reverse.Of([]int{3, 2, 1, 0, -1}))
}

func Test_Clone(t *testing.T) {
	type entity struct{ val string }
	var (
		first  = entity{"first"}
		second = entity{"second"}
		third  = entity{"third"}

		entities = []*entity{&first, &second, &third}
		copy     = clone.Of(entities)
	)

	assert.Equal(t, entities, copy)
	assert.NotSame(t, entities, copy)

	for i := range entities {
		assert.Same(t, entities[i], copy[i])
	}
}

func Test_DeepClone(t *testing.T) {
	type entity struct{ val string }
	var (
		first  = entity{"first"}
		second = entity{"second"}
		third  = entity{"third"}

		entities = []*entity{&first, &second, &third}
		copy     = clone.Deep(entities, func(e *entity) *entity { return ptr.Of(*e) })
	)

	assert.Equal(t, entities, copy)
	assert.NotSame(t, entities, copy)

	for i := range entities {
		assert.Equal(t, entities[i], copy[i])
		assert.NotSame(t, entities[i], copy[i])
	}
}

func Test_Convert(t *testing.T) {
	s := slice.Of(1, 3, 5, 7, 9, 11)
	r := slice.Convert(s, strconv.Itoa)
	assert.Equal(t, slice.Of("1", "3", "5", "7", "9", "11"), r)
}

var even = func(v int) bool { return v%2 == 0 }

func Test_ConvertFiltered(t *testing.T) {
	s := []int{1, 3, 4, 5, 7, 8, 9, 11}
	r := slice.ConvertFit(s, even, strconv.Itoa)
	assert.Equal(t, []string{"4", "8"}, r)
}

func Test_ConvertFilteredWithIndexInPlace(t *testing.T) {
	s := slice.Of(1, 3, 4, 5, 7, 8, 9, 11)
	r := slice.ConvertCheckIndexed(s, func(index int, elem int) (string, bool) { return strconv.Itoa(index + elem), even(elem) })
	assert.Equal(t, []string{"6", "13"}, r)
}

func Test_Slice_Filter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	f := slice.Filter(s, even)
	e := []int{2, 4, 6}
	assert.Equal(t, e, f)
}

func Test_Flatt(t *testing.T) {
	md := [][]int{{1, 2, 3}, {4}, {5, 6}}
	f := slice.Flatt(md, func(i []int) []int { return i })
	e := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, e, f)
}

func Test_Slice_Group(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	g := group.Of(s, even)
	e := map[bool][]int{false: {1, 3, 5}, true: {2, 4, 6}}
	assert.Equal(t, e, g)
}

func Test_Slice_ReduceSum(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	sum := slice.Reduce(s, op.Sum[int])
	e := 1 + 2 + 3 + 4 + 5 + 6
	assert.Equal(t, e, sum)
}

func Test_Slice_Sum(t *testing.T) {
	sum := sum.Of(1, 2, 3, 4, 5, 6)
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

func Test_First(t *testing.T) {
	r, ok := first.Of(1, 3, 5, 7, 9, 11).By(func(i int) bool { return i > 5 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)
}

func Test_Last(t *testing.T) {
	r, ok := last.Of(1, 3, 5, 7, 9, 11).By(func(i int) bool { return i < 9 })
	assert.True(t, ok)
	assert.Equal(t, 7, r)
}

func Test_BehaveAsStrings(t *testing.T) {
	type TypeBasedOnString string
	type ArrayTypeBasedOnString []TypeBasedOnString

	vals := ArrayTypeBasedOnString{"1", "2", "3"}
	strs := slice.BehaveAsStrings(vals)

	assert.Equal(t, []string{"1", "2", "3"}, strs)
}

type rows[T any] struct {
	row    []T
	cursor int
}

func (r *rows[T]) hasNext() bool    { return r.cursor < len(r.row) }
func (r *rows[T]) next() (T, error) { e := r.row[r.cursor]; r.cursor++; return e, nil }

func Test_OfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := slice.OfLoop(stream, (*rows[int]).hasNext, (*rows[int]).next)

	assert.Equal(t, slice.Of(1, 2, 3), result)
}

func Test_Generate(t *testing.T) {
	counter := 0
	result, _ := slice.Generate(func() (int, bool, error) { counter++; return counter, counter < 4, nil })

	assert.Equal(t, slice.Of(1, 2, 3), result)
}
```

## [Map API](./map_/api.go)

```go
package examples

import (
	"testing"

	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/clone"
	"github.com/m4gshm/gollections/map_/group"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone/sort"
	"github.com/stretchr/testify/assert"
)

type entity struct{ val string }

var (
	first  = entity{"first"}
	second = entity{"second"}
	third  = entity{"third"}

	entities = map[int]*entity{1: &first, 2: &second, 3: &third}
)

func Test_Clone(t *testing.T) {
	copy := clone.Of(entities)

	assert.Equal(t, entities, copy)
	assert.NotSame(t, entities, copy)

	for k := range entities {
		assert.Same(t, entities[k], copy[k])
	}
}

func Test_DeepClone(t *testing.T) {
	copy := clone.Deep(entities, func(e *entity) *entity { return ptr.Of(*e) })

	assert.Equal(t, entities, copy)
	assert.NotSame(t, entities, copy)

	for i := range entities {
		assert.Equal(t, entities[i], copy[i])
		assert.NotSame(t, entities[i], copy[i])
	}
}

func Test_Keys(t *testing.T) {
	keys := map_.Keys(entities)
	assert.Equal(t, slice.Of(1, 2, 3), sort.Of(keys))
}

func Test_Values(t *testing.T) {
	values := map_.Values(entities)
	assert.Equal(t, slice.Of(&first, &second, &third), sort.By(values, func(e *entity) string { return e.val }))
}

func Test_ConvertValues(t *testing.T) {
	var strValues map[int]string = map_.ConvertValues(entities, func(e *entity) string { return e.val })

	assert.Equal(t, "first", strValues[1])
	assert.Equal(t, "second", strValues[2])
	assert.Equal(t, "third", strValues[3])
}

func Test_ValuesConverted(t *testing.T) {
	var values []string = map_.ValuesConverted(entities, func(e *entity) string { return e.val })
	assert.Equal(t, slice.Of("1_first", "2_second", "3_third"), sort.Of(values))
}

type rows[T any] struct {
	in     []T
	cursor int
}

func (r *rows[T]) hasNext() bool    { return r.cursor < len(r.in) }
func (r *rows[T]) next() (T, error) { e := r.in[r.cursor]; r.cursor++; return e, nil }

func Test_OfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := map_.OfLoop(stream, (*rows[int]).hasNext, func(r *rows[int]) (bool, int, error) {
		n, err := r.next()
		return n%2 == 0, n, err
	})

	assert.Equal(t, 2, result[true])
	assert.Equal(t, 1, result[false])
}

func Test_Generate(t *testing.T) {
	counter := 0
	result, _ := map_.Generate(func() (bool, int, bool, error) { counter++; return counter%2 == 0, counter, counter < 4, nil })

	assert.Equal(t, 2, result[true])
	assert.Equal(t, 1, result[false])
}

func Test_GroupOfLoop(t *testing.T) {
	stream := &rows[int]{slice.Of(1, 2, 3), 0}
	result, _ := group.OfLoop(stream, (*rows[int]).hasNext, func(r *rows[int]) (bool, int, error) {
		n, err := r.next()
		return n%2 == 0, n, err
	})

	assert.Equal(t, slice.Of(2), result[true])
	assert.Equal(t, slice.Of(1, 3), result[false])
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