// Package get provides conditional expression builders
package get

// If builds a head of getter function call by condition.
// Looks like val := get.If(condition, getterFuncOnTrue).Else(defaltVal) tha can be rewrited by:
//
//		var val type
//	 if condtion {
//			val = getterFuncOnTrue()
//		} else {
//			val = defaltVal
//		}
func If[T any](condition bool, tru func() T) When[T] {
	return When[T]{condition, tru}
}

// IfErr is like If but aimed to use an error return function
func IfErr[T any](condition bool, tru func() (T, error)) WhenErr[T] {
	return WhenErr[T]{condition, tru}
}

// If_ is alias of IfErr
func If_[T any](condition bool, tru func() (T, error)) WhenErr[T] {
	return IfErr(condition, tru)
}
