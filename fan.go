package simpleflow

import "sync"

// FanOut reads from the `from` channel and publishes the data across all `to` channels
func FanOut[T any](from <-chan T, to ...chan<- T) {
	for v := range from {
		for _, ch := range to {
			ch <- v
		}
	}
}

// FanOutAndClose reads from the `from` channel and publishes the data across all `to` channels
// It closes each `to` channel once all messages are drained from the `from` channels.
func FanOutAndClose[T any](from <-chan T, to ...chan<- T) {
	FanOut(from, to...)
	CloseManyWriters(to...)
}

// FanIn reads from each `from` channel and writes to the `to` channel
func FanIn[T any](to chan<- T, from ...<-chan T) {
	var wg sync.WaitGroup
	wg.Add(len(from))
	for _, ch := range from {

		go func(ch <-chan T) {
			defer wg.Done()
			for v := range ch {
				to <- v
			}
		}(ch)
	}

	wg.Wait()
}

// FanInAndClose reads from each `from` channel and writes to the `to` channel
// It closes the `to` channel once all messages are drained from the `from` channels.
func FanInAndClose[T any](to chan<- T, from ...<-chan T) {
	FanIn(to, from...)
	close(to)
}
