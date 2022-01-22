package simpleflow

// BatchSlice takes a slice and breaks it up into sub-slices of `size` length each
func BatchSlice[T any](items []T, size int) [][]T {
	batches := make([][]T, 0, (len(items)/size)+1)

	batch := make([]T, 0, size)
	for ii := 0; ii < len(items); ii++ {
		batch = append(batch, items[ii])
		if len(batch) == size {
			batches = append(batches, batch)
			batch = make([]T, 0, size)
		}
	}
	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

// BatchMap takes a map and breaks it up into sub-maps of `size` keys each
func BatchMap[K comparable, V any](items map[K]V, size int) []map[K]V {
	batches := make([]map[K]V, 0, (len(items)/size)+1)

	batch := make(map[K]V, size)
	for k, v := range items {
		batch[k] = v
		if len(batch) == size {
			batches = append(batches, batch)
			batch = make(map[K]V, size)
		}
	}
	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

// BatchChan reads from a channel and pushes batches of size `size` onto the `to` channel
func BatchChan[T any](items <-chan T, size int, to chan []T) {
	batch := make([]T, 0, size)
	for v := range items {
		batch = append(batch, v)

		if len(batch) == size {
			to <- batch
			batch = make([]T, 0, size)
		}
	}
	if len(batch) > 0 {
		to <- batch
	}

	return
}
