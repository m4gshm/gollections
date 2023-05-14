// Package get provides conditional expression builders
package get

import "github.com/m4gshm/gollections/expr/use"

// If builds a head of getter function call by condition.
// Looks like val := get.If(condition, getterFuncOnTrue).Else(defaltVal) tha can be rewrited by:
//
//	var val type
//	if condtion {
//		val = getterFuncOnTrue()
//	} else {
//		val = defaltVal
//	}
func If[T any](condition bool, then func() T) use.When[T] {
	return use.IfGet(condition, then)
}

// IfErr is like If but aimed to use an error return function
func IfErr[T any](condition bool, then func() (T, error)) use.WhenErr[T] {
	return use.IfGetErr(condition, then)
}

// If_ is alias of IfErr
func If_[T any](condition bool, tru func() (T, error)) use.WhenErr[T] {
	return IfErr(condition, tru)
}
