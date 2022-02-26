package simpletime

import (
	"time"
)

// RangeFinder implements a greedy algorithm for finding the  earliest and latest values of a series of time.Time objects
type RangeFinder struct {
	earliest time.Time
	latest   time.Time
}

// Include adds the time.Time object to the RangeFinder, potentially increasing the limits.
func (l *RangeFinder) Include(t time.Time) {

	if t.Before(l.earliest) || l.earliest.IsZero() {
		l.earliest = t
	}
	if t.After(l.latest) || l.latest.IsZero() {
		l.latest = t
	}
}

// Range returns the range of the RangeFinder
func (l *RangeFinder) Range() Range {
	return Range{Start: l.earliest, End: l.latest}
}

// Earliest returns the earliest and latest limits
func (l *RangeFinder) Earliest() time.Time {
	return l.earliest
}

// Latest returns the earliest and latest limits
func (l *RangeFinder) Latest() time.Time {
	return l.latest
}
