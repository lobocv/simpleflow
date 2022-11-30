package simpleflow

// FilterSliceInplace filters the input slice with the function `fn` in-place
// The function `fn` accepts a value and should return true to keep the value in the slice
func FilterSliceInplace[V any](in []V, fn func(v V) bool) []V {
	return FilterSliceInto(in, in[:0], fn)
}

// FilterSlice filters the input slice with the function `fn` and return a new slice
// The function `fn` accepts a value and should return true to copy the value into the new slice
func FilterSlice[V any](in []V, fn func(v V) bool) []V {
	var out []V
	return FilterSliceInto(in, out, fn)
}

// FilterSliceInto filters the input slice `in` with the function `fn` into `out`
// The function `fn` accepts a value and should return true to copy the value into the `out` slice
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

// FilterMapInplace filters the input map with the function `fn` in-place
// The function `fn` accepts a key, value pair and should return true to keep the pair in the map
func FilterMapInplace[K comparable, V any](in map[K]V, fn func(k K, v V) bool) {

	for k, v := range in {
		if !fn(k, v) {
			delete(in, k)
		}

	}
}

// FilterMap filters the input map with the function `fn` and return a new map
// The function `fn` accepts a key, value pair and should return true to keep the pair in the map
func FilterMap[K comparable, V any](in map[K]V, fn func(k K, v V) bool) map[K]V {
	var out = map[K]V{}
	for k, v := range in {
		if fn(k, v) {
			out[k] = v
		}
	}

	return out
}
