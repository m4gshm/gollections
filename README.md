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

You can make clear code, extensive, but without dependencies:

``` go
var namesByRole = map[string][]string{}
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

Or using the devkit you can write more compact code like this:

``` go
var namesByRole = group.ByMultipleKeys(users, func(u User) []string {
    return convert.AndConvert(u.Roles(), Role.Name, strings.ToLower)
}, User.Name)

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

All packages consists of functions placed in the package and subpackages
aimed to make short aliases of that functions. For example the function
[slice.SortByOrdered](./slice/api.go#L459) has aliases
[sort.By](./slice/sort/api.go#L12) and
[sort.Of](./slice/sort/api.go#L23).

### [slice](./slice/api.go) and [map\_](./map_/api.go)

Contains utility functions of [converting](./slice/api.go#L156),
[filtering](./slice/api.go#L379) (searching),
[reducing](./slice/api.go#L464), [cloning](./map_/api.go#L90) elements
of embedded slices and maps.

Usage examples
[here](./internal/examples/sliceexamples/slice_examples_test.go) and
[here](./internal/examples/mapexamples/map_examples_test.go).

### [mutable](./collection/mutable/api.go) and [immutable](./collection/immutable/api.go) collections

Provides implelentations of [Vector](./collection/iface.go#L25),
[Set](./collection/iface.go#L35) and [Map](./collection/iface.go#L41).

Mutables support content appending, updating and deleting (the ordered
map implementation is not supported delete operations). Immutables are
read-only datasets.

Detailed description of implementations below.

### [predicate](./predicate/api.go) and breakable [break/predicate](./predicate/api.go)

Provides predicate builder api that used for filtering collection
elements.

### [loop](./loop/api.go), [kv/loop](./kv/loop/api.go) and breakable versions [break/loop](./break/loop/api.go), [break/kv/loop](./break/kv/loop/api.go)

TODO

### Short aliases for collection constructors

TODO

## Mutable collections

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
  supports [stream functions](#stream-functions).

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
//TODO
```
