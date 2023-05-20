// Package use provides conditional expression builders
package use

// If builds use.If(condition, tru).Else(fals) expression builder.
// Looks like val := use.If(condition, valOnTrue).Else(defaltVal) tha can be rewrited by:
//
//	var val type
//	if condtion {
//		val = valOnTrue
//	} else {
//		val = defaltVal
//	}
func If[T any](condition bool, tru T) When[T] {
	return newWhen(condition, tru)
}

// IfGet is like If but aimed to use an getter function
func IfGet[T any](condition bool, then func() T) When[T] {
	return If(condition, evaluate(condition, then))
}

// IfGetErr is like If but aimed to use an error return function
func IfGetErr[T any](condition bool, tru func() (T, error)) WhenErr[T] {
	then, err := evaluateErr(condition, tru)
	return newWhenErr(condition, then, err)
}

// If_ is alias of IfErr
func If_[T any](condition bool, tru func() (T, error)) WhenErr[T] {
	return IfGetErr(condition, tru)
}
