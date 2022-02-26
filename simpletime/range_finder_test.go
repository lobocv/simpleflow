package simpletime

import "time"

func (s TestSuite) TestRangeFinder() {
	var rf RangeFinder
	expectedStart := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2022, 1, 11, 0, 0, 0, 0, time.UTC)

	// Include the next 10 days
	for ii := 0; ii <= 10; ii++ {
		rf.Include(expectedStart.Add(time.Duration(24*ii) * time.Hour))
	}
	r := rf.Range()
	s.Equal(expectedStart, r.Start)
	s.Equal(expectedStart, rf.Earliest())

	s.Equal(expectedEnd, r.End)
	s.Equal(expectedEnd, rf.Latest())
}
