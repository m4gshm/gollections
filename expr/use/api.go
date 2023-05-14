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
	return When[T]{condition, tru}
}
