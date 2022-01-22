package concurgo

import (
	"sync"
)

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
