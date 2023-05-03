// Package predicate provides predicate builders
package predicate

import (
	"github.com/m4gshm/gollections/slice"
)

// Predicate tests value (converts to true or false).
type Predicate[T any] func(T) bool

// Or makes disjunction
func (p Predicate[T]) Or(or Predicate[T]) Predicate[T] { return Or(p, or) }

// And makes conjunction
func (p Predicate[T]) And(and Predicate[T]) Predicate[T] { return And(p, and) }

// Xor makes exclusive OR
func (p Predicate[T]) Xor(xor Predicate[T]) Predicate[T] { return Xor(p, xor) }

// Eq creates a predicate to test for equality
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

func Match[From, To any](convert func(From) To, predicate Predicate[To]) Predicate[From] {
	return func(from From) bool { return predicate(convert(from)) }
}

func MatchAny[From, To any](flatter func(From) []To, predicate Predicate[To]) Predicate[From] {
	return func(from From) bool {
		return slice.Has(flatter(from), predicate)
	}
}

func HasConverted[From, I, To any](flatter func(From) []I, convert func([]I) To, predicate Predicate[To]) Predicate[From] {
	return func(from From) bool {
		return predicate(convert(flatter(from)))
	}
}

func HasAnyConverted[From, I, To any](flatter func(From) []I, convert func(I) To, predicate Predicate[To]) Predicate[From] {
	return func(from From) bool {
		for _, f := range flatter(from) {
			if c := convert(f); predicate(c) {
				return true
			}
		}
		return false
	}
}

func ContainsConverted[From, I any, To comparable](flatter func(From) []I, convert func(I) To, expected To) Predicate[From] {
	return func(from From) bool {
		ff := flatter(from)
		for _, f := range ff {
			if c := convert(f); c == expected {
				return true
			}
		}
		return false
	}
}

func Contains[From any, To comparable](flatter func(From) []To, expected To) Predicate[From] {
	return func(from From) bool {
		return slice.Contains(flatter(from), expected)
	}
}

func Len[TS ~[]T, T any](predicate Predicate[int]) Predicate[TS] {
	return Match(slice.Len[TS], predicate)
}

func Any[TS ~[]T, T, C any](convert func(TS) C, predicate Predicate[C]) Predicate[TS] {
	return Match(convert, predicate)
}

func Empty[From, To any](flattener func(From) []To) Predicate[From] {
	return func(f From) bool { return slice.Empty(flattener(f)) }
}
