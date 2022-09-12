package simpletime

import (
	"time"
)

// Signed is a constraint that permits any signed integer type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

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

// AbsDelta returns the absolute difference in time between the two provided times
func AbsDelta(t1, t2 time.Time) time.Duration {
	d := t1.Sub(t2)
	if d < 0 {
		return -d
	}
	return d
}

// YearsAgo returns a time.Time n years back in time from now
func YearsAgo[T Signed](n T) time.Time {
	return time.Now().AddDate(int(-n), 0, 0)
}

// YearsFromNow returns a time.Time n years in the future from now
func YearsFromNow[T Signed](n T) time.Time {
	return time.Now().AddDate(int(n), 0, 0)
}

// MonthsAgo returns a time.Time n months back in time from now
func MonthsAgo[T Signed](n T) time.Time {
	return time.Now().AddDate(0, int(-n), 0)
}

// MonthsFromNow returns a time.Time n months in the future from now
func MonthsFromNow[T Signed](n T) time.Time {
	return time.Now().AddDate(0, int(n), 0)
}

// DaysAgo returns a time.Time n days back in time from now
func DaysAgo[T Signed](n T) time.Time {
	return time.Now().AddDate(0, 0, int(-n))
}

// DaysFromNow returns a time.Time n days in the future from now
func DaysFromNow[T Signed](n T) time.Time {
	return time.Now().AddDate(0, 0, int(n))
}

// HoursAgo returns a time.Time n hours back in time from now
func HoursAgo[T Signed](n T) time.Time {
	return time.Now().Add(time.Duration(-n) * time.Hour)
}

// HoursFromNow returns a time.Time n hours in the future from now
func HoursFromNow[T Signed](n T) time.Time {
	return time.Now().Add(time.Duration(n) * time.Hour)
}

// MinutesAgo returns a time.Time n minutes back in time from now
func MinutesAgo[T Signed](n T) time.Time {
	return time.Now().Add(time.Duration(-n) * time.Minute)
}

// MinutesFromNow returns a time.Time n minutes in the future from now
func MinutesFromNow[T Signed](n T) time.Time {
	return time.Now().Add(time.Duration(n) * time.Minute)
}

// SecondsAgo returns a time.Time n seconds back in time from now
func SecondsAgo[T Signed](n T) time.Time {
	return time.Now().Add(time.Duration(-n) * time.Second)
}

// SecondsFromNow returns a time.Time n seconds in the future from now
func SecondsFromNow[T Signed](n T) time.Time {
	return time.Now().Add(time.Duration(n) * time.Second)
}
