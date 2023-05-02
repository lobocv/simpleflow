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

// ExtractFirst returns the first element in the slice for which fn(element) == true.
// If no matches are found, the second return argument is false.
func ExtractFirst[T any](values []T, fn func(T) bool) (v T, exists bool) {
	for ii := 0; ii < len(values); ii++ {
		if fn(values[ii]) {
			return values[ii], true
		}
	}

	return v, false
}
