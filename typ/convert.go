package typ

//Converter convert From -> To
type Converter[From, To any] func(From) To

//Flatter extracts slice of To
type Flatter[From, To any] Converter[From, []To]

//Predicate tests value (converts to true or false)
type Predicate[T any] Converter[T, bool]
