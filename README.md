# Gollections

This is a development kit aimed at reducing boilerplate code when using
[slices](./slice/api.go), [maps](./map_/api.go) and extending
functionality by new collection implementations such as [ordered
map](./collection/collection/mutable/omap/api.go) or
[set](./collection/collection/mutable/oset/api.go).

Supports Go version 1.20.

For example, you want to group some users by their role names converted
to lowercase:

``` go
var users = []User{
    {name: "Bob", age: 26, roles: []Role{{"Admin"}, {"manager"}}},
    {name: "Alice", age: 35, roles: []Role{{"Manager"}}},
    {name: "Tom", age: 18},
}
```

You can write code like this:

``` go
   var namesByRole = group.ByMultipleKeys(users, func(u User) []string {
        return convert.AndConvert(u.Roles(), Role.Name, strings.ToLower)
    }, User.Name)

    assert.Equal(t, namesByRole[""], []string{"Tom"})
    assert.Equal(t, namesByRole["manager"], []string{"Bob", "Alice"})
    assert.Equal(t, namesByRole["admin"], []string{"Bob"})
```

Or you can make clearer code, more extensive, but without dependencies:

``` go
   namesByRole := map[string][]string{}
    for _, u := range users {
        roles := u.Roles()
        if len(roles) == 0 {
            lr := ""
            names := namesByRole[lr]
            names = append(names, u.Name())
            namesByRole[lr] = names
        } else {
            for _, r := range roles {
                lr := strings.ToLower(r.Name())
                names := namesByRole[lr]
                names = append(names, u.Name())
                namesByRole[lr] = names
            }
        }
    }

    assert.Equal(t, namesByRole[""], []string{"Tom"})
    assert.Equal(t, namesByRole["manager"], []string{"Bob", "Alice"})
    assert.Equal(t, namesByRole["admin"], []string{"Bob"})
```

## Installation

``` console
go get -u github.com/m4gshm/gollections
```

or

``` console
go get -u github.com/m4gshm/gollections@HEAD
```

## Main packages

### [slice](./slice/api.go)

The package provides helper subpackages and functions for using with
slices.  
Most helpers organized as pair of a main function and short aliases in a
subpackage. For example the function
[slice.SortByOrdered](./slice/api.go#L247) has aliases
[sort.By](./slice/sort/api.go#L12) and
[sort.Of](./slice/sort/api.go#L23).

Usage examples
[here](./internal/examples/sliceexamples/slice_examples_test.go).

### [map\_](./map_/api.go)

The package provides helper subpackages and functions for using with
maps.  

Usage examples
[here](./internal/examples/mapexamples/map_examples_test.go).

### [loop](./loop/api.go), [kv/loop](./kv/loop/api.go) and breakable versions [break/loop](./break/loop/api.go), [break/kv/loop](./break/kv/loop/api.go)

TODO

### [mutable](./collection/mutable/api.go) and [immutable](./collection/immutable/api.go) collections

TODO

TODO

### [predicate](./predicate/api.go) and breakable [break/predicate](./predicate/api.go)

TODO

### Short aliases for collection constructors

TODO

### Mutable collections

Supports write operations (append, delete, replace).

- [Vector](./collection/mutable/vector/api.go) - the simplest based on
  built-in slice collection.

<!-- -->

        _ immutable.Vector[int]  = vector.Of(1, 2, 3)
        _ collection.Vector[int] = immutable.NewVector([]int{1, 2, 3})

- [Set](./collection/mutable/set/api.go) - collection of unique items,
  prevents duplicates.

<!-- -->

        _ immutable.Set[int]  = set.Of(1, 2, 3)
        _ collection.Set[int] = immutable.NewSet([]int{1, 2, 3})

- [Map](./collection/mutable/map_/api.go) - built-in map wrapper that
  supports [container functions](#container-functions).

<!-- -->

        _ immutable.Map[int, string]  = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
        _ collection.Map[int, string] = immutable.NewMap(map[int]string{1: "2", 2: "2", 3: "3"})

- [OrderedSet](./collection/mutable/oset/api.go) - collection of unique
  items, prevents duplicates, provides iteration in order of addition.

<!-- -->

        _ ordered.Set[int]    = oset.Of(1, 2, 3)
        _ collection.Set[int] = ordered.NewSet([]int{1, 2, 3})

- [OrderedMap](./collection/mutable/omap/api.go) - same as the
  [Map](./collection/mutable/map_/api.go), but supports iteration in the
  order in which elements are added.

<!-- -->

        _ *ordered.Map[int, string]    = omap.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
        _ collection.Map[int, string] = ordered.NewMap(

### Immutable containers

The same interfaces as in the mutable package but for read-only
purposes.

### Containers creating

#### Mutable

``` go
package examples

import (
    "github.com/m4gshm/gollections/collection"
    "github.com/m4gshm/gollections/collection/mutable"
    mmap "github.com/m4gshm/gollections/collection/mutable/map_"
    "github.com/m4gshm/gollections/collection/mutable/omap"
    "github.com/m4gshm/gollections/collection/mutable/ordered"
    "github.com/m4gshm/gollections/collection/mutable/oset"
    "github.com/m4gshm/gollections/collection/mutable/set"
    "github.com/m4gshm/gollections/collection/mutable/vector"
    "github.com/m4gshm/gollections/k"
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

where [vector](./collection/mutable/vector/api.go),
[set](./collection/mutable/set/api.go),
[oset](./collection/mutable/oset/api.go),
[map\_](./collection/mutable/map_/api.go),
[omap](./collection/mutable/omap/api.go) are packages from
[github.com/m4gshm/gollections/collection/mutable](./collection/mutable/)
and [k.V](./k/v.go) is the method V from the package [k](./k/)

## Stream functions

There are three groups of operations:

- Immediate - retrieves the result in place
  ([Sort](./collection/immutable/vector.go#L112),
  [Reduce](./collection/immutable/vector.go#L107) (of containers),
  [Track](./collection/immutable/vector.go#L81),
  [TrackEach](./collection/immutable/ordered/map.go#L136),
  [For](./collection/immutable/vector.go#L89),
  [ForEach](./collection/immutable/ordered/map.go#L144))

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

    cGroup "github.com/m4gshm/gollections/collection/group"
    "github.com/m4gshm/gollections/collection/immutable"
    "github.com/m4gshm/gollections/collection/immutable/oset"
    "github.com/m4gshm/gollections/collection/immutable/set"
    "github.com/m4gshm/gollections/convert/as"
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
