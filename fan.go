package concurgo

import (
	"sync"
)

// FanOut reads from the `from` channel and publishes the data across all `to` channels
// It closes each `to` channel once all messages are drained from the `from` channels.
func FanOut[T any](from <-chan T, to ...chan<- T) {
	for v := range from {
		for _, ch := range to {
			ch <- v
		}
	}

	for _, ch := range to {
		close(ch)
	}
}

// FanIn reads from each `from` channel and writes to the `to` channel
// It closes the `to` channel once all messages are drained from the `from` channels.
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
	close(to)
}
