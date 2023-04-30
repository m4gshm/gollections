# Gollections

Golang generic containers and functions.

Supports Go version 1.20 or higher.

## Installation

``` console
go get -u github.com/m4gshm/gollections
```

or

``` console
go get -u github.com/m4gshm/gollections@HEAD
```

## Main packages

### Slices - [github.com/m4gshm/gollections/slice](./slice/api.go)

The package provides helper subpackages and functions for using with
slices.  
Most helpers organized as pair of a main function and short aliases in a
subpackage. For example the function
[slice.SortByOrdered](./slice/api.go#L247) has aliases
[sort.By](./slice/sort/api.go#L12) and
[sort.Of](./slice/sort/api.go#L23).

Usage examples
[here](./internal/examples/sliceexamples/slice_examples_test.go).

### Loops - [loop](./loop/api.go), [kv/loop](./kv/loop/api.go)

TODO

### Maps - [github.com/m4gshm/gollections/map\_](./map_/api.go)

The package provides helper subpackages and functions for using with
maps.  

Usage examples
[here](./internal/examples/mapexamples/map_examples_test.go).

### Mutable containers

Supports write operations (append, delete, replace).

- [Vector](./mutable/vector/api.go) - the simplest based on built-in
  slice collection.

- [Set](./mutable/set/api.go) - collection of unique items, prevents
  duplicates.

- [Map](./mutable/map_/api.go) - built-in map wrapper that supports
  [container functions](#container-functions).

- [OrderedSet](./mutable/oset/api.go) - collection of unique items,
  prevents duplicates, provides iteration in order of addition.

- [OrderedMap](./mutable/omap/api.go) - same as the
  [Map](./mutable/map_/api.go), but supports iteration in the order in
  which elements are added.

- [sync.Map](./mutable/sync/map.go) - generic wrapper of built-in
  embedded sync.Map.

### Immutable containers

The same interfaces as in the mutable package but for read-only
purposes.

### Containers creating

#### Immutable

``` go
// Package examples of collection constructors
package examples

import (
    "github.com/m4gshm/gollections/collection"
    "github.com/m4gshm/gollections/immutable"
    imap "github.com/m4gshm/gollections/immutable/map_"
    "github.com/m4gshm/gollections/immutable/omap"
    "github.com/m4gshm/gollections/immutable/ordered"
    "github.com/m4gshm/gollections/immutable/oset"
    "github.com/m4gshm/gollections/immutable/set"
    "github.com/m4gshm/gollections/immutable/vector"
    "github.com/m4gshm/gollections/k"
)

func _() {
    var (
        _ immutable.Vector[int]  = vector.Of(1, 2, 3)
        _ collection.Vector[int] = vector.New([]int{1, 2, 3})
    )
    var (
        _ immutable.Set[int]  = set.Of(1, 2, 3)
        _ collection.Set[int] = set.New([]int{1, 2, 3})
    )
    var (
        _ ordered.Set[int]    = oset.Of(1, 2, 3)
        _ collection.Set[int] = oset.New([]int{1, 2, 3})
    )
    var (
        _ immutable.Map[int, string]  = imap.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
        _ collection.Map[int, string] = imap.New(map[int]string{1: "2", 2: "2", 3: "3"})
    )
    var (
        _ ordered.Map[int, string]    = omap.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
        _ collection.Map[int, string] = omap.New(
            /*uniques*/ map[int]string{1: "2", 2: "2", 3: "3"} /*order*/, []int{3, 1, 2},
        )
    )
}
```

where [vector](./immutable/vector/api.go),
[set](./immutable/set/api.go), [oset](./immutable/oset/api.go),
[map\_](./immutable/map_/api.go), [omap](./immutable/omap/api.go) are
packages from [github.com/m4gshm/gollections/immutable](./immutable/)
and [k.V](./k/v.go) is the method V from the package [k](./k/)

#### Mutable

``` go
package examples

import (
    "github.com/m4gshm/gollections/collection"
    "github.com/m4gshm/gollections/k"
    "github.com/m4gshm/gollections/mutable"
    mmap "github.com/m4gshm/gollections/mutable/map_"
    "github.com/m4gshm/gollections/mutable/omap"
    "github.com/m4gshm/gollections/mutable/ordered"
    "github.com/m4gshm/gollections/mutable/oset"
    "github.com/m4gshm/gollections/mutable/set"
    "github.com/m4gshm/gollections/mutable/vector"
)

func _() {
    capacity := 10

    var (
        _ *mutable.Vector[int]   = vector.Of(1, 2, 3)
        _ *mutable.Vector[int]   = new(mutable.Vector[int])
        _ *mutable.Vector[int]   = vector.NewCap[int](capacity)
        _ collection.Vector[int] = vector.Empty[int]()
    )
    var (
        _ *mutable.Set[int]   = set.Of(1, 2, 3)
        _ *mutable.Set[int]   = new(mutable.Set[int])
        _ *mutable.Set[int]   = set.NewCap[int](capacity)
        _ collection.Set[int] = set.Empty[int]()
    )
    var (
        _ *ordered.Set[int]   = oset.Of(1, 2, 3)
        _ *ordered.Set[int]   = new(ordered.Set[int])
        _ *ordered.Set[int]   = oset.NewCap[int](capacity)
        _ collection.Set[int] = oset.Empty[int]()
    )
    var (
        _ *mutable.Map[int, string]   = mmap.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
        _ *mutable.Map[int, string]   = new(mutable.Map[int, string])
        _ *mutable.Map[int, string]   = mmap.New[int, string](capacity)
        _ collection.Map[int, string] = mmap.Empty[int, string]()
    )
    var (
        _ *ordered.Map[int, string]   = omap.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
        _ *ordered.Map[int, string]   = new(ordered.Map[int, string])
        _ *ordered.Map[int, string]   = omap.New[int, string](capacity)
        _ collection.Map[int, string] = omap.Empty[int, string]()
    )
}
```

where [vector](./mutable/vector/api.go), [set](./mutable/set/api.go),
[oset](./mutable/oset/api.go), [map\_](./mutable/map_/api.go),
[omap](./mutable/omap/api.go) are packages from
[github.com/m4gshm/gollections/mutable](./mutable/) and [k.V](./k/v.go)
is the method V from the package [k](./k/)

## Stream functions

There are three groups of operations:

- Immediate - retrieves the result in place
  ([Sort](./immutable/vector.go#L112),
  [Reduce](./immutable/vector.go#L107) (of containers),
  [Track](./immutable/vector.go#L81),
  [TrackEach](./immutable/ordered/map.go#L136),
  [For](./immutable/vector.go#L89),
  [ForEach](./immutable/ordered/map.go#L144))

- Intermediate - only defines a computation ([Wrap](./it/api.go#L17),
  [Map](./c/op/api.go#L11), [Flatt](./c/op/api.go#L21),
  [Filter](./c/op/api.go#L33), [Group](./c/op/api.go#L53)).

- Final - applies intermediates and retrieves a result
  ([ForEach](./it/api.go#L75), [Slice](./it/api.go#L65),
  [Reduce](./it/api.go#L55) (of iterators))

Intermediates should wrap one by one to make a lazy computation chain
that can be applied to the latest final operation.

``` go
//Example 'filter', 'map', 'reduce' for an iterative container of 'items'

var items immutable.Vector[Item]

var (
    condition = func(item Item) bool { ... }
    max       = func(attribute1 Attribute, attribute2 Attribute) Attribute { ... }
)

maxItemAttribute := it.Reduce(it.Convert(c.Filer(items, condition), Item.GetAttribute), max)
```

Functions grouped into packages by applicable type
([container](./c/api.go), [map](./c/map_/api.go),
[iterator](./it/api.go), [slice](slice/api.go))

## Additional packages

### [Common interfaces](./c/iface.go)

Iterator, Iterable, Container, Vector, Map, Set and so on.

### [Iterable container API](./c/op/api.go)

Declarative style API over 'Iterable' interface. Based on 'Iterator API'
(see below).

### [Iterator API](./it/api.go)

Declarative style API over 'Iterator' interface.

## Examples

``` go
package examples

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/m4gshm/gollections/as"
    cGroup "github.com/m4gshm/gollections/collection/group"
    "github.com/m4gshm/gollections/immutable"
    "github.com/m4gshm/gollections/immutable/oset"
    "github.com/m4gshm/gollections/immutable/set"
    "github.com/m4gshm/gollections/iter"
    "github.com/m4gshm/gollections/op"
    "github.com/m4gshm/gollections/predicate/more"
    "github.com/m4gshm/gollections/slice"
    sliceIter "github.com/m4gshm/gollections/slice/iter"
    "github.com/m4gshm/gollections/walk/group"
)

func Test_Set(t *testing.T) {
    var (
        s      immutable.Set[int] = set.Of(1, 1, 2, 4, 3, 1)
        values []int              = s.Slice()
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
    values := s.Slice()
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
    var groups = cGroup.Of(oset.Of(
        "seventh", "seventh", //duplicated
        "first", "second", "third", "fourth",
        "fifth", "sixth", "eighth",
        "ninth", "tenth", "one", "two", "three", "1",
        "second", //duplicate
    ), func(v string) int { return len(v) },
    ).FilterKey(
        more.Than(3),
    ).ConvertValue(
        func(v string) string { return v + "_" },
    ).Map()

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
    oddSum := iter.Reduce(iter.Filter(iter.Flatt(sliceIter.Flatt(multiDimension, as.Is[[][]int]), as.Is[[]int]), odds), op.Sum[int])
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
