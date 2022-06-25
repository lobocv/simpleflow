package simpleflow

// BatchSlice takes a slice and breaks it up into sub-slices of `size` length each
func BatchSlice[T any](items []T, size int) [][]T {
	batches := make([][]T, 0, (len(items)+size-1)/size)

	for size < len(items) {
		items, batches = items[size:], append(batches, items[0:size:size])
	}
	batches = append(batches, items)

	return batches
}

// BatchMap takes a map and breaks it up into sub-maps of `size` keys each
func BatchMap[K comparable, V any](items map[K]V, size int) []map[K]V {
	batches := make([]map[K]V, 0, (len(items)+size-1)/size)

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

// IncrementalBatchSlice incrementally builds slice batches of size `batchSize` by appending to a slice
// If the slice is larger than `batchSize` elements, a single batch is returned. The remaining
// elements of the slice are always returned. Batched items are returned from the head of the slice.
// To avoid errors on the caller side,  passing a batchSize < 1 will result in a batchSize of 1.
func IncrementalBatchSlice[T any](items []T, batchSize int, v T) (remaining, batch []T) {
	// prevent bugs on the caller side by using a minimum of 1 batchSize
	if batchSize < 1 {
		batchSize = 1
	}

	items = append(items, v)
	if len(items) >= batchSize {
		remaining = items[batchSize:]
		batch = items[:batchSize]

		return remaining, batch
	}

	return items, nil
}

// IncrementalBatchMap incrementally builds map batches of size `batchSize` by adding elements to a map
// If the map is larger than `batchSize` elements, a single batch is returned along with the remaining
// elements of the map. Batched items are chosen by iterating the (unordered) map and thus you cannot make
// assumptions on which keys will exist in the batch.
// To avoid errors on the caller side,  passing a batchSize < 1 will result in a batchSize of 1.
func IncrementalBatchMap[K comparable, V any](items map[K]V, batchSize int, k K, v V) (batch map[K]V) {
	// prevent bugs on the caller side by using a minimum of 1 batchSize
	if batchSize < 1 {
		batchSize = 1
	}
	items[k] = v
	if len(items) >= batchSize {
		batch = make(map[K]V, batchSize)
		var count int

		// iterate the map and pop off batchSize keys to return
		// Since map order is indeterminate, you cannot know which elements
		// will be returned in the batch
		for kk, vv := range items {
			batch[kk] = vv
			delete(items, kk)
			count++
			if count == batchSize {
				break
			}

		}
		return batch
	}

	return nil
}
