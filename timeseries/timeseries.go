package timeseries

import (
	"time"
)

// TimeTransformation converts one time to another. It is mainly used to set the granularity of a time series
type TimeTransformation func(time.Time) time.Time

// Entry is a specific entry in a time series, it contains a value at a specific time
type Entry[V any] struct {
	Time  time.Time
	Value V
}

// TimeSeries keeps track of values for a series of times
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

// Join merges the second time series values into the first time series
// Values are overwritten by the provided time series if they already exist
func Join[V any](ts, other TimeSeries[V]) {
	iter, done := IterateTimeSeries(other)
	defer done()
	for e := range iter {
		Set(ts, e.Time, e.Value)
	}
}

// Set sets a value at a specific date in the time series
func Set[V any](m TimeSeries[V], date time.Time, v V) {
	m.values[m.tf(date)] = v
}

// Unset removes a value at a specific date in the time series
func Unset[V any](m TimeSeries[V], date time.Time) {
	delete(m.values, m.tf(date))
}

// IterateTimeSeries returns a read-only channel and a close function to iterate time series values
// Note that this does not iterate in chronological order
func IterateTimeSeries[V any](ts TimeSeries[V]) (<-chan Entry[V], func()) {
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
