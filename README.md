# Gollections

Gollections is set of functions for [slices](#slices), [maps](#maps),
[iter.Seq, iter.Seq2](#seq-seq2-seqe) and additional implementations of
data structures such as [ordered map](#mutable-collections) or
[set](#mutable-collections) aimed to reduce boilerplate code.

Supports Go version 1.24.

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
import (
    "github.com/m4gshm/gollections/slice/convert"
    "github.com/m4gshm/gollections/slice/group"
)

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
//[1 2 3 4], invalid syntax

even := func(i int) bool { return i%2 == 0 }
result := slice.Reduce(slice.Convert(slice.Filter(data, even), strconv.Itoa), op.Sum) //24
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
consumption. Consider [seq](#seq) instead of slice API.

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

#### Collectors

##### group.Of

``` go
import (
    "github.com/m4gshm/gollections/convert/as"
    "github.com/m4gshm/gollections/expr/use"
    "github.com/m4gshm/gollections/slice/group"
)

var ageGroups map[string][]User = group.Of(users, func(u User) string {
    return use.If(u.age <= 20, "<=20").If(u.age <= 30, "<=30").Else(">30")
}, as.Is)

//map[<=20:[{Tom 18 []}] <=30:[{Bob 26 []}] >30:[{Alice 35 []} {Chris 41 []}]]
```

##### group.Order

``` go
import (
    "github.com/m4gshm/gollections/convert/as"
    "github.com/m4gshm/gollections/expr/use"
    "github.com/m4gshm/gollections/slice/group"
)

var order, ageGroups = group.Order(users, func(u User) string {
    return use.If(u.age <= 20, "<=20").If(u.age <= 30, "<=30").Else(">30")
}, as.Is)

//order     [<=30 >30 <=20]
//ageGroups map[<=20:[{Tom 18 []}] <=30:[{Bob 26 []}] >30:[{Alice 35 []} {Chris 41 []}]]
```

##### group.ByMultipleKeys

``` go
import (
    "github.com/m4gshm/gollections/slice/convert"
    "github.com/m4gshm/gollections/slice/group"
)

var namesByRole = group.ByMultipleKeys(users, func(u User) []string {
    return convert.AndConvert(u.Roles(), Role.Name, strings.ToLower)
}, User.Name)

// map[:[Tom] admin:[Bob] manager:[Bob Alice]]
```

##### slice.Map, slice.AppendMap

``` go
import "github.com/m4gshm/gollections/slice"

var agePerGroup map[string]int = slice.Map(users, User.Name, User.Age)

//"map[Alice:35 Bob:26 Chris:41 Tom:18]"
```

##### slice.MapOrder, slice.AppendMapOrder

``` go
import "github.com/m4gshm/gollections/slice"

var names, agePerName = slice.MapOrder(users, User.Name, User.Age)

//"[Bob Alice Tom Chris]"
//"map[Alice:35 Bob:26 Chris:41 Tom:18]"
```

##### slice.MapResolv, slice.AppendMapResolv

``` go
import (
    "github.com/m4gshm/gollections/map_/resolv"
    "github.com/m4gshm/gollections/op"
    "github.com/m4gshm/gollections/slice"
)

var ageGroupedSortedNames map[string][]string

ageGroupedSortedNames = slice.MapResolv(users, func(u User) string {
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

##### slice.Accum

``` go
import (
    "github.com/m4gshm/gollections/op"
    "github.com/m4gshm/gollections/slice"
)

var sum = slice.Accum(100, slice.Of(1, 2, 3, 4, 5, 6), op.Sum)
//121
```

##### slice.First

``` go
import (
    "github.com/m4gshm/gollections/predicate/more"
    "github.com/m4gshm/gollections/slice"
)

result, ok := slice.First([]int{1, 3, 5, 7, 9, 11}, more.Than(5)) //7, true
```

##### slice.Head

``` go
import (
    "github.com/m4gshm/gollections/slice"
)

result, ok := slice.Head([]int{1, 3, 5, 7, 9, 11}) //1, true
```

##### slice.Top

``` go
import (
    "github.com/m4gshm/gollections/slice"
)

result := slice.Top(3, []int{1, 3, 5, 7, 9, 11}) //[]int{1, 3, 5}
```

##### slice.Last

``` go
import (
    "github.com/m4gshm/gollections/predicate/less"
    "github.com/m4gshm/gollections/slice"
)

result, ok := slice.Last([]int{1, 3, 5, 7, 9, 11}, less.Than(9)) //7, true
```

##### slice.Tail

``` go
import (
    "github.com/m4gshm/gollections/slice"
)

result, ok := slice.Tail([]int{1, 3, 5, 7, 9, 11}) //11, true
```

#### Element converters

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

#### Slice converters

##### slice.Filter

``` go
import (
    "github.com/m4gshm/gollections/predicate/exclude"
    "github.com/m4gshm/gollections/predicate/one"
    "github.com/m4gshm/gollections/slice"
)

var f1 = slice.Filter([]int{1, 3, 5, 7, 9, 11}, one.Of(1, 7).Or(one.Of(11))) //[]int{1, 7, 11}
var f2 = slice.Filter([]int{1, 3, 5, 7, 9, 11}, exclude.All(1, 7, 11))       //[]int{3, 5, 9}
```

##### slice.Flat

``` go
import (
    "github.com/m4gshm/gollections/convert/as"
    "github.com/m4gshm/gollections/slice"
)

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

#### Collectors

##### map\_.Slice

``` go
var users = map_.Slice(employers, func(title string, employer map[string]string) User {
    return User{name: employer["name"], roles: []Role{{name: title}}}
})
//[{name:Bob age:0 roles:[{name:devops}]} {name:Tom age:0 roles:[{name:jun}]}]
```

##### map\_.Keys, map\_.Values, map\_.KeysConvert, map\_.ValuesConvert

``` go
var employers = map[string]map[string]string{
    "devops": {"name": "Bob"},
    "jun":    {"name": "Tom"},
}

keys := map_.Keys(employers)     //[devops jun]
values := map_.Values(employers) //[map[name:Bob] map[name:Tom]]
```

#### Element converters

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

## [seq](./seq/api.go), [seq2](./seq2/api.go), [seqe](./seqe/api.go)

API extends rangefunc iterator types `iter.Seq[V]`, `iter.Seq2[K,V]`
with utility functions kit.

``` go
even := func(i int) bool { return i%2 == 0 }
sequence := seq.Convert(seq.Filter(seq.Of(1, 2, 3, 4), even), strconv.Itoa)
var result []string = seq.Slice(sequence) //[2 4]
```

or

``` go
intSeq := seq.Conv(seq.Of("1", "2", "3", "ddd4", "5"), strconv.Atoi)
ints, err := seqe.Slice(intSeq) //[1 2 3], invalid syntax
```

### Sequence API

To use any collection or loop as a rangefunc sequecne just call
[All](#iterating-over-collections) method of that one.

In many cases the API likes the
[loop](#loop-kvloop-and-breakable-versions-breakloop-breakkvloop) API.

#### Instantiators

##### seq.Of, seq2.Of, seq2.OfMap

``` go
import(
    "github.com/m4gshm/gollections/seq"
    "github.com/m4gshm/gollections/seq2"
)

var (
    ints  iter.Seq[int]          = seq.Of(1, 2, 3)
    pairs iter.Seq2[string, int] = seq2.OfMap(map[string]int{"first": 1, "second": 2, "third": 3})
)
```

#### Collectors

##### seq.Slice

``` go
filter := func(u User) bool { return u.age <= 30 }
names := seq.Slice(seq.Convert(seq.Filter(seq.Of(users...), filter), User.Name))
//[Bob Tom]
```

##### seq2.Group, seq2.Map

``` go
import (
    "iter"

    "github.com/m4gshm/gollections/expr/use"
    "github.com/m4gshm/gollections/seq"
    "github.com/m4gshm/gollections/seq2"
    "github.com/m4gshm/gollections/slice"
    "github.com/m4gshm/gollections/slice/sort"
)

var users iter.Seq[User] = seq.Of(users...)
var groups iter.Seq2[string, User] = seq.ToSeq2(users, func(u User) (string, User) {
    return use.If(u.age <= 20, "<=20").If(u.age <= 30, "<=30").Else(">30"), u
})
var ageGroups map[string][]User = seq2.Group(groups)

//map[<=20:[{Tom 18 []}] <=30:[{Bob 26 []}] >30:[{Alice 35 []} {Chris 41 []}]]
```

#### Reducers

##### seq.Reduce

``` go
var sum = seq.Reduce(seq.Of(1, 2, 3, 4, 5, 6), func(i1, i2 int) int { return i1 + i2 })
//21
```

##### seq.ReduceOK

``` go
adder := func(i1, i2 int) int { return i1 + i2 }

sum, ok := seq.ReduceOK(seq.Of(1, 2, 3, 4, 5, 6), adder)
//21, true

emptyLoop := seq.Of[int]()
sum, ok = seq.ReduceOK(emptyLoop, adder)
//0, false
```

##### seq.First

``` go
import (
    "github.com/m4gshm/gollections/predicate/more"
    "github.com/m4gshm/gollections/seq"
)

result, ok := seq.First(seq.Of(1, 3, 5, 7, 9, 11), more.Than(5)) //7, true
```

#### Element converters

##### seq.Convert

``` go
var s []string = seq.Slice(seq.Convert(seq.Of(1, 3, 5, 7, 9, 11), strconv.Itoa))
//[]string{"1", "3", "5", "7", "9", "11"}
```

##### seq.Conv

``` go
result, err := seqe.Slice(seq.Conv(seq.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi))
//[]int{1, 3, 5}, ErrSyntax
```

#### Sequence converters

##### seq.Filter

``` go
import (
    "github.com/m4gshm/gollections/predicate/exclude"
    "github.com/m4gshm/gollections/predicate/one"
    "github.com/m4gshm/gollections/seq"
)

var f1 = seq.Slice(seq.Filter(seq.Of(1, 3, 5, 7, 9, 11), one.Of(1, 7).Or(one.Of(11))))
//[]int{1, 7, 11}

var f2 = seq.Slice(seq.Filter(seq.Of(1, 3, 5, 7, 9, 11), exclude.All(1, 7, 11)))
//[]int{3, 5, 9}
```

##### seq.Flat

``` go
import (
    "github.com/m4gshm/gollections/convert/as"
    "github.com/m4gshm/gollections/seq"
)

var i []int = seq.Slice(seq.Flat(seq.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...), as.Is))
//[]int{1, 2, 3, 4, 5, 6}
```

## [loop](./loop/api.go), [kv/loop](./kv/loop/api.go) and breakable versions [break/loop](./break/loop/api.go), [break/kv/loop](./break/kv/loop/api.go)

Deprecated: will be replaced by [seq](#seq-seq2-seqe) API. Legacy
iterators API based on the following functions:

``` go
type (
    Loop[T any]           func() (element T, ok bool)
    KVLoop[K, V any]      func() (key K, value V, ok bool)
    BreakLoop[T any]      func() (element T, ok bool, err error)
    BreakKVLoop[K, V any] func() (key K, value V, ok bool, err error)
)
```

The `Loop` function returns a next element from a dataset and returns
`ok==true` on success. `ok==false` means there are no more elements in
the dataset.  
The `KVLoop` behaves similar but returns key/value pairs.  

``` go
even := func(i int) bool { return i%2 == 0 }
seq := loop.Convert(loop.Filter(loop.Of(1, 2, 3, 4), even), strconv.Itoa)
var result []string = seq.Slice() //[2 4]
```

`BreakLoop` and `BreakKVLoop` are used for sources that can issue an
error.

``` go
intSeq := loop.Conv(loop.Of("1", "2", "3", "ddd4", "5"), strconv.Atoi)
ints, err := loop.Slice(intSeq) //[1 2 3], invalid syntax
```

The API in most cases is similar to the [slice](./slice/api.go) API but
with delayed computation which means that the methods don’t compute a
result but only return a loop provider. The loop provider is type with a
`Next` method that returns a next processed element.

### Main loop functions

#### Instantiators

##### loop.Of, loop.S

``` go
import "github.com/m4gshm/gollections/loop"

var (
    ints    = loop.Of(1, 2, 3)
    strings = loop.S([]string{"a", "b", "c"})
)
```

##### range\_.Of

``` go
import "github.com/m4gshm/gollections/loop/range_"

var increasing = range_.Of(-1, 3).Slice()    //[]int{-1, 0, 1, 2}
var decreasing = range_.Of('e', 'a').Slice() //[]rune{'e', 'd', 'c', 'b'}
var nothing = range_.Of(1, 1).Slice()        //nil
```

##### range\_.Closed

``` go
var increasing = range_.Closed(-1, 3).Slice()    //[]int{-1, 0, 1, 2, 3}
var decreasing = range_.Closed('e', 'a').Slice() //[]rune{'e', 'd', 'c', 'b', 'a'}
var one = range_.Closed(1, 1).Slice()            //[]int{1}
```

#### Collectors

##### loop.Slice

``` go
filter := func(u User) bool { return u.age <= 30 }
names := loop.Slice(loop.Convert(loop.Filter(loop.Of(users...), filter), User.Name))
//[Bob Tom]
```

##### group.Of

``` go
import (
    "github.com/m4gshm/gollections/convert/as"
    "github.com/m4gshm/gollections/expr/use"
    "github.com/m4gshm/gollections/loop"
    "github.com/m4gshm/gollections/loop/group"
)

var ageGroups map[string][]User = group.Of(loop.Of(users...), func(u User) string {
    return use.If(u.age <= 20, "<=20").If(u.age <= 30, "<=30").Else(">30")
}, as.Is)

//map[<=20:[{Tom 18 []}] <=30:[{Bob 26 []}] >30:[{Alice 35 []} {Chris 41 []}]]
```

##### loop.Map, loop.MapResolv

``` go
import (
    "github.com/m4gshm/gollections/map_/resolv"
    "github.com/m4gshm/gollections/op"
    "github.com/m4gshm/gollections/loop"
)

var ageGroupedSortedNames map[string][]string

ageGroupedSortedNames = loop.MapResolv(loop.Of(users...), func(u User) string {
    return op.IfElse(u.age <= 30, "<=30", ">30")
}, User.Name, resolv.SortedSlice)

//map[<=30:[Bob Tom] >30:[Alice Chris]]
```

#### Reducers

##### sum.Of

``` go
import (
    "github.com/m4gshm/gollections/loop"
    "github.com/m4gshm/gollections/loop/sum"
)

var sum = sum.Of(loop.Of(1, 2, 3, 4, 5, 6)) //21
```

##### loop.Reduce

``` go
var sum = loop.Reduce(loop.Of(1, 2, 3, 4, 5, 6), func(i1, i2 int) int { return i1 + i2 })
//21
```

##### loop.ReduceOK

``` go
adder := func(i1, i2 int) int { return i1 + i2 }

sum, ok := loop.ReduceOK(loop.Of(1, 2, 3, 4, 5, 6), adder)
//21, true

emptyLoop := loop.Of[int]()
sum, ok = loop.ReduceOK(emptyLoop, adder)
//0, false
```

##### loop.Accum

``` go
import (
    "github.com/m4gshm/gollections/loop"
    "github.com/m4gshm/gollections/op"
)

var sum = loop.Accum(100, loop.Of(1, 2, 3, 4, 5, 6), op.Sum)
//121
```

##### loop.First

``` go
import (
    "github.com/m4gshm/gollections/predicate/more"
    "github.com/m4gshm/gollections/loop"
)

result, ok := loop.First(loop.Of(1, 3, 5, 7, 9, 11), more.Than(5)) //7, true
```

#### Element converters

##### loop.Convert

``` go
var s []string = loop.Convert(loop.Of(1, 3, 5, 7, 9, 11), strconv.Itoa).Slice()
//[]string{"1", "3", "5", "7", "9", "11"}
```

##### loop.Conv

``` go
result, err := loop.Conv(loop.Of("1", "3", "5", "_7", "9", "11"), strconv.Atoi).Slice()
//[]int{1, 3, 5}, ErrSyntax
```

#### Loop converters

##### loop.Filter

``` go
import (
    "github.com/m4gshm/gollections/predicate/exclude"
    "github.com/m4gshm/gollections/predicate/one"
    "github.com/m4gshm/gollections/loop"
)

var f1 = loop.Filter(loop.Of(1, 3, 5, 7, 9, 11), one.Of(1, 7).Or(one.Of(11))).Slice()
//[]int{1, 7, 11}

var f2 = loop.Filter(loop.Of(1, 3, 5, 7, 9, 11), exclude.All(1, 7, 11)).Slice()
//[]int{3, 5, 9}
```

##### loop.Flat

``` go
import (
    "github.com/m4gshm/gollections/convert/as"
    "github.com/m4gshm/gollections/loop"
)

var i []int = loop.Flat(loop.Of([][]int{{1, 2, 3}, {4}, {5, 6}}...), as.Is).Slice()
//[]int{1, 2, 3, 4, 5, 6}
```

#### Operations chain functions

- convert.AndReduce, conv.AndReduce

- convert.AndFilter

- filter.AndConvert

These functions combine converters, filters and reducers.

### Iterating over loops

- Using rangefunc `All` like:

``` go
for i := range range_.Of(0, 100).All {
    doOp(i)
}
```

- Using `for` statement like:

``` go
next := range_.Of(0, 100)
for i, ok := next(); ok; i, ok = next() {
    doOp(i)
}
```

- or

``` go
for next, i, ok := range_.Of(0, 100).Crank(); ok; i, ok = next() {
    doOp(i)
}
```

- `ForEach` method

``` go
range_.Of(0, 100).ForEach(doOp)
```

- or `For` method that can be aborted by returning `Break` for expected
  completion, or another error otherwise.

``` go
range_.Of(0, 100).For(func(i int) error {
    if i > 22 {
        return loop.Break
    }
    doOp(i)
    return loop.Continue
})
```

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
import (
    "github.com/m4gshm/gollections/predicate/where"
    "github.com/m4gshm/gollections/slice"
)

bob, _ := slice.First(users, where.Eq(User.Name, "Bob"))
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

### Immutable collections

The same underlying interfaces but for read-only use cases.

### Iterating over collections

- Using rangefunc `All` like:

``` go
uniques := set.From(range_.Of(0, 100))
for i := range uniques.All {
    doOp(i)
}
```

- `ForEach` method

``` go
uniques := set.From(range_.Of(0, 100))
uniques.ForEach(doOp)
```

- or `For` method that can be aborted by returning `Break` for expected
  completion, or another error otherwise.

``` go
uniques := set.From(range_.Of(0, 100))
uniques.For(func(i int) error {
    if i > 22 {
        return loop.Break
    }
    doOp(i)
    return loop.Continue
})
```
