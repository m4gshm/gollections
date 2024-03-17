package examples

import (
	"github.com/m4gshm/gollections/collection/mutable"
	mmap "github.com/m4gshm/gollections/collection/mutable/map_"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	omap "github.com/m4gshm/gollections/collection/mutable/ordered/map_"
	oset "github.com/m4gshm/gollections/collection/mutable/ordered/set"
	"github.com/m4gshm/gollections/collection/mutable/set"
	"github.com/m4gshm/gollections/collection/mutable/vector"
	"github.com/m4gshm/gollections/k"
)

func _() {
	capacity := 10

	var (
		_ *mutable.Vector[int]                            = vector.Of(1, 2, 3)
		_ *mutable.Vector[int]                            = new(mutable.Vector[int])
		_ *mutable.Vector[int]                            = vector.NewCap[int](capacity)
	)
	var (
		_ *mutable.Set[int]   = set.Of(1, 2, 3)
		_ *mutable.Set[int]   = new(mutable.Set[int])
		_ *mutable.Set[int]   = set.NewCap[int](capacity)
		_ *mutable.Set[int] = set.Empty[int]()
	)
	var (
		_ *ordered.Set[int]   = oset.Of(1, 2, 3)
		_ *ordered.Set[int]   = new(ordered.Set[int])
		_ *ordered.Set[int]   = oset.NewCap[int](capacity)
		_ *ordered.Set[int] = oset.Empty[int]()
	)
	var (
		_ *mutable.Map[int, string]   = mmap.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
		_ *mutable.Map[int, string]   = new(mutable.Map[int, string])
		_ *mutable.Map[int, string]   = mmap.New[int, string](capacity)
		_ *mutable.Map[int, string] = mmap.Empty[int, string]()
	)
	var (
		_ *ordered.Map[int, string]   = omap.Of(k.V(1, "1"), k.V(2, "2"), k.V(3, "3"))
		_ *ordered.Map[int, string]   = new(ordered.Map[int, string])
		_ *ordered.Map[int, string]   = omap.New[int, string](capacity)
		_ *ordered.Map[int, string] = omap.Empty[int, string]()
	)
}
