package c

// Converter convert From -> To.
type Converter[From, To any] func(From) To

// BiConverter convert pairs of From -> To.
type BiConverter[From1, From2, To1, To2 any] func(From1, From2) (To1, To2)

// Flatter extracts slice of To.
type Flatter[From, To any] Converter[From, []To]

// Predicate tests value (converts to true or false).
type Predicate[T any] func(T) bool

// BiPredicate tests values pair (converts to true or false).
type BiPredicate[v1, v2 any] func(v1, v2) bool
