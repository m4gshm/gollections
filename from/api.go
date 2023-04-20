package from

// To transforms converers chain of From->Internal->To into From->To
func To[From, Internal, To any](from func(From) Internal, to func(Internal) To) func(From) To {
	return func(v From) To { return to(from(v)) }
}
