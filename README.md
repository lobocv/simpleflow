[![Go Reference](https://pkg.go.dev/badge/github.com/lobocv/simplerr.svg)](https://pkg.go.dev/github.com/lobocv/simpleflow)
[![Github tag](https://badgen.net/github/tag/lobocv/simpleflow)](https://github.com/lobocv/simpleflow/tags)
![Go version](https://img.shields.io/github/go-mod/go-version/lobocv/simpleflow)
![Build Status](https://github.com/lobocv/simpleflow/actions/workflows/build.yaml/badge.svg)
[![GoReportCard](https://goreportcard.com/badge/github.com/lobocv/simpleflow)](https://goreportcard.com/report/github.com/lobocv/simpleflow)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>
[<img src="https://img.shields.io/github/license/lobocv/simpleflow">](https://img.shields.io/github/license/lobocv/simpleflow)


# SimpleFlow

SimpleFlow is a a collection of generic functions and patterns that help building common workflows.
Please see the tests for examples on how to use these functions.

## Why should I use Simpleflow?

- A single library for common workflows so you do not have to reinvent the wheel, maintain your own library or
  copy-paste code.
- Simple and easy to use API.
- Detailed documentation and examples. 
- Worker pools are simple, worker pools with error handling are not.
- 100% test coverage

# Installation
```
go get -u github.com/lobocv/simpleflow
```

## Table of Contents

1. [Channels](https://github.com/lobocv/simpleflow#channels)
2. [Work Pools](https://github.com/lobocv/simpleflow#worker-pools)
   1. [Example](https://github.com/lobocv/simpleflow#workerpoolfromslice-example)
   2. [Canceling a running worker pool](https://github.com/lobocv/simpleflow#canceling-a-running-worker-pool)
3. [Fan-Out and Fan-In](https://github.com/lobocv/simpleflow#fan-out-and-fan-in)
4. [Round Robin](https://github.com/lobocv/simpleflow#round-robin)
5. [Batching](https://github.com/lobocv/simpleflow#batching)
6. [Incremental Batching](https://github.com/lobocv/simpleflow#incremental-batching)
7. [Filtering](https://github.com/lobocv/simpleflow#filtering)
8. [Extracting](https://github.com/lobocv/simpleflow#extracting)
9. [Segmenting](https://github.com/lobocv/simpleflow#segmenting)
10. [Deduplication](https://github.com/lobocv/simpleflow#deduplication)
11. [Counter](https://github.com/lobocv/simpleflow#counter)
12. [Time](https://github.com/lobocv/simpleflow#time)
13. [Time Series](https://github.com/lobocv/simpleflow#timeseries)

## Channels

Some common but tedious operations on channels are done by the channel functions:

Example:

```go
items := make(chan int, 3)
// push 1, 2, 3 onto the channel
LoadChannel(items, 1, 2, 3)
// Close the channel so ChannelToSlice doesn't block.
close(items) 
out := ChannelToSlice(items)
// out == []int{1, 2, 3}
```

## Worker Pools

Worker pools provide a way to spin up a finite set of go routines to process items in a collection.

- `WorkerPoolFromSlice` - Starts a fixed pool of workers that process elements in the `slice`
- `WorkerPoolFromMap` - Starts a fixed pool of workers that process key-value pairs in the `map`
- `WorkerPoolFromChan` - Starts a fixed pool of workers that process values read from a `channel`

These functions block until all workers finish processing.

### WorkerPoolFromSlice example

```go
ctx := context.Background()
items := []int{0, 1, 2, 3, 4, 5}
out := NewSyncMap(map[int]int{})
nWorkers := 2
f := func(_ context.Context, v int) error {
    out.Set(v, v*v)
    return nil
}
errors := WorkerPoolFromSlice(ctx, items, nWorkers, f)
// errors == []error{}
// out == map[int]int{0: 0, 1: 1, 2: 4, 3: 9, 4: 16, 5: 25}
```

### Canceling a running worker pool 

```go
// Create a cancel-able context
ctx, cancel := context.WithCancel(context.Background())

items := []int{0, 1, 2, 3, 4, 5}
out := NewSyncMap(map[int]int{}) // threadsafe map used in tests
nWorkers := 2

f := func(_ context.Context, v int) error {
    // Cancel as soon as we hit v > 2
    if v > 2 {
        cancel()
        return nil
    }
    out.Set(v, v*v)
    return nil
}
WorkerPoolFromSlice(ctx, items, nWorkers, f)
// errors == []error{}
// out == map[int]int{0: 0, 1: 1, 2: 4}
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
fanInResults := ChannelToSlice(fanInSink)
// fanInResults == []int{1, 2, 3, 1, 2, 3}
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
sink1 := make(chan int, N)
sink2 := make(chan int, N)
RoundRobin(source, fanoutSink1, sink2)
CloseManyWriters(fanoutSink1, sink2)

sink1Data := ChannelToSlice(sink1)
// sink1Data == []int{1, 3, 5}
sink2Data := ChannelToSlice(sink2)
// sink2Data == []int{2, 4}
```

## Batching

`BatchMap`, `BatchSlice` and `BatchChan` provide ways to break `maps`, `slices` and `channels` into smaller
components of at most `N` size.

Example:

```go
items := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
batches := BatchSlice(items, 2)
// batches == [][]int{{0, 1}, {2, 3}, {4, 5}, {6, 7}, {8, 9}}
```

```go
items := map[int]int{0: 0, 1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
batches := BatchMap(items, 3)
// batches == []map[int]int{ {0: 0, 3: 3, 4: 4}, {1: 1, 2: 2, 5: 5} }

```

## Incremental Batching

Batching can also be done incrementally by using `IncrementalBatchSlice` and `IncrementalBatchMap` functions.
These functions are meant to be called repeatedly, adding elements until a full batch can be processed, at which time,
the batch is returned.

Example:

```go
batchSize := 3
var items, batch []int
items, batch = IncrementalBatchSlice(items, batchSize, 1)
// items == []int{1}, batch == nil

items, batch = IncrementalBatchSlice(items, batchSize, 2)
// items == []int{1, 2}, batch == nil

items, batch = IncrementalBatchSlice(items, batchSize, 3)
// Batch size reached
// items == []int{}, batch == []int{1, 2, 3}

items, batch = IncrementalBatchSlice(items, batchSize, 4)
// items == []int{4}, batch == nil
```

## Filtering

Filtering operations allows you to remove elements from slices or maps. Filtering can done either in-place with
`FilterSliceInplace`, `FilterMapInPlace` or by creating a copy `FilterSlice`, `FilterMap`. 

The following example filters the positive numbers from the slice of integers:

```go
getPositive := func(t int) bool {
    return t > 0
}

out := FilterSlice([]int{5, -2, 3, 1, 0, -3, -5, -6}, getPositive)
// out == []int{5, 3, 1}
```

## Extracting

Extraction operations allow you to extract one or more elements from a slice. Extraction functions accept a function which takes in each element of the slice and returns the
element to extract (potentially different than the input type) and a boolean for whether the element should be extracted.

The following example extracts the `Name` field from a slice of `Object`.

```go
type Object struct {
    Name string
}

in := []Object{
    {Name: "John"},
    {Name: "Paul"},
    {Name: "George"},
    {Name: "Ringo"},
    {Name: "Bob"},
}
var names []string

fn := func(t Object) (string, bool) {
    if t.Name == "Bob" {
        return "", false
    }
    return t.Name, true
}

names = ExtractToSlice(in, fn, names)
// names == []string{"John", "Paul", "George", "Ringo"}
```

`ExtractFirst()` can be used to extract the first element in the slice. The following example extracts the first value
larger than 4.

```go
values := []int{4, 1, 5, 7}
fn := func(v int) bool {
    return v > 4
}

v, found := ExtractFirst(values, fn)
// v == 5, found == true
```

## Segmenting

`SegmentSlice`, `SegmentMap` and `SegmentChan` allow you to split a `slice` or `map` into sub-slices or maps based on the provided
segmentation function:

### Segmenting a slice into even and odd values
```go
items := []int{0, 1, 2, 3, 4, 5}

segments := SegmentSlice(items, func(v int) int {
    if v % 2 == 0 {
        return "even"
	}
        return "odd"
})
// segments == map[string][]int{"even": {0, 2, 4}, "odd": {1, 3, 5}}
```

## Deduplication
A series of values can be deduplicated using the `Deduplicator{}`. It can either accept the entire slice:

```go
deduped := Deduplicate([]int{1, 1, 2, 2, 3, 3})
// deduped == []int{1, 2, 3}
```
or iteratively deduplicate for situations where you want fine control with a `for` loop.
```go
dd := NewDeduplicator[int]()
values := []int{1, 1, 2, 3, 3}
deduped := []int{}

for _, v := range values {
    seen := dd.Seen(v) 
    // seen == true for index 1 and 4
    isNew := dd.Add(v) 
    // isNew == true for index 0, 2 and 3
    if isNew {
        deduped = append(deduped, v)	
    }
}
```

Complex objects can also be deduplicated using the `ObjectDeduplicator{}`, which requires providing a function that
creates unique IDs for the provided objects being deduplicated. This is useful for situations where the values being 
deduplicated are not comparable (ie, have a slice field) or if you want more fine control over just what constitutes a 
duplicate.

```go
// Object is a complex structure that cannot be used with a regular Deduplicator as it contains 
// a slice field, and thus is not `comparable`.
type Object struct {
    slice   []int
    pointer *int
    value   string
}

// Create a deduplicator that deduplicates Object's by their "value" field.
dd := NewObjectDeduplicator[Object](func(v Object) string {
        return v.value
    })
```


## Counter
The `Counter{}` and `ObjectCounter{}` can be used to count the number of occurrences
of values. Much like the `Deduplicator{}`, the `Counter{}` works well for simple types.

```go
counter := NewCounter[int]()
values := []int{1, 1, 2, 3, 3, 3, 3}

// Add the values to the counter, values can also be added individually with counter.Add()
currentCount := counter.AddMany(values) 

numberOfOnes := counter.Count(1) // returns 2
numberOfTwos := counter.Count(2) // returns 1
numberOfThrees := counter.Count(3) // returns 4
```

Complex objects can also be counted using the `ObjectCounter{}`, which requires providing a function that
creates buckets for the provided objects being deduplicated. This is useful for situations where the values being
counted are not comparable (ie, have a slice field) or if you want more fine control over the bucketing logic (ie
bucket objects by a certain field value).

```go
// Object is a complex structure that cannot be used with a regular Counter as it contains 
// a slice field, and thus is not `comparable`.
type Object struct {
    slice   []int
    pointer *int
    value   string
}

// Create a counter that counts Object's bucketed by their "value" field.
counter := NewObjectCounter[Object](func(v Object) string {
        return v.value
    })
```


## Time

The `simeplflow/time` package provides functions that assist with working with the standard library `time` package
and `time.Time` objects. The package contains functions to define, compare and iterate time ranges.

## Timeseries

The `simpleflow/timeseries` packages contains a generic `TimeSeries` object that allows you
to manipulate timestamped data. `TimeSeries` store unordered time series data in an underlying 
`map[time.Time]V`. Each `TimeSeries` is configured with a `TimeTransformation` which applies to each
`time.Time` key when accessed. This makes storing time series data with a particular time granularity
easy. For example, with a `TimeTransformation` that truncates to the day, any 
`time.Time` object in the given day will access the same key.

Example:

```go
// TF is a TimeTransformation that truncates the time to the start of the day
func TF(t time.Time) time.Time {
    return t.UTC().Truncate(24 * time.Hour)
}

// Day is a function to create a Time object on a given day offset from Jan 1st 2022 by the `i`th day
func Day(i int) time.Time {
    return time.Date(2022, 01, i, 0, 0, 0, 0, time.UTC)
}

func main() {
    data := map[time.Time]int{
            Day(0): 0, // Jan 1st 2022
            Day(1): 1, // Jan 2nd 2022
            Day(2): 2, // Jan 3rd 2022
        }
    ts := timeseries.NewTimeSeries(data, TF)
	
	// Get the value on Jan 2th at 4am and at 5 am
	// The values for `a` and `b` are both == 1 because the hour is irrelevant
	// when accessing data using the TF() time transform
	a := ts.Get(time.Date(2022, 01, 2, 4, 0, 0, 0, time.UTC))
	b := ts.Get(time.Date(2022, 01, 2, 5, 0, 0, 0, time.UTC))
	// a == b == 1 
}
```
