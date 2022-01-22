package concurgo

func ChannelToSlice[T any](ch chan T, out []T) []T {
	for v := range ch {
		out = append(out, v)
	}
	return out
}

func DumpChannel[T any](ch chan T) (out []T) {
	return ChannelToSlice(ch, out)
}

func LoadChannel[T any](ch chan<- T, items ...T) {
	for ii := 0; ii < len(items); ii++ {
		ch <- items[ii]
	}
	return
}
