# Gollections

Golang experiments with container data structures, slices, maps, and so on using generics.
Need use version 1.18 beta 1 or newer.

## Containers

### [Immutable](./immutable/api.go)

[Vector](./immutable/vector.go)

[OrderedMap](./immutable/map.go)

[OrderedSet](./immutable/set.go)


### [Mutable](./mutable/api.go)

[Vector](./mutable/vector.go)

[OrderedMap](./mutable/map.go)

[OrderedSet](./mutable/set.go)


## Packages
### [Interfaces](./typ/iface.go)
```go
package examples

import (
	"fmt"
	"testing"

	"github.com/m4gshm/container/conv"
	"github.com/m4gshm/container/immutable"
	it "github.com/m4gshm/container/iter"
	"github.com/m4gshm/container/op"
	slc "github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/walk"
	"github.com/stretchr/testify/assert"
)

func Test_OrderedSet(t *testing.T) {
	set := immutable.NewOrderedSet(1, 1, 2, 4, 3, 1)
	values := set.Values()
	fmt.Println(set) //[1, 2, 4, 3]

	assert.Equal(t, slc.Of(1, 2, 4, 3), values)
}

func Test_group_odd_even(t *testing.T) {
	var (
		even   = func(v int) bool { return v%2 == 0 }
		groups = walk.Group(immutable.NewOrderedSet(1, 1, 2, 4, 3, 1), even)
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