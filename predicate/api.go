package predicate

// Predicate tests value (converts to true or false).
type Predicate[T any] func(T) bool

func (p Predicate[T]) Or(or Predicate[T]) Predicate[T]   { return Or(p, or) }
func (p Predicate[T]) And(and Predicate[T]) Predicate[T] { return And(p, and) }
func (p Predicate[T]) Xor(and Predicate[T]) Predicate[T] { return Xor(p, and) }

// Eq makes a predicate to test for equality
func Eq[T comparable](v T) Predicate[T] {
	return func(c T) bool { return v == c }
}

// Not inverts a predicate
func Not[T any](p Predicate[T]) Predicate[T] {
	return func(v T) bool { return !p(v) }
}

// And makes a conjunction of two predicates
func And[T any](p1, p2 Predicate[T]) Predicate[T] {
	return func(v T) bool { return p1(v) && p2(v) }
}

// Or makes a disjunction of two predicates
func Or[T any](p1, p2 Predicate[T]) Predicate[T] {
	return func(v T) bool { return p1(v) || p2(v) }
}

// Xor makes an exclusive OR of two predicates
func Xor[T any](p1, p2 Predicate[T]) Predicate[T] {
	return func(v T) bool { return !(p1(v) == p2(v)) }
}

// Union applies And to predicates
func Union[T any](predicates ...Predicate[T]) Predicate[T] {
	l := len(predicates)
	if l == 0 {
		return func(_ T) bool { return false }
	} else if l == 1 {
		return predicates[0]
	} else if l == 2 {
		return And(predicates[0], predicates[1])
	}
	return func(v T) bool {
		for i := 0; i < len(predicates); i++ {
			if !predicates[i](v) {
				return false
			}
		}
		return true
	}
}

// Always returns v every time.
func Always[T any](v bool) Predicate[T] {
	return func(_ T) bool { return v }
}

// Never returns the negative of v every time
func Never[T any](v bool) Predicate[T] {
	return func(_ T) bool { return !v }
}

// BiPredicate tests values pair (converts to true or false).
type BiPredicate[v1, v2 any] func(v1, v2) bool
