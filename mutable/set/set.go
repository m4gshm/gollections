package set

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func Convert[T comparable](elements []T) *Set[T] {
	uniques := make(map[T]struct{}, 0)
	for _, v := range elements {
		uniques[v] = struct{}{}
	}
	return Wrap(uniques)
}

func Wrap[k comparable](uniques map[k]struct{}) *Set[k] {
	return &Set[k]{uniques: uniques}
}

type Set[k comparable] struct {
	uniques    map[k]struct{}
	err        error
	changeMark int32
}

var _ mutable.Set[any, mutable.Iterator[any]] = (*Set[any])(nil)
var _ typ.Set[any, mutable.Iterator[any]] = (*Set[any])(nil)
var _ fmt.Stringer = (*Set[any])(nil)
var _ fmt.GoStringer = (*Set[any])(nil)

func (s *Set[k]) Begin() mutable.Iterator[k] {
	return s.Iter()
}

func (s *Set[k]) Iter() *Iter[k] {
	return NewIter(s.uniques, s.DeleteOne)
}

func (s *Set[k]) Add(elements ...k) (bool, error) {
	return s.AddAll(elements)
}

func (s *Set[k]) AddAll(elements []k) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	markOnStart := s.changeMark
	uniques := s.uniques
	added := false
	for _, element := range elements {
		if _, ok := uniques[element]; !ok {
			uniques[element] = struct{}{}
			added = true
		}
	}
	if !added {
		return false, nil
	}
	return mutable.Commit(markOnStart, &s.changeMark, &s.err)
}

func (s *Set[k]) AddOne(element k) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	markOnStart := s.changeMark
	uniques := s.uniques
	if _, ok := uniques[element]; ok {
		return false, nil
	}
	uniques[element] = struct{}{}

	return mutable.Commit(markOnStart, &s.changeMark, &s.err)
}

func (s *Set[k]) Delete(elements ...k) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}

	markOnStart := s.changeMark
	uniques := s.uniques
	for _, element := range elements {
		if _, ok := uniques[element]; !ok {
			return false, nil
		}

		delete(uniques, element)
	}
	return mutable.Commit(markOnStart, &s.changeMark, &s.err)
}

func (s *Set[k]) DeleteOne(element k) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}

	markOnStart := s.changeMark
	uniques := s.uniques
	if _, ok := uniques[element]; !ok {
		return false, nil
	}

	delete(uniques, element)
	return mutable.Commit(markOnStart, &s.changeMark, &s.err)
}

func (s *Set[k]) Collect() []k {
	uniques := s.uniques
	out := make([]k, 0, len(uniques))
	for e := range uniques {
		out = append(out, e)
	}
	return out
}

func (s *Set[k]) For(walker func(k) error) error {
	return map_.ForKeys(s.uniques, walker)
}

func (s *Set[k]) ForEach(walker func(k)) error {
	return map_.ForEachKey(s.uniques, walker)
}

func (s *Set[k]) Filter(filter typ.Predicate[k]) typ.Pipe[k, []k, typ.Iterator[k]] {
	return it.NewPipe[k](it.Filter(s.Iter(), filter))
}

func (s *Set[k]) Map(by typ.Converter[k, k]) typ.Pipe[k, []k, typ.Iterator[k]] {
	return it.NewPipe[k](it.Map(s.Iter(), by))
}

func (s *Set[k]) Reduce(by op.Binary[k]) k {
	return it.Reduce(s.Iter(), by)
}

func (s *Set[k]) Len() int {
	return len(s.uniques)
}

func (s *Set[k]) Contains(val k) bool {
	_, ok := s.uniques[val]
	return ok
}

func (s *Set[k]) String() string {
	return s.GoString()
}

func (s *Set[k]) GoString() string {
	return slice.ToString(s.Collect())
}
