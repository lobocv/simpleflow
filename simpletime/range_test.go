package simpletime

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func Date(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

type TestSuite struct {
	suite.Suite
}

func TestTime(t *testing.T) {
	s := new(TestSuite)
	suite.Run(t, s)
}

func (s TestSuite) TestCombineRanges() {

	r1 := Range{Date(2021, 02, 02), Date(2021, 02, 10)}
	r2 := Range{Date(2021, 01, 20), Date(2021, 01, 30)}
	r3 := Range{Date(2021, 02, 1), Date(2021, 02, 21)}

	combined := r1.Combine(r2, r3)
	expected := Range{Date(2021, 01, 20), Date(2021, 02, 21)}
	s.Equal(expected, combined)
}

func (s TestSuite) TestCombineRangesEmptyRange() {

	r1 := Range{}
	r2 := Range{Date(2021, 02, 02), Date(2021, 02, 10)}

	expected := Range{Date(2021, 02, 02), Date(2021, 02, 10)}

	combined := r1.Combine(r2)
	s.Equal(expected, combined)

	combined = r2.Combine(r1)
	s.Equal(expected, combined)
}

func (s TestSuite) TestRangeIntersectsAndContains() {
	r0 := Range{Date(2021, 01, 01), Date(2021, 01, 5)}
	r1 := Range{Date(2021, 01, 05), Date(2021, 01, 10)}
	s.True(r0.Overlaps(r1, true))
	s.True(r1.Overlaps(r0, true))
	s.False(r0.Overlaps(r1, false))
	s.False(r1.Overlaps(r0, false))
	s.False(r0.ContainsRange(r1, false))
	s.False(r1.ContainsRange(r0, false))

	r2 := Range{Date(2021, 01, 10), Date(2021, 01, 15)}
	s.True(r1.Overlaps(r2, true))
	s.True(r2.Overlaps(r1, true))
	s.False(r1.Overlaps(r2, false))
	s.False(r2.Overlaps(r1, false))
	s.False(r1.ContainsRange(r2, false))
	s.False(r2.ContainsRange(r1, false))

	s.True(r1.ContainsRange(r1, true))
	s.False(r1.ContainsRange(r1, false))
}
