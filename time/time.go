package simpletime

import (
	"time"
)

// Date is a short form for creating a time.Time at the start of a day with UTC time
func Date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// Earliest returns the earliest time from the given times. It returns zero time if no arguments are passed
func Earliest(times ...time.Time) time.Time {
	var earliest time.Time
	for _, t := range times {
		if earliest.IsZero() || t.Before(earliest) {
			earliest = t
		}
	}
	return earliest
}

// Latest returns the latest time from the given times. It returns zero time if no arguments are passed
func Latest(times ...time.Time) time.Time {
	var latest time.Time
	for _, t := range times {
		if latest.IsZero() || t.After(latest) {
			latest = t
		}
	}
	return latest
}
