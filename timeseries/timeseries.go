package timeseries

import (
	"github.com/lobocv/simpleflow/v0/simpletime"
	"time"
)

// TimeTransformation converts one time to another. It is primarily used to set the granularity of a time series
type TimeTransformation func(time.Time) time.Time

// Entry is a specific entry in a time series, it contains a value at a specific time
type Entry[V any] struct {
	Time  time.Time
	Value V
}

// TimeSeries keeps track of values for a series of times. Values are not expected to be contiguous
// as they are stored in an underlying map. Time granularity can be enforced by providing a TimeTransformation
// function. This function can be used to round all values to their closest minute, hour or day.
type TimeSeries[V any] struct {
	values map[time.Time]V
	tf     TimeTransformation
}

// NewTimeSeries creates a new TimeSeries
// All times are stored in UTC regardless of their TimeTransformation. This is required in order to look up values in
// the map where time.Time is the key
func NewTimeSeries[V any](m map[time.Time]V, timeGranularity TimeTransformation) TimeSeries[V] {
	values := map[time.Time]V{}

	for date, metrics := range m {
		values[timeGranularity(date)] = metrics
	}
	return TimeSeries[V]{
		values: values,
		tf:     timeGranularity,
	}
}

// Length return the length of the time series
func Length[V any](m TimeSeries[V]) int {
	return len(m.values)
}

// Merge merges the second time series values into the first time series
// Values are overwritten by the provided time series if they already exist
func Merge[V any](ts, other TimeSeries[V]) {
	iter, done := Iterate(other)
	defer done()
	for e := range iter {
		Set(ts, e.Time, e.Value)
	}
}

// Get gets a value at a specific date in the time series
func Get[V any](m TimeSeries[V], date time.Time) (V, bool) {
	v, ok := m.values[m.tf(date)]
	return v, ok
}

// Set sets a value at a specific date in the time series
func Set[V any](m TimeSeries[V], date time.Time, v V) {
	m.values[m.tf(date)] = v
}

// Unset removes a value at a specific date in the time series
func Unset[V any](m TimeSeries[V], date time.Time) {
	delete(m.values, m.tf(date))
}

// Iterate returns a read-only channel and a close function to iterate through all time series values
// Note that this does not iterate in chronological order.
func Iterate[V any](ts TimeSeries[V]) (<-chan Entry[V], func()) {
	iter := make(chan Entry[V])
	done := make(chan struct{}, 2) // nolint

	// define a closer function which we can pass back in case the caller wants to end iteration early
	closer := func() {
		done <- struct{}{}
	}

	// start a go routine that iterates the map and sends the values on a channel
	go func() {
		defer close(iter)
		for date, m := range ts.values {
			dm := Entry[V]{
				Time:  date,
				Value: m,
			}
			select {
			case iter <- dm:
			case <-done:
				return
			}
		}
	}()
	return iter, closer
}

// OrderedIterate returns a read-only channel and a close function to iterate through time series values in
// the given range and step. For time entries that do not exist in the time series, nothing is returned.
func OrderedIterate[V any](ts TimeSeries[V], start, end time.Time, step time.Duration) (<-chan Entry[V], func()) {
	iter := make(chan Entry[V])
	done := make(chan struct{}, 2) // nolint

	// define a closer function which we can pass back in case the caller wants to end iteration early
	closer := func() {
		done <- struct{}{}
	}

	// start a go routine that iterates the map and sends the values on a channel
	go func() {
		defer close(iter)
		it := simpletime.NewIterator(start, end, step)
		for it.Next() {
			date := ts.tf(it.Current())
			v, ok := Get(ts, date)
			if !ok {
				continue
			}
			dm := Entry[V]{
				Time:  date,
				Value: v,
			}
			select {
			case iter <- dm:
			case <-done:
				return
			}
		}
	}()
	return iter, closer
}
