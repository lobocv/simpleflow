package simpleflow

// ChannelIntoSlice reads elements from the channel and returns appends them to the `out` slice.
// This operation will block until the channel is closed
func ChannelIntoSlice[T any](ch chan T, out []T) []T {
	for v := range ch {
		out = append(out, v)
	}
	return out
}

// ChannelToSlice reads elements from the channel and returns them as a slice.
// This operation will block until the channel is closed
func ChannelToSlice[T any](ch chan T) (out []T) {
	return ChannelIntoSlice(ch, out)
}

// LoadChannel puts all elements from `items` onto the channel `ch`
// This operation will block if not all items fit within the channel buffer or
// if there is not simultaneously another go routine reading from the channel.
func LoadChannel[T any](ch chan<- T, items ...T) {
	for ii := 0; ii < len(items); ii++ {
		ch <- items[ii]
	}
	return
}
