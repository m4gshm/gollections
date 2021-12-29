package iter

import (
	"github.com/m4gshm/container/typ"
)

type FlattenFit[From, To any] struct {
	Iter  typ.Iterator[From]
	Flatt typ.Flatter[From, To]
	Fit   typ.Predicate[From]

	iterTo []To
	indTo  int
	c      To
	err    error
}

var _ typ.Iterator[interface{}] = (*FlattenFit[interface{}, interface{}])(nil)

func (s *FlattenFit[From, To]) HasNext() bool {
	if iterTo := s.iterTo; len(iterTo) > 0 {
		if indTo := s.indTo; indTo < len(iterTo) {
			c := iterTo[indTo]
			s.c = c
			s.indTo = indTo + 1
			return true
		} else {
			s.indTo = 0
			s.iterTo = nil
		}
	}

	iter := s.Iter
	for {
		if !iter.HasNext() {
			s.err = Exhausted
			return false
		} else if err := iter.Err(); err != nil {
			s.err = err
			return true
		} else if v := iter.Get(); s.Fit(v) {
			if iterTo := s.Flatt(v); len(iterTo) > 0 {
				s.c = iterTo[0]
				s.iterTo = iterTo
				s.indTo = 1
				return true
			}
		}
	}
}

func (s *FlattenFit[From, To]) Get() To {
	return s.c
}

func (s *FlattenFit[From, To]) Err() error {
	return s.err
}

type Flatten[From, To any] struct {
	Iter   typ.Iterator[From]
	Flatt  typ.Flatter[From, To]
	iterTo []To
	indTo  int
	c      To
	err    error
}

var _ typ.Iterator[interface{}] = (*Flatten[interface{}, interface{}])(nil)

func (s *Flatten[From, To]) HasNext() bool {
	if iterTo := s.iterTo; len(iterTo) > 0 {
		if indTo := s.indTo; indTo < len(iterTo) {
			c := iterTo[indTo]
			s.c = c
			s.indTo = indTo + 1
			return true
		} else {
			s.indTo = 0
			s.iterTo = nil
		}
	}

	iter := s.Iter
	for {
		if ok := iter.HasNext(); !ok {
			s.err = Exhausted
			return false
		} else if err := iter.Err(); err != nil {
			s.err = iter.Err()
			return true
		}
		v := iter.Get()
		if iterTo := s.Flatt(v); len(iterTo) > 0 {
			s.c = iterTo[0]
			s.iterTo = iterTo
			s.indTo = 1
			return true
		}
	}
}

func (s *Flatten[From, To]) Get() To {
	return s.c
}

func (s *Flatten[From, To]) Err() error {
	return s.err
}

type FlattFitSlice[From, To any] struct {
	Elements []From
	Flatt    typ.Flatter[From, To]
	Fit      typ.Predicate[From]

	indFrom    int
	elementsTo []To
	indTo      int
	c          To
	err        error
}

var _ typ.Iterator[interface{}] = (*FlattFitSlice[interface{}, interface{}])(nil)

func (s *FlattFitSlice[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.c = c
			s.indTo = indTo + 1
			return true
		} else {
			s.indTo = 0
			s.elementsTo = nil
		}
	}

	elements := s.Elements
	le := len(elements)
	for indFrom := s.indFrom; indFrom < le; indFrom++ {
		s.indFrom = indFrom + 1
		if v := elements[indFrom]; s.Fit(v) {
			if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
				s.c = elementsTo[0]
				s.elementsTo = elementsTo
				s.indTo = 1
				return true
			}
		}
	}
	s.err = Exhausted
	return false
}

func (s *FlattFitSlice[From, To]) Get() To {
	if err := s.err; err != nil {
		panic(err)
	}
	return s.c
}

func (s *FlattFitSlice[From, To]) Err() error {
	return s.err
}

type FlattSlice[From, To any] struct {
	Elements []From
	Flatt    typ.Flatter[From, To]

	indFrom    int
	elementsTo []To
	indTo      int
	c          To
	err        error
}

var _ typ.Iterator[interface{}] = (*FlattSlice[interface{}, interface{}])(nil)

func (s *FlattSlice[From, To]) HasNext() bool {
	if elementsTo := s.elementsTo; len(elementsTo) > 0 {
		if indTo := s.indTo; indTo < len(elementsTo) {
			c := elementsTo[indTo]
			s.c = c
			s.indTo = indTo + 1
			return true
		} else {
			s.indTo = 0
			s.elementsTo = nil
		}
	}

	elements := s.Elements
	le := len(elements)
	for indFrom := s.indFrom; indFrom < le; indFrom++ {
		s.indFrom = indFrom + 1
		v := elements[indFrom]
		if elementsTo := s.Flatt(v); len(elementsTo) > 0 {
			s.c = elementsTo[0]
			s.elementsTo = elementsTo
			s.indTo = 1
			return true
		}
	}
	s.err = Exhausted
	return false
}

func (s *FlattSlice[From, To]) Get() To {
	if err := s.err; err != nil {
		panic(err)
	}
	return s.c
}

func (s *FlattSlice[From, To]) Err() error {
	return s.err
}
