package concurgo

func FanOut[T any](from <-chan T, to ...chan<- T) {
	v := <-from
	for ii := 0; ii < len(to); ii++ {
		to[ii] <- v
	}
}
