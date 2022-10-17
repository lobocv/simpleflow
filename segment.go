package simpleflow

// SegmentFunc is a function that determines how an item of type `T` is segmented into a segment of type `S`
type SegmentFunc[T any, S comparable] func(T) S

type SegmentFuncKV[K comparable, V any, S comparable] func(K, V) S

// SegmentSlice takes a slice and breaks it into smaller segmented slices using the provided function `f`
// The segments are returned in a map where the segment is the key.
func SegmentSlice[T any, S comparable](items []T, f SegmentFunc[T, S]) map[S][]T {

	segments := make(map[S][]T)

	for ii := 0; ii < len(items); ii++ {
		IncrementalSegmentSlice(segments, items[ii], f)
	}

	return segments
}

// IncrementalSegmentSlice adds the item to the correct segment inside `segments by calling `f` on `item`
func IncrementalSegmentSlice[T any, S comparable](segments map[S][]T, item T, f SegmentFunc[T, S]) {
	s := f(item)

	// Check if a segment has already been created, if not, create one
	if _, ok := segments[s]; !ok {
		segments[s] = []T{}
	}
	// Add the item to the segment
	segments[s] = append(segments[s], item)

}

// SegmentMap takes a map and breaks it into smaller segmented maps using the provided function `f`
// The segments are returned in a map where the segment is the key.
func SegmentMap[K comparable, V any, S comparable](items map[K]V, f SegmentFuncKV[K, V, S]) map[S]map[K]V {

	segments := make(map[S]map[K]V)

	for k, v := range items {
		IncrementalSegmentMap(segments, k, v, f)
	}

	return segments
}

// IncrementalSegmentMap adds the (key,value) pair to the correct segment inside `segments by calling `f` on `(k, v)`
func IncrementalSegmentMap[K comparable, V any, S comparable](segments map[S]map[K]V, k K, v V, f SegmentFuncKV[K, V, S]) {
	s := f(k, v)

	// Check if a segment has already been created, if not, create one
	if _, ok := segments[s]; !ok {
		segments[s] = map[K]V{}
	}
	// Add the item to the segment
	segments[s][k] = v
}

// SegmentChan takes a channel and breaks it into smaller segmented slices using the provided function `f`
// The segments are returned in a map where the segment is the key.
func SegmentChan[T any, S comparable](items <-chan T, f SegmentFunc[T, S]) map[S][]T {

	segments := make(map[S][]T)

	for item := range items {
		s := f(item)

		// Check if a segment has already been created, if not, create one
		if _, ok := segments[s]; !ok {
			segments[s] = []T{}
		}
		// Add the item to the segment
		segments[s] = append(segments[s], item)
	}

	return segments
}
