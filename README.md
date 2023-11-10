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

    assert.Equal(t, namesByRole[""], []string{"Tom"})
    assert.Equal(t, namesByRole["manager"], []string{"Bob", "Alice"})
    assert.Equal(t, namesByRole["admin"], []string{"Bob"})

}
```

Or you can write more compact code using the collections API, like:

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

## Slices

``` go
data, err := slice.Conv(slice.Of("1", "2", "3", "4", "_", "6"), strconv.Atoi)
even := func(i int) bool { return i%2 == 0 }

result := slice.Reduce(slice.Convert(slice.Filter(data, even), strconv.Itoa), op.Sum[string])

assert.ErrorIs(t, err, strconv.ErrSyntax)
assert.Equal(t, "24", result)
```

In the example is used only small set of slice functions as
[slice.Filter](#slicefilter), [slice.Conv](#sliceconv)
[Convert](./slice/api.go#L166), and [Reduce](#slicereduce). More you can
look in the [slice](./slice/api.go) package.

### Shortcut packages

``` go
result := sum.Of(filter.AndConvert(data, even, strconv.Itoa))
```

This is a shorter version of the previous example that used short
aliases [sum.Of](#sumof) and [filter.AndConvert](#filterandonvert).

#### Brief of usage

``` go
data := slice.Of("Bob", "Chris", "Alice") // constructor

sorted := sort.Asc(data) //sorting

reversed := reverse.Of(clone.Of(sorted)) //reversing of cloned slice

chris, ok := first.Of(reversed, func(name string) bool { return name[0] == 'C' }) //finding the first element by a predicate function

var lengthMap map[int][]string = group.Of(sorted, func(name string) int { return len(name) }, as.Is[string]) // converting to a map

assert.Equal(t, slice.Of("Alice", "Bob", "Chris"), sorted)
assert.Equal(t, "Chris", chris)
assert.True(t, true, ok)
assert.Equal(t, slice.Of("Alice", "Chris"), lengthMap[5])
```

More shortcuts you can find by exploring slices [subpackages](./slice).

### Main slice functions

#### slice.Of

``` go
"github.com/m4gshm/gollections/slice"

var s = slice.Of(1, 3, -1, 2, 0) //[]int{1, 3, -1, 2, 0}
```

#### sort.Asc, sort.Desc

``` go
// sorting in-place API
import "github.com/m4gshm/gollections/slice/sort"

var ascendingSorted = sort.Asc([]int{1, 3, -1, 2, 0})   //[]int{-1, 0, 1, 2, 3}
var descendingSorted = sort.Desc([]int{1, 3, -1, 2, 0}) //[]int{3, 2, 1, 0, -1}
```

#### sort.By, sort.ByDesc

``` go
// sorting copied slice API does not change original slice
import "github.com/m4gshm/gollections/slice/clone/sort" 

// see the User structure above
var users = []User{
    {name: "Bob", age: 26},
    {name: "Alice", age: 35},
    {name: "Tom", age: 18},
    {name: "Chris", age: 41},
}

var byName = sort.By(users, User.Name) //[{Alice 35 []} {Bob 26 []} {Chris 41 []} {Tom 18 []}]
var byAgeReverse = sort.DescBy(users, User.Age) //[{Chris 41 []} {Alice 35 []} {Bob 26 []} {Tom 18 []}]
```

#### group.Of

``` go
import (
    "github.com/m4gshm/gollections/convert/as"
    "github.com/m4gshm/gollections/expr/use"
    "github.com/m4gshm/gollections/slice/group"
)

var ageGroups map[string][]User

ageGroups = group.Of(users, func(u User) string {
    return use.If(u.age <= 20, "<=20").If(u.age <= 30, "<=30").Else(">30")
}, as.Is[User])

//map[<=20:[{Tom 18 []}] <=30:[{Bob 26 []}] >30:[{Alice 35 []} {Chris 41 []}]]
```

#### slice.ToMap, slice.ToMapResolv

``` go
   "github.com/stretchr/testify/assert"
)

