// Package predicate provides breakable predicate builders
package predicate

// Predicate tests value (converts to true or false) or aborts by an error.
type Predicate[T any] func(T) (bool, error)

// Or makes disjunction
func (p Predicate[T]) Or(or Predicate[T]) Predicate[T] { return Or(p, or) }

// And makes conjunction
func (p Predicate[T]) And(and Predicate[T]) Predicate[T] { return And(p, and) }

// Xor makes exclusive OR
func (p Predicate[T]) Xor(xor Predicate[T]) Predicate[T] { return Xor(p, xor) }

// Of breakable Predicate constructor
func Of[T comparable](predicate func(T) bool) Predicate[T] {
	return func(c T) (bool, error) { return predicate(c), nil }
}

// Wrap converts the specified predicate to the erroreable one
func Wrap[T any](predicate func(T) bool) Predicate[T] {
	return func(t T) (bool, error) { return predicate(t), nil }
}

// Eq creates a predicate to test for equality
func Eq[T comparable](v T) Predicate[T] {
	return func(c T) (bool, error) { return v == c, nil }
}

// Not inverts a predicate
func Not[T any](p Predicate[T]) Predicate[T] {
	return func(v T) (bool, error) {
		ok, err := p(v)
		return !ok, err
	}
}

// And makes a conjunction of two predicates
func And[T any](p1, p2 Predicate[T]) Predicate[T] {
	return func(v T) (bool, error) {
		if ok, err := p1(v); err != nil || !ok {
			return false, err
		}
		return p2(v)
	}
}

// Or makes a disjunction of two predicates
func Or[T any](p1, p2 Predicate[T]) Predicate[T] {
	return func(v T) (bool, error) {
		if ok, err := p1(v); err != nil || ok {
			return ok, err
		}
		return p2(v)
	}
}

// Xor makes an exclusive OR of two predicates
func Xor[T any](p1, p2 Predicate[T]) Predicate[T] {
	return func(v T) (bool, error) {
		ok, err := p1(v)
		if err != nil {
			return ok, err
		}
		ok2, err := p2(v)
		if err != nil {
			return ok2, err
		}
		return ok != ok2, nil
	}
}

// Union applies And to predicates
func Union[T any](predicates ...Predicate[T]) Predicate[T] {
	l := len(predicates)
	switch l {
	case 0:
		return func(_ T) (bool, error) { return false, nil }
	case 1:
		return predicates[0]
	case 2:
		return And(predicates[0], predicates[1])
	}
	return func(v T) (bool, error) {
		for i := range predicates {
			if ok, err := predicates[i](v); err != nil {
				return ok, err
			} else if !ok {
				return false, nil
			}
		}
		return true, nil
	}
}

// Always returns v every time.
func Always[T any](v bool) Predicate[T] {
	return func(_ T) (bool, error) { return v, nil }
}

// Never returns the negative of v every time
func Never[T any](v bool) Predicate[T] {
	return func(_ T) (bool, error) { return !v, nil }
}
