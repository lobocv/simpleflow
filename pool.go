package simpleflow

import (
	"context"
	"sync"
)

// WorkerPoolFromMap starts a worker pool of size `nWorkers` and calls the function `f` for each
// element in the `items` map
func WorkerPoolFromMap[K comparable, V any](ctx context.Context, items map[K]V, nWorkers int, f func(K, V)) {
	sem := make(chan struct{}, nWorkers)
	for ii := 0; ii < nWorkers; ii++ {
		sem <- struct{}{}
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for k, v := range items {
		wg.Add(1)
		go func(k K, v V) {
			select {
			case <-sem:
				f(k, v)
				sem <- struct{}{}
				wg.Done()
			case <-ctx.Done():
				return
			}

		}(k, v)
	}

	wg.Wait()
}

// WorkerPoolFromChan starts a worker pool of size `nWorkers` and calls the function `f` for each
// element in the `items` channel
func WorkerPoolFromChan[T any](ctx context.Context, items <-chan T, nWorkers int, f func(job T)) {
	sem := make(chan struct{}, nWorkers)
	for ii := 0; ii < nWorkers; ii++ {
		sem <- struct{}{}
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for item := range items {
		wg.Add(1)
		go func(v T) {
			select {
			case <-sem:
				f(v)
				sem <- struct{}{}
				wg.Done()
			case <-ctx.Done():
				return
			}

		}(item)
	}

	wg.Wait()
}

// WorkerPoolFromSlice starts a worker pool of size `nWorkers` and calls the function `f` for each
// element in the `items` slice
func WorkerPoolFromSlice[T any](ctx context.Context, items []T, nWorkers int, f func(v T)) {
	sem := make(chan struct{}, nWorkers)
	for ii := 0; ii < nWorkers; ii++ {
		sem <- struct{}{}
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for ii := 0; ii < len(items); ii++ {

		wg.Add(1)
		go func(v T) {
			select {
			case <-sem:
				f(v)
				sem <- struct{}{}
				wg.Done()
			case <-ctx.Done():
				return
			}

		}(items[ii])
	}

	wg.Wait()
}