func Test_Slice_ToMapResolv(t *testing.T) {

    ageGroupedSortedNames = slice.ToMapResolv(users, func(u User) string {
        return op.IfElse(u.age <= 30, "<=30", ">30")
    }, User.Name, resolv.SortedSlice)

    //map[<=30:[Bob Tom] >30:[Alice Chris]]

    assert.Equal(t, slice.Of("Bob", "Tom"), ageGroupedSortedNames["<=30"])
    assert.Equal(t, slice.Of("Alice", "Chris"), ageGroupedSortedNames[">30"])
}
```

#### sum.Of

``` go
import "github.com/m4gshm/gollections/op/sum"

var sum = sum.Of(1, 2, 3, 4, 5, 6)
//21
```

#### range\_.Of

``` go
import "github.com/m4gshm/gollections/slice/range_"

var increasing = range_.Of(-1, 3)    //[]int{-1, 0, 1, 2}
var decreasing = range_.Of('e', 'a') //[]rune{'e', 'd', 'c', 'b'}
var nothing = range_.Of(1, 1)        //nil
```

#### range\_.Closed

``` go
import "github.com/m4gshm/gollections/slice/range_"

var increasing = range_.Closed(-1, 3)    //[]int{-1, 0, 1, 2, 3}
var decreasing = range_.Closed('e', 'a') //[]rune{'e', 'd', 'c', 'b', 'a'}
var one = range_.Closed(1, 1)            //[]int{1}
```

#### first.Of

``` go
import (
    "github.com/m4gshm/gollections/expr/first"
    "github.com/m4gshm/gollections/predicate/more"
)

result, ok := first.Of(1, 3, 5, 7, 9, 11).By(more.Than(5)) //7, true
```

#### last.Of

``` go
import (
    "github.com/m4gshm/gollections/expr/last"
    "github.com/m4gshm/gollections/predicate/less"
)

result, ok := last.Of(1, 3, 5, 7, 9, 11).By(less.Than(9)) //7, true
```

#### slice.Convert

``` go
import "strconv"
import "github.com/m4gshm/gollections/slice"

var s []string = slice.Convert([]int{1, 3, 5, 7, 9, 11}, strconv.Itoa)
//[]string{"1", "3", "5", "7", "9", "11"}
```

#### slice.Conv

``` go
import "strconv"
import "github.com/m4gshm/gollections/slice"

result, err := slice.Conv(slice.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi)
//[]int{1, 3, 5}, ErrSyntax
```

#### slice.Filter

``` go
import (
    "github.com/m4gshm/gollections/predicate/exclude"
    "github.com/m4gshm/gollections/predicate/one"
    "github.com/m4gshm/gollections/slice"
)

var f1 = slice.Filter([]int{1, 3, 5, 7, 9, 11}, one.Of(1, 7).Or(one.Of(11))) //[]int{1, 7, 11}
var f2 = slice.Filter([]int{1, 3, 5, 7, 9, 11}, exclude.All(1, 7, 11))       //[]int{3, 5, 9}
```

#### convert.AndReduce

#### conv.AndReduce

#### convert.AndFilter

#### filter.AndConvert

## Maps

``` go
type Employer struct {
    name   string
    salary int
}

employers := map[string][]Employer{
    "internal": {{"Alice", 100}, {"Bob", 90}},
    "external": {{"Chris", 125}},
    "foreing":  {{"Mari", 99}},
}

noForeings := filter.Values(employers, func(employers []Employer) bool {
    return slice.Has(employers, func(e Employer) bool { return e.name != "Mari" })
})

assert.Equal(t, slice.Of("external", "internal"), sort.Asc(map_.Keys(noForeings)))

var (
    getSalary                 = func(e Employer) int { return e.salary }
    getDepartmentAndSalarySum = func(department string, e []Employer) (string, int) {
        return department, slice.ConvertAndReduce(e, getSalary, op.Sum[int])
    }
)

departmentSalary := map_.Convert(noForeings, getDepartmentAndSalarySum)

assert.Equal(t, 2, len(departmentSalary))
assert.Equal(t, 190, departmentSalary["internal"])
assert.Equal(t, 125, departmentSalary["external"])
```

More shortcuts are [here](./map_).

### Main map functions (TODO)

- creating - map\_.Of, clone.Of, map\_.DeepClone

- extract keys, values - map\_.Keys, map\_.Values

- converting, filtering, summarizing - map\_.Convert, map\_.ConvertKeys,
  map\_.ConvertValues, map\_.ToSlice, map\_.Reduce

- versions of aboves with possible result errors - map\_.Conv

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
bob, _ := slice.First(users, where.Eq(User.Name, "Bob"))

assert.Equal(t, "Bob", bob.Name())
```

