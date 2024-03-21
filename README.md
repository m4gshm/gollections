# Gollections

Gollections is set of functions for [slices](#slices), [maps](#maps) and
additional implementations of data structures such as [ordered
map](#mutable-collections) or [set](#mutable-collections) aimed to
reduce boilerplate code.

Supports Go version 1.21.

For example, it’s need to group some
[users](./internal/examples/boilerplate/user_type.go) by their role
names converted to lowercase:

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
add := func(role string, u User) {
    namesByRole[role] = append(namesByRole[role], u.Name())
}
for _, u := range users {
    if roles := u.Roles(); len(roles) == 0 {
        add("", u)
    } else {
        for _, r := range roles {
            add(strings.ToLower(r.Name()), u)
        }
    }
}
//map[:[Tom] admin:[Bob] manager:[Bob Alice]]
```

Or you can write more compact code using the collections API, like:

``` go
import "github.com/m4gshm/gollections/slice/convert"
import "github.com/m4gshm/gollections/slice/group"

var namesByRole = group.ByMultipleKeys(users, func(u User) []string {
    return convert.AndConvert(u.Roles(), Role.Name, strings.ToLower)
}, User.Name)
// map[:[Tom] admin:[Bob] manager:[Bob Alice]]
```

## Installation

``` console
go get -u github.com/m4gshm/gollections
```

## Slices

``` go
data, err := slice.Conv(slice.Of("1", "2", "3", "4", "_", "6"), strconv.Atoi)
even := func(i int) bool { return i%2 == 0 }

result := slice.Reduce(slice.Convert(slice.Filter(data, even), strconv.Itoa), op.Sum)

assert.ErrorIs(t, err, strconv.ErrSyntax)
assert.Equal(t, "24", result)
```

In the example is used only small set of slice functions as
[slice.Filter](#slicefilter), [slice.Conv](#sliceconv)
[slice.Convert](#sliceconvert#), and [slice.Reduce](#slicereduce). More
you can look in the [slice](./slice/api.go) package.

### Shortcut packages

``` go
result := sum.Of(filter.AndConvert(data, even, strconv.Itoa))
```

This is a shorter version of the previous example that used short
aliases [sum.Of](#sumof) and
[filter.AndConvert](#operations-chain-functions). More shortcuts you can
find by exploring slices [subpackages](./slice).

**Be careful** when use several slice functions subsequently like
`slice.Filter(slice.Convert(…​))`. This can lead to unnecessary RAM
consumption. Consider [chain functions](#operations-chain-functions)
instead, [loop](#additional-api) or [collections](#collection-functions)
for delayed operations.

### Main slice functions

#### Instantiators

##### slice.Of

``` go
var s = slice.Of(1, 3, -1, 2, 0) //[]int{1, 3, -1, 2, 0}
```

##### range\_.Of

``` go
import "github.com/m4gshm/gollections/slice/range_"

var increasing = range_.Of(-1, 3)    //[]int{-1, 0, 1, 2}
var decreasing = range_.Of('e', 'a') //[]rune{'e', 'd', 'c', 'b'}
var nothing = range_.Of(1, 1)        //nil
```

##### range\_.Closed

``` go
var increasing = range_.Closed(-1, 3)    //[]int{-1, 0, 1, 2, 3}
var decreasing = range_.Closed('e', 'a') //[]rune{'e', 'd', 'c', 'b', 'a'}
var one = range_.Closed(1, 1)            //[]int{1}
```

#### Sorters

##### sort.Asc, sort.Desc

``` go
// sorting in-place API
import "github.com/m4gshm/gollections/slice/sort"

var ascendingSorted = sort.Asc([]int{1, 3, -1, 2, 0})   //[]int{-1, 0, 1, 2, 3}
var descendingSorted = sort.Desc([]int{1, 3, -1, 2, 0}) //[]int{3, 2, 1, 0, -1}
```

##### sort.By, sort.ByDesc

``` go
// sorting copied slice API does not change the original slice
import "github.com/m4gshm/gollections/slice/clone/sort"

// see the User structure above
var users = []User{
    {name: "Bob", age: 26},
    {name: "Alice", age: 35},
    {name: "Tom", age: 18},
    {name: "Chris", age: 41},
}

var byName = sort.By(users, User.Name)
//[{Alice 35 []} {Bob 26 []} {Chris 41 []} {Tom 18 []}]

var byAgeReverse = sort.DescBy(users, User.Age)
//[{Chris 41 []} {Alice 35 []} {Bob 26 []} {Tom 18 []}]
```

#### To map converters

##### group.Of

``` go
import "github.com/m4gshm/gollections/convert/as"
import "github.com/m4gshm/gollections/expr/use"
import "github.com/m4gshm/gollections/slice/group"

var ageGroups = group.Of(users, func(u User) string {
    return use.If(u.age <= 20, "<=20").If(u.age <= 30, "<=30").Else(">30")
}, as.Is)

//map[<=20:[{Tom 18 []}] <=30:[{Bob 26 []}] >30:[{Alice 35 []} {Chris 41 []}]]
```

##### group.ByMultipleKeys

``` go
import "github.com/m4gshm/gollections/slice/convert"
import "github.com/m4gshm/gollections/slice/group"

var namesByRole = group.ByMultipleKeys(users, func(u User) []string {
    return convert.AndConvert(u.Roles(), Role.Name, strings.ToLower)
}, User.Name)
// map[:[Tom] admin:[Bob] manager:[Bob Alice]]
```

##### slice.ToMap, slice.ToMapResolv

``` go
import (
    "github.com/m4gshm/gollections/map_/resolv"
    "github.com/m4gshm/gollections/op"
    "github.com/m4gshm/gollections/slice"
)

var ageGroupedSortedNames map[string][]string

ageGroupedSortedNames = slice.ToMapResolv(users, func(u User) string {
    return op.IfElse(u.age <= 30, "<=30", ">30")
}, User.Name, resolv.SortedSlice)

//map[<=30:[Bob Tom] >30:[Alice Chris]]
```

#### Reducers

##### sum.Of

``` go
import "github.com/m4gshm/gollections/op/sum"

var sum = sum.Of(1, 2, 3, 4, 5, 6) //21
```

##### slice.Reduce

``` go
var sum = slice.Reduce([]int{1, 2, 3, 4, 5, 6}, func(i1, i2 int) int { return i1 + i2 })
//21
```

##### slice.First

``` go
import "github.com/m4gshm/gollections/predicate/more"
import "github.com/m4gshm/gollections/slice"

result, ok := slice.First([]int{1, 3, 5, 7, 9, 11}, more.Than(5)) //7, true
```

##### slice.Last

``` go
import "github.com/m4gshm/gollections/predicate/less"
import "github.com/m4gshm/gollections/slice"

result, ok := slice.Last([]int{1, 3, 5, 7, 9, 11}, less.Than(9)) //7, true
```

#### Converters

##### slice.Convert

``` go
var s []string = slice.Convert([]int{1, 3, 5, 7, 9, 11}, strconv.Itoa)
//[]string{"1", "3", "5", "7", "9", "11"}
```

##### slice.Conv

``` go
result, err := slice.Conv(slice.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi)
//[]int{1, 3, 5}, ErrSyntax
```

##### slice.Filter

``` go
import "github.com/m4gshm/gollections/predicate/exclude"
import "github.com/m4gshm/gollections/predicate/one"
import "github.com/m4gshm/gollections/slice"

var f1 = slice.Filter([]int{1, 3, 5, 7, 9, 11}, one.Of(1, 7).Or(one.Of(11))) //[]int{1, 7, 11}
var f2 = slice.Filter([]int{1, 3, 5, 7, 9, 11}, exclude.All(1, 7, 11))       //[]int{3, 5, 9}
```

##### slice.Flat

``` go
import "github.com/m4gshm/gollections/convert/as"
import "github.com/m4gshm/gollections/slice"

var i []int = slice.Flat([][]int{{1, 2, 3}, {4}, {5, 6}}, as.Is)
//[]int{1, 2, 3, 4, 5, 6}
```

#### Operations chain functions

- convert.AndReduce, conv.AndReduce

- convert.AndFilter

- filter.AndConvert

These functions combine converters, filters and reducers.

## Maps

### Main map functions

#### Instantiators

##### clone.Of

``` go
import "github.com/m4gshm/gollections/map_/clone"

var bob = map[string]string{"name": "Bob"}
var tom = map[string]string{"name": "Tom"}

var employers = map[string]map[string]string{
    "devops": bob,
    "jun":    tom,
}

copy := clone.Of(employers)
delete(copy, "jun")
bob["name"] = "Superbob"

fmt.Printf("%v\n", employers) //map[devops:map[name:Superbob] jun:map[name:Tom]]
fmt.Printf("%v\n", copy)      //map[devops:map[name:Superbob]]
```

##### clone.Deep

``` go
import "github.com/m4gshm/gollections/map_/clone"

var bob = map[string]string{"name": "Bob"}
var tom = map[string]string{"name": "Tom"}

var employers = map[string]map[string]string{
    "devops": bob,
    "jun":    tom,
}

copy := clone.Deep(employers, func(employer map[string]string) map[string]string {
    return clone.Of(employer)
})
delete(copy, "jun")
bob["name"] = "Superbob"

fmt.Printf("%v\n", employers) //map[devops:map[name:Superbob] jun:map[name:Tom]]
fmt.Printf("%v\n", copy)      //map[devops:map[name:Bob]]
```

#### Keys, values exrtractors

##### map\_.Keys, map\_.Values

``` go
var employers = map[string]map[string]string{
    "devops": {"name": "Bob"},
    "jun":    {"name": "Tom"},
}

keys := map_.Keys(employers)     //[devops jun]
values := map_.Values(employers) //[map[name:Bob] map[name:Tom]]
```

#### Converters

##### map\_.ConvertKeys

``` go
var keys = map_.ConvertKeys(employers, func(title string) string {
    return string([]rune(title)[0])
})
//map[d:map[name:Bob] j:map[name:Tom]]
```

##### map\_.ConvertValues

``` go
var vals = map_.ConvertValues(employers, func(employer map[string]string) string {
    return employer["name"]
})
//map[devops:Bob jun:Tom]
```

##### map\_.Convert

``` go
var all = map_.Convert(employers, func(title string, employer map[string]string) (string, string) {
    return string([]rune(title)[0]), employer["name"]
})
//map[d:Bob j:Tom]
```

##### map\_.Conv

``` go
var all, err = map_.Conv(employers, func(title string, employer map[string]string) (string, string, error) {
    return string([]rune(title)[0]), employer["name"], nil
})
//map[d:Bob j:Tom], nil
```

##### map\_.ToSlice

``` go
var users = map_.ToSlice(employers, func(title string, employer map[string]string) User {
    return User{name: employer["name"], roles: []Role{{name: title}}}
})
//[{name:Bob age:0 roles:[{name:devops}]} {name:Tom age:0 roles:[{name:jun}]}]
```

## [loop](./loop/api.go), [kv/loop](./kv/loop/api.go) and breakable versions [break/loop](./break/loop/api.go), [break/kv/loop](./break/kv/loop/api.go)

Low-level API for iteration based on next functions:

``` go
type (
    Loop[T any]           func() (element T, ok bool)
    KVLoop[K, V any]      func() (key K, value V, ok bool)
    BreakLoop[T any]      func() (element T, ok bool, err error)
    BreakKVLoop[K, V any] func() (key K, value V, ok bool, err error)
)
```

The `Loop` function retrieves a next element from a dataset and returns
`ok==true` if successful.  
The `KVLoop` behaves similar but returns a key/value pair.  

``` go
even := func(i int) bool { return i%2 == 0 }
stringSeq := loop.Convert(loop.Filter(loop.Of(1, 2, 3, 4), even).Next, strconv.Itoa)

assert.Equal(t, []string{"2", "4"}, loop.Slice(stringSeq.Next))
```

`BreakLoop` and `BreakKVLoop` are used for sources that can issue an
error.

``` go
intSeq := loop.Conv(loop.Of("1", "2", "3", "ddd4", "5"), strconv.Atoi)
ints, err := loop.Slice(intSeq.Next)

assert.Equal(t, []int{1, 2, 3}, ints)
assert.ErrorContains(t, err, "invalid syntax")
```

The API in most cases is similar to the [slice](./slice/api.go) API but
with delayed computation which means that the methods don’t compute a
result but only return a loop provider. The loop provider is type with a
`Next` method that returns a next processed element.

## Data structures

### [mutable](./collection/mutable/api.go) and [immutable](./collection/immutable/api.go) collections

Provides implelentations of [Vector](./collection/iface.go#L25),
[Set](./collection/iface.go#L35) and [Map](./collection/iface.go#L41).

Mutables support content appending, updating and deleting (the ordered
map implementation is not supported delete operations).  
Immutables are read-only datasets.

Detailed description of implementations [below](#mutable-collections).

## Additional API

### [predicate](./predicate/api.go) and breakable [break/predicate](./predicate/api.go)

Provides predicate builder api that used for filtering collection
elements.

``` go
import "github.com/m4gshm/gollections/predicate/where"
import "github.com/m4gshm/gollections/slice"

bob, _ := slice.First(users, where.Eq(User.Name, "Bob"))
```

It is used for computations where an error may occur.

``` go
   assert.Equal(t, []int{1, 2, 3}, ints)
    assert.ErrorContains(t, err, "invalid syntax")

}
```

### Expressions: [use.If](./expr/use/api.go), [get.If](./expr/get/api.go), [first.Of](#firstof), [last.Of](#lastof)

Aimed to evaluate a value using conditions. May cause to make code
shorter by not in all cases.  
As example:

``` go
import "github.com/m4gshm/gollections/expr/use"

user := User{name: "Bob", surname: "Smith"}

fullName := use.If(len(user.surname) == 0, user.name).If(len(user.name) == 0, user.surname).
    ElseGet(func() string { return user.name + " " + user.surname })

assert.Equal(t, "Bob Smith", fullName)
```

instead of:

``` go
user := User{name: "Bob", surname: "Smith"}

fullName := ""
if len(user.surname) == 0 {
    fullName = user.name
} else if len(user.name) == 0 {
    fullName = user.surname
} else {
    fullName = user.name + " " + user.surname
}

assert.Equal(t, "Bob Smith", fullName)
```

#### first.Of

``` go
import "github.com/m4gshm/gollections/expr/first"
import "github.com/m4gshm/gollections/predicate/more"

result, ok := first.Of(1, 3, 5, 7, 9, 11).By(more.Than(5)) //7, true
```

#### last.Of

``` go
import "github.com/m4gshm/gollections/expr/last"
import "github.com/m4gshm/gollections/predicate/less"

result, ok := last.Of(1, 3, 5, 7, 9, 11).By(less.Than(9)) //7, true
```

## Mutable collections

Supports write operations (append, delete, replace).

- [Vector](./collection/mutable/vector/api.go) - the simplest based on
  built-in slice collection.

``` go
_ *mutable.Vector[int] = vector.Of(1, 2, 3)
_ *mutable.Vector[int] = &mutable.Vector[int]{}
```

- [Set](./collection/mutable/set/api.go) - collection of unique items,
  prevents duplicates.

``` go
_ *mutable.Set[int] = set.Of(1, 2, 3)
_ *mutable.Set[int] = &mutable.Set[int]{}
```

- [Map](./collection/mutable/map_/api.go) - built-in map wrapper that
  supports [stream functions](#stream-functions).

``` go
_ *mutable.Map[int, string] = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
_ *mutable.Map[int, string] = mutable.NewMapOf(map[int]string{1: "2", 2: "2", 3: "3"})
```

- [OrderedSet](./collection/mutable/oset/api.go) - collection of unique
  items, prevents duplicates, provides iteration in order of addition.

``` go
_ *ordered.Set[int] = set.Of(1, 2, 3)
_ *ordered.Set[int] = &ordered.Set[int]{}
```

- [OrderedMap](./collection/mutable/omap/api.go) - same as the Map, but
  supports iteration in the order in which elements are added.

``` go
_ *ordered.Map[int, string] = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
_ *ordered.Map[int, string] = ordered.NewMapOf(
    /*order  */ []int{3, 1, 2},
    /*uniques*/ map[int]string{1: "2", 2: "2", 3: "3"},
)
```

### Immutable containers

The same underlying interfaces but for read-only use cases.

## Collection functions

There are three groups of operations:

- Immediate - retrieves the result in place
  ([Sort](./collection/mutable/vector.go#L322),
  [Reduce](./collection/immutable/vector.go#L152),
  [Track](./collection/immutable/vector.go#L109),
  [TrackEach](./collection/mutable/ordered/map.go#L155),
  [For](./collection/immutable/vector.go#L120),
  [ForEach](./collection/immutable/ordered/map.go#L151))

- Intermediate - only defines a computation
  ([Convert](./collection/api.go#22),
  [Filter](./collection/immutable/ordered/set.go#L108),
  [Flat](./collection/api.go#L41), [Group](./collection/api.go#L182)).

- Final - applies intermediates and retrieves a result
  ([First](./collection/api.go#L188),
  [Slice](./collection/immutable/ordered/set.go#L69),
  [Reduce](./collection/stream/iter.go#L76)).

Intermediates should wrap one by one to make a lazy computation chain
that can be applied to the latest final operation.

``` go
var groupedByLength = group.Of(set.Of(
    "seventh", "seventh", //duplicate
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
}, groupedByLength)
```
