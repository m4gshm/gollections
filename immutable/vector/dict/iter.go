package dict

import (
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/typ"
)

func NewIterator[k comparable, v any](elements []*k, uniques map[k]v) *Iterator[k, v] {
	return &Iterator[k, v]{elements: elements, uniques: uniques}
}

type Iterator[k comparable, v any] struct {
	elements []*k
	uniques  map[k]v

	err     error
	current int
}

var _ typ.Iterator[any] = (*Iterator[any, any])(nil)

func (s *Iterator[k, v]) HasNext() bool {
	return it.HasNext(s.elements, &s.current, &s.err)
}

func (s *Iterator[k, v]) Get() (v, error) {
	kref, err := it.Get(s.current, s.elements, s.err)
	if err != nil {
		var no v
		return no, err
	}
	return s.uniques[*kref], nil
}
