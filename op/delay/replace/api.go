// Package replace provides builders that specifies on value replacing
package replace

// By creates replacer a value by the specified one
func By[T any](r T) func(t T) T {
	return func(_ T) T { return r }
}
