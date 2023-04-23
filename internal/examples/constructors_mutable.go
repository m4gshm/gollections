package examples

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/k"
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
		_ *mutable.Vector[int] = new(mutable.Vector[int])
		_ c.Vector[int]        = vector.NewCap[int](capacity)
		_ c.Vector[int]        = vector.Empty[int]()
	)
	var (
		_ *mutable.Set[int] = set.Of(1, 2, 3)
		_ *mutable.Set[int] = new(mutable.Set[int])
		_ c.Set[int]        = set.NewCap[int](capacity)
		_ c.Set[int]        = set.Empty[int]()
	)
	var (
		_ *ordered.Set[int] = oset.Of(1, 2, 3)
		_ *ordered.Set[int] = new(ordered.Set[int])
		_ c.Set[int]        = oset.NewCap[int](capacity)
		_ c.Set[int]        = oset.Empty[int]()
	)
	var (
		_ *mutable.Map[int, string] = map_.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
		_ *mutable.Map[int, string] = new(mutable.Map[int, string])
		_ c.Map[int, string]        = map_.New[int, string](capacity)
		_ c.Map[int, string]        = map_.Empty[int, string]()
	)
	var (
		_ *ordered.Map[int, string] = omap.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
		_ *ordered.Map[int, string] = new(ordered.Map[int, string])
		_ c.Map[int, string]        = omap.New[int, string](capacity)
		_ c.Map[int, string]        = omap.Empty[int, string]()
	)
}
