package simpleflow

import (
	"context"
	"sync"
)

// Job is a function that the slice or channel worker pool executes
type Job[T any] func(ctx context.Context, item T) error

// JobKV is a function that the map worker pool executes
type JobKV[K comparable, V any] func(context.Context, K, V) error

// KeyValue is a tuple of key, value
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// WorkerPoolFromMap starts a worker pool of size `nWorkers` and calls the function `f` for each
// element in the `items` map
func WorkerPoolFromMap[K comparable, V any](ctx context.Context, items map[K]V, nWorkers int, f JobKV[K, V]) []error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var errChan = make(chan error)
	var wg sync.WaitGroup
	wg.Add(nWorkers + 1)

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
					err := f(ctx, v.Key, v.Value)
					if err != nil {
						errChan <- err
					}
				case <-ctx.Done():
					return
				}
			}

		}()
	}

	// Start a separate go routine that is pulling errors from the workers and appending them to a slice
	// This is required so that the error channel does not block. The alternative is to create a buffered
	// error channel with size len(items)
	var errors []error
	errWaitGroup := sync.WaitGroup{}
	errWaitGroup.Add(1)
	go func() {
		defer errWaitGroup.Done()
		errors = ChannelToSlice(errChan)
	}()

	// Push items to the workers in a separate go routine in case the channel gets blocked by a canceled worker
	go func() {
		defer wg.Done()
		for k, v := range items {
			select {
			case ch <- KeyValue[K, V]{Key: k, Value: v}:
			case <-ctx.Done():
				return
			}
		}
		close(ch)
	}()

	// Wait for the workers to finish processing their current items
	wg.Wait()
	// Close the error channel and wait for the errors to drain into the slice
	close(errChan)
	errWaitGroup.Wait()

	return errors
}

// WorkerPoolFromChan starts a worker pool of size `nWorkers` and calls the function `f` for each
// element in the `items` channel
func WorkerPoolFromChan[T any](ctx context.Context, items <-chan T, nWorkers int, f Job[T]) []error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var errChan = make(chan error)
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
					err := f(ctx, v)
					if err != nil {
						errChan <- err
					}
				case <-ctx.Done():
					return
				}
			}

		}(ii)
	}

	// Start a separate go routine that is pulling errors from the workers and appending them to a slice
	// This is required so that the error channel does not block. The alternative is to create a buffered
	// error channel with size len(items)
	var errors []error
	errWaitGroup := sync.WaitGroup{}
	errWaitGroup.Add(1)
	go func() {
		defer errWaitGroup.Done()
		errors = ChannelToSlice(errChan)
	}()

	// Wait for the workers to finish processing their current items
	wg.Wait()

	// Close the error channel and wait for the errors to drain into the slice
	close(errChan)
	errWaitGroup.Wait()

	return errors
}

// WorkerPoolFromSlice starts a worker pool of size `nWorkers` and calls the function `f` for each
// element in the `items` slice. It returns an array of errors from the jobs.
func WorkerPoolFromSlice[T any](ctx context.Context, items []T, nWorkers int, f Job[T]) []error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var errChan = make(chan error)
	var wg sync.WaitGroup
	wg.Add(nWorkers + 1)

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
					err := f(ctx, v)
					if err != nil {
						errChan <- err
					}

				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Start a separate go routine that is pulling errors from the workers and appending them to a slice
	// This is required so that the error channel does not block. The alternative is to create a buffered
	// error channel with size len(items)
	var errors []error
	errWaitGroup := sync.WaitGroup{}
	errWaitGroup.Add(1)
	go func() {
		defer errWaitGroup.Done()
		errors = ChannelToSlice(errChan)
	}()

	// Push items to the workers in a separate go routine in case the channel gets blocked by a canceled worker
	go func() {
		defer wg.Done()
		for ii := 0; ii < len(items); ii++ {
			select {
			case ch <- items[ii]:
			case <-ctx.Done():
				return
			}
		}
		close(ch)
	}()

	// Wait for the workers to finish processing their current items
	wg.Wait()
	// Close the error channel and wait for the errors to drain into the slie
	close(errChan)
	errWaitGroup.Wait()

	return errors
}
