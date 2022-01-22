package concurgo

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
