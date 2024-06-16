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

// Not creates a 'not p' predicate
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

// Match creates a predicate that tests whether a value of a structure property matches a specified condition
func Match[Entity, Property any](getter func(Entity) Property, condition Predicate[Property]) Predicate[Entity] {
	return func(entity Entity) bool { return condition(getter(entity)) }
}

// MatchAny creates a predicate that tests whether any value of a structure property matches a specified condition
// The property has a slice type.
func MatchAny[Entity, Property any](getter func(Entity) []Property, condition Predicate[Property]) Predicate[Entity] {
	return func(entity Entity) bool {
		return slice.Has(getter(entity), condition)
	}
}
