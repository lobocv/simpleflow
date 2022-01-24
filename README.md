[![GoReportCard](https://goreportcard.com/badge/github.com/lobocv/simpleflow)](https://goreportcard.com/report/github.com/lobocv/simpleflow)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-95%25-brightgreen.svg?longCache=true&style=flat)</a>



# SimpleFlow

SimpleFlow is a a collection of generic functions and patterns that help building common concurrent workflows.
Please see the tests for examples on how to use these functions.

## Channels

Some common but tedious operations on channels are done by the channel functions:

Example:

```go
    items := make(chan int, 3)
    LoadChannel(items, 1, 2, 3)  // pushes 1, 2, 3 onto the channel
	close(items) // Close the channel so ChannelToSlice doesn't block.
    out := ChannelToSlice(items) // out ---> []int{1, 2, 3}
```

## Worker Pools

Worker pools provide a way to spin up a finite set of go routines to process items in a collection.

- `WorkerPoolFromSlice` - Starts a fixed pool of workers that process elements in the `slice`
- `WorkerPoolFromMap` - Starts a fixed pool of workers that process key-value pairs in the `map`
- `WorkerPoolFromChan` - Starts a fixed pool of workers that process values read from a `channel`

These functions block until all workers finish processing.

Example:

```go

	ctx := context.Background()
	items := []int{0, 1, 2, 3, 4, 5}
	out := make([]int, 0, len(items))
	nWorkers := 2
	
	f := func(ctx context.Context, v int) {
		out = append(out, v*v)
	}
	WorkerPoolFromSlice(ctx, items, nWorkers, f)
	
	// out ---> []int{4, 1, 0, 9, 25, 16}
	 

```

## Fan-Out and Fan-In

`FanOut` and `FanIn` provide means of fanning-in and fanning-out channel to other channels. 

Example:

```go
// Generate some data on a channel (source for fan out)
N := 3
source := make(chan int, N)
data := []int{1, 2, 3}
for _, v := range data {
    source <- v
}
close(source)

// Fan out to two channels. Each will get a copy of the data
fanoutSink1 := make(chan int, N)
fanoutSink2 := make(chan int, N)
FanOutAndClose(source, fanoutSink1, fanoutSink2)

// Fan them back in to a single channel. We should get the original source data with two copies of each item
fanInSink := make(chan int, 2*N)
FanInAndClose(fanInSink, fanoutSink1, fanoutSink2)
faninResults := ChannelToSlice(fanInSink) // faninResults ---> []int{1, 2, 3, 1, 2, 3}
```

## Round Robin

`RoundRobin` distributes values from a channel over other channels in a round-robin fashion

Example:

```go
// Generate some data on a channel
N := 5
source := make(chan int, N)
data := []int{1, 2, 3, 4, 5}
for _, v := range data {
    source <- v
}
close(source)

// Round robin the data into two channels, each should have half the data
fanoutSink1 := make(chan int, N)
fanoutSink2 := make(chan int, N)
RoundRobin(source, fanoutSink1, fanoutSink2)
CloseManyWriters(fanoutSink1, fanoutSink2)

fanout1Data := ChannelToSlice(fanoutSink1) // fanout1Data ---> []int{1, 3, 5}
fanout2Data := ChannelToSlice(fanoutSink2) // fanout2Data ---> []int{2, 4}
```

## Batching

`BatchMap`, `BatchSlice` and `BatchChan` provide ways to break `maps`, `slices` and `channels` into smaller
components of at most `N` size.

Example:

```go
items := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
batches := BatchSlice(items, 2)
// batches ---> [][]int{{0, 1}, {2, 3}, {4, 5}, {6, 7}, {8, 9}}
```

```go
items := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
batches := BatchMap(items, 2)
// batches ---> []map[int]int{ {0: 0, 3: 3, 5: 5}, {1: 1, 2: 2, 4: 4}

```
