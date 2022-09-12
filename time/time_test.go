package simpletime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestTime(t *testing.T) {
	s := new(TestSuite)
	suite.Run(t, s)
}

func (s TestSuite) TestEarliestAndLatest() {
	times := []time.Time{
		DaysAgo(10),
		DaysAgo(5),
		DaysAgo(3),
		DaysAgo(0),
	}

	earliest := Earliest(times...)
	s.Equal(times[0], earliest)

	latest := Latest(times...)
	s.Equal(times[3], latest)

	// Check single input edge case
	s.Equal(times[0], Latest(times[0]))
	s.Equal(times[0], Earliest(times[0]))

	// Check edge case with no inputs
	s.Equal(time.Time{}, Latest())
	s.Equal(time.Time{}, Earliest())
}

func (s TestSuite) TestAgoAndFromNow() {
	testCases := []struct {
		fn            func(n int) time.Time
		expectedDiff  time.Duration
		expectedDelta time.Duration
	}{
		{fn: SecondsAgo[int], expectedDiff: -2 * time.Second},
		{fn: MinutesAgo[int], expectedDiff: -2 * time.Minute},
		{fn: HoursAgo[int], expectedDiff: -2 * time.Hour},
		{fn: DaysAgo[int], expectedDiff: -48 * time.Hour},
		{fn: MonthsAgo[int], expectedDiff: time.Now().AddDate(0, -2, 0).Sub(time.Now())},
		{fn: YearsAgo[int], expectedDiff: time.Now().AddDate(-2, 0, 0).Sub(time.Now())},
		{fn: SecondsFromNow[int], expectedDiff: 2 * time.Second},
		{fn: MinutesFromNow[int], expectedDiff: 2 * time.Minute},
		{fn: HoursFromNow[int], expectedDiff: 2 * time.Hour},
		{fn: DaysFromNow[int], expectedDiff: 48 * time.Hour},
		{fn: MonthsFromNow[int], expectedDiff: time.Now().AddDate(0, 2, 0).Sub(time.Now())},
		{fn: YearsFromNow[int], expectedDiff: time.Now().AddDate(2, 0, 0).Sub(time.Now())},
	}

	for ii, tc := range testCases {
		// Get the time from the function being tested
		t := tc.fn(2)
		// determine the expected time
		expect := time.Now().Add(tc.expectedDiff)
		// Check that they don't differ by very much (should be only the difference between the two calls to time.Now())
		deltaExpected := t.Sub(expect)
		if deltaExpected < 0 {
			deltaExpected *= -1
		}
		s.Less(deltaExpected, time.Millisecond, "Test case failed:", ii+1)
	}
}

func (s TestSuite) TestAbsDelta() {
	expectedDelta := 15 * time.Second
	t1 := time.Date(2022, 02, 02, 14, 15, 16, 0, time.UTC)
	t2 := t1.Add(expectedDelta)

	// Delta should be the same regardless of the order t1 and t2 are passed in. Value should be positive.
	delta := AbsDelta(t1, t2)
	s.Equal(delta, expectedDelta)
	delta = AbsDelta(t2, t1)
	s.Equal(delta, expectedDelta)
}
