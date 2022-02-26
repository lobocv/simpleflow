package simpletime

import "time"

func (s TestSuite) TestIterateDays() {

	r := Range{Start: Date(2021, 01, 01), End: Date(2021, 01, 15)}
	iter := r.IterateDays(1)
	var count int
	for iter.Next() {
		date := iter.Current()

		expected := Date(2021, 01, 01).AddDate(0, 0, count)
		s.Equal(expected, date)
		count++
	}
	s.Equal(count, 15) // 1st thru 15th inclusive
}

func (s TestSuite) TestIterateMinutes() {

	r := Range{
		Start: time.Date(2021, 01, 01, 1, 1, 2, 0, time.UTC),
		End:   time.Date(2021, 01, 01, 1, 16, 2, 0, time.UTC),
	}
	iter := r.IterateMinutes(2)
	var count int
	for iter.Next() {
		date := iter.Current()

		expected := r.Start.Add(time.Duration(2*count) * time.Minute)
		s.Equal(expected, date)
		count++
	}
	s.Equal(count, 8) // 1st thru 15th inclusive
}

func (s TestSuite) TestIterateSeconds() {

	r := Range{
		Start: time.Date(2021, 01, 01, 1, 1, 1, 0, time.UTC),
		End:   time.Date(2021, 01, 01, 1, 1, 15, 0, time.UTC),
	}
	iter := r.IterateSeconds(3)
	var count int
	for iter.Next() {
		date := iter.Current()

		expected := r.Start.Add(time.Duration(3*count) * time.Second)
		s.Equal(expected, date)
		count++
	}
	s.Equal(count, 5) // 1st thru 15th inclusive
}
