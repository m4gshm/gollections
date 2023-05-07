package curry

func Of[In1, Out any](f func(In1) Out, arg In1) func() Out {
	return func() Out { return f(arg) }
}

func First[In1, In2, Out any](f func(In1, In2) Out, first In1) func(In2) Out {
	return func(second In2) Out { return f(first, second) }
}

func Second[In1, In2, Out any](f func(In1, In2) Out, second In2) func(In1) Out {
	return func(first In1) Out { return f(first, second) }
}
