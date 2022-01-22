package simpleflow

// RoundRobin reads from the `from` channel and distributes the values in a round-robin fashion to the `to` channels.
func RoundRobin[T any](from <-chan T, to ...chan<- T) {
	if len(to) == 0 {
		return
	}

	var count int
	for v := range from {
		to[count%len(to)] <- v
		count++
	}

	for _, ch := range to {
		close(ch)
	}
}
