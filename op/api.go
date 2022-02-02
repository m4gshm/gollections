package op

//Binary is an operation with two arguments
type Binary[T any] func(T, T) T

//Quaternary is an operation with four arguments
type Quaternary[t1, t2 any] func(t1, t2, t1, t2) (t1, t2)
