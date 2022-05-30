package simpletime

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
}

func TestTime(t *testing.T) {
	s := new(TestSuite)
	suite.Run(t, s)
}

// daysAgo takes the current date (in UTC) and subtracts the input argument number of days
func daysAgo(nDays int) time.Time {
	return time.Now().UTC().AddDate(0, 0, -nDays)
}

func (s TestSuite) TestEarliestAndLatest() {
	times := []time.Time{
		daysAgo(10),
		daysAgo(5),
		daysAgo(3),
		daysAgo(0),
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
