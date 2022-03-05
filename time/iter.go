package simpletime

import "time"

// Iterator iterates between a start and end time at a fixed time.Duration step
type Iterator struct {
	current time.Time
	end     time.Time
	step    time.Duration
}

// NewIterator creates a new iterate that iterates between `start` and `end` by `step` intervals.
func NewIterator(start time.Time, end time.Time, step time.Duration) *Iterator {
	return &Iterator{
		// subtract one step so that the first call to Next() returns r.Earliest
		current: start.Add(-step),
		end:     end,
		step:    step,
	}
}

// Next moves the iterator to the next position and returns whether the iterator is exhausted
func (it *Iterator) Next() bool {
	it.current = it.current.Add(it.step)
	hasNext := !it.current.After(it.end) // !d1.After(d2) is the same as d1 <= d2
	return hasNext
}

// Current returns the current value of the iterator
func (it *Iterator) Current() time.Time {
	return it.current
}
