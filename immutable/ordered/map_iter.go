package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

func NewValIter[k comparable, v any](elements []k, uniques map[k]v) *ValIter[k, v] {
	return &ValIter[k, v]{elements: elements, uniques: uniques}
}

type ValIter[k comparable, v any] struct {
	elements []k
	uniques  map[k]v

	err     error
	current int
}

var _ c.Iterator[any] = (*ValIter[int, any])(nil)

func (s *ValIter[k, v]) HasNext() bool {
	return it.HasNext(&s.elements, &s.current, &s.err)
}

func (s *ValIter[k, v]) Get() (v, error) {
	key, err := it.Get(&s.elements, s.current, s.err)
	if err != nil {
		var no v
		return no, err
	}
	return s.uniques[key], nil
}

func (s *ValIter[k, v]) Next() v {
	return it.Next[v](s)
}
