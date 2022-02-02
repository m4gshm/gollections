//Package walk provides utilily functions for the interface Walker
package walk

import (
	"github.com/m4gshm/gollections/c"
)

//Group converts elements into the map containing slices of the elements separated by keys, which are retrieved using a Converter object.
func Group[T any, K comparable, W c.WalkEach[T]](elements W, by c.Converter[T, K]) map[K][]T {
	groups := map[K][]T{}
	elements.ForEach(func(e T) {
		key := by(e)
		group := groups[key]
		if group == nil {
			group = make([]T, 0)
		}
		groups[key] = append(group, e)
	})
	return groups
}
