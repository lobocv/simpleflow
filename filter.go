package simpleflow

// FilterSliceInplace filters the input slice with the function `fn` in-place
func FilterSliceInplace[V any](in []V, fn func(v V) bool) []V {
	return FilterSliceInto(in, in[:0], fn)
}

// FilterSlice filters the input slice with the function `fn` and return a new slice
func FilterSlice[V any](in []V, fn func(v V) bool) []V {
	var out []V
	return FilterSliceInto(in, out, fn)
}

// FilterSliceInto filters the input slice `in` with the function `fn` into `out`
func FilterSliceInto[V any](in, out []V, fn func(v V) bool) []V {

	var count int
	for ii := 0; ii < len(in); ii++ {

		if !fn(in[ii]) {
			continue
		}

		out = append(out, in[ii])
		count++
	}

	return out[:count]

}
