package simpleflow

// Transform applies a transformation function to each element in the input slice and returns
// a new slice
func Transform[T, V any](values []T, fn func(t T) V) []V {
	var out = make([]V, 0, len(values))
	for ii := 0; ii < len(values); ii++ {
		out = append(out, fn(values[ii]))
	}
	return out
}

// TransformAndFilter applies a transformation function to each element in the input slice.
// If the second return value of the transformation function is false, then the value will be
// omitted from the output.
func TransformAndFilter[T, V any](values []T, fn func(t T) (V, bool)) []V {
	var out = make([]V, 0, len(values))
	for ii := 0; ii < len(values); ii++ {
		v, ok := fn(values[ii])
		if !ok {
			continue
		}
		out = append(out, v)
	}
	return out
}