### [loop](./loop/api.go), [kv/loop](./kv/loop/api.go) and breakable versions [break/loop](./break/loop/api.go), [break/kv/loop](./break/kv/loop/api.go)

Low level iteration api based on `next` function.

``` go
type (
    next[T any]      func() (element T, ok bool)
    kvNext[K, V any] func() (key K, value V, ok bool)
)
```

The function retrieves a next element from a dataset and returns
`ok==true` if successful.  
The API in most cases is similar to the [slice](./slice/api.go) API but
with delayed computation which means that the methods don’t compute a
result but only return a loop provider. The loop provider is type with a
`Next` method that returns a next processed element.

``` go
even := func(i int) bool { return i%2 == 0 }
loopStream := loop.Convert(loop.Filter(loop.Of(1, 2, 3, 4), even).Next, strconv.Itoa)

assert.Equal(t, []string{"2", "4"}, loop.Slice(loopStream.Next))
```

Breakable loops additionaly have error returned value.

``` go
type (
    next[T any]      func() (element T, ok bool, err error)
    kvNext[K, V any] func() (key K, value V, ok bool, err error)
)
```

It is used for computations where an error may occur.

``` go
iter := loop.Conv(loop.Of("1", "2", "3", "ddd4", "5"), strconv.Atoi)
result, err := loop.Slice(iter.Next)

assert.Equal(t, []int{1, 2, 3}, result)
assert.ErrorContains(t, err, "invalid syntax")
```

### Expressions: [use](./expr/use/api.go), [get](./expr/get/api.go), [first](./expr/first/api.go), [last](./expr/last/api.go)

Aimed to evaluate a value using conditions. May cause to make code
shorter by not in all cases.  
As example:

``` go
user := User{name: "Bob", surname: "Smith"}

fullName := use.If(len(user.surname) == 0, user.name).If(len(user.name) == 0, user.surname).
    ElseGet(func() string { return user.name + " " + user.surname })

assert.Equal(t, "Bob Smith", fullName)
```

instead of:

``` go
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

## Mutable collections

Supports write operations (append, delete, replace).

- [Vector](./collection/mutable/vector/api.go) - the simplest based on
  built-in slice collection.

``` go
_ *mutable.Vector[int]   = vector.Of(1, 2, 3)
_ collection.Vector[int] = &mutable.Vector[int]{}
```

- [Set](./collection/mutable/set/api.go) - collection of unique items,
  prevents duplicates.

``` go
_ *mutable.Set[int]   = set.Of(1, 2, 3)
_ collection.Set[int] = &mutable.Set[int]{}
```

- [Map](./collection/mutable/map_/api.go) - built-in map wrapper that
  supports [stream functions](#stream-functions).

``` go
_ *mutable.Map[int, string]   = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
_ collection.Map[int, string] = mutable.NewMapOf(map[int]string{1: "2", 2: "2", 3: "3"})
```

- [OrderedSet](./collection/mutable/oset/api.go) - collection of unique
  items, prevents duplicates, provides iteration in order of addition.

``` go
_ *ordered.Set[int]   = set.Of(1, 2, 3)
_ collection.Set[int] = &ordered.Set[int]{}
```

- [OrderedMap](./collection/mutable/omap/api.go) - same as the Map, but
  supports iteration in the order in which elements are added.

``` go
_ *ordered.Map[int, string]   = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
_ collection.Map[int, string] = ordered.NewMapOf(
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
  [Reduce](./collection/immutable/vector.go#L154),
  [Track](./collection/immutable/vector.go#L111),
  [TrackEach](./collection/mutable/ordered/map.go#L182),
  [For](./collection/immutable/vector.go#L122),
  [ForEach](./collection/immutable/ordered/map.go#L175))

- Intermediate - only defines a computation
  ([Convert](./collection/api.go#L17),
  [Filter](./collection/immutable/ordered/set.go#L124),
  [Flatt](./collection/api.go#L36), [Group](./collection/api.go#L69)).

- Final - applies intermediates and retrieves a result
  ([First](./collection/api.go#L75),
  [Slice](./collection/immutable/ordered/set.go#L94),
  [Reduce](./collection/immutable/ordered/set.go#L146))

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
