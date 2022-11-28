package simpleflow

// ExtractToSlice calls the func `fn` for each element in `in` and appends the result to `out` only if the second
// return argument of `fn` is true.
func ExtractToSlice[T, V any](in []T, fn func(T) (V, bool), out []V) []V {
	for _, v := range in {
		s, ok := fn(v)
		if ok {
			out = append(out, s)
		}
	}
	return out
}

// ExtractToChannel calls the func `fn` for each element in `in` and pushes the result to `out` only if the second
// return argument of `fn` is true.
func ExtractToChannel[T, V any](in []T, fn func(T) (V, bool), out chan V) {
	for _, v := range in {
		s, ok := fn(v)
		if ok {
			out <- s
		}
	}
	return
}
