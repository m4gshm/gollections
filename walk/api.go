package walk

import "github.com/m4gshm/container/typ"

func Group[T any, K comparable, W typ.Walk[T]](elements W, by typ.Converter[T, K]) map[K][]T {
	groups := map[K][]T{}
	elements.ForEach(func (e T)  {
		key := by(e)
		group := groups[key]
		if group == nil {
			group = make([]T, 0)
		}
		groups[key] = append(group, e)
	})
	return groups
}