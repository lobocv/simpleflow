package simpleflow

import (
	"context"
	"sync"
)

// KeyValue is a tuple of key, value
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// WorkerPoolFromMap starts a worker pool of size `nWorkers` and calls the function `f` for each
// element in the `items` map
func WorkerPoolFromMap[K comparable, V any](ctx context.Context, items map[K]V, nWorkers int, f func(context.Context, K, V)) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(nWorkers)

	// Spawn a fixed pool of workers
	ch := make(chan KeyValue[K, V])
	for ii := 0; ii < nWorkers; ii++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case v, ok := <-ch:
					if !ok {
						// If the channel is closed, exit
						return
					}
					f(ctx, v.Key, v.Value)
				case <-ctx.Done():
					return
				}
			}

		}()
	}

	// Push items to the workers
	for k, v := range items {
		ch <- KeyValue[K, V]{Key: k, Value: v}
	}
	// Stop all the workers
	close(ch)

	// Wait for the workers to finish processing their current items
	wg.Wait()
}

// WorkerPoolFromChan starts a worker pool of size `nWorkers` and calls the function `f` for each
// element in the `items` channel
func WorkerPoolFromChan[T any](ctx context.Context, items <-chan T, nWorkers int, f func(ctx context.Context, item T)) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(nWorkers)

	// Spawn a fixed pool of workers
	for ii := 0; ii < nWorkers; ii++ {
		go func(ii int) {
			defer wg.Done()
			for {
				select {
				case v, ok := <-items:
					if !ok {
						return
					}
					f(ctx, v)
				case <-ctx.Done():
					return
				}
			}

		}(ii)
	}

	// Wait for the workers to finish processing their current items
	wg.Wait()
}

// WorkerPoolFromSlice starts a worker pool of size `nWorkers` and calls the function `f` for each
// element in the `items` slice
func WorkerPoolFromSlice[T any](ctx context.Context, items []T, nWorkers int, f func(ctx context.Context, item T)) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(nWorkers)

	// Spawn a fixed pool of workers
	ch := make(chan T)
	for ii := 0; ii < nWorkers; ii++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case v, ok := <-ch:
					if !ok {
						// If the channel is closed, exit
						return
					}
					f(ctx, v)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Push items to the workers
	for ii := 0; ii < len(items); ii++ {
		ch <- items[ii]
	}
	close(ch)

	// Wait for the workers to finish processing their current items
	wg.Wait()
}
