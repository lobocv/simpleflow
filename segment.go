package simpleflow

// SegmentFunc is a function that determines how an item of type `T` is segmented into a segment of type `S`
type SegmentFunc[T any, S comparable] func(T) S

type SegmentFuncKV[K comparable, V any, S comparable] func(K, V) S

// SegmentSlice takes a slice and breaks it into smaller segmented slices using the provided function `f`
func SegmentSlice[T any, S comparable](items []T, f SegmentFunc[T, S]) map[S][]T {

	segments := make(map[S][]T)

	for ii := 0; ii < len(items); ii++ {
		s := f(items[ii])

		// Check if a segment has already been created, if not, create one
		if _, ok := segments[s]; !ok {
			segments[s] = []T{}
		}
		// Add the item to the segment
		segments[s] = append(segments[s], items[ii])
	}

	return segments
}

// SegmentMap takes a map and breaks it into smaller segmented maps using the provided function `f`
func SegmentMap[K comparable, V any, S comparable](items map[K]V, f SegmentFuncKV[K, V, S]) map[S]map[K]V {

	segments := make(map[S]map[K]V)

	for k, v := range items {
		s := f(k, v)

		// Check if a segment has already been created, if not, create one
		if _, ok := segments[s]; !ok {
			segments[s] = map[K]V{}
		}
		// Add the item to the segment
		segments[s][k] = v
	}

	return segments
}
