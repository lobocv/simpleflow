package simpletime

import "time"

// Range is a simple tuple of earliest and latest time
type Range struct {
	Start time.Time
	End   time.Time
}

// Duration returns the duration between start and end of the range
func (r Range) Duration() time.Duration {
	return r.End.Sub(r.Start)
}

// Combine returns a Range which spans all the given ranges. The final range will be the earliest and latest value
// of all provided ranges.
func (r Range) Combine(others ...Range) Range {
	combined := r
	for _, other := range append(others, r) {
		// We need to watch out for zero value Earliest times
		if (other.Start.Before(combined.Start) || combined.Start.IsZero()) && !other.Start.IsZero() {
			combined.Start = other.Start
		}
		if other.End.After(combined.End) {
			combined.End = other.End
		}
	}
	return combined
}

// Contains checks whether the given time is within the range. If `inclusive` is set true, the edges of the range are
// included in the contains logic.
func (r Range) Contains(t time.Time, inclusive bool) bool {
	left := r.Start.Before(t)
	if inclusive {
		left = left || r.Start.Equal(t)
	}
	right := r.End.After(t)
	if inclusive {
		right = right || r.End.Equal(t)
	}

	return right && left
}

/*
ContainsRange returns whether the receiver range contains the other range. The `inclusive` flag determines if the
extents on the ranges are considered to be contained if they are equal in value.

Example of `r` contains `other`:
	r:		|-----------------------------|
	other:		 |-----------|
*/
func (r Range) ContainsRange(other Range, inclusive bool) bool {
	left := r.Start.Before(other.Start)
	if inclusive {
		left = left || r.Start.Equal(other.Start)
	}
	right := r.End.After(other.End)
	if inclusive {
		right = right || r.End.Equal(other.End)
	}
	return left && right
}

/*
Overlaps returns whether the receiver range overlaps the other range. The `inclusive` flag determines if the
extents on the ranges are considered to be overlapped if they are equal in value.

Example of `r` contains `other`:
	r:		|---------|
	other:        |-----------|
*/
func (r Range) Overlaps(other Range, inclusive bool) bool {
	return r.overlapsLeft(other, inclusive) || r.overlapsRight(other, inclusive) || r.ContainsRange(other, inclusive)
}

/*
	r:		|-------------|
	other:		 |-------------|
*/
func (r Range) overlapsLeft(other Range, inclusive bool) bool {
	overlap := r.Start.Before(other.Start)
	if inclusive {
		overlap = overlap && (r.End.After(other.Start) || r.End.Equal(other.Start))
	} else {
		overlap = overlap && r.End.After(other.Start)
	}
	return overlap
}

/*
	r:				|-------------|
	other:	 |-------------|
*/
func (r Range) overlapsRight(other Range, inclusive bool) bool {
	overlap := r.End.After(other.End)
	if inclusive {
		overlap = overlap && (r.Start.Before(other.End) || r.Start.Equal(other.End))
	} else {
		overlap = overlap && r.Start.Before(other.End)
	}
	return overlap
}

// Iterate returns an iterator that uses the extents of the Range.
func (r Range) Iterate(step time.Duration) *Iterator {
	return NewIterator(r.Start, r.End, step)
}

// IterateDays returns an iterator function that returns the nth day (n >= 1) and a boolean that signals that the iterator is done
func (r Range) IterateDays(n int) *Iterator {
	return r.IterateHours(n * 24)
}

func (r Range) IterateHours(n int) *Iterator {
	return r.Iterate(time.Duration(n) * time.Hour)
}

func (r Range) IterateMinutes(n int) *Iterator {
	return r.Iterate(time.Duration(n) * time.Minute)
}

func (r Range) IterateSeconds(n int) *Iterator {
	return r.Iterate(time.Duration(n) * time.Second)
}
