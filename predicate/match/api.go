package match

import "github.com/m4gshm/gollections/predicate"

func To[From, To any](convert func(From) To, matcher predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.Match(convert, matcher)
}

func Any[From, To any](flatter func(From) []To, matcher predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.MatchAny(flatter, matcher)
}
