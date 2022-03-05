package timeseries

import (
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
	"time"
)

func TF(t time.Time) time.Time {
	return t.UTC().Truncate(24 * time.Hour)
}

func Day(i int) time.Time {
	return time.Date(2022, 01, i, rand.Intn(23), rand.Intn(59), 7, 8, time.UTC)

}

type TestSuite struct {
	suite.Suite
}

func TestTimeSeries(t *testing.T) {
	s := new(TestSuite)
	suite.Run(t, s)
}

func (s *TestSuite) TestSetAndUnset() {

	ts := NewTimeSeries(map[time.Time]int{}, TF)
	s.Zero(Length(ts))
	Set(ts, Day(0), 0)
	Set(ts, Day(1), 1)
	Set(ts, Day(2), 2)
	s.Equal(3, Length(ts))
	s.Equal(map[time.Time]int{TF(Day(0)): 0, TF(Day(1)): 1, TF(Day(2)): 2}, ts.values)

	// Get individual dates
	for ii := 0; ii < 3; ii++ {
		got, ok := Get(ts, Day(ii))
		s.True(ok)
		s.Equal(got, ii)
	}

	// Get dates not in the timeseries
	for ii := 3; ii < 10; ii++ {
		got, ok := Get(ts, Day(ii))
		s.False(ok)
		s.Equal(got, 0)
	}

	Unset(ts, Day(2))
	s.Equal(2, Length(ts))
	s.Equal(map[time.Time]int{TF(Day(0)): 0, TF(Day(1)): 1}, ts.values)
}

func (s *TestSuite) TestIterate() {
	ts := NewTimeSeries(map[time.Time]int{
		Day(0): 0,
		Day(1): 1,
		Day(2): 2,
	},
		TF)
	s.Equal(3, Length(ts))

	ch, done := Iterate(ts)
	defer done()

	var got []Entry[int]
	for v := range ch {
		got = append(got, v)
	}
	expected := []Entry[int]{
		{Time: TF(Day(0)), Value: 0},
		{Time: TF(Day(1)), Value: 1},
		{Time: TF(Day(2)), Value: 2},
	}
	s.ElementsMatch(got, expected)
}

func (s *TestSuite) TestOrderedIterate() {
	ts := NewTimeSeries(map[time.Time]int{
		Day(0): 0,
		Day(1): 1,
		Day(2): 2,
		// Missing Day(3)
		Day(4): 4,
		Day(5): 5,
	},
		TF)

	start := time.Date(2022, 01, 1, 0, 0, 0, 0, time.UTC)
	stop := time.Date(2022, 01, 6, 5, 0, 0, 0, time.UTC)
	step := 24 * time.Hour
	ch, done := OrderedIterate(ts, start, stop, step)
	defer done()

	var got []Entry[int]
	for v := range ch {
		if v.Time.After(Day(3)) {
			done()
			break
		}
		got = append(got, v)
	}
	expected := []Entry[int]{
		{Time: TF(Day(1)), Value: 1},
		{Time: TF(Day(2)), Value: 2},
	}
	s.ElementsMatch(got, expected)
}

func (s *TestSuite) TestIterateExitEarly() {
	ts := NewTimeSeries(map[time.Time]int{
		Day(0): 0,
		Day(1): 1,
		Day(2): 2,
		Day(3): 3,
		Day(4): 3,
	},
		TF)
	s.Equal(5, Length(ts))

	ch, done := Iterate(ts)
	defer done()

	var got []Entry[int]

	count := 0
	for v := range ch {
		// Exit early
		if count == 3 {
			done()
			break
		}
		got = append(got, v)
		count++
	}
	s.Len(got, 3)
}

func (s *TestSuite) TestJoin() {
	ts1 := NewTimeSeries(map[time.Time]int{Day(0): 0, Day(1): 1, Day(2): 2}, TF)
	ts2 := NewTimeSeries(map[time.Time]int{Day(2): 20, Day(3): 30, Day(4): 40}, TF)

	Merge(ts1, ts2)
	s.Equal(5, Length(ts1))
	s.Equal(map[time.Time]int{
		TF(Day(0)): 0,
		TF(Day(1)): 1,
		TF(Day(2)): 20,
		TF(Day(3)): 30,
		TF(Day(4)): 40,
	}, ts1.values)
}
