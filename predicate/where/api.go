package where

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/eq"
)

func Match[From, To any](structGetter func(From) To, matcher predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.Match(structGetter, matcher)
}

func Any[From, To any](structGetter func(From) []To, matcher predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.MatchAny(structGetter, matcher)
}

func Eq[From any, To comparable](structGetter func(From) To, example To) predicate.Predicate[From] {
	return Match(structGetter, eq.To(example))
}

func Not[From, To any](structGetter func(From) To, matcher predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.Not(Match(structGetter, matcher))
}
