// Package examples of collection constructors
package examples

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/map_"
	"github.com/m4gshm/gollections/immutable/omap"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/immutable/oset"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/immutable/vector"
	"github.com/m4gshm/gollections/k"
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
		_ immutable.Map[int, string] = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
		_ c.Map[int, string]         = map_.New(map[int]string{1: "2", 2: "2", 3: "3"})
	)
	var (
		_ ordered.Map[int, string] = omap.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
		_ c.Map[int, string]       = omap.New(
			/*uniques*/ map[int]string{1: "2", 2: "2", 3: "3"} /*order*/, []int{3, 1, 2},
		)
	)
}
